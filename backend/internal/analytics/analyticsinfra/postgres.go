package analyticsinfra

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/analytics"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// PostgresRepo implements analytics.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed analytics repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// GetDashboardStats aggregates high-level metrics for the tenant.
func (r *PostgresRepo) GetDashboardStats(ctx context.Context, tenantID kernel.TenantID) (*analytics.DashboardStats, error) {
	stats := &analytics.DashboardStats{}
	tid := string(tenantID)

	// Total products.
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM products WHERE tenant_id = $1`, tid,
	).Scan(&stats.TotalProducts); err != nil {
		return nil, errx.Wrap(err, "counting products", errx.TypeInternal)
	}

	// Total orders.
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM orders WHERE tenant_id = $1`, tid,
	).Scan(&stats.TotalOrders); err != nil {
		return nil, errx.Wrap(err, "counting orders", errx.TypeInternal)
	}

	// Total customers.
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM customers WHERE tenant_id = $1`, tid,
	).Scan(&stats.TotalCustomers); err != nil {
		return nil, errx.Wrap(err, "counting customers", errx.TypeInternal)
	}

	// Total revenue and currency (sum of non-cancelled orders).
	row := r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(total_amount), 0), COALESCE(MAX(total_currency), '')
		   FROM orders
		  WHERE tenant_id = $1 AND status != 'cancelled'`, tid,
	)
	if err := row.Scan(&stats.TotalRevenue, &stats.Currency); err != nil {
		return nil, errx.Wrap(err, "summing revenue", errx.TypeInternal)
	}

	// Pending orders count.
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM orders WHERE tenant_id = $1 AND status = 'pending'`, tid,
	).Scan(&stats.PendingOrders); err != nil {
		return nil, errx.Wrap(err, "counting pending orders", errx.TypeInternal)
	}

	// Active promos count (active flag and within date range).
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM promos
		  WHERE tenant_id = $1
		    AND active = true
		    AND (starts_at IS NULL OR starts_at <= NOW())
		    AND (ends_at IS NULL OR ends_at >= NOW())`, tid,
	).Scan(&stats.ActivePromos); err != nil {
		return nil, errx.Wrap(err, "counting active promos", errx.TypeInternal)
	}

	// Pages pending review.
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM pages WHERE tenant_id = $1 AND status = 'pending_review'`, tid,
	).Scan(&stats.PendingPages); err != nil {
		return nil, errx.Wrap(err, "counting pending pages", errx.TypeInternal)
	}

	return stats, nil
}

// GetRevenueTimeline returns daily revenue totals for the last N days.
func (r *PostgresRepo) GetRevenueTimeline(ctx context.Context, tenantID kernel.TenantID, days int) ([]analytics.RevenuePoint, error) {
	since := time.Now().UTC().AddDate(0, 0, -days)

	type row struct {
		Date   time.Time `db:"date"`
		Amount int64     `db:"amount"`
	}

	var rows []row
	err := r.db.SelectContext(ctx, &rows, `
		SELECT DATE(created_at) AS date, SUM(total_amount) AS amount
		  FROM orders
		 WHERE tenant_id = $1
		   AND created_at >= $2
		   AND status != 'cancelled'
		 GROUP BY DATE(created_at)
		 ORDER BY date`,
		string(tenantID), since,
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying revenue timeline", errx.TypeInternal)
	}

	points := make([]analytics.RevenuePoint, len(rows))
	for i, rw := range rows {
		points[i] = analytics.RevenuePoint{
			Date:   rw.Date.Format("2006-01-02"),
			Amount: rw.Amount,
		}
	}
	return points, nil
}

// GetTopProducts returns the top N products by revenue from non-cancelled orders.
func (r *PostgresRepo) GetTopProducts(ctx context.Context, tenantID kernel.TenantID, limit int) ([]analytics.TopProduct, error) {
	type row struct {
		ProductID   string `db:"product_id"`
		ProductName string `db:"product_name"`
		TotalSold   int    `db:"total_sold"`
		Revenue     int64  `db:"revenue"`
	}

	var rows []row
	err := r.db.SelectContext(ctx, &rows, `
		SELECT oi.product_id, oi.product_name,
		       SUM(oi.quantity)     AS total_sold,
		       SUM(oi.total_amount) AS revenue
		  FROM order_items oi
		  JOIN orders o ON oi.order_id = o.id
		 WHERE o.tenant_id = $1
		   AND o.status != 'cancelled'
		 GROUP BY oi.product_id, oi.product_name
		 ORDER BY revenue DESC
		 LIMIT $2`,
		string(tenantID), limit,
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying top products", errx.TypeInternal)
	}

	products := make([]analytics.TopProduct, len(rows))
	for i, rw := range rows {
		products[i] = analytics.TopProduct{
			ProductID:   rw.ProductID,
			ProductName: rw.ProductName,
			TotalSold:   rw.TotalSold,
			Revenue:     rw.Revenue,
		}
	}
	return products, nil
}

// GetOrderStatusBreakdown returns order counts grouped by status for the tenant.
func (r *PostgresRepo) GetOrderStatusBreakdown(ctx context.Context, tenantID kernel.TenantID) ([]analytics.OrderStatusBreakdown, error) {
	type row struct {
		Status string `db:"status"`
		Count  int    `db:"count"`
	}

	var rows []row
	err := r.db.SelectContext(ctx, &rows, `
		SELECT status, COUNT(*) AS count
		  FROM orders
		 WHERE tenant_id = $1
		 GROUP BY status`,
		string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying order status breakdown", errx.TypeInternal)
	}

	breakdown := make([]analytics.OrderStatusBreakdown, len(rows))
	for i, rw := range rows {
		breakdown[i] = analytics.OrderStatusBreakdown{
			Status: rw.Status,
			Count:  rw.Count,
		}
	}
	return breakdown, nil
}

// GetRecentOrders returns the N most recent orders for the tenant.
func (r *PostgresRepo) GetRecentOrders(ctx context.Context, tenantID kernel.TenantID, limit int) ([]analytics.RecentOrder, error) {
	type row struct {
		ID          string    `db:"id"`
		CustomerID  string    `db:"customer_id"`
		Status      string    `db:"status"`
		TotalAmount int64     `db:"total_amount"`
		Currency    string    `db:"total_currency"`
		CreatedAt   time.Time `db:"created_at"`
	}

	var rows []row
	err := r.db.SelectContext(ctx, &rows, `
		SELECT id, customer_id, status, total_amount, total_currency, created_at
		  FROM orders
		 WHERE tenant_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2`,
		string(tenantID), limit,
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying recent orders", errx.TypeInternal)
	}

	orders := make([]analytics.RecentOrder, len(rows))
	for i, rw := range rows {
		orders[i] = analytics.RecentOrder{
			ID:          rw.ID,
			CustomerID:  rw.CustomerID,
			Status:      rw.Status,
			TotalAmount: rw.TotalAmount,
			Currency:    rw.Currency,
			CreatedAt:   rw.CreatedAt.UTC().Format(time.RFC3339),
		}
	}
	return orders, nil
}

// Ensure interface compliance at compile time.
var _ analytics.Repository = (*PostgresRepo)(nil)
