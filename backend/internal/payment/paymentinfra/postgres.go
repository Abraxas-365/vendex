package paymentinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/payment"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements payment.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed payment repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// ─── Payments ────────────────────────────────────────────────────────────────

func (r *PostgresRepo) CreatePayment(ctx context.Context, p *payment.Payment) error {
	providerData := p.ProviderData
	if providerData == nil {
		providerData = json.RawMessage("{}")
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO payments
			(id, tenant_id, order_id, amount_amount, amount_currency, status,
			 provider, provider_payment_id, provider_data, method, error_message, paid_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		string(p.ID), string(p.TenantID), string(p.OrderID),
		p.Amount.Amount, p.Amount.Currency,
		string(p.Status), p.Provider, p.ProviderPaymentID,
		string(providerData), p.Method, p.ErrorMessage,
		p.PaidAt, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting payment", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) GetPaymentByID(ctx context.Context, tenantID kernel.TenantID, id kernel.PaymentID) (*payment.Payment, error) {
	const q = `
		SELECT id, tenant_id, order_id, amount_amount, amount_currency, status,
		       provider, provider_payment_id, provider_data, method, error_message, paid_at, created_at, updated_at
		FROM payments
		WHERE id = $1 AND tenant_id = $2`

	row := r.db.QueryRowContext(ctx, q, string(id), string(tenantID))
	p, err := scanPayment(row)
	if err == sql.ErrNoRows {
		return nil, payment.ErrPaymentNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "getting payment by id", errx.TypeInternal)
	}
	return p, nil
}

func (r *PostgresRepo) GetPaymentByOrder(ctx context.Context, tenantID kernel.TenantID, orderID kernel.OrderID) (*payment.Payment, error) {
	const q = `
		SELECT id, tenant_id, order_id, amount_amount, amount_currency, status,
		       provider, provider_payment_id, provider_data, method, error_message, paid_at, created_at, updated_at
		FROM payments
		WHERE order_id = $1 AND tenant_id = $2
		ORDER BY created_at DESC
		LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, string(orderID), string(tenantID))
	p, err := scanPayment(row)
	if err == sql.ErrNoRows {
		return nil, payment.ErrPaymentNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "getting payment by order", errx.TypeInternal)
	}
	return p, nil
}

func (r *PostgresRepo) UpdatePayment(ctx context.Context, p *payment.Payment) error {
	providerData := p.ProviderData
	if providerData == nil {
		providerData = json.RawMessage("{}")
	}

	_, err := r.db.ExecContext(ctx, `
		UPDATE payments
		SET status=$1, provider_payment_id=$2, provider_data=$3, error_message=$4,
		    paid_at=$5, updated_at=$6
		WHERE id=$7 AND tenant_id=$8`,
		string(p.Status), p.ProviderPaymentID,
		string(providerData), p.ErrorMessage,
		p.PaidAt, p.UpdatedAt,
		string(p.ID), string(p.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating payment", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) ListPaymentsByOrder(ctx context.Context, tenantID kernel.TenantID, orderID kernel.OrderID) ([]payment.Payment, error) {
	const q = `
		SELECT id, tenant_id, order_id, amount_amount, amount_currency, status,
		       provider, provider_payment_id, provider_data, method, error_message, paid_at, created_at, updated_at
		FROM payments
		WHERE order_id = $1 AND tenant_id = $2
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, q, string(orderID), string(tenantID))
	if err != nil {
		return nil, errx.Wrap(err, "listing payments by order", errx.TypeInternal)
	}
	defer rows.Close()

	return scanPayments(rows)
}

// ─── Refunds ─────────────────────────────────────────────────────────────────

func (r *PostgresRepo) CreateRefund(ctx context.Context, ref *payment.Refund) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO refunds
			(id, tenant_id, payment_id, order_id, amount_amount, amount_currency,
			 reason, status, provider_refund_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		string(ref.ID), string(ref.TenantID), string(ref.PaymentID), string(ref.OrderID),
		ref.Amount.Amount, ref.Amount.Currency,
		ref.Reason, string(ref.Status), ref.ProviderRefundID,
		ref.CreatedAt, ref.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting refund", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) GetRefundByID(ctx context.Context, tenantID kernel.TenantID, id kernel.RefundID) (*payment.Refund, error) {
	const q = `
		SELECT id, tenant_id, payment_id, order_id, amount_amount, amount_currency,
		       reason, status, provider_refund_id, created_at, updated_at
		FROM refunds
		WHERE id = $1 AND tenant_id = $2`

	row := r.db.QueryRowContext(ctx, q, string(id), string(tenantID))
	ref, err := scanRefund(row)
	if err == sql.ErrNoRows {
		return nil, payment.ErrRefundNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "getting refund by id", errx.TypeInternal)
	}
	return ref, nil
}

