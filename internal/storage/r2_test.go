package storage

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Sall-lah/store_product/internal/config"
)

func TestStorageValidationAndKeys(t *testing.T) {
	cfg := &config.Config{
		R2BucketName:    "store-products-test",
		R2PublicBaseURL: "https://cdn.mystore.com",
	}

	client, err := NewR2StorageClient(context.Background(), cfg)
	if err != nil {
		t.Fatalf("unexpected error creating storage client: %v", err)
	}

	t.Run("ContentType validation allows whitelist and rejects others", func(t *testing.T) {
		validTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/webp", "image/avif", "image/gif", "IMAGE/PNG "}
		for _, vt := range validTypes {
			if !IsValidContentType(vt) {
				t.Errorf("expected %q to be valid content type", vt)
			}
		}

		invalidTypes := []string{"application/pdf", "text/plain", "image/svg+xml", "image/bmp", "video/mp4", ""}
		for _, it := range invalidTypes {
			if IsValidContentType(it) {
				t.Errorf("expected %q to be rejected as invalid content type", it)
			}
		}
	})

	t.Run("GenerateR2Key creates product-scoped path", func(t *testing.T) {
		key := client.GenerateR2Key("prod-123", nil, "photo.jpg")
		if !strings.HasPrefix(key, "products/prod-123/images/") {
			t.Errorf("expected product image prefix, got: %s", key)
		}
		if !strings.HasSuffix(key, ".jpg") {
			t.Errorf("expected .jpg extension, got: %s", key)
		}
	})

	t.Run("GenerateR2Key creates variant-scoped path", func(t *testing.T) {
		varID := "var-456"
		key := client.GenerateR2Key("prod-123", &varID, "front.webp")
		if !strings.HasPrefix(key, "products/prod-123/variants/var-456/") {
			t.Errorf("expected variant image prefix, got: %s", key)
		}
		if !strings.HasSuffix(key, ".webp") {
			t.Errorf("expected .webp extension, got: %s", key)
		}
	})

	t.Run("BuildPublicURL formats CDN URL correctly", func(t *testing.T) {
		url := client.BuildPublicURL("products/prod-123/images/img.webp")
		expected := "https://cdn.mystore.com/products/prod-123/images/img.webp"
		if url != expected {
			t.Errorf("expected %q, got %q", expected, url)
		}
	})

	t.Run("GeneratePresignedUploadURL mock mode", func(t *testing.T) {
		url, err := client.GeneratePresignedUploadURL(context.Background(), "test-key.jpg", "image/jpeg", 15*time.Minute)
		if err != nil {
			t.Fatalf("unexpected error in mock presign: %v", err)
		}
		if !strings.Contains(url, "test-key.jpg") {
			t.Errorf("expected mock url to contain key, got: %s", url)
		}
	})

	t.Run("DeleteObject mock mode does not panic", func(t *testing.T) {
		err := client.DeleteObject(context.Background(), "test-key.jpg")
		if err != nil {
			t.Errorf("unexpected error deleting object: %v", err)
		}
	})
}
