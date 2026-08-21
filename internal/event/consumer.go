package event

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Sall-lah/store_product/internal/config"
	"github.com/segmentio/kafka-go"
)

// Consumer manages a dedicated Kafka reader loop subscribing to order lifecycle events.
type Consumer struct {
	reader  *kafka.Reader
	handler *StockEventHandler
	cfg     *config.Config
}

// NewConsumer creates a configured Kafka consumer ready to start background message polling.
func NewConsumer(cfg *config.Config, handler *StockEventHandler) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.KafkaBrokers,
		GroupID:        cfg.KafkaConsumerGroup,
		Topic:          cfg.KafkaTopicOrderEvents,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		MaxWait:        1 * time.Second,
		CommitInterval: time.Second,
		StartOffset:    kafka.FirstOffset,
	})

	return &Consumer{
		reader:  reader,
		handler: handler,
		cfg:     cfg,
	}
}

// Start begins the polling loop in a blocking manner until the provided context is canceled.
func (c *Consumer) Start(ctx context.Context) {
	log.Printf("[INFO] Starting Kafka consumer for topic '%s' (Group: '%s', Brokers: %v)",
		c.cfg.KafkaTopicOrderEvents, c.cfg.KafkaConsumerGroup, c.cfg.KafkaBrokers)

	for {
		// FetchMessage allows manual offset management and context cancellation
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Println("[INFO] Kafka consumer worker stopped via context cancellation.")
				return
			}
			log.Printf("[WARN] Kafka consumer read error: %v", err)
			// Apply backoff to prevent tight CPU spin during broker disconnects
			select {
			case <-ctx.Done():
				return
			case <-time.After(500 * time.Millisecond):
			}
			continue
		}

		// Process message through domain handler
		if err := c.handler.ProcessMessage(ctx, msg.Value); err != nil {
			log.Printf("[ERROR] Failed processing Kafka message on topic %s at offset %d: %v",
				msg.Topic, msg.Offset, err)
		}

		// Commit offset explicitly after processing attempt
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("[WARN] Failed committing Kafka offset %d: %v", msg.Offset, err)
		}
	}
}

// Close terminates the underlying Kafka reader connection pool cleanly.
func (c *Consumer) Close() error {
	log.Println("[INFO] Closing Kafka consumer connection...")
	return c.reader.Close()
}
