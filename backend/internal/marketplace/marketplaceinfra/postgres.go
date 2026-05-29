package marketplaceinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/marketplace"
)

// PostgresRepo implements marketplace.Repository using database/sql.
type PostgresRepo struct {
	db *sql.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed marketplace repository.
func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// ---------------------------------------------------------------------------
// Plugins
// ---------------------------------------------------------------------------

func (r *PostgresRepo) CreatePlugin(ctx context.Context, p *marketplace.Plugin) error {
	tagsJSON, err := json.Marshal(p.Tags)
	if err != nil {
		return fmt.Errorf("marshaling tags: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO plugins (id, name, display_name, description, author, icon, category, tags, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		string(p.ID), p.Name, p.DisplayName, p.Description,
		p.Author, p.Icon, string(p.Category), string(tagsJSON),
		p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting plugin: %w", err)
	}
	return nil
}

func (r *PostgresRepo) GetPlugin(ctx context.Context, id kernel.PluginID) (*marketplace.Plugin, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, display_name, description, author, icon, category, tags, created_at, updated_at
		FROM plugins WHERE id = $1`,
		string(id),
	)
	return scanPlugin(row)
}

func (r *PostgresRepo) GetPluginByName(ctx context.Context, name string) (*marketplace.Plugin, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, display_name, description, author, icon, category, tags, created_at, updated_at
		FROM plugins WHERE name = $1`,
		name,
	)
	return scanPlugin(row)
}

func (r *PostgresRepo) ListPlugins(ctx context.Context, pg kernel.Pagination) (kernel.PaginatedResult[marketplace.Plugin], error) {
	return r.queryPlugins(ctx, pg, "", nil)
}

func (r *PostgresRepo) ListPluginsByCategory(ctx context.Context, cat marketplace.PluginCategory, pg kernel.Pagination) (kernel.PaginatedResult[marketplace.Plugin], error) {
	return r.queryPlugins(ctx, pg, "WHERE category = $1", []any{string(cat)})
}

func (r *PostgresRepo) queryPlugins(ctx context.Context, pg kernel.Pagination, where string, whereArgs []any) (kernel.PaginatedResult[marketplace.Plugin], error) {
	var zero kernel.PaginatedResult[marketplace.Plugin]

	countQuery := "SELECT COUNT(*) FROM plugins " + where
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&total); err != nil {
		return zero, fmt.Errorf("counting plugins: %w", err)
	}

	nextParam := len(whereArgs) + 1
	dataQuery := fmt.Sprintf(`
		SELECT id, name, display_name, description, author, icon, category, tags, created_at, updated_at
		FROM plugins %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, where, nextParam, nextParam+1)
	args := append(whereArgs, pg.Limit(), pg.Offset())

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return zero, fmt.Errorf("querying plugins: %w", err)
	}
	defer rows.Close()

	var plugins []marketplace.Plugin
	for rows.Next() {
		p, err := scanPluginRow(rows)
		if err != nil {
			return zero, err
		}
		plugins = append(plugins, *p)
	}
	if err := rows.Err(); err != nil {
		return zero, fmt.Errorf("iterating plugins: %w", err)
	}

	return kernel.NewPaginatedResult(plugins, total, pg), nil
}

// ---------------------------------------------------------------------------
// Versions
// ---------------------------------------------------------------------------

func (r *PostgresRepo) CreateVersion(ctx context.Context, v *marketplace.PluginVersion) error {
	permsJSON, err := json.Marshal(v.Permissions)
	if err != nil {
		return fmt.Errorf("marshaling permissions: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO plugin_versions (id, plugin_id, version, changelog, permissions, manifest_json, frontend_url, backend_entry, min_platform_ver, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		string(v.ID), string(v.PluginID), v.Version, v.Changelog,
		string(permsJSON), v.ManifestJSON, v.FrontendURL, v.BackendEntry,
		v.MinPlatformVer, v.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting plugin version: %w", err)
	}
	return nil
}

func (r *PostgresRepo) GetVersion(ctx context.Context, id kernel.PluginVersionID) (*marketplace.PluginVersion, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, plugin_id, version, changelog, permissions, manifest_json, frontend_url, backend_entry, min_platform_ver, created_at
		FROM plugin_versions WHERE id = $1`,
		string(id),
	)
	return scanVersion(row)
}

