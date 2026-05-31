// Package agentmemorysrv implements business logic for the agent memory domain.
package agentmemorysrv

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Abraxas-365/vendex/internal/agentmemory"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Service manages agent memory lifecycle.
type Service struct {
	repo agentmemory.Repository
}

// NewService creates a new agentmemory Service.
func NewService(repo agentmemory.Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new memory entry for the given tenant.
func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, req agentmemory.CreateMemoryRequest) (agentmemory.Memory, error) {
	if strings.TrimSpace(req.Title) == "" {
		return agentmemory.Memory{}, agentmemory.ErrInvalidInput
	}
	if strings.TrimSpace(req.Content) == "" {
		return agentmemory.Memory{}, agentmemory.ErrInvalidInput
	}

	source := req.Source
	if source == "" {
		source = "human"
	}

	tags := req.Tags
	if tags == nil {
		tags = []string{}
	}

	category := req.Category
	if category == "" {
		category = "general"
	}

	now := time.Now()
	m := agentmemory.Memory{
		ID:        kernel.AgentMemoryID(uuid.New().String()),
		TenantID:  tenantID,
		Category:  category,
		Title:     req.Title,
		Content:   req.Content,
		Tags:      tags,
		Source:    source,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.repo.Create(ctx, m)
}

// Get retrieves a single memory entry by ID, scoped to a tenant.
func (s *Service) Get(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentMemoryID) (agentmemory.Memory, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// Update applies a partial update to an existing memory entry.
func (s *Service) Update(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentMemoryID, req agentmemory.UpdateMemoryRequest) (agentmemory.Memory, error) {
	m, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return agentmemory.Memory{}, err
	}

	if req.Category != nil {
		m.Category = *req.Category
	}
	if req.Title != nil {
		m.Title = *req.Title
	}
	if req.Content != nil {
		m.Content = *req.Content
	}
	if req.Tags != nil {
		m.Tags = *req.Tags
	}
	m.UpdatedAt = time.Now()

	return s.repo.Update(ctx, m)
}

// Delete removes a memory entry by ID, scoped to a tenant.
func (s *Service) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentMemoryID) error {
	return s.repo.Delete(ctx, tenantID, id)
}

// List returns paginated memories for a tenant.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[agentmemory.Memory], error) {
	return s.repo.List(ctx, tenantID, p)
}

// Search returns memories matching the given options, scoped to a tenant.
func (s *Service) Search(ctx context.Context, tenantID kernel.TenantID, opts agentmemory.MemorySearchOptions, p kernel.PaginationOptions) (kernel.Paginated[agentmemory.Memory], error) {
	return s.repo.Search(ctx, tenantID, opts, p)
}

// GetContext searches memories matching the query and returns the top 10
// formatted as a plain-text context string suitable for injection into
// an agent system prompt.
func (s *Service) GetContext(ctx context.Context, tenantID kernel.TenantID, query string) (string, error) {
	opts := agentmemory.MemorySearchOptions{
		Query: query,
	}
	p := kernel.NewPaginationOptions(1, 10)

	result, err := s.repo.Search(ctx, tenantID, opts, p)
	if err != nil {
		return "", err
	}

	if len(result.Items) == 0 {
		return "", nil
	}

	var sb strings.Builder
	sb.WriteString("=== Store Knowledge Base ===\n\n")
	for i, m := range result.Items {
		fmt.Fprintf(&sb, "[%d] %s (%s)\n%s\n\n", i+1, m.Title, m.Category, m.Content)
	}
	return sb.String(), nil
}
