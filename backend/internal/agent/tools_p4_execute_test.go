package agent

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/Abraxas-365/vendex/internal/audit"
	"github.com/Abraxas-365/vendex/internal/audit/auditsrv"
	"github.com/Abraxas-365/vendex/internal/bundle"
	"github.com/Abraxas-365/vendex/internal/bundle/bundlesrv"
	"github.com/Abraxas-365/vendex/internal/dashboard"
	"github.com/Abraxas-365/vendex/internal/dashboard/dashboardsrv"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/inventory"
	"github.com/Abraxas-365/vendex/internal/inventory/inventorysrv"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/loyalty"
	"github.com/Abraxas-365/vendex/internal/loyalty/loyaltysrv"
	"github.com/Abraxas-365/vendex/internal/notification"
	"github.com/Abraxas-365/vendex/internal/notification/notificationsrv"
	"github.com/Abraxas-365/vendex/internal/returns"
	"github.com/Abraxas-365/vendex/internal/returns/returnssrv"
	"github.com/Abraxas-365/vendex/internal/review"
	"github.com/Abraxas-365/vendex/internal/review/reviewsrv"
	"github.com/Abraxas-365/vendex/internal/webhook"
	"github.com/Abraxas-365/vendex/internal/webhook/webhooksrv"
)

// ─── Inventory Stubs ────────────────────────────────────────────────────────

type stubInventoryRepo struct{}

func (s *stubInventoryRepo) CreateWarehouse(_ context.Context, w inventory.Warehouse) (inventory.Warehouse, error) {
	w.ID = "wh-1"
	return w, nil
}
func (s *stubInventoryRepo) GetWarehouse(_ context.Context, _ kernel.TenantID, id kernel.WarehouseID) (inventory.Warehouse, error) {
	return inventory.Warehouse{ID: id, TenantID: testTenant, Name: "Main WH"}, nil
}
func (s *stubInventoryRepo) UpdateWarehouse(_ context.Context, w inventory.Warehouse) (inventory.Warehouse, error) {
	return w, nil
}
func (s *stubInventoryRepo) DeleteWarehouse(_ context.Context, _ kernel.TenantID, _ kernel.WarehouseID) error {
	return nil
}
func (s *stubInventoryRepo) ListWarehouses(_ context.Context, _ kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[inventory.Warehouse], error) {
	return kernel.Paginated[inventory.Warehouse]{
		Items: []inventory.Warehouse{{ID: "wh-1", TenantID: testTenant, Name: "Main", Address: "123 St", IsDefault: true}},
		Total: 1, Page: pg.Page, TotalPages: 1,
	}, nil
}
func (s *stubInventoryRepo) GetStockLevel(_ context.Context, _ kernel.TenantID, _ kernel.ProductID, _ *kernel.VariantID, _ kernel.WarehouseID) (inventory.StockLevel, error) {
	return inventory.StockLevel{ProductID: "p-1", WarehouseID: "wh-1", Quantity: 100, Reserved: 5, LowStockThreshold: 10}, nil
}
func (s *stubInventoryRepo) UpsertStockLevel(_ context.Context, sl inventory.StockLevel) (inventory.StockLevel, error) {
	return sl, nil
}
func (s *stubInventoryRepo) ListStockLevels(_ context.Context, _ kernel.TenantID, _ kernel.ProductID) ([]inventory.StockLevel, error) {
	return []inventory.StockLevel{}, nil
}
func (s *stubInventoryRepo) GetLowStockItems(_ context.Context, _ kernel.TenantID) ([]inventory.StockLevel, error) {
	return []inventory.StockLevel{}, nil
}
func (s *stubInventoryRepo) CreateMovement(_ context.Context, m inventory.StockMovement) (inventory.StockMovement, error) {
	m.ID = "mv-1"
	return m, nil
}
func (s *stubInventoryRepo) ListMovements(_ context.Context, _ kernel.TenantID, _ kernel.ProductID, pg kernel.PaginationOptions) (kernel.Paginated[inventory.StockMovement], error) {
	return kernel.Paginated[inventory.StockMovement]{Items: []inventory.StockMovement{}, Total: 0, Page: pg.Page, TotalPages: 0}, nil
}

