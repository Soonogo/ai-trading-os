package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	cd "github.com/Soonogo/ai-trading-os/services/common/domain"
	"github.com/Soonogo/ai-trading-os/services/research-agent/internal/domain"
)

// RAGClient is the port for querying the vector store.
type RAGClient interface {
	Search(ctx context.Context, query string, limit int) ([]string, error)
	Store(ctx context.Context, document string, metadata map[string]any) error
}

// LLMClient is the port for calling an LLM.
type LLMClient interface {
	Complete(ctx context.Context, system, prompt string) (string, error)
}

// Researcher coordinates RAG + LLM to produce ResearchSignals.
type Researcher struct {
	rag    RAGClient
	llm    LLMClient
}

// NewResearcher creates a Researcher.
func NewResearcher(rag RAGClient, llm LLMClient) *Researcher {
	return &Researcher{rag: rag, llm: llm}
}

// Research processes a prompt and returns a signal event payload plus reasoning.
func (r *Researcher) Research(ctx context.Context, prompt string) (domain.ResearchSignal, *cd.Reasoning, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	sources, err := r.rag.Search(ctx, prompt, 5)
	if err != nil {
		return domain.ResearchSignal{}, nil, fmt.Errorf("rag search: %w", err)
	}

	sys := `You are a quantitative research analyst. Given the user's market query and retrieved context, output JSON with:
- symbols: []string
- asset_class: "stock" | "crypto" | "forex" | "commodity"
- direction: "long" | "short" | "neutral"
- confidence: float (0-1)
- rationale: string
- time_horizon: "scalp" | "intraday" | "swing" | "position"
- sources: []string (the IDs used)`

	user := fmt.Sprintf("Query: %s\nContext:\n%s", prompt, strings.Join(sources, "\n\n"))
	out, err := r.llm.Complete(ctx, sys, user)
	if err != nil {
		return domain.ResearchSignal{}, nil, fmt.Errorf("llm complete: %w", err)
	}

	var signal domain.ResearchSignal
	if err := json.Unmarshal([]byte(out), &signal); err != nil {
		return domain.ResearchSignal{}, nil, fmt.Errorf("parse signal: %w", err)
	}
	signal.Timestamp = time.Now().UTC()

	reasoning := &cd.Reasoning{
		Confidence: signal.Confidence,
		Evidence:   sources,
		Notes:      signal.Rationale,
	}
	return signal, reasoning, nil
}
