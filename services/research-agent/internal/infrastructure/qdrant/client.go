package qdrant

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Client is a minimal Qdrant HTTP client.
type Client struct {
	baseURL string
	http    *http.Client
}

// NewClient creates a Qdrant client.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

// Search runs a dummy keyword search. In production this uses vector similarity.
func (c *Client) Search(ctx context.Context, query string, limit int) ([]string, error) {
	// Placeholder: return static context when Qdrant is empty.
	return []string{
		"SPY daily ATR expanding, relative volume +22%.",
		"AAPL broke above 20-day VWAP with bullish MACD histogram.",
		"TSLA options flow shows call skew at 110% of put volume.",
	}, nil
}

// Store persists a document. In production this upserts vectors to a collection.
func (c *Client) Store(ctx context.Context, document string, metadata map[string]any) error {
	payload, _ := json.Marshal(map[string]any{
		"document": document,
		"metadata": metadata,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/collections/default/points", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = c.http.Do(req)
	return err
}
