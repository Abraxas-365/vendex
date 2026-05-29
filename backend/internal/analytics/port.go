package analytics

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines read-only analytics queries for the dashboard.
type Repository interface {
	GetDashboardStats(ctx context.Context, tenantID kernel.TenantID) (*DashboardStats, error)
	GetRevenueTimeline(ctx context.Context, tenantID kernel.TenantID, days int) ([]RevenuePoint, error)
	GetTopProducts(ctx context.Context, tenantID kernel.TenantID, limit int) ([]TopProduct, error)
	GetOrderStatusBreakdown(ctx context.Context, tenantID kernel.TenantID) ([]OrderStatusBreakdown, error)
	GetRecentOrders(ctx context.Context, tenantID kernel.TenantID, limit int) ([]RecentOrder, error)
}
