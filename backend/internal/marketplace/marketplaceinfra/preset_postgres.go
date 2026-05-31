package marketplaceinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/marketplace"
)

// ─── Preset Repository ──────────────────────────────────────────────────────

// PostgresPresetRepo implements marketplace.PresetRepository.
type PostgresPresetRepo struct{ db *sqlx.DB }

// NewPostgresPresetRepo creates a new PostgresPresetRepo.
func NewPostgresPresetRepo(db *sqlx.DB) *PostgresPresetRepo {
	return &PostgresPresetRepo{db: db}
}

// dbPreset is the sqlx-scannable row for a preset.
type dbPreset struct {
	ID            string    `db:"id"`
	TenantID      string    `db:"tenant_id"`
	Name          string    `db:"name"`
	Slug          string    `db:"slug"`
	Description   string    `db:"description"`
	Version       string    `db:"version"`
	Image         string    `db:"image"`
	FrontendPort  int       `db:"frontend_port"`
	SystemPrompt  string    `db:"system_prompt"`
	ToolsManifest []byte    `db:"tools_manifest"`
	Status        string    `db:"status"`
	Visibility    string    `db:"visibility"`
	Icon          string    `db:"icon"`
	Tags          []byte    `db:"tags"`
	InstallCount  int       `db:"install_count"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func fromDBPreset(row dbPreset) marketplace.Preset {
	var tags []string
	if len(row.Tags) > 0 {
		_ = json.Unmarshal(row.Tags, &tags)
	}
	if tags == nil {
		tags = []string{}
	}
	toolsManifest := json.RawMessage(row.ToolsManifest)
	if len(toolsManifest) == 0 {
		toolsManifest = json.RawMessage("[]")
	}
	return marketplace.Preset{
		ID:            kernel.PresetID(row.ID),
		TenantID:      kernel.TenantID(row.TenantID),
		Name:          row.Name,
		Slug:          row.Slug,
		Description:   row.Description,
		Version:       row.Version,
		Image:         row.Image,
		FrontendPort:  row.FrontendPort,
		SystemPrompt:  row.SystemPrompt,
		ToolsManifest: toolsManifest,
		Status:        marketplace.PresetStatus(row.Status),
		Visibility:    marketplace.PresetVisibility(row.Visibility),
		Icon:          row.Icon,
		Tags:          tags,
		InstallCount:  row.InstallCount,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

func (r *PostgresPresetRepo) Create(ctx context.Context, p marketplace.Preset) (marketplace.Preset, error) {
	tagsJSON, err := json.Marshal(p.Tags)
	if err != nil {
		return marketplace.Preset{}, errx.Wrap(err, "marshal preset tags", errx.TypeInternal)
	}
	toolsJSON := p.ToolsManifest
	if len(toolsJSON) == 0 {
		toolsJSON = json.RawMessage("[]")
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO presets
			(id, tenant_id, name, slug, description, version, image, frontend_port,
			 system_prompt, tools_manifest, status, visibility, icon, tags, install_count,
			 created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`,
		string(p.ID), string(p.TenantID), p.Name, p.Slug, p.Description,
		p.Version, p.Image, p.FrontendPort, p.SystemPrompt,
		[]byte(toolsJSON), string(p.Status), string(p.Visibility),
		p.Icon, tagsJSON, p.InstallCount, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return marketplace.Preset{}, marketplace.ErrPresetSlugTaken
		}
		return marketplace.Preset{}, errx.Wrap(err, "create preset", errx.TypeInternal)
	}
	return p, nil
}

func (r *PostgresPresetRepo) GetByID(ctx context.Context, id kernel.PresetID) (marketplace.Preset, error) {
	var row dbPreset
	err := r.db.GetContext(ctx, &row, `
		SELECT id, tenant_id, name, slug, description, version, image, frontend_port,
		       system_prompt, tools_manifest, status, visibility, icon, tags, install_count,
		       created_at, updated_at
		FROM presets WHERE id=$1`,
		string(id),
	)
	if err == sql.ErrNoRows {
		return marketplace.Preset{}, marketplace.ErrPresetNotFound
	}
	if err != nil {
		return marketplace.Preset{}, errx.Wrap(err, "get preset", errx.TypeInternal)
	}
	return fromDBPreset(row), nil
}

