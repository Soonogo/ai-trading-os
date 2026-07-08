package domain

import "time"

// Order represents a trade order.
type Order struct {
	TraceID   string    `json:"trace_id"`
	OrderID   string    `json:"order_id"`
	Symbol    string    `json:"symbol"`
	Side      string    `json:"side"`
	Quantity  float64   `json:"quantity"`
	OrderType string    `json:"order_type"`
	Status    string    `json:"status"`
	Paper     bool      `json:"paper"`
	Timestamp time.Time `json:"timestamp"`
}

// Fill is an execution fill.
type Fill struct {
	OrderID   string    `json:"order_id"`
	Symbol    string    `json:"symbol"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}
