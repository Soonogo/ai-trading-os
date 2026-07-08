package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EventType is a namespaced event type constant.
type EventType string

const (
	EventTypeMarketData       EventType = "market.data"
	EventTypeResearchSignal   EventType = "research.signal"
	EventTypeMacroSignal      EventType = "macro.signal"
	EventTypeNewsSignal       EventType = "news.signal"
	EventTypeTechnicalSignal  EventType = "technical.signal"
	EventTypeFundamentalSignal EventType = "fundamental.signal"
	EventTypeRiskSignal       EventType = "risk.signal"
	EventTypePortfolioDecision EventType = "portfolio.decision"
	EventTypeOrderSubmitted   EventType = "execution.order.submitted"
	EventTypeOrderFilled      EventType = "execution.order.filled"
	EventTypeReflectionReport EventType = "reflection.report"
	EventTypeDailyReport      EventType = "report.daily"
)

// AllowedTypes is the set of valid event type constants.
var AllowedTypes = map[EventType]bool{
	EventTypeMarketData:          true,
	EventTypeResearchSignal:    true,
	EventTypeMacroSignal:       true,
	EventTypeNewsSignal:        true,
	EventTypeTechnicalSignal:   true,
	EventTypeFundamentalSignal: true,
	EventTypeRiskSignal:        true,
	EventTypePortfolioDecision: true,
	EventTypeOrderSubmitted:    true,
	EventTypeOrderFilled:       true,
	EventTypeReflectionReport:  true,
	EventTypeDailyReport:       true,
}

// Event is the canonical envelope for all cross-service messages.
type Event struct {
	ID        string          `json:"event_id"`
	TraceID   string          `json:"trace_id"`
	Type      EventType       `json:"type"`
	Version   int             `json:"version"`
	Source    string          `json:"source_agent"`
	Timestamp time.Time       `json:"timestamp"`
	Reasoning *Reasoning      `json:"reasoning,omitempty"`
	Payload   json.RawMessage `json:"payload"`
}

// Reasoning records why a decision was made.
type Reasoning struct {
	Confidence   float64  `json:"confidence"`
	Evidence     []string `json:"evidence"`
	ModelCalls   []ModelCall `json:"model_calls,omitempty"`
	Notes        string   `json:"notes,omitempty"`
}

// ModelCall captures one LLM invocation.
type ModelCall struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Latency  int64  `json:"latency_ms"`
	Output   string `json:"output"`
}

// NewEvent creates a validated event envelope.
func NewEvent(traceID string, typ EventType, source string, payload any) (Event, error) {
	if !AllowedTypes[typ] {
		return Event{}, fmt.Errorf("unknown event type: %s", typ)
	}
	if traceID == "" {
		return Event{}, errors.New("trace_id is required")
	}
	if source == "" {
		return Event{}, errors.New("source_agent is required")
	}
	if payload == nil {
		return Event{}, errors.New("payload is required")
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return Event{}, fmt.Errorf("marshal payload: %w", err)
	}
	return Event{
		ID:        uuid.MustParse(uuid.NewString()).String(),
		TraceID:   traceID,
		Type:      typ,
		Version:   1,
		Source:    source,
		Timestamp: time.Now().UTC(),
		Payload:   raw,
	}, nil
}

// DecodePayload unmarshals the event payload into a concrete type.
func (e *Event) DecodePayload(out any) error {
	if e.Payload == nil {
		return errors.New("event has no payload")
	}
	return json.Unmarshal(e.Payload, out)
}
