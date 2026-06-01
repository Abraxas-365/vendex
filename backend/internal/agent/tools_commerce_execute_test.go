package agent

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/Abraxas-365/vendex/internal/cartrecovery"
	"github.com/Abraxas-365/vendex/internal/cartrecovery/cartrecoverysrv"
	"github.com/Abraxas-365/vendex/internal/currency"
	"github.com/Abraxas-365/vendex/internal/currency/currencysrv"
	"github.com/Abraxas-365/vendex/internal/customergroup"
	"github.com/Abraxas-365/vendex/internal/customergroup/customergroupsrv"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/giftcard"
	"github.com/Abraxas-365/vendex/internal/giftcard/giftcardsrv"
	"github.com/Abraxas-365/vendex/internal/i18n"
	"github.com/Abraxas-365/vendex/internal/i18n/i18nsrv"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/payment"
	"github.com/Abraxas-365/vendex/internal/payment/paymentsrv"
	"github.com/Abraxas-365/vendex/internal/product/productsrv"
	"github.com/Abraxas-365/vendex/internal/search"
	"github.com/Abraxas-365/vendex/internal/search/searchsrv"
	"github.com/Abraxas-365/vendex/internal/shipping"
	"github.com/Abraxas-365/vendex/internal/shipping/shippingsrv"
	"github.com/Abraxas-365/vendex/internal/subscription"
	"github.com/Abraxas-365/vendex/internal/subscription/subscriptionsrv"
	"github.com/Abraxas-365/vendex/internal/tax"
	"github.com/Abraxas-365/vendex/internal/tax/taxsrv"
)

// ---------------------------------------------------------------------------
// Stub: shipping.ZoneRepository
// ---------------------------------------------------------------------------

type stubZoneRepo struct{}

func (s *stubZoneRepo) Create(_ context.Context, z *shipping.ShippingZone) error { return nil }
func (s *stubZoneRepo) GetByID(_ context.Context, _ kernel.TenantID, _ kernel.ShippingZoneID) (*shipping.ShippingZone, error) {
	return &shipping.ShippingZone{ID: "zone-1", TenantID: "test-tenant", Name: "US"}, nil
}
func (s *stubZoneRepo) List(_ context.Context, _ kernel.TenantID) ([]shipping.ShippingZone, error) {
	return []shipping.ShippingZone{{ID: "zone-1", TenantID: "test-tenant", Name: "US", Countries: []string{"US"}}}, nil
}
func (s *stubZoneRepo) Update(_ context.Context, _ *shipping.ShippingZone) error { return nil }
func (s *stubZoneRepo) Delete(_ context.Context, _ kernel.TenantID, _ kernel.ShippingZoneID) error {
	return nil
}
func (s *stubZoneRepo) FindByAddress(_ context.Context, _ kernel.TenantID, _, _ string) ([]shipping.ShippingZone, error) {
	return []shipping.ShippingZone{{ID: "zone-1", TenantID: "test-tenant", Name: "US", Countries: []string{"US"}}}, nil
}

// ---------------------------------------------------------------------------
// Stub: shipping.RateRepository
// ---------------------------------------------------------------------------

type stubRateRepo struct{}

func (s *stubRateRepo) Create(_ context.Context, r *shipping.ShippingRate) error { return nil }
func (s *stubRateRepo) GetByID(_ context.Context, _ kernel.TenantID, _ kernel.ShippingRateID) (*shipping.ShippingRate, error) {
	return &shipping.ShippingRate{ID: "rate-1", ZoneID: "zone-1", Name: "Standard", Type: shipping.RateFlat, Price: kernel.Money{Amount: 500, Currency: "USD"}, Active: true}, nil
}
func (s *stubRateRepo) ListByZone(_ context.Context, _ kernel.TenantID, _ kernel.ShippingZoneID) ([]shipping.ShippingRate, error) {
	return []shipping.ShippingRate{{ID: "rate-1", ZoneID: "zone-1", Name: "Standard", Type: shipping.RateFlat, Price: kernel.Money{Amount: 500, Currency: "USD"}, Active: true}}, nil
}
func (s *stubRateRepo) Update(_ context.Context, _ *shipping.ShippingRate) error { return nil }
func (s *stubRateRepo) Delete(_ context.Context, _ kernel.TenantID, _ kernel.ShippingRateID) error {
	return nil
}

