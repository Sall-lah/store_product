## Why

In an API Gateway Offloading architecture, `store_gateway` verifies JWT tokens and injects verified upstream headers (`X-User-Role`, `X-User-Id`) to backend services. Because authentication microservices (`store_auth`) frequently emit uppercase role claims (such as `"ADMIN"`), `store_product`'s strict case check (`role != "admin"`) results in false `403 Forbidden` responses. Additionally, the OpenAPI 3.1.0 and Swagger documentation must accurately document this gateway offloading pattern and case-insensitive role handling to prevent integration mismatches.

## What Changes

- **Case-Insensitive Role Validation**: Update `RequireAdmin` middleware to accept case-insensitive `X-User-Role` values (`"admin"`, `"ADMIN"`, `"Admin"`) with whitespace trimming.
- **Unit Test Coverage**: Add test cases for uppercase, mixed-case, and whitespace-padded role headers.
- **OpenAPI 3.1.0 & Swagger Spec Update**: Clarify the Gateway Offloading identity injection headers and case-insensitivity in `docs/openapi.yaml` and `docs/openapi.json`.
- **Database Connection Documentation**: Clarify PostgreSQL connection parameters in `.env.example` for both local container networking and Supabase cloud poolers.

## Capabilities

### Modified Capabilities
- `gateway-auth-ratelimit`: Update role requirement to validate `X-User-Role` case-insensitively for admin endpoints.
- `openapi-documentation`: Update security scheme descriptions and documentation to specify Gateway Offloading header expectations and case-insensitive role values.

## Impact

- `internal/middleware/auth.go`: Middleware role comparison logic.
- `internal/middleware/middleware_test.go`: Auth test assertions.
- `docs/openapi.yaml` & `docs/openapi.json`: OpenAPI 3.1 schema definitions and security documentation.
- `.env.example`: Configuration guidance for PostgreSQL connection URLs.
