## Context

The `store_product` Go microservice is built using the Chi HTTP router, Prisma Client Go, and Redis. While the endpoints are implemented and verified, client developers and API gateways need a standardized OpenAPI 3.1.0 specification and interactive documentation UI to explore schemas, test endpoints, and view rate limit/authentication requirements.

## Goals / Non-Goals

**Goals:**
- Provide a comprehensive, strictly compliant **OpenAPI 3.1.0** document in both YAML (`docs/openapi.yaml`) and JSON (`docs/openapi.json`).
- Embed the OpenAPI 3.1 specification files directly into the compiled Go binary using `embed.FS` to prevent runtime missing-file errors in Docker/production deployments.
- Serve interactive **Swagger UI** (`/swagger` and `/docs`) and **Scalar UI** (`/docs/scalar` or unified `/docs`) with dark/light mode and dynamic "Try it out" capabilities.
- Fully document security schemes (`ApiKeyAuth` on `X-User-Role` and `X-User-Id`) and sliding window rate limiting response headers (`X-RateLimit-*`).
- Provide unit and integration tests verifying doc endpoint accessibility and content-type headers.

**Non-Goals:**
- Runtime request payload validation against OpenAPI schema (Prisma and Go structs already enforce strict type safety).
- Generating Go handler code from OpenAPI spec (the service is already cleanly implemented).

## Decisions

### 1. Spec-First Declarative YAML with Go `embed.FS`
- **Choice**: Handcraft a complete, modular OpenAPI 3.1.0 specification in `docs/openapi.yaml` and generate/embed it via `embed.FS` in `internal/handler/docs.go`.
- **Why**: Existing Go annotation tools like `swaggo/swag` are constrained to Swagger 2.0 / OpenAPI 3.0 and do not support OpenAPI 3.1 features (JSON Schema 2020-12, webhooks, complex `anyOf`/`oneOf`, `null` type unions). Handcrafting gives 100% control over examples, descriptions, and type precision without cluttering handler source code with 30-line comments.
- **Alternatives Considered**:
  - *Swaggo annotations*: Rejected due to Swagger 2.0 limitation and heavy code clutter.
  - *External Redoc/Swagger container*: Rejected to keep the microservice self-contained in a single executable.

### 2. Dual UI Rendering (Swagger UI & Scalar UI)
- **Choice**: Serve both **Swagger UI** (at `/swagger`) and modern **Scalar UI** (at `/docs`).
- **Why**: Swagger UI is universally recognized by enterprise tooling, while Scalar provides an ultra-modern developer experience with dark mode, instantaneous search, and auto-generated multi-language code snippets (cURL, Go, JS, Python). Both render client-side via lightweight HTML shells backed by CDN assets pointing to `/openapi.json`.
- **Alternatives Considered**:
  - *Swagger UI only*: Misses out on modern code snippet generation and OpenAPI 3.1 interactive features.

### 3. API Gateway Authentication in Documentation Console
- **Choice**: Define `ApiKeyAuth` security schemes for `X-User-Role` and `X-User-Id` in the OpenAPI `components.securitySchemes`.
- **Why**: Allows developers and testers to click "Authorize" in Swagger UI / Scalar UI, enter `admin` and `usr_admin_1`, and execute protected mutation endpoints (`POST`, `PUT`, `DELETE`) directly from the browser.

## Risks / Trade-offs

- **[Risk] Spec Drift**: Modifying Go handlers in the future without updating `docs/openapi.yaml`.
  - *Mitigation*: Comprehensive integration tests (`internal/handler/docs_test.go`) ensure that every registered route in Chi has a corresponding path in the OpenAPI spec.
- **[Risk] CDN Availability for UI**: HTML shells rely on unpkg/cdnjs for Swagger UI and Scalar assets.
  - *Mitigation*: Raw `/openapi.json` and `/openapi.yaml` endpoints are always available offline; UI shells use version-pinned, highly available CDNs with fallback links.

## File Structure

```
store_product/
├── docs/
│   ├── openapi.yaml       # Master OpenAPI 3.1.0 specification
│   └── openapi.json       # JSON formatted specification
├── internal/
│   └── handler/
│       ├── docs.go        # Embedded filesystem & HTTP handlers for /docs, /swagger, /openapi.json
│       ├── docs_test.go   # Integration tests for doc endpoints
│       └── router.go      # Route registration for documentation
```