// ---------------------------------------------------------------------------
// Stub: tax.Repository
// ---------------------------------------------------------------------------

type stubTaxRepo struct{}

func (s *stubTaxRepo) Create(_ context.Context, r *tax.TaxRate) error { return nil }
func (s *stubTaxRepo) GetByID(_ context.Context, _ kernel.TenantID, _ kernel.TaxRateID) (*tax.TaxRate, error) {
	return &tax.TaxRate{ID: "tax-1", Name: "Sales Tax", Rate: 0.08}, nil
}
func (s *stubTaxRepo) List(_ context.Context, _ kernel.TenantID) ([]tax.TaxRate, error) {
	return []tax.TaxRate{{ID: "tax-1", TenantID: "test-tenant", Name: "Sales Tax", Rate: 0.08, Country: "US", Active: true}}, nil
}
func (s *stubTaxRepo) Update(_ context.Context, _ *tax.TaxRate) error { return nil }
func (s *stubTaxRepo) Delete(_ context.Context, _ kernel.TenantID, _ kernel.TaxRateID) error {
	return nil
}
func (s *stubTaxRepo) FindByLocation(_ context.Context, _ kernel.TenantID, _, _, _, _ string) ([]tax.TaxRate, error) {
	return []tax.TaxRate{{ID: "tax-1", Name: "Sales Tax", Rate: 0.08, Country: "US", Active: true}}, nil
}

// ---------------------------------------------------------------------------
// Stub: payment.Repository
// ---------------------------------------------------------------------------

type stubPaymentRepo struct{}

func (s *stubPaymentRepo) CreatePayment(_ context.Context, _ *payment.Payment) error { return nil }
func (s *stubPaymentRepo) GetPaymentByID(_ context.Context, _ kernel.TenantID, _ kernel.PaymentID) (*payment.Payment, error) {
	return &payment.Payment{ID: "pay-1", OrderID: "ord-1", Amount: kernel.Money{Amount: 1000, Currency: "USD"}, Status: "completed", CreatedAt: time.Now()}, nil
}
func (s *stubPaymentRepo) GetPaymentByOrder(_ context.Context, _ kernel.TenantID, _ kernel.OrderID) (*payment.Payment, error) {
	return &payment.Payment{ID: "pay-1", OrderID: "ord-1", Amount: kernel.Money{Amount: 1000, Currency: "USD"}, Status: "completed", Provider: "manual", CreatedAt: time.Now()}, nil
}
func (s *stubPaymentRepo) UpdatePayment(_ context.Context, _ *payment.Payment) error { return nil }
func (s *stubPaymentRepo) ListPaymentsByOrder(_ context.Context, _ kernel.TenantID, _ kernel.OrderID) ([]payment.Payment, error) {
	return nil, nil
}
func (s *stubPaymentRepo) CreateRefund(_ context.Context, _ *payment.Refund) error { return nil }
func (s *stubPaymentRepo) GetRefundByID(_ context.Context, _ kernel.TenantID, _ kernel.RefundID) (*payment.Refund, error) {
	return nil, nil
}
func (s *stubPaymentRepo) ListRefundsByPayment(_ context.Context, _ kernel.TenantID, _ kernel.PaymentID) ([]payment.Refund, error) {
	return []payment.Refund{{ID: "ref-1", PaymentID: "pay-1", Amount: kernel.Money{Amount: 500, Currency: "USD"}, Status: "completed", CreatedAt: time.Now()}}, nil
}
func (s *stubPaymentRepo) UpdateRefund(_ context.Context, _ *payment.Refund) error { return nil }

// ---------------------------------------------------------------------------
// Stub: search.Repository
// ---------------------------------------------------------------------------

type stubSearchRepo struct{}

func (s *stubSearchRepo) Search(_ context.Context, _ kernel.TenantID, q search.SearchQuery) (*search.SearchResult, error) {
	return &search.SearchResult{
		Products: []search.ProductHit{{ID: "p-1", Name: "Widget", SKU: "WGT-1", Price: kernel.Money{Amount: 999, Currency: "USD"}}},
		Total:    1, Page: 1, PageSize: 20, TotalPages: 1, Query: q.Query,
	}, nil
}
func (s *stubSearchRepo) Suggest(_ context.Context, _ kernel.TenantID, prefix string, _ int) ([]search.SearchSuggestion, error) {
	return []search.SearchSuggestion{{Term: prefix + " widget", Count: 5}}, nil
}

