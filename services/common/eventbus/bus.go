package eventbus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Soonogo/ai-trading-os/services/common/domain"
)

// Publisher is implemented by message-bus backends.
type Publisher interface {
	Publish(ctx context.Context, event domain.Event) error
}

// Subscriber receives events matching a set of types.
type Subscriber interface {
	Subscribe(ctx context.Context, types []domain.EventType, handler Handler) error
}

// Bus combines publish and subscribe.
type Bus interface {
	Publisher
	Subscriber
	Close() error
}

// Handler processes a single event.
type Handler func(ctx context.Context, event domain.Event) error

// JSON returns an indented JSON string for logs and debugging.
func JSON(e domain.Event) string {
	b, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Sprintf("<unmarshalable event: %v>", err)
	}
	return string(b)
}
