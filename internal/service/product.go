package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Sall-lah/store_product/internal/cache"
	"github.com/Sall-lah/store_product/internal/repository"
)

const (
	// ProductDetailTTL controls how long individual product models remain in Redis cache.
	ProductDetailTTL = 30 * time.Minute
	// ProductListTTL controls how long catalog queries remain cached before background refresh.
	ProductListTTL = 3 * time.Minute
)

// ProductService coordinates domain logic, cache-aside read strategy, and internal cache invalidation.
type ProductService struct {
	repo        *repository.ProductRepository
	variantRepo *repository.VariantRepository
	cacheClient *cache.Client
}

// NewProductService initializes a new ProductService with its repository and cache dependencies.
func NewProductService(
	repo *repository.ProductRepository,
	variantRepo *repository.VariantRepository,
	cacheClient *cache.Client,
) *ProductService {
	return &ProductService{
		repo:        repo,
		variantRepo: variantRepo,
		cacheClient: cacheClient,
	}
}

// ListProducts retrieves a paginated and filtered catalog list with cache-aside read acceleration.
func (s *ProductService) ListProducts(ctx context.Context, filter repository.ProductFilter) (*repository.PaginatedProducts, error) {
	cacheKey := s.computeListCacheKey(filter)

	// 1. Attempt read from Redis
	var cachedResult repository.PaginatedProducts
	if s.cacheClient.GetJSON(ctx, cacheKey, &cachedResult) {
		return &cachedResult, nil
	}

	// 2. Cache miss or Redis unavailable -> query Supabase via Prisma
	result, err := s.repo.ListProducts(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 3. Populate Redis asynchronously
	s.cacheClient.SetJSON(ctx, cacheKey, result, ProductListTTL)

	return result, nil
}

// AdminListProducts retrieves a paginated and filtered catalog list directly from PostgreSQL.
// Bypassing Redis cache for admin queries ensures backoffice operators see immediate data changes and inactive/draft items.
func (s *ProductService) AdminListProducts(ctx context.Context, filter repository.ProductFilter) (*repository.PaginatedProducts, error) {
	filter.IncludeInactive = true
	return s.repo.ListProducts(ctx, filter)
}

// AdminGetProductByID retrieves a product and all its variants (active & inactive) directly from PostgreSQL.
// Bypassing Redis cache avoids leaking inactive product data into public cache keys while providing full visibility to admins.
func (s *ProductService) AdminGetProductByID(ctx context.Context, id string) (*repository.ProductDTO, error) {
	return s.repo.GetAdminProductByID(ctx, id)
}

// GetProductByID retrieves a product and its variants by ID with cache-aside pattern.
func (s *ProductService) GetProductByID(ctx context.Context, id string) (*repository.ProductDTO, error) {
	cacheKey := cache.ProductDetailKey(id)

	var cached repository.ProductDTO
	if s.cacheClient.GetJSON(ctx, cacheKey, &cached) {
		return &cached, nil
	}

	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache both by ID and by Slug for fast dual-indexing
	s.cacheClient.SetJSON(ctx, cacheKey, product, ProductDetailTTL)
	s.cacheClient.SetJSON(ctx, cache.ProductSlugKey(product.Slug), product, ProductDetailTTL)

	return product, nil
}

// GetProductBySlug retrieves a product and its variants by Slug with cache-aside pattern.
func (s *ProductService) GetProductBySlug(ctx context.Context, slug string) (*repository.ProductDTO, error) {
	cacheKey := cache.ProductSlugKey(slug)

	var cached repository.ProductDTO
	if s.cacheClient.GetJSON(ctx, cacheKey, &cached) {
		return &cached, nil
	}

	product, err := s.repo.GetProductBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	s.cacheClient.SetJSON(ctx, cacheKey, product, ProductDetailTTL)
	s.cacheClient.SetJSON(ctx, cache.ProductDetailKey(product.ID), product, ProductDetailTTL)

	return product, nil
}

// CreateProduct persists a new product and invalidates existing catalog list caches.
func (s *ProductService) CreateProduct(ctx context.Context, input repository.CreateProductInput) (*repository.ProductDTO, error) {
	product, err := s.repo.CreateProduct(ctx, input)
	if err != nil {
		return nil, err
	}

	// Invalidate catalog list query caches so newly added product surfaces immediately
	s.invalidateListCaches(ctx)

	// Pre-warm the single product cache
	s.cacheClient.SetJSON(ctx, cache.ProductDetailKey(product.ID), product, ProductDetailTTL)
	s.cacheClient.SetJSON(ctx, cache.ProductSlugKey(product.Slug), product, ProductDetailTTL)

	return product, nil
}

// UpdateProduct updates product attributes and purges associated Redis cache entries.
func (s *ProductService) UpdateProduct(ctx context.Context, id string, input repository.UpdateProductInput) (*repository.ProductDTO, error) {
	// Retrieve previous state to identify if the slug changed
	existing, _ := s.repo.GetProductByID(ctx, id)

	updated, err := s.repo.UpdateProduct(ctx, id, input)
	if err != nil {
		return nil, err
	}

	// Invalidate previous cache keys
	s.cacheClient.Del(ctx, cache.ProductDetailKey(id))
	if existing != nil {
		s.cacheClient.Del(ctx, cache.ProductSlugKey(existing.Slug))
	}
	s.cacheClient.Del(ctx, cache.ProductSlugKey(updated.Slug))

	// Invalidate catalog lists
	s.invalidateListCaches(ctx)

	// Store refreshed product in cache
	s.cacheClient.SetJSON(ctx, cache.ProductDetailKey(updated.ID), updated, ProductDetailTTL)
	s.cacheClient.SetJSON(ctx, cache.ProductSlugKey(updated.Slug), updated, ProductDetailTTL)

	return updated, nil
}

// DeleteProduct deletes a product and removes all related cache records.
func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	existing, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteProduct(ctx, id); err != nil {
		return err
	}

	// Invalidate single product detail caches
	s.cacheClient.Del(ctx, cache.ProductDetailKey(id))
	if existing != nil {
		s.cacheClient.Del(ctx, cache.ProductSlugKey(existing.Slug))
	}

	// Invalidate catalog lists
	s.invalidateListCaches(ctx)

	return nil
}

