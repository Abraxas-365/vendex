package agentmemory

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines persistence operations for Memory records.
type Repository interface {
	Create(ctx context.Context, m Memory) (Memory, error)
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentMemoryID) (Memory, error)
	Update(ctx context.Context, m Memory) (Memory, error)
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentMemoryID) error
	List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[Memory], error)
	Search(ctx context.Context, tenantID kernel.TenantID, opts MemorySearchOptions, p kernel.PaginationOptions) (kernel.Paginated[Memory], error)
}
