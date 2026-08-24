## ADDED Requirements

### Requirement: Atomic Stock Adjustment for Variants
The repository and service layer SHALL provide atomic stock adjustment operations (`AdjustStock` / `DecrementStock` / `IncrementStock`) to increment or decrement variant stock safely under concurrent execution without race conditions.

#### Scenario: Concurrent stock deduction
- **WHEN** multiple order events concurrently decrement stock for the same variant
- **THEN** the database executes atomic updates (`stock = stock - delta` where `stock >= delta`), ensuring stock never becomes negative and preventing race-condition overselling.

#### Scenario: Atomic stock increment on restock
- **WHEN** an order cancellation or expiration event increments stock for a variant
- **THEN** the database executes an atomic update (`stock = stock + delta`), correctly restoring inventory.
