package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sall-lah/store_product/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// AllowedContentTypes defines the whitelisted image MIME types acceptable for product catalog media.
var AllowedContentTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/webp": true,
	"image/avif": true,
	"image/gif":  true,
}

// MaxUploadSizeBytes defines the strict 10MB file size ceiling for Cloudflare R2 image presigning.
const MaxUploadSizeBytes int64 = 10 * 1024 * 1024

// R2StorageClient abstracts interactions with Cloudflare R2 object storage via AWS S3 v2 client.
type R2StorageClient struct {
	s3Client       *s3.Client
	presignClient  *s3.PresignClient
	bucket         string
	publicBaseURL  string
}

// NewR2StorageClient initializes an S3 client configured specifically for Cloudflare R2 endpoints.
// Custom endpoint and static credentials are required because R2 uses account-specific regional routing.
func NewR2StorageClient(ctx context.Context, cfg *config.Config) (*R2StorageClient, error) {
	if cfg.R2AccountID == "" || cfg.R2AccessKeyID == "" || cfg.R2SecretAccessKey == "" {
		// Return client with nil internals if unconfigured (enables mock/dev fallback gracefully)
		return &R2StorageClient{
			bucket:        cfg.R2BucketName,
			publicBaseURL: strings.TrimRight(cfg.R2PublicBaseURL, "/"),
		}, nil
	}

	r2Endpoint := cfg.R2Endpoint()
	staticCreds := credentials.NewStaticCredentialsProvider(
		cfg.R2AccessKeyID,
		cfg.R2SecretAccessKey,
		"",
	)

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               r2Endpoint,
			SigningRegion:     "auto",
			HostnameImmutable: true,
		}, nil
	})

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithCredentialsProvider(staticCreds),
		awsconfig.WithRegion("auto"),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS/R2 configuration: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	presignClient := s3.NewPresignClient(s3Client)

	return &R2StorageClient{
		s3Client:      s3Client,
		presignClient: presignClient,
		bucket:        cfg.R2BucketName,
		publicBaseURL: strings.TrimRight(cfg.R2PublicBaseURL, "/"),
	}, nil
}

// GeneratePresignedUploadURL creates a short-lived HTTP PUT URL enabling direct client-to-R2 upload.
// Direct uploading prevents large binary streams from saturating backend application RAM and bandwidth.
func (c *R2StorageClient) GeneratePresignedUploadURL(ctx context.Context, key, contentType string, lifetime time.Duration) (string, error) {
	if c.presignClient == nil {
		// Mock endpoint for local testing environments without R2 secrets configured
		return fmt.Sprintf("https://mock-upload.local/%s/%s", c.bucket, key), nil
	}

	req, err := c.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(lifetime))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	return req.URL, nil
}

// DeleteObject removes an image file from the Cloudflare R2 bucket.
// Immediate object deletion prevents orphaned media from lingering and accumulating storage space.
func (c *R2StorageClient) DeleteObject(ctx context.Context, key string) error {
	if c.s3Client == nil || key == "" {
		return nil
	}

	_, err := c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete R2 object for key %q: %w", key, err)
	}

	return nil
}

// BuildPublicURL formats the canonical CDN delivery URL for a given object key.
// Unifying URL formatting here ensures consistent CDN caching behavior across the frontend.
func (c *R2StorageClient) BuildPublicURL(key string) string {
	cleanKey := strings.TrimLeft(key, "/")
	if c.publicBaseURL == "" {
		return fmt.Sprintf("/%s", cleanKey)
	}
	return fmt.Sprintf("%s/%s", c.publicBaseURL, cleanKey)
}

// GenerateR2Key constructs an isolated, collision-free object key partitioned by product and variant ID.
// Deterministic partitioning organizes bucket structures and prevents concurrent upload filename conflicts.
func (c *R2StorageClient) GenerateR2Key(productID string, variantID *string, originalFileName string) string {
	ext := filepath.Ext(originalFileName)
	if ext == "" {
		ext = ".webp"
	}

	randomToken := make([]byte, 8)
	_, _ = rand.Read(randomToken)
	suffix := hex.EncodeToString(randomToken)

	timestamp := time.Now().Unix()

	if variantID != nil && *variantID != "" {
		return fmt.Sprintf("products/%s/variants/%s/%d_%s%s", productID, *variantID, timestamp, suffix, ext)
	}

	return fmt.Sprintf("products/%s/images/%d_%s%s", productID, timestamp, suffix, ext)
}

// IsValidContentType verifies whether an uploaded file MIME matches allowed image formats.
func IsValidContentType(contentType string) bool {
	cleanType := strings.ToLower(strings.TrimSpace(contentType))
	return AllowedContentTypes[cleanType]
}
