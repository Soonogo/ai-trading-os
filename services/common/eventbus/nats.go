package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Soonogo/ai-trading-os/services/common/domain"
	"github.com/nats-io/nats.go"
)

// NATSBus implements Bus using NATS JetStream (legacy API).
type NATSBus struct {
	conn   *nats.Conn
	js     nats.JetStreamContext
	stream string

	subs []*nats.Subscription
	mu   sync.Mutex
}

// NewNATSBus creates a JetStream-backed event bus.
func NewNATSBus(url, stream string) (*NATSBus, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}
	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("jetstream context: %w", err)
	}
	return &NATSBus{conn: nc, js: js, stream: stream}, nil
}

// Publish sends an event to a JetStream subject derived from the event type.
func (b *NATSBus) Publish(ctx context.Context, event domain.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	subject := fmt.Sprintf("events.%s", event.Type)
	_, err = b.js.Publish(subject, body)
	return err
}

// Subscribe registers a queue group handler for one or more event types.
func (b *NATSBus) Subscribe(ctx context.Context, types []domain.EventType, handler Handler) error {
	for _, t := range types {
		subject := fmt.Sprintf("events.%s", t)
		queue := b.stream + "-" + string(t)
		sub, err := b.js.Subscribe(subject, func(msg *nats.Msg) {
			var event domain.Event
			if err := json.Unmarshal(msg.Data, &event); err != nil {
				_ = msg.Nak()
				return
			}
			if err := handler(context.Background(), event); err != nil {
				_ = msg.Nak()
				return
			}
			_ = msg.Ack()
		},
			nats.Durable(queue),
			nats.ManualAck(),
			nats.DeliverAll(),
		)
		if err != nil {
			return fmt.Errorf("subscribe %s: %w", t, err)
		}
		b.mu.Lock()
		b.subs = append(b.subs, sub)
		b.mu.Unlock()
	}
	return nil
}

// Close drains subscriptions and closes the connection.
func (b *NATSBus) Close() error {
	b.mu.Lock()
	for _, sub := range b.subs {
		_ = sub.Unsubscribe()
	}
	b.subs = b.subs[:0]
	b.mu.Unlock()
	b.conn.Close()
	return nil
}
