package domain

import (
	"time"

	"github.com/Soonogo/ai-trading-os/services/common/sim"
)

// BacktestRequest is a request to run a historical simulation.
type BacktestRequest struct {
	TraceID    string             `json:"trace_id"`
	Symbol     string             `json:"symbol"`
	Strategy   string             `json:"strategy"`
	Config     sim.Config         `json:"config"`
	Candles    []sim.Candle       `json:"candles"`
	Parameters map[string]float64 `json:"parameters,omitempty"`
	Timestamp  time.Time          `json:"timestamp"`
}
