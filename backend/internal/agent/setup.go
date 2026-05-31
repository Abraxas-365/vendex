// Package agent wires all domain services into harness-compatible tools
// and provides an EventHandler for streaming agent output to the admin UI.
package agent

import (
	"context"
	"encoding/json"

	"github.com/Abraxas-365/hada-commerce/internal/abtest/abtestsrv"
	"github.com/Abraxas-365/hada-commerce/internal/audit/auditsrv"
	"github.com/Abraxas-365/hada-commerce/internal/blog/blogsrv"
	"github.com/Abraxas-365/hada-commerce/internal/bulkops/bulkopssrv"
	"github.com/Abraxas-365/hada-commerce/internal/bundle/bundlesrv"
	"github.com/Abraxas-365/hada-commerce/internal/cartrecovery/cartrecoverysrv"
	"github.com/Abraxas-365/hada-commerce/internal/catalog/catalogsrv"
	"github.com/Abraxas-365/hada-commerce/internal/collection/collectionsrv"
	"github.com/Abraxas-365/hada-commerce/internal/currency/currencysrv"
	"github.com/Abraxas-365/hada-commerce/internal/customergroup/customergroupsrv"
	"github.com/Abraxas-365/hada-commerce/internal/dashboard/dashboardsrv"
	"github.com/Abraxas-365/hada-commerce/internal/giftcard/giftcardsrv"
	"github.com/Abraxas-365/hada-commerce/internal/i18n/i18nsrv"
	"github.com/Abraxas-365/hada-commerce/internal/inventory/inventorysrv"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/loyalty/loyaltysrv"
	"github.com/Abraxas-365/hada-commerce/internal/multistore/multistoresrv"
	"github.com/Abraxas-365/hada-commerce/internal/notification/notificationsrv"
	"github.com/Abraxas-365/hada-commerce/internal/order/ordersrv"
	"github.com/Abraxas-365/hada-commerce/internal/payment/paymentsrv"
	"github.com/Abraxas-365/hada-commerce/internal/product/productsrv"
	"github.com/Abraxas-365/hada-commerce/internal/recommendation/recommendationsrv"
	"github.com/Abraxas-365/hada-commerce/internal/promo/promosrv"
	"github.com/Abraxas-365/hada-commerce/internal/returns/returnssrv"
	"github.com/Abraxas-365/hada-commerce/internal/review/reviewsrv"
	"github.com/Abraxas-365/hada-commerce/internal/search/searchsrv"
	"github.com/Abraxas-365/hada-commerce/internal/shipping/shippingsrv"
	"github.com/Abraxas-365/hada-commerce/internal/storefront/storefrontsrv"
	"github.com/Abraxas-365/hada-commerce/internal/subscription/subscriptionsrv"
	"github.com/Abraxas-365/hada-commerce/internal/tax/taxsrv"
	"github.com/Abraxas-365/hada-commerce/internal/theme/themesrv"
	"github.com/Abraxas-365/hada-commerce/internal/webhook/webhooksrv"
)

// Tool is the minimal interface that each agent tool must satisfy.
// cmd/main.go adapts these into harness.Tool when wiring harness.New().
type Tool interface {
	// Name returns the tool's unique identifier used in API calls.
	Name() string
	// Description is shown to the LLM to help it decide when to use the tool.
	Description() string
	// InputSchema returns the JSON Schema describing the tool's input object.
	InputSchema() map[string]any
	// Execute runs the tool with the given JSON input and returns a
	// human-readable result string. Both tool-level errors (wrong input,
	// service failures) and framework errors are returned as Go errors.
	Execute(ctx context.Context, input json.RawMessage) (string, error)
}

// Services bundles all domain services the agent tools need.
type Services struct {
	Storefront     *storefrontsrv.Service
	Products       *productsrv.Service
	Orders         *ordersrv.Service
	Promos         *promosrv.Service
	Catalog        *catalogsrv.Service
	Themes         *themesrv.Service
	Shipping       *shippingsrv.Service
	Tax            *taxsrv.Service
	Payment        *paymentsrv.Service
	Search         *searchsrv.Service
	CustomerGroups *customergroupsrv.Service
	GiftCards      *giftcardsrv.Service
	CartRecovery   *cartrecoverysrv.Service
	Currency       *currencysrv.Service
	I18n           *i18nsrv.Service
	Subscriptions  *subscriptionsrv.Service
	Inventory      *inventorysrv.Service
	Reviews        *reviewsrv.Service
	Returns        *returnssrv.Service
	Webhooks       *webhooksrv.Service
	Audit          *auditsrv.Service
	Loyalty        *loyaltysrv.Service
	Bundles        *bundlesrv.Service
	Dashboard      *dashboardsrv.Service
	Notifications   *notificationsrv.Service
	MultiStore      *multistoresrv.Service
	BulkOps         *bulkopssrv.Service
	Blog            *blogsrv.Service
	Collections     *collectionsrv.Service
	ABTest          *abtestsrv.Service
	Recommendations *recommendationsrv.Service
}

