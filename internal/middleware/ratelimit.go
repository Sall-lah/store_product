package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Sall-lah/store_product/internal/cache"
	"github.com/redis/go-redis/v9"
)

// RateLimiter returns HTTP middleware implementing a sliding-window rate limiter
// backed by Redis. Tiered limits protect distinct endpoints from resource exhaustion.
func RateLimiter(cacheClient *cache.Client, requestsPerMinute int, tierName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If Redis is offline or client is nil, fail open to maintain service availability
			// rather than denying legitimate traffic due to cache infrastructure issues.
			if cacheClient == nil || !cacheClient.IsAvailable() {
				next.ServeHTTP(w, r)
				return
			}

			identifier := extractClientIdentifier(r)
			key := fmt.Sprintf("ratelimit:%s:%s", tierName, identifier)
			now := time.Now()
			nowMs := now.UnixNano() / int64(time.Millisecond)
			windowStartMs := nowMs - int64(60*time.Second/time.Millisecond)

			rdb := cacheClient.Underlying()
			ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
			defer cancel()

			pipe := rdb.Pipeline()
			// 1. Remove expired timestamps outside the rolling 60-second window
			pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStartMs, 10))
			// 2. Count requests currently within the active window
			cardCmd := pipe.ZCard(ctx, key)
			// 3. Set key TTL so inactive user records naturally expire from memory
			pipe.Expire(ctx, key, 65*time.Second)

			_, err := pipe.Exec(ctx)
			if err != nil && err != redis.Nil {
				log.Printf("[WARN] Rate limiter pipeline error: %v. Allowing request.", err)
				next.ServeHTTP(w, r)
				return
			}

			currentCount := int(cardCmd.Val())
			resetTimestamp := now.Add(60 * time.Second).Unix()
			remaining := requestsPerMinute - currentCount

			if remaining < 0 {
				remaining = 0
			}

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(requestsPerMinute))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTimestamp, 10))

			if currentCount >= requestsPerMinute {
				w.Header().Set("Retry-After", "60")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error":   "too_many_requests",
					"message": fmt.Sprintf("Rate limit exceeded for %s. Please retry in 60 seconds.", tierName),
				})
				return
			}

			// Record current request timestamp in the sliding window
			if err := rdb.ZAdd(ctx, key, redis.Z{
				Score:  float64(nowMs),
				Member: fmt.Sprintf("%d", now.UnixNano()),
			}).Err(); err != nil {
				log.Printf("[WARN] Failed to record request in rate limiter ZSET: %v", err)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// extractClientIdentifier resolves a stable client identifier for rate limiting.
// Prioritizes authenticated user ID header from the API Gateway, falling back to IP.
func extractClientIdentifier(r *http.Request) string {
	if userID := r.Header.Get("X-User-Id"); userID != "" {
		return "user:" + userID
	}

	// Extract real client IP through reverse proxies / gateway
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return "ip:" + ip
			}
		}
	}

	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return "ip:" + strings.TrimSpace(xrip)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return "ip:" + host
	}

	return "ip:" + r.RemoteAddr
}
