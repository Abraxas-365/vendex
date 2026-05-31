package storefrontinfra

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/storefront"
)

// --- Page Repository ---

// PagePostgresRepo implements storefront.PageRepository.
type PagePostgresRepo struct {
	db *sqlx.DB
}

// NewPagePostgresRepo creates a new PostgreSQL-backed page repository.
func NewPagePostgresRepo(db *sqlx.DB) *PagePostgresRepo {
	return &PagePostgresRepo{db: db}
}

func (r *PagePostgresRepo) Create(ctx context.Context, page *storefront.Page) error {
	metaJSON, err := json.Marshal(page.Meta)
	if err != nil {
		return errx.Wrap(err, "marshaling page meta", errx.TypeInternal)
	}

	sections := page.Sections
	if sections == nil {
		sections = []storefront.Section{}
	}
	sectionsJSON, err := json.Marshal(sections)
	if err != nil {
		return errx.Wrap(err, "marshaling page sections", errx.TypeInternal)
	}

	contentType := page.ContentType
	if contentType == "" {
		contentType = storefront.ContentTypeHTML
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO pages (id, tenant_id, slug, title, html, css, meta, content_type, sections, status, version, created_by, published_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		string(page.ID), string(page.TenantID), page.Slug, page.Title,
		page.HTML, page.CSS, string(metaJSON), string(contentType), string(sectionsJSON),
		string(page.Status), page.Version, page.CreatedBy, page.PublishedAt, page.CreatedAt, page.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting page", errx.TypeInternal)
	}
	return nil
}

func (r *PagePostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.PageID) (*storefront.Page, error) {
	return r.queryOnePage(ctx, `
		SELECT id, tenant_id, slug, title, html, css, meta, content_type, sections, status, version, created_by, published_at, created_at, updated_at
		FROM pages WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
}

func (r *PagePostgresRepo) GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*storefront.Page, error) {
	return r.queryOnePage(ctx, `
		SELECT id, tenant_id, slug, title, html, css, meta, content_type, sections, status, version, created_by, published_at, created_at, updated_at
		FROM pages WHERE slug = $1 AND tenant_id = $2`,
		slug, string(tenantID),
	)
}

func (r *PagePostgresRepo) GetPublished(ctx context.Context, tenantID kernel.TenantID, slug string) (*storefront.Page, error) {
	return r.queryOnePage(ctx, `
		SELECT id, tenant_id, slug, title, html, css, meta, content_type, sections, status, version, created_by, published_at, created_at, updated_at
		FROM pages WHERE slug = $1 AND tenant_id = $2 AND status = 'published'`,
		slug, string(tenantID),
	)
}

func (r *PagePostgresRepo) queryOnePage(ctx context.Context, query string, args ...any) (*storefront.Page, error) {
	var page storefront.Page
	var id, tenantID, status, contentType string
	var metaJSON, sectionsJSON string

	err := r.db.QueryRowContext(ctx, query, args...).
		Scan(&id, &tenantID, &page.Slug, &page.Title, &page.HTML, &page.CSS,
			&metaJSON, &contentType, &sectionsJSON, &status, &page.Version, &page.CreatedBy,
			&page.PublishedAt, &page.CreatedAt, &page.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, storefront.ErrPageNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "querying page", errx.TypeInternal)
	}

	page.ID = kernel.PageID(id)
	page.TenantID = kernel.TenantID(tenantID)
	page.Status = storefront.PageStatus(status)
	page.ContentType = storefront.ContentType(contentType)
	_ = json.Unmarshal([]byte(metaJSON), &page.Meta)
	_ = json.Unmarshal([]byte(sectionsJSON), &page.Sections)

	return &page, nil
}

func (r *PagePostgresRepo) Update(ctx context.Context, page *storefront.Page) error {
	metaJSON, err := json.Marshal(page.Meta)
	if err != nil {
		return errx.Wrap(err, "marshaling page meta", errx.TypeInternal)
	}

	sections := page.Sections
	if sections == nil {
		sections = []storefront.Section{}
	}
	sectionsJSON, err := json.Marshal(sections)
	if err != nil {
		return errx.Wrap(err, "marshaling page sections", errx.TypeInternal)
	}

	contentType := page.ContentType
	if contentType == "" {
		contentType = storefront.ContentTypeHTML
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE pages
		SET slug=$1, title=$2, html=$3, css=$4, meta=$5, content_type=$6, sections=$7, status=$8, version=$9, published_at=$10, updated_at=$11
		WHERE id=$12 AND tenant_id=$13`,
		page.Slug, page.Title, page.HTML, page.CSS, string(metaJSON),
		string(contentType), string(sectionsJSON),
		string(page.Status), page.Version, page.PublishedAt, page.UpdatedAt,
		string(page.ID), string(page.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating page", errx.TypeInternal)
	}
	return nil
}

func (r *PagePostgresRepo) ListByStatus(ctx context.Context, tenantID kernel.TenantID, status storefront.PageStatus, p kernel.PaginationOptions) (kernel.Paginated[storefront.Page], error) {
	var zero kernel.Paginated[storefront.Page]

	var total int
	if err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM pages WHERE tenant_id = $1 AND status = $2",
		string(tenantID), string(status),
	).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting pages by status", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, slug, title, html, css, meta, content_type, sections, status, version, created_by, published_at, created_at, updated_at
		FROM pages WHERE tenant_id = $1 AND status = $2
		ORDER BY updated_at DESC
		LIMIT $3 OFFSET $4`,
		string(tenantID), string(status), p.Limit(), p.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "querying pages by status", errx.TypeInternal)
	}
	defer rows.Close()

	pages, err := scanPageRows(rows)
	if err != nil {
		return zero, err
	}

	return kernel.NewPaginated(pages, p.Page, p.PageSize, total), nil
}

func (r *PagePostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[storefront.Page], error) {
	var zero kernel.Paginated[storefront.Page]

	var total int
	if err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM pages WHERE tenant_id = $1", string(tenantID),
	).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting pages", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, slug, title, html, css, meta, content_type, sections, status, version, created_by, published_at, created_at, updated_at
		FROM pages WHERE tenant_id = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "querying pages", errx.TypeInternal)
	}
	defer rows.Close()

	pages, err := scanPageRows(rows)
	if err != nil {
		return zero, err
	}

	return kernel.NewPaginated(pages, p.Page, p.PageSize, total), nil
}

func scanPageRows(rows *sql.Rows) ([]storefront.Page, error) {
	var pages []storefront.Page
	for rows.Next() {
		var page storefront.Page
		var id, tenantID, status, contentType, metaJSON, sectionsJSON string
		err := rows.Scan(&id, &tenantID, &page.Slug, &page.Title, &page.HTML, &page.CSS,
			&metaJSON, &contentType, &sectionsJSON, &status, &page.Version, &page.CreatedBy,
			&page.PublishedAt, &page.CreatedAt, &page.UpdatedAt)
		if err != nil {
			return nil, errx.Wrap(err, "scanning page", errx.TypeInternal)
		}
		page.ID = kernel.PageID(id)
		page.TenantID = kernel.TenantID(tenantID)
		page.Status = storefront.PageStatus(status)
		page.ContentType = storefront.ContentType(contentType)
		_ = json.Unmarshal([]byte(metaJSON), &page.Meta)
		_ = json.Unmarshal([]byte(sectionsJSON), &page.Sections)
		pages = append(pages, page)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating pages", errx.TypeInternal)
	}
	return pages, nil
}

// Ensure interface compliance.
var _ storefront.PageRepository = (*PagePostgresRepo)(nil)

// --- PageVersion Repository ---

// PageVersionPostgresRepo implements storefront.PageVersionRepository.
type PageVersionPostgresRepo struct {
	db *sqlx.DB
}

// NewPageVersionPostgresRepo creates a new PostgreSQL-backed page version repository.
func NewPageVersionPostgresRepo(db *sqlx.DB) *PageVersionPostgresRepo {
	return &PageVersionPostgresRepo{db: db}
}

func (r *PageVersionPostgresRepo) Create(ctx context.Context, v *storefront.PageVersion) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO page_versions (id, page_id, tenant_id, version, html, css, edited_by, comment, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		string(v.ID), string(v.PageID), string(v.TenantID),
		v.Version, v.HTML, v.CSS, v.EditedBy, v.Comment, v.CreatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting page version", errx.TypeInternal)
	}
	return nil
}

func (r *PageVersionPostgresRepo) GetByVersion(ctx context.Context, tenantID kernel.TenantID, pageID kernel.PageID, version int) (*storefront.PageVersion, error) {
	var v storefront.PageVersion
	var id, pID, tID string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, page_id, tenant_id, version, html, css, edited_by, comment, created_at
		FROM page_versions WHERE page_id = $1 AND tenant_id = $2 AND version = $3`,
		string(pageID), string(tenantID), version,
	).Scan(&id, &pID, &tID, &v.Version, &v.HTML, &v.CSS, &v.EditedBy, &v.Comment, &v.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, storefront.ErrVersionNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "getting page version", errx.TypeInternal)
	}

	v.ID = kernel.PageVersionID(id)
	v.PageID = kernel.PageID(pID)
	v.TenantID = kernel.TenantID(tID)
	return &v, nil
}

