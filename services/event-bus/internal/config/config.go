package config

import (
	"fmt"
	"os"
)

// Config is the runtime configuration for the event bus API gateway.
type Config struct {
	HTTPAddr    string
	NATSURL     string
	ServiceName string
}

// Load reads configuration from environment with sane defaults for local dev.
func Load() (*Config, error) {
	cfg := &Config{
		HTTPAddr:    getEnv("HTTP_ADDR", ":8080"),
		NATSURL:     getEnv("NATS_URL", "nats://localhost:4222"),
		ServiceName: getEnv("SERVICE_NAME", "event-bus"),
	}
	if cfg.HTTPAddr == "" {
		return nil, fmt.Errorf("HTTP_ADDR must not be empty")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
