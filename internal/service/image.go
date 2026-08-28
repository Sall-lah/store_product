package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Sall-lah/store_product/internal/cache"
	"github.com/Sall-lah/store_product/internal/repository"
	"github.com/Sall-lah/store_product/internal/storage"
)

var (
	// ErrInvalidContentType indicates an uploaded file MIME does not match approved image types.
	ErrInvalidContentType = errors.New("invalid or unsupported image content type")
	// ErrMissingFileName indicates the presign request omitted a target filename.
	ErrMissingFileName = errors.New("file name is required")
	// ErrMissingImagePayload indicates the image confirmation omitted necessary URLs or keys.
	ErrMissingImagePayload = errors.New("image url and r2_key are required")
)

const (
	// PresignExpiryDuration defines the 15-minute validity window for client PUT upload URLs.
	PresignExpiryDuration = 15 * time.Minute
)

// ImageService coordinates business logic, presigned URL generation, persistence, and storage deletion.
type ImageService struct {
	repo        *repository.ImageRepository
	productRepo *repository.ProductRepository
	storage     *storage.R2StorageClient
	cacheClient *cache.Client
}

// NewImageService creates a new instance of ImageService.
func NewImageService(
	repo *repository.ImageRepository,
	productRepo *repository.ProductRepository,
	storageClient *storage.R2StorageClient,
	cacheClient *cache.Client,
) *ImageService {
	return &ImageService{
		repo:        repo,
		productRepo: productRepo,
		storage:     storageClient,
		cacheClient: cacheClient,
	}
}

// GeneratePresignedURL validates request parameters and returns a direct-to-R2 signed PUT URL.
// Verifying MIME types upfront prevents unapproved binary formats from being uploaded to the bucket.
func (s *ImageService) GeneratePresignedURL(ctx context.Context, productID string, input repository.PresignImageInput) (*repository.PresignImageResponse, error) {
	// Ensure parent product exists
	product, err := s.productRepo.GetAdminProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	cleanFileName := strings.TrimSpace(input.FileName)
	if cleanFileName == "" {
		return nil, ErrMissingFileName
	}

	cleanContentType := strings.ToLower(strings.TrimSpace(input.ContentType))
	if !storage.IsValidContentType(cleanContentType) {
		return nil, ErrInvalidContentType
	}

	r2Key := s.storage.GenerateR2Key(product.ID, input.VariantID, cleanFileName)

	uploadURL, err := s.storage.GeneratePresignedUploadURL(ctx, r2Key, cleanContentType, PresignExpiryDuration)
	if err != nil {
		return nil, err
	}

	publicURL := s.storage.BuildPublicURL(r2Key)

	return &repository.PresignImageResponse{
		UploadURL:        uploadURL,
		PublicURL:        publicURL,
		R2Key:            r2Key,
		ExpiresInSeconds: int64(PresignExpiryDuration.Seconds()),
	}, nil
}

// CreateImage registers an uploaded image in PostgreSQL and purges affected Redis product caches.
func (s *ImageService) CreateImage(ctx context.Context, productID string, input repository.CreateProductImageInput) (*repository.ProductImageDTO, error) {
	if strings.TrimSpace(input.URL) == "" || strings.TrimSpace(input.R2Key) == "" {
		return nil, ErrMissingImagePayload
	}

	product, err := s.productRepo.GetAdminProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	created, err := s.repo.Create(ctx, productID, input)
	if err != nil {
		return nil, err
	}

	// Invalidate caches so updated image gallery is immediately reflected
	s.invalidateProductCaches(ctx, product.ID, product.Slug)

	return created, nil
}

// ListImages retrieves all gallery assets for a given product with optional variant filtering.
func (s *ImageService) ListImages(ctx context.Context, productID string, variantID *string) ([]repository.ProductImageDTO, error) {
	_, err := s.productRepo.GetAdminProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return s.repo.ListByProduct(ctx, productID, variantID)
}

// UpdateImage updates image metadata and evicts stale cached product queries.
func (s *ImageService) UpdateImage(ctx context.Context, id string, input repository.UpdateProductImageInput) (*repository.ProductImageDTO, error) {
	updated, err := s.repo.Update(ctx, id, input)
	if err != nil {
		return nil, err
	}

	product, _ := s.productRepo.GetAdminProductByID(ctx, updated.ProductID)
	if product != nil {
		s.invalidateProductCaches(ctx, product.ID, product.Slug)
	}

	return updated, nil
}

// DeleteImage deletes the DB record, purges the R2 object, and evicts stale caches.
// Removing both database rows and bucket objects prevents storage leaks and orphaned assets.
func (s *ImageService) DeleteImage(ctx context.Context, id string) error {
	deleted, err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Clean up from R2 bucket
	_ = s.storage.DeleteObject(ctx, deleted.R2Key)

	product, _ := s.productRepo.GetAdminProductByID(ctx, deleted.ProductID)
	if product != nil {
		s.invalidateProductCaches(ctx, product.ID, product.Slug)
	}

	return nil
}

// invalidateProductCaches purges single-product cache keys and all paginated catalog list queries.
func (s *ImageService) invalidateProductCaches(ctx context.Context, productID, slug string) {
	s.cacheClient.Del(ctx, cache.ProductDetailKey(productID))
	if slug != "" {
		s.cacheClient.Del(ctx, cache.ProductSlugKey(slug))
	}
	s.cacheClient.DelPattern(ctx, "product:list:*")
}
