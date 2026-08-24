## Why

The e-commerce application requires a dedicated, high-performance Product Microservice to serve product catalog data, handle complex filtering and search, manage multi-attribute variants (size, color), and process admin modifications with low latency. By leveraging Redis caching, Supabase PostgreSQL, and Go, the service provides sub-millisecond cached responses while maintaining data consistency through internal cache invalidation.

## What Changes

- Scaffold a Go microservice project with modular architecture and Prisma Client Go.
- Implement public read endpoints for product catalog with high-performance cursor-based pagination, multi-field filtering, search, and detail retrieval.
- Implement product variant management supporting sizes, colors, individual SKUs, stock levels, and custom pricing.
- Implement Redis (port 6739 in dev) cache-aside layer for fast catalog reads and automatic internal cache invalidation on admin writes.
- Implement admin write endpoints (`POST`, `PUT`, `DELETE`) protected by API Gateway header authentication (`X-User-Role: admin`, `X-User-Id`).
- Implement Redis-backed sliding window rate limiting for public, search, and admin endpoints.

## Capabilities

### New Capabilities
- `product-catalog-api`: High-performance public API for querying products with cursor-based pagination, multi-attribute filtering, search, and single-item retrieval.
- `product-variant-management`: Data modeling and admin CRUD operations for products and their variants (sizes, colors, SKUs, inventory, prices).
- `cache-invalidation`: Redis caching layer with cache-aside read pattern and internal invalidation on create, update, and delete actions.
- `gateway-auth-ratelimit`: Middleware for API Gateway header-based role verification and Redis-backed sliding window rate limiting.

### Modified Capabilities
<!-- No existing capabilities modified; this is a new service initialization. -->

## Impact

- **New Services**: Creates the standalone `store_product` Go microservice repository.
- **Databases**: Establishes Supabase PostgreSQL schema for `Product` and `ProductVariant` via Prisma.
- **Dependencies**: Introduces Go runtime dependencies (`prisma-client-go`, `go-redis/v9`, HTTP router like Chi, etc.) and Redis (defaulting to port 6739 in local dev).
- **APIs**: Exposes `/api/v1/products` and `/api/v1/categories` endpoints consumed by API Gateway and frontend clients.
