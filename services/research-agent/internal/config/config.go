package config

import (
	"os"
)

type Config struct {
	NATSURL      string
	QdrantURL    string
	OpenAIKey    string
	HTTPAddr     string
}

func Load() (*Config, error) {
	return &Config{
		NATSURL:   getEnv("NATS_URL", "nats://localhost:4222"),
		QdrantURL: getEnv("QDRANT_URL", "http://localhost:6333"),
		OpenAIKey: getEnv("OPENAI_API_KEY", ""),
		HTTPAddr:  getEnv("HTTP_ADDR", ":8083"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
