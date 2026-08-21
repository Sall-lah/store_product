package event

import (
	"encoding/json"
	"testing"
)

func TestOrderEvent_EnvelopeParsing(t *testing.T) {
	rawJSON := `{
		"event_id": "evt-12345",
		"event_type": "order.created",
		"timestamp": "2026-08-21T12:00:00Z",
		"data": {
			"order_id": "ord-999",
			"user_id": "usr-888",
			"items": [
				{
					"variant_id": "var-1",
					"sku": "TSHIRT-RED-L",
					"quantity": 2
				},
				{
					"variant_id": "var-2",
					"sku": "TSHIRT-BLU-M",
					"quantity": 1
				}
			]
		}
	}`

	var event OrderEvent
	if err := json.Unmarshal([]byte(rawJSON), &event); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if event.GetEventID() != "evt-12345" {
		t.Errorf("expected EventID evt-12345, got %s", event.GetEventID())
	}
	if event.GetEventType() != "order.created" {
		t.Errorf("expected EventType order.created, got %s", event.GetEventType())
	}
	if event.GetOrderID() != "ord-999" {
		t.Errorf("expected OrderID ord-999, got %s", event.GetOrderID())
	}
	if len(event.GetItems()) != 2 {
		t.Fatalf("expected 2 items, got %d", len(event.GetItems()))
	}
	if event.GetItems()[0].VariantID != "var-1" || event.GetItems()[0].Quantity != 2 {
		t.Errorf("unexpected first item: %+v", event.GetItems()[0])
	}
}

func TestOrderEvent_FlattenedParsing(t *testing.T) {
	rawJSON := `{
		"id": "evt-flat-555",
		"type": "order.cancelled",
		"order_id": "ord-flat-777",
		"items": [
			{
				"variant_id": "var-99",
				"quantity": 3
			}
		]
	}`

	var event OrderEvent
	if err := json.Unmarshal([]byte(rawJSON), &event); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if event.GetEventID() != "evt-flat-555" {
		t.Errorf("expected EventID evt-flat-555, got %s", event.GetEventID())
	}
	if event.GetEventType() != "order.cancelled" {
		t.Errorf("expected EventType order.cancelled, got %s", event.GetEventType())
	}
	if event.GetOrderID() != "ord-flat-777" {
		t.Errorf("expected OrderID ord-flat-777, got %s", event.GetOrderID())
	}
	if len(event.GetItems()) != 1 {
		t.Fatalf("expected 1 item, got %d", len(event.GetItems()))
	}
	if event.GetItems()[0].VariantID != "var-99" || event.GetItems()[0].Quantity != 3 {
		t.Errorf("unexpected item: %+v", event.GetItems()[0])
	}
}

func TestOrderEvent_DeterministicFallbackID(t *testing.T) {
	rawJSON := `{
		"event_type": "order.expired",
		"data": {
			"order_id": "ord-expire-100",
			"items": [
				{
					"variant_id": "var-100",
					"quantity": 1
				}
			]
		}
	}`

	var event OrderEvent
	if err := json.Unmarshal([]byte(rawJSON), &event); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	id := event.GetEventID()
	if id == "" {
		t.Fatalf("expected generated deterministic ID, got empty string")
	}

	// Verify idempotency: running twice produces identical hash
	if event.GetEventID() != id {
		t.Errorf("expected deterministic fallback ID, got different hashes")
	}
}
