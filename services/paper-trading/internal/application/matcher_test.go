package application

import (
	"testing"
	"time"

	"github.com/Soonogo/ai-trading-os/services/common/sim"
)

func TestMatchOrder(t *testing.T) {
	eng := NewMatchingEngine()
	candle := sim.Candle{Symbol: "AAPL",
		Time:   time.Now(),
		Open:   100,
		High:   101,
		Low:    99,
		Close:  100,
		Volume: 1000,
	}

	eng.SetCandle(candle)
	order := sim.Order{OrderID: "o1", Symbol: "AAPL", Side: sim.SideBuy, Quantity: 10}
	fill, err := eng.MatchOrder(order)
	if err != nil {
		t.Fatalf("match: %v", err)
	}
	if fill.Price != 100.05 {
		t.Fatalf("expected buy fill price 100.05, got %f", fill.Price)
	}
}
