package returnsinfra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/returns"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements returns.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed returns repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Ensure interface compliance.
var _ returns.Repository = (*PostgresRepo)(nil)

func (r *PostgresRepo) Create(ctx context.Context, req *returns.ReturnRequest) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errx.Wrap(err, "beginning transaction", errx.TypeInternal)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO return_requests (
			id, tenant_id, order_id, customer_id, status, reason, notes, admin_notes,
			refund_amount_cents, refund_currency, resolution, tracking_number, carrier,
			created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
		string(req.ID), string(req.TenantID), string(req.OrderID), string(req.CustomerID),
		string(req.Status), req.Reason, nullableString(req.Notes), nullableString(req.AdminNotes),
		req.RefundAmount.Amount, req.RefundAmount.Currency,
		nullableString(string(req.Resolution)), nullableString(req.TrackingNumber), nullableString(req.Carrier),
		req.CreatedAt, req.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting return request", errx.TypeInternal)
	}

	for i := range req.Items {
		item := &req.Items[i]
		item.ReturnID = req.ID
		_, err = tx.ExecContext(ctx, `
			INSERT INTO return_items (
				id, return_id, tenant_id, product_id, variant_id, quantity, reason, condition, created_at
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
			string(item.ID), string(req.ID), string(req.TenantID),
			string(item.ProductID), nullableString(string(item.VariantID)),
			item.Quantity, nullableString(item.Reason), string(item.Condition), item.CreatedAt,
		)
		if err != nil {
			return errx.Wrap(err, "inserting return item", errx.TypeInternal)
		}
	}

	return tx.Commit()
}

func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ReturnID) (*returns.ReturnRequest, error) {
	req, err := r.scanRequest(ctx, `
		SELECT id, tenant_id, order_id, customer_id, status, reason, notes, admin_notes,
		       refund_amount_cents, refund_currency, resolution, tracking_number, carrier,
		       created_at, updated_at
		FROM return_requests
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return nil, err
	}

	items, err := r.loadItems(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	req.Items = items

	return req, nil
}

func (r *PostgresRepo) Update(ctx context.Context, req *returns.ReturnRequest) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE return_requests SET
			status=$1, admin_notes=$2, refund_amount_cents=$3, refund_currency=$4,
			resolution=$5, tracking_number=$6, carrier=$7, updated_at=$8
		WHERE id=$9 AND tenant_id=$10`,
		string(req.Status), nullableString(req.AdminNotes),
		req.RefundAmount.Amount, req.RefundAmount.Currency,
		nullableString(string(req.Resolution)), nullableString(req.TrackingNumber), nullableString(req.Carrier),
		req.UpdatedAt, string(req.ID), string(req.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating return request", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, status string, pg kernel.PaginationOptions) (kernel.Paginated[returns.ReturnRequest], error) {
	if status != "" {
		return r.queryRequests(ctx, tenantID, pg, "AND status = $3", []any{status})
	}
	return r.queryRequests(ctx, tenantID, pg, "", nil)
}

func (r *PostgresRepo) ListByOrder(ctx context.Context, tenantID kernel.TenantID, orderID kernel.OrderID, pg kernel.PaginationOptions) (kernel.Paginated[returns.ReturnRequest], error) {
	return r.queryRequests(ctx, tenantID, pg, "AND order_id = $3", []any{string(orderID)})
}

func (r *PostgresRepo) ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, pg kernel.PaginationOptions) (kernel.Paginated[returns.ReturnRequest], error) {
	return r.queryRequests(ctx, tenantID, pg, "AND customer_id = $3", []any{string(customerID)})
}

func (r *PostgresRepo) queryRequests(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions, extraWhere string, extraArgs []any) (kernel.Paginated[returns.ReturnRequest], error) {
	var zero kernel.Paginated[returns.ReturnRequest]

	baseArgs := []any{string(tenantID)}
	args := append(baseArgs, extraArgs...)

	var total int
	countQ := "SELECT COUNT(*) FROM return_requests WHERE tenant_id = $1 " + extraWhere
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting return requests", errx.TypeInternal)
	}

	nextParam := len(args) + 1
	dataQ := fmt.Sprintf(`
		SELECT id, tenant_id, order_id, customer_id, status, reason, notes, admin_notes,
		       refund_amount_cents, refund_currency, resolution, tracking_number, carrier,
		       created_at, updated_at
		FROM return_requests
		WHERE tenant_id = $1 %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, extraWhere, nextParam, nextParam+1)
	args = append(args, pg.Limit(), pg.Offset())

	rows, err := r.db.QueryContext(ctx, dataQ, args...)
	if err != nil {
		return zero, errx.Wrap(err, "querying return requests", errx.TypeInternal)
	}
	defer rows.Close()

	var reqs []returns.ReturnRequest
	for rows.Next() {
		req, err := r.scanRequestRow(rows)
		if err != nil {
			return zero, err
		}
		items, err := r.loadItems(ctx, req.ID)
		if err != nil {
			return zero, err
		}
		req.Items = items
		reqs = append(reqs, *req)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating return requests", errx.TypeInternal)
	}

	return kernel.NewPaginated(reqs, pg.Page, pg.PageSize, total), nil
}

