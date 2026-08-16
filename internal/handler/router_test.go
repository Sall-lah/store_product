package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sall-lah/store_product/internal/config"
	"github.com/Sall-lah/store_product/internal/repository"
	"github.com/Sall-lah/store_product/internal/service"
)

func TestHealthCheck(t *testing.T) {
	cfg := &config.Config{Port: "8080", RateLimitPublicRPM: 100, RateLimitAdminRPM: 50}
	svc := service.NewProductService(&repository.ProductRepository{}, &repository.VariantRepository{}, nil)
	h := NewProductHandler(svc)
	router := SetupRouter(cfg, h, nil)

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
	router := SetupRouter(cfg, h, nil)

	t.Run("Reject Unauthenticated POST /api/v1/products", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected status 403 Forbidden without admin header, got %d", rec.Code)
		}
	})

	t.Run("Reject Non-Admin User PUT /api/v1/products/123", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/products/123", nil)
		req.Header.Set("X-User-Role", "seller")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected status 403 Forbidden for role 'seller', got %d", rec.Code)
		}
	})
}
