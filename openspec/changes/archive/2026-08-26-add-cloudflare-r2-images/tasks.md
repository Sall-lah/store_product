## 1. Schema & Database Models

- [x] 1.1 Update `prisma/schema.prisma` with `ProductImage` model and relations to `Product` and `ProductVariant`
- [x] 1.2 Generate Prisma Go client via `go run github.com/steebchen/prisma-client-go generate`
- [x] 1.3 Add `ProductImageDTO`, `CreateProductImageInput`, `UpdateProductImageInput`, and `PresignImageInput` domain models to `internal/repository/models.go`

## 2. Configuration & AWS S3 / R2 Client

- [x] 2.1 Add R2 configuration keys (`R2_ACCOUNT_ID`, `R2_ACCESS_KEY_ID`, `R2_SECRET_ACCESS_KEY`, `R2_BUCKET_NAME`, `R2_PUBLIC_BASE_URL`) to `internal/config/config.go` and `.env.example`
- [x] 2.2 Add AWS SDK for Go v2 S3 dependencies to `go.mod`
- [x] 2.3 Implement modular R2 storage client in `internal/storage/r2.go` supporting S3 PresignClient (`PresignPutObject`) and `DeleteObject`

## 3. Repository Layer

- [x] 3.1 Implement `ImageRepository` in `internal/repository/image.go` with Create, FindByProduct, FindByID, Update, Delete, and Primary Image re-indexing
- [x] 3.2 Update `ProductRepository` queries in `internal/repository/product.go` to include `ProductImage` relations and populate `primary_image_url` and `images` DTO fields

## 4. Service Layer & Cache Invalidation

- [x] 4.1 Implement `ImageService` in `internal/service/image.go` for presigned URL generation, validation (10MB, MIME types), image registration, and R2 deletion
- [x] 4.2 Integrate Redis cache invalidation for product detail and catalog queries upon image mutations

## 5. HTTP Handlers & Routing

- [x] 5.1 Implement `ImageHandler` in `internal/handler/image.go` with admin routes for presign, create, list, update, and delete
- [x] 5.2 Mount admin image routes in `internal/handler/router.go` under `/api/v1/admin/products/{id}/images`
- [x] 5.3 Update OpenAPI 3.1 schema in `internal/handler/docs.go` with all new image endpoints and schemas

## 6. Testing & Verification

- [x] 6.1 Add unit and integration tests for R2 storage presigning, image repository CRUD, and handler validation
- [x] 6.2 Verify end-to-end flow: presign URL generation, mock R2 PUT, image registration, image listing with variant tags, and cache eviction
