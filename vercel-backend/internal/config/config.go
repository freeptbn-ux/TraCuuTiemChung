package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	X_API_KEY       string
	PORT            string
	PORTAL_USERNAME string
	PORTAL_PASSWORD string
}

func LoadConfig() *Config {
	// Load .env if exists
	_ = godotenv.Load()

	return &Config{
		X_API_KEY:       getEnv("X_API_KEY", ""),
		PORT:            getEnv("PORT", "8080"),
		PORTAL_USERNAME: getEnv("PORTAL_USERNAME", ""),
		PORTAL_PASSWORD: getEnv("PORTAL_PASSWORD", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
