## 1. Middleware Role Casing Fix

- [x] 1.1 Update `RequireAdmin` in `internal/middleware/auth.go` to use `strings.EqualFold(strings.TrimSpace(role), "admin")`
- [x] 1.2 Add unit tests in `internal/middleware/middleware_test.go` covering `"ADMIN"`, `"Admin"`, `"admin"`, and whitespace trimming

## 2. OpenAPI 3.1 & Swagger Documentation Updates

- [x] 2.1 Update `docs/openapi.yaml` with Gateway Offloading auth explanations, case-insensitive `X-User-Role` definitions, and updated schema examples
- [x] 2.2 Synchronize `docs/openapi.json` with `docs/openapi.yaml`
- [x] 2.3 Verify `DocsHandler` unit tests in `internal/handler/docs_test.go` pass

## 3. Database Connection Configuration & Verification

- [x] 3.1 Update `.env.example` to document PostgreSQL connection strings for both local container networking and Supabase pooler setups
- [x] 3.2 Run full test suite across `store_product` (`go test ./...`) to verify all packages pass
