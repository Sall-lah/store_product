## 1. Repository & Service Layer Updates

- [x] 1.1 Extend `ProductFilter` in `internal/repository/product.go` to support `IncludeInactive` and `IsActive` filtering
- [x] 1.2 Add admin product listing and single product retrieval methods in `internal/repository/product.go` and `internal/service/product.go`

## 2. Handler & Router Restructuring

- [x] 2.1 Add `AdminListProducts` and `AdminGetProductByID` handlers in `internal/handler/product.go`
- [x] 2.2 Restructure router groups in `internal/handler/router.go` into public `/api/v1/products` and admin `/api/v1/admin/products`
- [x] 2.3 Apply `RequireAdmin` and admin rate limiting exclusively to `/api/v1/admin/products`

## 3. OpenAPI Documentation Updates

- [x] 3.1 Update `docs/openapi.yaml` and `docs/openapi.json` to define `/api/v1/admin/products` endpoints and security schemes
- [x] 3.2 Update `docs/docs.go` with regenerated Swagger / Scalar metadata

## 4. Testing & Verification

- [x] 4.1 Update router and handler tests in `internal/handler/router_test.go` and `internal/handler/docs_test.go`
- [x] 4.2 Run test suite `go test -v ./...` to verify authorization, route splitting, and public catalog read safety
