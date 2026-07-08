package ports

import (
	"context"
	"log"
	"time"

	"github.com/Soonogo/ai-trading-os/services/backtest-engine/internal/application"
	"github.com/Soonogo/ai-trading-os/services/backtest-engine/internal/domain"
	cd "github.com/Soonogo/ai-trading-os/services/common/domain"
	eb "github.com/Soonogo/ai-trading-os/services/common/eventbus"
)

// Handler processes backtest requests.
type Handler struct {
	service *application.Service
	bus     eb.Bus
}

// NewHandler creates a backtest handler.
func NewHandler(bus eb.Bus) *Handler {
	return &Handler{
		service: application.NewService(),
		bus:     bus,
	}
}

// HandleEvent processes a single event.
func (h *Handler) HandleEvent(ctx context.Context, event cd.Event) error {
	if event.Type != cd.EventTypeBacktestRequest {
		return nil
	}
	var req domain.BacktestRequest
	if err := event.DecodePayload(&req); err != nil {
		log.Printf("decode backtest request: %v", err)
		return nil
	}
	if req.TraceID == "" {
		req.TraceID = event.TraceID
	}
	result, err := h.service.Run(ctx, req)
	if err != nil {
		log.Printf("run backtest: %v", err)
		result = map[string]any{"error": err.Error()}
	}
	out, err := h.service.ToEvent(req.TraceID, result)
	if err != nil {
		return err
	}
	out.Timestamp = time.Now().UTC()
	return h.bus.Publish(ctx, out)
}
