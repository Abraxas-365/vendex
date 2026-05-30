package dashboardsrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/dashboard"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Service provides dashboard reporting business logic.
type Service struct {
	repo dashboard.Repository
}

// New creates a new dashboard Service.
func New(repo dashboard.Repository) *Service {
	return &Service{repo: repo}
}

// validateRange checks that dr.From is before dr.To and applies defaults if zero.
func validateRange(dr *dashboard.DateRange) error {
	if dr.From.IsZero() {
		dr.From = time.Now().UTC().AddDate(0, -1, 0) // default: last 30 days
	}
	if dr.To.IsZero() {
		dr.To = time.Now().UTC()
	}
	if dr.From.After(dr.To) {
		return dashboard.ErrInvalidDateRange
	}
	return nil
}

// GetSalesOverview returns aggregated sales KPIs for the given date range.
func (s *Service) GetSalesOverview(ctx context.Context, tenantID kernel.TenantID, dr dashboard.DateRange) (dashboard.SalesOverview, error) {
	if err := validateRange(&dr); err != nil {
		return dashboard.SalesOverview{}, err
	}
	return s.repo.GetSalesOverview(ctx, tenantID, dr)
}

// GetTopProducts returns the top N products by revenue in the date range.
func (s *Service) GetTopProducts(ctx context.Context, tenantID kernel.TenantID, dr dashboard.DateRange, limit int) ([]dashboard.TopProduct, error) {
	if err := validateRange(&dr); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.GetTopProducts(ctx, tenantID, dr, limit)
}

// GetRevenueByDay returns a day-by-day revenue breakdown for the date range.
func (s *Service) GetRevenueByDay(ctx context.Context, tenantID kernel.TenantID, dr dashboard.DateRange) ([]dashboard.DailyRevenue, error) {
	if err := validateRange(&dr); err != nil {
		return nil, err
	}
	return s.repo.GetRevenueByDay(ctx, tenantID, dr)
}

// GetCustomerStats returns customer acquisition metrics for the date range.
func (s *Service) GetCustomerStats(ctx context.Context, tenantID kernel.TenantID, dr dashboard.DateRange) (dashboard.CustomerStats, error) {
	if err := validateRange(&dr); err != nil {
		return dashboard.CustomerStats{}, err
	}
	return s.repo.GetCustomerStats(ctx, tenantID, dr)
}

// GetConversionFunnel returns funnel metrics for the date range.
func (s *Service) GetConversionFunnel(ctx context.Context, tenantID kernel.TenantID, dr dashboard.DateRange) (dashboard.ConversionFunnel, error) {
	if err := validateRange(&dr); err != nil {
		return dashboard.ConversionFunnel{}, err
	}
	return s.repo.GetConversionFunnel(ctx, tenantID, dr)
}
