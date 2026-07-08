package sim

import (
	"fmt"
	"time"
)

// Position tracks an open position.
type Position struct {
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Cost     float64 `json:"cost"`
}

// PortfolioState holds cash and positions.
type PortfolioState struct {
	Cash      float64              `json:"cash"`
	Positions map[string]*Position `json:"positions"`
	Trades    []Trade              `json:"trades"`
	Equity    []EquityPoint        `json:"equity"`
}

// Trade records a closed transaction.
type Trade struct {
	Symbol    string    `json:"symbol"`
	Side      Side      `json:"side"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	Pnl       float64   `json:"pnl"`
	Timestamp time.Time `json:"timestamp"`
}

// EquityPoint is a portfolio value snapshot.
type EquityPoint struct {
	Time     time.Time `json:"time"`
	Value    float64   `json:"value"`
	Cash     float64   `json:"cash"`
	Invested float64   `json:"invested"`
}

// NewPortfolioState creates an initial portfolio.
func NewPortfolioState(initialCash float64) *PortfolioState {
	return &PortfolioState{
		Cash:      initialCash,
		Positions: make(map[string]*Position),
		Trades:    make([]Trade, 0),
		Equity:    make([]EquityPoint, 0),
	}
}

// ApplyFill updates the portfolio from an execution fill.
func (p *PortfolioState) ApplyFill(fill Fill) error {
	switch fill.Side {
	case SideBuy:
		cost := fill.Quantity * fill.Price
		if cost > p.Cash {
			return fmt.Errorf("insufficient cash for %s buy: need %.2f have %.2f", fill.Symbol, cost, p.Cash)
		}
		p.Cash -= cost
		pos, ok := p.Positions[fill.Symbol]
		if !ok {
			p.Positions[fill.Symbol] = &Position{Symbol: fill.Symbol, Quantity: fill.Quantity, Cost: cost}
		} else {
			totalCost := pos.Cost + cost
			totalQty := pos.Quantity + fill.Quantity
			pos.Cost = totalCost
			pos.Quantity = totalQty
		}
	case SideSell:
		pos, ok := p.Positions[fill.Symbol]
		if !ok || pos.Quantity < fill.Quantity {
			return fmt.Errorf("insufficient position for %s sell: have %.4f want %.4f", fill.Symbol, pos.Quantity, fill.Quantity)
		}
		avgCost := pos.Cost / pos.Quantity
		realized := (fill.Price - avgCost) * fill.Quantity
		pos.Cost -= avgCost * fill.Quantity
		pos.Quantity -= fill.Quantity
		proceeds := fill.Quantity * fill.Price
		p.Cash += proceeds
		p.Trades = append(p.Trades, Trade{
			Symbol:    fill.Symbol,
			Side:      fill.Side,
			Quantity:  fill.Quantity,
			Price:     fill.Price,
			Pnl:       realized,
			Timestamp: fill.Timestamp,
		})
		if pos.Quantity == 0 {
			delete(p.Positions, fill.Symbol)
		}
	default:
		return fmt.Errorf("unknown side: %s", fill.Side)
	}
	return nil
}

// TotalValue returns cash plus unrealized PnL at the given mark price map.
func (p *PortfolioState) TotalValue(prices map[string]float64) float64 {
	value := p.Cash
	for symbol, pos := range p.Positions {
		price, ok := prices[symbol]
		if !ok {
			continue
		}
		value += pos.Quantity * price
	}
	return value
}

// Snapshot records an equity point for the current prices.
func (p *PortfolioState) Snapshot(t time.Time, prices map[string]float64) {
	invested := 0.0
	for symbol, pos := range p.Positions {
		price, ok := prices[symbol]
		if !ok {
			continue
		}
		invested += pos.Quantity * price
	}
	p.Equity = append(p.Equity, EquityPoint{
		Time:     t,
		Value:    p.Cash + invested,
		Cash:     p.Cash,
		Invested: invested,
	})
}
