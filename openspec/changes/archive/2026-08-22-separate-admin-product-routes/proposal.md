## Why

Currently, admin product and variant mutations (`POST`, `PUT`, `DELETE`) share the `/api/v1/products` path with public catalog queries. This design prevents path-based API Gateway security policies, complicates rate-limiting rules, and makes it difficult for administrators to list or view inactive/draft products without contaminating public catalog cache keys. 

Separating admin functionality into a dedicated `/api/v1/admin/products` namespace isolates backoffice operations, simplifies API Gateway ACL routing, enables admin-specific query capabilities (including deactivated products and inventory metrics), and guarantees clean cache separation.

## What Changes

- **BREAKING**: Move admin product mutation endpoints (`POST /api/v1/products`, `PUT /api/v1/products/{id}`, `DELETE /api/v1/products/{id}`) to `/api/v1/admin/products`.
- **BREAKING**: Move admin variant mutation endpoints (`POST /api/v1/products/{id}/variants`, `PUT /api/v1/products/{id}/variants/{variantId}`, `DELETE /api/v1/products/{id}/variants/{variantId}`) to `/api/v1/admin/products/{id}/variants`.
- **NEW**: Add admin-specific catalog query endpoints (`GET /api/v1/admin/products` and `GET /api/v1/admin/products/{id}`) to allow admins to view all products including inactive, draft, and out-of-stock items.
- Public `/api/v1/products` routes become strictly read-only for customer storefront consumption (`GET` only, returning `isActive: true` records).
- Update OpenAPI/Swagger and Scalar documentation to document the isolated `/api/v1/admin/products` namespace.

## Capabilities

### New Capabilities
<!-- None: capabilities already exist in modified form -->

### Modified Capabilities
- `product-catalog-api`: Scope public endpoints strictly to read-only customer storefront operations (`GET /api/v1/products`, `GET /api/v1/products/{id}`, `GET /api/v1/products/slug/{slug}`).
- `product-variant-management`: Move admin CRUD operations for products and variants under the `/api/v1/admin/products` path prefix, and introduce admin listing capabilities that support inactive product retrieval.
- `gateway-auth-ratelimit`: Enforce upstream gateway role authorization (`X-User-Role: admin`) and admin rate limiting across the entire `/api/v1/admin/**` path hierarchy.

## Impact

- **Router**: `internal/handler/router.go` restructured into distinct public `/api/v1/products` and protected `/api/v1/admin/products` route groups.
- **Handlers**: `internal/handler/product.go` updated with admin list/detail handlers and routed to the new endpoints.
- **Repository**: `internal/repository/product.go` updated to allow optional inclusion of inactive products for admin queries (`IncludeInactive`).
- **Documentation**: `docs/openapi.yaml`, `docs/openapi.json`, and `docs/docs.go` updated to reflect the new paths and tags.
- **Tests**: `internal/handler/router_test.go` and related test suites updated to test the `/api/v1/admin/products` routes.
