## ADDED Requirements

### Requirement: Kafka Consumer for Order Events
The service SHALL run a background Kafka consumer subscribing to the configured topic (default `order.events`) using a dedicated consumer group (default `store_product_stock_worker`) to receive order lifecycle events.

#### Scenario: Background worker initialization and shutdown
- **WHEN** the application starts up
- **THEN** it initializes the Kafka consumer worker pool in the background and gracefully terminates consumer connections upon SIGINT/SIGTERM.

#### Scenario: Malformed event handling
- **WHEN** the consumer receives an unparseable or schema-invalid JSON message on `order.events`
- **THEN** it logs a warning/error and acknowledges/commits the offset or routes to dead-letter handling without crashing the consumer loop.

### Requirement: Order Placed Stock Deduction
The service SHALL consume `order.created` (or order placement) events and atomically reduce available stock for each variant item in the order.

#### Scenario: Sufficient stock reduction
- **WHEN** an `order.created` event is received containing variant IDs and quantities, and all variants have sufficient stock
- **THEN** the system atomically decreases the variant stock counts in PostgreSQL, records the processed event ID, and invalidates affected cache keys.

#### Scenario: Insufficient stock handling
- **WHEN** an `order.created` event is received for a variant whose remaining stock is less than the requested quantity
- **THEN** the system prevents negative stock, logs the stock exhaustion failure, and emits an alert or dead-letter notification.

### Requirement: Order Cancelled Stock Restocking
The service SHALL consume `order.cancelled` events and atomically increase variant stock quantities by the items in the cancelled order.

#### Scenario: Successful restocking on cancellation
- **WHEN** an `order.cancelled` event is received with valid variant items and quantities
- **THEN** the system atomically increments each variant's stock count in PostgreSQL, logs the event as processed, and purges Redis cache keys for the affected products.

### Requirement: Order Expired Stock Restocking
The service SHALL consume `order.expired` events and atomically increase variant stock quantities for unpaid or timed-out orders.

#### Scenario: Successful restocking on expiration
- **WHEN** an `order.expired` event is received with variant items and quantities
- **THEN** the system atomically increments each variant's stock count in PostgreSQL, logs the event as processed, and purges Redis cache keys for the affected products.

### Requirement: Event Processing Idempotency
The service SHALL track processed event IDs or order event states to ensure at-least-once Kafka deliveries do not result in duplicate stock deductions or restocking.

#### Scenario: Duplicate event redelivery
- **WHEN** a previously processed `order.cancelled` or `order.expired` event with the same `event_id` is re-received from Kafka
- **THEN** the system recognizes the duplicate event, skips re-incrementing the variant stock, and safely acknowledges the message offset.
