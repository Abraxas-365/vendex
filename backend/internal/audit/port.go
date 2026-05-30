package audit

import (
	"context"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines persistence operations for the audit domain.
type Repository interface {
	// Create persists a new audit entry and returns the saved entity.
	Create(ctx context.Context, entry AuditEntry) (AuditEntry, error)

	// GetByID retrieves a single audit entry scoped to the tenant.
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.AuditEntryID) (AuditEntry, error)

	// List returns a paginated, filtered list of audit entries for the tenant.
	List(ctx context.Context, tenantID kernel.TenantID, filter AuditFilter, p kernel.PaginationOptions) (kernel.Paginated[AuditEntry], error)

	// CountByAction returns the count of audit entries grouped by action for the given period.
	CountByAction(ctx context.Context, tenantID kernel.TenantID, from, to time.Time) ([]ActionStats, error)
}
