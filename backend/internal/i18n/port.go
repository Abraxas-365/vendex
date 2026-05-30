package i18n

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines the persistence interface for the i18n domain.
type Repository interface {
	// Upsert inserts or updates a translation (conflict on tenant_id, entity_type, entity_id, locale, field).
	Upsert(ctx context.Context, t *Translation) error

	// GetByEntity returns all translations for a specific entity+locale combination.
	GetByEntity(ctx context.Context, tenantID kernel.TenantID, entityType, entityID, locale string) ([]Translation, error)

	// GetBundle returns a TranslationBundle (grouped fields map) for a specific entity+locale.
	GetBundle(ctx context.Context, tenantID kernel.TenantID, entityType, entityID, locale string) (*TranslationBundle, error)

	// ListLocales returns all distinct locales available for a given entity.
	ListLocales(ctx context.Context, tenantID kernel.TenantID, entityType, entityID string) ([]string, error)

	// Delete removes a single translated field for a specific entity+locale+field.
	Delete(ctx context.Context, tenantID kernel.TenantID, entityType, entityID, locale, field string) error

	// DeleteAll removes all translations for a given entity (all locales, all fields).
	DeleteAll(ctx context.Context, tenantID kernel.TenantID, entityType, entityID string) error
}
