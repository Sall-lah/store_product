package repository

import (
	"time"

	"github.com/Sall-lah/store_product/internal/db"
)

// ProductDTO represents the decoupled domain representation of a product entity.
// Using explicit DTOs isolates transport and presentation layers from the Prisma ORM internal models.
type ProductDTO struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Slug        string       `json:"slug"`
	Description *string      `json:"description,omitempty"`
	BasePrice   float64      `json:"base_price"`
	Category    string       `json:"category"`
	IsActive    bool         `json:"is_active"`
	Variants    []VariantDTO `json:"variants"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// VariantDTO represents the decoupled domain representation of a product variant (size, color, sku, stock).
type VariantDTO struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	SKU       string    `json:"sku"`
	Size      *string   `json:"size,omitempty"`
	Color     *string   `json:"color,omitempty"`
	Price     *float64  `json:"price,omitempty"`
	Stock     int       `json:"stock"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	}

	if desc, ok := m.Description(); ok {
		dto.Description = &desc
	}

	if m.RelationsProduct.Variants != nil {
		for _, v := range m.RelationsProduct.Variants {
			dto.Variants = append(dto.Variants, ToVariantDTO(&v))
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

	return dto
}