// ─── Review Stubs ───────────────────────────────────────────────────────────

type stubReviewRepo struct{}

func (s *stubReviewRepo) Create(_ context.Context, r review.Review) (review.Review, error) {
	r.ID = "rv-1"
	return r, nil
}
func (s *stubReviewRepo) GetByID(_ context.Context, _ kernel.TenantID, id kernel.ReviewID) (review.Review, error) {
	return review.Review{ID: id, TenantID: testTenant, ProductID: "p-1", CustomerID: "c-1", Rating: 5, Title: "Great", Status: "pending"}, nil
}
func (s *stubReviewRepo) List(_ context.Context, _ kernel.TenantID, _ string, pg kernel.PaginationOptions) (kernel.Paginated[review.Review], error) {
	return kernel.Paginated[review.Review]{
		Items: []review.Review{{ID: "rv-1", TenantID: testTenant, ProductID: "p-1", CustomerID: "c-1", Rating: 5, Title: "Great", Status: "pending"}},
		Total: 1, Page: pg.Page, TotalPages: 1,
	}, nil
}
func (s *stubReviewRepo) ListByProduct(_ context.Context, _ kernel.TenantID, _ kernel.ProductID, _ string, pg kernel.PaginationOptions) (kernel.Paginated[review.Review], error) {
	return kernel.Paginated[review.Review]{Items: []review.Review{}, Total: 0, Page: pg.Page, TotalPages: 0}, nil
}
func (s *stubReviewRepo) ListByCustomer(_ context.Context, _ kernel.TenantID, _ kernel.CustomerID, pg kernel.PaginationOptions) (kernel.Paginated[review.Review], error) {
	return kernel.Paginated[review.Review]{Items: []review.Review{}, Total: 0, Page: pg.Page, TotalPages: 0}, nil
}
func (s *stubReviewRepo) UpdateStatus(_ context.Context, _ kernel.TenantID, id kernel.ReviewID, status review.ReviewStatus) (review.Review, error) {
	return review.Review{ID: id, TenantID: testTenant, ProductID: "p-1", Status: review.ReviewStatus(status)}, nil
}
func (s *stubReviewRepo) GetStats(_ context.Context, _ kernel.TenantID, _ kernel.ProductID) (review.ReviewStats, error) {
	return review.ReviewStats{}, nil
}
func (s *stubReviewRepo) IncrementHelpful(_ context.Context, _ kernel.TenantID, _ kernel.ReviewID) error {
	return nil
}
func (s *stubReviewRepo) SetAdminResponse(_ context.Context, _ kernel.TenantID, _ kernel.ReviewID, _ string) (review.Review, error) {
	return review.Review{}, nil
}

// ─── Returns Stubs ──────────────────────────────────────────────────────────

type stubReturnsRepo struct{}

func (s *stubReturnsRepo) Create(_ context.Context, r *returns.ReturnRequest) error {
	r.ID = "ret-1"
	return nil
}
func (s *stubReturnsRepo) GetByID(_ context.Context, _ kernel.TenantID, id kernel.ReturnID) (*returns.ReturnRequest, error) {
	return &returns.ReturnRequest{ID: id, TenantID: testTenant, OrderID: "o-1", CustomerID: "c-1", Status: "requested", Reason: "defective"}, nil
}
func (s *stubReturnsRepo) Update(_ context.Context, r *returns.ReturnRequest) error {
	return nil
}
func (s *stubReturnsRepo) List(_ context.Context, _ kernel.TenantID, _ string, pg kernel.PaginationOptions) (kernel.Paginated[returns.ReturnRequest], error) {
	return kernel.Paginated[returns.ReturnRequest]{
		Items: []returns.ReturnRequest{{ID: "ret-1", TenantID: testTenant, OrderID: "o-1", CustomerID: "c-1", Status: "pending", Reason: "defective"}},
		Total: 1, Page: pg.Page, TotalPages: 1,
	}, nil
}
func (s *stubReturnsRepo) ListByOrder(_ context.Context, _ kernel.TenantID, _ kernel.OrderID, pg kernel.PaginationOptions) (kernel.Paginated[returns.ReturnRequest], error) {
	return kernel.Paginated[returns.ReturnRequest]{Items: []returns.ReturnRequest{}, Total: 0, Page: pg.Page, TotalPages: 0}, nil
}
func (s *stubReturnsRepo) ListByCustomer(_ context.Context, _ kernel.TenantID, _ kernel.CustomerID, pg kernel.PaginationOptions) (kernel.Paginated[returns.ReturnRequest], error) {
	return kernel.Paginated[returns.ReturnRequest]{Items: []returns.ReturnRequest{}, Total: 0, Page: pg.Page, TotalPages: 0}, nil
}

