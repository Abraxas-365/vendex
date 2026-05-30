package dashboard

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines the read-only queries needed by the dashboard domain.
// All methods are scoped by TenantID.
type Repository interface {
	// GetSalesOverview returns sales KPIs for the given date range.
	GetSalesOverview(ctx context.Context, tenantID kernel.TenantID, dr DateRange) (SalesOverview, error)

	// GetTopProducts returns the top N products by revenue in the date range.
	GetTopProducts(ctx context.Context, tenantID kernel.TenantID, dr DateRange, limit int) ([]TopProduct, error)

	// GetRevenueByDay returns daily revenue breakdown for the date range.
	GetRevenueByDay(ctx context.Context, tenantID kernel.TenantID, dr DateRange) ([]DailyRevenue, error)

	// GetCustomerStats returns customer acquisition metrics for the date range.
	GetCustomerStats(ctx context.Context, tenantID kernel.TenantID, dr DateRange) (CustomerStats, error)

	// GetConversionFunnel returns funnel metrics for the date range.
	GetConversionFunnel(ctx context.Context, tenantID kernel.TenantID, dr DateRange) (ConversionFunnel, error)
}
