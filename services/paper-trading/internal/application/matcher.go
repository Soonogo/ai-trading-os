package application

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Soonogo/ai-trading-os/services/common/sim"
)

// MatchingEngine simulates paper exchange fills against market candles.
type MatchingEngine struct {
	cfg        sim.Config
	lastCandle map[string]sim.Candle
	cash       float64
	positions  map[string]*sim.Position
}

// NewMatchingEngine creates a paper matching engine with default settings.
func NewMatchingEngine() *MatchingEngine {
	cfg := sim.DefaultConfig()
	return &MatchingEngine{
		cfg:        cfg,
		lastCandle: make(map[string]sim.Candle),
		cash:       cfg.InitialCash,
		positions:  make(map[string]*sim.Position),
	}
}

// SetCandle updates the latest market data for a symbol.
func (m *MatchingEngine) SetCandle(candle sim.Candle) {
	m.lastCandle[candle.Symbol] = candle
}

// MatchOrder simulates an immediate market fill on the current candle.
func (m *MatchingEngine) MatchOrder(order sim.Order) (sim.Fill, error) {
	candle, ok := m.lastCandle[order.Symbol]
	if !ok {
		return sim.Fill{}, fmt.Errorf("no market data for %s", order.Symbol)
	}
	engine := sim.NewEngine(m.cfg)
	fill, err := engine.MatchMarketOrder(order, candle)
	if err != nil {
		return sim.Fill{}, err
	}
	return fill, nil
}

// ApplyFill updates the paper portfolio.
func (m *MatchingEngine) ApplyFill(fill sim.Fill) error {
	portfolio := sim.NewPortfolioState(m.cash)
	portfolio.Positions = m.positions
	if err := portfolio.ApplyFill(fill); err != nil {
		return err
	}
	m.cash = portfolio.Cash
	m.positions = portfolio.Positions
	return nil
}

// Snapshot returns the current account equity snapshot.
func (m *MatchingEngine) Snapshot() (map[string]any, error) {
	prices := make(map[string]float64)
	for sym, c := range m.lastCandle {
		prices[sym] = c.Close
	}
	portfolio := sim.NewPortfolioState(m.cash)
	portfolio.Positions = m.positions
	snap := map[string]any{
		"cash":        m.cash,
		"equity":      portfolio.TotalValue(prices),
		"positions":   m.positions,
		"last_candle": m.lastCandle,
	}
	return snap, nil
}

// ParseOrder unmarshals a raw order payload.
func ParseOrder(payload json.RawMessage) (sim.Order, error) {
	var order sim.Order
	if err := json.Unmarshal(payload, &order); err != nil {
		return sim.Order{}, err
	}
	if order.Timestamp.IsZero() {
		order.Timestamp = time.Now().UTC()
	}
	return order, nil
}

// ParseCandle unmarshals a market data candle payload.
func ParseCandle(payload json.RawMessage) (sim.Candle, error) {
	var candle sim.Candle
	if err := json.Unmarshal(payload, &candle); err != nil {
		return sim.Candle{}, err
	}
	return candle, nil
}