// ─── Webhook Stubs ──────────────────────────────────────────────────────────

type stubWebhookRepo struct{}

func (s *stubWebhookRepo) Create(_ context.Context, w *webhook.Webhook) error {
	w.ID = "wh-1"
	return nil
}
func (s *stubWebhookRepo) GetByID(_ context.Context, _ kernel.TenantID, id kernel.WebhookID) (*webhook.Webhook, error) {
	return &webhook.Webhook{ID: id, TenantID: testTenant, URL: "https://example.com/hook", Active: true, Events: []string{"order.created"}}, nil
}
func (s *stubWebhookRepo) Update(_ context.Context, w *webhook.Webhook) error {
	return nil
}
func (s *stubWebhookRepo) Delete(_ context.Context, _ kernel.TenantID, _ kernel.WebhookID) error {
	return nil
}
func (s *stubWebhookRepo) List(_ context.Context, _ kernel.TenantID, page, pageSize int) (kernel.Paginated[webhook.Webhook], error) {
	return kernel.Paginated[webhook.Webhook]{
		Items: []webhook.Webhook{{ID: "wh-1", TenantID: testTenant, URL: "https://example.com/hook", Active: true, Events: []string{"order.created"}}},
		Total: 1, Page: page, TotalPages: 1,
	}, nil
}
func (s *stubWebhookRepo) ListActiveByEvent(_ context.Context, _ kernel.TenantID, _ string) ([]webhook.Webhook, error) {
	return []webhook.Webhook{}, nil
}
func (s *stubWebhookRepo) CreateDelivery(_ context.Context, d *webhook.WebhookDelivery) error {
	return nil
}
func (s *stubWebhookRepo) UpdateDelivery(_ context.Context, d *webhook.WebhookDelivery) error {
	return nil
}
func (s *stubWebhookRepo) GetDelivery(_ context.Context, _ kernel.TenantID, id kernel.WebhookDeliveryID) (*webhook.WebhookDelivery, error) {
	panic("unused")
}
func (s *stubWebhookRepo) ListDeliveries(_ context.Context, _ kernel.TenantID, _ kernel.WebhookID, _, _ int) (kernel.Paginated[webhook.WebhookDelivery], error) {
	panic("unused")
}

// ─── Audit Stubs ────────────────────────────────────────────────────────────

type stubAuditRepo struct{}

func (s *stubAuditRepo) Create(_ context.Context, e audit.AuditEntry) (audit.AuditEntry, error) {
	e.ID = "au-1"
	return e, nil
}
func (s *stubAuditRepo) GetByID(_ context.Context, _ kernel.TenantID, id kernel.AuditEntryID) (audit.AuditEntry, error) {
	return audit.AuditEntry{ID: id, TenantID: testTenant, UserID: "u-1", Action: "product.create"}, nil
}
func (s *stubAuditRepo) List(_ context.Context, _ kernel.TenantID, _ audit.AuditFilter, pg kernel.PaginationOptions) (kernel.Paginated[audit.AuditEntry], error) {
	return kernel.Paginated[audit.AuditEntry]{
		Items: []audit.AuditEntry{{ID: "au-1", TenantID: testTenant, UserID: "u-1", Action: "product.create", ResourceType: "product"}},
		Total: 1, Page: pg.Page, TotalPages: 1,
	}, nil
}
func (s *stubAuditRepo) CountByAction(_ context.Context, _ kernel.TenantID, _, _ time.Time) ([]audit.ActionStats, error) {
	return []audit.ActionStats{}, nil
}