// ---------------------------------------------------------------------------
// Stub: customergroup.Repository
// ---------------------------------------------------------------------------

type stubCustomerGroupRepo struct{}

func (s *stubCustomerGroupRepo) Create(_ context.Context, g *customergroup.CustomerGroup) error {
	return nil
}
func (s *stubCustomerGroupRepo) GetByID(_ context.Context, _ kernel.TenantID, _ kernel.CustomerGroupID) (*customergroup.CustomerGroup, error) {
	return &customergroup.CustomerGroup{ID: "cg-1", Name: "VIP"}, nil
}
func (s *stubCustomerGroupRepo) List(_ context.Context, _ kernel.TenantID) ([]customergroup.CustomerGroup, error) {
	return []customergroup.CustomerGroup{{ID: "cg-1", TenantID: "test-tenant", Name: "VIP", Description: "VIP customers"}}, nil
}
func (s *stubCustomerGroupRepo) Update(_ context.Context, _ *customergroup.CustomerGroup) error {
	return nil
}
func (s *stubCustomerGroupRepo) Delete(_ context.Context, _ kernel.TenantID, _ kernel.CustomerGroupID) error {
	return nil
}
func (s *stubCustomerGroupRepo) AddMember(_ context.Context, _ *customergroup.GroupMembership) error {
	return nil
}
func (s *stubCustomerGroupRepo) RemoveMember(_ context.Context, _ kernel.TenantID, _ kernel.CustomerGroupID, _ kernel.CustomerID) error {
	return nil
}
func (s *stubCustomerGroupRepo) ListMembers(_ context.Context, _ kernel.TenantID, _ kernel.CustomerGroupID) ([]customergroup.GroupMembership, error) {
	return nil, nil
}
func (s *stubCustomerGroupRepo) GetCustomerGroups(_ context.Context, _ kernel.TenantID, _ kernel.CustomerID) ([]customergroup.CustomerGroup, error) {
	return nil, nil
}

// ---------------------------------------------------------------------------
// Stub: giftcard.Repository
// ---------------------------------------------------------------------------

type stubGiftCardRepo struct{}

func (s *stubGiftCardRepo) Create(_ context.Context, gc *giftcard.GiftCard) error { return nil }
func (s *stubGiftCardRepo) GetByID(_ context.Context, _ kernel.TenantID, _ kernel.GiftCardID) (*giftcard.GiftCard, error) {
	return &giftcard.GiftCard{ID: "gc-1", Code: "GIFT-100", Balance: kernel.Money{Amount: 10000, Currency: "USD"}, Active: true}, nil
}
func (s *stubGiftCardRepo) GetByCode(_ context.Context, _ kernel.TenantID, _ string) (*giftcard.GiftCard, error) {
	return &giftcard.GiftCard{ID: "gc-1", Code: "GIFT-100", Balance: kernel.Money{Amount: 10000, Currency: "USD"}, InitialAmount: kernel.Money{Amount: 10000, Currency: "USD"}, Active: true}, nil
}
func (s *stubGiftCardRepo) List(_ context.Context, _ kernel.TenantID, _ kernel.PaginationOptions) (kernel.Paginated[giftcard.GiftCard], error) {
	return kernel.Paginated[giftcard.GiftCard]{
		Items: []giftcard.GiftCard{{ID: "gc-1", Code: "GIFT-100", Balance: kernel.Money{Amount: 10000, Currency: "USD"}, Active: true}},
		Total: 1, Page: 1, TotalPages: 1,
	}, nil
}
func (s *stubGiftCardRepo) Update(_ context.Context, _ *giftcard.GiftCard) error    { return nil }
func (s *stubGiftCardRepo) Delete(_ context.Context, _ kernel.TenantID, _ kernel.GiftCardID) error { return nil }
func (s *stubGiftCardRepo) CreateTransaction(_ context.Context, _ *giftcard.GiftCardTransaction) error {
	return nil
}
func (s *stubGiftCardRepo) ListTransactions(_ context.Context, _ kernel.TenantID, _ kernel.GiftCardID) ([]giftcard.GiftCardTransaction, error) {
	return nil, nil
}

