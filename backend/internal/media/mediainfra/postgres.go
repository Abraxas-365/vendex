package mediainfra

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/media"
)

// PostgresMediaRepository implements media.MediaRepository using sqlx.
type PostgresMediaRepository struct {
	db *sqlx.DB
}

// NewPostgresMediaRepository creates a new PostgresMediaRepository.
func NewPostgresMediaRepository(db *sqlx.DB) *PostgresMediaRepository {
	return &PostgresMediaRepository{db: db}
}

// dbMedia is the sqlx-scannable row for a media record.
type dbMedia struct {
	ID          string `db:"id"`
	TenantID    string `db:"tenant_id"`
	Filename    string `db:"filename"`
	ContentType string `db:"content_type"`
	Size        int64  `db:"size"`
	URL         string `db:"url"`
	Alt         string `db:"alt"`
	UploadedBy  string `db:"uploaded_by"`
	CreatedAt   string `db:"created_at"`
}

func toMedia(row dbMedia) media.Media {
	return media.Media{
		ID:          kernel.MediaID(row.ID),
		TenantID:    kernel.TenantID(row.TenantID),
		Filename:    row.Filename,
		ContentType: row.ContentType,
		Size:        row.Size,
		URL:         row.URL,
		Alt:         row.Alt,
		UploadedBy:  row.UploadedBy,
	}
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
		return errx.Wrap(err, "insert media", errx.TypeInternal)
	}
	return nil
}

// GetByID retrieves a media record by primary key, scoped to the tenant.
func (r *PostgresMediaRepository) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.MediaID) (*media.Media, error) {
	var m media.Media
	err := r.db.GetContext(ctx, &m, `
		SELECT id, tenant_id, filename, content_type, size, url, alt, uploaded_by, created_at
		FROM media WHERE tenant_id=$1 AND id=$2`,
		string(tenantID), string(id),
	)
	if err == sql.ErrNoRows {
		return nil, media.ErrMediaNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "get media by id", errx.TypeInternal)
	}
	return &m, nil
}

// Delete removes a media metadata row.
func (r *PostgresMediaRepository) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.MediaID) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM media WHERE tenant_id=$1 AND id=$2`,
		string(tenantID), string(id),
	)
	if err != nil {
		return errx.Wrap(err, "delete media", errx.TypeInternal)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return media.ErrMediaNotFound
	}
	return nil
}

// List returns paginated media records for a tenant.
func (r *PostgresMediaRepository) List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[media.Media], error) {
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM media WHERE tenant_id=$1`, string(tenantID)).Scan(&total)
	if err != nil {
		return kernel.Paginated[media.Media]{}, errx.Wrap(err, "count media", errx.TypeInternal)
	}

	var items []media.Media
	err = r.db.SelectContext(ctx, &items, `
		SELECT id, tenant_id, filename, content_type, size, url, alt, uploaded_by, created_at
		FROM media WHERE tenant_id=$1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.Paginated[media.Media]{}, errx.Wrap(err, "list media", errx.TypeInternal)
	}

	if items == nil {
		items = []media.Media{}
	}
	return kernel.NewPaginated(items, p.Page, p.PageSize, total), nil
}

// Ensure interface compliance at compile time.
var _ media.MediaRepository = (*PostgresMediaRepository)(nil)
