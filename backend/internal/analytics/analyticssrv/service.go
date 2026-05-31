package analyticssrv

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/analytics"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Service handles analytics query orchestration.
type Service struct {
	repo analytics.Repository
}

// New creates a new analytics service.
func New(repo analytics.Repository) *Service {
	return &Service{repo: repo}
}

// GetDashboardStats returns aggregate dashboard metrics for the tenant.
func (s *Service) GetDashboardStats(ctx context.Context, tenantID kernel.TenantID) (*analytics.DashboardStats, error) {
	return s.repo.GetDashboardStats(ctx, tenantID)
}

// GetRevenueTimeline returns daily revenue points for the last N days.
func (s *Service) GetRevenueTimeline(ctx context.Context, tenantID kernel.TenantID, days int) ([]analytics.RevenuePoint, error) {
	return s.repo.GetRevenueTimeline(ctx, tenantID, days)
}

// GetTopProducts returns the top N products by revenue.
func (s *Service) GetTopProducts(ctx context.Context, tenantID kernel.TenantID, limit int) ([]analytics.TopProduct, error) {
	return s.repo.GetTopProducts(ctx, tenantID, limit)
}

// GetOrderStatusBreakdown returns order counts grouped by status.
func (s *Service) GetOrderStatusBreakdown(ctx context.Context, tenantID kernel.TenantID) ([]analytics.OrderStatusBreakdown, error) {
	return s.repo.GetOrderStatusBreakdown(ctx, tenantID)
}

// GetRecentOrders returns the N most recent orders.
func (s *Service) GetRecentOrders(ctx context.Context, tenantID kernel.TenantID, limit int) ([]analytics.RecentOrder, error) {
	return s.repo.GetRecentOrders(ctx, tenantID, limit)
}
