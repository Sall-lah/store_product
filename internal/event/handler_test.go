package event

import (
	"context"
	"testing"
)

func TestStockEventHandler_InvalidJSON(t *testing.T) {
	handler := &StockEventHandler{}
	ctx := context.Background()

	err := handler.ProcessMessage(ctx, []byte("invalid json string"))
	if err != nil {
		t.Errorf("expected nil error on poison pill invalid JSON, got: %v", err)
	}
}

func TestStockEventHandler_EmptyItems(t *testing.T) {
	handler := &StockEventHandler{}
	ctx := context.Background()

	err := handler.ProcessMessage(ctx, []byte(`{"event_id":"evt-1","event_type":"order.created","items":[]}`))
	if err != nil {
		t.Errorf("expected nil error on empty items, got: %v", err)
	}
}

func TestStockEventHandler_IgnoredEventType(t *testing.T) {
	handler := &StockEventHandler{}
	ctx := context.Background()

	err := handler.ProcessMessage(ctx, []byte(`{"event_id":"evt-2","event_type":"order.delivered","items":[{"variant_id":"v-1","quantity":1}]}`))
	if err != nil {
		t.Errorf("expected nil error on ignored event type, got: %v", err)
	}
}

func TestStockEventHandler_OrderCreated_NilRepo(t *testing.T) {
	handler := &StockEventHandler{}
	ctx := context.Background()

	err := handler.ProcessMessage(ctx, []byte(`{"event_id":"evt-create-1","event_type":"order.created","items":[{"variant_id":"v-1","quantity":2}]}`))
	if err != nil {
		t.Errorf("expected nil error when repo is nil in tests, got: %v", err)
	}
}

func TestStockEventHandler_OrderCancelled_NilRepo(t *testing.T) {
	handler := &StockEventHandler{}
	ctx := context.Background()

	err := handler.ProcessMessage(ctx, []byte(`{"event_id":"evt-cancel-1","event_type":"order.cancelled","items":[{"variant_id":"v-1","quantity":2}]}`))
	if err != nil {
		t.Errorf("expected nil error when repo is nil in tests, got: %v", err)
	}
}

func TestStockEventHandler_OrderExpired_NilRepo(t *testing.T) {
	handler := &StockEventHandler{}
	ctx := context.Background()

	err := handler.ProcessMessage(ctx, []byte(`{"event_id":"evt-exp-1","event_type":"order.expired","items":[{"variant_id":"v-1","quantity":2}]}`))
	if err != nil {
		t.Errorf("expected nil error when repo is nil in tests, got: %v", err)
	}
}

