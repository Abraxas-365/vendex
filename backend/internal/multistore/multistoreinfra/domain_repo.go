package multistoreinfra

import (
	"context"
	"database/sql"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/multistore"
)

// DomainRepo implements multistore.TenantDomainRepo for tenant resolution.
type DomainRepo struct {
	db *sql.DB
}

// NewDomainRepo creates a new domain resolution repository.
func NewDomainRepo(db *sql.DB) *DomainRepo {
	return &DomainRepo{db: db}
}

// GetTenantBySlug resolves a tenant ID from its slug (for subdomain resolution).
func (r *DomainRepo) GetTenantBySlug(ctx context.Context, slug string) (kernel.TenantID, error) {
	var id string
	err := r.db.QueryRowContext(ctx,
		`SELECT id FROM tenants WHERE slug = $1 AND status IN ('ACTIVE', 'TRIAL')`,
		slug,
	).Scan(&id)
	if err == sql.ErrNoRows {
		return "", multistore.ErrTenantNotResolved
	}
	if err != nil {
		return "", errx.Wrap(err, "querying tenant by slug", errx.TypeInternal)
	}
	return kernel.TenantID(id), nil
}

// GetTenantByDomain resolves a tenant ID from a custom domain.
func (r *DomainRepo) GetTenantByDomain(ctx context.Context, domain string) (kernel.TenantID, error) {
	var tenantID string
	err := r.db.QueryRowContext(ctx,
		`SELECT tenant_id FROM tenant_domains WHERE domain = $1 AND is_verified = TRUE`,
		domain,
	).Scan(&tenantID)
	if err == sql.ErrNoRows {
		return "", multistore.ErrTenantNotResolved
	}
	if err != nil {
		return "", errx.Wrap(err, "querying tenant by domain", errx.TypeInternal)
	}
	return kernel.TenantID(tenantID), nil
}
