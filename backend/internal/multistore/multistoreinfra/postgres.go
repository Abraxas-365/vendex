package multistoreinfra

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/multistore"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements multistore.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed multistore repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Compile-time interface check.
var _ multistore.Repository = (*PostgresRepo)(nil)

// ---------------------------------------------------------------------------
// Storefront CRUD
// ---------------------------------------------------------------------------

// Create persists a new storefront.
func (r *PostgresRepo) Create(ctx context.Context, sf *multistore.Storefront) error {
	settingsJSON, err := json.Marshal(sf.Settings)
	if err != nil {
		return errx.Wrap(err, "marshaling storefront settings", errx.TypeInternal)
	}

	_, dbErr := r.db.ExecContext(ctx, `
		INSERT INTO storefronts (
			id, tenant_id, name, slug, domain, description,
			theme_id, logo_url, default_locale, default_currency,
			is_active, is_default, settings, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13, $14, $15
		)`,
		string(sf.ID), string(sf.TenantID), sf.Name, sf.Slug, sf.Domain, sf.Description,
		sf.ThemeID, sf.LogoURL, sf.DefaultLocale, sf.DefaultCurrency,
		sf.IsActive, sf.IsDefault, settingsJSON, sf.CreatedAt, sf.UpdatedAt,
	)
	if dbErr != nil {
		if isUniqueViolation(dbErr) {
			return multistore.ErrSlugConflict
		}
		return errx.Wrap(dbErr, "inserting storefront", errx.TypeInternal)
	}
	return nil
}

// GetByID returns a storefront scoped to the tenant.
func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.StorefrontEntryID) (*multistore.Storefront, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, slug, domain, description,
		       theme_id, logo_url, default_locale, default_currency,
		       is_active, is_default, settings, created_at, updated_at
		FROM storefronts
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	sf, err := scanStorefront(row.Scan)
	if err == sql.ErrNoRows {
		return nil, multistore.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning storefront", errx.TypeInternal)
	}
	return sf, nil
}

// GetBySlug returns a storefront by slug within a tenant.
func (r *PostgresRepo) GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*multistore.Storefront, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, slug, domain, description,
		       theme_id, logo_url, default_locale, default_currency,
		       is_active, is_default, settings, created_at, updated_at
		FROM storefronts
		WHERE slug = $1 AND tenant_id = $2`,
		slug, string(tenantID),
	)
	sf, err := scanStorefront(row.Scan)
	if err == sql.ErrNoRows {
		return nil, multistore.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning storefront by slug", errx.TypeInternal)
	}
	return sf, nil
}

// GetByDomain returns a storefront by its custom domain (global lookup).
func (r *PostgresRepo) GetByDomain(ctx context.Context, domain string) (*multistore.Storefront, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, slug, domain, description,
		       theme_id, logo_url, default_locale, default_currency,
		       is_active, is_default, settings, created_at, updated_at
		FROM storefronts
		WHERE domain = $1`,
		domain,
	)
	sf, err := scanStorefront(row.Scan)
	if err == sql.ErrNoRows {
		return nil, multistore.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning storefront by domain", errx.TypeInternal)
	}
	return sf, nil
}

// List returns paginated storefronts for a tenant.
func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, page, pageSize int) (kernel.Paginated[multistore.Storefront], error) {
	var zero kernel.Paginated[multistore.Storefront]
	pg := kernel.NewPaginationOptions(page, pageSize)

	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM storefronts WHERE tenant_id = $1`,
		string(tenantID),
	).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting storefronts", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, name, slug, domain, description,
		       theme_id, logo_url, default_locale, default_currency,
		       is_active, is_default, settings, created_at, updated_at
		FROM storefronts
		WHERE tenant_id = $1
		ORDER BY is_default DESC, created_at ASC
		LIMIT $2 OFFSET $3`,
		string(tenantID), pg.Limit(), pg.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "querying storefronts", errx.TypeInternal)
	}
	defer rows.Close()

	var items []multistore.Storefront
	for rows.Next() {
		sf, err := scanStorefront(rows.Scan)
		if err != nil {
			return zero, errx.Wrap(err, "scanning storefront row", errx.TypeInternal)
		}
		items = append(items, *sf)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating storefronts", errx.TypeInternal)
	}
	if items == nil {
		items = []multistore.Storefront{}
	}

	return kernel.NewPaginated(items, pg.Page, pg.PageSize, total), nil
}

// Update persists changes to an existing storefront.
func (r *PostgresRepo) Update(ctx context.Context, sf *multistore.Storefront) error {
	settingsJSON, err := json.Marshal(sf.Settings)
	if err != nil {
		return errx.Wrap(err, "marshaling storefront settings", errx.TypeInternal)
	}

	_, dbErr := r.db.ExecContext(ctx, `
		UPDATE storefronts SET
			name = $1, domain = $2, description = $3,
			theme_id = $4, logo_url = $5, default_locale = $6,
			default_currency = $7, is_active = $8, settings = $9, updated_at = $10
		WHERE id = $11 AND tenant_id = $12`,
		sf.Name, sf.Domain, sf.Description,
		sf.ThemeID, sf.LogoURL, sf.DefaultLocale,
		sf.DefaultCurrency, sf.IsActive, settingsJSON, sf.UpdatedAt,
		string(sf.ID), string(sf.TenantID),
	)
	if dbErr != nil {
		if isUniqueViolation(dbErr) {
			return multistore.ErrDomainConflict
		}
		return errx.Wrap(dbErr, "updating storefront", errx.TypeInternal)
	}
	return nil
}

