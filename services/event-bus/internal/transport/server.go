package transport

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/Soonogo/ai-trading-os/services/common/domain"
	"github.com/Soonogo/ai-trading-os/services/common/eventbus"
	"github.com/Soonogo/ai-trading-os/services/event-bus/internal/config"
	"github.com/Soonogo/ai-trading-os/services/event-bus/internal/websocket"
	"github.com/gin-gonic/gin"
)

// Server is the HTTP/WS gateway.
type Server struct {
	cfg    *config.Config
	bus    eventbus.Bus
	hub    *websocket.Hub
	router *gin.Engine
	srv    *http.Server
	mu     sync.Mutex
}

// NewServer creates a Gin server with event routes.
func NewServer(cfg *config.Config, bus eventbus.Bus, hub *websocket.Hub) *Server {
	r := gin.New()
	r.Use(gin.Recovery())

	s := &Server{
		cfg:    cfg,
		bus:    bus,
		hub:    hub,
		router: r,
	}

	r.GET("/health", s.handleHealth)
	r.POST("/v1/events", s.handlePublish)
	r.POST("/v1/prompts", s.handlePrompt)
	r.GET("/v1/ws", s.handleWS)

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

type publishRequest struct {
	TraceID string          `json:"trace_id" binding:"required"`
	Type    string          `json:"type" binding:"required"`
	Source  string          `json:"source_agent" binding:"required"`
	Payload interface{}     `json:"payload" binding:"required"`
	Reason  *domain.Reasoning `json:"reasoning,omitempty"`
}

func (s *Server) handlePublish(c *gin.Context) {
	var req publishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event, err := domain.NewEvent(req.TraceID, domain.EventType(req.Type), req.Source, req.Payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event.Reasoning = req.Reason
	if err := s.bus.Publish(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, event)
}

type promptRequest struct {
	TraceID string `json:"trace_id" binding:"required"`
	Prompt  string `json:"prompt" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
}

func (s *Server) handlePrompt(c *gin.Context) {
	var req promptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event, err := domain.NewEvent(req.TraceID, domain.EventTypePortfolioDecision, "event-bus", promptRequest{
		TraceID: req.TraceID,
		Prompt:  req.Prompt,
		UserID:  req.UserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	event.Reasoning = &domain.Reasoning{
		Confidence: 1.0,
		Evidence:   []string{"user prompt accepted"},
		Notes:      req.Prompt,
	}
	if err := s.bus.Publish(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"trace_id":  req.TraceID,
		"status":    "dispatched",
		"timestamp": time.Now().UTC(),
	})
}

func (s *Server) handleWS(c *gin.Context) {
	if err := s.hub.Upgrade(c.Writer, c.Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
