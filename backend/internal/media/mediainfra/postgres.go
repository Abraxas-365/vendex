package mediainfra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/media"
)

// PostgresMediaRepository implements media.MediaRepository using PostgreSQL.
type PostgresMediaRepository struct {
	db *sql.DB
}

// NewPostgresMediaRepository creates a new PostgresMediaRepository.
func NewPostgresMediaRepository(db *sql.DB) *PostgresMediaRepository {
	return &PostgresMediaRepository{db: db}
}

// Create inserts a new media metadata row.
func (r *PostgresMediaRepository) Create(ctx context.Context, m *media.Media) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO media
			(id, tenant_id, filename, content_type, size, url, alt, uploaded_by, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		string(m.ID), string(m.TenantID), m.Filename, m.ContentType,
		m.Size, m.URL, m.Alt, m.UploadedBy, m.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert media: %w", err)
	}
	return nil
}

// GetByID retrieves a media record by primary key, scoped to the tenant.
func (r *PostgresMediaRepository) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.MediaID) (*media.Media, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, filename, content_type, size, url, alt, uploaded_by, created_at
		FROM media WHERE tenant_id=$1 AND id=$2`,
		string(tenantID), string(id),
	)
	return scanMedia(row)
}

// Delete removes a media metadata row.
func (r *PostgresMediaRepository) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.MediaID) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM media WHERE tenant_id=$1 AND id=$2`,
		string(tenantID), string(id),
	)
	if err != nil {
		return fmt.Errorf("delete media: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return media.ErrMediaNotFound
	}
	return nil
}

// List returns paginated media records for a tenant.
func (r *PostgresMediaRepository) List(ctx context.Context, tenantID kernel.TenantID, p kernel.Pagination) (kernel.PaginatedResult[media.Media], error) {
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM media WHERE tenant_id=$1`, string(tenantID)).Scan(&total)
	if err != nil {
		return kernel.PaginatedResult[media.Media]{}, fmt.Errorf("count media: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, filename, content_type, size, url, alt, uploaded_by, created_at
		FROM media WHERE tenant_id=$1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.PaginatedResult[media.Media]{}, fmt.Errorf("list media: %w", err)
	}
	defer rows.Close()

	var items []media.Media
	for rows.Next() {
		m, err := scanMediaRow(rows)
		if err != nil {
			return kernel.PaginatedResult[media.Media]{}, err
		}
		items = append(items, *m)
	}
	if err := rows.Err(); err != nil {
		return kernel.PaginatedResult[media.Media]{}, fmt.Errorf("iterate media: %w", err)
	}
	return kernel.NewPaginatedResult(items, total, p), nil
}

func scanMedia(row *sql.Row) (*media.Media, error) {
	var m media.Media
	var idStr, tenantStr string
	err := row.Scan(
		&idStr, &tenantStr, &m.Filename, &m.ContentType,
		&m.Size, &m.URL, &m.Alt, &m.UploadedBy, &m.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, media.ErrMediaNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan media: %w", err)
	}
	m.ID = kernel.MediaID(idStr)
	m.TenantID = kernel.TenantID(tenantStr)
	return &m, nil
}

func scanMediaRow(rows *sql.Rows) (*media.Media, error) {
	var m media.Media
	var idStr, tenantStr string
	err := rows.Scan(
		&idStr, &tenantStr, &m.Filename, &m.ContentType,
		&m.Size, &m.URL, &m.Alt, &m.UploadedBy, &m.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan media row: %w", err)
	}
	m.ID = kernel.MediaID(idStr)
	m.TenantID = kernel.TenantID(tenantStr)
	return &m, nil
}
