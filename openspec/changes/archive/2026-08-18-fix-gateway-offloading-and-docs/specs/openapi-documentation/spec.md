## MODIFIED Requirements

### Requirement: Comprehensive OpenAPI 3.1.0 Specification
The system SHALL provide an OpenAPI 3.1.0 specification defining all endpoints, query parameters, request bodies, response models, error structures, rate limit headers, and security schemes for the `store_product` microservice.

#### Scenario: OpenAPI specification includes all product endpoints
- **WHEN** a developer inspects the OpenAPI 3.1.0 specification
- **THEN** it SHALL define schemas and routes for `GET /health`, `GET /api/v1/products`, `GET /api/v1/products/{id}`, `GET /api/v1/products/slug/{slug}`, `POST /api/v1/products`, `PUT /api/v1/products/{id}`, `DELETE /api/v1/products/{id}`, `POST /api/v1/products/{id}/variants`, `PUT /api/v1/products/{id}/variants/{variantId}`, and `DELETE /api/v1/products/{id}/variants/{variantId}`

#### Scenario: OpenAPI specification documents Keyset Cursor Pagination
- **WHEN** client queries the `/api/v1/products` catalog endpoint schema
- **THEN** the spec SHALL describe the `cursor`, `limit`, `category`, `min_price`, `max_price`, and `search` query parameters, along with the `pageInfo` response object containing `hasMore` and `nextCursor`

#### Scenario: OpenAPI specification documents Security Schemes
- **WHEN** an authenticated admin endpoint is examined in the specification
- **THEN** it SHALL declare `ApiKeyAuth` requirements specifying `X-User-Role` and `X-User-Id` header parameters and document the Gateway Offloading header injection mechanism with case-insensitive `admin` role support

#### Scenario: OpenAPI specification documents Rate Limit Headers
- **WHEN** inspect response header definitions for public, search, and admin endpoints
- **THEN** the spec SHALL include `X-RateLimit-Limit`, `X-RateLimit-Remaining`, and `X-RateLimit-Reset` headers
