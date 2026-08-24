## Context

The `store_product` service is a standalone Go microservice responsible for the e-commerce product catalog. It interfaces with Supabase (PostgreSQL) as the persistent system of record via Prisma Client Go, and connects to Redis (configured on port `6739` in local development) as a high-speed caching and rate-limiting layer. Public clients query product catalogs, search, and details, while administrative operations (create, update, delete) are authorized via upstream API Gateway headers and trigger immediate internal cache invalidation.

## Goals / Non-Goals

**Goals:**
- Provide low-latency read APIs for product browsing, filtering, searching, and detail retrieval.
- Implement high-performance keyset (cursor-based) pagination that scales efficiently with millions of products.
- Support multi-attribute product variants (size, color, SKU, specific pricing, and inventory).
- Cache single product details and catalog queries in Redis with automatic internal cache invalidation on writes.
- Enforce API Gateway header-based role checking (`X-User-Role: admin`) on modifying endpoints.
- Enforce tiered sliding-window rate limiting using Redis (public catalog, search, and admin endpoints).
- Organize Go code in a clean, modular structure with TSDoc/JSDoc/GoDoc comments.

**Non-Goals:**
- Direct JWT token decryption or user identity management (handled upstream by API Gateway / Auth microservice).
- Order placement, cart checkout, or payment processing (handled by separate microservices).
- Heavy semantic AI search (basic full-text / trigram / ILIKE filtering on Postgres is sufficient for initial version).

## Decisions

### 1. Keyset Cursor Pagination
- **Decision**: Use base64-encoded composite cursor `(created_at, id)` for `GET /api/v1/products`.
- **Rationale**: Offset-based pagination (`OFFSET N LIMIT M`) degrades to $O(N)$ because the database must scan and discard $N$ rows. Keyset pagination jumps directly to index keys with $O(1)$ seek time.
- **Alternatives Considered**: Offset/Limit pagination (rejected due to poor scaling on large catalogs).

### 2. Upstream Gateway Header Authorization
- **Decision**: Validate admin access via `X-User-Role: admin` and extract actor ID from `X-User-Id`.
- **Rationale**: Keeps the service loosely coupled from identity provider internals and token secret rotation. The API Gateway serves as the single authentication perimeter.
- **Alternatives Considered**: In-service JWT decoding with shared Supabase JWT secret (adds unnecessary coupling and secret management).

### 3. Redis-Backed Sliding Window Rate Limiting
- **Decision**: Implement sliding window rate limiting in Go middleware using Redis atomic operations (`ZREMRANGEBYSCORE`, `ZADD`, `ZCARD`, `EXPIRE`).
- **Rationale**: Prevents burst spikes at window boundaries inherent to fixed-window counters while supporting distinct limits for public catalog (`120 req/min`), search (`60 req/min`), and admin writes (`30 req/min`).
- **Alternatives Considered**: In-memory token bucket (rejected because microservices run horizontally in multiple instances).

### 4. Cache-Aside with Internal Invalidation
- **Decision**:
  - `GET /api/v1/products/:id`: Read from Redis `product:detail:{id}` (TTL: 30m). On miss, query Supabase via Prisma and write to Redis.
  - `POST/PUT/DELETE /api/v1/products`: Execute write in Supabase via Prisma, then immediately delete `product:detail:{id}` and bump catalog cache version/tag.
- **Rationale**: Eliminates stale product details immediately upon admin updates without requiring external CDC or message brokers.

### 5. Go Project Layout & Router
- **Decision**: Use `go-chi/chi/v5` router with clean layered architecture:
  - `cmd/server/main.go`: Entry point & dependency injection.
  - `internal/config`: Environment variable parser (supports custom `REDIS_PORT=6739`, `DATABASE_URL`).
  - `internal/handler`: HTTP request/response handlers with JSON serialization.
  - `internal/service`: Core business logic, cache invalidation, and orchestration.
  - `internal/repository`: Prisma Client Go database operations.
  - `internal/cache`: Redis client wrapper and key helpers.
  - `internal/middleware`: Gateway auth validator, rate limiter, logger, and CORS.
- **Rationale**: Chi is lightweight, stdlib `net/http` compatible, fast, and promotes idiomatic Go.

## Risks / Trade-offs

- **[Redis Outage / Unavailability]** → Service middleware and cache helpers MUST catch Redis connection errors and gracefully fallback to direct Supabase queries without returning 500 errors to callers.
- **[Cache Stampede on Hot Products]** → Implement short mutex locks (singleflight) or moderate TTL jitter for high-traffic product details.
- **[Spoofed Gateway Headers]** → Product service should only be accessible within the private cluster network or require a shared internal Gateway secret header (`X-Gateway-Secret`) if exposed publicly.

## Migration & Deployment Plan

1. Initialize Go module `github.com/Sall-lah/store_product`.
2. Generate Prisma Client Go schema and apply migration to Supabase PostgreSQL database.
3. Verify local Redis connection on port `6739`.
4. Deploy containerized service behind API Gateway.
