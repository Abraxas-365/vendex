package webhookinfra

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/webhook"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// PostgresRepo implements webhook.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed webhook repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Ensure interface compliance at compile time.
var _ webhook.Repository = (*PostgresRepo)(nil)

// ---------------------------------------------------------------------------
// Webhook CRUD
// ---------------------------------------------------------------------------

func (r *PostgresRepo) Create(ctx context.Context, wh *webhook.Webhook) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO webhooks (id, tenant_id, url, secret, events, active, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		string(wh.ID), string(wh.TenantID),
		wh.URL, wh.Secret, pq.Array(wh.Events),
		wh.Active, wh.Description,
		wh.CreatedAt, wh.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting webhook", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.WebhookID) (*webhook.Webhook, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, url, secret, events, active, description, created_at, updated_at
		FROM webhooks
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)

	wh, err := scanWebhook(row.Scan)
	if err == sql.ErrNoRows {
		return nil, webhook.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning webhook", errx.TypeInternal)
	}
	return wh, nil
}

func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, page, pageSize int) (kernel.Paginated[webhook.Webhook], error) {
	var zero kernel.Paginated[webhook.Webhook]
	pg := kernel.NewPaginationOptions(page, pageSize)

	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM webhooks WHERE tenant_id = $1`, string(tenantID),
	).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting webhooks", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, url, secret, events, active, description, created_at, updated_at
		FROM webhooks
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), pg.Limit(), pg.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "querying webhooks", errx.TypeInternal)
	}
	defer rows.Close()

	items, err := scanWebhooks(rows)
	if err != nil {
		return zero, err
	}

	return kernel.NewPaginated(items, pg.Page, pg.PageSize, total), nil
}

func (r *PostgresRepo) Update(ctx context.Context, wh *webhook.Webhook) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE webhooks SET
			url = $1, secret = $2, events = $3,
			active = $4, description = $5, updated_at = $6
		WHERE id = $7 AND tenant_id = $8`,
		wh.URL, wh.Secret, pq.Array(wh.Events),
		wh.Active, wh.Description, wh.UpdatedAt,
		string(wh.ID), string(wh.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating webhook", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.WebhookID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM webhooks WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting webhook", errx.TypeInternal)
	}
	return nil
}