func (r *PostgresPresetRepo) GetBySlug(ctx context.Context, slug string) (marketplace.Preset, error) {
	var row dbPreset
	err := r.db.GetContext(ctx, &row, `
		SELECT id, tenant_id, name, slug, description, version, image, frontend_port,
		       system_prompt, tools_manifest, status, visibility, icon, tags, install_count,
		       created_at, updated_at
		FROM presets WHERE slug=$1`,
		slug,
	)
	if err == sql.ErrNoRows {
		return marketplace.Preset{}, marketplace.ErrPresetNotFound
	}
	if err != nil {
		return marketplace.Preset{}, errx.Wrap(err, "get preset by slug", errx.TypeInternal)
	}
	return fromDBPreset(row), nil
}

func (r *PostgresPresetRepo) Update(ctx context.Context, p marketplace.Preset) (marketplace.Preset, error) {
	tagsJSON, err := json.Marshal(p.Tags)
	if err != nil {
		return marketplace.Preset{}, errx.Wrap(err, "marshal preset tags", errx.TypeInternal)
	}
	toolsJSON := p.ToolsManifest
	if len(toolsJSON) == 0 {
		toolsJSON = json.RawMessage("[]")
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE presets
		SET name=$2, description=$3, version=$4, image=$5, frontend_port=$6,
		    system_prompt=$7, tools_manifest=$8, status=$9, visibility=$10,
		    icon=$11, tags=$12, updated_at=$13
		WHERE id=$1`,
		string(p.ID), p.Name, p.Description, p.Version, p.Image, p.FrontendPort,
		p.SystemPrompt, []byte(toolsJSON), string(p.Status), string(p.Visibility),
		p.Icon, tagsJSON, p.UpdatedAt,
	)
	if err != nil {
		return marketplace.Preset{}, errx.Wrap(err, "update preset", errx.TypeInternal)
	}
	return p, nil
}

func (r *PostgresPresetRepo) Delete(ctx context.Context, id kernel.PresetID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM presets WHERE id=$1`, string(id))
	if err != nil {
		return errx.Wrap(err, "delete preset", errx.TypeInternal)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return marketplace.ErrPresetNotFound
	}
	return nil
}

func (r *PostgresPresetRepo) List(ctx context.Context, opts marketplace.PresetListOptions) (kernel.Paginated[marketplace.Preset], error) {
	// Build dynamic WHERE clause
	where := "WHERE 1=1"
	args := []interface{}{}
	idx := 1

	if opts.Status != nil {
		where += " AND status=$" + itoa(idx)
		args = append(args, string(*opts.Status))
		idx++
	}
	if opts.Visibility != nil {
		where += " AND visibility=$" + itoa(idx)
		args = append(args, string(*opts.Visibility))
		idx++
	}
	if opts.Search != "" {
		where += " AND (name ILIKE $" + itoa(idx) + " OR description ILIKE $" + itoa(idx) + ")"
		args = append(args, "%"+opts.Search+"%")
		idx++
	}

	var total int
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM presets "+where, countArgs...).Scan(&total); err != nil {
		return kernel.Paginated[marketplace.Preset]{}, errx.Wrap(err, "count presets", errx.TypeInternal)
	}

	query := `SELECT id, tenant_id, name, slug, description, version, image, frontend_port,
	                 system_prompt, tools_manifest, status, visibility, icon, tags, install_count,
	                 created_at, updated_at
	          FROM presets ` + where + ` ORDER BY install_count DESC, created_at DESC LIMIT $` + itoa(idx) + ` OFFSET $` + itoa(idx+1)
	args = append(args, opts.Limit(), opts.Offset())

	var rows []dbPreset
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return kernel.Paginated[marketplace.Preset]{}, errx.Wrap(err, "list presets", errx.TypeInternal)
	}

	items := make([]marketplace.Preset, len(rows))
	for i, row := range rows {
		items[i] = fromDBPreset(row)
	}
	return kernel.NewPaginated(items, opts.Page, opts.PageSize, total), nil
}

func (r *PostgresPresetRepo) ListByTenant(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[marketplace.Preset], error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM presets WHERE tenant_id=$1`, string(tenantID)).Scan(&total); err != nil {
		return kernel.Paginated[marketplace.Preset]{}, errx.Wrap(err, "count tenant presets", errx.TypeInternal)
	}

	var rows []dbPreset
	err := r.db.SelectContext(ctx, &rows, `
		SELECT id, tenant_id, name, slug, description, version, image, frontend_port,
		       system_prompt, tools_manifest, status, visibility, icon, tags, install_count,
		       created_at, updated_at
		FROM presets WHERE tenant_id=$1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		string(tenantID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.Paginated[marketplace.Preset]{}, errx.Wrap(err, "list tenant presets", errx.TypeInternal)
	}

	items := make([]marketplace.Preset, len(rows))
	for i, row := range rows {
		items[i] = fromDBPreset(row)
	}
	return kernel.NewPaginated(items, p.Page, p.PageSize, total), nil
}

