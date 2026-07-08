package application

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Soonogo/ai-trading-os/services/execution-agent/internal/domain"
	"github.com/google/uuid"
)

// Decision is a portfolio decision payload.
type Decision struct {
	TraceID    string  `json:"trace_id"`
	Symbol     string  `json:"symbol"`
	Action     string  `json:"action"`
	Quantity   float64 `json:"quantity"`
	Confidence float64 `json:"confidence"`
}

// Executor simulates order slicing and submission.
type Executor struct {
	paper bool
}

// NewExecutor creates an executor.
func NewExecutor(paper bool) *Executor {
	return &Executor{paper: paper}
}

// Execute turns a decision into an order and a simulated fill.
func (e *Executor) Execute(decision Decision) (domain.Order, domain.Fill, error) {
	if decision.Symbol == "" {
		return domain.Order{}, domain.Fill{}, fmt.Errorf("symbol required")
	}
	if decision.Quantity <= 0 {
		return domain.Order{}, domain.Fill{}, fmt.Errorf("quantity must be positive")
	}
	side := "buy"
	if decision.Action == "sell" {
		side = "sell"
	}

	order := domain.Order{
		TraceID:   decision.TraceID,
		OrderID:   uuid.NewString(),
		Symbol:    decision.Symbol,
		Side:      side,
		Quantity:  decision.Quantity,
		OrderType: "market",
		Status:    "filled",
		Paper:     e.paper,
		Timestamp: time.Now().UTC(),
	}

	fill := domain.Fill{
		OrderID:   order.OrderID,
		Symbol:    order.Symbol,
		Quantity:  order.Quantity,
		Price:     100.0, // simulated market price
		Timestamp: time.Now().UTC(),
	}
	return order, fill, nil
}

// ParseDecision parses a portfolio decision payload.
func ParseDecision(payload json.RawMessage) (Decision, error) {
	var d Decision
	if err := json.Unmarshal(payload, &d); err != nil {
		return Decision{}, err
	}
	return d, nil
}