// ─── Loyalty Stubs ──────────────────────────────────────────────────────────

type stubLoyaltyRepo struct{}

func (s *stubLoyaltyRepo) GetOrCreateAccount(_ context.Context, _ kernel.TenantID, _ kernel.CustomerID) (loyalty.LoyaltyAccount, error) {
	return loyalty.LoyaltyAccount{ID: "la-1", TenantID: testTenant, CustomerID: "c-1", PointsBalance: 100, Tier: "bronze", LifetimePoints: 500}, nil
}
func (s *stubLoyaltyRepo) GetAccountByID(_ context.Context, _ kernel.TenantID, id kernel.LoyaltyAccountID) (loyalty.LoyaltyAccount, error) {
	return loyalty.LoyaltyAccount{ID: id, TenantID: testTenant, CustomerID: "c-1", PointsBalance: 100, Tier: "bronze", LifetimePoints: 500}, nil
}
func (s *stubLoyaltyRepo) GetAccountByCustomerID(_ context.Context, _ kernel.TenantID, _ kernel.CustomerID) (loyalty.LoyaltyAccount, error) {
	return loyalty.LoyaltyAccount{ID: "la-1", TenantID: testTenant, CustomerID: "c-1", PointsBalance: 100, Tier: "bronze", LifetimePoints: 500}, nil
}
func (s *stubLoyaltyRepo) UpdateAccount(_ context.Context, a loyalty.LoyaltyAccount) (loyalty.LoyaltyAccount, error) {
	return a, nil
}
func (s *stubLoyaltyRepo) ListAccounts(_ context.Context, _ kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[loyalty.LoyaltyAccount], error) {
	return kernel.Paginated[loyalty.LoyaltyAccount]{
		Items: []loyalty.LoyaltyAccount{{ID: "la-1", TenantID: testTenant, CustomerID: "c-1", PointsBalance: 100, Tier: "bronze", LifetimePoints: 500}},
		Total: 1, Page: pg.Page, TotalPages: 1,
	}, nil
}
func (s *stubLoyaltyRepo) CreateTransaction(_ context.Context, t loyalty.LoyaltyTransaction) (loyalty.LoyaltyTransaction, error) {
	t.ID = "lt-1"
	return t, nil
}
func (s *stubLoyaltyRepo) ListTransactions(_ context.Context, _ kernel.TenantID, _ kernel.LoyaltyAccountID, pg kernel.PaginationOptions) (kernel.Paginated[loyalty.LoyaltyTransaction], error) {
	return kernel.Paginated[loyalty.LoyaltyTransaction]{Items: []loyalty.LoyaltyTransaction{}, Total: 0, Page: pg.Page, TotalPages: 0}, nil
}
func (s *stubLoyaltyRepo) CreateReward(_ context.Context, r loyalty.LoyaltyReward) (loyalty.LoyaltyReward, error) {
	r.ID = "lr-1"
	return r, nil
}
func (s *stubLoyaltyRepo) GetRewardByID(_ context.Context, _ kernel.TenantID, id kernel.RewardID) (loyalty.LoyaltyReward, error) {
	return loyalty.LoyaltyReward{ID: id, TenantID: testTenant, Name: "Free Shipping", PointsCost: 500, RewardType: "free_shipping"}, nil
}
func (s *stubLoyaltyRepo) UpdateReward(_ context.Context, r loyalty.LoyaltyReward) (loyalty.LoyaltyReward, error) {
	return r, nil
}
func (s *stubLoyaltyRepo) ListRewards(_ context.Context, _ kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[loyalty.LoyaltyReward], error) {
	return kernel.Paginated[loyalty.LoyaltyReward]{
		Items: []loyalty.LoyaltyReward{{ID: "lr-1", TenantID: testTenant, Name: "Free Shipping", PointsCost: 500, RewardType: "free_shipping"}},
		Total: 1, Page: pg.Page, TotalPages: 1,
	}, nil
}

