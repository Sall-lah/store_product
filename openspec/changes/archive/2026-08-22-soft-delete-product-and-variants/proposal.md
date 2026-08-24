## Why

Currently, `DELETE /api/v1/admin/products/:id` and `DELETE /api/v1/admin/products/:id/variants/:variantId` permanently remove records from the database via hard delete (cascading across variants). This causes loss of historical data, potential foreign key integrity violations with past order records or Kafka inventory logs, and prevents simple product reactivation. Switching to soft delete preserves product identity, audit trails, and inventory history while safely hiding deleted items from the storefront.

## What Changes

- Change `ProductRepository.DeleteProduct` to soft-delete by setting `Product.isActive = false` and cascading `ProductVariant.isActive = false` across all associated variants.
- Change `VariantRepository.DeleteVariant` to soft-delete by setting `ProductVariant.isActive = false`.
- Keep product slugs and variant SKUs unchanged in the database (Strategy A: Keep Reserved) to prevent SKU divergence in downstream inventory/order event logs and enable clean reactivation.
- Update public storefront lookup (`GetProductByID` and `GetProductBySlug`) to return `404 Not Found` when a product is soft-deleted (`isActive == false`).
- Preserve Redis cache purging on soft deletion (`product:detail:<id>`, `product:slug:<slug>`, and `product:list:*`) so the storefront immediately stops serving soft-deleted items.
- Ensure admin queries (`AdminListProducts` and `AdminGetProductByID`) retain visibility into soft-deleted products for audit and reactivation purposes.

## Capabilities

### New Capabilities
<!-- No new standalone capabilities are introduced. -->

### Modified Capabilities
- `product-variant-management`: Update product and variant deletion requirements to specify soft deletion (`isActive = false`) instead of database record removal.
- `product-catalog-api`: Update single product retrieval requirements to ensure soft-deleted (`isActive = false`) products return 404 Not Found on public storefront endpoints.

## Impact

- **Affected code**: `internal/repository/product.go`, `internal/repository/variant.go`, `internal/service/product.go`, and corresponding repository/handler unit and integration tests.
- **Database**: No schema migration required (uses existing `isActive` boolean fields in `Product` and `ProductVariant` tables).
- **APIs**:
  - `DELETE /api/v1/admin/products/:id` continues returning `204 No Content` but performs soft delete.
  - `DELETE /api/v1/admin/products/:id/variants/:variantId` continues returning `204 No Content` but performs soft delete.
  - Public `GET /api/v1/products/:id` and `GET /api/v1/products/slug/:slug` will now return 404 if the target product is inactive/soft-deleted.
- **Dependencies/Systems**: Historical references in Kafka order/stock events and database relations remain preserved.
