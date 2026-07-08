package transport

import (
	"context"
	"net/http"
	"sync"

	cd "github.com/Soonogo/ai-trading-os/services/common/domain"
	"github.com/Soonogo/ai-trading-os/services/common/eventbus"
	"github.com/Soonogo/ai-trading-os/services/model-router/internal/config"
	mrd "github.com/Soonogo/ai-trading-os/services/model-router/internal/domain"
	"github.com/Soonogo/ai-trading-os/services/model-router/internal/usecase"
	"github.com/gin-gonic/gin"
)

// Server exposes model router HTTP endpoints.
type Server struct {
	cfg      *config.Config
	router   *gin.Engine
	srv      *http.Server
	ensemble *usecase.Ensemble
	bus      eventbus.Publisher
	mu       sync.Mutex
}

// NewServer creates the model router HTTP server.
func NewServer(cfg *config.Config, ensemble *usecase.Ensemble, bus eventbus.Publisher) *Server {
	r := gin.New()
	r.Use(gin.Recovery())

	s := &Server{
		cfg:      cfg,
		router:   r,
		ensemble: ensemble,
		bus:      bus,
	}

	r.GET("/health", s.handleHealth)
	r.POST("/v1/route", s.handleRoute)

	return s
}

// Run starts the HTTP server.
func (s *Server) Run(ctx context.Context) error {
	s.mu.Lock()
	s.srv = &http.Server{
		Addr:    s.cfg.HTTPAddr,
		Handler: s.router,
	}
	s.mu.Unlock()
	return s.srv.ListenAndServe()
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	srv := s.srv
	s.mu.Unlock()
	if srv == nil {
		return nil
	}
	return srv.Shutdown(ctx)
}

func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

type routeRequest struct {
	TraceID string `json:"trace_id" binding:"required"`
	Prompt  string `json:"prompt" binding:"required"`
	System  string `json:"system"`
	Format  string `json:"format"`
}

func (s *Server) handleRoute(c *gin.Context) {
	var req routeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prompt := mrd.Prompt{
		TraceID: req.TraceID,
		Content: req.Prompt,
		System:  req.System,
		Format:  req.Format,
	}

	result, err := s.ensemble.Route(c.Request.Context(), prompt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := cd.NewEvent(req.TraceID, cd.EventTypePortfolioDecision, "model-router", result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	event.Reasoning = &cd.Reasoning{
		Confidence: result.Confidence,
		Evidence:   []string{result.Reasoning},
		ModelCalls: []cd.ModelCall{
			{
				Provider: string(result.Selected.Provider),
				Model:    result.Selected.Model,
				Latency:  int64(result.Selected.Latency / 1e6),
				Output:   result.Selected.Output,
			},
		},
	}

	if err := s.bus.Publish(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusAccepted, gin.H{
			"result":        result,
			"published":     false,
			"publish_error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
