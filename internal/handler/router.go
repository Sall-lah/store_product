package handler

import (
	"github.com/Sall-lah/store_product/internal/cache"
	"github.com/Sall-lah/store_product/internal/config"
	"github.com/Sall-lah/store_product/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// SetupRouter initializes the HTTP routing hierarchy, middlewares, and endpoint boundaries.
// Separating router construction from main.go enables simple end-to-end integration testing.
func SetupRouter(cfg *config.Config, h *ProductHandler, imgHandler *ImageHandler, cacheClient *cache.Client) *chi.Mux {
	r := chi.NewRouter()

	// Global cross-cutting middlewares
	r.Use(middleware.SetupCORS())
	r.Use(middleware.Logger)
	r.Use(middleware.Recovery)

	// Liveness probe
	r.Get("/health", h.HealthCheck)

	// OpenAPI 3.1 & Interactive Documentation Endpoints
	docsHandler := NewDocsHandler()
	r.Get("/openapi.json", docsHandler.ServeOpenAPIJSON)
	r.Get("/openapi.yaml", docsHandler.ServeOpenAPIYAML)

	// Scalar UI & sub-route specs for relative path resolution
	r.Route("/docs", func(sub chi.Router) {
		sub.Get("/openapi.json", docsHandler.ServeOpenAPIJSON)
		sub.Get("/openapi.yaml", docsHandler.ServeOpenAPIYAML)
		sub.Get("/", docsHandler.ServeScalarUI)
		sub.Get("/*", docsHandler.ServeScalarUI)
	})
	r.Get("/docs", docsHandler.ServeScalarUI)

	// Swagger UI & sub-route specs for relative path resolution
	r.Route("/swagger", func(sub chi.Router) {
		sub.Get("/openapi.json", docsHandler.ServeOpenAPIJSON)
		sub.Get("/openapi.yaml", docsHandler.ServeOpenAPIYAML)
		sub.Get("/", docsHandler.ServeSwaggerUI)
		sub.Get("/*", docsHandler.ServeSwaggerUI)
	})
	r.Get("/swagger", docsHandler.ServeSwaggerUI)

	// API v1 Namespace
	r.Route("/api/v1", func(api chi.Router) {

		// Public Catalog Endpoints (Protected by sliding window rate limiter, customer-facing read-only)
		api.Route("/products", func(public chi.Router) {
			public.Use(middleware.RateLimiter(cacheClient, cfg.RateLimitPublicRPM, "public_catalog"))

			public.Get("/", h.ListProducts)
			public.Get("/{id}", h.GetProductByID)
			public.Get("/slug/{slug}", h.GetProductBySlug)
		})

		// Admin Backoffice Endpoints (Enforces API Gateway X-User-Role: admin & strict admin rate limit)
		api.Route("/admin/products", func(admin chi.Router) {
			admin.Use(middleware.RequireAdmin)
			admin.Use(middleware.RateLimiter(cacheClient, cfg.RateLimitAdminRPM, "admin_writes"))

			// Product Backoffice Management
			admin.Get("/", h.AdminListProducts)
			admin.Get("/{id}", h.AdminGetProductByID)
			admin.Post("/", h.CreateProduct)
			admin.Put("/{id}", h.UpdateProduct)
			admin.Delete("/{id}", h.DeleteProduct)

			// Variant Backoffice Management
			admin.Post("/{id}/variants", h.CreateVariant)
			admin.Put("/{id}/variants/{variantId}", h.UpdateVariant)
			admin.Delete("/{id}/variants/{variantId}", h.DeleteVariant)

			// Media & Image Gallery Backoffice Management (Cloudflare R2 Direct Upload & Gallery)
			if imgHandler != nil {
				admin.Post("/{id}/images/presign", imgHandler.PresignImage)
				admin.Post("/{id}/images", imgHandler.CreateImage)
				admin.Get("/{id}/images", imgHandler.ListImages)
				admin.Put("/{id}/images/{imageId}", imgHandler.UpdateImage)
				admin.Delete("/{id}/images/{imageId}", imgHandler.DeleteImage)
			}
		})
	})

	return r
}
