package agent

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/Abraxas-365/vendex/internal/agentmemory"
	"github.com/Abraxas-365/vendex/internal/agentmemory/agentmemorysrv"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// ── Agent memory stubs ──

type stubAgentMemoryRepo struct{}

func (s *stubAgentMemoryRepo) Create(_ context.Context, m agentmemory.Memory) (agentmemory.Memory, error) {
	m.ID = "mem-1"
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return m, nil
}
func (s *stubAgentMemoryRepo) GetByID(_ context.Context, _ kernel.TenantID, _ kernel.AgentMemoryID) (agentmemory.Memory, error) {
	return agentmemory.Memory{
		ID:       "mem-1",
		TenantID: testTenant,
		Title:    "Store Policy",
		Content:  "Free shipping over $50",
	}, nil
}
func (s *stubAgentMemoryRepo) Update(_ context.Context, m agentmemory.Memory) (agentmemory.Memory, error) {
	return m, nil
}
func (s *stubAgentMemoryRepo) Delete(_ context.Context, _ kernel.TenantID, _ kernel.AgentMemoryID) error {
	return nil
}
func (s *stubAgentMemoryRepo) List(_ context.Context, _ kernel.TenantID, _ kernel.PaginationOptions) (kernel.Paginated[agentmemory.Memory], error) {
	return kernel.Paginated[agentmemory.Memory]{
		Items: []agentmemory.Memory{
			{ID: "mem-1", Title: "Store Policy", Content: "Free shipping over $50"},
		},
		Total: 1, Page: 1, TotalPages: 1,
	}, nil
}
func (s *stubAgentMemoryRepo) Search(_ context.Context, _ kernel.TenantID, _ agentmemory.MemorySearchOptions, _ kernel.PaginationOptions) (kernel.Paginated[agentmemory.Memory], error) {
	return kernel.Paginated[agentmemory.Memory]{
		Items: []agentmemory.Memory{
			{ID: "mem-1", Title: "Store Policy", Content: "Free shipping over $50"},
		},
		Total: 1, Page: 1, TotalPages: 1,
	}, nil
}


// ── Tests ──

func TestSearchMemoryTool_Execute(t *testing.T) {
	svc := agentmemorysrv.NewService(&stubAgentMemoryRepo{})
	tool := &searchMemoryTool{svc: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"query":"shipping","category":"policy"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Store Policy") {
		t.Errorf("expected 'Store Policy' in result, got: %s", result)
	}
}

func TestSearchMemoryTool_Execute_InvalidJSON(t *testing.T) {
	svc := agentmemorysrv.NewService(&stubAgentMemoryRepo{})
	tool := &searchMemoryTool{svc: svc, tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{bad`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "invalid input") {
		t.Errorf("expected 'invalid input' in result, got: %s", result)
	}
}

func TestSaveMemoryTool_Execute(t *testing.T) {
	svc := agentmemorysrv.NewService(&stubAgentMemoryRepo{})
	tool := &saveMemoryTool{svc: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"title":"Return Policy","content":"30 day returns","category":"policy","tags":["returns","policy"]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Memory saved") {
		t.Errorf("expected 'Memory saved' in result, got: %s", result)
	}
}

func TestSaveMemoryTool_Execute_InvalidJSON(t *testing.T) {
	svc := agentmemorysrv.NewService(&stubAgentMemoryRepo{})
	tool := &saveMemoryTool{svc: svc, tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{bad`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "invalid input") {
		t.Errorf("expected 'invalid input' in result, got: %s", result)
	}
}

func TestGetMemoryContextTool_Execute(t *testing.T) {
	svc := agentmemorysrv.NewService(&stubAgentMemoryRepo{})
	tool := &getMemoryContextTool{svc: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"topic":"shipping policy"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "shipping") {
		t.Errorf("expected 'shipping' in result, got: %s", result)
	}
}

func TestGetMemoryContextTool_Execute_InvalidJSON(t *testing.T) {
	svc := agentmemorysrv.NewService(&stubAgentMemoryRepo{})
	tool := &getMemoryContextTool{svc: svc, tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{bad`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "invalid input") {
		t.Errorf("expected 'invalid input' in result, got: %s", result)
	}
}
