package config

import (
	"os"
)

type Config struct {
	NATSURL      string
	PaperTrading bool
}

func Load() (*Config, error) {
	return &Config{
		NATSURL:      getEnv("NATS_URL", "nats://localhost:4222"),
		PaperTrading: getEnv("PAPER_TRADING", "true") == "true",
	}, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
