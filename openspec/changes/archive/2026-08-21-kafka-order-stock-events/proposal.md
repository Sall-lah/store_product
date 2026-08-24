## Why

Currently, the `store_product` service only provides HTTP endpoints for static catalog querying and admin stock CRUD operations. When orders are created, cancelled, or expired in the Order Service, variant stock quantities in `store_product` are not synchronized asynchronously. Integrating Kafka consumer workers on topic `order.events` enables real-time stock reduction upon order placement, atomic restocking upon order cancellation or expiration, prevents overselling, and guarantees inventory consistency across microservices.

## What Changes

- Introduce a Kafka consumer background worker subscribing to topic `order.events`.
- Process `order.created` (or order placement) events to atomically reduce variant stock quantities.
- Process `order.cancelled` and `order.expired` events to atomically restock variant stock quantities.
- Implement an event idempotency mechanism (`processed_events` tracking) to ensure at-least-once message delivery does not cause duplicate stock deductions or restocking.
- Trigger Redis cache invalidation for product details (`product:detail:*`, `product:slug:*`) and catalog query lists (`product:list:*`) whenever stock updates occur.
- Add Kafka broker and consumer group configuration settings to environment variables.

## Capabilities

### New Capabilities
- `kafka-stock-events`: Consumes `order.events` from Kafka to execute idempotent stock reduction on order creation and stock replenishment on order cancellation or expiration.

### Modified Capabilities
- `product-variant-management`: Expands variant stock management with atomic decrement/increment operations for high-concurrency order lifecycle events.
- `cache-invalidation`: Expands cache invalidation triggers to include async Kafka-driven stock mutations.

## Impact

- **Dependencies**: Adds pure-Go Kafka client library (e.g. `segmentio/kafka-go`).
- **Data Models / Persistence**: Adds atomic database query helpers for variant stock mutation and an idempotency log / tracking table in PostgreSQL / Prisma schema.
- **Service Lifecycle**: Background consumer goroutine initialized and shut down gracefully in `cmd/server/main.go`.
- **Configuration**: New environment variables: `KAFKA_BROKERS`, `KAFKA_TOPIC_ORDER_EVENTS`, `KAFKA_CONSUMER_GROUP`.