func (r *PostgresRepo) scanRequest(ctx context.Context, query string, args ...any) (*returns.ReturnRequest, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	return r.scanRequestFromRow(row)
}

func (r *PostgresRepo) scanRequestFromRow(row *sql.Row) (*returns.ReturnRequest, error) {
	var req returns.ReturnRequest
	var id, tenantID, orderID, customerID, status, reason string
	var notes, adminNotes, resolution, trackingNumber, carrier sql.NullString
	var refundCurrency string

	err := row.Scan(
		&id, &tenantID, &orderID, &customerID, &status, &reason,
		&notes, &adminNotes,
		&req.RefundAmount.Amount, &refundCurrency,
		&resolution, &trackingNumber, &carrier,
		&req.CreatedAt, &req.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, returns.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning return request", errx.TypeInternal)
	}

	req.ID = kernel.ReturnID(id)
	req.TenantID = kernel.TenantID(tenantID)
	req.OrderID = kernel.OrderID(orderID)
	req.CustomerID = kernel.CustomerID(customerID)
	req.Status = returns.ReturnStatus(status)
	req.Reason = reason
	req.RefundAmount.Currency = refundCurrency
	if notes.Valid {
		req.Notes = notes.String
	}
	if adminNotes.Valid {
		req.AdminNotes = adminNotes.String
	}
	if resolution.Valid {
		req.Resolution = returns.Resolution(resolution.String)
	}
	if trackingNumber.Valid {
		req.TrackingNumber = trackingNumber.String
	}
	if carrier.Valid {
		req.Carrier = carrier.String
	}

	return &req, nil
}

func (r *PostgresRepo) scanRequestRow(rows *sql.Rows) (*returns.ReturnRequest, error) {
	var req returns.ReturnRequest
	var id, tenantID, orderID, customerID, status, reason string
	var notes, adminNotes, resolution, trackingNumber, carrier sql.NullString
	var refundCurrency string

	err := rows.Scan(
		&id, &tenantID, &orderID, &customerID, &status, &reason,
		&notes, &adminNotes,
		&req.RefundAmount.Amount, &refundCurrency,
		&resolution, &trackingNumber, &carrier,
		&req.CreatedAt, &req.UpdatedAt,
	)
	if err != nil {
		return nil, errx.Wrap(err, "scanning return request row", errx.TypeInternal)
	}

	req.ID = kernel.ReturnID(id)
	req.TenantID = kernel.TenantID(tenantID)
	req.OrderID = kernel.OrderID(orderID)
	req.CustomerID = kernel.CustomerID(customerID)
	req.Status = returns.ReturnStatus(status)
	req.Reason = reason
	req.RefundAmount.Currency = refundCurrency
	if notes.Valid {
		req.Notes = notes.String
	}
	if adminNotes.Valid {
		req.AdminNotes = adminNotes.String
	}
	if resolution.Valid {
		req.Resolution = returns.Resolution(resolution.String)
	}
	if trackingNumber.Valid {
		req.TrackingNumber = trackingNumber.String
	}
	if carrier.Valid {
		req.Carrier = carrier.String
	}

	return &req, nil
}

func (r *PostgresRepo) loadItems(ctx context.Context, returnID kernel.ReturnID) ([]returns.ReturnItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, return_id, tenant_id, product_id, variant_id, quantity, reason, condition, created_at
		FROM return_items WHERE return_id = $1`,
		string(returnID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "loading return items", errx.TypeInternal)
	}
	defer rows.Close()

	var items []returns.ReturnItem
	for rows.Next() {
		var item returns.ReturnItem
		var id, returnIDStr, tenantID, productID string
		var variantID, reason sql.NullString
		var condition string

		err := rows.Scan(
			&id, &returnIDStr, &tenantID, &productID,
			&variantID, &item.Quantity, &reason, &condition, &item.CreatedAt,
		)
		if err != nil {
			return nil, errx.Wrap(err, "scanning return item", errx.TypeInternal)
		}

		item.ID = kernel.ReturnItemID(id)
		item.ReturnID = kernel.ReturnID(returnIDStr)
		item.TenantID = kernel.TenantID(tenantID)
		item.ProductID = kernel.ProductID(productID)
		item.Condition = returns.ItemCondition(condition)
		if variantID.Valid {
			item.VariantID = kernel.VariantID(variantID.String)
		}
		if reason.Valid {
			item.Reason = reason.String
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating return items", errx.TypeInternal)
	}

	if items == nil {
		items = []returns.ReturnItem{}
	}
	return items, nil
}

func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
