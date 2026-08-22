package event

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// Event type constants representing order lifecycle stages published across topics.
const (
	EventTypeOrderCreated   = "order.created"
	EventTypeOrderPlaced    = "order.placed"
	EventTypeOrderCancelled = "order.cancelled"
	EventTypeOrderCanceled  = "order.canceled"
	EventTypeOrderExpired   = "order.expired"
)

// OrderItem represents an individual line item variant and quantity within an order.
type OrderItem struct {
	ProductID string `json:"product_id,omitempty"`
	VariantID string `json:"variant_id"`
	SKU       string `json:"sku,omitempty"`
	Quantity  int    `json:"quantity"`
}

// OrderEventData holds domain-specific payload attributes when wrapped in an event envelope.
type OrderEventData struct {
	OrderID string      `json:"order_id"`
	UserID  string      `json:"user_id,omitempty"`
	Status  string      `json:"status,omitempty"`
	Items   []OrderItem `json:"items"`
	Reason  string      `json:"reason,omitempty"`
}

// OrderEvent represents the top-level message received from Kafka order event topics.
// It supports both envelope-wrapped payloads (`data.items`) and flattened payload structures,
// as well as flexible timestamp representations (ISO 8601 strings or Unix epoch integers).
type OrderEvent struct {
	EventID   string         `json:"event_id,omitempty"`
	ID        string         `json:"id,omitempty"`
	EventType string         `json:"event_type,omitempty"`
	Type      string         `json:"type,omitempty"`
	Timestamp interface{}    `json:"timestamp,omitempty"`
	OrderID   string         `json:"order_id,omitempty"`
	Items     []OrderItem    `json:"items,omitempty"`
	Data      OrderEventData `json:"data,omitempty"`
}

// GetEventID extracts a canonical event ID or generates a deterministic fallback to ensure idempotency.
func (e *OrderEvent) GetEventID() string {
	if e.EventID != "" {
		return e.EventID
	}
	if e.ID != "" {
		return e.ID
	}
	orderID := e.GetOrderID()
	eventType := e.GetEventType()
	if orderID != "" && eventType != "" {
		// Fallback deterministic ID based on order + event type
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s", orderID, eventType)))
		return fmt.Sprintf("gen_%s", hex.EncodeToString(hash[:12]))
	}
	return ""
}

// GetEventType extracts the normalized event type identifier across varying schema envelope conventions.
func (e *OrderEvent) GetEventType() string {
	if e.EventType != "" {
		return strings.ToLower(strings.TrimSpace(e.EventType))
	}
	if e.Type != "" {
		return strings.ToLower(strings.TrimSpace(e.Type))
	}
	return ""
}

// GetOrderID retrieves the order ID from either direct or nested envelope payload structures.
func (e *OrderEvent) GetOrderID() string {
	if e.Data.OrderID != "" {
		return e.Data.OrderID
	}
	return e.OrderID
}

// GetItems extracts the list of line items regardless of whether data was nested in an envelope.
func (e *OrderEvent) GetItems() []OrderItem {
	if len(e.Data.Items) > 0 {
		return e.Data.Items
	}
	return e.Items
}
