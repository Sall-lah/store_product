## Why

The e-commerce store catalog currently only supports textual product data and pricing without media assets. Modern storefronts require high-quality imagery for products and variant-specific visuals (e.g. distinct colorways) to drive customer engagement and conversions. 

Cloudflare R2 provides S3-compatible, zero-egress fee object storage paired with Cloudflare CDN, making it cost-effective and performant for storing master images and serving responsive images at the edge.

## What Changes

- Add **Cloudflare R2 storage integration** using the AWS SDK v2 Go client for generating presigned PUT upload URLs and deleting objects.
- Introduce the **`ProductImage` database entity** in Prisma/PostgreSQL to support a unified product gallery with optional variant tagging, primary image flag, and sort order.
- Implement **direct-to-R2 presigned upload flow** (`POST /api/v1/admin/products/{id}/images/presign` and `POST /api/v1/admin/products/{id}/images`) to prevent server memory/bandwidth bottlenecks.
- Implement **admin gallery management operations** (List, update metadata, reorder, and delete images with automatic R2 object cleanup).
- Update **public catalog endpoints** (`GET /api/v1/products`, `GET /api/v1/products/{id}`, `GET /api/v1/products/slug/{slug}`) to return primary thumbnail URLs in listing cards and full gallery arrays with variant associations in product details.
- Leverage **Cloudflare Image Resizing at the CDN edge** (`/cdn-cgi/image/width=...`) so clients can dynamically fetch high-res zoom or low-res thumbnails from a single uploaded master image.
- Enforce **file constraints**: Allowed MIME types (`image/jpeg`, `image/png`, `image/webp`, `image/avif`, `image/gif`) and max size of 10MB.
- Maintain **Redis cache invalidation** on image mutations and update OpenAPI 3.1 documentation.

## Capabilities

### New Capabilities
- `product-image-management`: Presigned upload URL generation, file type/size validation, image registration, metadata updates, sort ordering, and R2 object deletion.

### Modified Capabilities
- `product-catalog-api`: Public product listings and detail endpoints include `primary_image_url` and full gallery `images` array with variant associations.

## Impact

- **Database**: New `ProductImage` model in `prisma/schema.prisma` with foreign keys to `Product` and `ProductVariant`.
- **Backend**: New storage client package (`internal/storage/r2.go`), image repository (`internal/repository/image.go`), image service (`internal/service/image.go`), and image HTTP handlers (`internal/handler/image.go`).
- **Configuration**: New environment variables for R2 credentials (`R2_ACCOUNT_ID`, `R2_ACCESS_KEY_ID`, `R2_SECRET_ACCESS_KEY`, `R2_BUCKET_NAME`, `R2_PUBLIC_BASE_URL`).
- **Dependencies**: AWS SDK for Go v2 S3 client (`github.com/aws/aws-sdk-go-v2/service/s3` and config/credentials).
- **APIs**: New admin image routes under `/api/v1/admin/products/{id}/images` and updated responses for public `/api/v1/products`.
- **Docs**: OpenAPI 3.1 spec update in `internal/handler/docs.go`.
