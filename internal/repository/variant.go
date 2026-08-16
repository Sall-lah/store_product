package repository

import (
	"context"
	"errors"

	"github.com/Sall-lah/store_product/internal/db"
)

var (
	// ErrVariantNotFound indicates the requested variant entity does not exist.
	ErrVariantNotFound = errors.New("variant not found")
	// ErrDuplicateSKU indicates a duplicate stock keeping unit was supplied.
	ErrDuplicateSKU = errors.New("a variant with this SKU already exists")
)

// VariantRepository provides persistence operations for product variants (size/color combinations).
type VariantRepository struct {
	client *db.PrismaClient
}

// NewVariantRepository creates a new instance of VariantRepository.
func NewVariantRepository(client *db.PrismaClient) *VariantRepository {
	return &VariantRepository{client: client}
}

// GetVariantByID retrieves an individual variant entity by its ID.
func (r *VariantRepository) GetVariantByID(ctx context.Context, id string) (*VariantDTO, error) {
	record, err := r.client.ProductVariant.FindUnique(
		db.ProductVariant.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, ErrVariantNotFound
		}
		return nil, err
	}

	dto := ToVariantDTO(record)
	return &dto, nil
}

// CreateVariant creates and attaches a new variant to a parent product.
func (r *VariantRepository) CreateVariant(ctx context.Context, productID string, input CreateVariantInput) (*VariantDTO, error) {
	createParams := []db.ProductVariantSetParam{
		db.ProductVariant.Stock.Set(input.Stock),
	}

	if input.Size != nil {
		createParams = append(createParams, db.ProductVariant.Size.Set(*input.Size))
	}
	if input.Color != nil {
		createParams = append(createParams, db.ProductVariant.Color.Set(*input.Color))
	}
	if input.Price != nil {
		createParams = append(createParams, db.ProductVariant.Price.Set(*input.Price))
	}
	if input.IsActive != nil {
		createParams = append(createParams, db.ProductVariant.IsActive.Set(*input.IsActive))
	}

	record, err := r.client.ProductVariant.CreateOne(
		db.ProductVariant.Product.Link(db.Product.ID.Equals(productID)),
		db.ProductVariant.Sku.Set(input.SKU),
		createParams...,
	).Exec(ctx)

	if err != nil {
		if _, isUnique := db.IsErrUniqueConstraint(err); isUnique {
			return nil, ErrDuplicateSKU
		}
		return nil, err
	}

	dto := ToVariantDTO(record)
	return &dto, nil
}

// UpdateVariant modifies attributes of an existing variant.
func (r *VariantRepository) UpdateVariant(ctx context.Context, id string, input CreateVariantInput) (*VariantDTO, error) {
	var updateParams []db.ProductVariantSetParam

	if input.SKU != "" {
		updateParams = append(updateParams, db.ProductVariant.Sku.Set(input.SKU))
	}
	if input.Size != nil {
		updateParams = append(updateParams, db.ProductVariant.Size.Set(*input.Size))
	}
	if input.Color != nil {
		updateParams = append(updateParams, db.ProductVariant.Color.Set(*input.Color))
	}
	if input.Price != nil {
		updateParams = append(updateParams, db.ProductVariant.Price.Set(*input.Price))
	}
	if input.Stock >= 0 {
		updateParams = append(updateParams, db.ProductVariant.Stock.Set(input.Stock))
	}
	if input.IsActive != nil {
		updateParams = append(updateParams, db.ProductVariant.IsActive.Set(*input.IsActive))
	}

	record, err := r.client.ProductVariant.FindUnique(
		db.ProductVariant.ID.Equals(id),
	).Update(updateParams...).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, ErrVariantNotFound
		}
		if _, isUnique := db.IsErrUniqueConstraint(err); isUnique {
			return nil, ErrDuplicateSKU
		}
		return nil, err
	}

	dto := ToVariantDTO(record)
	return &dto, nil
}

// DeleteVariant deletes an existing variant entity from the database.
func (r *VariantRepository) DeleteVariant(ctx context.Context, id string) error {
	_, err := r.client.ProductVariant.FindUnique(
		db.ProductVariant.ID.Equals(id),
	).Delete().Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return ErrVariantNotFound
		}
		return err
	}

	return nil
}
