package handler

import (
	"github.com/Sall-lah/store_product/internal/cache"
	"github.com/Sall-lah/store_product/internal/config"
	"github.com/Sall-lah/store_product/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// SetupRouter initializes the HTTP routing hierarchy, middlewares, and endpoint boundaries.
// Separating router construction from main.go enables simple end-to-end integration testing.
func SetupRouter(cfg *config.Config, h *ProductHandler, cacheClient *cache.Client) *chi.Mux {
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
	r.Get("/swagger", docsHandler.ServeSwaggerUI)
	r.Get("/swagger/*", docsHandler.ServeSwaggerUI)
	r.Get("/docs", docsHandler.ServeScalarUI)
	r.Get("/docs/*", docsHandler.ServeScalarUI)

	// API v1 Namespace
	r.Route("/api/v1", func(api chi.Router) {

		// Public Catalog Endpoints (Protected by sliding window rate limiter)
		api.Group(func(public chi.Router) {
			public.Use(middleware.RateLimiter(cacheClient, cfg.RateLimitPublicRPM, "public_catalog"))

			public.Get("/products", h.ListProducts)
			public.Get("/products/{id}", h.GetProductByID)
			public.Get("/products/slug/{slug}", h.GetProductBySlug)
		})

		// Admin Mutation Endpoints (Enforces API Gateway X-User-Role: admin & strict rate limit)
		api.Group(func(admin chi.Router) {
			admin.Use(middleware.RequireAdmin)
			admin.Use(middleware.RateLimiter(cacheClient, cfg.RateLimitAdminRPM, "admin_writes"))

			// Product CRUD
			admin.Post("/products", h.CreateProduct)
			admin.Put("/products/{id}", h.UpdateProduct)
			admin.Delete("/products/{id}", h.DeleteProduct)

			// Variant CRUD
			admin.Post("/products/{id}/variants", h.CreateVariant)
			admin.Put("/products/{id}/variants/{variantId}", h.UpdateVariant)
			admin.Delete("/products/{id}/variants/{variantId}", h.DeleteVariant)
		})
	})

	return r
}
