// Package main contains a route contract test that verifies all expected API
// routes are registered on the Fiber app. This prevents frontend→backend path
// mismatches from going undetected.
//
// The test instantiates each handler with nil services (RegisterRoutes never
// calls the service, so nil is safe) and checks that every route the frontend
// depends on is present in app.GetRoutes(true).
//
// Run with:
//
//	go test ./cmd/ -run TestRouteContract -v
package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/vendex/internal/abtest/abtestapi"
	"github.com/Abraxas-365/vendex/internal/agentmemory/agentmemoryapi"
	"github.com/Abraxas-365/vendex/internal/agentsession/agentsessionapi"
	"github.com/Abraxas-365/vendex/internal/agenttrigger/agenttriggerapi"
	"github.com/Abraxas-365/vendex/internal/analytics/analyticsapi"
	"github.com/Abraxas-365/vendex/internal/approval/approvalapi"
	"github.com/Abraxas-365/vendex/internal/audit/auditapi"
	"github.com/Abraxas-365/vendex/internal/blog/blogapi"
	"github.com/Abraxas-365/vendex/internal/bulkops/bulkopsapi"
	"github.com/Abraxas-365/vendex/internal/bundle/bundleapi"
	"github.com/Abraxas-365/vendex/internal/cartrecovery/cartrecoveryapi"
	"github.com/Abraxas-365/vendex/internal/catalog/catalogapi"
	"github.com/Abraxas-365/vendex/internal/collection/collectionapi"
	"github.com/Abraxas-365/vendex/internal/currency/currencyapi"
	"github.com/Abraxas-365/vendex/internal/customer/customerapi"
	"github.com/Abraxas-365/vendex/internal/customergroup/customergroupapi"
	"github.com/Abraxas-365/vendex/internal/dashboard/dashboardapi"
	"github.com/Abraxas-365/vendex/internal/giftcard/giftcardapi"
	"github.com/Abraxas-365/vendex/internal/i18n/i18napi"
	"github.com/Abraxas-365/vendex/internal/importexport"
	"github.com/Abraxas-365/vendex/internal/inventory/inventoryapi"
	"github.com/Abraxas-365/vendex/internal/loyalty/loyaltyapi"
	"github.com/Abraxas-365/vendex/internal/marketplace/marketplaceapi"
	"github.com/Abraxas-365/vendex/internal/media/mediaapi"
	"github.com/Abraxas-365/vendex/internal/multistore/multistoreapi"
	"github.com/Abraxas-365/vendex/internal/notification/notificationapi"
	"github.com/Abraxas-365/vendex/internal/order/orderapi"
	"github.com/Abraxas-365/vendex/internal/payment/paymentapi"
	"github.com/Abraxas-365/vendex/internal/plugin/pluginapi"
	"github.com/Abraxas-365/vendex/internal/product/productapi"
	"github.com/Abraxas-365/vendex/internal/promo/promoapi"
	"github.com/Abraxas-365/vendex/internal/recommendation/recommendationapi"
	"github.com/Abraxas-365/vendex/internal/returns/returnsapi"
	"github.com/Abraxas-365/vendex/internal/review/reviewapi"
	"github.com/Abraxas-365/vendex/internal/settings/settingsapi"
	"github.com/Abraxas-365/vendex/internal/shipping/shippingapi"
	"github.com/Abraxas-365/vendex/internal/storefront/storefrontapi"
	"github.com/Abraxas-365/vendex/internal/subscription/subscriptionapi"
	"github.com/Abraxas-365/vendex/internal/tax/taxapi"
	"github.com/Abraxas-365/vendex/internal/theme/themeapi"
	"github.com/Abraxas-365/vendex/internal/webhook/webhookapi"
)

// routeKey is a method+path pair used as a map key.
// paths are stored without trailing slashes for consistent matching.
type routeKey struct{ method, path string }

