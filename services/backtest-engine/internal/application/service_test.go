package application

import (
	"context"
	"testing"
	"time"

	"github.com/Soonogo/ai-trading-os/services/backtest-engine/internal/domain"
	"github.com/Soonogo/ai-trading-os/services/common/sim"
)

func makeCandles(prices []float64) []sim.Candle {
	out := make([]sim.Candle, len(prices))
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i, p := range prices {
		out[i] = sim.Candle{
			Time:   base.AddDate(0, 0, i),
			Open:   p,
			High:   p * 1.01,
			Low:    p * 0.99,
			Close:  p,
			Volume: 1000,
		}
	}
	return out
}

func TestServiceRunBacktest(t *testing.T) {
	prices := []float64{100, 101, 102, 103, 104, 105}
	req := domain.BacktestRequest{
		TraceID:  "t1",
		Symbol:   "AAPL",
		Strategy: "buy-and-hold",
		Config:   sim.DefaultConfig(),
		Candles:  makeCandles(prices),
	}
	svc := NewService()
	res, err := svc.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("run backtest: %v", err)
	}
	if res["error"] != nil && res["error"] != "" {
		t.Fatalf("unexpected error in result: %v", res["error"])
	}
	if res["final_equity"].(float64) == 0 {
		t.Fatalf("expected nonzero final equity")
	}
}
