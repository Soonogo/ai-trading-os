package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Soonogo/ai-trading-os/services/common/domain"
	"github.com/Soonogo/ai-trading-os/services/common/eventbus"
	"github.com/Soonogo/ai-trading-os/services/event-bus/internal/config"
	"github.com/Soonogo/ai-trading-os/services/event-bus/internal/transport"
	"github.com/Soonogo/ai-trading-os/services/event-bus/internal/websocket"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	bus, err := eventbus.NewNATSBus(cfg.NATSURL, cfg.ServiceName)
	if err != nil {
		log.Fatalf("event bus: %v", err)
	}
	defer bus.Close()

	hub := websocket.NewHub(bus)

	srv := transport.NewServer(cfg, bus, hub)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := hub.Start(ctx); err != nil {
		log.Fatalf("hub start: %v", err)
	}

	go func() {
		if err := srv.Run(ctx); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	// Subscribe to all event types and broadcast to websocket clients.
	var types []domain.EventType
	for t := range domain.AllowedTypes {
		types = append(types, t)
	}
	_ = bus.Subscribe(ctx, types, hub.Broadcast)

	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
	log.Println("event-bus stopped")
}
