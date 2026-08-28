package repository

import (
	"time"

	"github.com/Sall-lah/store_product/internal/db"
)

// ProductDTO represents the decoupled domain representation of a product entity.
// Using explicit DTOs isolates transport and presentation layers from the Prisma ORM internal models.
type ProductDTO struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Slug            string            `json:"slug"`
	Description     *string           `json:"description,omitempty"`
	BasePrice       float64           `json:"base_price"`
	Category        string            `json:"category"`
	IsActive        bool              `json:"is_active"`
	PrimaryImageURL *string           `json:"primary_image_url,omitempty"`
	Images          []ProductImageDTO `json:"images,omitempty"`
	Variants        []VariantDTO      `json:"variants"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

// VariantDTO represents the decoupled domain representation of a product variant (size, color, sku, stock).
type VariantDTO struct {
	ID        string            `json:"id"`
	ProductID string            `json:"product_id"`
	SKU       string            `json:"sku"`
	Size      *string           `json:"size,omitempty"`
	Color     *string           `json:"color,omitempty"`
	Price     *float64          `json:"price,omitempty"`
	Stock     int               `json:"stock"`
	IsActive  bool              `json:"is_active"`
	Images    []ProductImageDTO `json:"images,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// ProductImageDTO represents the domain representation of a product media asset stored in Cloudflare R2.
type ProductImageDTO struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	VariantID *string   `json:"variant_id,omitempty"`
	URL       string    `json:"url"`
	R2Key     string    `json:"r2_key"`
	AltText   *string   `json:"alt_text,omitempty"`
	IsPrimary bool      `json:"is_primary"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PresignImageInput defines the request body for generating a direct-to-R2 upload URL.
type PresignImageInput struct {
	FileName    string  `json:"file_name"`
	ContentType string  `json:"content_type"`
	VariantID   *string `json:"variant_id,omitempty"`
}

// PresignImageResponse encapsulates the signed PUT URL and asset public reference.
type PresignImageResponse struct {
	UploadURL        string `json:"upload_url"`
	PublicURL        string `json:"public_url"`
	R2Key            string `json:"r2_key"`
	ExpiresInSeconds int64  `json:"expires_in_seconds"`
}

// CreateProductImageInput defines the confirmation payload after an asset is uploaded to R2.
type CreateProductImageInput struct {
	URL       string  `json:"url"`
	R2Key     string  `json:"r2_key"`
	VariantID *string `json:"variant_id,omitempty"`
	AltText   *string `json:"alt_text,omitempty"`
	IsPrimary *bool   `json:"is_primary,omitempty"`
	SortOrder *int    `json:"sort_order,omitempty"`
}

// UpdateProductImageInput defines mutable metadata fields on an existing product image.
type UpdateProductImageInput struct {
	VariantID *string `json:"variant_id,omitempty"`
	AltText   *string `json:"alt_text,omitempty"`
	IsPrimary *bool   `json:"is_primary,omitempty"`
	SortOrder *int    `json:"sort_order,omitempty"`
}

// ProductFilter encapsulates query parameters for searching, filtering, and keyset cursor pagination.
type ProductFilter struct {
	Category        *string
	MinPrice        *float64
	MaxPrice        *float64
	Size            *string
	Color           *string
	Search          *string
	Cursor          *string
	Limit           int
	IncludeInactive bool
	IsActive        *bool
}

// PaginatedProducts encapsulates the keyset-paginated result payload.
type PaginatedProducts struct {
	Items      []ProductDTO `json:"items"`
	NextCursor string       `json:"next_cursor,omitempty"`
	HasMore    bool         `json:"has_more"`
	Limit      int          `json:"limit"`
}

// CreateVariantInput defines payload attributes for adding a product variant.
type CreateVariantInput struct {
	SKU      string   `json:"sku"`
	Size     *string  `json:"size,omitempty"`
	Color    *string  `json:"color,omitempty"`
	Price    *float64 `json:"price,omitempty"`
	Stock    int      `json:"stock"`
	IsActive *bool    `json:"is_active,omitempty"`
}

// CreateProductInput defines payload attributes for creating a new product with optional initial variants.
type CreateProductInput struct {
	Name        string               `json:"name"`
	Slug        string               `json:"slug"`
	Description *string              `json:"description,omitempty"`
	BasePrice   float64              `json:"base_price"`
	Category    string               `json:"category"`
	IsActive    *bool                `json:"is_active,omitempty"`
	Variants    []CreateVariantInput `json:"variants,omitempty"`
}

// UpdateProductInput defines mutable attributes on an existing product entity.
type UpdateProductInput struct {
	Name        *string  `json:"name,omitempty"`
	Slug        *string  `json:"slug,omitempty"`
	Description *string  `json:"description,omitempty"`
	BasePrice   *float64 `json:"base_price,omitempty"`
	Category    *string  `json:"category,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

