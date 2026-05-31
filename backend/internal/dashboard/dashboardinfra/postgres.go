package dashboardinfra

import (
	"context"
	"time"

	"github.com/Abraxas-365/vendex/internal/dashboard"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements dashboard.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed dashboard repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// GetSalesOverview aggregates sales KPIs for the given date range.
func (r *PostgresRepo) GetSalesOverview(ctx context.Context, tenantID kernel.TenantID, dr dashboard.DateRange) (dashboard.SalesOverview, error) {
	var out dashboard.SalesOverview
	tid := string(tenantID)

	type row struct {
		OrderCount        int    `db:"order_count"`
		TotalRevenue      int64  `db:"total_revenue"`
		AverageOrderValue int64  `db:"avg_order_value"`
		Currency          string `db:"currency"`
	}
	var r1 row
	err := r.db.QueryRowxContext(ctx, `
		SELECT
			COUNT(*)                                  AS order_count,
			COALESCE(SUM(total_amount), 0)            AS total_revenue,
			COALESCE(AVG(total_amount)::BIGINT, 0)    AS avg_order_value,
			COALESCE(MAX(total_currency), 'USD')      AS currency
		FROM orders
		WHERE tenant_id = $1
		  AND created_at BETWEEN $2 AND $3
		  AND status != 'cancelled'`,
		tid, dr.From, dr.To,
	).StructScan(&r1)
	if err != nil {
		return out, errx.Wrap(err, "querying sales overview", errx.TypeInternal)
	}

	// Refund total = sum of cancelled orders in the same window.
	var refund int64
	if err := r.db.QueryRowxContext(ctx, `
		SELECT COALESCE(SUM(total_amount), 0)
		FROM orders
		WHERE tenant_id = $1
		  AND created_at BETWEEN $2 AND $3
		  AND status = 'cancelled'`,
		tid, dr.From, dr.To,
	).Scan(&refund); err != nil {
		return out, errx.Wrap(err, "querying refund total", errx.TypeInternal)
	}

	out = dashboard.SalesOverview{
		TotalRevenue:      r1.TotalRevenue,
		OrderCount:        r1.OrderCount,
		AverageOrderValue: r1.AverageOrderValue,
		RefundTotal:       refund,
		Currency:          r1.Currency,
	}
	return out, nil
}

// GetTopProducts returns the top N products by revenue in the date range.
func (r *PostgresRepo) GetTopProducts(ctx context.Context, tenantID kernel.TenantID, dr dashboard.DateRange, limit int) ([]dashboard.TopProduct, error) {
	type row struct {
		ProductID string `db:"product_id"`
		Name      string `db:"product_name"`
		Revenue   int64  `db:"revenue"`
		Quantity  int    `db:"quantity"`
	}

	var rows []row
	err := r.db.SelectContext(ctx, &rows, `
		SELECT
			oi.product_id,
			oi.product_name,
			SUM(oi.total_amount)  AS revenue,
			SUM(oi.quantity)      AS quantity
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		WHERE o.tenant_id = $1
		  AND o.created_at BETWEEN $2 AND $3
		  AND o.status != 'cancelled'
		GROUP BY oi.product_id, oi.product_name
		ORDER BY revenue DESC
		LIMIT $4`,
		string(tenantID), dr.From, dr.To, limit,
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying top products", errx.TypeInternal)
	}

	products := make([]dashboard.TopProduct, len(rows))
	for i, rw := range rows {
		products[i] = dashboard.TopProduct{
			ProductID: rw.ProductID,
			Name:      rw.Name,
			Revenue:   rw.Revenue,
			Quantity:  rw.Quantity,
		}
	}
	return products, nil
}

// GetRevenueByDay returns daily revenue for the date range.
func (r *PostgresRepo) GetRevenueByDay(ctx context.Context, tenantID kernel.TenantID, dr dashboard.DateRange) ([]dashboard.DailyRevenue, error) {
	type row struct {
		Date       time.Time `db:"date"`
		Revenue    int64     `db:"revenue"`
		OrderCount int       `db:"order_count"`
	}

	var rows []row
	err := r.db.SelectContext(ctx, &rows, `
		SELECT
			DATE(created_at)                   AS date,
			COALESCE(SUM(total_amount), 0)     AS revenue,
			COUNT(*)                           AS order_count
		FROM orders
		WHERE tenant_id = $1
		  AND created_at BETWEEN $2 AND $3
		  AND status != 'cancelled'
		GROUP BY DATE(created_at)
		ORDER BY date`,
		string(tenantID), dr.From, dr.To,
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying revenue by day", errx.TypeInternal)
	}

	points := make([]dashboard.DailyRevenue, len(rows))
	for i, rw := range rows {
		points[i] = dashboard.DailyRevenue{
			Date:       rw.Date.Format("2006-01-02"),
			Revenue:    rw.Revenue,
			OrderCount: rw.OrderCount,
		}
	}
	return points, nil
}

// GetCustomerStats returns customer acquisition metrics for the date range.
func (r *PostgresRepo) GetCustomerStats(ctx context.Context, tenantID kernel.TenantID, dr dashboard.DateRange) (dashboard.CustomerStats, error) {
	var total int
	if err := r.db.QueryRowxContext(ctx, `
		SELECT COUNT(*) FROM customers WHERE tenant_id = $1`,
		string(tenantID),
	).Scan(&total); err != nil {
		return dashboard.CustomerStats{}, errx.Wrap(err, "counting total customers", errx.TypeInternal)
	}

	var newCustomers int
	if err := r.db.QueryRowxContext(ctx, `
		SELECT COUNT(*) FROM customers
		WHERE tenant_id = $1 AND created_at BETWEEN $2 AND $3`,
		string(tenantID), dr.From, dr.To,
	).Scan(&newCustomers); err != nil {
		return dashboard.CustomerStats{}, errx.Wrap(err, "counting new customers", errx.TypeInternal)
	}

	returning := 0
	if total > newCustomers {
		returning = total - newCustomers
	}

	return dashboard.CustomerStats{
		TotalCustomers:     total,
		NewCustomers:       newCustomers,
		ReturningCustomers: returning,
	}, nil
}

// GetConversionFunnel returns funnel metrics for the date range.
// Visitors and checkout-started counts are not stored in the DB;
// they are returned as 0 and can be enriched by external analytics.
func (r *PostgresRepo) GetConversionFunnel(ctx context.Context, tenantID kernel.TenantID, dr dashboard.DateRange) (dashboard.ConversionFunnel, error) {
	// Carts created in window.
	var cartsCreated int
	if err := r.db.QueryRowxContext(ctx, `
		SELECT COUNT(*) FROM carts
		WHERE tenant_id = $1 AND created_at BETWEEN $2 AND $3`,
		string(tenantID), dr.From, dr.To,
	).Scan(&cartsCreated); err != nil {
		return dashboard.ConversionFunnel{}, errx.Wrap(err, "counting carts created", errx.TypeInternal)
	}

	// Orders completed (non-cancelled) in window.
	var ordersCompleted int
	if err := r.db.QueryRowxContext(ctx, `
		SELECT COUNT(*) FROM orders
		WHERE tenant_id = $1
		  AND created_at BETWEEN $2 AND $3
		  AND status != 'cancelled'`,
		string(tenantID), dr.From, dr.To,
	).Scan(&ordersCompleted); err != nil {
		return dashboard.ConversionFunnel{}, errx.Wrap(err, "counting completed orders", errx.TypeInternal)
	}

	return dashboard.ConversionFunnel{
		Visitors:         0, // not tracked in DB
		CartsCreated:     cartsCreated,
		CheckoutsStarted: 0, // not tracked in DB
		OrdersCompleted:  ordersCompleted,
	}, nil
}

// Compile-time interface check.
var _ dashboard.Repository = (*PostgresRepo)(nil)
