package theme

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// ThemeRepository defines the persistence interface for themes.
type ThemeRepository interface {
	Create(ctx context.Context, t *Theme) error
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ThemeID) (*Theme, error)
	GetActive(ctx context.Context, tenantID kernel.TenantID) (*Theme, error)
	List(ctx context.Context, tenantID kernel.TenantID) ([]Theme, error)
	Update(ctx context.Context, t *Theme) error
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.ThemeID) error
	DeactivateAll(ctx context.Context, tenantID kernel.TenantID) error
}
