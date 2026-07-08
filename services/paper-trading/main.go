package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	cd "github.com/Soonogo/ai-trading-os/services/common/domain"
	"github.com/Soonogo/ai-trading-os/services/common/eventbus"
	"github.com/Soonogo/ai-trading-os/services/paper-trading/internal/config"
	"github.com/Soonogo/ai-trading-os/services/paper-trading/internal/ports"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	bus, err := eventbus.NewNATSBus(cfg.NATSURL, "paper-trading")
	if err != nil {
		log.Fatalf("event bus: %v", err)
	}
	defer bus.Close()

	handler := ports.NewHandler(bus)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := bus.Subscribe(ctx, []cd.EventType{cd.EventTypePaperOrder, cd.EventTypeMarketData}, handler.HandleEvent); err != nil {
		log.Fatalf("subscribe: %v", err)
	}

	<-ctx.Done()
	log.Println("paper-trading stopped")
}
