package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Soonogo/ai-trading-os/services/common/eventbus"
	"github.com/Soonogo/ai-trading-os/services/model-router/internal/config"
	"github.com/Soonogo/ai-trading-os/services/model-router/internal/infrastructure/llm"
	"github.com/Soonogo/ai-trading-os/services/model-router/internal/transport"
	"github.com/Soonogo/ai-trading-os/services/model-router/internal/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	bus, err := eventbus.NewNATSBus(cfg.NATSURL, "model-router")
	if err != nil {
		log.Fatalf("event bus: %v", err)
	}
	defer bus.Close()

	client := llm.NewClientSet(cfg.OpenAIKey, cfg.AnthropicKey, cfg.GroqKey)
	ensemble := usecase.NewEnsemble(client, usecase.DefaultRegistry())

	srv := transport.NewServer(cfg, ensemble, bus)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := srv.Run(ctx); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
	log.Println("model-router stopped")
}
