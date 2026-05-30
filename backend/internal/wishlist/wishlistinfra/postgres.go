package wishlistinfra

import (
	"context"
	"database/sql"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/wishlist"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements wishlist.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed wishlist repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Ensure interface compliance.
var _ wishlist.Repository = (*PostgresRepo)(nil)

// Create inserts a new wishlist row.
func (r *PostgresRepo) Create(ctx context.Context, w *wishlist.Wishlist) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO wishlists (id, tenant_id, customer_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`,
		string(w.ID), string(w.TenantID), string(w.CustomerID), w.CreatedAt, w.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting wishlist", errx.TypeInternal)
	}
	return nil
}

// GetByCustomer retrieves a wishlist by customer ID scoped to tenant and loads its items.
func (r *PostgresRepo) GetByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) (*wishlist.Wishlist, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, customer_id, created_at, updated_at
		FROM wishlists
		WHERE tenant_id = $1 AND customer_id = $2`,
		string(tenantID), string(customerID),
	)
	w, err := scanWishlist(row)
	if err == sql.ErrNoRows {
		return nil, wishlist.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning wishlist by customer", errx.TypeInternal)
	}

	items, err := r.loadItems(ctx, w.ID)
	if err != nil {
		return nil, err
	}
	w.Items = items
	return w, nil
}

// GetByID retrieves a wishlist by ID scoped to tenant and loads its items.
func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.WishlistID) (*wishlist.Wishlist, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, customer_id, created_at, updated_at
		FROM wishlists
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	w, err := scanWishlist(row)
	if err == sql.ErrNoRows {
		return nil, wishlist.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning wishlist by id", errx.TypeInternal)
	}

	items, err := r.loadItems(ctx, w.ID)
	if err != nil {
		return nil, err
	}
	w.Items = items
	return w, nil
}

// AddItem inserts a new item into the wishlist.
func (r *PostgresRepo) AddItem(ctx context.Context, wishlistID kernel.WishlistID, item *wishlist.WishlistItem) error {
	variantID := nullableString(item.VariantID)
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO wishlist_items (id, wishlist_id, product_id, variant_id, added_at)
		VALUES ($1, $2, $3, $4, $5)`,
		string(item.ID), string(wishlistID), string(item.ProductID), variantID, item.AddedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting wishlist item", errx.TypeInternal)
	}
	return nil
}

// RemoveItem deletes a wishlist item by ID.
func (r *PostgresRepo) RemoveItem(ctx context.Context, wishlistID kernel.WishlistID, itemID kernel.WishlistItemID) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM wishlist_items WHERE id = $1 AND wishlist_id = $2`,
		string(itemID), string(wishlistID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting wishlist item", errx.TypeInternal)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errx.Wrap(err, "checking rows affected", errx.TypeInternal)
	}
	if rows == 0 {
		return wishlist.ErrItemNotFound
	}
	return nil
}

// Delete removes a wishlist and all its items (cascade via FK).
func (r *PostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.WishlistID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM wishlists WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting wishlist", errx.TypeInternal)
	}
	return nil
}

// loadItems fetches all items for a given wishlist.
func (r *PostgresRepo) loadItems(ctx context.Context, wishlistID kernel.WishlistID) ([]wishlist.WishlistItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, wishlist_id, product_id, COALESCE(variant_id, ''), added_at
		FROM wishlist_items
		WHERE wishlist_id = $1
		ORDER BY added_at ASC`,
		string(wishlistID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying wishlist items", errx.TypeInternal)
	}
	defer rows.Close()

	var items []wishlist.WishlistItem
	for rows.Next() {
		var item wishlist.WishlistItem
		var id, wid, pid, vid string
		if err := rows.Scan(&id, &wid, &pid, &vid, &item.AddedAt); err != nil {
			return nil, errx.Wrap(err, "scanning wishlist item", errx.TypeInternal)
		}
		item.ID = kernel.WishlistItemID(id)
		item.WishlistID = kernel.WishlistID(wid)
		item.ProductID = kernel.ProductID(pid)
		item.VariantID = vid
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating wishlist items", errx.TypeInternal)
	}

	if items == nil {
		items = []wishlist.WishlistItem{}
	}
	return items, nil
}

// scanWishlist scans a single wishlist row.
func scanWishlist(row *sql.Row) (*wishlist.Wishlist, error) {
	var w wishlist.Wishlist
	var id, tenantID, customerID string
	var createdAt, updatedAt time.Time

	err := row.Scan(&id, &tenantID, &customerID, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	w.ID = kernel.WishlistID(id)
	w.TenantID = kernel.TenantID(tenantID)
	w.CustomerID = kernel.CustomerID(customerID)
	w.CreatedAt = createdAt
	w.UpdatedAt = updatedAt
	w.Items = []wishlist.WishlistItem{}
	return &w, nil
}

// nullableString converts an empty string to a SQL NULL-compatible value.
func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
