# Product Variant Management Specification

## Purpose
Defines data modeling and admin CRUD operations for products and their multi-attribute variants (sizes, colors, SKUs, inventory, and pricing).

## Requirements

### Requirement: Admin Create Product with Variants
The service SHALL provide `POST /api/v1/admin/products` allowing an admin to create a new product along with one or more initial variants.

#### Scenario: Successful product creation
- **WHEN** an admin sends a valid product payload with name, category, basePrice, and an array of variants to `POST /api/v1/admin/products` with `X-User-Role: admin`
- **THEN** the system persists the product and variants in Supabase via Prisma and returns HTTP 201 Created with the generated product entity.

#### Scenario: Duplicate SKU validation
- **WHEN** an admin attempts to create a variant with an SKU that already exists
- **THEN** the system rejects the operation and returns HTTP 409 Conflict with an explanatory error message.

### Requirement: Admin Update Product and Variants
The service SHALL provide `PUT /api/v1/admin/products/:id` to update core product fields and modify variants.

#### Scenario: Successful product update
- **WHEN** an admin submits updated details for an existing product to `PUT /api/v1/admin/products/:id` with `X-User-Role: admin`
- **THEN** the system updates the database record, modifies variant properties, purges Redis caches, and returns HTTP 200 OK.

### Requirement: Admin Delete Product
The service SHALL provide `DELETE /api/v1/admin/products/:id` to soft-delete a product and all its associated variants by marking their active status as false while preserving database records, slugs, and SKUs.

#### Scenario: Successful product soft deletion
- **WHEN** an admin issues `DELETE /api/v1/admin/products/:id` with `X-User-Role: admin`
- **THEN** the system sets `isActive = false` on the product, sets `isActive = false` on all its associated variants in Supabase, invalidates detail and list caches, and returns HTTP 204 No Content.

#### Scenario: Deleting non-existent product
- **WHEN** an admin issues `DELETE /api/v1/admin/products/unknown_id` with `X-User-Role: admin`
- **THEN** the system returns HTTP 404 Not Found.

### Requirement: Admin List and Inspect All Products
The service SHALL provide `GET /api/v1/admin/products` allowing authenticated administrators to retrieve paginated products including inactive, draft, and out-of-stock items.

#### Scenario: Admin lists products with inactive items included
- **WHEN** an admin sends `GET /api/v1/admin/products` with `X-User-Role: admin`
- **THEN** the system returns active and inactive products matching query filters without filtering out `isActive = false` records.

#### Scenario: Admin filters by active status
- **WHEN** an admin sends `GET /api/v1/admin/products?is_active=false` with `X-User-Role: admin`
- **THEN** the system returns only inactive/draft products.

#### Scenario: Admin views inactive product
- **WHEN** an admin requests `GET /api/v1/admin/products/:id` for an inactive product with `X-User-Role: admin`
- **THEN** the system returns the product details and HTTP status 200 OK.

### Requirement: Admin Retrieve Single Product by ID
The service SHALL provide `GET /api/v1/admin/products/:id` returning complete product data and all variants regardless of product or variant active state.

#### Scenario: Admin views inactive product
- **WHEN** an admin requests `GET /api/v1/admin/products/:id` for an inactive product with `X-User-Role: admin`
- **THEN** the system returns the product details and HTTP status 200 OK.

### Requirement: Admin Variant CRUD Operations
The service SHALL provide dedicated variant mutation endpoints under the `/api/v1/admin/products/:id/variants` path.

#### Scenario: Admin creates a variant
- **WHEN** an admin sends `POST /api/v1/admin/products/:id/variants` with variant SKU, price, size, color, and stock
- **THEN** the system creates the variant linked to the product, invalidates parent product caches, and returns HTTP 201 Created.

#### Scenario: Admin updates a variant
- **WHEN** an admin sends `PUT /api/v1/admin/products/:id/variants/:variantId` with updated variant properties
- **THEN** the system updates the variant record, invalidates parent product caches, and returns HTTP 200 OK.

#### Scenario: Admin deletes a variant
- **WHEN** an admin sends `DELETE /api/v1/admin/products/:id/variants/:variantId`
- **THEN** the system soft-deletes the variant by setting its `isActive = false`, invalidates parent product caches, and returns HTTP 204 No Content.

### Requirement: Atomic Stock Adjustment for Variants
The repository and service layer SHALL provide atomic stock adjustment operations (`AdjustStock` / `DecrementStock` / `IncrementStock`) to increment or decrement variant stock safely under concurrent execution without race conditions.

#### Scenario: Concurrent stock deduction
- **WHEN** multiple order events concurrently decrement stock for the same variant
- **THEN** the database executes atomic updates (`stock = stock - delta` where `stock >= delta`), ensuring stock never becomes negative and preventing race-condition overselling.

#### Scenario: Atomic stock increment on restock
- **WHEN** an order cancellation or expiration event increments stock for a variant
- **THEN** the database executes an atomic update (`stock = stock + delta`), correctly restoring inventory.
