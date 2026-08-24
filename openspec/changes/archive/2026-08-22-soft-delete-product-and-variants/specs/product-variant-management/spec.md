## MODIFIED Requirements

### Requirement: Admin Delete Product
The service SHALL provide `DELETE /api/v1/admin/products/:id` to soft-delete a product and all its associated variants by marking their active status as false while preserving database records, slugs, and SKUs.

#### Scenario: Successful product soft deletion
- **WHEN** an admin issues `DELETE /api/v1/admin/products/:id` with `X-User-Role: admin`
- **THEN** the system sets `isActive = false` on the product, sets `isActive = false` on all its associated variants in Supabase, invalidates detail and list caches, and returns HTTP 204 No Content.

#### Scenario: Deleting non-existent product
- **WHEN** an admin issues `DELETE /api/v1/admin/products/unknown_id` with `X-User-Role: admin`
- **THEN** the system returns HTTP 404 Not Found.

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