func (r *PostgresRepo) GetLatestVersion(ctx context.Context, pluginID kernel.PluginID) (*marketplace.PluginVersion, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, plugin_id, version, changelog, permissions, manifest_json, frontend_url, backend_entry, min_platform_ver, created_at
		FROM plugin_versions WHERE plugin_id = $1
		ORDER BY created_at DESC LIMIT 1`,
		string(pluginID),
	)
	return scanVersion(row)
}

func (r *PostgresRepo) ListVersions(ctx context.Context, pluginID kernel.PluginID) ([]marketplace.PluginVersion, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, plugin_id, version, changelog, permissions, manifest_json, frontend_url, backend_entry, min_platform_ver, created_at
		FROM plugin_versions WHERE plugin_id = $1
		ORDER BY created_at DESC`,
		string(pluginID),
	)
	if err != nil {
		return nil, fmt.Errorf("querying versions: %w", err)
	}
	defer rows.Close()

	var versions []marketplace.PluginVersion
	for rows.Next() {
		v, err := scanVersionRow(rows)
		if err != nil {
			return nil, err
		}
		versions = append(versions, *v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating versions: %w", err)
	}
	return versions, nil
}

// ---------------------------------------------------------------------------
// Installations
// ---------------------------------------------------------------------------

func (r *PostgresRepo) CreateInstallation(ctx context.Context, inst *marketplace.Installation) error {
	settingsJSON, err := json.Marshal(inst.Settings)
	if err != nil {
		return fmt.Errorf("marshaling settings: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO plugin_installations (id, tenant_id, plugin_id, version_id, status, settings, installed_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		string(inst.ID), string(inst.TenantID), string(inst.PluginID), string(inst.VersionID),
		string(inst.Status), string(settingsJSON), inst.InstalledAt, inst.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting installation: %w", err)
	}
	return nil
}

func (r *PostgresRepo) GetInstallation(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID) (*marketplace.Installation, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, plugin_id, version_id, status, settings, installed_at, updated_at
		FROM plugin_installations WHERE tenant_id = $1 AND plugin_id = $2`,
		string(tenantID), string(pluginID),
	)
	return scanInstallation(row)
}

func (r *PostgresRepo) UpdateInstallation(ctx context.Context, inst *marketplace.Installation) error {
	settingsJSON, err := json.Marshal(inst.Settings)
	if err != nil {
		return fmt.Errorf("marshaling settings: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE plugin_installations SET version_id=$1, status=$2, settings=$3, updated_at=$4
		WHERE tenant_id=$5 AND plugin_id=$6`,
		string(inst.VersionID), string(inst.Status), string(settingsJSON), inst.UpdatedAt,
		string(inst.TenantID), string(inst.PluginID),
	)
	if err != nil {
		return fmt.Errorf("updating installation: %w", err)
	}
	return nil
}

func (r *PostgresRepo) DeleteInstallation(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM plugin_installations WHERE tenant_id = $1 AND plugin_id = $2`,
		string(tenantID), string(pluginID),
	)
	if err != nil {
		return fmt.Errorf("deleting installation: %w", err)
	}
	return nil
}

func (r *PostgresRepo) ListInstallations(ctx context.Context, tenantID kernel.TenantID) ([]marketplace.Installation, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, plugin_id, version_id, status, settings, installed_at, updated_at
		FROM plugin_installations WHERE tenant_id = $1
		ORDER BY installed_at DESC`,
		string(tenantID),
	)
	if err != nil {
		return nil, fmt.Errorf("querying installations: %w", err)
	}
	defer rows.Close()

	var installations []marketplace.Installation
	for rows.Next() {
		inst, err := scanInstallationRow(rows)
		if err != nil {
			return nil, err
		}
		installations = append(installations, *inst)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating installations: %w", err)
	}
	return installations, nil
}

