package themeinfra

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/theme"
)

// Compile-time interface check.
var _ theme.ThemeRepository = (*PostgresRepo)(nil)

// PostgresRepo implements theme.ThemeRepository using PostgreSQL.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed theme repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Create inserts a new theme into the database.
func (r *PostgresRepo) Create(ctx context.Context, t *theme.Theme) error {
	tokensJSON, err := json.Marshal(t.Tokens)
	if err != nil {
		return errx.Wrap(err, "marshaling theme tokens", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO themes (id, tenant_id, name, is_active, tokens, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		string(t.ID), string(t.TenantID), t.Name, t.IsActive,
		string(tokensJSON), t.CreatedAt, t.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting theme", errx.TypeInternal)
	}
	return nil
}

// GetByID retrieves a theme by its ID, scoped to a tenant.
func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ThemeID) (*theme.Theme, error) {
	return r.queryOne(ctx,
		`SELECT id, tenant_id, name, is_active, tokens, created_at, updated_at
		 FROM themes WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
}

// GetActive retrieves the currently active theme for a tenant.
func (r *PostgresRepo) GetActive(ctx context.Context, tenantID kernel.TenantID) (*theme.Theme, error) {
	return r.queryOne(ctx,
		`SELECT id, tenant_id, name, is_active, tokens, created_at, updated_at
		 FROM themes WHERE tenant_id = $1 AND is_active = true LIMIT 1`,
		string(tenantID),
	)
}

// List returns all themes for a tenant, ordered by creation date descending.
func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID) ([]theme.Theme, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, tenant_id, name, is_active, tokens, created_at, updated_at
		 FROM themes WHERE tenant_id = $1 ORDER BY created_at DESC`,
		string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "listing themes", errx.TypeInternal)
	}
	defer rows.Close()

	var themes []theme.Theme
	for rows.Next() {
		t, err := r.scanRow(rows)
		if err != nil {
			return nil, errx.Wrap(err, "scanning theme row", errx.TypeInternal)
		}
		themes = append(themes, *t)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating theme rows", errx.TypeInternal)
	}
	return themes, nil
}

// Update persists changes to an existing theme.
func (r *PostgresRepo) Update(ctx context.Context, t *theme.Theme) error {
	tokensJSON, err := json.Marshal(t.Tokens)
	if err != nil {
		return errx.Wrap(err, "marshaling theme tokens", errx.TypeInternal)
	}

	result, err := r.db.ExecContext(ctx, `
		UPDATE themes
		SET name = $1, is_active = $2, tokens = $3, updated_at = $4
		WHERE id = $5 AND tenant_id = $6`,
		t.Name, t.IsActive, string(tokensJSON), t.UpdatedAt,
		string(t.ID), string(t.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating theme", errx.TypeInternal)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errx.Wrap(err, "checking rows affected", errx.TypeInternal)
	}
	if rows == 0 {
		return theme.ErrThemeNotFound
	}
	return nil
}

// Delete removes a theme from the database.
func (r *PostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.ThemeID) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM themes WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting theme", errx.TypeInternal)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errx.Wrap(err, "checking rows affected", errx.TypeInternal)
	}
	if rows == 0 {
		return theme.ErrThemeNotFound
	}
	return nil
}

// DeactivateAll sets is_active = false for all themes belonging to a tenant.
func (r *PostgresRepo) DeactivateAll(ctx context.Context, tenantID kernel.TenantID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE themes SET is_active = false WHERE tenant_id = $1`,
		string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deactivating all themes", errx.TypeInternal)
	}
	return nil
}

// --- helpers ---

type rowScanner interface {
	Scan(dest ...any) error
}

func (r *PostgresRepo) queryOne(ctx context.Context, query string, args ...any) (*theme.Theme, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	t, err := r.scanRow(row)
	if err == sql.ErrNoRows {
		return nil, theme.ErrThemeNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "querying theme", errx.TypeInternal)
	}
	return t, nil
}

func (r *PostgresRepo) scanRow(scanner rowScanner) (*theme.Theme, error) {
	var t theme.Theme
	var id, tenantID, tokensJSON string

	err := scanner.Scan(
		&id, &tenantID, &t.Name, &t.IsActive,
		&tokensJSON, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	t.ID = kernel.ThemeID(id)
	t.TenantID = kernel.TenantID(tenantID)
	_ = json.Unmarshal([]byte(tokensJSON), &t.Tokens)

	return &t, nil
}
