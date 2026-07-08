package sim

import (
	"testing"
	"time"
)

func makeCandles(prices []float64) []Candle {
	out := make([]Candle, len(prices))
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i, p := range prices {
		out[i] = Candle{
			Symbol:   "AAPL",
			Time:     base.AddDate(0, 0, i),
			Open:     p,
			High:     p * 1.01,
			Low:      p * 0.99,
			Close:    p,
			Volume:   1000,
		}
	}
	return out
}

func TestRunBacktestBuyAndHold(t *testing.T) {
	prices := []float64{100, 101, 102, 103, 104}
	candles := makeCandles(prices)
	strategy := func(symbol string, c Candle, state *PortfolioState) *Signal {
		if c.Time.Equal(candles[0].Time) {
			qty := state.Cash / c.Open
			return &Signal{Symbol: symbol, Side: SideBuy, Quantity: qty}
		}
		return nil
	}
	en := NewEngine(DefaultConfig())
	res := en.RunBacktest("AAPL", candles, strategy)
	if res.Error != "" {
		t.Fatalf("unexpected error: %s", res.Error)
	}
	if res.FinalEquity <= 0 {
		t.Fatalf("expected positive equity, got %f", res.FinalEquity)
	}
	if len(res.Equity) != len(candles) {
		t.Fatalf("expected %d equity points, got %d", len(candles), len(res.Equity))
	}
}

func TestMatchMarketOrder(t *testing.T) {
	en := NewEngine(DefaultConfig())
	order := Order{OrderID: "o1", Symbol: "AAPL", Side: SideBuy, Quantity: 10}
	candle := Candle{Symbol: "AAPL", Open: 100, Close: 101}
	fill, err := en.MatchMarketOrder(order, candle)
	if err != nil {
		t.Fatalf("match: %v", err)
	}
	if fill.Price != 100.05 {
		t.Fatalf("expected slippage adjusted buy price 100.05, got %f", fill.Price)
	}
}