// ---------------------------------------------------------------------------
// Stub: cartrecovery.Repository
// ---------------------------------------------------------------------------

type stubCartRecoveryRepo struct{}

func (s *stubCartRecoveryRepo) Create(_ context.Context, _ *cartrecovery.RecoveryEmail) error {
	return nil
}
func (s *stubCartRecoveryRepo) GetByID(_ context.Context, _ kernel.TenantID, _ kernel.RecoveryID) (*cartrecovery.RecoveryEmail, error) {
	return nil, nil
}
func (s *stubCartRecoveryRepo) GetByCartID(_ context.Context, _ kernel.TenantID, _ kernel.CartID) ([]cartrecovery.RecoveryEmail, error) {
	return nil, nil
}
func (s *stubCartRecoveryRepo) ListPending(_ context.Context, _ kernel.TenantID) ([]cartrecovery.RecoveryEmail, error) {
	return nil, nil
}
func (s *stubCartRecoveryRepo) List(_ context.Context, _ kernel.TenantID, _, _ int) (kernel.Paginated[cartrecovery.RecoveryEmail], error) {
	return kernel.Paginated[cartrecovery.RecoveryEmail]{
		Items: []cartrecovery.RecoveryEmail{{ID: "rec-1", Email: "john@example.com", Status: "pending"}},
		Total: 1, Page: 1, TotalPages: 1,
	}, nil
}
func (s *stubCartRecoveryRepo) Update(_ context.Context, _ *cartrecovery.RecoveryEmail) error {
	return nil
}
func (s *stubCartRecoveryRepo) GetStats(_ context.Context, _ kernel.TenantID) (*cartrecovery.RecoveryStats, error) {
	return &cartrecovery.RecoveryStats{Total: 10, Pending: 3, Sent: 4, Clicked: 2, Converted: 1, ConversionRate: 10.0}, nil
}

// ---------------------------------------------------------------------------
// Stub: currency.Repository
// ---------------------------------------------------------------------------

type stubCurrencyRepo struct {
	rates map[string]*currency.CurrencyRate // key = "base:target"
}

func newStubCurrencyRepo() *stubCurrencyRepo {
	return &stubCurrencyRepo{rates: map[string]*currency.CurrencyRate{}}
}

func (s *stubCurrencyRepo) Create(_ context.Context, r *currency.CurrencyRate) error {
	s.rates[r.BaseCurrency+":"+r.TargetCurrency] = r
	return nil
}
func (s *stubCurrencyRepo) GetByPair(_ context.Context, _ kernel.TenantID, base, target string) (*currency.CurrencyRate, error) {
	if r, ok := s.rates[base+":"+target]; ok {
		return r, nil
	}
	return nil, currency.ErrRateNotFound
}
func (s *stubCurrencyRepo) GetByID(_ context.Context, _ kernel.TenantID, _ kernel.CurrencyRateID) (*currency.CurrencyRate, error) {
	return nil, currency.ErrRateNotFound
}
func (s *stubCurrencyRepo) List(_ context.Context, _ kernel.TenantID) ([]currency.CurrencyRate, error) {
	var out []currency.CurrencyRate
	for _, r := range s.rates {
		out = append(out, *r)
	}
	return out, nil
}
func (s *stubCurrencyRepo) Update(_ context.Context, r *currency.CurrencyRate) error {
	s.rates[r.BaseCurrency+":"+r.TargetCurrency] = r
	return nil
}
func (s *stubCurrencyRepo) Delete(_ context.Context, _ kernel.TenantID, _ kernel.CurrencyRateID) error {
	return nil
}

// ---------------------------------------------------------------------------
// Stub: i18n.Repository
// ---------------------------------------------------------------------------

type stubI18nRepo struct{}

