## 1. Dependencies & Configuration

- [x] 1.1 Add `github.com/segmentio/kafka-go` dependency to `go.mod`
- [x] 1.2 Add Kafka configuration parameters (`KAFKA_BROKERS`, `KAFKA_TOPIC_ORDER_EVENTS`, `KAFKA_CONSUMER_GROUP`) to `internal/config/config.go`
- [x] 1.3 Update `.env.example` with Kafka configuration settings

## 2. Persistence & Repository Layer

- [x] 2.1 Add `processed_events` table or idempotency model to track consumed event IDs
- [x] 2.2 Implement atomic stock deduction and restocking methods in `internal/repository/variant.go`
- [x] 2.3 Implement event idempotency repository methods in `internal/repository/event.go`

## 3. Kafka Event Consumer & Handler

- [x] 3.1 Define order event payload models (`OrderEventPayload`, `OrderItemPayload`) in `internal/event/models.go`
- [x] 3.2 Implement order event processing service in `internal/event/handler.go` to coordinate stock updates and cache invalidation
- [x] 3.3 Implement Kafka consumer worker loop with retry handling and graceful shutdown in `internal/event/consumer.go`

## 4. Server Integration & Verification

- [x] 4.1 Initialize Kafka consumer worker in `cmd/server/main.go` alongside HTTP server
- [x] 4.2 Write unit tests for event unmarshaling, stock deduction, restocking, and idempotency logic
- [x] 4.3 Update project documentation and architecture diagrams to reflect Kafka event integration
