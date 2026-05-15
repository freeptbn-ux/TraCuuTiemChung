package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Mock environment variables
	os.Setenv("X_API_KEY", "test-key")
	os.Setenv("PORT", "9090")
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	defer os.Unsetenv("X_API_KEY")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("REDIS_URL")

	cfg := LoadConfig()

	if cfg.X_API_KEY != "test-key" {
		t.Errorf("Expected X_API_KEY to be 'test-key', got '%s'", cfg.X_API_KEY)
	}

	if cfg.PORT != "9090" {
		t.Errorf("Expected PORT to be '9090', got '%s'", cfg.PORT)
	}

	if cfg.Redis.URL != "redis://localhost:6379" {
		t.Errorf("Expected Redis.URL to be 'redis://localhost:6379', got '%s'", cfg.Redis.URL)
	}
}

func TestDefaultConfig(t *testing.T) {
	// Ensure env is clean
	os.Unsetenv("X_API_KEY")
	os.Unsetenv("PORT")

	cfg := LoadConfig()

	if cfg.X_API_KEY != "" {
		t.Errorf("Expected default X_API_KEY to be empty, got '%s'", cfg.X_API_KEY)
	}

	if cfg.PORT != "8080" {
		t.Errorf("Expected default PORT to be '8080', got '%s'", cfg.PORT)
	}
}
