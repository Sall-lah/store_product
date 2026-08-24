## ADDED Requirements

### Requirement: Cache Invalidation on Event-Driven Stock Changes
The service SHALL invalidate Redis product detail caches (`product:detail:id:{id}` and `product:detail:slug:{slug}`) and catalog list query caches (`product:list:*`) whenever variant stock is adjusted via Kafka events.

#### Scenario: Invalidation following order stock deduction or restock
- **WHEN** an `order.created`, `order.cancelled`, or `order.expired` event modifies the stock quantity of one or more variants
- **THEN** the system identifies the affected parent products, deletes their cached entries in Redis, and purges catalog list caches to ensure real-time inventory visibility.
