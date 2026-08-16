package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireAdmin(t *testing.T) {
	adminHandler := RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r.Context())
		if !ok || user.Role != "admin" {
			http.Error(w, "missing context user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("admin_ok"))
	}))

	t.Run("Authorized Admin Request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", nil)
		req.Header.Set("X-User-Role", "admin")
		req.Header.Set("X-User-Id", "usr_admin_1")
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
