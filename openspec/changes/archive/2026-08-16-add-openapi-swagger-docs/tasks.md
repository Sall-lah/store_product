## 1. OpenAPI 3.1.0 Specification Creation

- [x] 1.1 Create `docs/openapi.yaml` defining OpenAPI 3.1.0 root metadata, servers, info, and `ApiKeyAuth` security schemes (`X-User-Role`, `X-User-Id`).
- [x] 1.2 Define complete schemas for Product, ProductVariant, Cursor PageInfo, ErrorResponse, and rate-limiting headers in `docs/openapi.yaml`.
- [x] 1.3 Document all public catalog routes (`/health`, `GET /api/v1/products`, `GET /api/v1/products/{id}`, `GET /api/v1/products/slug/{slug}`) with keyset pagination parameters and filters.
- [x] 1.4 Document all admin product and variant mutation routes (`POST /api/v1/products`, `PUT /api/v1/products/{id}`, `DELETE /api/v1/products/{id}`, `POST /api/v1/products/{id}/variants`, `PUT /api/v1/products/{id}/variants/{variantId}`, `DELETE /api/v1/products/{id}/variants/{variantId}`).
- [x] 1.5 Provide `docs/openapi.json` formatted specification.

## 2. Documentation Handler & Go Embedded Filesystem

- [x] 2.1 Implement `internal/handler/docs.go` with Go standard `embed.FS` to bundle OpenAPI 3.1 specifications into the compiled binary.
- [x] 2.2 Implement HTTP handler to serve raw OpenAPI specifications (`GET /openapi.json` and `GET /openapi.yaml`).
- [x] 2.3 Implement Swagger UI HTML renderer for `GET /swagger` with interactive console and header authorization support.
- [x] 2.4 Implement Scalar UI HTML renderer for `GET /docs` with modern themes and auto-generated multi-language code snippets.

## 3. Route Registration & Microservice Integration

- [x] 3.1 Mount documentation routes in `internal/handler/router.go` under public access.
- [x] 3.2 Ensure proper Content-Type headers (`text/html; charset=utf-8`, `application/json`, `application/yaml`).

## 4. Automated Testing & Verification

- [x] 4.1 Create `internal/handler/docs_test.go` to test `/docs`, `/swagger`, `/openapi.json`, and `/openapi.yaml` endpoints.
- [x] 4.2 Run unit and integration tests (`go test ./...`) to verify 100% test pass.
- [x] 4.3 Validate OpenAPI 3.1 schema correctness and interactive UI rendering against the running server.
