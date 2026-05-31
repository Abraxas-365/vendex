package storefrontinfra

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/storefront/renderer"
)

// TemplateOverridePostgresRepo implements renderer.TemplateOverrideRepository using PostgreSQL.
type TemplateOverridePostgresRepo struct {
	db *sqlx.DB
}

// NewTemplateOverridePostgresRepo creates a new PostgreSQL-backed template override repository.
func NewTemplateOverridePostgresRepo(db *sqlx.DB) *TemplateOverridePostgresRepo {
	return &TemplateOverridePostgresRepo{db: db}
}

// GetByBlockType returns the tenant-specific template override for the given block type,
// or nil (no error) if no override exists.
func (r *TemplateOverridePostgresRepo) GetByBlockType(ctx context.Context, tenantID kernel.TenantID, blockType string) (*renderer.TemplateOverride, error) {
	const q = `
		SELECT block_type, template
		FROM template_overrides
		WHERE tenant_id  = $1
		  AND block_type = $2
		LIMIT 1
	`

	var override renderer.TemplateOverride
	err := r.db.QueryRowContext(ctx, q, string(tenantID), blockType).Scan(
		&override.BlockType,
		&override.Template,
	)
	if err == sql.ErrNoRows {
		return nil, nil // no override — caller falls back to built-in template
	}
	if err != nil {
		return nil, err
	}
	return &override, nil
}
