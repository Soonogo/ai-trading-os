package config

import (
	"os"
)

type Config struct {
	NATSURL string
	HTTPAddr string
}

func Load() (*Config, error) {
	return &Config{
		NATSURL:  getEnv("NATS_URL", "nats://localhost:4222"),
		HTTPAddr: getEnv("HTTP_ADDR", ":8084"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
