package domain

import "time"

// PortfolioDecision is the output of the portfolio manager.
type PortfolioDecision struct {
	TraceID    string    `json:"trace_id"`
	Symbol     string    `json:"symbol"`
	Action     string    `json:"action"`
	Quantity   float64   `json:"quantity"`
	Confidence float64   `json:"confidence"`
	Rationale  string    `json:"rationale"`
	Timestamp  time.Time `json:"timestamp"`
}

// RiskBudget caps exposure per symbol.
type RiskBudget struct {
	MaxPositionSize float64
	MaxDrawdownPct  float64
}
