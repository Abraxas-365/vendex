package bulkops

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Repository defines the persistence contract for the bulk operations domain.
type Repository interface {
	// Create persists a new bulk operation and its items atomically.
	Create(ctx context.Context, op *BulkOperation, items []BulkOperationItem) error

	// GetByID returns a bulk operation scoped to the tenant.
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.BulkOperationID) (*BulkOperation, error)

	// List returns paginated bulk operations for a tenant.
	List(ctx context.Context, tenantID kernel.TenantID, page, pageSize int) (kernel.Paginated[BulkOperation], error)

	// ListItems returns paginated items for a specific bulk operation.
	ListItems(ctx context.Context, tenantID kernel.TenantID, operationID kernel.BulkOperationID, page, pageSize int) (kernel.Paginated[BulkOperationItem], error)

	// UpdateStatus updates the status (and optional started_at / completed_at) of a bulk operation.
	UpdateStatus(ctx context.Context, tenantID kernel.TenantID, id kernel.BulkOperationID, status OperationStatus) error

	// UpdateOperation persists the full mutable state of a bulk operation (counters, errors, timestamps).
	UpdateOperation(ctx context.Context, op *BulkOperation) error

	// UpdateItem persists the result of a single bulk operation item.
	UpdateItem(ctx context.Context, item *BulkOperationItem) error
}
