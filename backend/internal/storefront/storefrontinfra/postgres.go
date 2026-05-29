package storefrontinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/storefront"
)

// PostgresPageRepository implements storefront.PageRepository using PostgreSQL.
type PostgresPageRepository struct {
	db *sql.DB
}

// NewPostgresPageRepository creates a new PostgresPageRepository.
func NewPostgresPageRepository(db *sql.DB) *PostgresPageRepository {
	return &PostgresPageRepository{db: db}
}

// Create inserts a new page row.
func (r *PostgresPageRepository) Create(ctx context.Context, page *storefront.Page) error {
	metaJSON, err := json.Marshal(page.Meta)
	if err != nil {
		return fmt.Errorf("marshal meta: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO storefront_pages
			(id, tenant_id, slug, title, html, css, meta, status, version, created_by, published_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		string(page.ID), string(page.TenantID), page.Slug, page.Title,
		page.HTML, page.CSS, metaJSON, string(page.Status),
		page.Version, page.CreatedBy, page.PublishedAt,
		page.CreatedAt, page.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert page: %w", err)
	}
	return nil
}

// GetByID retrieves a page by primary key, scoped to the tenant.
func (r *PostgresPageRepository) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.PageID) (*storefront.Page, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, slug, title, html, css, meta, status, version, created_by, published_at, created_at, updated_at
		FROM storefront_pages
		WHERE tenant_id = $1 AND id = $2`,
		string(tenantID), string(id),
	)
	return scanPage(row)
}

// GetBySlug retrieves a page by its slug, scoped to the tenant.
func (r *PostgresPageRepository) GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*storefront.Page, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, slug, title, html, css, meta, status, version, created_by, published_at, created_at, updated_at
		FROM storefront_pages
		WHERE tenant_id = $1 AND slug = $2`,
		string(tenantID), slug,
	)
	return scanPage(row)
}

// GetPublished retrieves a published page by slug for public serving.
func (r *PostgresPageRepository) GetPublished(ctx context.Context, tenantID kernel.TenantID, slug string) (*storefront.Page, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, slug, title, html, css, meta, status, version, created_by, published_at, created_at, updated_at
		FROM storefront_pages
		WHERE tenant_id = $1 AND slug = $2 AND status = 'published'`,
		string(tenantID), slug,
	)
	return scanPage(row)
}

// Update persists mutations to an existing page row.
func (r *PostgresPageRepository) Update(ctx context.Context, page *storefront.Page) error {
	metaJSON, err := json.Marshal(page.Meta)
	if err != nil {
		return fmt.Errorf("marshal meta: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE storefront_pages
		SET slug=$3, title=$4, html=$5, css=$6, meta=$7, status=$8, version=$9, published_at=$10, updated_at=$11
		WHERE tenant_id=$1 AND id=$2`,
		string(page.TenantID), string(page.ID),
		page.Slug, page.Title, page.HTML, page.CSS, metaJSON,
		string(page.Status), page.Version, page.PublishedAt, page.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update page: %w", err)
	}
	return nil
}

// ListByStatus returns pages for a tenant filtered by status, with pagination.
func (r *PostgresPageRepository) ListByStatus(ctx context.Context, tenantID kernel.TenantID, status storefront.PageStatus, p kernel.Pagination) (kernel.PaginatedResult[storefront.Page], error) {
	var total int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM storefront_pages WHERE tenant_id=$1 AND status=$2`,
		string(tenantID), string(status),
	).Scan(&total)
	if err != nil {
		return kernel.PaginatedResult[storefront.Page]{}, fmt.Errorf("count pages by status: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, slug, title, html, css, meta, status, version, created_by, published_at, created_at, updated_at
		FROM storefront_pages
		WHERE tenant_id=$1 AND status=$2
		ORDER BY updated_at DESC
		LIMIT $3 OFFSET $4`,
		string(tenantID), string(status), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.PaginatedResult[storefront.Page]{}, fmt.Errorf("list pages by status: %w", err)
	}
	defer rows.Close()

	pages, err := scanPages(rows)
	if err != nil {
		return kernel.PaginatedResult[storefront.Page]{}, err
	}
	return kernel.NewPaginatedResult(pages, total, p), nil
}

