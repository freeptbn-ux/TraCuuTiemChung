package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type RedisConfig struct {
	URL string
}

type Config struct {
	X_API_KEY       string
	PORT            string
	PORTAL_USERNAME string
	PORTAL_PASSWORD string
	Redis           RedisConfig
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	redisURL := getEnv("REDIS_URL", getEnv("UPSTASH_REDIS_URL", ""))
	
	// Tự động lắp ghép URL từ biến Upstash REST nếu thiếu REDIS_URL
	if redisURL == "" {
		restURL := getEnv("UPSTASH_REDIS_REST_URL", "")
		restToken := getEnv("UPSTASH_REDIS_REST_TOKEN", "")
		if restURL != "" && restToken != "" {
			// restURL thường có dạng: https://host
			host := strings.TrimPrefix(restURL, "https://")
			redisURL = fmt.Sprintf("rediss://default:%s@%s:6379", restToken, host)
		}
	}

	return &Config{
		X_API_KEY:       getEnv("X_API_KEY", ""),
		PORT:            getEnv("PORT", "8080"),
		PORTAL_USERNAME: getEnv("PORTAL_USERNAME", ""),
		PORTAL_PASSWORD: getEnv("PORTAL_PASSWORD", ""),
		Redis: RedisConfig{
			URL: redisURL,
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
