package repository

import (
	"context"
	"errors"

	"github.com/Sall-lah/store_product/internal/db"
)

var (
	// ErrImageNotFound indicates the requested product image asset does not exist in the database.
	ErrImageNotFound = errors.New("product image not found")
)

// ImageRepository provides persistence operations for product gallery and variant media records.
type ImageRepository struct {
	client *db.PrismaClient
}

// NewImageRepository initializes an ImageRepository instance with the active Prisma client.
func NewImageRepository(client *db.PrismaClient) *ImageRepository {
	return &ImageRepository{client: client}
}

// Create persists a new ProductImage record and maintains single-primary image invariants per product.
// Atomic resetting of previous primary flags prevents inconsistent thumbnail states.
func (r *ImageRepository) Create(ctx context.Context, productID string, input CreateProductImageInput) (*ProductImageDTO, error) {
	isPrimary := false
	if input.IsPrimary != nil && *input.IsPrimary {
		isPrimary = true
		// Reset any previously flagged primary image for this product
		_, _ = r.client.ProductImage.FindMany(
			db.ProductImage.ProductID.Equals(productID),
			db.ProductImage.IsPrimary.Equals(true),
		).Update(
			db.ProductImage.IsPrimary.Set(false),
		).Exec(ctx)
	}

	sortOrder := 0
	if input.SortOrder != nil {
		sortOrder = *input.SortOrder
	}

	var optionalParams []db.ProductImageSetParam
	if input.VariantID != nil && *input.VariantID != "" {
		optionalParams = append(optionalParams, db.ProductImage.Variant.Link(db.ProductVariant.ID.Equals(*input.VariantID)))
	}
	if input.AltText != nil {
		optionalParams = append(optionalParams, db.ProductImage.AltText.Set(*input.AltText))
	}
	optionalParams = append(optionalParams, db.ProductImage.IsPrimary.Set(isPrimary))
	optionalParams = append(optionalParams, db.ProductImage.SortOrder.Set(sortOrder))

	created, err := r.client.ProductImage.CreateOne(
		db.ProductImage.Product.Link(db.Product.ID.Equals(productID)),
		db.ProductImage.URL.Set(input.URL),
		db.ProductImage.R2Key.Set(input.R2Key),
		optionalParams...,
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	dto := ToProductImageDTO(created)
	return &dto, nil
}

// GetByID retrieves a single product image record by its unique identifier.
func (r *ImageRepository) GetByID(ctx context.Context, id string) (*ProductImageDTO, error) {
	record, err := r.client.ProductImage.FindUnique(
		db.ProductImage.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, ErrImageNotFound
		}
		return nil, err
	}

	dto := ToProductImageDTO(record)
	return &dto, nil
}

// ListByProduct retrieves all images associated with a product, ordered by display sequence.
// Optional filtering by variant allows querying images scoped to a specific SKU colorway.
func (r *ImageRepository) ListByProduct(ctx context.Context, productID string, variantID *string) ([]ProductImageDTO, error) {
	filters := []db.ProductImageWhereParam{
		db.ProductImage.ProductID.Equals(productID),
	}

	if variantID != nil && *variantID != "" {
		filters = append(filters, db.ProductImage.VariantID.Equals(*variantID))
	}

	records, err := r.client.ProductImage.FindMany(
		filters...,
	).OrderBy(
		db.ProductImage.SortOrder.Order(db.SortOrderAsc),
		db.ProductImage.CreatedAt.Order(db.SortOrderAsc),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	dtos := make([]ProductImageDTO, 0, len(records))
	for _, rec := range records {
		dtos = append(dtos, ToProductImageDTO(&rec))
	}

	return dtos, nil
}

// Update modifies metadata fields on an existing product image record.
// When setting an image to primary, any existing primary flag on the parent product is reset.
func (r *ImageRepository) Update(ctx context.Context, id string, input UpdateProductImageInput) (*ProductImageDTO, error) {
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var updateParams []db.ProductImageSetParam

	if input.IsPrimary != nil {
		if *input.IsPrimary {
			// Reset any existing primary image for this parent product
			_, _ = r.client.ProductImage.FindMany(
				db.ProductImage.ProductID.Equals(existing.ProductID),
				db.ProductImage.IsPrimary.Equals(true),
			).Update(
				db.ProductImage.IsPrimary.Set(false),
			).Exec(ctx)
		}
		updateParams = append(updateParams, db.ProductImage.IsPrimary.Set(*input.IsPrimary))
	}

	if input.VariantID != nil {
		if *input.VariantID != "" {
			updateParams = append(updateParams, db.ProductImage.Variant.Link(db.ProductVariant.ID.Equals(*input.VariantID)))
		} else {
			updateParams = append(updateParams, db.ProductImage.Variant.Unlink())
		}
	}

	if input.AltText != nil {
		updateParams = append(updateParams, db.ProductImage.AltText.Set(*input.AltText))
	}

	if input.SortOrder != nil {
		updateParams = append(updateParams, db.ProductImage.SortOrder.Set(*input.SortOrder))
	}

	if len(updateParams) == 0 {
		return existing, nil
	}

	updated, err := r.client.ProductImage.FindUnique(
		db.ProductImage.ID.Equals(id),
	).Update(
		updateParams...,
	).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, ErrImageNotFound
		}
		return nil, err
	}

	dto := ToProductImageDTO(updated)
	return &dto, nil
}

// Delete removes an image record from the database and returns the deleted entity.
// Returning the deleted entity supplies the caller with the R2Key needed for R2 bucket file cleanup.
func (r *ImageRepository) Delete(ctx context.Context, id string) (*ProductImageDTO, error) {
	deleted, err := r.client.ProductImage.FindUnique(
		db.ProductImage.ID.Equals(id),
	).Delete().Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, ErrImageNotFound
		}
		return nil, err
	}

	dto := ToProductImageDTO(deleted)
	return &dto, nil
}
