package plugininfra

import (
	"context"
	"database/sql"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/plugin"
	"github.com/jmoiron/sqlx"
)

// ============================================================================
// PluginRepo — global plugin catalogue
// ============================================================================

// PluginRepo implements plugin.PluginRepository using sqlx.
type PluginRepo struct {
	db *sqlx.DB
}

// NewPluginRepo creates a new Postgres-backed plugin repository.
func NewPluginRepo(db *sqlx.DB) *PluginRepo {
	return &PluginRepo{db: db}
}

func (r *PluginRepo) Create(ctx context.Context, p *plugin.Plugin) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO plugins (id, name, display_name, description, author, icon, category, tags, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		string(p.ID), p.Name, p.DisplayName, p.Description,
		p.Author, p.Icon, p.Category, string(p.Tags),
		p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting plugin", errx.TypeInternal)
	}
	return nil
}

func (r *PluginRepo) GetByID(ctx context.Context, id kernel.PluginID) (*plugin.Plugin, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, display_name, description, author, icon, category, tags, created_at, updated_at
		FROM plugins WHERE id = $1`,
		string(id),
	)
	p, err := scanPlugin(row)
	if err == sql.ErrNoRows {
		return nil, plugin.ErrPluginNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning plugin", errx.TypeInternal)
	}
	return p, nil
}

func (r *PluginRepo) Update(ctx context.Context, p *plugin.Plugin) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE plugins
		SET name=$1, display_name=$2, description=$3, author=$4, icon=$5, category=$6, tags=$7, updated_at=$8
		WHERE id=$9`,
		p.Name, p.DisplayName, p.Description, p.Author,
		p.Icon, p.Category, string(p.Tags), p.UpdatedAt, string(p.ID),
	)
	if err != nil {
		return errx.Wrap(err, "updating plugin", errx.TypeInternal)
	}
	return nil
}

func (r *PluginRepo) Delete(ctx context.Context, id kernel.PluginID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM plugins WHERE id = $1`, string(id))
	if err != nil {
		return errx.Wrap(err, "deleting plugin", errx.TypeInternal)
	}
	return nil
}

func (r *PluginRepo) List(ctx context.Context, pg kernel.PaginationOptions) (kernel.Paginated[plugin.Plugin], error) {
	var zero kernel.Paginated[plugin.Plugin]

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM plugins`).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting plugins", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, display_name, description, author, icon, category, tags, created_at, updated_at
		FROM plugins
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`,
		pg.Limit(), pg.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "listing plugins", errx.TypeInternal)
	}
	defer rows.Close()

	var plugins []plugin.Plugin
	for rows.Next() {
		p, err := scanPluginRow(rows)
		if err != nil {
			return zero, errx.Wrap(err, "scanning plugin row", errx.TypeInternal)
		}
		plugins = append(plugins, *p)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating plugins", errx.TypeInternal)
	}
	if plugins == nil {
		plugins = []plugin.Plugin{}
	}

	return kernel.NewPaginated(plugins, pg.Page, pg.PageSize, total), nil
}

// Compile-time interface check.
var _ plugin.PluginRepository = (*PluginRepo)(nil)

// ============================================================================
// VersionRepo — plugin versions
// ============================================================================

// VersionRepo implements plugin.PluginVersionRepository using sqlx.
type VersionRepo struct {
	db *sqlx.DB
}

// NewVersionRepo creates a new Postgres-backed plugin version repository.
func NewVersionRepo(db *sqlx.DB) *VersionRepo {
	return &VersionRepo{db: db}
}

func (r *VersionRepo) Create(ctx context.Context, v *plugin.PluginVersion) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO plugin_versions (id, plugin_id, version, changelog, permissions, manifest_json, frontend_url, backend_entry, min_platform_ver, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		string(v.ID), string(v.PluginID), v.Version, v.Changelog,
		string(v.Permissions), v.ManifestJSON, v.FrontendURL,
		v.BackendEntry, v.MinPlatformVer, v.CreatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting plugin version", errx.TypeInternal)
	}
	return nil
}

func (r *VersionRepo) GetByID(ctx context.Context, id kernel.PluginVersionID) (*plugin.PluginVersion, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, plugin_id, version, changelog, permissions, manifest_json, frontend_url, backend_entry, min_platform_ver, created_at
		FROM plugin_versions WHERE id = $1`,
		string(id),
	)
	v, err := scanVersion(row)
	if err == sql.ErrNoRows {
		return nil, plugin.ErrVersionNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning plugin version", errx.TypeInternal)
	}
	return v, nil
}