// ---------------------------------------------------------------------------
// Scanners
// ---------------------------------------------------------------------------

// scanner is implemented by both *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...any) error
}

func scanPlugin(row *sql.Row) (*marketplace.Plugin, error) {
	p, err := scanPluginFields(row)
	if err == sql.ErrNoRows {
		return nil, marketplace.ErrPluginNotFound
	}
	return p, err
}

func scanPluginRow(rows *sql.Rows) (*marketplace.Plugin, error) {
	return scanPluginFields(rows)
}

func scanPluginFields(s scanner) (*marketplace.Plugin, error) {
	var p marketplace.Plugin
	var id, category string
	var tagsJSON string
	var createdAt, updatedAt time.Time

	err := s.Scan(
		&id, &p.Name, &p.DisplayName, &p.Description,
		&p.Author, &p.Icon, &category, &tagsJSON,
		&createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	p.ID = kernel.PluginID(id)
	p.Category = marketplace.PluginCategory(category)
	p.CreatedAt = createdAt
	p.UpdatedAt = updatedAt

	_ = json.Unmarshal([]byte(tagsJSON), &p.Tags)
	if p.Tags == nil {
		p.Tags = []string{}
	}

	return &p, nil
}

func scanVersion(row *sql.Row) (*marketplace.PluginVersion, error) {
	v, err := scanVersionFields(row)
	if err == sql.ErrNoRows {
		return nil, marketplace.ErrVersionNotFound
	}
	return v, err
}

func scanVersionRow(rows *sql.Rows) (*marketplace.PluginVersion, error) {
	return scanVersionFields(rows)
}

func scanVersionFields(s scanner) (*marketplace.PluginVersion, error) {
	var v marketplace.PluginVersion
	var id, pluginID string
	var permsJSON string
	var createdAt time.Time

	err := s.Scan(
		&id, &pluginID, &v.Version, &v.Changelog,
		&permsJSON, &v.ManifestJSON, &v.FrontendURL, &v.BackendEntry,
		&v.MinPlatformVer, &createdAt,
	)
	if err != nil {
		return nil, err
	}

	v.ID = kernel.PluginVersionID(id)
	v.PluginID = kernel.PluginID(pluginID)
	v.CreatedAt = createdAt

	_ = json.Unmarshal([]byte(permsJSON), &v.Permissions)
	if v.Permissions == nil {
		v.Permissions = []string{}
	}

	return &v, nil
}

func scanInstallation(row *sql.Row) (*marketplace.Installation, error) {
	inst, err := scanInstallationFields(row)
	if err == sql.ErrNoRows {
		return nil, marketplace.ErrNotInstalled
	}
	return inst, err
}

func scanInstallationRow(rows *sql.Rows) (*marketplace.Installation, error) {
	return scanInstallationFields(rows)
}

func scanInstallationFields(s scanner) (*marketplace.Installation, error) {
	var inst marketplace.Installation
	var id, tenantID, pluginID, versionID, status string
	var settingsJSON string
	var installedAt, updatedAt time.Time

	err := s.Scan(
		&id, &tenantID, &pluginID, &versionID,
		&status, &settingsJSON, &installedAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	inst.ID = kernel.InstallationID(id)
	inst.TenantID = kernel.TenantID(tenantID)
	inst.PluginID = kernel.PluginID(pluginID)
	inst.VersionID = kernel.PluginVersionID(versionID)
	inst.Status = marketplace.InstallationStatus(status)
	inst.InstalledAt = installedAt
	inst.UpdatedAt = updatedAt

	_ = json.Unmarshal([]byte(settingsJSON), &inst.Settings)
	if inst.Settings == nil {
		inst.Settings = map[string]any{}
	}

	return &inst, nil
}

// Ensure interface compliance.
var _ marketplace.Repository = (*PostgresRepo)(nil)