// Ensure interface compliance.
var _ marketplace.PresetRepository = (*PostgresPresetRepo)(nil)

// ─── PresetInstall Repository ────────────────────────────────────────────────

// PostgresPresetInstallRepo implements marketplace.PresetInstallRepository.
type PostgresPresetInstallRepo struct{ db *sqlx.DB }

// NewPostgresPresetInstallRepo creates a new PostgresPresetInstallRepo.
func NewPostgresPresetInstallRepo(db *sqlx.DB) *PostgresPresetInstallRepo {
	return &PostgresPresetInstallRepo{db: db}
}

type dbPresetInstall struct {
	ID          string    `db:"id"`
	TenantID    string    `db:"tenant_id"`
	PresetID    string    `db:"preset_id"`
	InstalledAt time.Time `db:"installed_at"`
	Config      []byte    `db:"config"`
}

func fromDBPresetInstall(row dbPresetInstall) marketplace.PresetInstall {
	cfg := json.RawMessage(row.Config)
	if len(cfg) == 0 {
		cfg = json.RawMessage("{}")
	}
	return marketplace.PresetInstall{
		ID:          row.ID,
		TenantID:    kernel.TenantID(row.TenantID),
		PresetID:    kernel.PresetID(row.PresetID),
		InstalledAt: row.InstalledAt,
		Config:      cfg,
	}
}

func (r *PostgresPresetInstallRepo) Install(ctx context.Context, install marketplace.PresetInstall) (marketplace.PresetInstall, error) {
	cfg := install.Config
	if len(cfg) == 0 {
		cfg = json.RawMessage("{}")
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO preset_installs (id, tenant_id, preset_id, installed_at, config)
		VALUES ($1,$2,$3,$4,$5)`,
		install.ID, string(install.TenantID), string(install.PresetID),
		install.InstalledAt, []byte(cfg),
	)
	if err != nil {
		if isUniqueViolation(err) {
			return marketplace.PresetInstall{}, errx.New("preset already installed", errx.TypeConflict)
		}
		return marketplace.PresetInstall{}, errx.Wrap(err, "install preset", errx.TypeInternal)
	}
	return install, nil
}

func (r *PostgresPresetInstallRepo) Uninstall(ctx context.Context, tenantID kernel.TenantID, presetID kernel.PresetID) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM preset_installs WHERE tenant_id=$1 AND preset_id=$2`,
		string(tenantID), string(presetID),
	)
	if err != nil {
		return errx.Wrap(err, "uninstall preset", errx.TypeInternal)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return marketplace.ErrPresetNotInstalled
	}
	return nil
}

func (r *PostgresPresetInstallRepo) ListByTenant(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[marketplace.PresetInstall], error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM preset_installs WHERE tenant_id=$1`, string(tenantID)).Scan(&total); err != nil {
		return kernel.Paginated[marketplace.PresetInstall]{}, errx.Wrap(err, "count installs", errx.TypeInternal)
	}

	var rows []dbPresetInstall
	err := r.db.SelectContext(ctx, &rows, `
		SELECT id, tenant_id, preset_id, installed_at, config
		FROM preset_installs WHERE tenant_id=$1
		ORDER BY installed_at DESC LIMIT $2 OFFSET $3`,
		string(tenantID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.Paginated[marketplace.PresetInstall]{}, errx.Wrap(err, "list installs", errx.TypeInternal)
	}

	items := make([]marketplace.PresetInstall, len(rows))
	for i, row := range rows {
		items[i] = fromDBPresetInstall(row)
	}
	return kernel.NewPaginated(items, p.Page, p.PageSize, total), nil
}

func (r *PostgresPresetInstallRepo) IsInstalled(ctx context.Context, tenantID kernel.TenantID, presetID kernel.PresetID) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM preset_installs WHERE tenant_id=$1 AND preset_id=$2`,
		string(tenantID), string(presetID),
	).Scan(&count)
	if err != nil {
		return false, errx.Wrap(err, "check install", errx.TypeInternal)
	}
	return count > 0, nil
}

// Ensure interface compliance.
var _ marketplace.PresetInstallRepository = (*PostgresPresetInstallRepo)(nil)

// ─── Helpers ─────────────────────────────────────────────────────────────────

// itoa converts an int to its decimal string representation (avoids strconv import).
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for ; i > 0; i /= 10 {
		pos--
		buf[pos] = byte(i%10) + '0'
	}
	return string(buf[pos:])
}

// isUniqueViolation checks if the error is a PostgreSQL unique constraint violation.
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	return len(err.Error()) > 0 && containsStr(err.Error(), "unique") || containsStr(err.Error(), "23505")
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
