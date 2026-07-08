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
	"github.com/Soonogo/ai-trading-os/services/execution-agent/internal/application"
	"github.com/Soonogo/ai-trading-os/services/execution-agent/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	bus, err := eventbus.NewNATSBus(cfg.NATSURL, "execution-agent")
	if err != nil {
		log.Fatalf("event bus: %v", err)
	}
	defer bus.Close()

	exec := application.NewExecutor(cfg.PaperTrading)

	handler := func(ctx context.Context, event cd.Event) error {
		if event.Type != cd.EventTypePortfolioDecision {
			return nil
		}
		decision, err := application.ParseDecision(event.Payload)
		if err != nil {
			return err
		}
		order, fill, err := exec.Execute(decision)
		if err != nil {
			log.Printf("execute skipped: %v", err)
			return nil
		}

		orderPayload, _ := json.Marshal(order)
		orderEvent := cd.Event{
			TraceID:    event.TraceID,
			Type:       cd.EventTypeOrderSubmitted,
			Version:    1,
			Source:     "execution-agent",
			Timestamp:  time.Now().UTC(),
			Reasoning:  event.Reasoning,
			Payload:    orderPayload,
		}
		_ = bus.Publish(ctx, orderEvent)

		fillPayload, _ := json.Marshal(fill)
		fillEvent := cd.Event{
			TraceID:    event.TraceID,
			Type:       cd.EventTypeOrderFilled,
			Version:    1,
			Source:     "execution-agent",
			Timestamp:  time.Now().UTC(),
			Reasoning:  event.Reasoning,
			Payload:    fillPayload,
		}
		return bus.Publish(ctx, fillEvent)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := bus.Subscribe(ctx, []cd.EventType{cd.EventTypePortfolioDecision}, handler); err != nil {
		log.Fatalf("subscribe: %v", err)
	}

	<-ctx.Done()
	log.Println("execution-agent stopped")
}