// Delete removes a storefront.
func (r *PostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.StorefrontEntryID) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM storefronts WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting storefront", errx.TypeInternal)
	}
	return nil
}

// ClearDefault clears the is_default flag for all storefronts of a tenant.
func (r *PostgresRepo) ClearDefault(ctx context.Context, tenantID kernel.TenantID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE storefronts SET is_default = false WHERE tenant_id = $1`,
		string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "clearing default storefronts", errx.TypeInternal)
	}
	return nil
}

// SetDefault marks one storefront as default and clears all others in a single transaction.
func (r *PostgresRepo) SetDefault(ctx context.Context, tenantID kernel.TenantID, id kernel.StorefrontEntryID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errx.Wrap(err, "beginning transaction", errx.TypeInternal)
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE storefronts SET is_default = false WHERE tenant_id = $1`,
		string(tenantID),
	); err != nil {
		_ = tx.Rollback()
		return errx.Wrap(err, "clearing default storefronts", errx.TypeInternal)
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE storefronts SET is_default = true WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	); err != nil {
		_ = tx.Rollback()
		return errx.Wrap(err, "setting default storefront", errx.TypeInternal)
	}

	if err := tx.Commit(); err != nil {
		return errx.Wrap(err, "committing set-default transaction", errx.TypeInternal)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Catalog links
// ---------------------------------------------------------------------------

// AddCatalog links a catalog to a storefront.
func (r *PostgresRepo) AddCatalog(ctx context.Context, sc *multistore.StorefrontCatalog) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO storefront_catalogs (
			id, tenant_id, storefront_id, catalog_id, sort_order, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)`,
		string(sc.ID), string(sc.TenantID), string(sc.StorefrontID),
		sc.CatalogID, sc.SortOrder, sc.CreatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return multistore.ErrCatalogConflict
		}
		return errx.Wrap(err, "inserting storefront catalog", errx.TypeInternal)
	}
	return nil
}

// RemoveCatalog removes a catalog link from a storefront.
func (r *PostgresRepo) RemoveCatalog(ctx context.Context, tenantID kernel.TenantID, storefrontID kernel.StorefrontEntryID, catalogID string) error {
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM storefront_catalogs
		WHERE tenant_id = $1 AND storefront_id = $2 AND catalog_id = $3`,
		string(tenantID), string(storefrontID), catalogID,
	)
	if err != nil {
		return errx.Wrap(err, "removing storefront catalog", errx.TypeInternal)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return multistore.ErrCatalogNotFound
	}
	return nil
}

// ListCatalogs returns all catalog links for a storefront.
func (r *PostgresRepo) ListCatalogs(ctx context.Context, tenantID kernel.TenantID, storefrontID kernel.StorefrontEntryID) ([]multistore.StorefrontCatalog, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, storefront_id, catalog_id, sort_order, created_at
		FROM storefront_catalogs
		WHERE tenant_id = $1 AND storefront_id = $2
		ORDER BY sort_order ASC, created_at ASC`,
		string(tenantID), string(storefrontID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying storefront catalogs", errx.TypeInternal)
	}
	defer rows.Close()

	var items []multistore.StorefrontCatalog
	for rows.Next() {
		var sc multistore.StorefrontCatalog
		var id, tenantIDStr, storefrontIDStr string
		if err := rows.Scan(&id, &tenantIDStr, &storefrontIDStr, &sc.CatalogID, &sc.SortOrder, &sc.CreatedAt); err != nil {
			return nil, errx.Wrap(err, "scanning storefront catalog row", errx.TypeInternal)
		}
		sc.ID = kernel.NewStorefrontCatalogID(id)
		sc.TenantID = kernel.TenantID(tenantIDStr)
		sc.StorefrontID = kernel.NewStorefrontEntryID(storefrontIDStr)
		items = append(items, sc)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating storefront catalogs", errx.TypeInternal)
	}
	if items == nil {
		items = []multistore.StorefrontCatalog{}
	}
	return items, nil
}

// ---------------------------------------------------------------------------
// Scan helpers
// ---------------------------------------------------------------------------

type scanFunc func(dest ...any) error

func scanStorefront(scan scanFunc) (*multistore.Storefront, error) {
	var sf multistore.Storefront
	var id, tenantIDStr string
	var settingsRaw []byte
	var domain sql.NullString

	err := scan(
		&id, &tenantIDStr, &sf.Name, &sf.Slug, &domain, &sf.Description,
		&sf.ThemeID, &sf.LogoURL, &sf.DefaultLocale, &sf.DefaultCurrency,
		&sf.IsActive, &sf.IsDefault, &settingsRaw, &sf.CreatedAt, &sf.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	sf.ID = kernel.NewStorefrontEntryID(id)
	sf.TenantID = kernel.TenantID(tenantIDStr)
	if domain.Valid {
		sf.Domain = &domain.String
	}

	if len(settingsRaw) > 0 {
		if err := json.Unmarshal(settingsRaw, &sf.Settings); err != nil {
			sf.Settings = map[string]interface{}{}
		}
	} else {
		sf.Settings = map[string]interface{}{}
	}

	return &sf, nil
}

// isUniqueViolation reports whether err is a PostgreSQL unique-constraint violation.
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	return containsStr(err.Error(), "23505") || containsStr(err.Error(), "unique constraint") || containsStr(err.Error(), "duplicate key")
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsSubstring(s, sub))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