func (s *stubI18nRepo) Upsert(_ context.Context, _ *i18n.Translation) error { return nil }
func (s *stubI18nRepo) GetByEntity(_ context.Context, _ kernel.TenantID, _, _, _ string) ([]i18n.Translation, error) {
	return nil, nil
}
func (s *stubI18nRepo) GetBundle(_ context.Context, _ kernel.TenantID, entityType, entityID, locale string) (*i18n.TranslationBundle, error) {
	return &i18n.TranslationBundle{
		EntityType: entityType, EntityID: entityID, Locale: locale,
		Fields: map[string]string{"name": "Widget", "description": "A widget"},
	}, nil
}
func (s *stubI18nRepo) ListLocales(_ context.Context, _ kernel.TenantID, _, _ string) ([]string, error) {
	return []string{"en", "es"}, nil
}
func (s *stubI18nRepo) Delete(_ context.Context, _ kernel.TenantID, _, _, _, _ string) error {
	return nil
}
func (s *stubI18nRepo) DeleteAll(_ context.Context, _ kernel.TenantID, _, _ string) error {
	return nil
}

// ---------------------------------------------------------------------------
// Stub: subscription.Repository
// ---------------------------------------------------------------------------

type stubSubscriptionRepo struct {
	subs map[kernel.SubscriptionID]*subscription.Subscription
}

func newStubSubscriptionRepo() *stubSubscriptionRepo {
	return &stubSubscriptionRepo{subs: map[kernel.SubscriptionID]*subscription.Subscription{}}
}

func (s *stubSubscriptionRepo) Create(_ context.Context, sub *subscription.Subscription) error {
	s.subs[sub.ID] = sub
	return nil
}
func (s *stubSubscriptionRepo) GetByID(_ context.Context, _ kernel.TenantID, id kernel.SubscriptionID) (*subscription.Subscription, error) {
	if sub, ok := s.subs[id]; ok {
		return sub, nil
	}
	return nil, subscription.ErrNotFound
}
func (s *stubSubscriptionRepo) Update(_ context.Context, sub *subscription.Subscription) error {
	s.subs[sub.ID] = sub
	return nil
}
func (s *stubSubscriptionRepo) Delete(_ context.Context, _ kernel.TenantID, _ kernel.SubscriptionID) error {
	return nil
}
func (s *stubSubscriptionRepo) List(_ context.Context, _ kernel.TenantID, _, _ int) (kernel.Paginated[subscription.Subscription], error) {
	items := make([]subscription.Subscription, 0, len(s.subs))
	for _, sub := range s.subs {
		items = append(items, *sub)
	}
	return kernel.Paginated[subscription.Subscription]{Items: items, Total: len(items), Page: 1, TotalPages: 1}, nil
}
func (s *stubSubscriptionRepo) ListByCustomer(_ context.Context, _ kernel.TenantID, _ kernel.CustomerID) ([]subscription.Subscription, error) {
	return nil, nil
}
func (s *stubSubscriptionRepo) ListDueBilling(_ context.Context, _ kernel.TenantID, _ time.Time) ([]subscription.Subscription, error) {
	return nil, nil
}
func (s *stubSubscriptionRepo) CreateBillingRecord(_ context.Context, _ *subscription.BillingRecord) error {
	return nil
}
func (s *stubSubscriptionRepo) ListBillingRecords(_ context.Context, _ kernel.TenantID, _ kernel.SubscriptionID, _, _ int) (kernel.Paginated[subscription.BillingRecord], error) {
	return kernel.Paginated[subscription.BillingRecord]{Items: []subscription.BillingRecord{}, Total: 0, Page: 1, TotalPages: 0}, nil
}

// ===========================================================================
// Shipping Tests
// ===========================================================================

func TestListShippingZonesTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := shippingsrv.New(&stubZoneRepo{}, &stubRateRepo{}, bus)
	tool := &ListShippingZonesTool{shipping: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "US") {
		t.Errorf("expected zone name 'US' in result, got: %s", result)
	}
}

func TestCreateShippingZoneTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := shippingsrv.New(&stubZoneRepo{}, &stubRateRepo{}, bus)
	tool := &CreateShippingZoneTool{shipping: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"name":"Europe","countries":["DE","FR"]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Europe") {
		t.Errorf("expected 'Europe' in result, got: %s", result)
	}
}

func TestCalculateShippingTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := shippingsrv.New(&stubZoneRepo{}, &stubRateRepo{}, bus)
	tool := &CalculateShippingTool{shipping: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"country":"US","state":"CA","order_amount_cents":5000,"weight_kg":2.5}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Standard") {
		t.Errorf("expected rate name 'Standard' in result, got: %s", result)
	}
}

func TestCalculateShippingTool_Execute_InvalidJSON(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := shippingsrv.New(&stubZoneRepo{}, &stubRateRepo{}, bus)
	tool := &CalculateShippingTool{shipping: svc, tenantID: "test-tenant"}
	_, err := tool.Execute(context.Background(), json.RawMessage(`{bad`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// ===========================================================================
// Tax Tests
// ===========================================================================

func TestListTaxRatesTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := taxsrv.New(&stubTaxRepo{}, bus)
	tool := &ListTaxRatesTool{tax: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Sales Tax") {
		t.Errorf("expected 'Sales Tax' in result, got: %s", result)
	}
}

func TestCreateTaxRateTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := taxsrv.New(&stubTaxRepo{}, bus)
	tool := &CreateTaxRateTool{tax: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"name":"VAT","rate":0.20,"country":"DE"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "VAT") {
		t.Errorf("expected 'VAT' in result, got: %s", result)
	}
}

func TestCalculateTaxTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := taxsrv.New(&stubTaxRepo{}, bus)
	tool := &CalculateTaxTool{tax: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"subtotal_cents":10000,"shipping_cents":500,"country":"US","state":"CA"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "tax") || !strings.Contains(result, "800") {
		t.Errorf("expected tax result with 800 cents (8%% of 10000), got: %s", result)
	}
}

// ===========================================================================
// Payment Tests
// ===========================================================================

func TestGetOrderPaymentTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := paymentsrv.New(&stubPaymentRepo{}, bus, map[string]payment.PaymentProvider{})
	tool := &GetOrderPaymentTool{payment: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"order_id":"ord-1"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "pay-1") {
		t.Errorf("expected payment ID in result, got: %s", result)
	}
}

func TestListRefundsTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := paymentsrv.New(&stubPaymentRepo{}, bus, map[string]payment.PaymentProvider{})
	tool := &ListRefundsTool{payment: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"payment_id":"pay-1"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "ref-1") {
		t.Errorf("expected refund ID in result, got: %s", result)
	}
}

// ===========================================================================
// Search Tests
// ===========================================================================

func TestSearchProductsTool_Execute(t *testing.T) {
	svc := searchsrv.New(&stubSearchRepo{})
	tool := &SearchProductsTool{search: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"query":"widget"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Widget") {
		t.Errorf("expected 'Widget' in result, got: %s", result)
	}
}

func TestSearchSuggestionsTool_Execute(t *testing.T) {
	svc := searchsrv.New(&stubSearchRepo{})
	tool := &SearchSuggestionsTool{search: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"prefix":"wid","limit":5}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "widget") {
		t.Errorf("expected suggestion containing 'widget', got: %s", result)
	}
}

// ===========================================================================
// Customer Group Tests
// ===========================================================================

func TestListCustomerGroupsTool_Execute(t *testing.T) {
	svc := customergroupsrv.New(&stubCustomerGroupRepo{})
	tool := &ListCustomerGroupsTool{groups: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "VIP") {
		t.Errorf("expected 'VIP' in result, got: %s", result)
	}
}

func TestCreateCustomerGroupTool_Execute(t *testing.T) {
	svc := customergroupsrv.New(&stubCustomerGroupRepo{})
	tool := &CreateCustomerGroupTool{groups: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"name":"Wholesale","description":"Wholesale buyers"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Wholesale") {
		t.Errorf("expected 'Wholesale' in result, got: %s", result)
	}
}

func TestAddGroupMemberTool_Execute(t *testing.T) {
	svc := customergroupsrv.New(&stubCustomerGroupRepo{})
	tool := &AddGroupMemberTool{groups: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"group_id":"cg-1","customer_id":"cust-1"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.ToLower(result), "added") && !strings.Contains(strings.ToLower(result), "member") {
		t.Errorf("expected membership confirmation in result, got: %s", result)
	}
}

// ===========================================================================
// Gift Card Tests
// ===========================================================================

