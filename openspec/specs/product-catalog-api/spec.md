# Product Catalog API Specification

## Purpose
Governs public read endpoints for querying products with cursor-based pagination, multi-field filtering, search, and detail retrieval.

## Requirements

### Requirement: Keyset Cursor Pagination for Product Listings
The service SHALL provide an endpoint `GET /api/v1/products` that returns a paginated list of active products using cursor-based pagination with limit support.

#### Scenario: Initial page request without cursor
- **WHEN** a client sends `GET /api/v1/products?limit=20`
- **THEN** the system returns up to 20 active products ordered by creation date descending, along with a `next_cursor` and `has_more` boolean flag.

#### Scenario: Subsequent page request using cursor
- **WHEN** a client sends `GET /api/v1/products?limit=20&cursor=<encoded_cursor>`
- **THEN** the system returns the next 20 products created before the cursor timestamp without scanning skipped records.

### Requirement: Multi-Attribute Filtering
The service SHALL allow filtering product listings by category, minimum price, maximum price, size, and color.

#### Scenario: Filter by category and price range
- **WHEN** a client sends `GET /api/v1/products?category=shoes&min_price=50&max_price=150`
- **THEN** the system returns only active products matching category "shoes" with base or variant price between 50 and 150.

#### Scenario: Filter by variant attributes
- **WHEN** a client sends `GET /api/v1/products?size=L&color=Black`
- **THEN** the system returns only active products that possess at least one active variant with size "L" and color "Black".

### Requirement: Search Products by Keyword
The service SHALL allow keyword search across product name, description, and variant SKU.

#### Scenario: Search by query keyword
- **WHEN** a client sends `GET /api/v1/products?search=running`
- **THEN** the system returns active products where name, description, or variant SKU matches "running".

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
