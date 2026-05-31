package orderinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/order"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements order.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed order repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Create(ctx context.Context, o *order.Order) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errx.Wrap(err, "beginning transaction", errx.TypeInternal)
	}
	defer tx.Rollback()

	addrJSON, err := json.Marshal(o.ShippingAddress)
	if err != nil {
		return errx.Wrap(err, "marshaling address", errx.TypeInternal)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO orders (id, tenant_id, customer_id, status, total_amount, total_currency, shipping_address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		string(o.ID), string(o.TenantID), string(o.CustomerID),
		string(o.Status), o.TotalAmount.Amount, o.TotalAmount.Currency,
		string(addrJSON), o.CreatedAt, o.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting order", errx.TypeInternal)
	}

	for _, item := range o.Items {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO order_items (id, order_id, product_id, product_name, quantity, unit_price_amount, unit_price_currency, total_amount, total_currency)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			string(item.ID), string(o.ID), string(item.ProductID),
			item.ProductName, item.Quantity,
			item.UnitPrice.Amount, item.UnitPrice.Currency,
			item.Total.Amount, item.Total.Currency,
		)
		if err != nil {
			return errx.Wrap(err, "inserting order item", errx.TypeInternal)
		}
	}

	return tx.Commit()
}

func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.OrderID) (*order.Order, error) {
	o, err := r.scanOrder(ctx, `
		SELECT id, tenant_id, customer_id, status, total_amount, total_currency, shipping_address, created_at, updated_at
		FROM orders WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return nil, err
	}

	items, err := r.loadItems(ctx, o.ID)
	if err != nil {
		return nil, err
	}
	o.Items = items

	return o, nil
}

func (r *PostgresRepo) Update(ctx context.Context, o *order.Order) error {
	addrJSON, err := json.Marshal(o.ShippingAddress)
	if err != nil {
		return errx.Wrap(err, "marshaling address", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE orders SET status=$1, total_amount=$2, total_currency=$3, shipping_address=$4, updated_at=$5
		WHERE id=$6 AND tenant_id=$7`,
		string(o.Status), o.TotalAmount.Amount, o.TotalAmount.Currency,
		string(addrJSON), o.UpdatedAt,
		string(o.ID), string(o.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating order", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) UpdateCheckoutFields(ctx context.Context, o *order.Order) error {
	var billingJSON []byte
	if o.BillingAddress != nil {
		var err error
		billingJSON, err = json.Marshal(o.BillingAddress)
		if err != nil {
			return errx.Wrap(err, "marshaling billing address", errx.TypeInternal)
		}
	}

	_, err := r.db.ExecContext(ctx, `
		UPDATE orders SET
			subtotal_amount=$1, subtotal_currency=$2,
			shipping_amount=$3, shipping_currency=$4,
			tax_amount=$5, tax_currency=$6,
			discount_amount=$7, discount_currency=$8,
			shipping_method=$9, billing_address=$10,
			payment_status=$11, payment_method=$12,
			promo_code=$13, cart_id=$14,
			updated_at=$15
		WHERE id=$16 AND tenant_id=$17`,
		o.SubtotalAmount.Amount, o.SubtotalAmount.Currency,
		o.ShippingAmount.Amount, o.ShippingAmount.Currency,
		o.TaxAmount.Amount, o.TaxAmount.Currency,
		o.DiscountAmount.Amount, o.DiscountAmount.Currency,
		nullableString(o.ShippingMethod), billingJSON,
		nullableString(o.PaymentStatus), nullableString(o.PaymentMethod),
		nullableString(o.PromoCode), nullableString(o.CartID),
		o.UpdatedAt,
		string(o.ID), string(o.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating order checkout fields", errx.TypeInternal)
	}
	return nil
}

// nullableString returns nil for empty strings (stored as NULL in DB).
func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[order.Order], error) {
	return r.queryOrders(ctx, tenantID, pg, "", nil)
}

func (r *PostgresRepo) ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, pg kernel.PaginationOptions) (kernel.Paginated[order.Order], error) {
	return r.queryOrders(ctx, tenantID, pg, "AND customer_id = $3", []any{string(customerID)})
}

func (r *PostgresRepo) queryOrders(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions, extraWhere string, extraArgs []any) (kernel.Paginated[order.Order], error) {
	var zero kernel.Paginated[order.Order]

	baseArgs := []any{string(tenantID)}
	args := append(baseArgs, extraArgs...)

	var total int
	countQ := "SELECT COUNT(*) FROM orders WHERE tenant_id = $1 " + extraWhere
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting orders", errx.TypeInternal)
	}

	nextParam := len(args) + 1
	dataQ := fmt.Sprintf(`
		SELECT id, tenant_id, customer_id, status, total_amount, total_currency, shipping_address, created_at, updated_at
		FROM orders WHERE tenant_id = $1 %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, extraWhere, nextParam, nextParam+1)
	args = append(args, pg.Limit(), pg.Offset())

	rows, err := r.db.QueryContext(ctx, dataQ, args...)
	if err != nil {
		return zero, errx.Wrap(err, "querying orders", errx.TypeInternal)
	}
	defer rows.Close()

	var orders []order.Order
	for rows.Next() {
		o, err := r.scanOrderRow(rows)
		if err != nil {
			return zero, err
		}
		items, err := r.loadItems(ctx, o.ID)
		if err != nil {
			return zero, err
		}
		o.Items = items
		orders = append(orders, *o)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating orders", errx.TypeInternal)
	}

	return kernel.NewPaginated(orders, pg.Page, pg.PageSize, total), nil
}

func (r *PostgresRepo) scanOrder(ctx context.Context, query string, args ...any) (*order.Order, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	var o order.Order
	var id, tenantID, customerID, status, addrJSON string

	err := row.Scan(&id, &tenantID, &customerID, &status,
		&o.TotalAmount.Amount, &o.TotalAmount.Currency,
		&addrJSON, &o.CreatedAt, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, order.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning order", errx.TypeInternal)
	}

	o.ID = kernel.OrderID(id)
	o.TenantID = kernel.TenantID(tenantID)
	o.CustomerID = kernel.CustomerID(customerID)
	o.Status = order.OrderStatus(status)
	_ = json.Unmarshal([]byte(addrJSON), &o.ShippingAddress)

	return &o, nil
}

func (r *PostgresRepo) scanOrderRow(rows *sql.Rows) (*order.Order, error) {
	var o order.Order
	var id, tenantID, customerID, status, addrJSON string

	err := rows.Scan(&id, &tenantID, &customerID, &status,
		&o.TotalAmount.Amount, &o.TotalAmount.Currency,
		&addrJSON, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, errx.Wrap(err, "scanning order row", errx.TypeInternal)
	}

	o.ID = kernel.OrderID(id)
	o.TenantID = kernel.TenantID(tenantID)
	o.CustomerID = kernel.CustomerID(customerID)
	o.Status = order.OrderStatus(status)
	_ = json.Unmarshal([]byte(addrJSON), &o.ShippingAddress)

	return &o, nil
}

func (r *PostgresRepo) loadItems(ctx context.Context, orderID kernel.OrderID) ([]order.OrderItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, product_id, product_name, quantity, unit_price_amount, unit_price_currency, total_amount, total_currency
		FROM order_items WHERE order_id = $1`,
		string(orderID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "loading order items", errx.TypeInternal)
	}
	defer rows.Close()

	var items []order.OrderItem
	for rows.Next() {
		var item order.OrderItem
		var id, productID string
		err := rows.Scan(&id, &productID, &item.ProductName, &item.Quantity,
			&item.UnitPrice.Amount, &item.UnitPrice.Currency,
			&item.Total.Amount, &item.Total.Currency)
		if err != nil {
			return nil, errx.Wrap(err, "scanning order item", errx.TypeInternal)
		}
		item.ID = kernel.OrderItemID(id)
		item.ProductID = kernel.ProductID(productID)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating order items", errx.TypeInternal)
	}

	if items == nil {
		items = []order.OrderItem{}
	}
	return items, nil
}

// Ensure interface compliance.
var _ order.Repository = (*PostgresRepo)(nil)
