package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config encapsulates runtime configuration loaded from environment variables.
type Config struct {
	Port         string
	Environment  string
	DatabaseURL  string
	RedisHost    string
	RedisPort    string
	RedisPassword string
	RedisDB      int

	// Rate limiting limits in requests per minute
	RateLimitPublicRPM int
	RateLimitSearchRPM int
	RateLimitAdminRPM  int

	// Kafka cluster connection and consumer parameters for asynchronous event processing
	KafkaBrokers           []string
	KafkaTopicOrderEvents  string
	KafkaConsumerGroup     string
}

// Load reads configuration from the environment and optional .env file.
// Defaults are provided for local development convenience, particularly
// binding Redis to port 6739 as specified by the microservices dev architecture.
func Load() (*Config, error) {
	// Intentionally ignore .env loading error because in containerized/production
	// environments, configuration is injected directly via system environment variables.
	_ = godotenv.Load()

	cfg := &Config{
		Port:               getEnv("PORT", "8080"),
		Environment:        getEnv("ENV", "development"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		RedisHost:          getEnv("REDIS_HOST", "localhost"),
		RedisPort:          getEnv("REDIS_PORT", "6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		RedisDB:            getEnvAsInt("REDIS_DB", 0),
		RateLimitPublicRPM:    getEnvAsInt("RATE_LIMIT_PUBLIC_RPM", 120),
		RateLimitSearchRPM:    getEnvAsInt("RATE_LIMIT_SEARCH_RPM", 60),
		RateLimitAdminRPM:     getEnvAsInt("RATE_LIMIT_ADMIN_RPM", 30),
		KafkaBrokers:          getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
		KafkaTopicOrderEvents: getEnv("KAFKA_TOPIC_ORDER_EVENTS", "order.events"),
		KafkaConsumerGroup:    getEnv("KAFKA_CONSUMER_GROUP", "store_product_stock_worker"),
	}

	return cfg, nil
}

// RedisAddr formats the Redis host and port into a standard connection target.
// Centralizing this prevents repetitive string formatting across the cache and middleware modules.
func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

// getEnv retrieves an environment variable or returns the fallback if empty.
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

// getEnvAsInt parses an integer environment variable or returns the fallback.
func getEnvAsInt(key string, fallback int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return fallback
}

// getEnvAsSlice splits comma-delimited environment variable strings to enable cluster broker address injection.
func getEnvAsSlice(key string, fallback []string) []string {
	valStr := getEnv(key, "")
	if valStr == "" {
		return fallback
	}
	rawParts := strings.Split(valStr, ",")
	var result []string
	for _, part := range rawParts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		return fallback
	}
	return result
}
