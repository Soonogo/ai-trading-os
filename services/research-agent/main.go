package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	cd "github.com/Soonogo/ai-trading-os/services/common/domain"
	"github.com/Soonogo/ai-trading-os/services/common/eventbus"
	"github.com/Soonogo/ai-trading-os/services/research-agent/internal/application"
	"github.com/Soonogo/ai-trading-os/services/research-agent/internal/config"
	"github.com/Soonogo/ai-trading-os/services/research-agent/internal/infrastructure/llm"
	"github.com/Soonogo/ai-trading-os/services/research-agent/internal/infrastructure/qdrant"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	bus, err := eventbus.NewNATSBus(cfg.NATSURL, "research-agent")
	if err != nil {
		log.Fatalf("event bus: %v", err)
	}
	defer bus.Close()

	rag := qdrant.NewClient(cfg.QdrantURL)
	llmClient := llm.NewClient(cfg.OpenAIKey, "")
	researcher := application.NewResearcher(rag, llmClient)

	handler := func(ctx context.Context, event cd.Event) error {
		if event.Type != cd.EventTypeMarketData {
			return nil
		}
		var md marketData
		if err := event.DecodePayload(&md); err != nil {
			return err
		}
		signal, reasoning, err := researcher.Research(ctx, fmt.Sprintf("Find the highest probability swing trade setup for %s", md.Symbol))
		if err != nil {
			log.Printf("research error: %v", err)
			return nil
		}

		payload, _ := json.Marshal(signal)
		e := cd.Event{
			TraceID:   event.TraceID,
			Type:      cd.EventTypeResearchSignal,
			Version:   1,
			Source:    "research-agent",
			Timestamp: time.Now().UTC(),
			Reasoning: reasoning,
			Payload:   payload,
		}
		return bus.Publish(ctx, e)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := bus.Subscribe(ctx, []cd.EventType{cd.EventTypeMarketData}, handler); err != nil {
		log.Fatalf("subscribe: %v", err)
	}

	<-ctx.Done()
	log.Println("research-agent stopped")
}

type marketData struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
}
