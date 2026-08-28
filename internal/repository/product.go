package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/Sall-lah/store_product/internal/db"
	"github.com/Sall-lah/store_product/internal/pkg/cursor"
)

var (
	// ErrProductNotFound indicates the requested product does not exist in the database.
	ErrProductNotFound = errors.New("product not found")
	// ErrDuplicateSlug indicates another product already claims the provided slug.
	ErrDuplicateSlug = errors.New("product with this slug already exists")
)

// ProductRepository provides type-safe database queries for product management.
type ProductRepository struct {
	client *db.PrismaClient
}

// NewProductRepository creates a new instance of ProductRepository.
func NewProductRepository(client *db.PrismaClient) *ProductRepository {
	return &ProductRepository{client: client}
}

// ListProducts executes a keyset-paginated query against Supabase with multi-field filters.
func (r *ProductRepository) ListProducts(ctx context.Context, filter ProductFilter) (*PaginatedProducts, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	var where []db.ProductWhereParam

	if filter.IsActive != nil {
		where = append(where, db.Product.IsActive.Equals(*filter.IsActive))
	} else if !filter.IncludeInactive {
		where = append(where, db.Product.IsActive.Equals(true))
	}

	if filter.Category != nil && strings.TrimSpace(*filter.Category) != "" {
		where = append(where, db.Product.Category.Equals(*filter.Category))
	}

	if filter.MinPrice != nil {
		where = append(where, db.Product.BasePrice.Gte(*filter.MinPrice))
	}

	if filter.MaxPrice != nil {
		where = append(where, db.Product.BasePrice.Lte(*filter.MaxPrice))
	}

	// Variant-level filtering (size, color)
	var variantWhere []db.ProductVariantWhereParam
	if !filter.IncludeInactive {
		variantWhere = append(variantWhere, db.ProductVariant.IsActive.Equals(true))
	}
	hasVariantFilter := false

	if filter.Size != nil && strings.TrimSpace(*filter.Size) != "" {
		variantWhere = append(variantWhere, db.ProductVariant.Size.Equals(*filter.Size))
		hasVariantFilter = true
	}
	if filter.Color != nil && strings.TrimSpace(*filter.Color) != "" {
		variantWhere = append(variantWhere, db.ProductVariant.Color.Equals(*filter.Color))
		hasVariantFilter = true
	}
	if hasVariantFilter {
		where = append(where, db.Product.Variants.Some(variantWhere...))
	}

	// Search by name, description, or variant SKU
	if filter.Search != nil && strings.TrimSpace(*filter.Search) != "" {
		searchTerm := strings.TrimSpace(*filter.Search)
		where = append(where, db.Product.Or(
			db.Product.Name.Contains(searchTerm),
			db.Product.Description.Contains(searchTerm),
			db.Product.Variants.Some(db.ProductVariant.Sku.Contains(searchTerm)),
		))
	}

	// Apply keyset cursor: (created_at, id) < (cursor_created_at, cursor_id)
	if filter.Cursor != nil && strings.TrimSpace(*filter.Cursor) != "" {
		decoded, err := cursor.Decode(*filter.Cursor)
		if err == nil && decoded != nil {
			where = append(where, db.Product.Or(
				db.Product.CreatedAt.Lt(decoded.CreatedAt),
				db.Product.And(
					db.Product.CreatedAt.Equals(decoded.CreatedAt),
					db.Product.ID.Lt(decoded.ID),
				),
			))
		}
	}

	query := r.client.Product.FindMany(where...)
	if filter.IncludeInactive {
		query = query.With(
			db.Product.Variants.Fetch(),
			db.Product.Images.Fetch().OrderBy(
				db.ProductImage.SortOrder.Order(db.SortOrderAsc),
				db.ProductImage.CreatedAt.Order(db.SortOrderAsc),
			),
		)
	} else {
		query = query.With(
			db.Product.Variants.Fetch(db.ProductVariant.IsActive.Equals(true)),
			db.Product.Images.Fetch().OrderBy(
				db.ProductImage.SortOrder.Order(db.SortOrderAsc),
				db.ProductImage.CreatedAt.Order(db.SortOrderAsc),
			),
		)
	}

	// Fetch limit + 1 items to determine if a subsequent page exists without a separate COUNT(*) query
	records, err := query.
		OrderBy(
			db.Product.CreatedAt.Order(db.SortOrderDesc),
			db.Product.ID.Order(db.SortOrderDesc),
		).
		Take(limit + 1).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	hasMore := len(records) > limit
	var items []ProductDTO
	count := len(records)
	if hasMore {
		count = limit
	}

	for i := 0; i < count; i++ {
		items = append(items, ToProductDTO(&records[i]))
	}

	var nextCursor string
	if hasMore && len(items) > 0 {
		lastItem := items[len(items)-1]
		nextCursor = cursor.Encode(lastItem.CreatedAt, lastItem.ID)
	}

	return &PaginatedProducts{
		Items:      items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
		Limit:      limit,
	}, nil
}

