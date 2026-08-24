## MODIFIED Requirements

### Requirement: Upstream Gateway Header Role Authorization
The service SHALL enforce that write endpoints (`POST`, `PUT`, `DELETE` on `/api/v1/products`) verify the presence of `X-User-Role: admin` using case-insensitive comparison (accepting "admin", "ADMIN", or mixed-case).

#### Scenario: Authorized admin request with lowercase role
- **WHEN** a client or gateway forwards a write request with header `X-User-Role: admin` and `X-User-Id: <user_id>`
- **THEN** the middleware passes the request to the handler and attaches the user identity to the Go request context.

#### Scenario: Authorized admin request with uppercase role
- **WHEN** a client or gateway forwards a write request with header `X-User-Role: ADMIN` and `X-User-Id: <user_id>`
- **THEN** the middleware passes the request to the handler and attaches the user identity to the Go request context.

#### Scenario: Unauthorized or missing role
- **WHEN** a client sends a write request with missing `X-User-Role` or `X-User-Role` value that does not match "admin" case-insensitively (e.g., "customer", "seller", or empty)
- **THEN** the middleware rejects the request with HTTP 403 Forbidden without processing the business logic.
