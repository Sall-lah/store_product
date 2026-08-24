# Cache Invalidation Specification

## Purpose
Governs Redis cache-aside reads, TTL management, proactive internal cache invalidation upon writes, and database fallbacks.

## Requirements

### Requirement: Cache-Aside Reads for Product Detail
The service SHALL cache product detail responses in Redis on port 6379 with key pattern `product:detail:id:{id}` and `product:detail:slug:{slug}` with an expiration TTL.

#### Scenario: Cache hit on product detail
- **WHEN** a client requests `GET /api/v1/products/:id` and the key exists in Redis
- **THEN** the system returns cached data directly from Redis without issuing queries to Supabase.

#### Scenario: Cache miss on product detail
- **WHEN** a client requests `GET /api/v1/products/:id` and the key is not in Redis
- **THEN** the system fetches the record from Supabase, serializes and stores it in Redis with TTL, and returns the response to the caller.

### Requirement: Internal Cache Invalidation on Writes
The service SHALL immediately delete associated product detail cache keys and invalidate list caches when an admin creates, updates, or deletes a product.

#### Scenario: Invalidation on update
- **WHEN** an admin updates a product via `PUT /api/v1/products/:id`
- **THEN** the system deletes `product:detail:id:{id}` and associated slug keys from Redis before completing the HTTP response.

#### Scenario: Invalidation on deletion
- **WHEN** an admin deletes a product via `DELETE /api/v1/products/:id`
- **THEN** the system deletes all corresponding Redis keys for that product ID and slug.

### Requirement: Graceful Cache Fallback
The service SHALL gracefully fallback to database queries if Redis becomes unavailable or encounters connection timeouts.

#### Scenario: Redis connection failure
- **WHEN** Redis is down or unreachable
- **THEN** the system logs the Redis error and directly queries Supabase, successfully returning product data without returning HTTP 500.

### Requirement: Cache Invalidation on Event-Driven Stock Changes
The service SHALL invalidate Redis product detail caches (`product:detail:id:{id}` and `product:detail:slug:{slug}`) and catalog list query caches (`product:list:*`) whenever variant stock is adjusted via Kafka events.

#### Scenario: Invalidation following order stock deduction or restock
- **WHEN** an `order.created`, `order.cancelled`, or `order.expired` event modifies the stock quantity of one or more variants
- **THEN** the system identifies the affected parent products, deletes their cached entries in Redis, and purges catalog list caches to ensure real-time inventory visibility.