func (r *PostgresRepo) ListRefundsByPayment(ctx context.Context, tenantID kernel.TenantID, paymentID kernel.PaymentID) ([]payment.Refund, error) {
	const q = `
		SELECT id, tenant_id, payment_id, order_id, amount_amount, amount_currency,
		       reason, status, provider_refund_id, created_at, updated_at
		FROM refunds
		WHERE payment_id = $1 AND tenant_id = $2
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, q, string(paymentID), string(tenantID))
	if err != nil {
		return nil, errx.Wrap(err, "listing refunds by payment", errx.TypeInternal)
	}
	defer rows.Close()

	return scanRefunds(rows)
}

func (r *PostgresRepo) UpdateRefund(ctx context.Context, ref *payment.Refund) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE refunds
		SET status=$1, provider_refund_id=$2, updated_at=$3
		WHERE id=$4 AND tenant_id=$5`,
		string(ref.Status), ref.ProviderRefundID, ref.UpdatedAt,
		string(ref.ID), string(ref.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating refund", errx.TypeInternal)
	}
	return nil
}

// ─── Scan helpers ─────────────────────────────────────────────────────────────

type rowScanner interface {
	Scan(dest ...any) error
}

func scanPayment(row rowScanner) (*payment.Payment, error) {
	var p payment.Payment
	var id, tenantID, orderID, status, provider, providerPaymentID, method, errorMessage string
	var providerDataStr string
	var paidAt *time.Time

	err := row.Scan(
		&id, &tenantID, &orderID,
		&p.Amount.Amount, &p.Amount.Currency,
		&status, &provider, &providerPaymentID,
		&providerDataStr, &method, &errorMessage,
		&paidAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	p.ID = kernel.PaymentID(id)
	p.TenantID = kernel.TenantID(tenantID)
	p.OrderID = kernel.OrderID(orderID)
	p.Status = payment.PaymentStatus(status)
	p.Provider = provider
	p.ProviderPaymentID = providerPaymentID
	p.Method = method
	p.ErrorMessage = errorMessage
	p.PaidAt = paidAt
	if providerDataStr != "" {
		p.ProviderData = json.RawMessage(providerDataStr)
	}

	return &p, nil
}

func scanPayments(rows *sql.Rows) ([]payment.Payment, error) {
	var payments []payment.Payment
	for rows.Next() {
		p, err := scanPayment(rows)
		if err != nil {
			return nil, errx.Wrap(err, "scanning payment row", errx.TypeInternal)
		}
		payments = append(payments, *p)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating payments", errx.TypeInternal)
	}
	if payments == nil {
		payments = []payment.Payment{}
	}
	return payments, nil
}

func scanRefund(row rowScanner) (*payment.Refund, error) {
	var ref payment.Refund
	var id, tenantID, paymentID, orderID, reason, status, providerRefundID string

	err := row.Scan(
		&id, &tenantID, &paymentID, &orderID,
		&ref.Amount.Amount, &ref.Amount.Currency,
		&reason, &status, &providerRefundID,
		&ref.CreatedAt, &ref.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	ref.ID = kernel.RefundID(id)
	ref.TenantID = kernel.TenantID(tenantID)
	ref.PaymentID = kernel.PaymentID(paymentID)
	ref.OrderID = kernel.OrderID(orderID)
	ref.Reason = reason
	ref.Status = payment.RefundStatus(status)
	ref.ProviderRefundID = providerRefundID

	return &ref, nil
}

func scanRefunds(rows *sql.Rows) ([]payment.Refund, error) {
	var refunds []payment.Refund
	for rows.Next() {
		ref, err := scanRefund(rows)
		if err != nil {
			return nil, errx.Wrap(err, "scanning refund row", errx.TypeInternal)
		}
		refunds = append(refunds, *ref)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating refunds", errx.TypeInternal)
	}
	if refunds == nil {
		refunds = []payment.Refund{}
	}
	return refunds, nil
}

// Ensure interface compliance.
var _ payment.Repository = (*PostgresRepo)(nil)
