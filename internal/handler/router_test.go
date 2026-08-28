package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Sall-lah/store_product/internal/config"
	"github.com/Sall-lah/store_product/internal/repository"
	"github.com/Sall-lah/store_product/internal/service"
)

func TestHealthCheck(t *testing.T) {
	cfg := &config.Config{Port: "8080", RateLimitPublicRPM: 100, RateLimitAdminRPM: 50}
	svc := service.NewProductService(&repository.ProductRepository{}, &repository.VariantRepository{}, nil)
	h := NewProductHandler(svc)
	router := SetupRouter(cfg, h, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var res map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode json response: %v", err)
	}

	if res["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got '%s'", res["status"])
	}
}

func TestAdminRoutesProtection(t *testing.T) {
	cfg := &config.Config{Port: "8080", RateLimitPublicRPM: 100, RateLimitAdminRPM: 50}
	svc := service.NewProductService(&repository.ProductRepository{}, &repository.VariantRepository{}, nil)
	h := NewProductHandler(svc)
	router := SetupRouter(cfg, h, nil, nil)

	t.Run("Reject Unauthenticated POST /api/v1/admin/products", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/products", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected status 403 Forbidden without admin header, got %d", rec.Code)
		}
	})

	t.Run("Reject Unauthenticated GET /api/v1/admin/products", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/products", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected status 403 Forbidden without admin header, got %d", rec.Code)
		}
	})

	t.Run("Reject Non-Admin User PUT /api/v1/admin/products/123", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/products/123", nil)
		req.Header.Set("X-User-Role", "seller")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected status 403 Forbidden for role 'seller', got %d", rec.Code)
		}
	})

	t.Run("Accept Authenticated Admin with Uppercase Role POST /api/v1/admin/products", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/products", strings.NewReader(`{}`))
		req.Header.Set("X-User-Role", "ADMIN")
		req.Header.Set("X-User-Id", "usr_admin_1")
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code == http.StatusForbidden {
			t.Errorf("expected request to pass admin middleware for uppercase ADMIN role, got 403 Forbidden")
		}
	})

	t.Run("Reject Unauthenticated DELETE /api/v1/admin/products/123", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/products/123", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected status 403 Forbidden without admin header, got %d", rec.Code)
		}
	})

	t.Run("Reject Unauthenticated Variant POST /api/v1/admin/products/123/variants", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/products/123/variants", strings.NewReader(`{}`))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected status 403 Forbidden without admin header, got %d", rec.Code)
		}
	})

	t.Run("Public Routes Are Accessible Without Admin Headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		// Public endpoint should not return 403 Forbidden
		if rec.Code == http.StatusForbidden {
			t.Errorf("expected public catalog to be accessible without admin headers, got %d", rec.Code)
		}
	})
}
