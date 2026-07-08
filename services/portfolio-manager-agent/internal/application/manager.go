package application

import (
	"encoding/json"
	"fmt"
	"time"

	cd "github.com/Soonogo/ai-trading-os/services/common/domain"
	"github.com/Soonogo/ai-trading-os/services/portfolio-manager-agent/internal/domain"
)

// Signal represents a generic research signal.
type Signal struct {
	TraceID      string  `json:"trace_id"`
	Symbols      []string `json:"symbols"`
	Direction    string  `json:"direction"`
	Confidence   float64 `json:"confidence"`
	Rationale    string  `json:"rationale"`
	TimeHorizon  string  `json:"time_horizon"`
}

// Manager converts signals into portfolio decisions with sizing and risk checks.
type Manager struct {
	budget domain.RiskBudget
}

// NewManager creates a portfolio manager.
func NewManager(budget domain.RiskBudget) *Manager {
	return &Manager{budget: budget}
}

// Decide creates a portfolio decision from a signal and model-router result.
func (m *Manager) Decide(signal Signal, modelOutput string, traceID string) (domain.PortfolioDecision, *cd.Reasoning, error) {
	if len(signal.Symbols) == 0 {
		return domain.PortfolioDecision{}, nil, fmt.Errorf("no symbols in signal")
	}
	if signal.Confidence < 0.5 {
		return domain.PortfolioDecision{}, nil, fmt.Errorf("confidence below threshold")
	}
	symbol := signal.Symbols[0]
	action := "buy"
	if signal.Direction == "short" {
		action = "sell"
	}
	qty := m.budget.MaxPositionSize * signal.Confidence
	if qty > m.budget.MaxPositionSize {
		qty = m.budget.MaxPositionSize
	}

	decision := domain.PortfolioDecision{
		TraceID:    traceID,
		Symbol:     symbol,
		Action:     action,
		Quantity:   qty,
		Confidence: signal.Confidence,
		Rationale:  fmt.Sprintf("%s | model: %s", signal.Rationale, modelOutput),
		Timestamp:  time.Now().UTC(),
	}

	reasoning := &cd.Reasoning{
		Confidence: signal.Confidence,
		Evidence:   signal.Symbols,
		Notes:      decision.Rationale,
	}
	return decision, reasoning, nil
}

// ParseSignal parses a generic signal JSON payload.
func ParseSignal(payload json.RawMessage) (Signal, error) {
	var s Signal
	if err := json.Unmarshal(payload, &s); err != nil {
		return Signal{}, err
	}
	return s, nil
}