// GetProductByID fetches an individual active product and its active variants by UUID.
// Returns ErrProductNotFound if the product does not exist or has been soft-deleted (isActive = false).
func (r *ProductRepository) GetProductByID(ctx context.Context, id string) (*ProductDTO, error) {
	record, err := r.client.Product.FindUnique(
		db.Product.ID.Equals(id),
	).With(
		db.Product.Variants.Fetch(db.ProductVariant.IsActive.Equals(true)),
		db.Product.Images.Fetch().OrderBy(
			db.ProductImage.SortOrder.Order(db.SortOrderAsc),
			db.ProductImage.CreatedAt.Order(db.SortOrderAsc),
		),
	).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if !record.IsActive {
		return nil, ErrProductNotFound
	}

	dto := ToProductDTO(record)
	return &dto, nil
}

// GetAdminProductByID fetches an individual product and all its variants (active and inactive) by UUID for backoffice inspection.
func (r *ProductRepository) GetAdminProductByID(ctx context.Context, id string) (*ProductDTO, error) {
	record, err := r.client.Product.FindUnique(
		db.Product.ID.Equals(id),
	).With(
		db.Product.Variants.Fetch(),
		db.Product.Images.Fetch().OrderBy(
			db.ProductImage.SortOrder.Order(db.SortOrderAsc),
			db.ProductImage.CreatedAt.Order(db.SortOrderAsc),
		),
	).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	dto := ToProductDTO(record)
	return &dto, nil
}

// GetProductBySlug fetches an individual active product and its variants by human-readable URL slug.
// Returns ErrProductNotFound if the product does not exist or has been soft-deleted (isActive = false).
func (r *ProductRepository) GetProductBySlug(ctx context.Context, slug string) (*ProductDTO, error) {
	record, err := r.client.Product.FindUnique(
		db.Product.Slug.Equals(slug),
	).With(
		db.Product.Variants.Fetch(db.ProductVariant.IsActive.Equals(true)),
		db.Product.Images.Fetch().OrderBy(
			db.ProductImage.SortOrder.Order(db.SortOrderAsc),
			db.ProductImage.CreatedAt.Order(db.SortOrderAsc),
		),
	).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if !record.IsActive {
		return nil, ErrProductNotFound
	}

	dto := ToProductDTO(record)
	return &dto, nil
}

