package agent

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

// mockTool is a test double for the Tool interface.
type mockTool struct {
	name        string
	description string
	schema      map[string]any
	execResult  string
	execErr     error
}

func (m *mockTool) Name() string        { return m.name }
func (m *mockTool) Description() string { return m.description }
func (m *mockTool) InputSchema() map[string]any { return m.schema }
func (m *mockTool) Execute(_ context.Context, _ json.RawMessage) (string, error) {
	return m.execResult, m.execErr
}

// newValidMock returns a mockTool with sensible defaults for happy-path tests.
func newValidMock() *mockTool {
	return &mockTool{
		name:        "test_tool",
		description: "A test tool",
		schema: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
	}
}

// ----- AdaptTool tests -----

func TestAdaptTool_Name(t *testing.T) {
	mock := newValidMock()
	mock.name = "my_tool"
	adapted := AdaptTool(mock)
	if adapted.Name() != "my_tool" {
		t.Errorf("Name() = %q, want %q", adapted.Name(), "my_tool")
	}
}

func TestAdaptTool_Description(t *testing.T) {
	mock := newValidMock()
	mock.description = "does something useful"
	adapted := AdaptTool(mock)
	if adapted.Description() != "does something useful" {
		t.Errorf("Description() = %q, want %q", adapted.Description(), "does something useful")
	}
}

func TestAdaptTool_InputSchema(t *testing.T) {
	mock := newValidMock()
	mock.schema = map[string]any{
		"type": "object",
		"properties": map[string]any{
			"foo": map[string]any{"type": "string"},
		},
	}
	adapted := AdaptTool(mock)

	raw := adapted.InputSchema()
	if len(raw) == 0 {
		t.Fatal("InputSchema() returned empty bytes")
	}

	// Must be valid JSON.
	var parsed map[string]any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatalf("InputSchema() is not valid JSON: %v", err)
	}

	// Must preserve the type field.
	if parsed["type"] != "object" {
		t.Errorf("InputSchema() 'type' = %v, want 'object'", parsed["type"])
	}

	// Must preserve the properties key.
	if _, ok := parsed["properties"]; !ok {
		t.Error("InputSchema() missing 'properties' key")
	}
}

func TestAdaptTool_Execute_Success(t *testing.T) {
	mock := newValidMock()
	mock.execResult = "ok"
	mock.execErr = nil

	adapted := AdaptTool(mock)
	result, err := adapted.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("Execute() returned unexpected Go error: %v", err)
	}
	if result == nil {
		t.Fatal("Execute() returned nil result")
	}
	if result.IsError {
		t.Errorf("Execute() result.IsError = true, want false")
	}
	if result.Content != "ok" {
		t.Errorf("Execute() result.Content = %q, want %q", result.Content, "ok")
	}
}

func TestAdaptTool_Execute_Error(t *testing.T) {
	mock := newValidMock()
	mock.execResult = ""
	mock.execErr = errors.New("fail")

	adapted := AdaptTool(mock)
	result, err := adapted.Execute(context.Background(), json.RawMessage(`{}`))
	// The adapter absorbs the inner error and surfaces it as a tool-level error,
	// NOT a Go error, so the agent loop keeps running.
	if err != nil {
		t.Fatalf("Execute() returned unexpected Go error: %v (adapter should absorb inner errors)", err)
	}
	if result == nil {
		t.Fatal("Execute() returned nil result")
	}
	if !result.IsError {
		t.Errorf("Execute() result.IsError = false, want true")
	}
	if result.Content != "fail" {
		t.Errorf("Execute() result.Content = %q, want %q", result.Content, "fail")
	}
}

func TestAdaptTool_IsReadOnly(t *testing.T) {
	adapted := AdaptTool(newValidMock())
	if adapted.IsReadOnly() {
		t.Error("IsReadOnly() = true, want false (domain tools may mutate state)")
	}
}

func TestAdaptTool_RequiresApproval(t *testing.T) {
	adapted := AdaptTool(newValidMock())
	if adapted.RequiresApproval(nil) {
		t.Error("RequiresApproval(nil) = true, want false (agent runs server-side in trusted context)")
	}
}

// ----- AdaptTools (slice) test -----

func TestAdaptTools_Slice(t *testing.T) {
	mocks := []Tool{
		&mockTool{name: "tool_a", description: "a", schema: map[string]any{"type": "object", "properties": map[string]any{}}},
		&mockTool{name: "tool_b", description: "b", schema: map[string]any{"type": "object", "properties": map[string]any{}}},
		&mockTool{name: "tool_c", description: "c", schema: map[string]any{"type": "object", "properties": map[string]any{}}},
	}

	adapted := AdaptTools(mocks)

	if len(adapted) != len(mocks) {
		t.Fatalf("AdaptTools() returned %d tools, want %d", len(adapted), len(mocks))
	}

	for i, a := range adapted {
		want := mocks[i].Name()
		if a.Name() != want {
			t.Errorf("AdaptTools()[%d].Name() = %q, want %q", i, a.Name(), want)
		}
	}
}

// TestAdaptTool_InputSchema_MarshalFallback verifies that even if the schema
// map contains an unmarshalable value, InputSchema() falls back to a valid
// empty object schema rather than returning nil or panicking.
func TestAdaptTool_InputSchema_MarshalFallback(t *testing.T) {
	// json.Marshal cannot marshal a channel, so use one to trigger the fallback.
	mock := &mockTool{
		name:        "bad_schema_tool",
		description: "tool with unmarshalable schema",
		schema: map[string]any{
			"type":       "object",
			"properties": make(chan int), // not JSON-serialisable
		},
	}

	adapted := AdaptTool(mock)
	raw := adapted.InputSchema()

	if len(raw) == 0 {
		t.Fatal("InputSchema() returned empty bytes even on fallback")
	}

	var parsed map[string]any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatalf("InputSchema() fallback is not valid JSON: %v", err)
	}

	// Fallback must still be a valid object schema.
	if parsed["type"] != "object" {
		t.Errorf("InputSchema() fallback 'type' = %v, want 'object'", parsed["type"])
	}
}