// Setup constructs and returns all agent tools wired to the given services.
// Every tool operates on behalf of tenantID; pass the store's tenant.
//
// The returned slice is ready to be adapted into harness tools in cmd/main.go.
// The ordering is stable so callers can safely index into it.
func Setup(tenantID kernel.TenantID, svc Services) []Tool {
	sf := svc.Storefront
	return []Tool{
		// Storefront / CMS
		&CreatePageTool{sf: sf, tenantID: tenantID},
		&UpdatePageTool{sf: sf, tenantID: tenantID},
		&ListPagesTool{sf: sf, tenantID: tenantID},

		// Block types (storefront editor)
		&ListBlockTypesTool{sf: sf},
		&CreateBlockTypeTool{sf: sf},

		// Products
		&CreateProductTool{products: svc.Products, tenantID: tenantID},
		&ListProductsTool{products: svc.Products, tenantID: tenantID},

		// Promos
		&CreatePromoTool{promos: svc.Promos, tenantID: tenantID},

		// Orders (read-only query)
		&QueryOrdersTool{orders: svc.Orders, tenantID: tenantID},

		// Catalog
		&SearchCatalogTool{catalog: svc.Catalog, tenantID: tenantID},

		// Themes
		&ListThemesTool{themes: svc.Themes, tenantID: tenantID},
		&GetActiveThemeTool{themes: svc.Themes, tenantID: tenantID},
		&CreateThemeTool{themes: svc.Themes, tenantID: tenantID},
		&UpdateThemeTool{themes: svc.Themes, tenantID: tenantID},
		&ActivateThemeTool{themes: svc.Themes, tenantID: tenantID},

		// Shipping
		&ListShippingZonesTool{shipping: svc.Shipping, tenantID: tenantID},
		&CreateShippingZoneTool{shipping: svc.Shipping, tenantID: tenantID},
		&CalculateShippingTool{shipping: svc.Shipping, tenantID: tenantID},

		// Tax
		&ListTaxRatesTool{tax: svc.Tax, tenantID: tenantID},
		&CreateTaxRateTool{tax: svc.Tax, tenantID: tenantID},
		&CalculateTaxTool{tax: svc.Tax, tenantID: tenantID},

		// Payments
		&GetOrderPaymentTool{payment: svc.Payment, tenantID: tenantID},
		&ListRefundsTool{payment: svc.Payment, tenantID: tenantID},

		// Search
		&SearchProductsTool{search: svc.Search, tenantID: tenantID},
		&SearchSuggestionsTool{search: svc.Search, tenantID: tenantID},

		// Product variants
		&CreateProductOptionTool{products: svc.Products, tenantID: tenantID},
		&ListProductOptionsTool{products: svc.Products, tenantID: tenantID},
		&CreateProductVariantTool{products: svc.Products, tenantID: tenantID},
		&ListProductVariantsTool{products: svc.Products, tenantID: tenantID},

		// Customer groups
		&ListCustomerGroupsTool{groups: svc.CustomerGroups, tenantID: tenantID},
		&CreateCustomerGroupTool{groups: svc.CustomerGroups, tenantID: tenantID},
		&AddGroupMemberTool{groups: svc.CustomerGroups, tenantID: tenantID},

		// Gift cards
		&ListGiftCardsTool{giftcards: svc.GiftCards, tenantID: tenantID},
		&CreateGiftCardTool{giftcards: svc.GiftCards, tenantID: tenantID},
		&CheckGiftCardBalanceTool{giftcards: svc.GiftCards, tenantID: tenantID},
		&RedeemGiftCardTool{giftcards: svc.GiftCards, tenantID: tenantID},

		// Cart recovery
		&ListRecoveryEmailsTool{recovery: svc.CartRecovery, tenantID: tenantID},
		&GetRecoveryStatsTool{recovery: svc.CartRecovery, tenantID: tenantID},

		// Currency
		&ListCurrencyRatesTool{currency: svc.Currency, tenantID: tenantID},
		&SetCurrencyRateTool{currency: svc.Currency, tenantID: tenantID},
		&ConvertCurrencyTool{currency: svc.Currency, tenantID: tenantID},

		// I18n
		&SetTranslationsTool{i18n: svc.I18n, tenantID: tenantID},
		&GetTranslationsTool{i18n: svc.I18n, tenantID: tenantID},
		&ListSupportedLocalesTool{i18n: svc.I18n},

		// Subscriptions
		&ListSubscriptionsTool{subscriptions: svc.Subscriptions, tenantID: tenantID},
		&CreateSubscriptionTool{subscriptions: svc.Subscriptions, tenantID: tenantID},
		&CancelSubscriptionTool{subscriptions: svc.Subscriptions, tenantID: tenantID},

		// Inventory
		&ListWarehousesTool{inventory: svc.Inventory, tenantID: tenantID},
		&CreateWarehouseTool{inventory: svc.Inventory, tenantID: tenantID},
		&AdjustStockTool{inventory: svc.Inventory, tenantID: tenantID},
		&GetLowStockTool{inventory: svc.Inventory, tenantID: tenantID},

		// Reviews
		&ListReviewsTool{reviews: svc.Reviews, tenantID: tenantID},
		&ApproveReviewTool{reviews: svc.Reviews, tenantID: tenantID},
		&RejectReviewTool{reviews: svc.Reviews, tenantID: tenantID},

		// Returns
		&ListReturnsTool{returns: svc.Returns, tenantID: tenantID},
		&ApproveReturnTool{returns: svc.Returns, tenantID: tenantID},

		// Webhooks
		&ListWebhooksTool{webhooks: svc.Webhooks, tenantID: tenantID},
		&CreateWebhookTool{webhooks: svc.Webhooks, tenantID: tenantID},
		&ToggleWebhookTool{webhooks: svc.Webhooks, tenantID: tenantID},

		// Audit
		&ListAuditLogsTool{audit: svc.Audit, tenantID: tenantID},

		// Loyalty
		&ListLoyaltyAccountsTool{loyalty: svc.Loyalty, tenantID: tenantID},
		&EarnLoyaltyPointsTool{loyalty: svc.Loyalty, tenantID: tenantID},
		&ListLoyaltyRewardsTool{loyalty: svc.Loyalty, tenantID: tenantID},
		&CreateLoyaltyRewardTool{loyalty: svc.Loyalty, tenantID: tenantID},

		// Bundles
		&ListBundlesTool{bundles: svc.Bundles, tenantID: tenantID},
		&CreateBundleTool{bundles: svc.Bundles, tenantID: tenantID},

		// Dashboard
		&GetSalesOverviewTool{dashboard: svc.Dashboard, tenantID: tenantID},
		&GetTopProductsTool{dashboard: svc.Dashboard, tenantID: tenantID},
		&GetRevenueByDayTool{dashboard: svc.Dashboard, tenantID: tenantID},

		// Notifications
		&GetUnreadNotificationCountTool{notifications: svc.Notifications, tenantID: tenantID},
		&MarkAllNotificationsReadTool{notifications: svc.Notifications, tenantID: tenantID},

		// Multi-storefront
		&ListStorefrontsTool{multistore: svc.MultiStore, tenantID: tenantID},
		&CreateStorefrontTool{multistore: svc.MultiStore, tenantID: tenantID},

		// Bulk operations
		&ListBulkOperationsTool{bulkops: svc.BulkOps, tenantID: tenantID},
		&CreateBulkOperationTool{bulkops: svc.BulkOps, tenantID: tenantID},

		// Blog
		&ListBlogPostsTool{blog: svc.Blog, tenantID: tenantID},
		&CreateBlogPostTool{blog: svc.Blog, tenantID: tenantID},
		&PublishBlogPostTool{blog: svc.Blog, tenantID: tenantID},

		// Collections
		&ListCollectionsTool{collections: svc.Collections, tenantID: tenantID},
		&CreateCollectionTool{collections: svc.Collections, tenantID: tenantID},
		&AddCollectionProductTool{collections: svc.Collections, tenantID: tenantID},

		// A/B Testing
		&ListExperimentsTool{abtest: svc.ABTest, tenantID: tenantID},
		&CreateExperimentTool{abtest: svc.ABTest, tenantID: tenantID},
		&GetExperimentResultsTool{abtest: svc.ABTest, tenantID: tenantID},

		// Recommendations
		&ListRecommendationRulesTool{recs: svc.Recommendations, tenantID: tenantID},
		&GetTrendingProductsTool{recs: svc.Recommendations, tenantID: tenantID},
	}
}
