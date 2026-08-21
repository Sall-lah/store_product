package config

import (
	"os"
	"testing"
)

func TestConfigDefaults(t *testing.T) {
	os.Unsetenv("REDIS_PORT")
	os.Unsetenv("PORT")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	if cfg.RedisPort != "6379" {
		t.Errorf("expected default Redis port to be 6379, got %s", cfg.RedisPort)
	}

	if cfg.RedisAddr() != "localhost:6379" {
		t.Errorf("expected RedisAddr to be localhost:6379, got %s", cfg.RedisAddr())
	}

	if cfg.Port != "8080" {
		t.Errorf("expected default Port to be 8080, got %s", cfg.Port)
	}

	if len(cfg.KafkaBrokers) != 1 || cfg.KafkaBrokers[0] != "localhost:9092" {
		t.Errorf("expected default KafkaBrokers to be [localhost:9092], got %v", cfg.KafkaBrokers)
	}

	if cfg.KafkaTopicOrderEvents != "order.events" {
		t.Errorf("expected default KafkaTopicOrderEvents to be order.events, got %s", cfg.KafkaTopicOrderEvents)
	}

	if cfg.KafkaConsumerGroup != "store_product_stock_worker" {
		t.Errorf("expected default KafkaConsumerGroup to be store_product_stock_worker, got %s", cfg.KafkaConsumerGroup)
	}
}

func TestConfigCustomEnv(t *testing.T) {
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("PORT", "9000")
	os.Setenv("KAFKA_BROKERS", "kafka1:9092, kafka2:9092")
	os.Setenv("KAFKA_TOPIC_ORDER_EVENTS", "custom.order.events")
	os.Setenv("KAFKA_CONSUMER_GROUP", "custom_stock_group")
	defer func() {
		os.Unsetenv("REDIS_PORT")
		os.Unsetenv("PORT")
		os.Unsetenv("KAFKA_BROKERS")
		os.Unsetenv("KAFKA_TOPIC_ORDER_EVENTS")
		os.Unsetenv("KAFKA_CONSUMER_GROUP")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	if cfg.RedisPort != "6379" {
		t.Errorf("expected Redis port 6379, got %s", cfg.RedisPort)
	}

	if cfg.Port != "9000" {
		t.Errorf("expected Port 9000, got %s", cfg.Port)
	}

	if len(cfg.KafkaBrokers) != 2 || cfg.KafkaBrokers[0] != "kafka1:9092" || cfg.KafkaBrokers[1] != "kafka2:9092" {
		t.Errorf("expected custom KafkaBrokers [kafka1:9092 kafka2:9092], got %v", cfg.KafkaBrokers)
	}

	if cfg.KafkaTopicOrderEvents != "custom.order.events" {
		t.Errorf("expected KafkaTopicOrderEvents custom.order.events, got %s", cfg.KafkaTopicOrderEvents)
	}

	if cfg.KafkaConsumerGroup != "custom_stock_group" {
		t.Errorf("expected KafkaConsumerGroup custom_stock_group, got %s", cfg.KafkaConsumerGroup)
	}
}
