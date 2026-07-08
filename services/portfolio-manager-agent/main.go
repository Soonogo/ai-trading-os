package main

import (
	"context"
	"encoding/json"
	"log"
	"os/signal"
	"syscall"
	"time"

	cd "github.com/Soonogo/ai-trading-os/services/common/domain"
	"github.com/Soonogo/ai-trading-os/services/common/eventbus"
	"github.com/Soonogo/ai-trading-os/services/portfolio-manager-agent/internal/application"
	"github.com/Soonogo/ai-trading-os/services/portfolio-manager-agent/internal/config"
	"github.com/Soonogo/ai-trading-os/services/portfolio-manager-agent/internal/domain"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	bus, err := eventbus.NewNATSBus(cfg.NATSURL, "portfolio-manager-agent")
	if err != nil {
		log.Fatalf("event bus: %v", err)
	}
	defer bus.Close()

	manager := application.NewManager(domain.RiskBudget{
		MaxPositionSize: 100,
		MaxDrawdownPct:  0.05,
	})

	handler := func(ctx context.Context, event cd.Event) error {
		if event.Type != cd.EventTypeResearchSignal && event.Type != cd.EventTypePortfolioDecision {
			return nil
		}
		signal, err := application.ParseSignal(event.Payload)
		if err != nil {
			return err
		}
		// If this is a model-router decision, the payload itself may not be a research signal.
		// For now treat model-router output as empty.
		decision, reasoning, err := manager.Decide(signal, "", event.TraceID)
		if err != nil {
			log.Printf("decide skipped: %v", err)
			return nil
		}
		payload, _ := json.Marshal(decision)
		e := cd.Event{
			TraceID:    event.TraceID,
			Type:       cd.EventTypePortfolioDecision,
			Version:    1,
			Source:     "portfolio-manager-agent",
			Timestamp:  time.Now().UTC(),
			Reasoning:  reasoning,
			Payload:    payload,
		}
		return bus.Publish(ctx, e)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := bus.Subscribe(ctx, []cd.EventType{cd.EventTypeResearchSignal, cd.EventTypePortfolioDecision}, handler); err != nil {
		log.Fatalf("subscribe: %v", err)
	}

	<-ctx.Done()
	log.Println("portfolio-manager-agent stopped")
}
