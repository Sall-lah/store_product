## 1. Repository Soft-Delete Implementation

- [x] 1.1 Update `ProductRepository.DeleteProduct` in `internal/repository/product.go` to set `isActive = false` on product and cascade `isActive = false` to all associated variants
- [x] 1.2 Update `VariantRepository.DeleteVariant` in `internal/repository/variant.go` to set `isActive = false` on the variant record
- [x] 1.3 Update `ProductRepository.GetProductByID` and `GetProductBySlug` in `internal/repository/product.go` to return `ErrProductNotFound` if `product.IsActive == false`

## 2. Service & Handler Alignment

- [x] 2.1 Verify `ProductService.DeleteProduct` and `ProductService.DeleteVariant` in `internal/service/product.go` purge Redis cache keys on soft delete
- [x] 2.2 Verify admin endpoints (`AdminListProducts` and `AdminGetProductByID`) retain full access to soft-deleted items

## 3. Testing & Quality Assurance

- [x] 3.1 Update repository and service tests to validate soft-delete behavior instead of row removal
- [x] 3.2 Update handler tests to verify `DELETE /api/v1/admin/products/:id` returns 204 and subsequent public lookup returns 404
- [x] 3.3 Run `go test ./...` across all packages to verify passing tests and zero regressions
