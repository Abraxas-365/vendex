package agent

import (
	"fmt"
	"log"
	"time"
)

// EventKind identifies the type of an agent streaming event.
type EventKind string

const (
	// EventTextDelta is emitted when the LLM streams a chunk of text.
	EventTextDelta EventKind = "text_delta"
	// EventToolStart is emitted when the agent begins executing a tool.
	EventToolStart EventKind = "tool_start"
	// EventToolEnd is emitted when a tool execution completes.
	EventToolEnd EventKind = "tool_end"
	// EventTurnEnd is emitted when the agent's full turn is complete.
	EventTurnEnd EventKind = "turn_end"
	// EventError is emitted when the agent encounters an unrecoverable error.
	EventError EventKind = "error"
)

// Event represents a single streaming event from the agent.
type Event struct {
	Kind      EventKind `json:"kind"`
	Text      string    `json:"text,omitempty"`
	ToolName  string    `json:"tool_name,omitempty"`
	ToolInput string    `json:"tool_input,omitempty"`
	Result    string    `json:"result,omitempty"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// EventHandler collects and dispatches agent streaming events.
// It is safe for concurrent use from a single agent turn.
//
// In the current implementation events are logged to stderr.
// The WebSocket broadcasting layer will be added when the admin UI
// WebSocket endpoint is wired up in cmd/main.go.
type EventHandler struct {
	// SessionID identifies the admin UI session receiving these events.
	// Used for routing when multiple admin sessions are active.
	SessionID string

	// OnEvent is an optional callback invoked for every event.
	// Set it to forward events over a WebSocket connection.
	// If nil, events are only logged.
	OnEvent func(e Event)
}

// NewEventHandler creates an EventHandler for the given admin session.
func NewEventHandler(sessionID string) *EventHandler {
	return &EventHandler{SessionID: sessionID}
}

// emit dispatches an event to the optional callback and logs it.
func (h *EventHandler) emit(e Event) {
	e.Timestamp = time.Now().UTC()
	log.Printf("[agent][%s] event=%s tool=%q text=%q err=%q",
		h.SessionID, e.Kind, e.ToolName, truncate(e.Text, 80), e.Error)
	if h.OnEvent != nil {
		h.OnEvent(e)
	}
}

// OnTextDelta is called by the harness streaming loop for each text chunk.
func (h *EventHandler) OnTextDelta(delta string) {
	h.emit(Event{Kind: EventTextDelta, Text: delta})
}

// OnToolStart is called when the agent begins a tool call.
func (h *EventHandler) OnToolStart(toolName, inputJSON string) {
	h.emit(Event{Kind: EventToolStart, ToolName: toolName, ToolInput: inputJSON})
}

// OnToolEnd is called when a tool call completes.
func (h *EventHandler) OnToolEnd(toolName, result string) {
	h.emit(Event{Kind: EventToolEnd, ToolName: toolName, Result: result})
}

// OnTurnEnd is called when the agent's full response turn finishes.
func (h *EventHandler) OnTurnEnd() {
	h.emit(Event{Kind: EventTurnEnd})
}

// OnError is called when the agent encounters an error it cannot recover from.
func (h *EventHandler) OnError(err error) {
	h.emit(Event{Kind: EventError, Error: fmt.Sprintf("%v", err)})
}

// truncate shortens s to at most n runes for log output.
func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "…"
}
