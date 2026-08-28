package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Sall-lah/store_product/internal/config"
	"github.com/Sall-lah/store_product/internal/repository"
	"github.com/Sall-lah/store_product/internal/storage"
)

func TestImageServiceValidation(t *testing.T) {
	cfg := &config.Config{
		R2BucketName:    "store-products-test",
		R2PublicBaseURL: "https://cdn.mystore.com",
	}

	storageClient, err := storage.NewR2StorageClient(context.Background(), cfg)
	if err != nil {
		t.Fatalf("unexpected error creating storage client: %v", err)
	}

	imgSvc := NewImageService(&repository.ImageRepository{}, &repository.ProductRepository{}, storageClient, nil)

	t.Run("CreateImage rejects empty URL or empty r2Key", func(t *testing.T) {
		_, err := imgSvc.CreateImage(context.Background(), "prod-1", repository.CreateProductImageInput{
			URL:   "",
			R2Key: "",
		})
		if !errors.Is(err, ErrMissingImagePayload) {
			t.Errorf("expected ErrMissingImagePayload, got %v", err)
		}
	})
}