// CreateProduct creates a new product in Supabase.
func (r *ProductRepository) CreateProduct(ctx context.Context, input CreateProductInput) (*ProductDTO, error) {
	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	createParams := []db.ProductSetParam{
		db.Product.Name.Set(input.Name),
		db.Product.Slug.Set(input.Slug),
		db.Product.BasePrice.Set(input.BasePrice),
		db.Product.Category.Set(input.Category),
		db.Product.IsActive.Set(isActive),
	}

	if input.Description != nil {
		createParams = append(createParams, db.Product.Description.Set(*input.Description))
	}

	record, err := r.client.Product.CreateOne(
		db.Product.Name.Set(input.Name),
		db.Product.Slug.Set(input.Slug),
		db.Product.BasePrice.Set(input.BasePrice),
		db.Product.Category.Set(input.Category),
		createParams...,
	).Exec(ctx)

	if err != nil {
		if _, isUnique := db.IsErrUniqueConstraint(err); isUnique {
			return nil, ErrDuplicateSlug
		}
		return nil, err
	}

	// Create initial variants if provided
	if len(input.Variants) > 0 {
		for _, v := range input.Variants {
			varCreateParams := []db.ProductVariantSetParam{
				db.ProductVariant.Stock.Set(v.Stock),
			}
			if v.Size != nil {
				varCreateParams = append(varCreateParams, db.ProductVariant.Size.Set(*v.Size))
			}
			if v.Color != nil {
				varCreateParams = append(varCreateParams, db.ProductVariant.Color.Set(*v.Color))
			}
			if v.Price != nil {
				varCreateParams = append(varCreateParams, db.ProductVariant.Price.Set(*v.Price))
			}
			if v.IsActive != nil {
				varCreateParams = append(varCreateParams, db.ProductVariant.IsActive.Set(*v.IsActive))
			}

			_, varErr := r.client.ProductVariant.CreateOne(
				db.ProductVariant.Product.Link(db.Product.ID.Equals(record.ID)),
				db.ProductVariant.Sku.Set(v.SKU),
				varCreateParams...,
			).Exec(ctx)
			if varErr != nil {
				return nil, varErr
			}
		}
	}

	// Fetch newly created product with relations populated
	return r.GetAdminProductByID(ctx, record.ID)
}

// UpdateProduct updates mutable fields on an existing product.
func (r *ProductRepository) UpdateProduct(ctx context.Context, id string, input UpdateProductInput) (*ProductDTO, error) {
	var updateParams []db.ProductSetParam

	if input.Name != nil {
		updateParams = append(updateParams, db.Product.Name.Set(*input.Name))
	}
	if input.Slug != nil {
		updateParams = append(updateParams, db.Product.Slug.Set(*input.Slug))
	}
	if input.Description != nil {
		updateParams = append(updateParams, db.Product.Description.Set(*input.Description))
	}
	if input.BasePrice != nil {
		updateParams = append(updateParams, db.Product.BasePrice.Set(*input.BasePrice))
	}
	if input.Category != nil {
		updateParams = append(updateParams, db.Product.Category.Set(*input.Category))
	}
	if input.IsActive != nil {
		updateParams = append(updateParams, db.Product.IsActive.Set(*input.IsActive))
	}

	if len(updateParams) == 0 {
		return r.GetAdminProductByID(ctx, id)
	}

	record, err := r.client.Product.FindUnique(
		db.Product.ID.Equals(id),
	).Update(updateParams...).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, ErrProductNotFound
		}
		if _, isUnique := db.IsErrUniqueConstraint(err); isUnique {
			return nil, ErrDuplicateSlug
		}
		return nil, err
	}

	return r.GetAdminProductByID(ctx, record.ID)
}

// DeleteProduct performs a soft delete by setting isActive to false on the product and all associated variants.
// Soft deletion preserves audit history, SKU integrity, and historical order references.
func (r *ProductRepository) DeleteProduct(ctx context.Context, id string) error {
	_, err := r.client.Product.FindUnique(
		db.Product.ID.Equals(id),
	).Update(
		db.Product.IsActive.Set(false),
	).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	// Cascade soft-deactivation to all child variants to prevent orphaned active variants
	_, err = r.client.Prisma.ExecuteRaw(
		`UPDATE "ProductVariant" SET "isActive" = false, "updatedAt" = NOW() WHERE "productId" = $1`,
		id,
	).Exec(ctx)

	return err
}