func (r *PageVersionPostgresRepo) ListByPage(ctx context.Context, tenantID kernel.TenantID, pageID kernel.PageID) ([]storefront.PageVersion, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, page_id, tenant_id, version, html, css, edited_by, comment, created_at
		FROM page_versions WHERE page_id = $1 AND tenant_id = $2
		ORDER BY version DESC`,
		string(pageID), string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying page versions", errx.TypeInternal)
	}
	defer rows.Close()

	var versions []storefront.PageVersion
	for rows.Next() {
		var v storefront.PageVersion
		var id, pID, tID string
		err := rows.Scan(&id, &pID, &tID, &v.Version, &v.HTML, &v.CSS, &v.EditedBy, &v.Comment, &v.CreatedAt)
		if err != nil {
			return nil, errx.Wrap(err, "scanning page version", errx.TypeInternal)
		}
		v.ID = kernel.PageVersionID(id)
		v.PageID = kernel.PageID(pID)
		v.TenantID = kernel.TenantID(tID)
		versions = append(versions, v)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating page versions", errx.TypeInternal)
	}

	return versions, nil
}

// Ensure interface compliance.
var _ storefront.PageVersionRepository = (*PageVersionPostgresRepo)(nil)

// --- BlockType Repository ---

// BlockTypePostgresRepo implements storefront.BlockTypeRepository.
type BlockTypePostgresRepo struct {
	db *sqlx.DB
}

// NewBlockTypePostgresRepo creates a new PostgreSQL-backed block type repository.
func NewBlockTypePostgresRepo(db *sqlx.DB) *BlockTypePostgresRepo {
	return &BlockTypePostgresRepo{db: db}
}

func (r *BlockTypePostgresRepo) Create(ctx context.Context, bt *storefront.BlockType) error {
	schemaJSON, err := json.Marshal(bt.Schema)
	if err != nil {
		return errx.Wrap(err, "marshaling block type schema", errx.TypeInternal)
	}
	defaultSettingsJSON, err := json.Marshal(bt.DefaultSettings)
	if err != nil {
		return errx.Wrap(err, "marshaling block type default settings", errx.TypeInternal)
	}

	var pluginID *string
	if bt.PluginID != nil {
		s := bt.PluginID.String()
		pluginID = &s
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO block_types (id, name, display_name, category, schema, default_settings, icon, plugin_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		string(bt.ID), bt.Name, bt.DisplayName, bt.Category,
		string(schemaJSON), string(defaultSettingsJSON),
		bt.Icon, pluginID, bt.CreatedAt, bt.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting block type", errx.TypeInternal)
	}
	return nil
}

func (r *BlockTypePostgresRepo) GetByID(ctx context.Context, id kernel.BlockTypeID) (*storefront.BlockType, error) {
	return r.queryOneBlockType(ctx,
		`SELECT id, name, display_name, category, schema, default_settings, icon, plugin_id, created_at, updated_at
		 FROM block_types WHERE id = $1`,
		string(id),
	)
}

func (r *BlockTypePostgresRepo) GetByName(ctx context.Context, name string) (*storefront.BlockType, error) {
	return r.queryOneBlockType(ctx,
		`SELECT id, name, display_name, category, schema, default_settings, icon, plugin_id, created_at, updated_at
		 FROM block_types WHERE name = $1`,
		name,
	)
}

func (r *BlockTypePostgresRepo) queryOneBlockType(ctx context.Context, query string, args ...any) (*storefront.BlockType, error) {
	var bt storefront.BlockType
	var id string
	var schemaJSON, defaultSettingsJSON string
	var pluginID *string

	err := r.db.QueryRowContext(ctx, query, args...).
		Scan(&id, &bt.Name, &bt.DisplayName, &bt.Category,
			&schemaJSON, &defaultSettingsJSON, &bt.Icon,
			&pluginID, &bt.CreatedAt, &bt.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, storefront.ErrBlockTypeNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "querying block type", errx.TypeInternal)
	}

	bt.ID = kernel.BlockTypeID(id)
	_ = json.Unmarshal([]byte(schemaJSON), &bt.Schema)
	_ = json.Unmarshal([]byte(defaultSettingsJSON), &bt.DefaultSettings)
	if pluginID != nil {
		pid := kernel.PluginID(*pluginID)
		bt.PluginID = &pid
	}

	return &bt, nil
}

func (r *BlockTypePostgresRepo) List(ctx context.Context, category string) ([]storefront.BlockType, error) {
	var rows *sql.Rows
	var err error

	if category != "" {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, name, display_name, category, schema, default_settings, icon, plugin_id, created_at, updated_at
			FROM block_types WHERE category = $1
			ORDER BY category, name`,
			category,
		)
	} else {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, name, display_name, category, schema, default_settings, icon, plugin_id, created_at, updated_at
			FROM block_types
			ORDER BY category, name`,
		)
	}
	if err != nil {
		return nil, errx.Wrap(err, "querying block types", errx.TypeInternal)
	}
	defer rows.Close()

	var blockTypes []storefront.BlockType
	for rows.Next() {
		var bt storefront.BlockType
		var id string
		var schemaJSON, defaultSettingsJSON string
		var pluginID *string

		err := rows.Scan(&id, &bt.Name, &bt.DisplayName, &bt.Category,
			&schemaJSON, &defaultSettingsJSON, &bt.Icon,
			&pluginID, &bt.CreatedAt, &bt.UpdatedAt)
		if err != nil {
			return nil, errx.Wrap(err, "scanning block type", errx.TypeInternal)
		}

		bt.ID = kernel.BlockTypeID(id)
		_ = json.Unmarshal([]byte(schemaJSON), &bt.Schema)
		_ = json.Unmarshal([]byte(defaultSettingsJSON), &bt.DefaultSettings)
		if pluginID != nil {
			pid := kernel.PluginID(*pluginID)
			bt.PluginID = &pid
		}
		blockTypes = append(blockTypes, bt)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating block types", errx.TypeInternal)
	}

	if blockTypes == nil {
		blockTypes = []storefront.BlockType{}
	}
	return blockTypes, nil
}

func (r *BlockTypePostgresRepo) Update(ctx context.Context, bt *storefront.BlockType) error {
	schemaJSON, err := json.Marshal(bt.Schema)
	if err != nil {
		return errx.Wrap(err, "marshaling block type schema", errx.TypeInternal)
	}
	defaultSettingsJSON, err := json.Marshal(bt.DefaultSettings)
	if err != nil {
		return errx.Wrap(err, "marshaling block type default settings", errx.TypeInternal)
	}

	var pluginID *string
	if bt.PluginID != nil {
		s := bt.PluginID.String()
		pluginID = &s
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE block_types
		SET name=$1, display_name=$2, category=$3, schema=$4, default_settings=$5, icon=$6, plugin_id=$7, updated_at=$8
		WHERE id=$9`,
		bt.Name, bt.DisplayName, bt.Category,
		string(schemaJSON), string(defaultSettingsJSON),
		bt.Icon, pluginID, bt.UpdatedAt,
		string(bt.ID),
	)
	if err != nil {
		return errx.Wrap(err, "updating block type", errx.TypeInternal)
	}
	return nil
}

func (r *BlockTypePostgresRepo) Delete(ctx context.Context, id kernel.BlockTypeID) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM block_types WHERE id = $1`,
		string(id),
	)
	if err != nil {
		return errx.Wrap(err, "deleting block type", errx.TypeInternal)
	}
	return nil
}

// Ensure interface compliance.
var _ storefront.BlockTypeRepository = (*BlockTypePostgresRepo)(nil)