// buildRouteSet registers all handlers onto a fresh Fiber app and returns a
// set of every registered method+path combination.
func buildRouteSet(t *testing.T) map[routeKey]bool {
	t.Helper()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	api := app.Group("/api/v1")

	// --- Products ---
	productapi.NewHandler(nil).RegisterRoutes(api)

	// --- Orders ---
	orderapi.NewHandler(nil).RegisterRoutes(api)

	// --- Catalog (categories + collections public) ---
	catalogapi.NewHandler(nil).RegisterRoutes(api)

	// --- Collections (admin) ---
	collectionapi.NewHandler(nil).RegisterRoutes(api)

	// --- Storefront pages + block-types ---
	storefrontapi.NewHandler(nil, nil).RegisterRoutes(api)

	// --- Promos ---
	promoapi.New(nil).RegisterRoutes(api)

	// --- Media ---
	mediaapi.NewHandler(nil).RegisterRoutes(api)

	// --- Marketplace ---
	marketplaceapi.NewHandler(nil).RegisterRoutes(api)

	// --- Analytics ---
	analyticsapi.NewHandler(nil).RegisterRoutes(api)

	// --- Dashboard ---
	dashboardapi.NewHandler(nil).RegisterRoutes(api)

	// --- Settings ---
	settingsapi.NewHandler(nil).RegisterRoutes(api)

	// --- Themes ---
	themeapi.NewHandler(nil).RegisterRoutes(api)

	// --- Shipping ---
	shippingapi.NewHandler(nil).RegisterRoutes(api)

	// --- Tax ---
	taxapi.NewHandler(nil).RegisterRoutes(api)

	// --- Payments ---
	paymentapi.NewHandler(nil).RegisterRoutes(api)

	// --- Customer groups ---
	customergroupapi.NewHandler(nil).RegisterRoutes(api)

	// --- Gift cards ---
	giftcardapi.New(nil).RegisterRoutes(api)

	// --- Customers ---
	customerapi.NewHandler(nil).RegisterRoutes(api)

	// --- Currency ---
	currencyapi.NewHandler(nil).RegisterRoutes(api)

	// --- Import/Export ---
	importexport.NewHandler(nil, nil).RegisterRoutes(api)

	// --- Subscriptions ---
	subscriptionapi.NewHandler(nil).RegisterRoutes(api)

	// --- Inventory ---
	inventoryapi.NewHandler(nil).RegisterRoutes(api)

	// --- Reviews ---
	reviewapi.NewHandler(nil).RegisterRoutes(api)

	// --- Returns ---
	returnsapi.NewHandler(nil).RegisterRoutes(api)

	// --- Webhooks ---
	webhookapi.NewHandler(nil).RegisterRoutes(api)

	// --- Audit ---
	auditapi.NewHandler(nil).RegisterRoutes(api)

	// --- Loyalty ---
	loyaltyapi.New(nil).RegisterRoutes(api)

	// --- Bundles ---
	bundleapi.NewHandler(nil).RegisterRoutes(api)

	// --- Notifications ---
	notificationapi.NewHandler(nil).RegisterRoutes(api)

	// --- Multistore (storefronts) ---
	multistoreapi.NewHandler(nil).RegisterRoutes(api)

	// --- Bulk operations ---
	bulkopsapi.NewHandler(nil).RegisterRoutes(api)

	// --- Blog ---
	blogapi.NewHandler(nil).RegisterRoutes(api)

	// --- A/B Tests (experiments) ---
	abtestapi.NewHandler(nil).RegisterRoutes(api)

	// --- Recommendations ---
	recommendationapi.New(nil).RegisterRoutes(api)

	// --- Plugins ---
	pluginapi.NewHandler(nil).RegisterRoutes(api)

	// --- Approvals ---
	approvalapi.NewHandler(nil).RegisterRoutes(api.Group("/approvals"))

	// --- Cart recovery ---
	cartrecoveryapi.NewHandler(nil).RegisterRoutes(api)

	// --- i18n ---
	i18napi.NewHandler(nil).RegisterRoutes(api)

	// --- Agent sessions ---
	agentsessionapi.NewHandler(nil).RegisterRoutes(api.Group("/agent/sessions"))

	// --- Agent memory ---
	agentmemoryapi.NewHandler(nil).RegisterRoutes(api.Group("/agent/memories"))

	// --- Agent triggers ---
	agenttriggerapi.NewHandler(nil).RegisterRoutes(api.Group("/agent/triggers"))

	registered := make(map[routeKey]bool)
	for _, r := range app.GetRoutes(true) {
		// Normalize: strip trailing slash so "/products/" == "/products"
		path := strings.TrimRight(r.Path, "/")
		if path == "" {
			path = "/"
		}
		registered[routeKey{r.Method, path}] = true
	}
	return registered
}

