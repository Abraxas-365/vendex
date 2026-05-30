// Package agent wires all domain services into harness-compatible tools
// and provides an EventHandler for streaming agent output to the admin UI.
package agent

import (
	"context"
	"encoding/json"

	"github.com/Abraxas-365/hada-commerce/internal/catalog/catalogsrv"
	"github.com/Abraxas-365/hada-commerce/internal/customergroup/customergroupsrv"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/order/ordersrv"
	"github.com/Abraxas-365/hada-commerce/internal/payment/paymentsrv"
	"github.com/Abraxas-365/hada-commerce/internal/product/productsrv"
	"github.com/Abraxas-365/hada-commerce/internal/promo/promosrv"
	"github.com/Abraxas-365/hada-commerce/internal/search/searchsrv"
	"github.com/Abraxas-365/hada-commerce/internal/shipping/shippingsrv"
	"github.com/Abraxas-365/hada-commerce/internal/storefront/storefrontsrv"
	"github.com/Abraxas-365/hada-commerce/internal/tax/taxsrv"
	"github.com/Abraxas-365/hada-commerce/internal/theme/themesrv"
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
	}
}
