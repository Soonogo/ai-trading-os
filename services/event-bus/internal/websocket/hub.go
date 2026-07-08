package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Soonogo/ai-trading-os/services/common/domain"
	"github.com/Soonogo/ai-trading-os/services/common/eventbus"
	"github.com/gorilla/websocket"
)

// Hub manages websocket clients and fan-out from the event bus.
type Hub struct {
	bus       eventbus.Publisher
	clients   map[*client]bool
	broadcast chan domain.Event
	register  chan *client
	unregister chan *client
	upgrader  websocket.Upgrader
	mu        sync.RWMutex
}

// NewHub creates a websocket hub.
func NewHub(bus eventbus.Publisher) *Hub {
	return &Hub{
		bus:        bus,
		clients:    make(map[*client]bool),
		broadcast:  make(chan domain.Event, 256),
		register:   make(chan *client),
		unregister: make(chan *client),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

// Start runs the hub event loop.
func (h *Hub) Start(ctx context.Context) error {
	go h.run(ctx)
	return nil
}

func (h *Hub) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			h.closeAll()
			return
		case c := <-h.register:
			h.clients[c] = true
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
		case event := <-h.broadcast:
			h.mu.RLock()
			for c := range h.clients {
				select {
				case c.send <- event:
				default:
					close(c.send)
					delete(h.clients, c)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) closeAll() {
	for c := range h.clients {
		delete(h.clients, c)
		close(c.send)
	}
}

// Upgrade upgrades an HTTP connection to websocket.
func (h *Hub) Upgrade(w http.ResponseWriter, r *http.Request) error {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("upgrade: %w", err)
	}
	c := &client{hub: h, conn: conn, send: make(chan domain.Event, 128)}
	h.register <- c
	go c.writePump()
	go c.readPump()
	return nil
}

// Broadcast is a Handler-compatible fan-out function.
func (h *Hub) Broadcast(ctx context.Context, event domain.Event) error {
	select {
	case h.broadcast <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan domain.Event
}

func (c *client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(8192)
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("websocket error: %v\n", err)
			}
			break
		}
	}
}

func (c *client) writePump() {
	defer func() {
		_ = c.conn.Close()
	}()
	for event := range c.send {
		body, err := json.Marshal(event)
		if err != nil {
			continue
		}
		_ = c.conn.WriteMessage(websocket.TextMessage, body)
	}
}
