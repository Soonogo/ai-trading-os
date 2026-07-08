package sim

import "time"

// Side is the direction of an order.
type Side string

const (
	SideBuy  Side = "buy"
	SideSell Side = "sell"
)

// Order is a simulated order.
type Order struct {
	OrderID   string    `json:"order_id"`
	TraceID   string    `json:"trace_id"`
	Symbol    string    `json:"symbol"`
	Side      Side      `json:"side"`
	Quantity  float64   `json:"quantity"`
	OrderType string    `json:"order_type"`
	Paper     bool      `json:"paper"`
	Timestamp time.Time `json:"timestamp"`
}

// Fill is the result of an order execution.
type Fill struct {
	OrderID   string    `json:"order_id"`
	TraceID   string    `json:"trace_id"`
	Symbol    string    `json:"symbol"`
	Side      Side      `json:"side"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	Slippage  float64   `json:"slippage"`
	Timestamp time.Time `json:"timestamp"`
}