func (r *VersionRepo) ListByPlugin(ctx context.Context, pluginID kernel.PluginID, pg kernel.PaginationOptions) (kernel.Paginated[plugin.PluginVersion], error) {
	var zero kernel.Paginated[plugin.PluginVersion]

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM plugin_versions WHERE plugin_id = $1`, string(pluginID)).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting plugin versions", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, plugin_id, version, changelog, permissions, manifest_json, frontend_url, backend_entry, min_platform_ver, created_at
		FROM plugin_versions WHERE plugin_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		string(pluginID), pg.Limit(), pg.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "listing plugin versions", errx.TypeInternal)
	}
	defer rows.Close()

	var versions []plugin.PluginVersion
	for rows.Next() {
		v, err := scanVersionRow(rows)
		if err != nil {
			return zero, errx.Wrap(err, "scanning version row", errx.TypeInternal)
		}
		versions = append(versions, *v)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating versions", errx.TypeInternal)
	}
	if versions == nil {
		versions = []plugin.PluginVersion{}
	}

	return kernel.NewPaginated(versions, pg.Page, pg.PageSize, total), nil
}

func (r *VersionRepo) GetLatest(ctx context.Context, pluginID kernel.PluginID) (*plugin.PluginVersion, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, plugin_id, version, changelog, permissions, manifest_json, frontend_url, backend_entry, min_platform_ver, created_at
		FROM plugin_versions WHERE plugin_id = $1
		ORDER BY created_at DESC
		LIMIT 1`,
		string(pluginID),
	)
	v, err := scanVersion(row)
	if err == sql.ErrNoRows {
		return nil, plugin.ErrVersionNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning latest plugin version", errx.TypeInternal)
	}
	return v, nil
}

// Compile-time interface check.
var _ plugin.PluginVersionRepository = (*VersionRepo)(nil)

// ============================================================================
// InstallationRepo — per-tenant plugin installations
// ============================================================================

// InstallationRepo implements plugin.InstallationRepository using sqlx.
type InstallationRepo struct {
	db *sqlx.DB
}

// NewInstallationRepo creates a new Postgres-backed installation repository.
func NewInstallationRepo(db *sqlx.DB) *InstallationRepo {
	return &InstallationRepo{db: db}
}

func (r *InstallationRepo) Create(ctx context.Context, i *plugin.PluginInstallation) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO plugin_installations (id, tenant_id, plugin_id, version_id, status, settings, installed_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		string(i.ID), string(i.TenantID), string(i.PluginID),
		string(i.VersionID), string(i.Status), string(i.Settings),
		i.InstalledAt, i.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting plugin installation", errx.TypeInternal)
	}
	return nil
}

func (r *InstallationRepo) GetByTenantAndPlugin(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID) (*plugin.PluginInstallation, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, plugin_id, version_id, status, settings, installed_at, updated_at
		FROM plugin_installations WHERE tenant_id = $1 AND plugin_id = $2`,
		string(tenantID), string(pluginID),
	)
	inst, err := scanInstallation(row)
	if err == sql.ErrNoRows {
		return nil, plugin.ErrInstallationNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning installation", errx.TypeInternal)
	}
	return inst, nil
}

func (r *InstallationRepo) ListByTenant(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[plugin.PluginInstallation], error) {
	var zero kernel.Paginated[plugin.PluginInstallation]

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM plugin_installations WHERE tenant_id = $1`, string(tenantID)).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting installations", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, plugin_id, version_id, status, settings, installed_at, updated_at
		FROM plugin_installations WHERE tenant_id = $1
		ORDER BY installed_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), pg.Limit(), pg.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "listing installations", errx.TypeInternal)
	}
	defer rows.Close()

	var installations []plugin.PluginInstallation
	for rows.Next() {
		inst, err := scanInstallationRow(rows)
		if err != nil {
			return zero, errx.Wrap(err, "scanning installation row", errx.TypeInternal)
		}
		installations = append(installations, *inst)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating installations", errx.TypeInternal)
	}
	if installations == nil {
		installations = []plugin.PluginInstallation{}
	}

	return kernel.NewPaginated(installations, pg.Page, pg.PageSize, total), nil
}

