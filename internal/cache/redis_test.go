package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/Sall-lah/store_product/internal/cache"
	"github.com/Sall-lah/store_product/internal/config"
	"github.com/joho/godotenv"
)

// TestCacheKeyFormatters verifies that Redis key generator functions construct
// consistent prefixes and formats across application layers to prevent cache namespace collisions.
func TestCacheKeyFormatters(t *testing.T) {
	if got := cache.ProductDetailKey("prod-123"); got != "product:detail:id:prod-123" {
		t.Errorf("ProductDetailKey() = %v, want %v", got, "product:detail:id:prod-123")
	}

	if got := cache.ProductSlugKey("sample-item"); got != "product:detail:slug:sample-item" {
		t.Errorf("ProductSlugKey() = %v, want %v", got, "product:detail:slug:sample-item")
	}

	if got := cache.ProductListKey("hash-abc"); got != "product:list:hash-abc" {
		t.Errorf("ProductListKey() = %v, want %v", got, "product:list:hash-abc")
	}
}

// TestNilOrUnavailableClient verifies that all cache helper methods fail silently
// and safely return fallback defaults without panicking when Redis is disconnected.
func TestNilOrUnavailableClient(t *testing.T) {
	var nilClient *cache.Client
	ctx := context.Background()

	if !nilClient.IsAvailable() {
		// expected
	} else {
		t.Errorf("expected nil client IsAvailable to be false")
	}
	if nilClient.Underlying() != nil {
		t.Errorf("expected nil client Underlying to be nil")
	}

	var dest map[string]string
	if nilClient.GetJSON(ctx, "any_key", &dest) {
		t.Errorf("expected GetJSON on nil client to return false")
	}

	// Del, SetJSON, and DelPattern should be no-ops and not panic
	nilClient.SetJSON(ctx, "any_key", map[string]string{"foo": "bar"}, time.Minute)
	nilClient.Del(ctx, "any_key")
	nilClient.DelPattern(ctx, "any_pattern*")
}

// TestRedisLiveOperations performs live integration tests against the configured
// Redis instance to validate ping, get, set, delete, and pattern scan operations.
func TestRedisLiveOperations(t *testing.T) {
	_ = godotenv.Load("../../.env", ".env")
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	client := cache.NewClient(cfg)
	if !client.IsAvailable() {
		t.Fatalf("Redis is not reachable at %s with current config", cfg.RedisAddr())
	}

	ctx := context.Background()
	testKey := "test:integration:item:99"
	type samplePayload struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	payload := samplePayload{ID: "99", Name: "Test Product"}

	// 1. Set JSON
	client.SetJSON(ctx, testKey, payload, 30*time.Second)

	// 2. Get JSON
	var retrieved samplePayload
	found := client.GetJSON(ctx, testKey, &retrieved)
	if !found {
		t.Fatalf("expected key %s to be retrieved from Redis", testKey)
	}
	if retrieved.ID != payload.ID || retrieved.Name != payload.Name {
		t.Fatalf("retrieved value %+v did not match expected %+v", retrieved, payload)
	}

	// 3. Del key
	client.Del(ctx, testKey)
	var afterDelete samplePayload
	if client.GetJSON(ctx, testKey, &afterDelete) {
		t.Fatalf("expected key %s to be deleted from Redis", testKey)
	}

	// 4. Pattern deletion test
	pKey1 := "test:integration:pattern:1"
	pKey2 := "test:integration:pattern:2"
	client.SetJSON(ctx, pKey1, payload, 30*time.Second)
	client.SetJSON(ctx, pKey2, payload, 30*time.Second)

	client.DelPattern(ctx, "test:integration:pattern:*")

	if client.GetJSON(ctx, pKey1, &retrieved) || client.GetJSON(ctx, pKey2, &retrieved) {
		t.Fatalf("expected all pattern keys to be deleted by DelPattern")
	}
}
