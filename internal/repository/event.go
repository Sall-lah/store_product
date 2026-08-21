package repository

import (
	"context"
	"errors"

	"github.com/Sall-lah/store_product/internal/db"
)

// EventRepository manages idempotency records to prevent duplicate event consumption.
type EventRepository struct {
	client *db.PrismaClient
}

// NewEventRepository creates a new instance of EventRepository.
func NewEventRepository(client *db.PrismaClient) *EventRepository {
	return &EventRepository{client: client}
}

// IsEventProcessed checks if an event with the given ID has already been recorded in the database.
func (r *EventRepository) IsEventProcessed(ctx context.Context, eventID string) (bool, error) {
	_, err := r.client.ProcessedEvent.FindUnique(
		db.ProcessedEvent.EventID.Equals(eventID),
	).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// MarkEventProcessed records a processed event to enforce at-least-once message deduplication.
func (r *EventRepository) MarkEventProcessed(ctx context.Context, eventID, eventType string, orderID *string) error {
	var optionalParams []db.ProcessedEventSetParam
	if orderID != nil && *orderID != "" {
		optionalParams = append(optionalParams, db.ProcessedEvent.OrderID.Set(*orderID))
	}

	_, err := r.client.ProcessedEvent.CreateOne(
		db.ProcessedEvent.EventID.Set(eventID),
		db.ProcessedEvent.EventType.Set(eventType),
		optionalParams...,
	).Exec(ctx)

	if err != nil {
		if _, isUnique := db.IsErrUniqueConstraint(err); isUnique {
			// Already recorded by a concurrent consumer worker
			return nil
		}
		return err
	}

	return nil
}
