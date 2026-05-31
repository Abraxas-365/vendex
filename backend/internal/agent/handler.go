package agent

import (
	"fmt"
	"log"
	"time"

	"github.com/Abraxas-365/harness/api"
	"github.com/Abraxas-365/harness/query"
	harnesstools "github.com/Abraxas-365/harness/tools"
)

// Compile-time check: EventHandler must implement the harness query.EventHandler interface.
var _ query.EventHandler = (*EventHandler)(nil)

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
// It implements the harness query.EventHandler interface so it can be plugged
// directly into a harness Session.
//
// All e-commerce domain tools are auto-approved (no user confirmation required)
// since the agent runs server-side in a trusted context.
type EventHandler struct {
	// SessionID identifies the admin UI session receiving these events.
	// Used for routing when multiple admin sessions are active.
	SessionID string

	// OnEvent is an optional callback invoked for every event.
	// Set it to forward events over a WebSocket or SSE connection.
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
func (h *EventHandler) OnTextDelta(text string) {
	h.emit(Event{Kind: EventTextDelta, Text: text})
}

// OnThinkingDelta is called when the model emits a thinking/reasoning chunk.
// We surface it as a text_delta so the frontend can render it without extra handling.
func (h *EventHandler) OnThinkingDelta(text string) {
	// Thinking deltas are internal reasoning — log only, do not emit to frontend.
	log.Printf("[agent][%s] thinking: %s", h.SessionID, truncate(text, 120))
}

// OnToolUseStart is called when the agent begins a tool call.
func (h *EventHandler) OnToolUseStart(tu harnesstools.ToolUse) {
	h.emit(Event{
		Kind:      EventToolStart,
		ToolName:  tu.Name,
		ToolInput: string(tu.Input),
	})
}

// OnToolUseEnd is called when a tool call completes.
func (h *EventHandler) OnToolUseEnd(tu harnesstools.ToolUse, result *harnesstools.Result) {
	content := ""
	if result != nil {
		content = result.Content
	}
	h.emit(Event{
		Kind:     EventToolEnd,
		ToolName: tu.Name,
		Result:   content,
	})
}

// OnTurnComplete is called when the agent's full response turn finishes.
// The usage parameter carries token counts for the turn.
func (h *EventHandler) OnTurnComplete(usage api.Usage) {
	log.Printf("[agent][%s] turn_complete input_tokens=%d output_tokens=%d cache_read=%d cache_create=%d",
		h.SessionID, usage.InputTokens, usage.OutputTokens, usage.CacheRead, usage.CacheCreate)
	h.emit(Event{Kind: EventTurnEnd})
}

// OnError is called when the agent encounters an error it cannot recover from.
func (h *EventHandler) OnError(err error) {
	h.emit(Event{Kind: EventError, Error: fmt.Sprintf("%v", err)})
}

// OnRetry is called when the engine silently retries a request (e.g. after
// hitting max_tokens mid-stream). We log and take no other action.
func (h *EventHandler) OnRetry(toolUses []harnesstools.ToolUse) {
	names := make([]string, len(toolUses))
	for i, tu := range toolUses {
		names[i] = tu.Name
	}
	log.Printf("[agent][%s] retry for tool_uses=%v", h.SessionID, names)
}

// OnToolApprovalNeeded is called when a tool requires user approval.
// We auto-approve all tools since the e-commerce agent runs server-side
// in a trusted, headless context.
func (h *EventHandler) OnToolApprovalNeeded(tu harnesstools.ToolUse) bool {
	log.Printf("[agent][%s] auto-approving tool=%q", h.SessionID, tu.Name)
	return true
}

// OnCostConfirmNeeded is called when the session cost exceeds the threshold.
// We auto-confirm to keep server-side sessions running unattended.
func (h *EventHandler) OnCostConfirmNeeded(currentCost, threshold float64) bool {
	log.Printf("[agent][%s] auto-confirming cost=%.4f threshold=%.4f", h.SessionID, currentCost, threshold)
	return true
}

// OnBgTaskComplete is called when a background task reaches a terminal state.
// Headless server handlers only need to log this.
func (h *EventHandler) OnBgTaskComplete(taskID, output string, exitCode int, errStr string, isSubAgent bool, agentName string) {
	log.Printf("[agent][%s] bg_task_complete task_id=%q exit_code=%d is_sub_agent=%v agent=%q err=%q",
		h.SessionID, taskID, exitCode, isSubAgent, agentName, truncate(errStr, 80))
}

// truncate shortens s to at most n runes for log output.
func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "…"
}
