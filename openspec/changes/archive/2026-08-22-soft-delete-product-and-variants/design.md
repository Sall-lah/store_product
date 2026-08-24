## Context

In the product catalog service, products and variants were initially deleted using hard delete (`DELETE FROM "Product"` / `DELETE FROM "ProductVariant"`). While simple, hard deletes permanently purge row records and cascade-delete child variants. In production e-commerce systems, deleting products hard breaks downstream audit logs, causes foreign key issues with historical order events, and permanently removes valuable SKU/stock activity records.

This design transitions product and variant deletion mutations to soft deletion using the existing `isActive` column, while maintaining cache invalidation and ensuring public storefront endpoints hide inactive records.

## Goals / Non-Goals

**Goals:**
- Transition `DELETE /api/v1/admin/products/:id` to soft-delete by setting `isActive = false` on the product and all associated variants.
- Transition `DELETE /api/v1/admin/products/:id/variants/:variantId` to soft-delete by setting `isActive = false` on the target variant.
- Retain `slug` and `SKU` values unmodified in PostgreSQL to protect historical integrity and ensure predictable reactivation.
- Prevent storefront data leaks: ensure public endpoints (`GET /api/v1/products/:id` and `GET /api/v1/products/slug/:slug`) return 404 for soft-deleted/inactive products.
- Invalidate Redis caches (`product:detail:*`, `product:slug:*`, `product:list:*`) on soft delete.
- Retain admin visibility via `GET /api/v1/admin/products` and `GET /api/v1/admin/products/:id`.

**Non-Goals:**
- Adding a `deleted_at` timestamp column to the database (existing `is_active` boolean satisfies requirements without schema migration).
- Automatically renaming slugs/SKUs upon deletion.
- Providing a dedicated `/restore` endpoint (reactivation is handled via existing `PUT /api/v1/admin/products/:id` with `{"is_active": true}`).

## Decisions

### 1. Flag-Based Soft Delete via `isActive`
* **Decision**: Update `db.Product.IsActive.Set(false)` and `db.ProductVariant.IsActive.Set(false)` instead of issuing database delete queries.
* **Alternatives Considered**:
  - *Add `deleted_at DateTime?` column*: Requires a Prisma schema migration and DB downtime/migration scripts. `isActive` is already present on both tables and indexed.
  - *Hard delete with archive table*: Excessively complex with duplicate schema maintenance.

### 2. Strategy A: Keep Slug & SKU Reserved
* **Decision**: Preserve original `slug` and `SKU` on soft-deleted rows without appending random deletion suffixes.
* **Rationale**: SKUs correspond to real-world barcodes and historical order line items. Keeping them immutable prevents SKU confusion across Kafka order/stock events. If an admin wishes to reuse a slug or SKU, they should reactivate or update the soft-deleted product.

### 3. Cascade Soft Delete to Variants
* **Decision**: When deleting a parent product, update all associated variant rows for that `productId` to `isActive = false`.
* **Rationale**: Prevents orphaned active variants from remaining searchable or filterable if the parent product is soft-deleted.

### 4. Public Endpoint 404 Guard
* **Decision**: In `ProductRepository.GetProductByID` and `GetProductBySlug`, verify that the found product record has `IsActive == true`. If `IsActive == false`, return `ErrProductNotFound` (HTTP 404).
* **Rationale**: Public customers should not be able to view direct product links for soft-deleted/inactive products.

## Risks / Trade-offs

- **[Risk] Duplicate Slug/SKU conflicts when creating a new product with the same name as a soft-deleted one**
  ➔ *Mitigation*: The API will return `409 Conflict` (duplicate slug/SKU). Admins can inspect soft-deleted products in the backoffice and re-activate or rename them.
- **[Risk] Cache desynchronization if Redis deletion fails**
  ➔ *Mitigation*: Existing Redis cache-aside design with short TTLs and explicit cache key purges (`Del` on detail/slug + `DelPattern` on list queries) minimizes staleness window.
- **[Risk] Race condition during cascading variant deactivation**
  ➔ *Mitigation*: Repository updates the product and its variants within a single database batch or query block.

## Migration Plan

No database migrations required. The application code changes can be deployed directly with zero database downtime.
