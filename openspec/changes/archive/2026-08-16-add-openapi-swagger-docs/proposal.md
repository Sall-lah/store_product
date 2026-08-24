## Why

As the `store_product` microservice expands and integrates with API Gateways, frontend applications, and external client developers, an interactive, standardized, and machine-readable API contract is essential. Providing an OpenAPI 3.1.0 specification with interactive Swagger UI and modern Scalar documentation ensures effortless API exploration, accurate client generation, and clear visibility into authentication and rate-limiting behaviors.

## What Changes

- Author an **OpenAPI 3.1.0** specification (`docs/openapi.yaml` and `docs/openapi.json`) covering all catalog queries, cursor pagination metadata, filtering parameters, admin product/variant CRUD mutations, and error response schemas.
- Document Gateway Header security schemes (`X-User-Role`, `X-User-Id`) and sliding window rate limiting response headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`).
- Embed the OpenAPI 3.1 specification and UI assets into the Go binary using Go standard `embed.FS`.
- Expose HTTP documentation endpoints in Chi router:
  - `GET /docs` & `GET /swagger`: Interactive documentation UI (Swagger UI and Scalar UI).
  - `GET /openapi.json` & `GET /openapi.yaml`: Raw machine-readable OpenAPI 3.1 specifications.
- Add comprehensive integration tests verifying doc endpoints return valid JSON/YAML and HTML payloads with correct HTTP 200 OK statuses and headers.

## Capabilities

### New Capabilities
- `openapi-documentation`: Full OpenAPI 3.1.0 specification describing all endpoints, models, parameters, error formats, and security schemes for the product microservice.
- `swagger-scalar-ui`: Embedded HTTP endpoints (`/docs`, `/swagger`, `/openapi.json`, `/openapi.yaml`) serving interactive API documentation and raw schemas.

### Modified Capabilities
<!-- No existing capabilities modified; this adds documentation capabilities to the service. -->

## Impact

- **New Files**: `docs/openapi.yaml`, `docs/openapi.json`, and `internal/handler/docs.go` for serving embedded documentation.
- **Router Changes**: Registers `/docs`, `/swagger`, `/openapi.json`, and `/openapi.yaml` in `internal/handler/router.go`.
- **Dependencies**: No external runtime CGO dependencies needed; uses Go standard `embed.FS` and CDN-backed Swagger UI / Scalar UI renderer.
- **Testing**: Adds `internal/handler/docs_test.go` to test doc routes and schema accessibility.