// ─── Bundle Stubs ───────────────────────────────────────────────────────────

type stubBundleRepo struct{}

func (s *stubBundleRepo) Create(_ context.Context, b *bundle.Bundle) error {
	b.ID = "bd-1"
	return nil
}
func (s *stubBundleRepo) GetByID(_ context.Context, _ kernel.TenantID, id kernel.BundleID) (*bundle.Bundle, error) {
	return &bundle.Bundle{ID: id, TenantID: testTenant, Name: "Starter Pack", Slug: "starter-pack", DiscountType: "percentage", DiscountValue: 10, Active: true}, nil
}
func (s *stubBundleRepo) GetBySlug(_ context.Context, _ kernel.TenantID, slug string) (*bundle.Bundle, error) {
	return nil, bundle.ErrNotFound
}
func (s *stubBundleRepo) List(_ context.Context, _ kernel.TenantID, _ bool, pg kernel.PaginationOptions) (kernel.Paginated[bundle.Bundle], error) {
	return kernel.Paginated[bundle.Bundle]{
		Items: []bundle.Bundle{{ID: "bd-1", TenantID: testTenant, Name: "Starter Pack", Slug: "starter-pack", DiscountType: "percentage", DiscountValue: 10, Active: true}},
		Total: 1, Page: pg.Page, TotalPages: 1,
	}, nil
}
func (s *stubBundleRepo) Update(_ context.Context, b *bundle.Bundle) error {
	return nil
}
func (s *stubBundleRepo) Delete(_ context.Context, _ kernel.TenantID, _ kernel.BundleID) error {
	return nil
}
func (s *stubBundleRepo) ListItems(_ context.Context, _ kernel.TenantID, _ kernel.BundleID) ([]bundle.BundleItem, error) {
	return []bundle.BundleItem{}, nil
}
func (s *stubBundleRepo) AddItem(_ context.Context, item *bundle.BundleItem) error {
	item.ID = "bi-1"
	return nil
}
func (s *stubBundleRepo) GetItemByID(_ context.Context, _ kernel.TenantID, _ kernel.BundleItemID) (*bundle.BundleItem, error) {
	return &bundle.BundleItem{ID: "bi-1"}, nil
}
func (s *stubBundleRepo) RemoveItem(_ context.Context, _ kernel.TenantID, _ kernel.BundleItemID) error {
	return nil
}

// ─── Dashboard Stubs ────────────────────────────────────────────────────────

type stubDashboardRepo struct{}

func (s *stubDashboardRepo) GetSalesOverview(_ context.Context, _ kernel.TenantID, _ dashboard.DateRange) (dashboard.SalesOverview, error) {
	return dashboard.SalesOverview{TotalRevenue: 150000, OrderCount: 25, AverageOrderValue: 6000}, nil
}
func (s *stubDashboardRepo) GetTopProducts(_ context.Context, _ kernel.TenantID, _ dashboard.DateRange, limit int) ([]dashboard.TopProduct, error) {
	return []dashboard.TopProduct{{ProductID: "p-1", Name: "Widget", Quantity: 50, Revenue: 49500}}, nil
}
func (s *stubDashboardRepo) GetRevenueByDay(_ context.Context, _ kernel.TenantID, _ dashboard.DateRange) ([]dashboard.DailyRevenue, error) {
	return []dashboard.DailyRevenue{{Date: "2024-01-15", Revenue: 5000, OrderCount: 3}}, nil
}
func (s *stubDashboardRepo) GetCustomerStats(_ context.Context, _ kernel.TenantID, _ dashboard.DateRange) (dashboard.CustomerStats, error) {
	return dashboard.CustomerStats{}, nil
}
func (s *stubDashboardRepo) GetConversionFunnel(_ context.Context, _ kernel.TenantID, _ dashboard.DateRange) (dashboard.ConversionFunnel, error) {
	return dashboard.ConversionFunnel{}, nil
}

// ─── Notification Stubs ─────────────────────────────────────────────────────

type stubNotificationRepo struct{}