// ToProductDTO converts an internal Prisma ProductModel into a clean domain ProductDTO.
func ToProductDTO(m *db.ProductModel) ProductDTO {
	if m == nil {
		return ProductDTO{}
	}

	dto := ProductDTO{
		ID:        m.InnerProduct.ID,
		Name:      m.InnerProduct.Name,
		Slug:      m.InnerProduct.Slug,
		BasePrice: m.InnerProduct.BasePrice,
		Category:  m.InnerProduct.Category,
		IsActive:  m.InnerProduct.IsActive,
		CreatedAt: m.InnerProduct.CreatedAt,
		UpdatedAt: m.InnerProduct.UpdatedAt,
		Variants:  make([]VariantDTO, 0),
		Images:    make([]ProductImageDTO, 0),
	}

	if desc, ok := m.Description(); ok {
		dto.Description = &desc
	}

	if m.RelationsProduct.Variants != nil {
		for _, v := range m.RelationsProduct.Variants {
			dto.Variants = append(dto.Variants, ToVariantDTO(&v))
		}
	}

	if m.RelationsProduct.Images != nil {
		for _, img := range m.RelationsProduct.Images {
			imgDTO := ToProductImageDTO(&img)
			dto.Images = append(dto.Images, imgDTO)
			if imgDTO.IsPrimary && dto.PrimaryImageURL == nil {
				urlCopy := imgDTO.URL
				dto.PrimaryImageURL = &urlCopy
			}
		}
		// Fallback: if no primary image is explicitly flagged, use the first gallery image
		if dto.PrimaryImageURL == nil && len(dto.Images) > 0 {
			firstURL := dto.Images[0].URL
			dto.PrimaryImageURL = &firstURL
		}
	}

	return dto
}

// ToVariantDTO converts an internal Prisma ProductVariantModel into a domain VariantDTO.
func ToVariantDTO(v *db.ProductVariantModel) VariantDTO {
	if v == nil {
		return VariantDTO{}
	}

	dto := VariantDTO{
		ID:        v.InnerProductVariant.ID,
		ProductID: v.InnerProductVariant.ProductID,
		SKU:       v.InnerProductVariant.Sku,
		Stock:     v.InnerProductVariant.Stock,
		IsActive:  v.InnerProductVariant.IsActive,
		CreatedAt: v.InnerProductVariant.CreatedAt,
		UpdatedAt: v.InnerProductVariant.UpdatedAt,
		Images:    make([]ProductImageDTO, 0),
	}

	if size, ok := v.Size(); ok {
		dto.Size = &size
	}
	if color, ok := v.Color(); ok {
		dto.Color = &color
	}
	if price, ok := v.Price(); ok {
		dto.Price = &price
	}

	if v.RelationsProductVariant.Images != nil {
		for _, img := range v.RelationsProductVariant.Images {
			dto.Images = append(dto.Images, ToProductImageDTO(&img))
		}
	}

	return dto
}

// ToProductImageDTO converts an internal Prisma ProductImageModel into a domain ProductImageDTO.
func ToProductImageDTO(m *db.ProductImageModel) ProductImageDTO {
	if m == nil {
		return ProductImageDTO{}
	}

	dto := ProductImageDTO{
		ID:        m.InnerProductImage.ID,
		ProductID: m.InnerProductImage.ProductID,
		URL:       m.InnerProductImage.URL,
		R2Key:     m.InnerProductImage.R2Key,
		IsPrimary: m.InnerProductImage.IsPrimary,
		SortOrder: m.InnerProductImage.SortOrder,
		CreatedAt: m.InnerProductImage.CreatedAt,
		UpdatedAt: m.InnerProductImage.UpdatedAt,
	}

	if variantID, ok := m.VariantID(); ok {
		dto.VariantID = &variantID
	}
	if altText, ok := m.AltText(); ok {
		dto.AltText = &altText
	}

	return dto
}
