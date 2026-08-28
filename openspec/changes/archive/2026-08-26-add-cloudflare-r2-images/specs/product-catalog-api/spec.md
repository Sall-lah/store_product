## MODIFIED Requirements

### Requirement: Keyset Cursor Pagination for Product Listings
The service SHALL provide an endpoint `GET /api/v1/products` that returns a paginated list of active products using cursor-based pagination with limit support, including the `primary_image_url` for fast catalog card rendering.

#### Scenario: Initial page request without cursor
- **WHEN** a client sends `GET /api/v1/products?limit=20`
- **THEN** the system returns up to 20 active products ordered by creation date descending, along with `primary_image_url`, a `next_cursor`, and `has_more` boolean flag.

#### Scenario: Subsequent page request using cursor
- **WHEN** a client sends `GET /api/v1/products?limit=20&cursor=<encoded_cursor>`
- **THEN** the system returns the next 20 products created before the cursor timestamp without scanning skipped records, including `primary_image_url` for each product.

### Requirement: Retrieve Single Product Detail
The service SHALL provide `GET /api/v1/products/:id` and `GET /api/v1/products/slug/:slug` returning complete product data including all associated active variants and the full array of gallery images with variant associations for active products only. Soft-deleted and inactive products SHALL return HTTP 404 Not Found.

#### Scenario: Existing active product retrieval
- **WHEN** a client requests `GET /api/v1/products/prod_123`
- **THEN** the system returns product information, base price, array of active variants, and array of gallery images ordered by sort order with HTTP status 200.

#### Scenario: Non-existent product retrieval
- **WHEN** a client requests `GET /api/v1/products/unknown_id`
- **THEN** the system returns a 404 Not Found JSON response.

#### Scenario: Soft-deleted or inactive product retrieval
- **WHEN** a public client requests `GET /api/v1/products/inactive_prod_id` or `GET /api/v1/products/slug/inactive-prod-slug` for a product where `isActive = false`
- **THEN** the system returns a 404 Not Found JSON response.