func (r *InstallationRepo) ListActiveByTenant(ctx context.Context, tenantID kernel.TenantID) ([]plugin.PluginInstallation, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, plugin_id, version_id, status, settings, installed_at, updated_at
		FROM plugin_installations WHERE tenant_id = $1 AND status = 'active'
		ORDER BY installed_at DESC`,
		string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "listing active installations", errx.TypeInternal)
	}
	defer rows.Close()

	var installations []plugin.PluginInstallation
	for rows.Next() {
		inst, err := scanInstallationRow(rows)
		if err != nil {
			return nil, errx.Wrap(err, "scanning active installation row", errx.TypeInternal)
		}
		installations = append(installations, *inst)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating active installations", errx.TypeInternal)
	}
	if installations == nil {
		installations = []plugin.PluginInstallation{}
	}

	return installations, nil
}

func (r *InstallationRepo) GetJSManifestData(ctx context.Context, tenantID kernel.TenantID) ([]plugin.PluginScript, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			pi.plugin_id,
			p.name,
			pv.version,
			pv.frontend_url,
			pi.settings
		FROM plugin_installations pi
		JOIN plugins p ON p.id = pi.plugin_id
		JOIN plugin_versions pv ON pv.id = pi.version_id
		WHERE pi.tenant_id = $1
		  AND pi.status = 'active'
		  AND pv.frontend_url <> ''
		ORDER BY pi.installed_at DESC`,
		string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying js manifest data", errx.TypeInternal)
	}
	defer rows.Close()

	var scripts []plugin.PluginScript
	for rows.Next() {
		var pluginID, name, version, frontendURL, settings string
		if err := rows.Scan(&pluginID, &name, &version, &frontendURL, &settings); err != nil {
			return nil, errx.Wrap(err, "scanning js manifest row", errx.TypeInternal)
		}
		scripts = append(scripts, plugin.PluginScript{
			PluginID:   pluginID,
			PluginName: name,
			Version:    version,
			Src:        frontendURL,
			Config:     []byte(settings),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating js manifest rows", errx.TypeInternal)
	}
	if scripts == nil {
		scripts = []plugin.PluginScript{}
	}
	return scripts, nil
}

func (r *InstallationRepo) Update(ctx context.Context, i *plugin.PluginInstallation) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE plugin_installations
		SET version_id=$1, status=$2, settings=$3, updated_at=$4
		WHERE tenant_id=$5 AND plugin_id=$6`,
		string(i.VersionID), string(i.Status), string(i.Settings), i.UpdatedAt,
		string(i.TenantID), string(i.PluginID),
	)
	if err != nil {
		return errx.Wrap(err, "updating plugin installation", errx.TypeInternal)
	}
	return nil
}

func (r *InstallationRepo) Delete(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM plugin_installations WHERE tenant_id = $1 AND plugin_id = $2`,
		string(tenantID), string(pluginID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting plugin installation", errx.TypeInternal)
	}
	return nil
}

// Compile-time interface check.
var _ plugin.InstallationRepository = (*InstallationRepo)(nil)

// ============================================================================
// Scanner helpers
// ============================================================================

type rowScanner interface {
	Scan(dest ...any) error
}

func scanPlugin(row rowScanner) (*plugin.Plugin, error) {
	var p plugin.Plugin
	var id, tags string
	err := row.Scan(
		&id, &p.Name, &p.DisplayName, &p.Description,
		&p.Author, &p.Icon, &p.Category, &tags,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	p.ID = kernel.PluginID(id)
	p.Tags = []byte(tags)
	return &p, nil
}

func scanPluginRow(rows *sql.Rows) (*plugin.Plugin, error) {
	var p plugin.Plugin
	var id, tags string
	err := rows.Scan(
		&id, &p.Name, &p.DisplayName, &p.Description,
		&p.Author, &p.Icon, &p.Category, &tags,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	p.ID = kernel.PluginID(id)
	p.Tags = []byte(tags)
	return &p, nil
}

func scanVersion(row rowScanner) (*plugin.PluginVersion, error) {
	var v plugin.PluginVersion
	var id, pluginID, permissions string
	err := row.Scan(
		&id, &pluginID, &v.Version, &v.Changelog,
		&permissions, &v.ManifestJSON, &v.FrontendURL,
		&v.BackendEntry, &v.MinPlatformVer, &v.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	v.ID = kernel.PluginVersionID(id)
	v.PluginID = kernel.PluginID(pluginID)
	v.Permissions = []byte(permissions)
	return &v, nil
}

func scanVersionRow(rows *sql.Rows) (*plugin.PluginVersion, error) {
	var v plugin.PluginVersion
	var id, pluginID, permissions string
	err := rows.Scan(
		&id, &pluginID, &v.Version, &v.Changelog,
		&permissions, &v.ManifestJSON, &v.FrontendURL,
		&v.BackendEntry, &v.MinPlatformVer, &v.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	v.ID = kernel.PluginVersionID(id)
	v.PluginID = kernel.PluginID(pluginID)
	v.Permissions = []byte(permissions)
	return &v, nil
}

func scanInstallation(row rowScanner) (*plugin.PluginInstallation, error) {
	var i plugin.PluginInstallation
	var id, tenantID, pluginID, versionID, status, settings string
	err := row.Scan(
		&id, &tenantID, &pluginID, &versionID,
		&status, &settings, &i.InstalledAt, &i.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	i.ID = kernel.InstallationID(id)
	i.TenantID = kernel.TenantID(tenantID)
	i.PluginID = kernel.PluginID(pluginID)
	i.VersionID = kernel.PluginVersionID(versionID)
	i.Status = plugin.InstallationStatus(status)
	i.Settings = []byte(settings)
	return &i, nil
}

func scanInstallationRow(rows *sql.Rows) (*plugin.PluginInstallation, error) {
	var i plugin.PluginInstallation
	var id, tenantID, pluginID, versionID, status, settings string
	err := rows.Scan(
		&id, &tenantID, &pluginID, &versionID,
		&status, &settings, &i.InstalledAt, &i.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	i.ID = kernel.InstallationID(id)
	i.TenantID = kernel.TenantID(tenantID)
	i.PluginID = kernel.PluginID(pluginID)
	i.VersionID = kernel.PluginVersionID(versionID)
	i.Status = plugin.InstallationStatus(status)
	i.Settings = []byte(settings)
	return &i, nil
}
