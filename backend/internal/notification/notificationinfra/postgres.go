package notificationinfra

import (
	"context"
	"database/sql"
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/notification"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements notification.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed notification repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Ensure interface compliance at compile time.
var _ notification.Repository = (*PostgresRepo)(nil)

// Create persists a new notification record.
func (r *PostgresRepo) Create(ctx context.Context, n *notification.Notification) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO notifications (
			id, tenant_id, user_id, title, body, type,
			resource_type, resource_id, read, read_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		string(n.ID), string(n.TenantID), string(n.UserID),
		n.Title, n.Body, string(n.Type),
		nullString(n.ResourceType), nullString(n.ResourceID),
		n.Read, n.ReadAt, n.CreatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting notification", errx.TypeInternal)
	}
	return nil
}

// GetByID returns a notification scoped to the tenant.
func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.NotificationID) (*notification.Notification, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, user_id, title, body, type,
		       resource_type, resource_id, read, read_at, created_at
		FROM notifications
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)

	n, err := scanNotification(row.Scan)
	if err == sql.ErrNoRows {
		return nil, notification.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning notification", errx.TypeInternal)
	}
	return n, nil
}

// List returns paginated notifications for a user within a tenant.
func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, userID kernel.UserID, unreadOnly bool, page, pageSize int) (kernel.Paginated[notification.Notification], error) {
	var zero kernel.Paginated[notification.Notification]
	pg := kernel.NewPaginationOptions(page, pageSize)

	var total int
	var countErr error
	if unreadOnly {
		countErr = r.db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM notifications WHERE tenant_id = $1 AND user_id = $2 AND read = false`,
			string(tenantID), string(userID),
		).Scan(&total)
	} else {
		countErr = r.db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM notifications WHERE tenant_id = $1 AND user_id = $2`,
			string(tenantID), string(userID),
		).Scan(&total)
	}
	if countErr != nil {
		return zero, errx.Wrap(countErr, "counting notifications", errx.TypeInternal)
	}

	var rows *sql.Rows
	var queryErr error
	if unreadOnly {
		rows, queryErr = r.db.QueryContext(ctx, `
			SELECT id, tenant_id, user_id, title, body, type,
			       resource_type, resource_id, read, read_at, created_at
			FROM notifications
			WHERE tenant_id = $1 AND user_id = $2 AND read = false
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4`,
			string(tenantID), string(userID), pg.Limit(), pg.Offset(),
		)
	} else {
		rows, queryErr = r.db.QueryContext(ctx, `
			SELECT id, tenant_id, user_id, title, body, type,
			       resource_type, resource_id, read, read_at, created_at
			FROM notifications
			WHERE tenant_id = $1 AND user_id = $2
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4`,
			string(tenantID), string(userID), pg.Limit(), pg.Offset(),
		)
	}
	if queryErr != nil {
		return zero, errx.Wrap(queryErr, "querying notifications", errx.TypeInternal)
	}
	defer rows.Close()

	var items []notification.Notification
	for rows.Next() {
		n, err := scanNotification(rows.Scan)
		if err != nil {
			return zero, errx.Wrap(err, "scanning notification row", errx.TypeInternal)
		}
		items = append(items, *n)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating notifications", errx.TypeInternal)
	}
	if items == nil {
		items = []notification.Notification{}
	}

	return kernel.NewPaginated(items, pg.Page, pg.PageSize, total), nil
}

// MarkRead marks a single notification as read.
func (r *PostgresRepo) MarkRead(ctx context.Context, tenantID kernel.TenantID, id kernel.NotificationID) error {
	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, `
		UPDATE notifications SET read = true, read_at = $1
		WHERE id = $2 AND tenant_id = $3 AND read = false`,
		now, string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "marking notification as read", errx.TypeInternal)
	}
	return nil
}

// MarkAllRead marks all notifications for a user as read.
func (r *PostgresRepo) MarkAllRead(ctx context.Context, tenantID kernel.TenantID, userID kernel.UserID) error {
	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, `
		UPDATE notifications SET read = true, read_at = $1
		WHERE tenant_id = $2 AND user_id = $3 AND read = false`,
		now, string(tenantID), string(userID),
	)
	if err != nil {
		return errx.Wrap(err, "marking all notifications as read", errx.TypeInternal)
	}
	return nil
}

// GetUnreadCount returns the count of unread notifications for a user.
func (r *PostgresRepo) GetUnreadCount(ctx context.Context, tenantID kernel.TenantID, userID kernel.UserID) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM notifications WHERE tenant_id = $1 AND user_id = $2 AND read = false`,
		string(tenantID), string(userID),
	).Scan(&count)
	if err != nil {
		return 0, errx.Wrap(err, "counting unread notifications", errx.TypeInternal)
	}
	return count, nil
}

// Delete removes a notification.
func (r *PostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.NotificationID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM notifications WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting notification", errx.TypeInternal)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Scan helpers
// ---------------------------------------------------------------------------

type scanFunc func(dest ...any) error

func scanNotification(scan scanFunc) (*notification.Notification, error) {
	var n notification.Notification
	var id, tenantID, userID, notifType string
	var body, resourceType, resourceID sql.NullString

	err := scan(
		&id, &tenantID, &userID,
		&n.Title, &body, &notifType,
		&resourceType, &resourceID,
		&n.Read, &n.ReadAt, &n.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	n.ID = kernel.NewNotificationID(id)
	n.TenantID = kernel.TenantID(tenantID)
	n.UserID = kernel.UserID(userID)
	n.Type = notification.Type(notifType)

	if body.Valid {
		n.Body = body.String
	}
	if resourceType.Valid {
		n.ResourceType = resourceType.String
	}
	if resourceID.Valid {
		n.ResourceID = resourceID.String
	}

	return &n, nil
}

// nullString converts an empty string to a sql.NullString.
func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
