// Package agent provides the AI store assistant harness, tools, and session management.
package agent

import (
	"context"
	"time"

	"github.com/Abraxas-365/vendex/internal/dashboard"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// StoreContext holds dynamic store statistics that are injected into the agent
// system prompt at session creation time.  All fields have sensible zero values
// so the prompt remains useful even when individual queries fail.
type StoreContext struct {
	// Counts
	ProductCount  int
	OrderCount    int
	CategoryCount int
	PromoCount    int

	// Revenue summary (last 30 days)
	RevenueLastMonth int64  // cents
	Currency         string // e.g. "USD"
}

// BuildStoreContext queries a best-effort snapshot of the tenant's store stats.
// Failures on any individual query are silently skipped so that session creation
// is never blocked by a transient DB error or an empty store.
func BuildStoreContext(ctx context.Context, tenantID kernel.TenantID, svc Services) StoreContext {
	var sc StoreContext

	// A tiny pagination request is enough — we only need the Total field.
	oneItem := kernel.NewPaginationOptions(1, 1)

	// Product count.
	if svc.Products != nil {
		if res, err := svc.Products.List(ctx, tenantID, oneItem); err == nil {
			sc.ProductCount = res.Total
		}
	}

	// Order count + revenue overview (last 30 days).
	if svc.Orders != nil {
		if res, err := svc.Orders.List(ctx, tenantID, oneItem); err == nil {
			sc.OrderCount = res.Total
		}
	}
	if svc.Dashboard != nil {
		dr := dashboard.DateRange{
			From: time.Now().AddDate(0, -1, 0),
			To:   time.Now(),
		}
		if overview, err := svc.Dashboard.GetSalesOverview(ctx, tenantID, dr); err == nil {
			sc.RevenueLastMonth = overview.TotalRevenue
			if overview.Currency != "" {
				sc.Currency = overview.Currency
			}
		}
	}
	if sc.Currency == "" {
		sc.Currency = "USD"
	}

	// Category count.
	if svc.Catalog != nil {
		if res, err := svc.Catalog.ListCategories(ctx, tenantID, oneItem); err == nil {
			sc.CategoryCount = res.Total
		}
	}

	// Promo count.
	if svc.Promos != nil {
		if res, err := svc.Promos.List(ctx, tenantID, oneItem); err == nil {
			sc.PromoCount = res.Total
		}
	}

	return sc
}