func TestListGiftCardsTool_Execute(t *testing.T) {
	svc := giftcardsrv.New(&stubGiftCardRepo{})
	tool := &ListGiftCardsTool{giftcards: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "GIFT-100") {
		t.Errorf("expected gift card code in result, got: %s", result)
	}
}

func TestCreateGiftCardTool_Execute(t *testing.T) {
	svc := giftcardsrv.New(&stubGiftCardRepo{})
	tool := &CreateGiftCardTool{giftcards: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"code":"NEW-CARD","amount_cents":5000,"currency":"USD"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "NEW-CARD") {
		t.Errorf("expected 'NEW-CARD' in result, got: %s", result)
	}
}

func TestCheckGiftCardBalanceTool_Execute(t *testing.T) {
	svc := giftcardsrv.New(&stubGiftCardRepo{})
	tool := &CheckGiftCardBalanceTool{giftcards: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"code":"GIFT-100"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "10000") || !strings.Contains(result, "GIFT-100") {
		t.Errorf("expected balance and code in result, got: %s", result)
	}
}

func TestRedeemGiftCardTool_Execute(t *testing.T) {
	svc := giftcardsrv.New(&stubGiftCardRepo{})
	tool := &RedeemGiftCardTool{giftcards: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"code":"GIFT-100","amount_cents":2000,"currency":"USD"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "GIFT-100") {
		t.Errorf("expected gift card code in result, got: %s", result)
	}
}

// ===========================================================================
// Cart Recovery Tests
// ===========================================================================

func TestListRecoveryEmailsTool_Execute(t *testing.T) {
	svc := cartrecoverysrv.New(&stubCartRecoveryRepo{})
	tool := &ListRecoveryEmailsTool{recovery: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "john@example.com") {
		t.Errorf("expected email in result, got: %s", result)
	}
}

func TestGetRecoveryStatsTool_Execute(t *testing.T) {
	svc := cartrecoverysrv.New(&stubCartRecoveryRepo{})
	tool := &GetRecoveryStatsTool{recovery: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "10") {
		t.Errorf("expected total count in result, got: %s", result)
	}
}

// ===========================================================================
// Currency Tests
// ===========================================================================

func TestSetCurrencyRateTool_Execute(t *testing.T) {
	repo := newStubCurrencyRepo()
	svc := currencysrv.New(repo)
	tool := &SetCurrencyRateTool{currency: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"base_currency":"USD","target_currency":"EUR","rate":0.85}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "USD") || !strings.Contains(result, "EUR") {
		t.Errorf("expected currency codes in result, got: %s", result)
	}
}

func TestConvertCurrencyTool_Execute(t *testing.T) {
	repo := newStubCurrencyRepo()
	// Pre-seed a rate
	repo.rates["USD:EUR"] = &currency.CurrencyRate{
		ID: "cr-1", TenantID: "test-tenant",
		BaseCurrency: "USD", TargetCurrency: "EUR", Rate: 0.85,
	}
	svc := currencysrv.New(repo)
	tool := &ConvertCurrencyTool{currency: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"from_currency":"USD","target_currency":"EUR","amount":10000}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "8500") {
		t.Errorf("expected converted amount 8500 in result, got: %s", result)
	}
}

func TestListCurrencyRatesTool_Execute(t *testing.T) {
	repo := newStubCurrencyRepo()
	repo.rates["USD:EUR"] = &currency.CurrencyRate{
		ID: "cr-1", TenantID: "test-tenant",
		BaseCurrency: "USD", TargetCurrency: "EUR", Rate: 0.85,
	}
	svc := currencysrv.New(repo)
	tool := &ListCurrencyRatesTool{currency: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "USD") {
		t.Errorf("expected 'USD' in result, got: %s", result)
	}
}

// ===========================================================================
// i18n Tests
// ===========================================================================

func TestSetTranslationsTool_Execute(t *testing.T) {
	svc := i18nsrv.New(&stubI18nRepo{})
	tool := &SetTranslationsTool{i18n: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"entity_type":"product","entity_id":"p-1","locale":"es","fields":{"name":"Widget ES","description":"A widget in Spanish"}}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.ToLower(result), "translation") || !strings.Contains(result, "es") {
		t.Errorf("expected translation confirmation, got: %s", result)
	}
}

