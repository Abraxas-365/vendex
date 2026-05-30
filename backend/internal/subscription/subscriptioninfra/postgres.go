package subscriptioninfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/subscription"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements subscription.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed subscription repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Ensure interface compliance at compile time.
var _ subscription.Repository = (*PostgresRepo)(nil)

// ---------------------------------------------------------------------------
// Subscription CRUD
// ---------------------------------------------------------------------------

func (r *PostgresRepo) Create(ctx context.Context, sub *subscription.Subscription) error {
	metaJSON, err := json.Marshal(sub.Metadata)
	if err != nil {
		return errx.Wrap(err, "marshaling subscription metadata", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO subscriptions (
			id, tenant_id, customer_id, product_id, variant_id,
			price_amount, price_currency, interval, status,
			next_billing_date, last_billed_at, cancelled_at, paused_at, trial_ends_at,
			metadata, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12, $13, $14,
			$15, $16, $17
		)`,
		string(sub.ID), string(sub.TenantID), string(sub.CustomerID), string(sub.ProductID),
		nullableID(sub.VariantID),
		sub.Price.Amount, sub.Price.Currency, string(sub.Interval), string(sub.Status),
		sub.NextBillingDate, sub.LastBilledAt, sub.CancelledAt, sub.PausedAt, sub.TrialEndsAt,
		string(metaJSON), sub.CreatedAt, sub.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting subscription", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.SubscriptionID) (*subscription.Subscription, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, customer_id, product_id, variant_id,
		       price_amount, price_currency, interval, status,
		       next_billing_date, last_billed_at, cancelled_at, paused_at, trial_ends_at,
		       metadata, created_at, updated_at
		FROM subscriptions
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)

	sub, err := scanSubscription(row.Scan)
	if err == sql.ErrNoRows {
		return nil, subscription.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning subscription", errx.TypeInternal)
	}
	return sub, nil
}

func (r *PostgresRepo) Update(ctx context.Context, sub *subscription.Subscription) error {
	metaJSON, err := json.Marshal(sub.Metadata)
	if err != nil {
		return errx.Wrap(err, "marshaling subscription metadata", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE subscriptions SET
			status = $1,
			next_billing_date = $2,
			last_billed_at = $3,
			cancelled_at = $4,
			paused_at = $5,
			trial_ends_at = $6,
			metadata = $7,
			updated_at = $8
		WHERE id = $9 AND tenant_id = $10`,
		string(sub.Status),
		sub.NextBillingDate, sub.LastBilledAt, sub.CancelledAt, sub.PausedAt, sub.TrialEndsAt,
		string(metaJSON), sub.UpdatedAt,
		string(sub.ID), string(sub.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating subscription", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.SubscriptionID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM subscriptions WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting subscription", errx.TypeInternal)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Listing
// ---------------------------------------------------------------------------

func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, page, pageSize int) (kernel.Paginated[subscription.Subscription], error) {
	var zero kernel.Paginated[subscription.Subscription]
	pg := kernel.NewPaginationOptions(page, pageSize)

	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM subscriptions WHERE tenant_id = $1`, string(tenantID),
	).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting subscriptions", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, customer_id, product_id, variant_id,
		       price_amount, price_currency, interval, status,
		       next_billing_date, last_billed_at, cancelled_at, paused_at, trial_ends_at,
		       metadata, created_at, updated_at
		FROM subscriptions
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), pg.Limit(), pg.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "querying subscriptions", errx.TypeInternal)
	}
	defer rows.Close()

	subs, err := scanSubscriptions(rows)
	if err != nil {
		return zero, err
	}

	return kernel.NewPaginated(subs, pg.Page, pg.PageSize, total), nil
}

func (r *PostgresRepo) ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) ([]subscription.Subscription, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, customer_id, product_id, variant_id,
		       price_amount, price_currency, interval, status,
		       next_billing_date, last_billed_at, cancelled_at, paused_at, trial_ends_at,
		       metadata, created_at, updated_at
		FROM subscriptions
		WHERE tenant_id = $1 AND customer_id = $2
		ORDER BY created_at DESC`,
		string(tenantID), string(customerID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying customer subscriptions", errx.TypeInternal)
	}
	defer rows.Close()

	return scanSubscriptions(rows)
}

func (r *PostgresRepo) ListDueBilling(ctx context.Context, tenantID kernel.TenantID, before time.Time) ([]subscription.Subscription, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, customer_id, product_id, variant_id,
		       price_amount, price_currency, interval, status,
		       next_billing_date, last_billed_at, cancelled_at, paused_at, trial_ends_at,
		       metadata, created_at, updated_at
		FROM subscriptions
		WHERE tenant_id = $1 AND status = 'active' AND next_billing_date <= $2
		ORDER BY next_billing_date ASC`,
		string(tenantID), before,
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying due subscriptions", errx.TypeInternal)
	}
	defer rows.Close()

	return scanSubscriptions(rows)
}

// ---------------------------------------------------------------------------
// Billing records
// ---------------------------------------------------------------------------

func (r *PostgresRepo) CreateBillingRecord(ctx context.Context, record *subscription.BillingRecord) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO billing_records (
			id, subscription_id, tenant_id,
			amount_cents, amount_currency,
			status, order_id, failure_reason,
			billed_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		string(record.ID), string(record.SubscriptionID), string(record.TenantID),
		record.Amount.Amount, record.Amount.Currency,
		record.Status, nullableOrderID(record.OrderID), record.FailureReason,
		record.BilledAt, record.CreatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting billing record", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) ListBillingRecords(ctx context.Context, tenantID kernel.TenantID, subID kernel.SubscriptionID, page, pageSize int) (kernel.Paginated[subscription.BillingRecord], error) {
	var zero kernel.Paginated[subscription.BillingRecord]
	pg := kernel.NewPaginationOptions(page, pageSize)

	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM billing_records WHERE tenant_id = $1 AND subscription_id = $2`,
		string(tenantID), string(subID),
	).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting billing records", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, subscription_id, tenant_id, amount_cents, amount_currency,
		       status, order_id, failure_reason, billed_at, created_at
		FROM billing_records
		WHERE tenant_id = $1 AND subscription_id = $2
		ORDER BY billed_at DESC
		LIMIT $3 OFFSET $4`,
		string(tenantID), string(subID), pg.Limit(), pg.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "querying billing records", errx.TypeInternal)
	}
	defer rows.Close()

	var records []subscription.BillingRecord
	for rows.Next() {
		rec, err := scanBillingRecord(rows)
		if err != nil {
			return zero, err
		}
		records = append(records, *rec)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating billing records", errx.TypeInternal)
	}
	if records == nil {
		records = []subscription.BillingRecord{}
	}

	return kernel.NewPaginated(records, pg.Page, pg.PageSize, total), nil
}

// ---------------------------------------------------------------------------
// Scan helpers
// ---------------------------------------------------------------------------

type scanFunc func(dest ...any) error

func scanSubscription(scan scanFunc) (*subscription.Subscription, error) {
	var sub subscription.Subscription
	var id, tenantID, customerID, productID, interval, status string
	var variantID sql.NullString
	var metaJSON string

	err := scan(
		&id, &tenantID, &customerID, &productID, &variantID,
		&sub.Price.Amount, &sub.Price.Currency, &interval, &status,
		&sub.NextBillingDate, &sub.LastBilledAt, &sub.CancelledAt, &sub.PausedAt, &sub.TrialEndsAt,
		&metaJSON, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	sub.ID = kernel.NewSubscriptionID(id)
	sub.TenantID = kernel.TenantID(tenantID)
	sub.CustomerID = kernel.NewCustomerID(customerID)
	sub.ProductID = kernel.NewProductID(productID)
	sub.Interval = subscription.BillingInterval(interval)
	sub.Status = subscription.SubscriptionStatus(status)

	if variantID.Valid {
		v := kernel.NewVariantID(variantID.String)
		sub.VariantID = &v
	}

	sub.Metadata = map[string]string{}
	_ = json.Unmarshal([]byte(metaJSON), &sub.Metadata)

	return &sub, nil
}

func scanSubscriptions(rows *sql.Rows) ([]subscription.Subscription, error) {
	var subs []subscription.Subscription
	for rows.Next() {
		sub, err := scanSubscription(rows.Scan)
		if err != nil {
			return nil, errx.Wrap(err, "scanning subscription row", errx.TypeInternal)
		}
		subs = append(subs, *sub)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating subscriptions", errx.TypeInternal)
	}
	if subs == nil {
		subs = []subscription.Subscription{}
	}
	return subs, nil
}

func scanBillingRecord(rows *sql.Rows) (*subscription.BillingRecord, error) {
	var rec subscription.BillingRecord
	var id, subID, tenantID string
	var orderID sql.NullString

	err := rows.Scan(
		&id, &subID, &tenantID,
		&rec.Amount.Amount, &rec.Amount.Currency,
		&rec.Status, &orderID, &rec.FailureReason,
		&rec.BilledAt, &rec.CreatedAt,
	)
	if err != nil {
		return nil, errx.Wrap(err, "scanning billing record", errx.TypeInternal)
	}

	rec.ID = kernel.NewBillingRecordID(id)
	rec.SubscriptionID = kernel.NewSubscriptionID(subID)
	rec.TenantID = kernel.TenantID(tenantID)

	if orderID.Valid {
		o := kernel.NewOrderID(orderID.String)
		rec.OrderID = &o
	}

	return &rec, nil
}

// ---------------------------------------------------------------------------
// Null helpers
// ---------------------------------------------------------------------------

func nullableID(v *kernel.VariantID) any {
	if v == nil {
		return nil
	}
	return string(*v)
}

func nullableOrderID(v *kernel.OrderID) any {
	if v == nil {
		return nil
	}
	return string(*v)
}
