package domain

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Provider is an LLM provider identifier.
type Provider string

const (
	ProviderOpenAI    Provider = "openai"
	ProviderAnthropic Provider = "anthropic"
	ProviderGroq      Provider = "groq"
)

// Prompt is a user or agent prompt.
type Prompt struct {
	TraceID string
	Content string
	System  string
	Format  string // e.g. "json" or "text"
}

// Validate ensures the prompt is well-formed.
func (p Prompt) Validate() error {
	if p.TraceID == "" {
		return fmt.Errorf("trace_id is required")
	}
	if strings.TrimSpace(p.Content) == "" {
		return fmt.Errorf("prompt content is required")
	}
	return nil
}

// ModelCall captures a single LLM invocation result.
type ModelCall struct {
	Provider  Provider
	Model     string
	Output    string
	Latency   time.Duration
	Error     string
	Timestamp time.Time
}

// ConfidenceScore is a normalized score from 0 to 1.
type ConfidenceScore struct {
	Provider Provider
	Model    string
	Score    float64
	Reason   string
}

// EnsembleResult is the output of the model router.
type EnsembleResult struct {
	TraceID      string
	Selected     ModelCall
	Confidence   float64
	AllCalls     []ModelCall
	Scores       []ConfidenceScore
	Reasoning    string
	Timestamp    time.Time
}

// LLMClient is the port for calling an LLM provider.
type LLMClient interface {
	Complete(ctx context.Context, provider Provider, model, system, prompt string) (string, error)
}
