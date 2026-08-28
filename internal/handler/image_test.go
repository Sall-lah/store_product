package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sall-lah/store_product/internal/config"
	"github.com/Sall-lah/store_product/internal/repository"
	"github.com/Sall-lah/store_product/internal/service"
	"github.com/Sall-lah/store_product/internal/storage"
)

func TestImageHandlerRoutesAndAuth(t *testing.T) {
	cfg := &config.Config{
		Port:               "8080",
		RateLimitPublicRPM: 100,
		RateLimitAdminRPM:  50,
		R2BucketName:       "test-bucket",
		R2PublicBaseURL:    "https://cdn.mystore.com",
	}

	storageClient, _ := storage.NewR2StorageClient(t.Context(), cfg)
	svc := service.NewProductService(&repository.ProductRepository{}, &repository.VariantRepository{}, nil)
	imgSvc := service.NewImageService(&repository.ImageRepository{}, &repository.ProductRepository{}, storageClient, nil)

	h := NewProductHandler(svc)
	imgHandler := NewImageHandler(imgSvc)
	router := SetupRouter(cfg, h, imgHandler, nil)

	t.Run("POST /api/v1/admin/products/{id}/images/presign rejects unauthenticated requests", func(t *testing.T) {
		body, _ := json.Marshal(repository.PresignImageInput{
			FileName:    "photo.webp",
			ContentType: "image/webp",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/products/prod-1/images/presign", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected 403 Forbidden without admin headers, got %d", rec.Code)
		}
	})

	t.Run("POST /api/v1/admin/products/{id}/images rejects unauthenticated requests", func(t *testing.T) {
		body, _ := json.Marshal(repository.CreateProductImageInput{
			URL:   "https://cdn.mystore.com/img.webp",
			R2Key: "products/1/images/img.webp",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/products/prod-1/images", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected 403 Forbidden without admin headers, got %d", rec.Code)
		}
	})

	t.Run("GET /api/v1/admin/products/{id}/images rejects unauthenticated requests", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/products/prod-1/images", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected 403 Forbidden without admin headers, got %d", rec.Code)
		}
	})

	t.Run("PUT /api/v1/admin/products/{id}/images/{imageId} rejects unauthenticated requests", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/products/prod-1/images/img-1", bytes.NewReader([]byte("{}")))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected 403 Forbidden without admin headers, got %d", rec.Code)
		}
	})

	t.Run("DELETE /api/v1/admin/products/{id}/images/{imageId} rejects unauthenticated requests", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/products/prod-1/images/img-1", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected 403 Forbidden without admin headers, got %d", rec.Code)
		}
	})

	t.Run("POST presign with admin headers passes auth gate", func(t *testing.T) {
		body, _ := json.Marshal(repository.PresignImageInput{
			FileName:    "photo.webp",
			ContentType: "image/webp",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/products/prod-1/images/presign", bytes.NewReader(body))
		req.Header.Set("X-User-Role", "ADMIN")
		req.Header.Set("X-User-Id", "usr_admin_1")
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		// Without live DB, will fail past auth with 500 or 404 (product not found), but must pass the 403 auth gate
		if rec.Code == http.StatusForbidden {
			t.Errorf("expected request with ADMIN role to pass auth gate, got 403")
		}
	})
}
