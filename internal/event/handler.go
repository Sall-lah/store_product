package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/Sall-lah/store_product/internal/cache"
	"github.com/Sall-lah/store_product/internal/repository"
)

// StockEventHandler coordinates inventory adjustments and cache purges driven by async order events.
type StockEventHandler struct {
	variantRepo *repository.VariantRepository
	eventRepo   *repository.EventRepository
	productRepo *repository.ProductRepository
	cacheClient *cache.Client
}

// NewStockEventHandler constructs an event handler instance with required database and cache dependencies.
func NewStockEventHandler(
	variantRepo *repository.VariantRepository,
	eventRepo *repository.EventRepository,
	productRepo *repository.ProductRepository,
	cacheClient *cache.Client,
) *StockEventHandler {
	return &StockEventHandler{
		variantRepo: variantRepo,
		eventRepo:   eventRepo,
		productRepo: productRepo,
		cacheClient: cacheClient,
	}
}

// ProcessMessage decodes an order event payload, verifies idempotency, applies inventory mutations, and flushes cache keys.
func (h *StockEventHandler) ProcessMessage(ctx context.Context, msgBytes []byte) error {
	var event OrderEvent
	if err := json.Unmarshal(msgBytes, &event); err != nil {
		log.Printf("[WARN] Failed to unmarshal message payload from Kafka: %v", err)
		// Return nil for unparseable poison messages to allow consumer offset commit
		return nil
	}

	eventID := event.GetEventID()
	eventType := event.GetEventType()
	orderID := event.GetOrderID()
	items := event.GetItems()

	if eventType == "" || len(items) == 0 {
		// Event does not contain inventory actionable items or recognized type
		return nil
	}

	// 1. Filter out non-inventory event types before querying database
	switch eventType {
	case EventTypeOrderCreated, EventTypeOrderPlaced, EventTypeOrderCancelled, EventTypeOrderCanceled, EventTypeOrderExpired:
		// Actionable inventory events
	default:
		// Other event types (e.g. order.shipped, order.completed) do not affect variant stock
		return nil
	}

	// 2. Idempotency Check to prevent duplicate stock updates from at-least-once deliveries
	if eventID != "" && h.eventRepo != nil {
		processed, err := h.eventRepo.IsEventProcessed(ctx, eventID)
		if err != nil {
			return fmt.Errorf("failed checking event idempotency: %w", err)
		}
		if processed {
			log.Printf("[INFO] Skipping duplicate event %s (Type: %s, Order: %s)", eventID, eventType, orderID)
			return nil
		}
	}

	affectedProductIDs := make(map[string]struct{})

	// 2. Dispatch stock modifications based on order lifecycle stage
	switch eventType {
	case EventTypeOrderCreated, EventTypeOrderPlaced:
		if h.variantRepo != nil {
			for _, item := range items {
				variant, err := h.variantRepo.DecrementStock(ctx, item.VariantID, item.Quantity)
				if err != nil {
					log.Printf("[ERROR] Failed to decrement stock for variant %s (Order: %s, Qty: %d): %v",
						item.VariantID, orderID, item.Quantity, err)
					if errors.Is(err, repository.ErrInsufficientStock) {
						// In production, publish a compensating 'stock.reservation_failed' event here
						continue
					}
					return fmt.Errorf("failed decrementing variant stock: %w", err)
				}
				if variant != nil && variant.ProductID != "" {
					affectedProductIDs[variant.ProductID] = struct{}{}
				}
			}
		}

	case EventTypeOrderCancelled, EventTypeOrderCanceled, EventTypeOrderExpired:
		if h.variantRepo != nil {
			for _, item := range items {
				variant, err := h.variantRepo.IncrementStock(ctx, item.VariantID, item.Quantity)
				if err != nil {
					log.Printf("[ERROR] Failed to restock variant %s (Order: %s, Qty: %d): %v",
						item.VariantID, orderID, item.Quantity, err)
					return fmt.Errorf("failed restocking variant: %w", err)
				}
				if variant != nil && variant.ProductID != "" {
					affectedProductIDs[variant.ProductID] = struct{}{}
				}
			}
		}

	default:
		// Other event types (e.g. order.shipped, order.completed) do not require variant quantity adjustments
		return nil
	}

	// 3. Purge Redis caches for affected parent products and catalog search lists
	h.invalidateAffectedProducts(ctx, affectedProductIDs)

	// 4. Mark event as processed in the idempotency ledger
	if eventID != "" && h.eventRepo != nil {
		var optOrderID *string
		if orderID != "" {
			optOrderID = &orderID
		}
		if err := h.eventRepo.MarkEventProcessed(ctx, eventID, eventType, optOrderID); err != nil {
			log.Printf("[WARN] Failed marking event %s as processed: %v", eventID, err)
		}
	}

	log.Printf("[INFO] Successfully processed %s for order %s (%d items)", eventType, orderID, len(items))
	return nil
}

// invalidateAffectedProducts purges cached product detail and list records in Redis to maintain cache consistency.
func (h *StockEventHandler) invalidateAffectedProducts(ctx context.Context, productIDs map[string]struct{}) {
	if h.productRepo != nil && h.cacheClient != nil {
		for productID := range productIDs {
			product, err := h.productRepo.GetProductByID(ctx, productID)
			if err == nil && product != nil {
				h.cacheClient.Del(ctx, cache.ProductDetailKey(product.ID), cache.ProductSlugKey(product.Slug))
			}
		}
	}

	if len(productIDs) > 0 && h.cacheClient != nil {
		h.cacheClient.DelPattern(ctx, "product:list:*")
	}
}
