package cartinfra

import (
	"context"
	"database/sql"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/cart"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements cart.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed cart repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Ensure interface compliance.
var _ cart.Repository = (*PostgresRepo)(nil)

// Create inserts a new cart row (without items).
func (r *PostgresRepo) Create(ctx context.Context, c *cart.Cart) error {
	customerID := nullableString(string(c.CustomerID))
	sessionID := nullableString(c.SessionID)

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO carts (id, tenant_id, customer_id, session_id, currency, created_at, updated_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		string(c.ID), string(c.TenantID), customerID, sessionID,
		c.Currency, c.CreatedAt, c.UpdatedAt, c.ExpiresAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting cart", errx.TypeInternal)
	}
	return nil
}

// GetByID retrieves a cart by ID, scoped to tenant, and loads its items.
func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CartID) (*cart.Cart, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, COALESCE(customer_id,''), COALESCE(session_id,''), currency, created_at, updated_at, expires_at
		FROM carts WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	c, err := scanCart(row)
	if err == sql.ErrNoRows {
		return nil, cart.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning cart", errx.TypeInternal)
	}

	items, err := r.loadItems(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Items = items
	return c, nil
}

// GetBySession retrieves a cart by session ID, scoped to tenant.
func (r *PostgresRepo) GetBySession(ctx context.Context, tenantID kernel.TenantID, sessionID string) (*cart.Cart, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, COALESCE(customer_id,''), COALESCE(session_id,''), currency, created_at, updated_at, expires_at
		FROM carts WHERE tenant_id = $1 AND session_id = $2
		ORDER BY created_at DESC LIMIT 1`,
		string(tenantID), sessionID,
	)
	c, err := scanCart(row)
	if err == sql.ErrNoRows {
		return nil, cart.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning cart by session", errx.TypeInternal)
	}

	items, err := r.loadItems(ctx, c.ID)
	if err != nil {
		return nil, err
	}
	c.Items = items
	return c, nil
}

// GetByCustomer retrieves a cart by customer ID, scoped to tenant.
func (r *PostgresRepo) GetByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) (*cart.Cart, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, COALESCE(customer_id,''), COALESCE(session_id,''), currency, created_at, updated_at, expires_at
		FROM carts WHERE tenant_id = $1 AND customer_id = $2
		ORDER BY created_at DESC LIMIT 1`,
		string(tenantID), string(customerID),
	)
	c, err := scanCart(row)
	if err == sql.ErrNoRows {
		return nil, cart.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning cart by customer", errx.TypeInternal)
	}

	items, err := r.loadItems(ctx, c.ID)
	if err != nil {
		return nil, err
	}
	c.Items = items
	return c, nil
}

// Update persists changes to an existing cart row.
func (r *PostgresRepo) Update(ctx context.Context, c *cart.Cart) error {
	customerID := nullableString(string(c.CustomerID))
	sessionID := nullableString(c.SessionID)

	_, err := r.db.ExecContext(ctx, `
		UPDATE carts SET customer_id=$1, session_id=$2, currency=$3, updated_at=$4, expires_at=$5
		WHERE id=$6 AND tenant_id=$7`,
		customerID, sessionID, c.Currency, c.UpdatedAt, c.ExpiresAt,
		string(c.ID), string(c.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating cart", errx.TypeInternal)
	}
	return nil
}

// Delete removes a cart (cascade deletes items via FK).
func (r *PostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CartID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM carts WHERE id=$1 AND tenant_id=$2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting cart", errx.TypeInternal)
	}
	return nil
}

// SaveItems replaces all items for a cart: DELETE existing, then INSERT new ones.
func (r *PostgresRepo) SaveItems(ctx context.Context, cartID kernel.CartID, items []cart.CartItem) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errx.Wrap(err, "beginning transaction", errx.TypeInternal)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `DELETE FROM cart_items WHERE cart_id = $1`, string(cartID)); err != nil {
		return errx.Wrap(err, "deleting cart items", errx.TypeInternal)
	}

	for _, item := range items {
		variantID := nullableString(item.VariantID)
		_, err := tx.ExecContext(ctx, `
			INSERT INTO cart_items (id, cart_id, tenant_id, product_id, variant_id, quantity, unit_price_amount, unit_price_currency, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			string(item.ID), string(item.CartID), string(item.TenantID),
			string(item.ProductID), variantID, item.Quantity,
			item.UnitPrice.Amount, item.UnitPrice.Currency,
			item.CreatedAt, item.UpdatedAt,
		)
		if err != nil {
			return errx.Wrap(err, "inserting cart item", errx.TypeInternal)
		}
	}

	if err := tx.Commit(); err != nil {
		return errx.Wrap(err, "committing cart items", errx.TypeInternal)
	}
	return nil
}

// loadItems fetches all items for a given cart.
func (r *PostgresRepo) loadItems(ctx context.Context, cartID kernel.CartID) ([]cart.CartItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, cart_id, tenant_id, product_id, COALESCE(variant_id,''), quantity, unit_price_amount, unit_price_currency, created_at, updated_at
		FROM cart_items WHERE cart_id = $1 ORDER BY created_at ASC`,
		string(cartID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying cart items", errx.TypeInternal)
	}
	defer rows.Close()

	var items []cart.CartItem
	for rows.Next() {
		var item cart.CartItem
		var id, cid, tid, pid, vid string
		err := rows.Scan(
			&id, &cid, &tid, &pid, &vid,
			&item.Quantity,
			&item.UnitPrice.Amount, &item.UnitPrice.Currency,
			&item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, errx.Wrap(err, "scanning cart item", errx.TypeInternal)
		}
		item.ID = kernel.CartItemID(id)
		item.CartID = kernel.CartID(cid)
		item.TenantID = kernel.TenantID(tid)
		item.ProductID = kernel.ProductID(pid)
		item.VariantID = vid
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating cart items", errx.TypeInternal)
	}

	if items == nil {
		items = []cart.CartItem{}
	}
	return items, nil
}

// scanCart scans a single cart row from a *sql.Row.
func scanCart(row *sql.Row) (*cart.Cart, error) {
	var c cart.Cart
	var id, tenantID, customerID, sessionID string
	var createdAt, updatedAt, expiresAt time.Time

	err := row.Scan(&id, &tenantID, &customerID, &sessionID, &c.Currency, &createdAt, &updatedAt, &expiresAt)
	if err != nil {
		return nil, err
	}

	c.ID = kernel.CartID(id)
	c.TenantID = kernel.TenantID(tenantID)
	c.CustomerID = kernel.CustomerID(customerID)
	c.SessionID = sessionID
	c.CreatedAt = createdAt
	c.UpdatedAt = updatedAt
	c.ExpiresAt = expiresAt
	c.Items = []cart.CartItem{}
	return &c, nil
}

// nullableString converts an empty string to a sql NULL-compatible value.
func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
