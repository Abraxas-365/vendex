package orderinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/order"
)

// PostgresRepo implements order.Repository using database/sql.
type PostgresRepo struct {
	db *sql.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed order repository.
func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Create(ctx context.Context, o *order.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	addrJSON, err := json.Marshal(o.ShippingAddress)
	if err != nil {
		return fmt.Errorf("marshaling address: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO orders (id, tenant_id, customer_id, status, total_amount, total_currency, shipping_address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		string(o.ID), string(o.TenantID), string(o.CustomerID),
		string(o.Status), o.TotalAmount.Amount, o.TotalAmount.Currency,
		string(addrJSON), o.CreatedAt, o.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting order: %w", err)
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
			return fmt.Errorf("inserting order item: %w", err)
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
		return fmt.Errorf("marshaling address: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE orders SET status=$1, total_amount=$2, total_currency=$3, shipping_address=$4, updated_at=$5
		WHERE id=$6 AND tenant_id=$7`,
		string(o.Status), o.TotalAmount.Amount, o.TotalAmount.Currency,
		string(addrJSON), o.UpdatedAt,
		string(o.ID), string(o.TenantID),
	)
	if err != nil {
		return fmt.Errorf("updating order: %w", err)
	}
	return nil
}

func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, pg kernel.Pagination) (kernel.PaginatedResult[order.Order], error) {
	return r.queryOrders(ctx, tenantID, pg, "", nil)
}

func (r *PostgresRepo) ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, pg kernel.Pagination) (kernel.PaginatedResult[order.Order], error) {
	return r.queryOrders(ctx, tenantID, pg, "AND customer_id = $3", []any{string(customerID)})
}

func (r *PostgresRepo) queryOrders(ctx context.Context, tenantID kernel.TenantID, pg kernel.Pagination, extraWhere string, extraArgs []any) (kernel.PaginatedResult[order.Order], error) {
	var zero kernel.PaginatedResult[order.Order]

	baseArgs := []any{string(tenantID)}
	args := append(baseArgs, extraArgs...)

	var total int
	countQ := "SELECT COUNT(*) FROM orders WHERE tenant_id = $1 " + extraWhere
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return zero, fmt.Errorf("counting orders: %w", err)
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
		return zero, fmt.Errorf("querying orders: %w", err)
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
		return zero, fmt.Errorf("iterating orders: %w", err)
	}

	return kernel.NewPaginatedResult(orders, total, pg), nil
}

func (r *PostgresRepo) scanOrder(ctx context.Context, query string, args ...any) (*order.Order, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	var o order.Order
	var id, tenantID, customerID, status, addrJSON string
	var createdAt, updatedAt time.Time

	err := row.Scan(&id, &tenantID, &customerID, &status,
		&o.TotalAmount.Amount, &o.TotalAmount.Currency,
		&addrJSON, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, order.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scanning order: %w", err)
	}

	o.ID = kernel.OrderID(id)
	o.TenantID = kernel.TenantID(tenantID)
	o.CustomerID = kernel.CustomerID(customerID)
	o.Status = order.OrderStatus(status)
	o.CreatedAt = createdAt
	o.UpdatedAt = updatedAt
	_ = json.Unmarshal([]byte(addrJSON), &o.ShippingAddress)

	return &o, nil
}

func (r *PostgresRepo) scanOrderRow(rows *sql.Rows) (*order.Order, error) {
	var o order.Order
	var id, tenantID, customerID, status, addrJSON string
	var createdAt, updatedAt time.Time

	err := rows.Scan(&id, &tenantID, &customerID, &status,
		&o.TotalAmount.Amount, &o.TotalAmount.Currency,
		&addrJSON, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("scanning order row: %w", err)
	}

	o.ID = kernel.OrderID(id)
	o.TenantID = kernel.TenantID(tenantID)
	o.CustomerID = kernel.CustomerID(customerID)
	o.Status = order.OrderStatus(status)
	o.CreatedAt = createdAt
	o.UpdatedAt = updatedAt
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
		return nil, fmt.Errorf("loading order items: %w", err)
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
			return nil, fmt.Errorf("scanning order item: %w", err)
		}
		item.ID = kernel.OrderItemID(id)
		item.ProductID = kernel.ProductID(productID)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating order items: %w", err)
	}

	if items == nil {
		items = []order.OrderItem{}
	}
	return items, nil
}

// Ensure interface compliance.
var _ order.Repository = (*PostgresRepo)(nil)
