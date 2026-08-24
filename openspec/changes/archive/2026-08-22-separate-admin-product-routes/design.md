## Context

Currently, the `store_product` service routes all product interactions through `/api/v1/products`. Public users access `GET` methods, while administrative actions (`POST`, `PUT`, `DELETE`) are nested within the same route group protected by the `RequireAdmin` middleware. 

This coupling creates challenges for path-based API Gateway routing, complicates rate limiting, and restricts administrative visibility into draft, deactivated, or out-of-stock products because the public repository queries hardcode `isActive = true`.

## Goals / Non-Goals

**Goals:**
- Separate public catalog operations (`GET /api/v1/products`) from administrative operations (`/api/v1/admin/products`).
- Enable administrative users to query, search, and inspect all products regardless of `isActive` status.
- Allow API Gateways to apply path-based authorization (`/api/v1/admin/**`) and rate-limiting policies cleanly.
- Maintain cache consistency: admin mutations immediately purge public Redis catalog and detail caches.
- Update OpenAPI 3.1 specifications and Swagger/Scalar documentation.

**Non-Goals:**
- Changing database schema or Prisma model definitions (the existing `Product` and `ProductVariant` models already support `isActive`).
- Changing Kafka asynchronous stock event handling.

## Decisions

### 1. Chi Route Hierarchy Restructuring
Separate the Chi router into two distinct route groups within `/api/v1`:

```
/api/v1
├── /products (Public Catalog)
│   ├── RateLimiter(PublicRPM)
│   ├── GET /
│   ├── GET /{id}
│   └── GET /slug/{slug}
│
└── /admin/products (Admin Backoffice)
    ├── RequireAdmin
    ├── RateLimiter(AdminRPM)
    ├── GET /                     (List all products, active/inactive)
    ├── GET /{id}                 (Get single product by ID)
    ├── POST /                    (Create product)
    ├── PUT /{id}                 (Update product)
    ├── DELETE /{id}              (Delete product)
    ├── POST /{id}/variants       (Create variant)
    ├── PUT /{id}/variants/{vId}  (Update variant)
    └── DELETE /{id}/variants/{vId} (Delete variant)
```

*Rationale*: Dedicated path prefixes allow upstream gateways (e.g. Kong, Traefik, AWS API Gateway) to enforce JWT validation, role checks, and IP whitelisting purely by URL pattern `/api/v1/admin/*`.

### 2. ProductFilter Extension in Repository
Extend `repository.ProductFilter` to include:
- `IncludeInactive bool`: When `true`, omits the `db.Product.IsActive.Equals(true)` restriction.
- `IsActive *bool`: When specified, filters explicitly by active or inactive status.

*Rationale*: Public endpoints always set `IncludeInactive = false`, ensuring customers never see unreleased or deactivated items, while admin endpoints set `IncludeInactive = true` by default.

### 3. Caching Strategy
- **Public Catalog**: Retains aggressive Redis cache-aside caching (`ProductListTTL = 3m`, `ProductDetailTTL = 30m`).
- **Admin Endpoints**: Bypass Redis read caches and query PostgreSQL directly to provide real-time accuracy for backoffice operations.
- **Admin Mutations**: Continue invalidating Redis cache keys (`product:detail:*`, `product:slug:*`, and pattern `product:list:*`).

## Risks / Trade-offs

- **[Breaking Change for Admin Clients]** → Frontend admin dashboards or scripts calling `/api/v1/products` with `POST`/`PUT`/`DELETE` must update their target URLs to `/api/v1/admin/products`.
- **[Direct DB Load on Admin Queries]** → Bypassing Redis cache for admin `GET` queries could increase PostgreSQL load if admin traffic spikes. *Mitigation*: Admin queries are protected by strict admin sliding-window rate limiters (`cfg.RateLimitAdminRPM`).