// assertRoute fails the test if method+path is not in the registered set.
func assertRoute(t *testing.T, registered map[routeKey]bool, method, path string) {
	t.Helper()
	full := strings.TrimRight(fmt.Sprintf("/api/v1%s", path), "/")
	if full == "" {
		full = "/"
	}
	if !registered[routeKey{method, full}] {
		t.Errorf("MISSING route: %s %s", method, full)
	}
}

// TestRouteContract verifies every frontend-expected route is registered.
func TestRouteContract(t *testing.T) {
	registered := buildRouteSet(t)

	t.Run("Products", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/products")
		assertRoute(t, registered, "POST", "/products")
		assertRoute(t, registered, "GET", "/products/:id")
		assertRoute(t, registered, "PUT", "/products/:id")
		assertRoute(t, registered, "DELETE", "/products/:id")
		assertRoute(t, registered, "GET", "/products/:id/options")
		assertRoute(t, registered, "POST", "/products/:id/options")
		assertRoute(t, registered, "PUT", "/products/options/:optionId")
		assertRoute(t, registered, "DELETE", "/products/options/:optionId")
		assertRoute(t, registered, "GET", "/products/:id/variants")
		assertRoute(t, registered, "POST", "/products/:id/variants")
		assertRoute(t, registered, "PUT", "/products/variants/:variantId")
		assertRoute(t, registered, "DELETE", "/products/variants/:variantId")
	})

	t.Run("Orders", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/orders")
		assertRoute(t, registered, "POST", "/orders")
		assertRoute(t, registered, "GET", "/orders/:id")
		assertRoute(t, registered, "PUT", "/orders/:id/status")
		assertRoute(t, registered, "POST", "/orders/:id/cancel")
	})

	t.Run("Categories", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/categories")
		assertRoute(t, registered, "POST", "/categories")
		assertRoute(t, registered, "GET", "/categories/:id")
		assertRoute(t, registered, "DELETE", "/categories/:id")
	})

	t.Run("Collections", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/collections")
		assertRoute(t, registered, "POST", "/collections")
		assertRoute(t, registered, "GET", "/collections/:id")
		assertRoute(t, registered, "PUT", "/collections/:id")
		assertRoute(t, registered, "DELETE", "/collections/:id")
		assertRoute(t, registered, "GET", "/collections/:id/products")
		assertRoute(t, registered, "POST", "/collections/:id/products")
		assertRoute(t, registered, "DELETE", "/collections/:id/products/:productId")
	})

	t.Run("StorefrontPages", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/storefront/pages")
		assertRoute(t, registered, "POST", "/storefront/pages")
		assertRoute(t, registered, "GET", "/storefront/pages/:id")
		assertRoute(t, registered, "PUT", "/storefront/pages/:id")
		assertRoute(t, registered, "DELETE", "/storefront/pages/:id")
		assertRoute(t, registered, "GET", "/storefront/pages/by-slug/:slug")
		assertRoute(t, registered, "POST", "/storefront/pages/:id/publish")
		assertRoute(t, registered, "POST", "/storefront/pages/:id/unpublish")
		assertRoute(t, registered, "POST", "/storefront/pages/:id/archive")
		assertRoute(t, registered, "GET", "/storefront/pages/:id/versions")
		assertRoute(t, registered, "GET", "/storefront/pages/:id/versions/:version")
	})

	t.Run("BlockTypes", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/storefront/block-types")
		assertRoute(t, registered, "POST", "/storefront/block-types")
		assertRoute(t, registered, "GET", "/storefront/block-types/:id")
		assertRoute(t, registered, "PUT", "/storefront/block-types/:id")
		assertRoute(t, registered, "DELETE", "/storefront/block-types/:id")
	})

	t.Run("Promos", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/admin/promos")
		assertRoute(t, registered, "POST", "/admin/promos")
		assertRoute(t, registered, "GET", "/admin/promos/:id")
		assertRoute(t, registered, "PUT", "/admin/promos/:id")
		assertRoute(t, registered, "POST", "/admin/promos/:id/deactivate")
		assertRoute(t, registered, "POST", "/promos/validate")
	})

	t.Run("Media", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/media")
		assertRoute(t, registered, "POST", "/media/upload")
		assertRoute(t, registered, "DELETE", "/media/:id")
	})

	t.Run("Marketplace", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/marketplace/vendors")
		assertRoute(t, registered, "POST", "/marketplace/vendors")
		assertRoute(t, registered, "GET", "/marketplace/vendors/:id")
		assertRoute(t, registered, "PUT", "/marketplace/vendors/:id")
		assertRoute(t, registered, "DELETE", "/marketplace/vendors/:id")
	})

	t.Run("Analytics", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/analytics/dashboard")
		assertRoute(t, registered, "GET", "/analytics/revenue")
		assertRoute(t, registered, "GET", "/analytics/top-products")
		assertRoute(t, registered, "GET", "/analytics/order-status")
		assertRoute(t, registered, "GET", "/analytics/recent-orders")
	})

	t.Run("Dashboard", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/dashboard/sales")
		assertRoute(t, registered, "GET", "/dashboard/top-products")
		assertRoute(t, registered, "GET", "/dashboard/revenue")
		assertRoute(t, registered, "GET", "/dashboard/customers")
		assertRoute(t, registered, "GET", "/dashboard/funnel")
	})

	t.Run("Settings", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/settings")
		assertRoute(t, registered, "PUT", "/settings")
	})

	t.Run("Themes", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/themes")
		assertRoute(t, registered, "POST", "/themes")
		assertRoute(t, registered, "GET", "/themes/active")
		assertRoute(t, registered, "GET", "/themes/:id")
		assertRoute(t, registered, "PUT", "/themes/:id")
		assertRoute(t, registered, "DELETE", "/themes/:id")
		assertRoute(t, registered, "POST", "/themes/:id/activate")
	})

	t.Run("Shipping", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/shipping/zones")
		assertRoute(t, registered, "POST", "/shipping/zones")
		assertRoute(t, registered, "GET", "/shipping/zones/:id")
		assertRoute(t, registered, "PUT", "/shipping/zones/:id")
		assertRoute(t, registered, "DELETE", "/shipping/zones/:id")
		assertRoute(t, registered, "GET", "/shipping/zones/:zoneId/rates")
		assertRoute(t, registered, "POST", "/shipping/zones/:zoneId/rates")
		assertRoute(t, registered, "PUT", "/shipping/rates/:id")
		assertRoute(t, registered, "DELETE", "/shipping/rates/:id")
	})

	t.Run("Tax", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/tax/rates")
		assertRoute(t, registered, "POST", "/tax/rates")
		assertRoute(t, registered, "GET", "/tax/rates/:id")
		assertRoute(t, registered, "PUT", "/tax/rates/:id")
		assertRoute(t, registered, "DELETE", "/tax/rates/:id")
	})

	t.Run("Payments", func(t *testing.T) {
		assertRoute(t, registered, "POST", "/payments")
		assertRoute(t, registered, "GET", "/payments/:id")
		assertRoute(t, registered, "POST", "/payments/:id/process")
		assertRoute(t, registered, "POST", "/payments/:id/refund")
	})

	t.Run("CustomerGroups", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/customer-groups")
		assertRoute(t, registered, "POST", "/customer-groups")
		assertRoute(t, registered, "GET", "/customer-groups/:id")
		assertRoute(t, registered, "PUT", "/customer-groups/:id")
		assertRoute(t, registered, "DELETE", "/customer-groups/:id")
	})

	t.Run("GiftCards", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/admin/gift-cards")
		assertRoute(t, registered, "POST", "/admin/gift-cards")
		assertRoute(t, registered, "GET", "/admin/gift-cards/:id")
		assertRoute(t, registered, "PUT", "/admin/gift-cards/:id")
		assertRoute(t, registered, "POST", "/admin/gift-cards/:id/disable")
	})

	t.Run("Customers", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/customers")
		assertRoute(t, registered, "POST", "/customers")
		assertRoute(t, registered, "GET", "/customers/:id")
		assertRoute(t, registered, "PUT", "/customers/:id")
		assertRoute(t, registered, "DELETE", "/customers/:id")
	})

	t.Run("Currency", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/currencies")
		assertRoute(t, registered, "GET", "/currency-rates")
		assertRoute(t, registered, "POST", "/currency-rates")
		assertRoute(t, registered, "DELETE", "/currency-rates/:id")
	})

	t.Run("ImportExport", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/import-export/products/export")
		assertRoute(t, registered, "GET", "/import-export/orders/export")
		assertRoute(t, registered, "GET", "/import-export/customers/export")
		assertRoute(t, registered, "POST", "/import-export/products/import")
	})

	t.Run("Subscriptions", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/subscriptions")
		assertRoute(t, registered, "POST", "/subscriptions")
		assertRoute(t, registered, "GET", "/subscriptions/:id")
		assertRoute(t, registered, "POST", "/subscriptions/:id/cancel")
		assertRoute(t, registered, "POST", "/subscriptions/:id/pause")
		assertRoute(t, registered, "POST", "/subscriptions/:id/resume")
	})

	t.Run("Inventory", func(t *testing.T) {
		assertRoute(t, registered, "POST", "/inventory/warehouses")
		assertRoute(t, registered, "GET", "/inventory/warehouses")
		assertRoute(t, registered, "GET", "/inventory/warehouses/:id")
		assertRoute(t, registered, "PUT", "/inventory/warehouses/:id")
		assertRoute(t, registered, "DELETE", "/inventory/warehouses/:id")
		assertRoute(t, registered, "GET", "/inventory/stock/:productId")
		assertRoute(t, registered, "POST", "/inventory/stock/adjust")
		assertRoute(t, registered, "GET", "/inventory/stock/low")
	})

	t.Run("Reviews", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/reviews")
		assertRoute(t, registered, "GET", "/reviews/:id")
		assertRoute(t, registered, "PUT", "/reviews/:id/approve")
		assertRoute(t, registered, "PUT", "/reviews/:id/reject")
	})

	t.Run("Returns", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/returns")
		assertRoute(t, registered, "POST", "/returns")
		assertRoute(t, registered, "GET", "/returns/:id")
		assertRoute(t, registered, "PUT", "/returns/:id/approve")
		assertRoute(t, registered, "PUT", "/returns/:id/reject")
	})

	t.Run("Webhooks", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/webhooks")
		assertRoute(t, registered, "POST", "/webhooks")
		assertRoute(t, registered, "GET", "/webhooks/:id")
		assertRoute(t, registered, "PUT", "/webhooks/:id")
		assertRoute(t, registered, "DELETE", "/webhooks/:id")
	})

	t.Run("Audit", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/audit")
	})

	t.Run("Loyalty", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/admin/loyalty/rewards")
		assertRoute(t, registered, "POST", "/admin/loyalty/rewards")
		assertRoute(t, registered, "PUT", "/admin/loyalty/rewards/:id")
		assertRoute(t, registered, "GET", "/admin/loyalty/accounts")
		assertRoute(t, registered, "GET", "/admin/loyalty/accounts/:id")
		assertRoute(t, registered, "POST", "/admin/loyalty/accounts/:id/adjust")
		assertRoute(t, registered, "GET", "/admin/loyalty/accounts/:id/transactions")
	})

	t.Run("Bundles", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/bundles")
		assertRoute(t, registered, "POST", "/bundles")
		assertRoute(t, registered, "GET", "/bundles/:id")
		assertRoute(t, registered, "PUT", "/bundles/:id")
		assertRoute(t, registered, "DELETE", "/bundles/:id")
	})

	t.Run("Notifications", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/notifications")
		assertRoute(t, registered, "PUT", "/notifications/:id/read")
		assertRoute(t, registered, "DELETE", "/notifications/:id")
	})

	t.Run("Storefronts", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/storefronts")
		assertRoute(t, registered, "POST", "/storefronts")
		assertRoute(t, registered, "GET", "/storefronts/:id")
		assertRoute(t, registered, "PUT", "/storefronts/:id")
		assertRoute(t, registered, "DELETE", "/storefronts/:id")
	})

	t.Run("BulkOperations", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/bulk-operations")
		assertRoute(t, registered, "POST", "/bulk-operations")
		assertRoute(t, registered, "GET", "/bulk-operations/:id")
	})

	t.Run("Blog", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/blog/posts")
		assertRoute(t, registered, "POST", "/blog/posts")
		assertRoute(t, registered, "GET", "/blog/posts/:id")
		assertRoute(t, registered, "PUT", "/blog/posts/:id")
		assertRoute(t, registered, "DELETE", "/blog/posts/:id")
		assertRoute(t, registered, "PUT", "/blog/posts/:id/publish")
		assertRoute(t, registered, "PUT", "/blog/posts/:id/archive")
		assertRoute(t, registered, "GET", "/blog/categories")
		assertRoute(t, registered, "POST", "/blog/categories")
		assertRoute(t, registered, "PUT", "/blog/categories/:id")
		assertRoute(t, registered, "DELETE", "/blog/categories/:id")
	})

	t.Run("Experiments", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/experiments")
		assertRoute(t, registered, "POST", "/experiments")
		assertRoute(t, registered, "GET", "/experiments/:id")
		assertRoute(t, registered, "PUT", "/experiments/:id")
		assertRoute(t, registered, "DELETE", "/experiments/:id")
		assertRoute(t, registered, "PUT", "/experiments/:id/start")
		assertRoute(t, registered, "PUT", "/experiments/:id/pause")
		assertRoute(t, registered, "PUT", "/experiments/:id/complete")
	})

	t.Run("Recommendations", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/recommendations/rules")
		assertRoute(t, registered, "POST", "/recommendations/rules")
		assertRoute(t, registered, "PUT", "/recommendations/rules/:id")
		assertRoute(t, registered, "DELETE", "/recommendations/rules/:id")
	})

	t.Run("Plugins", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/plugins")
		assertRoute(t, registered, "GET", "/plugins/installed")
		assertRoute(t, registered, "POST", "/plugins/install")
		assertRoute(t, registered, "GET", "/plugins/:id")
		assertRoute(t, registered, "GET", "/plugins/:id/manifest")
		assertRoute(t, registered, "POST", "/plugins/:id/uninstall")
		assertRoute(t, registered, "PUT", "/plugins/:id/settings")
		assertRoute(t, registered, "POST", "/plugins/:id/enable")
		assertRoute(t, registered, "POST", "/plugins/:id/disable")
	})

	t.Run("Approvals", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/approvals")
		assertRoute(t, registered, "GET", "/approvals/:id")
		assertRoute(t, registered, "POST", "/approvals/:id/approve")
		assertRoute(t, registered, "POST", "/approvals/:id/reject")
	})

	t.Run("CartRecovery", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/cart-recovery")
		assertRoute(t, registered, "GET", "/cart-recovery/stats")
	})

	t.Run("I18n", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/i18n/supported-locales")
		assertRoute(t, registered, "GET", "/i18n/:entityType/:entityId/:locale")
		assertRoute(t, registered, "PUT", "/i18n/:entityType/:entityId/:locale")
	})

	t.Run("AgentSessions", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/agent/sessions")
		assertRoute(t, registered, "POST", "/agent/sessions")
		assertRoute(t, registered, "GET", "/agent/sessions/:id")
		assertRoute(t, registered, "DELETE", "/agent/sessions/:id")
		assertRoute(t, registered, "POST", "/agent/sessions/:id/stop")
		assertRoute(t, registered, "GET", "/agent/sessions/:id/messages")
		assertRoute(t, registered, "POST", "/agent/sessions/:id/messages")
	})

	t.Run("AgentMemory", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/agent/memories")
		assertRoute(t, registered, "POST", "/agent/memories")
		assertRoute(t, registered, "GET", "/agent/memories/:id")
		assertRoute(t, registered, "PUT", "/agent/memories/:id")
		assertRoute(t, registered, "DELETE", "/agent/memories/:id")
		assertRoute(t, registered, "GET", "/agent/memories/search")
	})

	t.Run("AgentTriggers", func(t *testing.T) {
		assertRoute(t, registered, "GET", "/agent/triggers")
		assertRoute(t, registered, "POST", "/agent/triggers")
		assertRoute(t, registered, "GET", "/agent/triggers/:id")
		assertRoute(t, registered, "PUT", "/agent/triggers/:id")
		assertRoute(t, registered, "DELETE", "/agent/triggers/:id")
		assertRoute(t, registered, "POST", "/agent/triggers/:id/enable")
		assertRoute(t, registered, "POST", "/agent/triggers/:id/disable")
		assertRoute(t, registered, "GET", "/agent/triggers/:id/logs")
		assertRoute(t, registered, "GET", "/agent/triggers/event-types")
	})
}
