package agent

import (
	"context"
	"encoding/json"

	harnesstools "github.com/Abraxas-365/harness/tools"
)

// HarnessTool wraps an agent.Tool to implement the harness tools.Tool interface.
// This bridge converts the vendex Tool interface (map[string]any InputSchema,
// string Execute result) into the harness-compatible interface (json.RawMessage
// InputSchema, *tools.Result Execute result).
type HarnessTool struct {
	inner Tool
}

// AdaptTool wraps an agent.Tool into a harness-compatible tools.Tool.
func AdaptTool(t Tool) harnesstools.Tool {
	return &HarnessTool{inner: t}
}

// AdaptTools wraps a slice of agent.Tool into harness tools.
func AdaptTools(tt []Tool) []harnesstools.Tool {
	out := make([]harnesstools.Tool, len(tt))
	for i, t := range tt {
		out[i] = AdaptTool(t)
	}
	return out
}

// Name delegates to the inner tool.
func (h *HarnessTool) Name() string {
	return h.inner.Name()
}

// Description delegates to the inner tool.
func (h *HarnessTool) Description() string {
	return h.inner.Description()
}

// InputSchema marshals the inner tool's map[string]any schema to json.RawMessage.
// Falls back to an empty object schema if marshalling fails.
func (h *HarnessTool) InputSchema() json.RawMessage {
	raw, err := json.Marshal(h.inner.InputSchema())
	if err != nil {
		return json.RawMessage(`{"type":"object","properties":{}}`)
	}
	return raw
}

// Execute calls the inner tool's Execute method and wraps the string result in
// a *tools.Result. On error, the error message is returned as a tool error
// result rather than propagating as a Go error — this keeps the agent loop
// running instead of aborting the turn.
func (h *HarnessTool) Execute(ctx context.Context, input json.RawMessage) (*harnesstools.Result, error) {
	result, err := h.inner.Execute(ctx, input)
	if err != nil {
		return &harnesstools.Result{
			Content: err.Error(),
			IsError: true,
		}, nil
	}
	return &harnesstools.Result{
		Content: result,
	}, nil
}

// IsReadOnly returns false. Domain tools may create, update, or delete records.
func (h *HarnessTool) IsReadOnly() bool {
	return false
}

// RequiresApproval returns false. All e-commerce domain tools are auto-approved
// since the agent runs server-side in a trusted context.
func (h *HarnessTool) RequiresApproval(_ json.RawMessage) bool {
	return false
}
