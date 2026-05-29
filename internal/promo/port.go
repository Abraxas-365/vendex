package promo

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// PromoRepository defines persistence operations for Promo entities.
// All operations are scoped by TenantID.
type PromoRepository interface {
	// Create persists a new promo.
	Create(ctx context.Context, p *Promo) error
	// GetByID retrieves a promo by its primary key.
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.PromoID) (*Promo, error)
	// GetByCode retrieves a promo by its discount code (case-insensitive lookup).
	GetByCode(ctx context.Context, tenantID kernel.TenantID, code string) (*Promo, error)
	// Update persists mutations (e.g. deactivation, used_count increment).
	Update(ctx context.Context, p *Promo) error
	// IncrementUsedCount atomically increments the used_count for a promo.
	IncrementUsedCount(ctx context.Context, tenantID kernel.TenantID, id kernel.PromoID) error
	// List returns all promos for a tenant with pagination.
	List(ctx context.Context, tenantID kernel.TenantID, p kernel.Pagination) (kernel.PaginatedResult[Promo], error)
}
