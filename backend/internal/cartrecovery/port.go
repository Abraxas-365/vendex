package cartrecovery

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Repository defines the persistence interface for the cart recovery domain.
type Repository interface {
	// Create persists a new recovery email record.
	Create(ctx context.Context, email *RecoveryEmail) error

	// GetByID retrieves a recovery email by its ID, scoped to a tenant.
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.RecoveryID) (*RecoveryEmail, error)

	// GetByCartID returns all recovery emails for a given cart, scoped to a tenant.
	GetByCartID(ctx context.Context, tenantID kernel.TenantID, cartID kernel.CartID) ([]RecoveryEmail, error)

	// ListPending returns all pending recovery emails for a tenant.
	ListPending(ctx context.Context, tenantID kernel.TenantID) ([]RecoveryEmail, error)

	// List returns a paginated list of recovery emails for a tenant.
	List(ctx context.Context, tenantID kernel.TenantID, page, pageSize int) (kernel.Paginated[RecoveryEmail], error)

	// Update persists changes to an existing recovery email record.
	Update(ctx context.Context, email *RecoveryEmail) error

	// GetStats returns aggregate recovery statistics for a tenant.
	GetStats(ctx context.Context, tenantID kernel.TenantID) (*RecoveryStats, error)
}
