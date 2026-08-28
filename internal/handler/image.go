package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Sall-lah/store_product/internal/repository"
	"github.com/Sall-lah/store_product/internal/service"
	"github.com/go-chi/chi/v5"
)

// ImageHandler handles admin HTTP requests for presigning, creating, listing, updating, and deleting product images.
type ImageHandler struct {
	imageService *service.ImageService
}

// NewImageHandler constructs an ImageHandler instance.
func NewImageHandler(imageService *service.ImageService) *ImageHandler {
	return &ImageHandler{imageService: imageService}
}

// PresignImage handles POST /api/v1/admin/products/{id}/images/presign.
// Generates a short-lived direct upload URL to Cloudflare R2 after validating MIME type.
func (h *ImageHandler) PresignImage(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	if strings.TrimSpace(productID) == "" {
		renderError(w, http.StatusBadRequest, "invalid_product_id", "Product ID is required.")
		return
	}

	var input repository.PresignImageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		renderError(w, http.StatusBadRequest, "invalid_payload", "Invalid JSON request body.")
		return
	}

	resp, err := h.imageService.GeneratePresignedURL(r.Context(), productID, input)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			renderError(w, http.StatusNotFound, "product_not_found", "Product not found.")
			return
		}
		if errors.Is(err, service.ErrInvalidContentType) {
			renderError(w, http.StatusBadRequest, "invalid_content_type", "Invalid or unsupported image content type. Allowed: JPEG, PNG, WEBP, AVIF, GIF.")
			return
		}
		if errors.Is(err, service.ErrMissingFileName) {
			renderError(w, http.StatusBadRequest, "missing_filename", "File name is required.")
			return
		}
		renderError(w, http.StatusInternalServerError, "presign_failed", "Failed to generate presigned upload URL.")
		return
	}

	renderJSON(w, http.StatusOK, resp)
}

// CreateImage handles POST /api/v1/admin/products/{id}/images.
// Confirms and registers the uploaded image in PostgreSQL and sets primary/order metadata.
func (h *ImageHandler) CreateImage(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	if strings.TrimSpace(productID) == "" {
		renderError(w, http.StatusBadRequest, "invalid_product_id", "Product ID is required.")
		return
	}

	var input repository.CreateProductImageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		renderError(w, http.StatusBadRequest, "invalid_payload", "Invalid JSON request body.")
		return
	}

	created, err := h.imageService.CreateImage(r.Context(), productID, input)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			renderError(w, http.StatusNotFound, "product_not_found", "Product not found.")
			return
		}
		if errors.Is(err, service.ErrMissingImagePayload) {
			renderError(w, http.StatusBadRequest, "invalid_payload", "URL and r2_key are required.")
			return
		}
		renderError(w, http.StatusInternalServerError, "create_image_failed", "Failed to register product image.")
		return
	}

	renderJSON(w, http.StatusCreated, created)
}

// ListImages handles GET /api/v1/admin/products/{id}/images.
// Lists all media assets attached to a product, optionally filtered by variant ID.
func (h *ImageHandler) ListImages(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	if strings.TrimSpace(productID) == "" {
		renderError(w, http.StatusBadRequest, "invalid_product_id", "Product ID is required.")
		return
	}

	var variantIDPtr *string
	if variantID := strings.TrimSpace(r.URL.Query().Get("variant_id")); variantID != "" {
		variantIDPtr = &variantID
	}

	images, err := h.imageService.ListImages(r.Context(), productID, variantIDPtr)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			renderError(w, http.StatusNotFound, "product_not_found", "Product not found.")
			return
		}
		renderError(w, http.StatusInternalServerError, "list_images_failed", "Failed to retrieve product images.")
		return
	}

	renderJSON(w, http.StatusOK, images)
}

// UpdateImage handles PUT /api/v1/admin/products/{id}/images/{imageId}.
// Updates image metadata such as alt text, primary flag, display sort order, and variant link.
func (h *ImageHandler) UpdateImage(w http.ResponseWriter, r *http.Request) {
	imageID := chi.URLParam(r, "imageId")
	if strings.TrimSpace(imageID) == "" {
		renderError(w, http.StatusBadRequest, "invalid_image_id", "Image ID is required.")
		return
	}

	var input repository.UpdateProductImageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		renderError(w, http.StatusBadRequest, "invalid_payload", "Invalid JSON request body.")
		return
	}

	updated, err := h.imageService.UpdateImage(r.Context(), imageID, input)
	if err != nil {
		if errors.Is(err, repository.ErrImageNotFound) {
			renderError(w, http.StatusNotFound, "image_not_found", "Product image not found.")
			return
		}
		renderError(w, http.StatusInternalServerError, "update_image_failed", "Failed to update product image.")
		return
	}

	renderJSON(w, http.StatusOK, updated)
}

// DeleteImage handles DELETE /api/v1/admin/products/{id}/images/{imageId}.
// Deletes the database record and removes the physical asset from Cloudflare R2.
func (h *ImageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	imageID := chi.URLParam(r, "imageId")
	if strings.TrimSpace(imageID) == "" {
		renderError(w, http.StatusBadRequest, "invalid_image_id", "Image ID is required.")
		return
	}

	if err := h.imageService.DeleteImage(r.Context(), imageID); err != nil {
		if errors.Is(err, repository.ErrImageNotFound) {
			renderError(w, http.StatusNotFound, "image_not_found", "Product image not found.")
			return
		}
		renderError(w, http.StatusInternalServerError, "delete_image_failed", "Failed to delete product image.")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
