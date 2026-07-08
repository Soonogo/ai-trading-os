package application

import (
	"fmt"

	"github.com/Soonogo/ai-trading-os/services/common/sim"
)

// StrategyFunc is a named strategy constructor.
type StrategyFunc func(params map[string]float64) sim.StrategyFn

// Strategies is the registry of built-in strategies.
var Strategies = map[string]StrategyFunc{
	"buy-and-hold": BuyAndHoldStrategy,
	"sma-cross":    SmaCrossStrategy,
}

// BuyAndHoldStrategy buys at the first candle and holds.
func BuyAndHoldStrategy(_ map[string]float64) sim.StrategyFn {
	bought := false
	return func(symbol string, candle sim.Candle, state *sim.PortfolioState) *sim.Signal {
		if !bought {
			bought = true
			qty := state.Cash / candle.Open
			return &sim.Signal{Symbol: symbol, Side: sim.SideBuy, Quantity: qty}
		}
		return nil
	}
}

// SmaCrossStrategy uses fast and slow simple moving averages.
func SmaCrossStrategy(params map[string]float64) sim.StrategyFn {
	fast := int(params["fast"])
	if fast == 0 {
		fast = 10
	}
	slow := int(params["slow"])
	if slow == 0 {
		slow = 30
	}
	closes := make([]float64, 0, slow)
	return func(symbol string, candle sim.Candle, state *sim.PortfolioState) *sim.Signal {
		closes = append(closes, candle.Close)
		if len(closes) < slow {
			return nil
		}
		fastSum, slowSum := 0.0, 0.0
		for i := 0; i < slow; i++ {
			idx := len(closes) - 1 - i
			slowSum += closes[idx]
			if i < fast {
				fastSum += closes[idx]
			}
		}
		fastMA := fastSum / float64(fast)
		slowMA := slowSum / float64(slow)
		qty := 0.0
		if fastMA > slowMA {
			equity := state.TotalValue(map[string]float64{symbol: candle.Close})
			if equity > 0 {
				qty = (equity * 0.2) / candle.Open
			}
			return &sim.Signal{Symbol: symbol, Side: sim.SideBuy, Quantity: qty}
		}
		if fastMA < slowMA {
			pos, ok := state.Positions[symbol]
			if ok {
				return &sim.Signal{Symbol: symbol, Side: sim.SideSell, Quantity: pos.Quantity}
			}
		}
		return nil
	}
}

// ResolveStrategy returns the named strategy or an error.
func ResolveStrategy(name string, params map[string]float64) (sim.StrategyFn, error) {
	fn, ok := Strategies[name]
	if !ok {
		return nil, fmt.Errorf("unknown strategy: %s", name)
	}
	return fn(params), nil
}
