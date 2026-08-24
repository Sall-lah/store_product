## 1. Project Scaffolding & Database Setup

- [x] 1.1 Initialize Go module and install dependencies (`go-chi/chi/v5`, `go-redis/v9`, `joho/godotenv`, `prisma-client-go`).
- [x] 1.2 Define `prisma/schema.prisma` with `Product` and `ProductVariant` models with proper relations and indexes.
- [x] 1.3 Create configuration package in `internal/config` supporting `DATABASE_URL`, `PORT`, and `REDIS_PORT` (defaulting to 6739).

## 2. Middleware & Cache Infrastructure

- [x] 2.1 Implement Redis client setup in `internal/cache` with connection pooling and graceful error handling.
- [x] 2.2 Implement Redis sliding-window rate limiting middleware in `internal/middleware` with `X-RateLimit-*` response headers.
- [x] 2.3 Implement API Gateway header authorization middleware in `internal/middleware` verifying `X-User-Role: admin`.
- [x] 2.4 Implement standard recovery, structured JSON logging, and CORS middlewares.

## 3. Repositories & Keyset Pagination

- [x] 3.1 Implement base64 cursor encoding and decoding helpers for `(created_at, id)` keyset pagination in `internal/pkg/cursor`.
- [x] 3.2 Implement Prisma-based Product repository in `internal/repository` supporting cursor pagination, category/price/variant filters, and search.
- [x] 3.3 Implement ProductVariant repository methods for transactional variant creation, updates, and cascading deletes.

## 4. Service Layer & Cache Invalidation

- [x] 4.1 Implement product query service with Redis cache-aside reads for product details (`product:detail:{id}`).
- [x] 4.2 Implement admin write business logic with immediate internal cache invalidation for updated/deleted product keys.
- [x] 4.3 Implement resilient Redis fallback allowing the service to query Supabase directly if Redis is unreachable.

## 5. HTTP Handlers & Route Registration

- [x] 5.1 Implement public HTTP handlers for `GET /api/v1/products` and `GET /api/v1/products/:id` in `internal/handler`.
- [x] 5.2 Implement admin-protected HTTP handlers for `POST`, `PUT`, and `DELETE` on `/api/v1/products` in `internal/handler`.
- [x] 5.3 Register all public and protected routes with Chi router and initialize server entry point in `cmd/server/main.go`.

## 6. Verification & Testing

- [x] 6.1 Implement unit tests for keyset cursor pagination encoding/decoding and query filters.
- [x] 6.2 Implement integration tests for cache-aside hit/miss behavior, write invalidation, and Redis fallback.
- [x] 6.3 Implement test suite for API Gateway header role verification and Redis sliding-window rate limiter.
