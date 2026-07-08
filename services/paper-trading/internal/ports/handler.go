package ports

import (
	"context"
	"encoding/json"
	"log"
	"time"

	cd "github.com/Soonogo/ai-trading-os/services/common/domain"
	eb "github.com/Soonogo/ai-trading-os/services/common/eventbus"
	"github.com/Soonogo/ai-trading-os/services/paper-trading/internal/application"
)

// Handler processes paper orders and market data.
type Handler struct {
	engine *application.MatchingEngine
	bus    eb.Bus
}

// NewHandler creates a paper trading handler.
func NewHandler(bus eb.Bus) *Handler {
	return &Handler{
		engine: application.NewMatchingEngine(),
		bus:    bus,
	}
}

// HandleEvent processes a single event.
func (h *Handler) HandleEvent(ctx context.Context, event cd.Event) error {
	switch event.Type {
	case cd.EventTypePaperOrder:
		order, err := application.ParseOrder(event.Payload)
		if err != nil {
			log.Printf("parse paper order: %v", err)
			return nil
		}
		fill, err := h.engine.MatchOrder(order)
		if err != nil {
			log.Printf("match paper order: %v", err)
			return nil
		}
		if err := h.engine.ApplyFill(fill); err != nil {
			log.Printf("apply paper fill: %v", err)
			return nil
		}
		payload, _ := json.Marshal(fill)
		out := cd.Event{
			TraceID:   event.TraceID,
			Type:      cd.EventTypePaperFill,
			Version:   1,
			Source:    "paper-trading",
			Timestamp: time.Now().UTC(),
			Payload:   payload,
		}
		return h.bus.Publish(ctx, out)
	case cd.EventTypeMarketData:
		candle, err := application.ParseCandle(event.Payload)
		if err != nil {
			log.Printf("parse market data: %v", err)
			return nil
		}
		if candle.Symbol == "" {
			log.Printf("market data missing symbol")
			return nil
		}
		h.engine.SetCandle(candle)
	}
	return nil
}
