package config

import (
	"os"
)

// Config is the runtime configuration for the model router.
type Config struct {
	HTTPAddr    string
	NATSURL     string
	OpenAIKey   string
	AnthropicKey string
	GroqKey     string
	DefaultProvider string
	FallbackProvider string
}

// Load reads configuration from environment.
func Load() (*Config, error) {
	cfg := &Config{
		HTTPAddr:         getEnv("HTTP_ADDR", ":8082"),
		NATSURL:          getEnv("NATS_URL", "nats://localhost:4222"),
		OpenAIKey:        getEnv("OPENAI_API_KEY", ""),
		AnthropicKey:     getEnv("ANTHROPIC_API_KEY", ""),
		GroqKey:          getEnv("GROQ_API_KEY", ""),
		DefaultProvider:  getEnv("MODEL_ROUTER_DEFAULT_PROVIDER", "openai"),
		FallbackProvider: getEnv("MODEL_ROUTER_FALLBACK_PROVIDER", "anthropic"),
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
