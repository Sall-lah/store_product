package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Sall-lah/store_product/internal/repository"
	"github.com/Sall-lah/store_product/internal/service"
	"github.com/go-chi/chi/v5"
)

// ProductHandler handles incoming HTTP requests for product queries and mutations.
type ProductHandler struct {
	service *service.ProductService
}

// NewProductHandler constructs a new ProductHandler instance.
func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{service: svc}
}

// ListProducts handles GET /api/v1/products with filtering, searching, and keyset pagination.
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := repository.ProductFilter{
		Limit: 20,
	}

	if limitStr := q.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	if cat := strings.TrimSpace(q.Get("category")); cat != "" {
		filter.Category = &cat
	}

	if minPriceStr := q.Get("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filter.MinPrice = &minPrice
		}
	}

	if maxPriceStr := q.Get("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filter.MaxPrice = &maxPrice
		}
	}

	if size := strings.TrimSpace(q.Get("size")); size != "" {
		filter.Size = &size
	}

	if color := strings.TrimSpace(q.Get("color")); color != "" {
		filter.Color = &color
	}

	if search := strings.TrimSpace(q.Get("search")); search != "" {
		filter.Search = &search
	}

	if cursorStr := strings.TrimSpace(q.Get("cursor")); cursorStr != "" {
		filter.Cursor = &cursorStr
	}

	result, err := h.service.ListProducts(r.Context(), filter)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "query_failed", "Failed to retrieve products.")
		return
	}

	renderJSON(w, http.StatusOK, result)
}

// GetProductByID handles GET /api/v1/products/{id}.
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		renderError(w, http.StatusBadRequest, "invalid_id", "Product ID is required.")
		return
	}

	product, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			renderError(w, http.StatusNotFound, "product_not_found", "Product not found.")
			return
		}
		renderError(w, http.StatusInternalServerError, "query_failed", "Failed to fetch product.")
		return
	}

	renderJSON(w, http.StatusOK, product)
}

// GetProductBySlug handles GET /api/v1/products/slug/{slug}.
func (h *ProductHandler) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		renderError(w, http.StatusBadRequest, "invalid_slug", "Product slug is required.")
		return
	}

	product, err := h.service.GetProductBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			renderError(w, http.StatusNotFound, "product_not_found", "Product not found.")
			return
		}
		renderError(w, http.StatusInternalServerError, "query_failed", "Failed to fetch product.")
		return
	}

	renderJSON(w, http.StatusOK, product)
}

// CreateProduct handles POST /api/v1/products. Admin only.
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var input repository.CreateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		renderError(w, http.StatusBadRequest, "invalid_payload", "Invalid JSON payload.")
		return
	}

	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Slug) == "" || strings.TrimSpace(input.Category) == "" {
		renderError(w, http.StatusBadRequest, "validation_error", "Name, slug, and category are required.")
		return
	}

	if input.BasePrice < 0 {
		renderError(w, http.StatusBadRequest, "validation_error", "Base price must not be negative.")
		return
	}

	product, err := h.service.CreateProduct(r.Context(), input)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateSlug) {
			renderError(w, http.StatusConflict, "duplicate_slug", "A product with this slug already exists.")
			return
		}
		if errors.Is(err, repository.ErrDuplicateSKU) {
			renderError(w, http.StatusConflict, "duplicate_sku", "A variant with this SKU already exists.")
			return
		}
		renderError(w, http.StatusInternalServerError, "create_failed", "Failed to create product.")
		return
	}

	renderJSON(w, http.StatusCreated, product)
}

// UpdateProduct handles PUT /api/v1/products/{id}. Admin only.
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		renderError(w, http.StatusBadRequest, "invalid_id", "Product ID is required.")
		return
	}

	var input repository.UpdateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		renderError(w, http.StatusBadRequest, "invalid_payload", "Invalid JSON payload.")
		return
	}

	product, err := h.service.UpdateProduct(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			renderError(w, http.StatusNotFound, "product_not_found", "Product not found.")
			return
		}
		if errors.Is(err, repository.ErrDuplicateSlug) {
			renderError(w, http.StatusConflict, "duplicate_slug", "A product with this slug already exists.")
			return
		}
		renderError(w, http.StatusInternalServerError, "update_failed", "Failed to update product.")
		return
	}

	renderJSON(w, http.StatusOK, product)
}

// DeleteProduct handles DELETE /api/v1/products/{id}. Admin only.
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		renderError(w, http.StatusBadRequest, "invalid_id", "Product ID is required.")
		return
	}

	if err := h.service.DeleteProduct(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			renderError(w, http.StatusNotFound, "product_not_found", "Product not found.")
			return
		}
		renderError(w, http.StatusInternalServerError, "delete_failed", "Failed to delete product.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreateVariant handles POST /api/v1/products/{id}/variants. Admin only.
func (h *ProductHandler) CreateVariant(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	if productID == "" {
		renderError(w, http.StatusBadRequest, "invalid_id", "Product ID is required.")
		return
	}

	var input repository.CreateVariantInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		renderError(w, http.StatusBadRequest, "invalid_payload", "Invalid JSON payload.")
		return
	}

	if strings.TrimSpace(input.SKU) == "" {
		renderError(w, http.StatusBadRequest, "validation_error", "SKU is required for a variant.")
		return
	}

	variant, err := h.service.CreateVariant(r.Context(), productID, input)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateSKU) {
			renderError(w, http.StatusConflict, "duplicate_sku", "A variant with this SKU already exists.")
			return
		}
		renderError(w, http.StatusInternalServerError, "create_variant_failed", "Failed to create variant.")
		return
	}

	renderJSON(w, http.StatusCreated, variant)
}

// UpdateVariant handles PUT /api/v1/products/{id}/variants/{variantId}. Admin only.
func (h *ProductHandler) UpdateVariant(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	variantID := chi.URLParam(r, "variantId")

	var input repository.CreateVariantInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		renderError(w, http.StatusBadRequest, "invalid_payload", "Invalid JSON payload.")
		return
	}

	variant, err := h.service.UpdateVariant(r.Context(), productID, variantID, input)
	if err != nil {
		if errors.Is(err, repository.ErrVariantNotFound) {
			renderError(w, http.StatusNotFound, "variant_not_found", "Variant not found.")
			return
		}
		if errors.Is(err, repository.ErrDuplicateSKU) {
			renderError(w, http.StatusConflict, "duplicate_sku", "A variant with this SKU already exists.")
			return
		}
		renderError(w, http.StatusInternalServerError, "update_variant_failed", "Failed to update variant.")
		return
	}

	renderJSON(w, http.StatusOK, variant)
}

// DeleteVariant handles DELETE /api/v1/products/{id}/variants/{variantId}. Admin only.
func (h *ProductHandler) DeleteVariant(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	variantID := chi.URLParam(r, "variantId")

	if err := h.service.DeleteVariant(r.Context(), productID, variantID); err != nil {
		if errors.Is(err, repository.ErrVariantNotFound) {
			renderError(w, http.StatusNotFound, "variant_not_found", "Variant not found.")
			return
		}
		renderError(w, http.StatusInternalServerError, "delete_variant_failed", "Failed to delete variant.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HealthCheck provides a liveness and readiness probe for container orchestrators.
func (h *ProductHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	renderJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

// renderJSON serializes an object to JSON and writes it with the appropriate HTTP status code.
func renderJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// renderError emits a structured JSON error response.
func renderError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":   code,
		"message": message,
	})
}
