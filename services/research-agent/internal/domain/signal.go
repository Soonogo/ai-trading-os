package domain

import "time"

// AssetClass is the type of market being analyzed.
type AssetClass string

const (
	AssetStock    AssetClass = "stock"
	AssetCrypto   AssetClass = "crypto"
	AssetForex    AssetClass = "forex"
	AssetCommodity AssetClass = "commodity"
)

// ResearchSignal is a domain event payload produced by the Research Agent.
type ResearchSignal struct {
	TraceID      string     `json:"trace_id"`
	Symbols      []string   `json:"symbols"`
	AssetClass   AssetClass `json:"asset_class"`
	Direction    string     `json:"direction"`
	Confidence   float64    `json:"confidence"`
	Rationale    string     `json:"rationale"`
	TimeHorizon  string     `json:"time_horizon"`
	Sources      []string   `json:"sources"`
	Timestamp    time.Time  `json:"timestamp"`
}

// Prompt is the agent's query.
type Prompt struct {
	TraceID string
	Query   string
}
