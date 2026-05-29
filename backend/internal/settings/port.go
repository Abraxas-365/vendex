package settings

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines persistence operations for store settings.
type Repository interface {
	// Get retrieves settings for the given tenant.
	// Returns ErrNotFound if no settings have been saved yet.
	Get(ctx context.Context, tenantID kernel.TenantID) (*StoreSettings, error)

	// Upsert creates or replaces the settings row for the given tenant.
	Upsert(ctx context.Context, s *StoreSettings) error
}
