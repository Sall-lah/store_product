## ADDED Requirements

### Requirement: Admin Create Product with Variants
The service SHALL provide `POST /api/v1/products` allowing an admin to create a new product along with one or more initial variants.

#### Scenario: Successful product creation
- **WHEN** an admin sends a valid product payload with name, category, basePrice, and an array of variants (size, color, sku, stock, price)
- **THEN** the system persists the product and variants in Supabase via Prisma and returns HTTP 201 Created with the generated product entity.

#### Scenario: Duplicate SKU validation
- **WHEN** an admin attempts to create a variant with an SKU that already exists
- **THEN** the system rejects the operation and returns HTTP 409 Conflict with an explanatory error message.

### Requirement: Admin Update Product and Variants
The service SHALL provide `PUT /api/v1/products/:id` to update core product fields and upsert/modify variants.

#### Scenario: Successful product update
- **WHEN** an admin submits updated details for an existing product
- **THEN** the system updates the database record, modifies variant properties, updates timestamps, and returns HTTP 200 OK.

### Requirement: Admin Delete Product
The service SHALL provide `DELETE /api/v1/products/:id` to remove or soft-delete a product and its associated variants.

#### Scenario: Successful product deletion
- **WHEN** an admin issues `DELETE /api/v1/products/:id` for an existing product
- **THEN** the system cascades deletion to product variants in Supabase and returns HTTP 204 No Content.
