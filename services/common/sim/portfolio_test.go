package sim

import (
	"testing"
	"time"
)

func TestPortfolioApplyFill(t *testing.T) {
	p := NewPortfolioState(10000)
	fill := Fill{OrderID: "o1", Symbol: "AAPL", Side: SideBuy, Quantity: 10, Price: 100, Timestamp: time.Now()}
	if err := p.ApplyFill(fill); err != nil {
		t.Fatalf("buy fill: %v", err)
	}
	if p.Cash != 9000 {
		t.Fatalf("expected cash 9000, got %f", p.Cash)
	}
	sell := Fill{OrderID: "o2", Symbol: "AAPL", Side: SideSell, Quantity: 10, Price: 110, Timestamp: time.Now()}
	if err := p.ApplyFill(sell); err != nil {
		t.Fatalf("sell fill: %v", err)
	}
	if p.Cash != 10100 {
		t.Fatalf("expected cash 10100, got %f", p.Cash)
	}
}