func (s *stubNotificationRepo) Create(_ context.Context, n *notification.Notification) error {
	n.ID = "nt-1"
	return nil
}
func (s *stubNotificationRepo) GetByID(_ context.Context, _ kernel.TenantID, id kernel.NotificationID) (*notification.Notification, error) {
	return &notification.Notification{ID: id, TenantID: testTenant}, nil
}
func (s *stubNotificationRepo) List(_ context.Context, _ kernel.TenantID, _ kernel.UserID, _ bool, page, pageSize int) (kernel.Paginated[notification.Notification], error) {
	return kernel.Paginated[notification.Notification]{Items: []notification.Notification{}, Total: 0, Page: page, TotalPages: 0}, nil
}
func (s *stubNotificationRepo) MarkRead(_ context.Context, _ kernel.TenantID, _ kernel.NotificationID) error {
	return nil
}
func (s *stubNotificationRepo) MarkAllRead(_ context.Context, _ kernel.TenantID, _ kernel.UserID) error {
	return nil
}
func (s *stubNotificationRepo) GetUnreadCount(_ context.Context, _ kernel.TenantID, _ kernel.UserID) (int, error) {
	return 5, nil
}
func (s *stubNotificationRepo) Delete(_ context.Context, _ kernel.TenantID, _ kernel.NotificationID) error {
	return nil
}

// ═══════════════════════════════════════════════════════════════════════════
// TESTS
// ═══════════════════════════════════════════════════════════════════════════

// ─── Inventory Tests ────────────────────────────────────────────────────────

func newInventoryService() *inventorysrv.Service {
	return inventorysrv.NewService(&stubInventoryRepo{}, eventbus.NewInMemoryBus())
}

func TestListWarehousesTool_Execute(t *testing.T) {
	tool := &ListWarehousesTool{inventory: newInventoryService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "wh-1") {
		t.Errorf("expected warehouse id in result, got: %s", result)
	}
}

func TestCreateWarehouseTool_Execute(t *testing.T) {
	tool := &CreateWarehouseTool{inventory: newInventoryService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"name":"Main Warehouse","address":"123 St"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Created warehouse") {
		t.Errorf("expected 'Created warehouse' in result, got: %s", result)
	}
}

func TestAdjustStockTool_Execute(t *testing.T) {
	tool := &AdjustStockTool{inventory: newInventoryService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"product_id":"p-1","warehouse_id":"wh-1","quantity":10,"type":"received"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Stock adjusted") {
		t.Errorf("expected 'Stock adjusted' in result, got: %s", result)
	}
}

func TestGetLowStockTool_Execute(t *testing.T) {
	tool := &GetLowStockTool{inventory: newInventoryService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "No low stock") {
		t.Errorf("expected 'No low stock' in result, got: %s", result)
	}
}

