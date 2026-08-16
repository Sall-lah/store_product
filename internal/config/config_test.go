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
}

func TestConfigCustomEnv(t *testing.T) {
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("PORT", "9000")
	defer func() {
		os.Unsetenv("REDIS_PORT")
		os.Unsetenv("PORT")
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
}
