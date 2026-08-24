# Gateway Authentication and Rate Limiting Specification

## Purpose
Defines API Gateway header validation (`X-User-Role`, `X-User-Id`) and Redis sliding window rate limiting.

## Requirements

### Requirement: Upstream Gateway Header Role Authorization
The service SHALL enforce that all endpoints under the `/api/v1/admin/**` hierarchy verify the presence of `X-User-Role: admin` using case-insensitive comparison (accepting "admin", "ADMIN", or mixed-case).

#### Scenario: Authorized admin request with lowercase role
- **WHEN** a client or gateway forwards a request to `/api/v1/admin/**` with header `X-User-Role: admin` and `X-User-Id: <user_id>`
- **THEN** the middleware passes the request to the handler and attaches the user identity to the Go request context.

#### Scenario: Authorized admin request with uppercase role
- **WHEN** a client or gateway forwards a request to `/api/v1/admin/**` with header `X-User-Role: ADMIN` and `X-User-Id: <user_id>`
- **THEN** the middleware passes the request to the handler and attaches the user identity to the Go request context.

#### Scenario: Unauthorized or missing role on admin path
- **WHEN** a client sends a request to `/api/v1/admin/**` with missing `X-User-Role` or `X-User-Role` value that does not match "admin" case-insensitively (e.g., "customer", "seller", or empty)
- **THEN** the middleware rejects the request with HTTP 403 Forbidden without processing the business logic.

### Requirement: Redis-Backed Sliding Window Rate Limiter
The service SHALL limit incoming requests using a Redis-backed sliding window counter algorithm and attach standard `X-RateLimit-*` headers to responses.

#### Scenario: Within allowable rate limits
- **WHEN** a client makes requests within the configured threshold (e.g., 60 req/min for search, 120 req/min for catalog)
- **THEN** the system executes the request and returns `X-RateLimit-Limit`, `X-RateLimit-Remaining`, and `X-RateLimit-Reset` headers.

#### Scenario: Exceeded rate limit threshold
- **WHEN** a client exceeds the maximum allowed requests in the rolling 60-second window
- **THEN** the system immediately rejects the request with HTTP 429 Too Many Requests and a `Retry-After` header.
