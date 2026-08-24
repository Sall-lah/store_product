## Context

The `store_product` microservice is responsible for product catalog management, variants, and stock levels. Currently, stock is only mutated through synchronous admin REST endpoints. The Order Service publishes order lifecycle events (`order.created`, `order.cancelled`, `order.expired`) to a Kafka topic named `order.events`. To maintain real-time inventory consistency, prevent overselling, and restock expired/cancelled reservations, `store_product` must consume these events asynchronously and reliably.

## Goals / Non-Goals

**Goals:**
- Implement a robust, non-blocking Kafka consumer worker in Go subscribing to `order.events`.
- Handle `order.created` (reduce variant stock), `order.cancelled` (restock variant stock), and `order.expired` (restock variant stock).
- Guarantee event processing idempotency using a database tracking table (`processed_events`).
- Execute atomic stock updates in PostgreSQL to avoid race conditions and negative inventory.
- Ensure Redis cache invalidation for parent products and catalog query caches upon stock changes.
- Graceful shutdown of Kafka consumer goroutines during service termination.

**Non-Goals:**
- Creating a full payment saga coordinator (Order Service orchestrates sagas; `store_product` acts as an inventory participant).
- Replacing REST APIs for admin product management.

## Decisions

### Decision 1: Kafka Client Library (`segmentio/kafka-go`)
- **Choice**: Use `github.com/segmentio/kafka-go`.
- **Rationale**: `kafka-go` is pure Go, requires no CGO or native `librdkafka` dependencies (essential for clean Windows cross-compilation and lightweight Alpine Docker images), provides idiomatic Reader APIs, and supports consumer groups seamlessly.
- **Alternatives Considered**:
  - `confluent-kafka-go`: High performance but requires CGO and `librdkafka`, complicating Windows dev setups and multi-stage Docker builds.
  - `IBM/sarama`: Robust pure-Go client, but has higher boilerplate configuration and complex consumer group setup compared to `kafka-go`.

### Decision 2: Idempotency Tracking Table (`processed_events`)
- **Choice**: Persist processed events in a dedicated PostgreSQL table with schema `(event_id VARCHAR PRIMARY KEY, event_type VARCHAR, order_id VARCHAR, processed_at TIMESTAMP)`.
- **Rationale**: Kafka operates on at-least-once delivery semantics. Consumer rebalances or crashes can cause redelivery of events. Checking and inserting `event_id` in the same database transaction as the stock mutation guarantees exact-once execution semantics.
- **Alternatives Considered**:
  - Redis SET with TTL: Fast, but volatile during Redis restarts and doesn't participate in atomic database transactions.

### Decision 3: Atomic Stock Mutations via Raw SQL / Prisma Execute
- **Choice**: Execute atomic SQL queries:
  - Decrement: `UPDATE "ProductVariant" SET stock = stock - $1 WHERE id = $2 AND stock >= $1`
  - Increment: `UPDATE "ProductVariant" SET stock = stock + $1 WHERE id = $2`
- **Rationale**: Prevents read-modify-write race conditions when concurrent orders target the same variant SKU.

### Decision 4: Architecture & Worker Organization
- **Choice**: Dedicated package `internal/event/consumer.go` and `internal/event/handler.go` with dependency-injected services. Main server in `cmd/server/main.go` launches the consumer loop in a managed background goroutine with context cancellation.

## Risks / Trade-offs

- **[Risk] Out-of-order events (e.g. `order.cancelled` arrives before `order.created`)**:
  → *Mitigation*: The Order Service includes full line items `[{variant_id, quantity}]` in `order.cancelled` and `order.expired` events. In a direct stock model, restocking restores the quantity safely.
- **[Risk] Insufficient stock on `order.created`**:
  → *Mitigation*: Atomic SQL update checks `stock >= quantity`. If rows affected is 0, record the stock failure and log an alert.
- **[Risk] Redis cache stampede after massive restocking**:
  → *Mitigation*: Invalidate keys with single deletes (`Del`) and pattern deletes (`DelPattern`) asynchronously so catalog reads rebuild cache on demand.
