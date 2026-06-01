package agent

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/Abraxas-365/vendex/internal/catalog"
	"github.com/Abraxas-365/vendex/internal/catalog/catalogsrv"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/order"
	"github.com/Abraxas-365/vendex/internal/order/ordersrv"
	"github.com/Abraxas-365/vendex/internal/product"
	"github.com/Abraxas-365/vendex/internal/product/productsrv"
	"github.com/Abraxas-365/vendex/internal/promo"
	"github.com/Abraxas-365/vendex/internal/promo/promosrv"
	"github.com/Abraxas-365/vendex/internal/storefront"
	"github.com/Abraxas-365/vendex/internal/storefront/storefrontsrv"
	"github.com/Abraxas-365/vendex/internal/theme"
	"github.com/Abraxas-365/vendex/internal/theme/themesrv"
)

const testTenant = kernel.TenantID("test-tenant")

// ─── Mock: storefront.PageRepository ─────────────────────────────────────────

type stubPageRepo struct{}

func (r *stubPageRepo) Create(_ context.Context, p *storefront.Page) error { return nil }
func (r *stubPageRepo) GetByID(_ context.Context, _ kernel.TenantID, id kernel.PageID) (*storefront.Page, error) {
	return &storefront.Page{ID: id, TenantID: testTenant, Slug: "test", Title: "Test", Status: storefront.PageStatusDraft, Version: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}
func (r *stubPageRepo) GetBySlug(_ context.Context, _ kernel.TenantID, _ string) (*storefront.Page, error) {
	return nil, nil
}
func (r *stubPageRepo) GetPublished(_ context.Context, _ kernel.TenantID, _ string) (*storefront.Page, error) {
	return nil, nil
}
func (r *stubPageRepo) Update(_ context.Context, _ *storefront.Page) error { return nil }
func (r *stubPageRepo) List(_ context.Context, _ kernel.TenantID, _ kernel.PaginationOptions) (kernel.Paginated[storefront.Page], error) {
	return kernel.Paginated[storefront.Page]{
		Items:      []storefront.Page{{ID: "pg-1", TenantID: testTenant, Slug: "home", Title: "Home", Status: storefront.PageStatusDraft, Version: 1}},
		Total:      1,
		Page:       1,
		TotalPages: 1,
	}, nil
}
func (r *stubPageRepo) ListByStatus(_ context.Context, _ kernel.TenantID, _ storefront.PageStatus, _ kernel.PaginationOptions) (kernel.Paginated[storefront.Page], error) {
	return kernel.Paginated[storefront.Page]{
		Items:      []storefront.Page{{ID: "pg-1", TenantID: testTenant, Slug: "home", Title: "Home", Status: storefront.PageStatusDraft, Version: 1}},
		Total:      1,
		Page:       1,
		TotalPages: 1,
	}, nil
}
func (r *stubPageRepo) Delete(_ context.Context, _ kernel.TenantID, _ kernel.PageID) error {
	return nil
}

// ─── Mock: storefront.PageVersionRepository ──────────────────────────────────

type stubVersionRepo struct{}

func (r *stubVersionRepo) Create(_ context.Context, _ *storefront.PageVersion) error { return nil }
func (r *stubVersionRepo) GetByVersion(_ context.Context, _ kernel.TenantID, _ kernel.PageID, _ int) (*storefront.PageVersion, error) {
	return nil, nil
}
func (r *stubVersionRepo) ListByPage(_ context.Context, _ kernel.TenantID, _ kernel.PageID) ([]storefront.PageVersion, error) {
	return nil, nil
}

// ─── Mock: storefront.BlockTypeRepository ────────────────────────────────────

type stubBlockTypeRepo struct{}

func (r *stubBlockTypeRepo) Create(_ context.Context, _ *storefront.BlockType) error { return nil }
func (r *stubBlockTypeRepo) GetByID(_ context.Context, _ kernel.BlockTypeID) (*storefront.BlockType, error) {
	return nil, nil
}
func (r *stubBlockTypeRepo) GetByName(_ context.Context, _ string) (*storefront.BlockType, error) {
	return nil, nil
}
func (r *stubBlockTypeRepo) List(_ context.Context, _ string) ([]storefront.BlockType, error) {
	return []storefront.BlockType{
		{ID: "bt-1", Name: "hero", DisplayName: "Hero Banner", Category: "content"},
	}, nil
}
func (r *stubBlockTypeRepo) Update(_ context.Context, _ *storefront.BlockType) error { return nil }
func (r *stubBlockTypeRepo) Delete(_ context.Context, _ kernel.BlockTypeID) error    { return nil }

// ─── Mock: theme.ThemeRepository ─────────────────────────────────────────────

type stubThemeRepo struct {
	themes map[kernel.ThemeID]*theme.Theme
}

func newStubThemeRepo() *stubThemeRepo {
	return &stubThemeRepo{themes: map[kernel.ThemeID]*theme.Theme{}}
}

func (r *stubThemeRepo) Create(_ context.Context, t *theme.Theme) error {
	r.themes[t.ID] = t
	return nil
}
func (r *stubThemeRepo) GetByID(_ context.Context, _ kernel.TenantID, id kernel.ThemeID) (*theme.Theme, error) {
	if t, ok := r.themes[id]; ok {
		return t, nil
	}
	return nil, theme.ErrThemeNotFound
}
func (r *stubThemeRepo) GetActive(_ context.Context, _ kernel.TenantID) (*theme.Theme, error) {
	for _, t := range r.themes {
		if t.IsActive {
			return t, nil
		}
	}
	return nil, theme.ErrThemeNotFound
}
func (r *stubThemeRepo) List(_ context.Context, _ kernel.TenantID) ([]theme.Theme, error) {
	var out []theme.Theme
	for _, t := range r.themes {
		out = append(out, *t)
	}
	return out, nil
}
func (r *stubThemeRepo) Update(_ context.Context, t *theme.Theme) error {
	r.themes[t.ID] = t
	return nil
}
func (r *stubThemeRepo) Delete(_ context.Context, _ kernel.TenantID, id kernel.ThemeID) error {
	delete(r.themes, id)
	return nil
}
func (r *stubThemeRepo) DeactivateAll(_ context.Context, _ kernel.TenantID) error {
	for _, t := range r.themes {
		t.IsActive = false
	}
	return nil
}

// ─── Mock: product.Repository ────────────────────────────────────────────────

type stubProductRepo struct{}

func (r *stubProductRepo) Create(_ context.Context, _ *product.Product) error { return nil }
func (r *stubProductRepo) GetByID(_ context.Context, _ kernel.TenantID, _ kernel.ProductID) (*product.Product, error) {
	return nil, product.ErrNotFound
}
func (r *stubProductRepo) GetBySKU(_ context.Context, _ kernel.TenantID, _ string) (*product.Product, error) {
	return nil, product.ErrNotFound
}
func (r *stubProductRepo) GetBySlug(_ context.Context, _ kernel.TenantID, _ string) (*product.Product, error) {
	return nil, product.ErrNotFound
}
func (r *stubProductRepo) Update(_ context.Context, _ *product.Product) error { return nil }
func (r *stubProductRepo) Delete(_ context.Context, _ kernel.TenantID, _ kernel.ProductID) error {
	return nil
}
func (r *stubProductRepo) List(_ context.Context, _ kernel.TenantID, _ kernel.PaginationOptions) (kernel.Paginated[product.Product], error) {
	return kernel.Paginated[product.Product]{
		Items: []product.Product{
			{ID: "p-1", TenantID: testTenant, Name: "Widget", SKU: "WGT-1", Price: kernel.Money{Amount: 999, Currency: "USD"}, Stock: 10, Status: product.StatusDraft},
		},
		Total: 1, Page: 1, TotalPages: 1,
	}, nil
}
func (r *stubProductRepo) ListByCategory(_ context.Context, _ kernel.TenantID, _ kernel.CategoryID, _ kernel.PaginationOptions) (kernel.Paginated[product.Product], error) {
	return kernel.Paginated[product.Product]{Items: []product.Product{}, Total: 0, Page: 1, TotalPages: 0}, nil
}

// ─── Mock: product.VariantRepository ─────────────────────────────────────────

type stubVariantRepo struct{}

func (r *stubVariantRepo) CreateOption(_ context.Context, _ *product.ProductOption) error {
	return nil
}
func (r *stubVariantRepo) ListOptions(_ context.Context, _ kernel.TenantID, _ kernel.ProductID) ([]product.ProductOption, error) {
	return []product.ProductOption{}, nil
}
func (r *stubVariantRepo) UpdateOption(_ context.Context, _ *product.ProductOption) error {
	return nil
}
func (r *stubVariantRepo) DeleteOption(_ context.Context, _ kernel.TenantID, _ kernel.OptionID) error {
	return nil
}
func (r *stubVariantRepo) CreateVariant(_ context.Context, _ *product.ProductVariant) error {
	return nil
}
func (r *stubVariantRepo) GetVariantByID(_ context.Context, _ kernel.TenantID, _ kernel.VariantID) (*product.ProductVariant, error) {
	return nil, product.ErrVariantNotFound
}
func (r *stubVariantRepo) ListVariants(_ context.Context, _ kernel.TenantID, _ kernel.ProductID) ([]product.ProductVariant, error) {
	return []product.ProductVariant{}, nil
}
func (r *stubVariantRepo) UpdateVariant(_ context.Context, _ *product.ProductVariant) error {
	return nil
}
func (r *stubVariantRepo) DeleteVariant(_ context.Context, _ kernel.TenantID, _ kernel.VariantID) error {
	return nil
}
func (r *stubVariantRepo) GetVariantBySKU(_ context.Context, _ kernel.TenantID, _ string) (*product.ProductVariant, error) {
	return nil, product.ErrVariantNotFound
}

// ─── Mock: promo.PromoRepository ─────────────────────────────────────────────

type stubPromoRepo struct{}

func (r *stubPromoRepo) Create(_ context.Context, p *promo.Promo) error { return nil }
func (r *stubPromoRepo) GetByID(_ context.Context, _ kernel.TenantID, _ kernel.PromoID) (*promo.Promo, error) {
	return nil, promo.ErrPromoNotFound
}
func (r *stubPromoRepo) GetByCode(_ context.Context, _ kernel.TenantID, _ string) (*promo.Promo, error) {
	return nil, promo.ErrPromoNotFound
}
func (r *stubPromoRepo) List(_ context.Context, _ kernel.TenantID, _ kernel.PaginationOptions) (kernel.Paginated[promo.Promo], error) {
	return kernel.Paginated[promo.Promo]{Items: []promo.Promo{}, Total: 0, Page: 1, TotalPages: 0}, nil
}
func (r *stubPromoRepo) Update(_ context.Context, _ *promo.Promo) error { return nil }
func (r *stubPromoRepo) Delete(_ context.Context, _ kernel.TenantID, _ kernel.PromoID) error {
	return nil
}
func (r *stubPromoRepo) IncrementUsedCount(_ context.Context, _ kernel.TenantID, _ kernel.PromoID) error {
	return nil
}

// ─── Mock: order.Repository ──────────────────────────────────────────────────

type stubOrderRepo struct{}

func (r *stubOrderRepo) Create(_ context.Context, _ *order.Order) error { return nil }
func (r *stubOrderRepo) GetByID(_ context.Context, _ kernel.TenantID, _ kernel.OrderID) (*order.Order, error) {
	return nil, order.ErrNotFound
}
func (r *stubOrderRepo) Update(_ context.Context, _ *order.Order) error { return nil }
func (r *stubOrderRepo) List(_ context.Context, _ kernel.TenantID, _ kernel.PaginationOptions) (kernel.Paginated[order.Order], error) {
	now := time.Now()
	return kernel.Paginated[order.Order]{
		Items: []order.Order{
			{
				ID: "ord-1", TenantID: testTenant, CustomerID: "cust-1",
				Status:      order.StatusPending,
				Items:       []order.OrderItem{{ID: "oi-1", ProductID: "p-1", ProductName: "Widget", Quantity: 2, UnitPrice: kernel.Money{Amount: 999, Currency: "USD"}}},
				TotalAmount: kernel.Money{Amount: 1998, Currency: "USD"},
				CreatedAt:   now, UpdatedAt: now,
			},
		},
		Total: 1, Page: 1, TotalPages: 1,
	}, nil
}
func (r *stubOrderRepo) UpdateCheckoutFields(_ context.Context, _ *order.Order) error {
	return nil
}
func (r *stubOrderRepo) ListByCustomer(_ context.Context, _ kernel.TenantID, _ kernel.CustomerID, _ kernel.PaginationOptions) (kernel.Paginated[order.Order], error) {
	return kernel.Paginated[order.Order]{Items: []order.Order{}, Total: 0, Page: 1, TotalPages: 0}, nil
}

// ─── Mock: catalog.CategoryRepository ────────────────────────────────────────

type stubCategoryRepo struct{}

func (r *stubCategoryRepo) Create(_ context.Context, _ *catalog.Category) error { return nil }
func (r *stubCategoryRepo) GetByID(_ context.Context, _ kernel.TenantID, _ kernel.CategoryID) (*catalog.Category, error) {
	return nil, catalog.ErrCategoryNotFound
}
func (r *stubCategoryRepo) GetBySlug(_ context.Context, _ kernel.TenantID, _ string) (*catalog.Category, error) {
	return nil, catalog.ErrCategoryNotFound
}
func (r *stubCategoryRepo) List(_ context.Context, _ kernel.TenantID, _ kernel.PaginationOptions) (kernel.Paginated[catalog.Category], error) {
	return kernel.Paginated[catalog.Category]{
		Items: []catalog.Category{
			{ID: "cat-1", TenantID: testTenant, Name: "Electronics", Slug: "electronics"},
		},
		Total: 1, Page: 1, TotalPages: 1,
	}, nil
}
func (r *stubCategoryRepo) Update(_ context.Context, _ *catalog.Category) error { return nil }
func (r *stubCategoryRepo) Delete(_ context.Context, _ kernel.TenantID, _ kernel.CategoryID) error {
	return nil
}
func (r *stubCategoryRepo) ListByParent(_ context.Context, _ kernel.TenantID, _ *kernel.CategoryID, _ kernel.PaginationOptions) (kernel.Paginated[catalog.Category], error) {
	return kernel.Paginated[catalog.Category]{Items: []catalog.Category{}, Total: 0, Page: 1, TotalPages: 0}, nil
}

// ─── Mock: catalog.CollectionRepository ──────────────────────────────────────

type stubCollectionRepo struct{}

func (r *stubCollectionRepo) Create(_ context.Context, _ *catalog.Collection) error { return nil }
func (r *stubCollectionRepo) GetByID(_ context.Context, _ kernel.TenantID, _ kernel.CollectionID) (*catalog.Collection, error) {
	return nil, catalog.ErrCollectionNotFound
}
func (r *stubCollectionRepo) GetBySlug(_ context.Context, _ kernel.TenantID, _ string) (*catalog.Collection, error) {
	return nil, catalog.ErrCollectionNotFound
}
func (r *stubCollectionRepo) List(_ context.Context, _ kernel.TenantID, _ kernel.PaginationOptions) (kernel.Paginated[catalog.Collection], error) {
	return kernel.Paginated[catalog.Collection]{
		Items: []catalog.Collection{
			{ID: "col-1", TenantID: testTenant, Name: "Summer Sale", Slug: "summer-sale"},
		},
		Total: 1, Page: 1, TotalPages: 1,
	}, nil
}
func (r *stubCollectionRepo) Update(_ context.Context, _ *catalog.Collection) error { return nil }
func (r *stubCollectionRepo) Delete(_ context.Context, _ kernel.TenantID, _ kernel.CollectionID) error {
	return nil
}

// ─── Helper builders ─────────────────────────────────────────────────────────

func newStorefrontSvc() *storefrontsrv.Service {
	return storefrontsrv.New(&stubPageRepo{}, &stubVersionRepo{}, &stubBlockTypeRepo{}, eventbus.NewInMemoryBus())
}

func newThemeSvc() (*themesrv.Service, *stubThemeRepo) {
	repo := newStubThemeRepo()
	return themesrv.New(repo, eventbus.NewInMemoryBus()), repo
}

func newProductSvc() *productsrv.Service {
	return productsrv.New(&stubProductRepo{}, &stubVariantRepo{}, eventbus.NewInMemoryBus())
}

func newPromoSvc() *promosrv.Service {
	return promosrv.New(&stubPromoRepo{})
}

func newOrderSvc() *ordersrv.Service {
	return ordersrv.New(&stubOrderRepo{}, eventbus.NewInMemoryBus())
}

func newCatalogSvc() *catalogsrv.Service {
	return catalogsrv.New(&stubCategoryRepo{}, &stubCollectionRepo{}, eventbus.NewInMemoryBus())
}

// ═══════════════════════════════════════════════════════════════════════════════
// Tests for tools.go — Storefront tools
// ═══════════════════════════════════════════════════════════════════════════════

func TestCreatePageTool_Execute(t *testing.T) {
	svc := newStorefrontSvc()
	tool := &CreatePageTool{sf: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"slug":"about","title":"About Us","html":"<h1>About</h1>"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "pending_review") {
		t.Errorf("expected pending_review in result, got: %s", result)
	}
	if !strings.Contains(result, "About Us") {
		t.Errorf("expected title in result, got: %s", result)
	}
}

func TestCreatePageTool_Execute_InvalidJSON(t *testing.T) {
	svc := newStorefrontSvc()
	tool := &CreatePageTool{sf: svc, tenantID: testTenant}
	_, err := tool.Execute(context.Background(), json.RawMessage(`{bad`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestUpdatePageTool_Execute(t *testing.T) {
	svc := newStorefrontSvc()
	tool := &UpdatePageTool{sf: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"page_id":"pg-1","html":"<h1>Updated</h1>"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "pg-1") {
		t.Errorf("expected page ID in result, got: %s", result)
	}
}

func TestListPagesTool_Execute(t *testing.T) {
	svc := newStorefrontSvc()
	tool := &ListPagesTool{sf: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "home") {
		t.Errorf("expected page slug in result, got: %s", result)
	}
}

func TestListPagesTool_Execute_WithStatus(t *testing.T) {
	svc := newStorefrontSvc()
	tool := &ListPagesTool{sf: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"status":"draft"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "home") {
		t.Errorf("expected page slug in result, got: %s", result)
	}
}

func TestListBlockTypesTool_Execute(t *testing.T) {
	svc := newStorefrontSvc()
	tool := &ListBlockTypesTool{sf: svc}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "hero") {
		t.Errorf("expected block type name in result, got: %s", result)
	}
}

func TestCreateBlockTypeTool_Execute(t *testing.T) {
	svc := newStorefrontSvc()
	tool := &CreateBlockTypeTool{sf: svc}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"name":"banner","display_name":"Banner","category":"content"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "banner") {
		t.Errorf("expected block type name in result, got: %s", result)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Tests for tools.go — Theme tools
// ═══════════════════════════════════════════════════════════════════════════════

func TestListThemesTool_Execute_Empty(t *testing.T) {
	svc, _ := newThemeSvc()
	tool := &ListThemesTool{themes: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "No themes") {
		t.Errorf("expected empty message, got: %s", result)
	}
}

func TestGetActiveThemeTool_Execute_CreatesDefault(t *testing.T) {
	svc, _ := newThemeSvc()
	tool := &GetActiveThemeTool{themes: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Default") {
		t.Errorf("expected default theme name in result, got: %s", result)
	}
}

func TestCreateThemeTool_Execute(t *testing.T) {
	svc, _ := newThemeSvc()
	tool := &CreateThemeTool{themes: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"name":"Dark Mode"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Dark Mode") {
		t.Errorf("expected theme name in result, got: %s", result)
	}
}

func TestCreateThemeTool_Execute_EmptyName(t *testing.T) {
	svc, _ := newThemeSvc()
	tool := &CreateThemeTool{themes: svc, tenantID: testTenant}

	_, err := tool.Execute(context.Background(), json.RawMessage(`{"name":""}`))
	if err == nil {
		t.Fatal("expected error for empty theme name")
	}
}

func TestUpdateThemeTool_Execute(t *testing.T) {
	svc, repo := newThemeSvc()
	// Pre-seed a theme.
	repo.themes["th-1"] = &theme.Theme{
		ID: "th-1", TenantID: testTenant, Name: "Light",
		Tokens: theme.DefaultTokens(), CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	tool := &UpdateThemeTool{themes: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"id":"th-1","name":"Light v2"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Light v2") {
		t.Errorf("expected updated name in result, got: %s", result)
	}
}

func TestActivateThemeTool_Execute(t *testing.T) {
	svc, repo := newThemeSvc()
	repo.themes["th-1"] = &theme.Theme{
		ID: "th-1", TenantID: testTenant, Name: "Dark", IsActive: false,
		Tokens: theme.DefaultTokens(), CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	tool := &ActivateThemeTool{themes: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"id":"th-1"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "th-1") {
		t.Errorf("expected theme ID in result, got: %s", result)
	}
}

func TestActivateThemeTool_Execute_NotFound(t *testing.T) {
	svc, _ := newThemeSvc()
	tool := &ActivateThemeTool{themes: svc, tenantID: testTenant}

	_, err := tool.Execute(context.Background(), json.RawMessage(`{"theme_id":"nonexistent"}`))
	if err == nil {
		t.Fatal("expected error for nonexistent theme")
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Tests for tools.go — Product tools
// ═══════════════════════════════════════════════════════════════════════════════

func TestCreateProductTool_Execute(t *testing.T) {
	svc := newProductSvc()
	tool := &CreateProductTool{products: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"name":"Widget","sku":"WGT-1","price_cents":999,"currency":"USD","stock":10}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Widget") {
		t.Errorf("expected product name in result, got: %s", result)
	}
	if !strings.Contains(result, "999") {
		t.Errorf("expected price in result, got: %s", result)
	}
}

func TestCreateProductTool_Execute_ZeroPrice(t *testing.T) {
	svc := newProductSvc()
	tool := &CreateProductTool{products: svc, tenantID: testTenant}

	_, err := tool.Execute(context.Background(), json.RawMessage(`{"name":"Free","price_cents":0,"currency":"USD"}`))
	if err == nil {
		t.Fatal("expected error for zero price")
	}
}

func TestListProductsTool_Execute(t *testing.T) {
	svc := newProductSvc()
	tool := &ListProductsTool{products: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Widget") {
		t.Errorf("expected product name in result, got: %s", result)
	}
	if !strings.Contains(result, "WGT-1") {
		t.Errorf("expected SKU in result, got: %s", result)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Tests for tools.go — Promo tools
// ═══════════════════════════════════════════════════════════════════════════════

func TestCreatePromoTool_Execute(t *testing.T) {
	svc := newPromoSvc()
	tool := &CreatePromoTool{promos: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"code":"SUMMER20","type":"percentage","value":20}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "SUMMER20") {
		t.Errorf("expected promo code in result, got: %s", result)
	}
	if !strings.Contains(result, "percentage") {
		t.Errorf("expected promo type in result, got: %s", result)
	}
}

func TestCreatePromoTool_Execute_InvalidJSON(t *testing.T) {
	svc := newPromoSvc()
	tool := &CreatePromoTool{promos: svc, tenantID: testTenant}
	_, err := tool.Execute(context.Background(), json.RawMessage(`{bad`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Tests for tools.go — Order tools
// ═══════════════════════════════════════════════════════════════════════════════

func TestQueryOrdersTool_Execute(t *testing.T) {
	svc := newOrderSvc()
	tool := &QueryOrdersTool{orders: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "ord-1") {
		t.Errorf("expected order ID in result, got: %s", result)
	}
	if !strings.Contains(result, "pending") {
		t.Errorf("expected order status in result, got: %s", result)
	}
}

func TestQueryOrdersTool_Execute_InvalidJSON(t *testing.T) {
	svc := newOrderSvc()
	tool := &QueryOrdersTool{orders: svc, tenantID: testTenant}
	_, err := tool.Execute(context.Background(), json.RawMessage(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Tests for tools.go — Catalog tools
// ═══════════════════════════════════════════════════════════════════════════════

func TestSearchCatalogTool_Execute_Categories(t *testing.T) {
	svc := newCatalogSvc()
	tool := &SearchCatalogTool{catalog: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"target":"categories"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Electronics") {
		t.Errorf("expected category name in result, got: %s", result)
	}
}

func TestSearchCatalogTool_Execute_Collections(t *testing.T) {
	svc := newCatalogSvc()
	tool := &SearchCatalogTool{catalog: svc, tenantID: testTenant}

	result, err := tool.Execute(context.Background(), json.RawMessage(`{"target":"collections"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Summer Sale") {
		t.Errorf("expected collection name in result, got: %s", result)
	}
}

func TestSearchCatalogTool_Execute_InvalidTarget(t *testing.T) {
	svc := newCatalogSvc()
	tool := &SearchCatalogTool{catalog: svc, tenantID: testTenant}

	_, err := tool.Execute(context.Background(), json.RawMessage(`{"target":"invalid"}`))
	if err == nil {
		t.Fatal("expected error for invalid target")
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Tests for Name() and Description() on all tools.go tools
// ═══════════════════════════════════════════════════════════════════════════════

func TestToolsGo_NameAndDescription(t *testing.T) {
	sfSvc := newStorefrontSvc()
	themeSvc, _ := newThemeSvc()
	productSvc := newProductSvc()
	promoSvc := newPromoSvc()
	orderSvc := newOrderSvc()
	catalogSvc := newCatalogSvc()

	tools := []Tool{
		&CreatePageTool{sf: sfSvc, tenantID: testTenant},
		&UpdatePageTool{sf: sfSvc, tenantID: testTenant},
		&ListPagesTool{sf: sfSvc, tenantID: testTenant},
		&ListBlockTypesTool{sf: sfSvc},
		&CreateBlockTypeTool{sf: sfSvc},
		&ListThemesTool{themes: themeSvc, tenantID: testTenant},
		&GetActiveThemeTool{themes: themeSvc, tenantID: testTenant},
		&CreateThemeTool{themes: themeSvc, tenantID: testTenant},
		&UpdateThemeTool{themes: themeSvc, tenantID: testTenant},
		&ActivateThemeTool{themes: themeSvc, tenantID: testTenant},
		&CreateProductTool{products: productSvc, tenantID: testTenant},
		&ListProductsTool{products: productSvc, tenantID: testTenant},
		&CreatePromoTool{promos: promoSvc, tenantID: testTenant},
		&QueryOrdersTool{orders: orderSvc, tenantID: testTenant},
		&SearchCatalogTool{catalog: catalogSvc, tenantID: testTenant},
	}

	for _, tool := range tools {
		t.Run(tool.Name(), func(t *testing.T) {
			if tool.Name() == "" {
				t.Error("Name() should not be empty")
			}
			if tool.Description() == "" {
				t.Error("Description() should not be empty")
			}
			schema := tool.InputSchema()
			if schema == nil {
				t.Error("InputSchema() should not be nil")
			}
		})
	}
}
