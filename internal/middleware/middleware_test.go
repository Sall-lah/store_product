package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Sall-lah/store_product/internal/cache"
	"github.com/Sall-lah/store_product/internal/config"
	"github.com/joho/godotenv"
)

func TestRequireAdmin(t *testing.T) {
	adminHandler := RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r.Context())
		if !ok || !strings.EqualFold(strings.TrimSpace(user.Role), "admin") {
			http.Error(w, "missing context user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("admin_ok"))
	}))

	t.Run("Authorized Admin Request (lowercase)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", nil)
		req.Header.Set("X-User-Role", "admin")
		req.Header.Set("X-User-Id", "usr_admin_1")
		rec := httptest.NewRecorder()

		adminHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("Authorized Admin Request (uppercase from gateway/auth)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", nil)
		req.Header.Set("X-User-Role", "ADMIN")
		req.Header.Set("X-User-Id", "usr_admin_2")
		rec := httptest.NewRecorder()

		adminHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("Authorized Admin Request (mixed case with whitespace)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", nil)
		req.Header.Set("X-User-Role", "  Admin  ")
		req.Header.Set("X-User-Id", "usr_admin_3")
		rec := httptest.NewRecorder()

		adminHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("Unauthorized Request Missing Role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", nil)
		rec := httptest.NewRecorder()

		adminHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected status %d, got %d", http.StatusForbidden, rec.Code)
		}

		var body map[string]string
		_ = json.NewDecoder(rec.Body).Decode(&body)
		if body["error"] != "forbidden" {
			t.Errorf("expected error code forbidden, got %s", body["error"])
		}
	})

	t.Run("Forbidden Customer Role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", nil)
		req.Header.Set("X-User-Role", "customer")
		rec := httptest.NewRecorder()

		adminHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestRecoveryMiddleware(t *testing.T) {
	panicHandler := Recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("unexpected critical failure")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test-panic", nil)
	rec := httptest.NewRecorder()

	panicHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500 on panic recovery, got %d", rec.Code)
	}

	var body map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&body)
	if body["error"] != "internal_server_error" {
		t.Errorf("expected error internal_server_error, got %s", body["error"])
	}
}

func TestRateLimiterFallbackWhenNil(t *testing.T) {
	limiter := RateLimiter(nil, 60, "test_tier")
	nextCalled := false
	handler := limiter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !nextCalled {
		t.Errorf("expected handler to execute when rate limiter fails open")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestRateLimiterWithLiveRedis(t *testing.T) {
	_ = godotenv.Load("../../.env", ".env")
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	cacheClient := cache.NewClient(cfg)
	if !cacheClient.IsAvailable() {
		t.Skip("Redis is not available, skipping live rate limiter test")
	}

	// Limit to 2 requests per minute for this test tier
	limiter := RateLimiter(cacheClient, 2, "test_limiter_tier")
	handler := limiter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	req1.Header.Set("X-User-Id", "test_user_rate_limit")
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Fatalf("first request expected 200, got %d", rec1.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	req2.Header.Set("X-User-Id", "test_user_rate_limit")
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("second request expected 200, got %d", rec2.Code)
	}

	req3 := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	req3.Header.Set("X-User-Id", "test_user_rate_limit")
	rec3 := httptest.NewRecorder()
	handler.ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusTooManyRequests {
		t.Fatalf("third request expected 429 Too Many Requests, got %d", rec3.Code)
	}

	// Clean up key
	cacheClient.Del(req1.Context(), "ratelimit:test_limiter_tier:user:test_user_rate_limit")
}
