package i18ninfra

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/i18n"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements i18n.Repository using PostgreSQL.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed i18n repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Upsert inserts or updates a translation row using ON CONFLICT.
func (r *PostgresRepo) Upsert(ctx context.Context, t *i18n.Translation) error {
	const q = `
		INSERT INTO translations (id, tenant_id, entity_type, entity_id, locale, field, value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (tenant_id, entity_type, entity_id, locale, field)
		DO UPDATE SET value = $7, updated_at = $9`

	_, err := r.db.ExecContext(ctx, q,
		string(t.ID),
		string(t.TenantID),
		t.EntityType,
		t.EntityID,
		t.Locale,
		t.Field,
		t.Value,
		t.CreatedAt,
		t.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "upserting translation", errx.TypeInternal)
	}
	return nil
}

// GetByEntity returns all translations for a given entity+locale combination.
func (r *PostgresRepo) GetByEntity(ctx context.Context, tenantID kernel.TenantID, entityType, entityID, locale string) ([]i18n.Translation, error) {
	const q = `
		SELECT id, tenant_id, entity_type, entity_id, locale, field, value, created_at, updated_at
		FROM translations
		WHERE tenant_id = $1 AND entity_type = $2 AND entity_id = $3 AND locale = $4
		ORDER BY field ASC`

	rows, err := r.db.QueryContext(ctx, q, string(tenantID), entityType, entityID, locale)
	if err != nil {
		return nil, errx.Wrap(err, "querying translations", errx.TypeInternal)
	}
	defer rows.Close()

	var results []i18n.Translation
	for rows.Next() {
		var t i18n.Translation
		if err := rows.Scan(
			&t.ID, &t.TenantID, &t.EntityType, &t.EntityID,
			&t.Locale, &t.Field, &t.Value, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, errx.Wrap(err, "scanning translation row", errx.TypeInternal)
		}
		results = append(results, t)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating translation rows", errx.TypeInternal)
	}
	return results, nil
}

// GetBundle returns a TranslationBundle (fields map) for a given entity+locale.
func (r *PostgresRepo) GetBundle(ctx context.Context, tenantID kernel.TenantID, entityType, entityID, locale string) (*i18n.TranslationBundle, error) {
	translations, err := r.GetByEntity(ctx, tenantID, entityType, entityID, locale)
	if err != nil {
		return nil, err
	}

	bundle := &i18n.TranslationBundle{
		EntityType: entityType,
		EntityID:   entityID,
		Locale:     locale,
		Fields:     make(map[string]string),
	}
	for _, t := range translations {
		bundle.Fields[t.Field] = t.Value
	}
	return bundle, nil
}

// ListLocales returns distinct locales available for the given entity.
func (r *PostgresRepo) ListLocales(ctx context.Context, tenantID kernel.TenantID, entityType, entityID string) ([]string, error) {
	const q = `
		SELECT DISTINCT locale
		FROM translations
		WHERE tenant_id = $1 AND entity_type = $2 AND entity_id = $3
		ORDER BY locale ASC`

	rows, err := r.db.QueryContext(ctx, q, string(tenantID), entityType, entityID)
	if err != nil {
		return nil, errx.Wrap(err, "listing locales", errx.TypeInternal)
	}
	defer rows.Close()

	var locales []string
	for rows.Next() {
		var locale string
		if err := rows.Scan(&locale); err != nil {
			return nil, errx.Wrap(err, "scanning locale row", errx.TypeInternal)
		}
		locales = append(locales, locale)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating locale rows", errx.TypeInternal)
	}
	return locales, nil
}

// Delete removes a single translated field for an entity+locale+field.
func (r *PostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, entityType, entityID, locale, field string) error {
	const q = `
		DELETE FROM translations
		WHERE tenant_id = $1 AND entity_type = $2 AND entity_id = $3 AND locale = $4 AND field = $5`

	result, err := r.db.ExecContext(ctx, q, string(tenantID), entityType, entityID, locale, field)
	if err != nil {
		return errx.Wrap(err, "deleting translation", errx.TypeInternal)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errx.Wrap(err, "checking rows affected", errx.TypeInternal)
	}
	if rowsAffected == 0 {
		return i18n.ErrNotFound
	}
	return nil
}

// DeleteAll removes all translations for an entity (all locales, all fields).
func (r *PostgresRepo) DeleteAll(ctx context.Context, tenantID kernel.TenantID, entityType, entityID string) error {
	const q = `
		DELETE FROM translations
		WHERE tenant_id = $1 AND entity_type = $2 AND entity_id = $3`

	_, err := r.db.ExecContext(ctx, q, string(tenantID), entityType, entityID)
	if err != nil {
		return errx.Wrap(err, "deleting all translations for entity", errx.TypeInternal)
	}
	return nil
}

// ensure interface is satisfied at compile time
var _ i18n.Repository = (*PostgresRepo)(nil)
