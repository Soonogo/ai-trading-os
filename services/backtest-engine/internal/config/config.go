package config

import (
	"fmt"
	"os"
)

// Config holds backtest-engine settings.
type Config struct {
	NATSURL  string
	HTTPAddr string
}

// Load reads configuration from environment.
func Load() (Config, error) {
	cfg := Config{
		NATSURL:  getEnvOr("NATS_URL", "nats://localhost:4222"),
		HTTPAddr: getEnvOr("HTTP_ADDR", ":8080"),
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