// List returns all pages for a tenant with pagination.
func (r *PostgresPageRepository) List(ctx context.Context, tenantID kernel.TenantID, p kernel.Pagination) (kernel.PaginatedResult[storefront.Page], error) {
	var total int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM storefront_pages WHERE tenant_id=$1`,
		string(tenantID),
	).Scan(&total)
	if err != nil {
		return kernel.PaginatedResult[storefront.Page]{}, fmt.Errorf("count pages: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, slug, title, html, css, meta, status, version, created_by, published_at, created_at, updated_at
		FROM storefront_pages
		WHERE tenant_id=$1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.PaginatedResult[storefront.Page]{}, fmt.Errorf("list pages: %w", err)
	}
	defer rows.Close()

	pages, err := scanPages(rows)
	if err != nil {
		return kernel.PaginatedResult[storefront.Page]{}, err
	}
	return kernel.NewPaginatedResult(pages, total, p), nil
}

// scanPage scans a single page row.
func scanPage(row *sql.Row) (*storefront.Page, error) {
	var p storefront.Page
	var idStr, tenantStr, statusStr string
	var metaJSON []byte
	var publishedAt sql.NullTime

	err := row.Scan(
		&idStr, &tenantStr, &p.Slug, &p.Title,
		&p.HTML, &p.CSS, &metaJSON, &statusStr,
		&p.Version, &p.CreatedBy, &publishedAt,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, storefront.ErrPageNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan page: %w", err)
	}

	p.ID = kernel.PageID(idStr)
	p.TenantID = kernel.TenantID(tenantStr)
	p.Status = storefront.PageStatus(statusStr)
	if publishedAt.Valid {
		t := publishedAt.Time
		p.PublishedAt = &t
	}
	if err := json.Unmarshal(metaJSON, &p.Meta); err != nil {
		return nil, fmt.Errorf("unmarshal meta: %w", err)
	}
	return &p, nil
}

// scanPages scans multiple page rows.
func scanPages(rows *sql.Rows) ([]storefront.Page, error) {
	var pages []storefront.Page
	for rows.Next() {
		var p storefront.Page
		var idStr, tenantStr, statusStr string
		var metaJSON []byte
		var publishedAt sql.NullTime

		err := rows.Scan(
			&idStr, &tenantStr, &p.Slug, &p.Title,
			&p.HTML, &p.CSS, &metaJSON, &statusStr,
			&p.Version, &p.CreatedBy, &publishedAt,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan page row: %w", err)
		}

		p.ID = kernel.PageID(idStr)
		p.TenantID = kernel.TenantID(tenantStr)
		p.Status = storefront.PageStatus(statusStr)
		if publishedAt.Valid {
			t := publishedAt.Time
			p.PublishedAt = &t
		}
		if err := json.Unmarshal(metaJSON, &p.Meta); err != nil {
			return nil, fmt.Errorf("unmarshal meta: %w", err)
		}
		pages = append(pages, p)
	}
	return pages, rows.Err()
}

// PostgresPageVersionRepository implements storefront.PageVersionRepository using PostgreSQL.
type PostgresPageVersionRepository struct {
	db *sql.DB
}

// NewPostgresPageVersionRepository creates a new PostgresPageVersionRepository.
func NewPostgresPageVersionRepository(db *sql.DB) *PostgresPageVersionRepository {
	return &PostgresPageVersionRepository{db: db}
}

// Create inserts a new page version snapshot (append-only).
func (r *PostgresPageVersionRepository) Create(ctx context.Context, v *storefront.PageVersion) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO storefront_page_versions
			(id, page_id, tenant_id, version, html, css, edited_by, comment, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		string(v.ID), string(v.PageID), string(v.TenantID),
		v.Version, v.HTML, v.CSS, v.EditedBy, v.Comment, v.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert page version: %w", err)
	}
	return nil
}

// GetByVersion retrieves a specific version of a page.
func (r *PostgresPageVersionRepository) GetByVersion(ctx context.Context, tenantID kernel.TenantID, pageID kernel.PageID, version int) (*storefront.PageVersion, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, page_id, tenant_id, version, html, css, edited_by, comment, created_at
		FROM storefront_page_versions
		WHERE tenant_id=$1 AND page_id=$2 AND version=$3`,
		string(tenantID), string(pageID), version,
	)
	return scanPageVersion(row)
}

// ListByPage returns all version snapshots for a page, newest first.
func (r *PostgresPageVersionRepository) ListByPage(ctx context.Context, tenantID kernel.TenantID, pageID kernel.PageID) ([]storefront.PageVersion, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, page_id, tenant_id, version, html, css, edited_by, comment, created_at
		FROM storefront_page_versions
		WHERE tenant_id=$1 AND page_id=$2
		ORDER BY version DESC`,
		string(tenantID), string(pageID),
	)
	if err != nil {
		return nil, fmt.Errorf("list page versions: %w", err)
	}
	defer rows.Close()

	var versions []storefront.PageVersion
	for rows.Next() {
		v, err := scanPageVersionRow(rows)
		if err != nil {
			return nil, err
		}
		versions = append(versions, *v)
	}
	return versions, rows.Err()
}

func scanPageVersion(row *sql.Row) (*storefront.PageVersion, error) {
	var v storefront.PageVersion
	var idStr, pageIDStr, tenantStr string
	err := row.Scan(
		&idStr, &pageIDStr, &tenantStr,
		&v.Version, &v.HTML, &v.CSS,
		&v.EditedBy, &v.Comment, &v.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, storefront.ErrVersionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan page version: %w", err)
	}
	v.ID = kernel.PageVersionID(idStr)
	v.PageID = kernel.PageID(pageIDStr)
	v.TenantID = kernel.TenantID(tenantStr)
	return &v, nil
}

func scanPageVersionRow(rows *sql.Rows) (*storefront.PageVersion, error) {
	var v storefront.PageVersion
	var idStr, pageIDStr, tenantStr string
	err := rows.Scan(
		&idStr, &pageIDStr, &tenantStr,
		&v.Version, &v.HTML, &v.CSS,
		&v.EditedBy, &v.Comment, &v.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan page version row: %w", err)
	}
	v.ID = kernel.PageVersionID(idStr)
	v.PageID = kernel.PageID(pageIDStr)
	v.TenantID = kernel.TenantID(tenantStr)
	return &v, nil
}

// Ensure time is imported (used in sql.NullTime).
var _ = time.Time{}
