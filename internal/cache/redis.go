package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Sall-lah/store_product/internal/config"
	"github.com/redis/go-redis/v9"
)

// Client wraps the Go Redis client providing domain-specific caching helpers
// and resilient error-handling to prevent Redis outages from taking down the API.
type Client struct {
	rdb       *redis.Client
	available bool
}

// NewClient initializes a connection to Redis using the application configuration.
// It executes a short-timeout ping to verify reachability without stalling startup.
func NewClient(cfg *config.Config) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr(),
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		PoolSize:     20,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client := &Client{rdb: rdb, available: false}
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("[WARN] Redis is not reachable at %s: %v. Running in degraded cache-fallback mode.", cfg.RedisAddr(), err)
	} else {
		log.Printf("[INFO] Successfully connected to Redis at %s", cfg.RedisAddr())
		client.available = true
	}

	return client
}

// IsAvailable returns true if the Redis client was successfully initialized.
// Handlers and services check this to bypass cache calls when Redis is unavailable.
func (c *Client) IsAvailable() bool {
	return c != nil && c.available && c.rdb != nil
}

// Underlying returns the raw redis.Client for low-level operations such as rate limiter scripts.
func (c *Client) Underlying() *redis.Client {
	if c == nil {
		return nil
	}
	return c.rdb
}

// GetJSON retrieves a cached value by key and unmarshals it into the target destination.
// Returns false if key was not found or if Redis is unavailable.
func (c *Client) GetJSON(ctx context.Context, key string, dest interface{}) bool {
	if !c.IsAvailable() {
		return false
	}

	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return false
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		log.Printf("[WARN] Failed to unmarshal cached key %s: %v", key, err)
		return false
	}

	return true
}

// SetJSON marshals an object into JSON and stores it in Redis with the given TTL.
// Fails silently with a log message to ensure caching failures never block write transactions.
func (c *Client) SetJSON(ctx context.Context, key string, data interface{}, ttl time.Duration) {
	if !c.IsAvailable() {
		return
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("[WARN] Failed to marshal object for cache key %s: %v", key, err)
		return
	}

	if err := c.rdb.Set(ctx, key, bytes, ttl).Err(); err != nil {
		log.Printf("[WARN] Failed to set Redis cache key %s: %v", key, err)
	}
}

// Del removes one or more keys from Redis.
func (c *Client) Del(ctx context.Context, keys ...string) {
	if !c.IsAvailable() || len(keys) == 0 {
		return
	}

	if err := c.rdb.Del(ctx, keys...).Err(); err != nil {
		log.Printf("[WARN] Failed to delete Redis keys %v: %v", keys, err)
	}
}

// DelPattern invalidates all keys matching a specific pattern using SCAN.
// SCAN is used instead of KEYS to avoid blocking the Redis server in high-throughput environments.
func (c *Client) DelPattern(ctx context.Context, pattern string) {
	if !c.IsAvailable() {
		return
	}

	var cursor uint64
	for {
		keys, nextCursor, err := c.rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			log.Printf("[WARN] Error scanning pattern %s: %v", pattern, err)
			break
		}

		if len(keys) > 0 {
			c.Del(ctx, keys...)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
}

// ProductDetailKey formats a cache key for single product lookups by ID.
func ProductDetailKey(id string) string {
	return fmt.Sprintf("product:detail:id:%s", id)
}

// ProductSlugKey formats a cache key for single product lookups by Slug.
func ProductSlugKey(slug string) string {
	return fmt.Sprintf("product:detail:slug:%s", slug)
}

// ProductListKey formats a cache key for catalog query results using a hashed query string.
func ProductListKey(hash string) string {
	return fmt.Sprintf("product:list:%s", hash)
}
