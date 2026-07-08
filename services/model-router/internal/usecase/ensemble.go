package usecase

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Soonogo/ai-trading-os/services/model-router/internal/domain"
)

// ModelRegistry defines the set of models to invoke per provider.
type ModelRegistry struct {
	ModelsByProvider map[domain.Provider][]string
}

// DefaultRegistry returns a sensible default multi-model registry.
func DefaultRegistry() ModelRegistry {
	return ModelRegistry{
		ModelsByProvider: map[domain.Provider][]string{
			domain.ProviderOpenAI:    {"gpt-4o-mini", "gpt-4o"},
			domain.ProviderAnthropic: {"claude-3-5-sonnet-20241022"},
			domain.ProviderGroq:      {"llama3-70b-8192"},
		},
	}
}

// Ensemble is the multi-model debate and confidence-fusion service.
type Ensemble struct {
	client   domain.LLMClient
	registry ModelRegistry
}

// NewEnsemble creates an ensemble router.
func NewEnsemble(client domain.LLMClient, registry ModelRegistry) *Ensemble {
	return &Ensemble{client: client, registry: registry}
}

// Route executes the prompt across all registered models and returns a fused result.
func (e *Ensemble) Route(ctx context.Context, prompt domain.Prompt) (domain.EnsembleResult, error) {
	if err := prompt.Validate(); err != nil {
		return domain.EnsembleResult{}, fmt.Errorf("validate prompt: %w", err)
	}

	var calls []domain.ModelCall
	var mu sync.Mutex
	var wg sync.WaitGroup

	for provider, models := range e.registry.ModelsByProvider {
		for _, model := range models {
			wg.Add(1)
			go func(p domain.Provider, m string) {
				defer wg.Done()
				start := time.Now()
				output, err := e.client.Complete(ctx, p, m, prompt.System, injectFormat(prompt))
				call := domain.ModelCall{
					Provider:  p,
					Model:     m,
					Output:    output,
					Latency:   time.Since(start),
					Timestamp: time.Now().UTC(),
				}
				if err != nil {
					call.Error = err.Error()
				}
				mu.Lock()
				calls = append(calls, call)
				mu.Unlock()
			}(provider, model)
		}
	}
	wg.Wait()

	scores := scoreCalls(calls)
	best := selectBest(calls, scores)

	result := domain.EnsembleResult{
		TraceID:    prompt.TraceID,
		Selected:   best,
		Confidence: aggregateConfidence(scores),
		AllCalls:   calls,
		Scores:     scores,
		Reasoning:  buildReasoning(calls, best, scores),
		Timestamp:  time.Now().UTC(),
	}
	return result, nil
}

func injectFormat(p domain.Prompt) string {
	if p.Format == "json" {
		return p.Content + "\n\nRespond with valid JSON and include a top-level 'confidence' float between 0 and 1."
	}
	return p.Content
}

var confidenceRe = regexp.MustCompile(`"confidence"\s*[:=]\s*(0\.\d+|1\.?0?|[\d.]+)`)

func scoreCalls(calls []domain.ModelCall) []domain.ConfidenceScore {
	scores := make([]domain.ConfidenceScore, 0, len(calls))
	for _, call := range calls {
		if call.Error != "" {
			scores = append(scores, domain.ConfidenceScore{
				Provider: call.Provider,
				Model:    call.Model,
				Score:    0,
				Reason:   "call failed: " + call.Error,
			})
			continue
		}

		raw := confidenceRe.FindStringSubmatch(call.Output)
		score := 0.5
		reason := "no explicit confidence; default 0.5"
		if len(raw) > 1 {
			if parsed, err := strconv.ParseFloat(raw[1], 64); err == nil {
				score = parsed
				reason = "parsed model-reported confidence"
			}
		}
		// Boost for shorter latency as a weak proxy for model readiness.
		if call.Latency < time.Second {
			score = min(score*1.05, 1.0)
		}
		scores = append(scores, domain.ConfidenceScore{
			Provider: call.Provider,
			Model:    call.Model,
			Score:    score,
			Reason:   reason,
		})
	}
	return scores
}

func selectBest(calls []domain.ModelCall, scores []domain.ConfidenceScore) domain.ModelCall {
	bestIndex := -1
	bestScore := -1.0
	scoreMap := make(map[string]float64)
	for _, s := range scores {
		key := string(s.Provider) + "/" + s.Model
		scoreMap[key] = s.Score
	}
	for i, call := range calls {
		key := string(call.Provider) + "/" + call.Model
		if scoreMap[key] > bestScore {
			bestScore = scoreMap[key]
			bestIndex = i
		}
	}
	if bestIndex == -1 {
		return domain.ModelCall{}
	}
	return calls[bestIndex]
}

func aggregateConfidence(scores []domain.ConfidenceScore) float64 {
	if len(scores) == 0 {
		return 0
	}
	var sum float64
	for _, s := range scores {
		sum += s.Score
	}
	return sum / float64(len(scores))
}

func buildReasoning(calls []domain.ModelCall, best domain.ModelCall, scores []domain.ConfidenceScore) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Selected %s/%s. ", best.Provider, best.Model)
	for _, s := range scores {
		fmt.Fprintf(&b, "%s/%s=%.2f (%s); ", s.Provider, s.Model, s.Score, s.Reason)
	}
	return b.String()
}
