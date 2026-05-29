package analyticsinfra

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/analytics"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// PostgresRepo implements analytics.Repository using database/sql.
type PostgresRepo struct {
	db *sql.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed analytics repository.
func NewPostgresRepo(db *sql.DB) *PostgresRepo {
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
		return nil, fmt.Errorf("counting products: %w", err)
	}

	// Total orders.
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM orders WHERE tenant_id = $1`, tid,
	).Scan(&stats.TotalOrders); err != nil {
		return nil, fmt.Errorf("counting orders: %w", err)
	}

	// Total customers.
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM customers WHERE tenant_id = $1`, tid,
	).Scan(&stats.TotalCustomers); err != nil {
		return nil, fmt.Errorf("counting customers: %w", err)
	}

	// Total revenue and currency (sum of non-cancelled orders).
	var revenue sql.NullInt64
	var currency sql.NullString
	if err := r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(total_amount), 0), MAX(total_currency)
		   FROM orders
		  WHERE tenant_id = $1 AND status != 'cancelled'`, tid,
	).Scan(&revenue, &currency); err != nil {
		return nil, fmt.Errorf("summing revenue: %w", err)
	}
	if revenue.Valid {
		stats.TotalRevenue = revenue.Int64
	}
	if currency.Valid {
		stats.Currency = currency.String
	}

	// Pending orders count.
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM orders WHERE tenant_id = $1 AND status = 'pending'`, tid,
	).Scan(&stats.PendingOrders); err != nil {
		return nil, fmt.Errorf("counting pending orders: %w", err)
	}

	// Active promos count (active flag and within date range).
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM promos
		  WHERE tenant_id = $1
		    AND active = true
		    AND (starts_at IS NULL OR starts_at <= NOW())
		    AND (ends_at IS NULL OR ends_at >= NOW())`, tid,
	).Scan(&stats.ActivePromos); err != nil {
		return nil, fmt.Errorf("counting active promos: %w", err)
	}

	// Pages pending review.
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM pages WHERE tenant_id = $1 AND status = 'pending_review'`, tid,
	).Scan(&stats.PendingPages); err != nil {
		return nil, fmt.Errorf("counting pending pages: %w", err)
	}

	return stats, nil
}

// GetRevenueTimeline returns daily revenue totals for the last N days.
func (r *PostgresRepo) GetRevenueTimeline(ctx context.Context, tenantID kernel.TenantID, days int) ([]analytics.RevenuePoint, error) {
	since := time.Now().UTC().AddDate(0, 0, -days)

	rows, err := r.db.QueryContext(ctx, `
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
		return nil, fmt.Errorf("querying revenue timeline: %w", err)
	}
	defer rows.Close()

	var points []analytics.RevenuePoint
	for rows.Next() {
		var p analytics.RevenuePoint
		var date time.Time
		if err := rows.Scan(&date, &p.Amount); err != nil {
			return nil, fmt.Errorf("scanning revenue point: %w", err)
		}
		p.Date = date.Format("2006-01-02")
		points = append(points, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating revenue timeline: %w", err)
	}

	if points == nil {
		points = []analytics.RevenuePoint{}
	}
	return points, nil
}

// GetTopProducts returns the top N products by revenue from non-cancelled orders.
func (r *PostgresRepo) GetTopProducts(ctx context.Context, tenantID kernel.TenantID, limit int) ([]analytics.TopProduct, error) {
	rows, err := r.db.QueryContext(ctx, `
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
		return nil, fmt.Errorf("querying top products: %w", err)
	}
	defer rows.Close()

	var products []analytics.TopProduct
	for rows.Next() {
		var p analytics.TopProduct
		if err := rows.Scan(&p.ProductID, &p.ProductName, &p.TotalSold, &p.Revenue); err != nil {
			return nil, fmt.Errorf("scanning top product: %w", err)
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating top products: %w", err)
	}

	if products == nil {
		products = []analytics.TopProduct{}
	}
	return products, nil
}

// GetOrderStatusBreakdown returns order counts grouped by status for the tenant.
func (r *PostgresRepo) GetOrderStatusBreakdown(ctx context.Context, tenantID kernel.TenantID) ([]analytics.OrderStatusBreakdown, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT status, COUNT(*) AS count
		  FROM orders
		 WHERE tenant_id = $1
		 GROUP BY status`,
		string(tenantID),
	)
	if err != nil {
		return nil, fmt.Errorf("querying order status breakdown: %w", err)
	}
	defer rows.Close()

	var breakdown []analytics.OrderStatusBreakdown
	for rows.Next() {
		var b analytics.OrderStatusBreakdown
		if err := rows.Scan(&b.Status, &b.Count); err != nil {
			return nil, fmt.Errorf("scanning order status: %w", err)
		}
		breakdown = append(breakdown, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating order status breakdown: %w", err)
	}

	if breakdown == nil {
		breakdown = []analytics.OrderStatusBreakdown{}
	}
	return breakdown, nil
}

// GetRecentOrders returns the N most recent orders for the tenant.
func (r *PostgresRepo) GetRecentOrders(ctx context.Context, tenantID kernel.TenantID, limit int) ([]analytics.RecentOrder, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, customer_id, status, total_amount, total_currency, created_at
		  FROM orders
		 WHERE tenant_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2`,
		string(tenantID), limit,
	)
	if err != nil {
		return nil, fmt.Errorf("querying recent orders: %w", err)
	}
	defer rows.Close()

	var orders []analytics.RecentOrder
	for rows.Next() {
		var o analytics.RecentOrder
		var createdAt time.Time
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.Status, &o.TotalAmount, &o.Currency, &createdAt); err != nil {
			return nil, fmt.Errorf("scanning recent order: %w", err)
		}
		o.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating recent orders: %w", err)
	}

	if orders == nil {
		orders = []analytics.RecentOrder{}
	}
	return orders, nil
}

// Ensure interface compliance at compile time.
var _ analytics.Repository = (*PostgresRepo)(nil)