// CreateVariant adds a variant to a product and purges parent product caches.
func (s *ProductService) CreateVariant(ctx context.Context, productID string, input repository.CreateVariantInput) (*repository.VariantDTO, error) {
	variant, err := s.variantRepo.CreateVariant(ctx, productID, input)
	if err != nil {
		return nil, err
	}

	s.invalidateProductCaches(ctx, productID)
	return variant, nil
}

// UpdateVariant updates a variant and purges parent product caches.
func (s *ProductService) UpdateVariant(ctx context.Context, productID, variantID string, input repository.CreateVariantInput) (*repository.VariantDTO, error) {
	variant, err := s.variantRepo.UpdateVariant(ctx, variantID, input)
	if err != nil {
		return nil, err
	}

	s.invalidateProductCaches(ctx, productID)
	return variant, nil
}

// DeleteVariant deletes a variant and purges parent product caches.
func (s *ProductService) DeleteVariant(ctx context.Context, productID, variantID string) error {
	if err := s.variantRepo.DeleteVariant(ctx, variantID); err != nil {
		return err
	}

	s.invalidateProductCaches(ctx, productID)
	return nil
}

// invalidateProductCaches purges detail and list caches when a product's variants change.
func (s *ProductService) invalidateProductCaches(ctx context.Context, productID string) {
	product, err := s.repo.GetProductByID(ctx, productID)
	if err == nil && product != nil {
		s.cacheClient.Del(ctx, cache.ProductDetailKey(product.ID), cache.ProductSlugKey(product.Slug))
	}
	s.invalidateListCaches(ctx)
}

// invalidateListCaches removes all cached catalog list queries.
func (s *ProductService) invalidateListCaches(ctx context.Context) {
	s.cacheClient.DelPattern(ctx, "product:list:*")
}

// computeListCacheKey creates a deterministic MD5 hash of the query filter parameters.
func (s *ProductService) computeListCacheKey(f repository.ProductFilter) string {
	raw := fmt.Sprintf("cat=%v|min=%v|max=%v|sz=%v|col=%v|q=%v|c=%v|l=%d",
		safeStr(f.Category),
		safeFloat(f.MinPrice),
		safeFloat(f.MaxPrice),
		safeStr(f.Size),
		safeStr(f.Color),
		safeStr(f.Search),
		safeStr(f.Cursor),
		f.Limit,
	)

	hash := md5.Sum([]byte(raw))
	return cache.ProductListKey(hex.EncodeToString(hash[:]))
}

func safeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func safeFloat(f *float64) string {
	if f == nil {
		return ""
	}
	return fmt.Sprintf("%.2f", *f)
}
