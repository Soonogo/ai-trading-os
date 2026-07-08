package config

import (
	"fmt"
	"os"
)

// Config holds paper-trading settings.
type Config struct {
	NATSURL string
}

// Load reads configuration from environment.
func Load() (Config, error) {
	cfg := Config{
		NATSURL: getEnvOr("NATS_URL", "nats://localhost:4222"),
	}
	if cfg.NATSURL == "" {
		return Config{}, fmt.Errorf("NATS_URL is required")
	}
	return cfg, nil
}

func getEnvOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
