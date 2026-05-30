package cartrecoveryinfra

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/cartrecovery"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// PostgresRepo implements cartrecovery.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed cart recovery repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Ensure interface compliance at compile time.
var _ cartrecovery.Repository = (*PostgresRepo)(nil)

// Create inserts a new recovery email row.
func (r *PostgresRepo) Create(ctx context.Context, email *cartrecovery.RecoveryEmail) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO cart_recovery_emails
			(id, tenant_id, cart_id, customer_id, email, step, status, discount_code, sent_at, clicked_at, converted_at, created_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		string(email.ID),
		string(email.TenantID),
		string(email.CartID),
		string(email.CustomerID),
		email.Email,
		email.Step,
		string(email.Status),
		email.DiscountCode,
		email.SentAt,
		email.ClickedAt,
		email.ConvertedAt,
		email.CreatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting cart recovery email", errx.TypeInternal)
	}
	return nil
}

// GetByID retrieves a recovery email by ID, scoped to a tenant.
func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.RecoveryID) (*cartrecovery.RecoveryEmail, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, cart_id, customer_id, email, step, status,
		       discount_code, sent_at, clicked_at, converted_at, created_at
		FROM cart_recovery_emails
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)

	rec, err := scanRecoveryEmail(row)
	if err == sql.ErrNoRows {
		return nil, cartrecovery.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning cart recovery email", errx.TypeInternal)
	}
	return rec, nil
}

// GetByCartID returns all recovery emails for a given cart, scoped to a tenant.
func (r *PostgresRepo) GetByCartID(ctx context.Context, tenantID kernel.TenantID, cartID kernel.CartID) ([]cartrecovery.RecoveryEmail, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, cart_id, customer_id, email, step, status,
		       discount_code, sent_at, clicked_at, converted_at, created_at
		FROM cart_recovery_emails
		WHERE tenant_id = $1 AND cart_id = $2
		ORDER BY created_at ASC`,
		string(tenantID), string(cartID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying cart recovery emails by cart", errx.TypeInternal)
	}
	defer rows.Close()

	return scanRecoveryEmailRows(rows)
}

// ListPending returns all pending recovery emails for a tenant.
func (r *PostgresRepo) ListPending(ctx context.Context, tenantID kernel.TenantID) ([]cartrecovery.RecoveryEmail, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, cart_id, customer_id, email, step, status,
		       discount_code, sent_at, clicked_at, converted_at, created_at
		FROM cart_recovery_emails
		WHERE tenant_id = $1 AND status = 'pending'
		ORDER BY created_at ASC`,
		string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying pending cart recovery emails", errx.TypeInternal)
	}
	defer rows.Close()

	return scanRecoveryEmailRows(rows)
}

// List returns a paginated list of recovery emails for a tenant.
func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, page, pageSize int) (kernel.Paginated[cartrecovery.RecoveryEmail], error) {
	p := kernel.NewPagination(page, pageSize)

	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM cart_recovery_emails WHERE tenant_id = $1`,
		string(tenantID),
	).Scan(&total); err != nil {
		return kernel.Paginated[cartrecovery.RecoveryEmail]{}, errx.Wrap(err, "counting cart recovery emails", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, cart_id, customer_id, email, step, status,
		       discount_code, sent_at, clicked_at, converted_at, created_at
		FROM cart_recovery_emails
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.Paginated[cartrecovery.RecoveryEmail]{}, errx.Wrap(err, "listing cart recovery emails", errx.TypeInternal)
	}
	defer rows.Close()

	items, err := scanRecoveryEmailRows(rows)
	if err != nil {
		return kernel.Paginated[cartrecovery.RecoveryEmail]{}, err
	}

	return kernel.NewPaginated(items, p.Page, p.PageSize, total), nil
}

// Update persists changes to an existing recovery email.
func (r *PostgresRepo) Update(ctx context.Context, email *cartrecovery.RecoveryEmail) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE cart_recovery_emails
		SET status = $1, discount_code = $2, sent_at = $3, clicked_at = $4, converted_at = $5
		WHERE id = $6 AND tenant_id = $7`,
		string(email.Status),
		email.DiscountCode,
		email.SentAt,
		email.ClickedAt,
		email.ConvertedAt,
		string(email.ID),
		string(email.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating cart recovery email", errx.TypeInternal)
	}
	return nil
}

// GetStats returns aggregate recovery statistics using conditional counts.
func (r *PostgresRepo) GetStats(ctx context.Context, tenantID kernel.TenantID) (*cartrecovery.RecoveryStats, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*)                                          AS total,
			COUNT(*) FILTER (WHERE status = 'pending')       AS pending,
			COUNT(*) FILTER (WHERE status = 'sent')          AS sent,
			COUNT(*) FILTER (WHERE status = 'clicked')       AS clicked,
			COUNT(*) FILTER (WHERE status = 'converted')     AS converted
		FROM cart_recovery_emails
		WHERE tenant_id = $1`,
		string(tenantID),
	)

	var stats cartrecovery.RecoveryStats
	if err := row.Scan(&stats.Total, &stats.Pending, &stats.Sent, &stats.Clicked, &stats.Converted); err != nil {
		return nil, errx.Wrap(err, "scanning cart recovery stats", errx.TypeInternal)
	}

	if stats.Total > 0 {
		stats.ConversionRate = float64(stats.Converted) / float64(stats.Total) * 100
	}

	return &stats, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

type singleRow interface {
	Scan(dest ...any) error
}

func scanRecoveryEmail(row singleRow) (*cartrecovery.RecoveryEmail, error) {
	var rec cartrecovery.RecoveryEmail
	var id, tenantID, cartID, customerID, status string

	err := row.Scan(
		&id, &tenantID, &cartID, &customerID,
		&rec.Email, &rec.Step, &status,
		&rec.DiscountCode, &rec.SentAt, &rec.ClickedAt, &rec.ConvertedAt,
		&rec.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	rec.ID = kernel.NewRecoveryID(id)
	rec.TenantID = kernel.TenantID(tenantID)
	rec.CartID = kernel.CartID(cartID)
	rec.CustomerID = kernel.CustomerID(customerID)
	rec.Status = cartrecovery.RecoveryStatus(status)

	return &rec, nil
}

func scanRecoveryEmailRows(rows *sql.Rows) ([]cartrecovery.RecoveryEmail, error) {
	var items []cartrecovery.RecoveryEmail
	for rows.Next() {
		rec, err := scanRecoveryEmail(rows)
		if err != nil {
			return nil, errx.Wrap(err, "scanning cart recovery email row", errx.TypeInternal)
		}
		items = append(items, *rec)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating cart recovery email rows", errx.TypeInternal)
	}
	if items == nil {
		items = []cartrecovery.RecoveryEmail{}
	}
	return items, nil
}
