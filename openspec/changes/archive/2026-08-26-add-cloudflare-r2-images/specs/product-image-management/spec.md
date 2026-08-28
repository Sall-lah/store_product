## ADDED Requirements

### Requirement: Presigned Upload URL Generation
The service SHALL provide an admin endpoint `POST /api/v1/admin/products/:id/images/presign` that validates file parameters and generates a Cloudflare R2 presigned PUT URL with an expiry duration of 15 minutes (900 seconds).

#### Scenario: Successful presign request with valid MIME and size
- **WHEN** an authenticated admin requests a presigned URL with a valid MIME type (`image/jpeg`, `image/png`, `image/webp`, `image/avif`, `image/gif`) and file size <= 10MB (10485760 bytes)
- **THEN** the system returns HTTP 200 with `upload_url`, `public_url`, `r2_key`, and `expires_in_seconds`.

#### Scenario: Rejection of invalid MIME type
- **WHEN** an admin requests a presigned URL with an unapproved MIME type (e.g. `application/pdf` or `image/bmp`)
- **THEN** the system returns HTTP 400 Bad Request with error code `invalid_content_type`.

#### Scenario: Rejection of oversized file request
- **WHEN** an admin requests a presigned URL with file size exceeding 10MB
- **THEN** the system returns HTTP 400 Bad Request with error code `file_too_large`.

### Requirement: Product Image Registration and Association
The service SHALL provide an admin endpoint `POST /api/v1/admin/products/:id/images` to confirm and persist an uploaded image in PostgreSQL linked to the product and optional variant.

#### Scenario: Register general product image
- **WHEN** an admin submits image confirmation with `url`, `r2_key`, `is_primary: true`, and `sort_order: 0` without a `variant_id`
- **THEN** the system persists the image record, marks previous primary images as `is_primary: false`, invalidates the product Redis cache, and returns HTTP 201 Created.

#### Scenario: Register variant-specific image
- **WHEN** an admin submits image confirmation with a valid `variant_id` belonging to the product
- **THEN** the system persists the image linked to both product and variant, invalidates Redis cache, and returns HTTP 201 Created.

### Requirement: Admin Image Listing and Filtering
The service SHALL provide `GET /api/v1/admin/products/:id/images` returning all images attached to a product, ordered by `sort_order ASC`, with optional filtering by `variant_id`.

#### Scenario: List all images for product
- **WHEN** an admin requests `GET /api/v1/admin/products/:id/images`
- **THEN** the system returns all general and variant-tagged images for the product.

### Requirement: Image Metadata Modification
The service SHALL provide `PUT /api/v1/admin/products/:id/images/:imageId` allowing modification of `alt_text`, `is_primary`, `sort_order`, and `variant_id`.

#### Scenario: Update image sort order and primary flag
- **WHEN** an admin updates an image setting `is_primary: true` and `sort_order: 1`
- **THEN** the system updates the record, unsets previous primary flag if changed, invalidates Redis cache, and returns HTTP 200 OK.

### Requirement: Image Deletion and R2 Storage Cleanup
The service SHALL provide `DELETE /api/v1/admin/products/:id/images/:imageId` that deletes the record from PostgreSQL and deletes the object from the Cloudflare R2 bucket.

#### Scenario: Successful image deletion
- **WHEN** an admin sends a DELETE request for an existing image ID
- **THEN** the system deletes the database row, invokes S3 `DeleteObject` on R2, invalidates Redis cache, and returns HTTP 200 OK or 204 No Content.
