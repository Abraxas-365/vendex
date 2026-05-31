package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Abraxas-365/hada-commerce/internal/agentmemory"
	"github.com/Abraxas-365/hada-commerce/internal/agentmemory/agentmemorysrv"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// ──────────────────────────────────────────────────────────────────────────────
// Memory tools — let the agent read/write persistent knowledge
// ──────────────────────────────────────────────────────────────────────────────

type searchMemoryTool struct {
	tenantID kernel.TenantID
	svc      *agentmemorysrv.Service
}

func (t *searchMemoryTool) Name() string        { return "search_memory" }
func (t *searchMemoryTool) Description() string { return "Search the store's knowledge base for relevant context (brand guidelines, product info, past decisions)" }
func (t *searchMemoryTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"query":    map[string]any{"type": "string", "description": "Search query"},
			"category": map[string]any{"type": "string", "description": "Filter by category: brand, product, seo, general, decision", "enum": []string{"brand", "product", "seo", "general", "decision"}},
		},
		"required": []string{"query"},
	}
}

func (t *searchMemoryTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var req struct {
		Query    string `json:"query"`
		Category string `json:"category"`
	}
	if err := json.Unmarshal(input, &req); err != nil {
		return "invalid input: " + err.Error(), nil
	}
	opts := agentmemory.MemorySearchOptions{
		Query:    req.Query,
		Category: req.Category,
	}
	result, err := t.svc.Search(ctx, t.tenantID, opts, kernel.NewPaginationOptions(1, 10))
	if err != nil {
		return "search failed: " + err.Error(), nil
	}
	if len(result.Items) == 0 {
		return "No memories found matching your query.", nil
	}
	out, _ := json.Marshal(result.Items)
	return string(out), nil
}

type saveMemoryTool struct {
	tenantID kernel.TenantID
	svc      *agentmemorysrv.Service
}

func (t *saveMemoryTool) Name() string        { return "save_memory" }
func (t *saveMemoryTool) Description() string { return "Save a new knowledge entry to the store's memory (for brand guidelines, decisions, product insights)" }
func (t *saveMemoryTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"title":    map[string]any{"type": "string", "description": "Short title for the memory"},
			"content":  map[string]any{"type": "string", "description": "Full knowledge content"},
			"category": map[string]any{"type": "string", "description": "Category", "enum": []string{"brand", "product", "seo", "general", "decision"}},
			"tags":     map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Searchable tags"},
		},
		"required": []string{"title", "content", "category"},
	}
}

func (t *saveMemoryTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var req struct {
		Title    string   `json:"title"`
		Content  string   `json:"content"`
		Category string   `json:"category"`
		Tags     []string `json:"tags"`
	}
	if err := json.Unmarshal(input, &req); err != nil {
		return "invalid input: " + err.Error(), nil
	}
	mem, err := t.svc.Create(ctx, t.tenantID, agentmemory.CreateMemoryRequest{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Tags:     req.Tags,
		Source:   "agent",
	})
	if err != nil {
		return "failed to save memory: " + err.Error(), nil
	}
	return fmt.Sprintf("Memory saved (id: %s): %s", mem.ID, mem.Title), nil
}

type getMemoryContextTool struct {
	tenantID kernel.TenantID
	svc      *agentmemorysrv.Service
}

func (t *getMemoryContextTool) Name() string        { return "get_memory_context" }
func (t *getMemoryContextTool) Description() string { return "Get relevant memories formatted as context for a specific topic (useful before making decisions)" }
func (t *getMemoryContextTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"topic": map[string]any{"type": "string", "description": "Topic to retrieve context for"},
		},
		"required": []string{"topic"},
	}
}

func (t *getMemoryContextTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var req struct {
		Topic string `json:"topic"`
	}
	if err := json.Unmarshal(input, &req); err != nil {
		return "invalid input: " + err.Error(), nil
	}
	result, err := t.svc.GetContext(ctx, t.tenantID, req.Topic)
	if err != nil {
		return "failed to get context: " + err.Error(), nil
	}
	if result == "" {
		return "No relevant memories found for this topic.", nil
	}
	return result, nil
}
