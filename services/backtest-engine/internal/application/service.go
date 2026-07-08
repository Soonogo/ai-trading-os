package application

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Soonogo/ai-trading-os/services/backtest-engine/internal/domain"
	cd "github.com/Soonogo/ai-trading-os/services/common/domain"
	"github.com/Soonogo/ai-trading-os/services/common/sim"
)

// Service runs backtest and paper simulations.
type Service struct {
	engine *sim.Engine
}

// NewService creates a backtest service with default configuration.
func NewService() *Service {
	return &Service{engine: sim.NewEngine(sim.DefaultConfig())}
}

// Run executes a backtest request and returns a domain result payload.
func (s *Service) Run(ctx context.Context, req domain.BacktestRequest) (map[string]any, error) {
	cfg := req.Config
	if cfg.InitialCash == 0 {
		cfg = sim.DefaultConfig()
	}
	engine := sim.NewEngine(cfg)
	strategy, err := ResolveStrategy(req.Strategy, req.Parameters)
	if err != nil {
		return nil, err
	}
	res := engine.RunBacktest(req.Symbol, req.Candles, strategy)
	payload, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("marshal result: %w", err)
	}
	var out map[string]any
	if err := json.Unmarshal(payload, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ToEvent wraps a backtest result into a domain event.
func (s *Service) ToEvent(traceID string, result map[string]any) (cd.Event, error) {
	payload, err := json.Marshal(result)
	if err != nil {
		return cd.Event{}, err
	}
	return cd.NewEvent(traceID, cd.EventTypeBacktestResult, "backtest-engine", payload)
}
