## Context

The store product catalog requires media management for base products and variant-specific colorways. To achieve zero-egress cost and high-speed global delivery, we integrate Cloudflare R2 object storage paired with a custom CDN domain. We use a presigned direct-to-R2 upload workflow so client uploads bypass server memory and network bandwidth limits.

## Goals / Non-Goals

**Goals:**
- Provide direct-to-R2 presigned upload URL generation with MIME type whitelist and 10MB size limits.
- Persist product image records in PostgreSQL via Prisma with support for product-level and variant-level associations.
- Support single primary thumbnail flags and sequence ordering (`sortOrder`).
- Implement automatic R2 object deletion when an image record is removed by admin.
- Return image URLs and galleries in public and admin product endpoints with Redis cache invalidation.
- Support dynamic edge image transformation via Cloudflare's `/cdn-cgi/image/` path formatting.

**Non-Goals:**
- On-server image processing/resizing libraries (e.g. libvips or imagemagick on the Go backend). All image resizing is offloaded to Cloudflare edge transformations.
- Multi-bucket or multi-cloud storage abstraction beyond S3-compatible R2 API.

## Decisions

### Decision 1: AWS SDK for Go v2 for Cloudflare R2
- **Choice**: Use `github.com/aws/aws-sdk-go-v2/service/s3` with `s3.NewPresignClient`.
- **Rationale**: Cloudflare R2 is fully S3-compatible. AWS SDK v2 is well-tested, robust, and supports custom endpoint routing (`https://<account_id>.r2.cloudflarestorage.com`).
- **Alternatives Considered**: Raw HTTP client with manual AWS Signature v4 (unnecessarily complex) or MinIO client SDK (heavier dependency).

### Decision 2: Direct-to-R2 Presigned Upload (3-Step Flow)
- **Choice**:
  1. Client calls `POST /api/v1/admin/products/{id}/images/presign` with filename and MIME type.
  2. Client uploads raw binary directly to R2 with HTTP PUT.
  3. Client calls `POST /api/v1/admin/products/{id}/images` to persist the database record.
- **Rationale**: Prevents multi-megabyte image payloads from consuming backend Go application RAM and network threads.

### Decision 3: Unified Product Gallery with Optional Variant Tagging
- **Choice**: Dedicated `ProductImage` table containing `productId` and nullable `variantId`.
- **Rationale**: Shopify-style unified gallery allows products to have base hero shots, lifestyle images, and variant-specific colorways in one coherent table.
- **Alternatives Considered**: Direct `imageUrl` columns on Product/Variant (cannot support multiple angles or galleries) or Many-to-Many join table (unneeded schema overhead).

### Decision 4: Edge Image Resizing via Cloudflare CDN
- **Choice**: Store only original master images in R2; clients request dynamic transformations via `/cdn-cgi/image/width=...,quality=.../products/...`.
- **Rationale**: Eliminates storage duplication and server-side CPU load while allowing infinite viewport adaptability.

## Risks / Trade-offs

- **[Risk] Orphaned R2 Uploads**: A client generates a presigned URL and uploads to R2, but fails to call the confirmation endpoint.
  - *Mitigation*: R2 object keys include product IDs and timestamps; optional Cloudflare bucket lifecycle rules can purge unreferenced uploads if needed.
- **[Risk] Multiple Primary Images**: Concurrent requests attempt to mark different images as primary.
  - *Mitigation*: Database transactions/Prisma queries atomically unset existing `isPrimary: true` for the product when a new primary image is designated.
- **[Risk] R2 Deletion Failure**: Database record is deleted but R2 `DeleteObject` fails due to transient network error.
  - *Mitigation*: Service logs error with structured logging; S3 client includes automatic retries with exponential backoff.

## Migration Plan

1. Update `prisma/schema.prisma` with `ProductImage` model and generate Prisma client.
2. Run database migration via `npx prisma db push` or migration scripts.
3. Configure environment variables in `.env` (`R2_ACCOUNT_ID`, `R2_ACCESS_KEY_ID`, `R2_SECRET_ACCESS_KEY`, `R2_BUCKET_NAME`, `R2_PUBLIC_BASE_URL`).
4. Deploy updated Go backend binary.

## Open Questions

None. All architecture choices, MIME whitelists (10MB limit), and edge resizing methods have been finalized.
