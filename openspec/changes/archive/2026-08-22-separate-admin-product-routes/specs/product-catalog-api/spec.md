## MODIFIED Requirements

### Requirement: Keyset Cursor Pagination for Product Listings
The service SHALL provide an endpoint `GET /api/v1/products` that returns a paginated list of active customer-facing products using cursor-based pagination with limit support.

#### Scenario: Initial page request without cursor
- **WHEN** a client sends `GET /api/v1/products?limit=20`
- **THEN** the system returns up to 20 active products ordered by creation date descending, along with a `next_cursor` and `has_more` boolean flag.

#### Scenario: Subsequent page request using cursor
- **WHEN** a client sends `GET /api/v1/products?limit=20&cursor=<encoded_cursor>`
- **THEN** the system returns the next 20 products created before the cursor timestamp without scanning skipped records.
