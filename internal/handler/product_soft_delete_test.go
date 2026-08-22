package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sall-lah/store_product/internal/config"
	"github.com/Sall-lah/store_product/internal/repository"
	"github.com/Sall-lah/store_product/internal/service"
)

// TestProductSoftDeleteRoutes verifies routing and authorization for soft delete endpoints.
func TestProductSoftDeleteRoutes(t *testing.T) {
	cfg := &config.Config{Port: "8080", RateLimitPublicRPM: 100, RateLimitAdminRPM: 50}
	svc := service.NewProductService(&repository.ProductRepository{}, &repository.VariantRepository{}, nil)
	h := NewProductHandler(svc)
	router := SetupRouter(cfg, h, nil)

	t.Run("Admin DELETE /api/v1/admin/products/{id} requires admin role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/products/prod-123", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected 403 Forbidden for unauthenticated request, got %d", rec.Code)
		}
	})

	t.Run("Admin DELETE /api/v1/admin/products/{id}/variants/{variantId} requires admin role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/products/prod-123/variants/var-456", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected 403 Forbidden for unauthenticated request, got %d", rec.Code)
		}
	})

	t.Run("Authenticated Admin DELETE /api/v1/admin/products/{id} passes auth middleware", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/products/prod-123", nil)
		req.Header.Set("X-User-Role", "admin")
		req.Header.Set("X-User-Id", "admin-1")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		// Passes auth gate (may fail with 500/nil client in unit test without db, but should not be 403)
		if rec.Code == http.StatusForbidden {
			t.Errorf("expected request to pass admin auth gate, got 403 Forbidden")
		}
	})

	t.Run("Authenticated Admin DELETE variant passes auth middleware", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/products/prod-123/variants/var-456", nil)
		req.Header.Set("X-User-Role", "admin")
		req.Header.Set("X-User-Id", "admin-1")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code == http.StatusForbidden {
			t.Errorf("expected request to pass admin auth gate, got 403 Forbidden")
		}
	})
}
