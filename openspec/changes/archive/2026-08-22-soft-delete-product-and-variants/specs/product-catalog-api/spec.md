## MODIFIED Requirements

### Requirement: Retrieve Single Product Detail
The service SHALL provide `GET /api/v1/products/:id` and `GET /api/v1/products/slug/:slug` returning complete product data including all associated active variants for active products only. Soft-deleted and inactive products SHALL return HTTP 404 Not Found.

#### Scenario: Existing active product retrieval
- **WHEN** a client requests `GET /api/v1/products/prod_123`
- **THEN** the system returns product information, base price, and array of active variants with HTTP status 200.

#### Scenario: Non-existent product retrieval
- **WHEN** a client requests `GET /api/v1/products/unknown_id`
- **THEN** the system returns a 404 Not Found JSON response.

#### Scenario: Soft-deleted or inactive product retrieval
- **WHEN** a public client requests `GET /api/v1/products/inactive_prod_id` or `GET /api/v1/products/slug/inactive-prod-slug` for a product where `isActive = false`
- **THEN** the system returns a 404 Not Found JSON response.
