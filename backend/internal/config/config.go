package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string

	// AI / Harness
	AnthropicKey string
	AIModel      string

	// Storage
	MediaStorage string // "local" or "s3"
	S3Bucket     string
	S3Region     string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:         envOr("PORT", "8080"),
		DatabaseURL:  envOr("DATABASE_URL", "postgres://hada:hada@localhost:5432/hada?sslmode=disable"),
		RedisURL:     envOr("REDIS_URL", "redis://localhost:6379"),
		AnthropicKey: os.Getenv("ANTHROPIC_API_KEY"),
		AIModel:      envOr("AI_MODEL", "claude-sonnet-4-6"),
		MediaStorage: envOr("MEDIA_STORAGE", "local"),
		S3Bucket:     os.Getenv("S3_BUCKET"),
		S3Region:     os.Getenv("S3_REGION"),
	}

	if cfg.AnthropicKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is required")
	}

	return cfg, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