// ListActiveByEvent returns all active webhooks for a tenant subscribed to the given event type.
func (r *PostgresRepo) ListActiveByEvent(ctx context.Context, tenantID kernel.TenantID, eventType string) ([]webhook.Webhook, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, url, secret, events, active, description, created_at, updated_at
		FROM webhooks
		WHERE tenant_id = $1 AND active = true AND $2 = ANY(events)
		ORDER BY created_at ASC`,
		string(tenantID), eventType,
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying active webhooks by event", errx.TypeInternal)
	}
	defer rows.Close()

	return scanWebhooks(rows)
}

// ---------------------------------------------------------------------------
// Delivery CRUD
// ---------------------------------------------------------------------------

func (r *PostgresRepo) CreateDelivery(ctx context.Context, d *webhook.WebhookDelivery) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO webhook_deliveries (
			id, tenant_id, webhook_id, event_type, payload,
			response_status, response_body, status,
			attempts, max_attempts, next_retry_at, delivered_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		string(d.ID), string(d.TenantID), string(d.WebhookID),
		d.EventType, string(d.Payload),
		d.ResponseStatus, d.ResponseBody, string(d.Status),
		d.Attempts, d.MaxAttempts, d.NextRetryAt, d.DeliveredAt, d.CreatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting webhook delivery", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) GetDelivery(ctx context.Context, tenantID kernel.TenantID, id kernel.WebhookDeliveryID) (*webhook.WebhookDelivery, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, webhook_id, event_type, payload,
		       response_status, response_body, status,
		       attempts, max_attempts, next_retry_at, delivered_at, created_at
		FROM webhook_deliveries
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)

	d, err := scanDelivery(row.Scan)
	if err == sql.ErrNoRows {
		return nil, webhook.ErrDeliveryNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning webhook delivery", errx.TypeInternal)
	}
	return d, nil
}

func (r *PostgresRepo) UpdateDelivery(ctx context.Context, d *webhook.WebhookDelivery) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE webhook_deliveries SET
			response_status = $1, response_body = $2,
			status = $3, attempts = $4,
			next_retry_at = $5, delivered_at = $6
		WHERE id = $7 AND tenant_id = $8`,
		d.ResponseStatus, d.ResponseBody,
		string(d.Status), d.Attempts,
		d.NextRetryAt, d.DeliveredAt,
		string(d.ID), string(d.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating webhook delivery", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) ListDeliveries(ctx context.Context, tenantID kernel.TenantID, webhookID kernel.WebhookID, page, pageSize int) (kernel.Paginated[webhook.WebhookDelivery], error) {
	var zero kernel.Paginated[webhook.WebhookDelivery]
	pg := kernel.NewPaginationOptions(page, pageSize)

	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM webhook_deliveries WHERE tenant_id = $1 AND webhook_id = $2`,
		string(tenantID), string(webhookID),
	).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting webhook deliveries", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, webhook_id, event_type, payload,
		       response_status, response_body, status,
		       attempts, max_attempts, next_retry_at, delivered_at, created_at
		FROM webhook_deliveries
		WHERE tenant_id = $1 AND webhook_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`,
		string(tenantID), string(webhookID), pg.Limit(), pg.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "querying webhook deliveries", errx.TypeInternal)
	}
	defer rows.Close()

	var items []webhook.WebhookDelivery
	for rows.Next() {
		d, err := scanDelivery(rows.Scan)
		if err != nil {
			return zero, errx.Wrap(err, "scanning webhook delivery row", errx.TypeInternal)
		}
		items = append(items, *d)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating webhook deliveries", errx.TypeInternal)
	}
	if items == nil {
		items = []webhook.WebhookDelivery{}
	}

	return kernel.NewPaginated(items, pg.Page, pg.PageSize, total), nil
}

// ---------------------------------------------------------------------------
// Scan helpers
// ---------------------------------------------------------------------------

type scanFunc func(dest ...any) error

func scanWebhook(scan scanFunc) (*webhook.Webhook, error) {
	var wh webhook.Webhook
	var id, tenantID string
	var events pq.StringArray

	err := scan(
		&id, &tenantID, &wh.URL, &wh.Secret, &events,
		&wh.Active, &wh.Description, &wh.CreatedAt, &wh.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	wh.ID = kernel.NewWebhookID(id)
	wh.TenantID = kernel.TenantID(tenantID)
	wh.Events = []string(events)

	return &wh, nil
}

func scanWebhooks(rows *sql.Rows) ([]webhook.Webhook, error) {
	var items []webhook.Webhook
	for rows.Next() {
		wh, err := scanWebhook(rows.Scan)
		if err != nil {
			return nil, errx.Wrap(err, "scanning webhook row", errx.TypeInternal)
		}
		items = append(items, *wh)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating webhooks", errx.TypeInternal)
	}
	if items == nil {
		items = []webhook.Webhook{}
	}
	return items, nil
}

func scanDelivery(scan scanFunc) (*webhook.WebhookDelivery, error) {
	var d webhook.WebhookDelivery
	var id, tenantID, webhookID, status string
	var payloadJSON string

	err := scan(
		&id, &tenantID, &webhookID, &d.EventType, &payloadJSON,
		&d.ResponseStatus, &d.ResponseBody, &status,
		&d.Attempts, &d.MaxAttempts, &d.NextRetryAt, &d.DeliveredAt, &d.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	d.ID = kernel.NewWebhookDeliveryID(id)
	d.TenantID = kernel.TenantID(tenantID)
	d.WebhookID = kernel.NewWebhookID(webhookID)
	d.Status = webhook.DeliveryStatus(status)
	d.Payload = json.RawMessage(payloadJSON)

	return &d, nil
}
