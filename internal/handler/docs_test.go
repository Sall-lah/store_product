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

func TestDocsEndpoints(t *testing.T) {
	cfg := &config.Config{Port: "8080", RateLimitPublicRPM: 100, RateLimitAdminRPM: 50}
	svc := service.NewProductService(&repository.ProductRepository{}, &repository.VariantRepository{}, nil)
	h := NewProductHandler(svc)
	router := SetupRouter(cfg, h, nil, nil)

	t.Run("GET /openapi.json returns valid OpenAPI 3.1 schema", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}

		contentType := rec.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			t.Errorf("expected application/json content-type, got %s", contentType)
		}

		var spec map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&spec); err != nil {
			t.Fatalf("failed to decode openapi.json response: %v", err)
		}

		if spec["openapi"] != "3.1.0" {
			t.Errorf("expected openapi 3.1.0, got %v", spec["openapi"])
		}

		paths, ok := spec["paths"].(map[string]interface{})
		if !ok || paths["/api/v1/products"] == nil {
			t.Errorf("expected spec to contain /api/v1/products path")
		}
	})

	t.Run("GET /openapi.yaml returns YAML specification", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}

		contentType := rec.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/yaml") {
			t.Errorf("expected application/yaml content-type, got %s", contentType)
		}

		body := rec.Body.String()
		if !strings.Contains(body, "openapi: 3.1.0") {
			t.Errorf("expected YAML to include 'openapi: 3.1.0'")
		}
	})

	t.Run("GET /swagger renders Swagger UI HTML", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}

		contentType := rec.Header().Get("Content-Type")
		if !strings.Contains(contentType, "text/html") {
			t.Errorf("expected text/html content-type, got %s", contentType)
		}

		body := rec.Body.String()
		if !strings.Contains(body, "swagger-ui") {
			t.Errorf("expected Swagger UI body to contain swagger-ui DOM container")
		}
	})

	t.Run("GET /docs renders Scalar UI HTML", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}

		contentType := rec.Header().Get("Content-Type")
		if !strings.Contains(contentType, "text/html") {
			t.Errorf("expected text/html content-type, got %s", contentType)
		}

		body := rec.Body.String()
		if !strings.Contains(body, "@scalar/api-reference") {
			t.Errorf("expected Scalar UI body to contain @scalar/api-reference")
		}
	})
}