func TestListWarehousesTool_Execute_InvalidJSON(t *testing.T) {
	tool := &ListWarehousesTool{inventory: newInventoryService(), tenantID: testTenant}
	_, err := tool.Execute(context.Background(), json.RawMessage(`{bad`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// ─── Review Tests ───────────────────────────────────────────────────────────

func newReviewService() *reviewsrv.Service {
	return reviewsrv.New(&stubReviewRepo{}, eventbus.NewInMemoryBus())
}

func TestListReviewsTool_Execute(t *testing.T) {
	tool := &ListReviewsTool{reviews: newReviewService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "rv-1") {
		t.Errorf("expected review id in result, got: %s", result)
	}
}

func TestApproveReviewTool_Execute(t *testing.T) {
	tool := &ApproveReviewTool{reviews: newReviewService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"review_id":"rv-1"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "approved") {
		t.Errorf("expected 'approved' in result, got: %s", result)
	}
}

func TestRejectReviewTool_Execute(t *testing.T) {
	tool := &RejectReviewTool{reviews: newReviewService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"review_id":"rv-1"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "rejected") {
		t.Errorf("expected 'rejected' in result, got: %s", result)
	}
}

func TestListReviewsTool_Execute_InvalidJSON(t *testing.T) {
	tool := &ListReviewsTool{reviews: newReviewService(), tenantID: testTenant}
	_, err := tool.Execute(context.Background(), json.RawMessage(`{bad`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// ─── Returns Tests ──────────────────────────────────────────────────────────

func newReturnsService() *returnssrv.Service {
	return returnssrv.New(&stubReturnsRepo{}, eventbus.NewInMemoryBus())
}

func TestListReturnsTool_Execute(t *testing.T) {
	tool := &ListReturnsTool{returns: newReturnsService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "ret-1") {
		t.Errorf("expected return id in result, got: %s", result)
	}
}

func TestApproveReturnTool_Execute(t *testing.T) {
	tool := &ApproveReturnTool{returns: newReturnsService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"return_id":"ret-1"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "approved") || !strings.Contains(result, "ret-1") {
		t.Errorf("expected approved return in result, got: %s", result)
	}
}

func TestListReturnsTool_Execute_InvalidJSON(t *testing.T) {
	tool := &ListReturnsTool{returns: newReturnsService(), tenantID: testTenant}
	_, err := tool.Execute(context.Background(), json.RawMessage(`{bad`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// ─── Webhook Tests ──────────────────────────────────────────────────────────

func newWebhookService() *webhooksrv.Service {
	return webhooksrv.New(&stubWebhookRepo{})
}

func TestListWebhooksTool_Execute(t *testing.T) {
	tool := &ListWebhooksTool{webhooks: newWebhookService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "example.com") {
		t.Errorf("expected webhook URL in result, got: %s", result)
	}
}

func TestCreateWebhookTool_Execute(t *testing.T) {
	tool := &CreateWebhookTool{webhooks: newWebhookService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"url":"https://example.com/hook","events":["order.created"]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Created webhook") {
		t.Errorf("expected 'Created webhook' in result, got: %s", result)
	}
}

func TestToggleWebhookTool_Execute(t *testing.T) {
	tool := &ToggleWebhookTool{webhooks: newWebhookService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"webhook_id":"wh-1","active":false}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "active=false") {
		t.Errorf("expected 'active=false' in result, got: %s", result)
	}
}

func TestCreateWebhookTool_Execute_InvalidJSON(t *testing.T) {
	tool := &CreateWebhookTool{webhooks: newWebhookService(), tenantID: testTenant}
	_, err := tool.Execute(context.Background(), json.RawMessage(`{bad`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// ─── Audit Tests ────────────────────────────────────────────────────────────

func newAuditService() *auditsrv.Service {
	return auditsrv.New(&stubAuditRepo{})
}

func TestListAuditLogsTool_Execute(t *testing.T) {
	tool := &ListAuditLogsTool{audit: newAuditService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "au-1") || !strings.Contains(result, "product.create") {
		t.Errorf("expected audit entry in result, got: %s", result)
	}
}

func TestListAuditLogsTool_Execute_InvalidJSON(t *testing.T) {
	tool := &ListAuditLogsTool{audit: newAuditService(), tenantID: testTenant}
	_, err := tool.Execute(context.Background(), json.RawMessage(`{bad`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// ─── Loyalty Tests ──────────────────────────────────────────────────────────

func newLoyaltyService() *loyaltysrv.Service {
	return loyaltysrv.New(&stubLoyaltyRepo{}, eventbus.NewInMemoryBus())
}

func TestListLoyaltyAccountsTool_Execute(t *testing.T) {
	tool := &ListLoyaltyAccountsTool{loyalty: newLoyaltyService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "la-1") {
		t.Errorf("expected loyalty account id in result, got: %s", result)
	}
}

func TestEarnLoyaltyPointsTool_Execute(t *testing.T) {
	tool := &EarnLoyaltyPointsTool{loyalty: newLoyaltyService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"CustomerID":"c-1","Points":50,"Reference":"order-123"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Awarded") || !strings.Contains(result, "50") {
		t.Errorf("expected 'Awarded 50' in result, got: %s", result)
	}
}

func TestListLoyaltyRewardsTool_Execute(t *testing.T) {
	tool := &ListLoyaltyRewardsTool{loyalty: newLoyaltyService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "lr-1") || !strings.Contains(result, "Free Shipping") {
		t.Errorf("expected reward in result, got: %s", result)
	}
}

func TestCreateLoyaltyRewardTool_Execute(t *testing.T) {
	tool := &CreateLoyaltyRewardTool{loyalty: newLoyaltyService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"Name":"10% Off","PointsCost":200,"RewardType":"discount","ValueCents":1000}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Created reward") {
		t.Errorf("expected 'Created reward' in result, got: %s", result)
	}
}

func TestEarnLoyaltyPointsTool_Execute_InvalidJSON(t *testing.T) {
	tool := &EarnLoyaltyPointsTool{loyalty: newLoyaltyService(), tenantID: testTenant}
	_, err := tool.Execute(context.Background(), json.RawMessage(`{bad`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// ─── Bundle Tests ───────────────────────────────────────────────────────────

func newBundleService() *bundlesrv.Service {
	return bundlesrv.New(&stubBundleRepo{}, eventbus.NewInMemoryBus())
}

func TestListBundlesTool_Execute(t *testing.T) {
	tool := &ListBundlesTool{bundles: newBundleService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "bd-1") || !strings.Contains(result, "Starter Pack") {
		t.Errorf("expected bundle in result, got: %s", result)
	}
}

func TestCreateBundleTool_Execute(t *testing.T) {
	tool := &CreateBundleTool{bundles: newBundleService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"Name":"Summer Pack","Slug":"summer-pack","DiscountType":"percentage","DiscountValue":15,"Active":true}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Created bundle") {
		t.Errorf("expected 'Created bundle' in result, got: %s", result)
	}
}

func TestCreateBundleTool_Execute_InvalidJSON(t *testing.T) {
	tool := &CreateBundleTool{bundles: newBundleService(), tenantID: testTenant}
	_, err := tool.Execute(context.Background(), json.RawMessage(`{bad`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// ─── Dashboard Tests ────────────────────────────────────────────────────────

func newDashboardService() *dashboardsrv.Service {
	return dashboardsrv.New(&stubDashboardRepo{})
}

func TestGetSalesOverviewTool_Execute(t *testing.T) {
	tool := &GetSalesOverviewTool{dashboard: newDashboardService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"days":30}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Revenue") || !strings.Contains(result, "Orders") {
		t.Errorf("expected sales overview in result, got: %s", result)
	}
}

func TestGetTopProductsTool_Execute(t *testing.T) {
	tool := &GetTopProductsTool{dashboard: newDashboardService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"days":30,"limit":5}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Widget") {
		t.Errorf("expected product name in result, got: %s", result)
	}
}

func TestGetRevenueByDayTool_Execute(t *testing.T) {
	tool := &GetRevenueByDayTool{dashboard: newDashboardService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"days":7}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Daily revenue") {
		t.Errorf("expected 'Daily revenue' in result, got: %s", result)
	}
}

func TestGetSalesOverviewTool_Execute_InvalidJSON(t *testing.T) {
	tool := &GetSalesOverviewTool{dashboard: newDashboardService(), tenantID: testTenant}
	_, err := tool.Execute(context.Background(), json.RawMessage(`{bad`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// ─── Notification Tests ─────────────────────────────────────────────────────

func newNotificationService() *notificationsrv.Service {
	return notificationsrv.New(&stubNotificationRepo{})
}

func TestGetUnreadNotificationCountTool_Execute(t *testing.T) {
	tool := &GetUnreadNotificationCountTool{notifications: newNotificationService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"user_id":"u-1"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "5") {
		t.Errorf("expected '5' unread count in result, got: %s", result)
	}
}

func TestMarkAllNotificationsReadTool_Execute(t *testing.T) {
	tool := &MarkAllNotificationsReadTool{notifications: newNotificationService(), tenantID: testTenant}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"user_id":"u-1"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "marked as read") {
		t.Errorf("expected 'marked as read' in result, got: %s", result)
	}
}

func TestGetUnreadNotificationCountTool_Execute_InvalidJSON(t *testing.T) {
	tool := &GetUnreadNotificationCountTool{notifications: newNotificationService(), tenantID: testTenant}
	_, err := tool.Execute(context.Background(), json.RawMessage(`{bad`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