func TestGetTranslationsTool_Execute(t *testing.T) {
	svc := i18nsrv.New(&stubI18nRepo{})
	tool := &GetTranslationsTool{i18n: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"entity_type":"product","entity_id":"p-1","locale":"en"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Widget") {
		t.Errorf("expected 'Widget' in result, got: %s", result)
	}
}

func TestListSupportedLocalesTool_Execute(t *testing.T) {
	svc := i18nsrv.New(&stubI18nRepo{})
	tool := &ListSupportedLocalesTool{i18n: svc}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "en") || !strings.Contains(result, "es") {
		t.Errorf("expected supported locales in result, got: %s", result)
	}
}

// ===========================================================================
// Subscription Tests
// ===========================================================================

func TestListSubscriptionsTool_Execute(t *testing.T) {
	repo := newStubSubscriptionRepo()
	bus := eventbus.NewInMemoryBus()
	svc := subscriptionsrv.New(repo, bus)
	tool := &ListSubscriptionsTool{subscriptions: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "subscription") && !strings.Contains(result, "0 subscription") && !strings.Contains(result, "Subscription") {
		t.Errorf("expected subscription info in result, got: %s", result)
	}
}

func TestCreateSubscriptionTool_Execute(t *testing.T) {
	repo := newStubSubscriptionRepo()
	bus := eventbus.NewInMemoryBus()
	svc := subscriptionsrv.New(repo, bus)
	tool := &CreateSubscriptionTool{subscriptions: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"customer_id":"cust-1","product_id":"prod-1","price_cents":2999,"currency":"USD","interval":"monthly"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "cust-1") || !strings.Contains(result, "monthly") {
		t.Errorf("expected subscription details in result, got: %s", result)
	}
}

func TestCancelSubscriptionTool_Execute(t *testing.T) {
	repo := newStubSubscriptionRepo()
	// Pre-seed an active subscription
	repo.subs["sub-1"] = &subscription.Subscription{
		ID:         "sub-1",
		TenantID:   "test-tenant",
		CustomerID: "cust-1",
		ProductID:  "prod-1",
		Status:     subscription.StatusActive,
		Price:      kernel.Money{Amount: 2999, Currency: "USD"},
		Interval:   subscription.IntervalMonthly,
	}
	bus := eventbus.NewInMemoryBus()
	svc := subscriptionsrv.New(repo, bus)
	tool := &CancelSubscriptionTool{subscriptions: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"subscription_id":"sub-1","reason":"Too expensive"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "cancel") || !strings.Contains(result, "sub-1") {
		t.Errorf("expected cancellation confirmation, got: %s", result)
	}
}

// ===========================================================================
// Product Variant Tests (from tools_commerce.go, using productsrv)
// ===========================================================================

func TestCreateProductOptionTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := productsrv.New(&stubProductRepo{}, &stubVariantRepo{}, bus)
	tool := &CreateProductOptionTool{products: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"product_id":"p-1","name":"Size","values":["S","M","L"]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Size") {
		t.Errorf("expected 'Size' in result, got: %s", result)
	}
}

func TestListProductOptionsTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := productsrv.New(&stubProductRepo{}, &stubVariantRepo{}, bus)
	tool := &ListProductOptionsTool{products: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"product_id":"p-1"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "No options") && !strings.Contains(result, "option") {
		t.Errorf("expected options info in result, got: %s", result)
	}
}

func TestCreateProductVariantTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := productsrv.New(&stubProductRepo{}, &stubVariantRepo{}, bus)
	tool := &CreateProductVariantTool{products: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"product_id":"p-1","sku":"WGT-RED","price":1099,"currency":"USD","options":{"color":"red"}}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "WGT-RED") {
		t.Errorf("expected variant SKU in result, got: %s", result)
	}
}

func TestListProductVariantsTool_Execute(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	svc := productsrv.New(&stubProductRepo{}, &stubVariantRepo{}, bus)
	tool := &ListProductVariantsTool{products: svc, tenantID: "test-tenant"}
	result, err := tool.Execute(context.Background(), json.RawMessage(`{"product_id":"p-1"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "No variants") && !strings.Contains(result, "variant") {
		t.Errorf("expected variants info in result, got: %s", result)
	}
}
