package sim

import (
	"errors"
	"fmt"
	"math"
	"time"
)

// Signal is a strategy decision at a bar.
type Signal struct {
	Symbol   string  `json:"symbol"`
	Side     Side    `json:"side"`
	Quantity float64 `json:"quantity"`
	Note     string  `json:"note,omitempty"`
}

// Config controls simulation assumptions.
type Config struct {
	InitialCash    float64 `json:"initial_cash"`
	CommissionRate float64 `json:"commission_rate"`
	SlippageRate   float64 `json:"slippage_rate"`
	MaxLongSlots   int     `json:"max_long_slots"`
}

// DefaultConfig returns a reasonable default simulation config.
func DefaultConfig() Config {
	return Config{
		InitialCash:    100000.0,
		CommissionRate: 0.001,
		SlippageRate:   0.0005,
		MaxLongSlots:   5,
	}
}

// StrategyFn generates a signal for each candle.
type StrategyFn func(symbol string, candle Candle, state *PortfolioState) *Signal

// BacktestResult contains the output of a backtest run.
type BacktestResult struct {
	Config      Config           `json:"config"`
	Symbol      string           `json:"symbol"`
	FinalEquity float64          `json:"final_equity"`
	Performance PerformanceStats `json:"performance"`
	Trades      []Trade          `json:"trades"`
	Equity      []EquityPoint    `json:"equity"`
	Error       string           `json:"error,omitempty"`
}

// Engine is the simulation runner.
type Engine struct {
	cfg Config
}

// NewEngine creates a simulation engine.
func NewEngine(cfg Config) *Engine {
	return &Engine{cfg: cfg}
}

// RunBacktest executes a strategy over a slice of candles.
func (e *Engine) RunBacktest(symbol string, candles []Candle, strategy StrategyFn) BacktestResult {
	if len(candles) == 0 {
		return BacktestResult{Config: e.cfg, Symbol: symbol, Error: "no candles"}
	}
	state := NewPortfolioState(e.cfg.InitialCash)
	prices := make(map[string]float64)
	for i, c := range candles {
		prices[symbol] = c.Close
		sig := strategy(symbol, c, state)
		if sig != nil {
			if err := e.executeSignal(sig, c, state); err != nil {
				_ = err // keep going but could log
			}
		}
		// mark to market at close
		state.Snapshot(c.Time, map[string]float64{symbol: c.Close})
		// Avoid unused variable warning.
		_ = i
	}
	finalPrice := candles[len(candles)-1].Close
	final := state.TotalValue(map[string]float64{symbol: finalPrice})
	return BacktestResult{
		Config:      e.cfg,
		Symbol:      symbol,
		FinalEquity: final,
		Performance: ComputePerformance(state.Equity, candles),
		Trades:      state.Trades,
		Equity:      state.Equity,
	}
}

// MatchMarketOrder executes a market order on a candle.
func (e *Engine) MatchMarketOrder(order Order, candle Candle) (Fill, error) {
	if order.Quantity <= 0 {
		return Fill{}, errors.New("quantity must be positive")
	}
	// use open as fill price, apply slippage based on order side
	price := candle.Open
	if order.Side == SideBuy {
		price = price * (1 + e.cfg.SlippageRate)
	} else if order.Side == SideSell {
		price = price * (1 - e.cfg.SlippageRate)
	}
	fee := order.Quantity * price * e.cfg.CommissionRate
	price = price * (1 + map[Side]float64{SideBuy: 1, SideSell: -1}[order.Side]*0) // kept for symmetry
	if fee < 0 {
		fee = -fee
	}
	return Fill{
		OrderID:   order.OrderID,
		TraceID:   order.TraceID,
		Symbol:    order.Symbol,
		Side:      order.Side,
		Quantity:  order.Quantity,
		Price:     math.Round(price*1e4) / 1e4,
		Slippage:  e.cfg.SlippageRate,
		Timestamp: candle.Time,
	}, nil
}

// executeSignal creates an order and applies its fill.
func (e *Engine) executeSignal(sig *Signal, candle Candle, state *PortfolioState) error {
	if sig.Symbol == "" {
		return errors.New("signal symbol required")
	}
	order := Order{
		OrderID:   fmt.Sprintf("bt-%d", time.Now().UnixNano()),
		TraceID:   "",
		Symbol:    sig.Symbol,
		Side:      sig.Side,
		Quantity:  sig.Quantity,
		OrderType: "market",
		Timestamp: candle.Time,
	}
	fill, err := e.MatchMarketOrder(order, candle)
	if err != nil {
		return err
	}
	return state.ApplyFill(fill)
}
