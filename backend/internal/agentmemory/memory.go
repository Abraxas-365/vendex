// Package agentmemory defines the agent memory domain for persistent,
// searchable knowledge entries scoped per tenant.
package agentmemory

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Memory represents a persistent knowledge entry for a tenant.
type Memory struct {
	ID        kernel.AgentMemoryID `json:"id"`
	TenantID  kernel.TenantID      `json:"tenant_id"`
	Category  string               `json:"category"` // e.g. "brand", "product", "seo", "general", "decision"
	Title     string               `json:"title"`    // short summary
	Content   string               `json:"content"`  // full knowledge text
	Tags      []string             `json:"tags"`     // searchable tags
	Source    string               `json:"source"`   // "agent" (auto-created) or "human" (manually added)
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

// CreateMemoryRequest holds input for creating a new memory entry.
type CreateMemoryRequest struct {
	Category string   `json:"category"`
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags"`
	Source   string   `json:"source"` // defaults to "human"
}

// UpdateMemoryRequest holds input for updating an existing memory entry.
type UpdateMemoryRequest struct {
	Category *string   `json:"category,omitempty"`
	Title    *string   `json:"title,omitempty"`
	Content  *string   `json:"content,omitempty"`
	Tags     *[]string `json:"tags,omitempty"`
}

// MemorySearchOptions holds search and filter options for querying memories.
type MemorySearchOptions struct {
	Query    string   // text search against title and content
	Category string   // filter by category (empty = all)
	Tags     []string // filter by tags (AND match, empty = ignore)
}
