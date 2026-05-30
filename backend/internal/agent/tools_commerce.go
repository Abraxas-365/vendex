package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Abraxas-365/hada-commerce/internal/cartrecovery/cartrecoverysrv"
	"github.com/Abraxas-365/hada-commerce/internal/currency/currencysrv"
	"github.com/Abraxas-365/hada-commerce/internal/customergroup"
	"github.com/Abraxas-365/hada-commerce/internal/customergroup/customergroupsrv"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/giftcard"
	"github.com/Abraxas-365/hada-commerce/internal/giftcard/giftcardsrv"
	"github.com/Abraxas-365/hada-commerce/internal/i18n/i18nsrv"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/payment/paymentsrv"
	"github.com/Abraxas-365/hada-commerce/internal/product/productsrv"
	"github.com/Abraxas-365/hada-commerce/internal/search"
	"github.com/Abraxas-365/hada-commerce/internal/search/searchsrv"
	"github.com/Abraxas-365/hada-commerce/internal/shipping/shippingsrv"
	"github.com/Abraxas-365/hada-commerce/internal/subscription"
	"github.com/Abraxas-365/hada-commerce/internal/subscription/subscriptionsrv"
	"github.com/Abraxas-365/hada-commerce/internal/tax/taxsrv"
)

// ─── ListShippingZonesTool ────────────────────────────────────────────────────

type ListShippingZonesTool struct {
	shipping *shippingsrv.Service
	tenantID kernel.TenantID
}

func (t *ListShippingZonesTool) Name() string { return "list_shipping_zones" }

func (t *ListShippingZonesTool) Description() string {
	return "List all shipping zones for the store, showing countries and states each zone covers."
}

func (t *ListShippingZonesTool) InputSchema() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

func (t *ListShippingZonesTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	zones, err := t.shipping.ListZones(ctx, t.tenantID)
	if err != nil {
		return "", errx.Wrap(err, "list_shipping_zones", errx.TypeInternal)
	}
	if len(zones) == 0 {
		return "No shipping zones found.", nil
	}
	out := fmt.Sprintf("Shipping zones (%d):\n\n", len(zones))
	for _, z := range zones {
		out += fmt.Sprintf("- ID: %s | Name: %s | Countries: %v | States: %v\n",
			z.ID, z.Name, z.Countries, z.States)
	}
	return out, nil
}

// ─── CreateShippingZoneTool ───────────────────────────────────────────────────

type CreateShippingZoneTool struct {
	shipping *shippingsrv.Service
	tenantID kernel.TenantID
}

type createShippingZoneInput struct {
	Name      string   `json:"name"`
	Countries []string `json:"countries"`
	States    []string `json:"states"`
}

func (t *CreateShippingZoneTool) Name() string { return "create_shipping_zone" }

func (t *CreateShippingZoneTool) Description() string {
	return "Create a new shipping zone for the store. A zone groups countries/states for rate assignment."
}

func (t *CreateShippingZoneTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":      map[string]any{"type": "string", "description": "Zone name, e.g. 'Domestic' or 'EU'"},
			"countries": map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "ISO country codes, e.g. ['US','CA']"},
			"states":    map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "State/province codes (optional)"},
		},
		"required": []string{"name", "countries"},
	}
}

func (t *CreateShippingZoneTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in createShippingZoneInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "create_shipping_zone: unmarshal input", errx.TypeValidation)
	}
	zone, err := t.shipping.CreateZone(ctx, t.tenantID, in.Name, in.Countries, in.States)
	if err != nil {
		return "", errx.Wrap(err, "create_shipping_zone", errx.TypeInternal)
	}
	return fmt.Sprintf("Shipping zone created.\nID: %s\nName: %s\nCountries: %v\nStates: %v",
		zone.ID, zone.Name, zone.Countries, zone.States), nil
}

// ─── CalculateShippingTool ────────────────────────────────────────────────────

type CalculateShippingTool struct {
	shipping *shippingsrv.Service
	tenantID kernel.TenantID
}

type calculateShippingInput struct {
	Country     string  `json:"country"`
	State       string  `json:"state"`
	OrderAmount int64   `json:"order_amount"`
	Weight      float64 `json:"weight"`
}

func (t *CalculateShippingTool) Name() string { return "calculate_shipping" }

func (t *CalculateShippingTool) Description() string {
	return "Calculate available shipping rates for a given destination and order. Returns matching rates with prices."
}

func (t *CalculateShippingTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"country":      map[string]any{"type": "string", "description": "ISO country code, e.g. 'US'"},
			"state":        map[string]any{"type": "string", "description": "State/province code (optional)"},
			"order_amount": map[string]any{"type": "integer", "description": "Order subtotal in cents"},
			"weight":       map[string]any{"type": "number", "description": "Total weight in kg"},
		},
		"required": []string{"country", "order_amount"},
	}
}

func (t *CalculateShippingTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in calculateShippingInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "calculate_shipping: unmarshal input", errx.TypeValidation)
	}
	rates, err := t.shipping.CalculateShipping(ctx, t.tenantID, in.Country, in.State, in.OrderAmount, in.Weight)
	if err != nil {
		return "", errx.Wrap(err, "calculate_shipping", errx.TypeInternal)
	}
	if len(rates) == 0 {
		return "No shipping rates available for this destination.", nil
	}
	out := fmt.Sprintf("Available shipping rates (%d):\n\n", len(rates))
	for _, r := range rates {
		out += fmt.Sprintf("- %s | %d %s", r.Name, r.Price.Amount, r.Price.Currency)
		if r.EstDaysMin != nil && r.EstDaysMax != nil {
			out += fmt.Sprintf(" | %d-%d days", *r.EstDaysMin, *r.EstDaysMax)
		}
		out += "\n"
	}
	return out, nil
}

// ─── ListTaxRatesTool ─────────────────────────────────────────────────────────

type ListTaxRatesTool struct {
	tax      *taxsrv.Service
	tenantID kernel.TenantID
}

func (t *ListTaxRatesTool) Name() string { return "list_tax_rates" }

func (t *ListTaxRatesTool) Description() string {
	return "List all tax rates configured for the store."
}

func (t *ListTaxRatesTool) InputSchema() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

func (t *ListTaxRatesTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	rates, err := t.tax.ListRates(ctx, t.tenantID)
	if err != nil {
		return "", errx.Wrap(err, "list_tax_rates", errx.TypeInternal)
	}
	if len(rates) == 0 {
		return "No tax rates configured.", nil
	}
	out := fmt.Sprintf("Tax rates (%d):\n\n", len(rates))
	for _, r := range rates {
		compound := ""
		if r.Compound {
			compound = " [compound]"
		}
		out += fmt.Sprintf("- ID: %s | %s: %.2f%%%s | %s %s %s | Active: %v\n",
			r.ID, r.Name, r.Rate*100, compound, r.Country, r.State, r.City, r.Active)
	}
	return out, nil
}

// ─── CreateTaxRateTool ────────────────────────────────────────────────────────

type CreateTaxRateTool struct {
	tax      *taxsrv.Service
	tenantID kernel.TenantID
}

type createTaxRateInput struct {
	Name             string  `json:"name"`
	Rate             float64 `json:"rate"`
	Country          string  `json:"country"`
	State            string  `json:"state"`
	City             string  `json:"city"`
	ZipCode          string  `json:"zip_code"`
	Priority         int     `json:"priority"`
	Compound         bool    `json:"compound"`
	IncludesShipping bool    `json:"includes_shipping"`
	Active           bool    `json:"active"`
}

func (t *CreateTaxRateTool) Name() string { return "create_tax_rate" }

func (t *CreateTaxRateTool) Description() string {
	return "Create a new tax rate for a specific region (country/state/city). Rate is a decimal fraction, e.g. 0.08 for 8%."
}

func (t *CreateTaxRateTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":              map[string]any{"type": "string", "description": "Tax rate name, e.g. 'State Sales Tax'"},
			"rate":              map[string]any{"type": "number", "description": "Tax rate as decimal fraction, e.g. 0.08 for 8%"},
			"country":           map[string]any{"type": "string", "description": "ISO country code"},
			"state":             map[string]any{"type": "string", "description": "State/province code (optional)"},
			"city":              map[string]any{"type": "string", "description": "City name (optional)"},
			"zip_code":          map[string]any{"type": "string", "description": "ZIP/postal code (optional)"},
			"priority":          map[string]any{"type": "integer", "description": "Evaluation priority (lower = first)"},
			"compound":          map[string]any{"type": "boolean", "description": "If true, tax compounds on previous tax amounts"},
			"includes_shipping": map[string]any{"type": "boolean", "description": "If true, tax applies to shipping charges too"},
			"active":            map[string]any{"type": "boolean", "description": "Whether this rate is active"},
		},
		"required": []string{"name", "rate", "country"},
	}
}

func (t *CreateTaxRateTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in createTaxRateInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "create_tax_rate: unmarshal input", errx.TypeValidation)
	}
	rate, err := t.tax.CreateRate(ctx, t.tenantID, taxsrv.CreateRateInput{
		Name:             in.Name,
		Rate:             in.Rate,
		Country:          in.Country,
		State:            in.State,
		City:             in.City,
		ZipCode:          in.ZipCode,
		Priority:         in.Priority,
		Compound:         in.Compound,
		IncludesShipping: in.IncludesShipping,
		Active:           in.Active,
	})
	if err != nil {
		return "", errx.Wrap(err, "create_tax_rate", errx.TypeInternal)
	}
	return fmt.Sprintf("Tax rate created.\nID: %s\nName: %s\nRate: %.2f%%\nCountry: %s\nState: %s\nActive: %v",
		rate.ID, rate.Name, rate.Rate*100, rate.Country, rate.State, rate.Active), nil
}

// ─── CalculateTaxTool ─────────────────────────────────────────────────────────

type CalculateTaxTool struct {
	tax      *taxsrv.Service
	tenantID kernel.TenantID
}

type calculateTaxInput struct {
	SubtotalCents int64  `json:"subtotal_cents"`
	ShippingCents int64  `json:"shipping_cents"`
	Country       string `json:"country"`
	State         string `json:"state"`
	City          string `json:"city"`
	ZipCode       string `json:"zip_code"`
}

func (t *CalculateTaxTool) Name() string { return "calculate_tax" }

func (t *CalculateTaxTool) Description() string {
	return "Calculate tax for an order given subtotal, shipping, and destination address."
}

func (t *CalculateTaxTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"subtotal_cents": map[string]any{"type": "integer", "description": "Order subtotal in cents"},
			"shipping_cents": map[string]any{"type": "integer", "description": "Shipping cost in cents"},
			"country":        map[string]any{"type": "string", "description": "ISO country code"},
			"state":          map[string]any{"type": "string", "description": "State/province code"},
			"city":           map[string]any{"type": "string", "description": "City name (optional)"},
			"zip_code":       map[string]any{"type": "string", "description": "ZIP/postal code (optional)"},
		},
		"required": []string{"subtotal_cents", "country"},
	}
}

func (t *CalculateTaxTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in calculateTaxInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "calculate_tax: unmarshal input", errx.TypeValidation)
	}
	result, err := t.tax.CalculateTax(ctx, t.tenantID, in.SubtotalCents, in.ShippingCents,
		in.Country, in.State, in.City, in.ZipCode)
	if err != nil {
		return "", errx.Wrap(err, "calculate_tax", errx.TypeInternal)
	}
	out := fmt.Sprintf("Tax calculation:\nTotal tax: %d cents\n\nBreakdown:\n", result.TotalTax)
	for _, item := range result.TaxBreakdown {
		out += fmt.Sprintf("- %s: %d cents (rate: %.2f%%)\n", item.Name, item.Amount, item.Rate*100)
	}
	return out, nil
}

// ─── GetOrderPaymentTool ──────────────────────────────────────────────────────

type GetOrderPaymentTool struct {
	payment  *paymentsrv.Service
	tenantID kernel.TenantID
}

type getOrderPaymentInput struct {
	OrderID string `json:"order_id"`
}

func (t *GetOrderPaymentTool) Name() string { return "get_order_payment" }

func (t *GetOrderPaymentTool) Description() string {
	return "Get payment details for a specific order."
}

func (t *GetOrderPaymentTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"order_id": map[string]any{"type": "string", "description": "The order ID"},
		},
		"required": []string{"order_id"},
	}
}

func (t *GetOrderPaymentTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in getOrderPaymentInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "get_order_payment: unmarshal input", errx.TypeValidation)
	}
	p, err := t.payment.GetPaymentByOrder(ctx, t.tenantID, kernel.OrderID(in.OrderID))
	if err != nil {
		return "", errx.Wrap(err, "get_order_payment", errx.TypeInternal)
	}
	return fmt.Sprintf("Payment for order %s:\nID: %s\nStatus: %s\nAmount: %d %s\nProvider: %s\nMethod: %s\nCreated: %s",
		in.OrderID, p.ID, p.Status, p.Amount.Amount, p.Amount.Currency,
		p.Provider, p.Method, p.CreatedAt.Format("2006-01-02 15:04")), nil
}

// ─── ListRefundsTool ──────────────────────────────────────────────────────────

type ListRefundsTool struct {
	payment  *paymentsrv.Service
	tenantID kernel.TenantID
}

type listRefundsInput struct {
	PaymentID string `json:"payment_id"`
}

func (t *ListRefundsTool) Name() string { return "list_refunds" }

func (t *ListRefundsTool) Description() string {
	return "List all refunds for a specific payment."
}

func (t *ListRefundsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"payment_id": map[string]any{"type": "string", "description": "The payment ID"},
		},
		"required": []string{"payment_id"},
	}
}

func (t *ListRefundsTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in listRefundsInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "list_refunds: unmarshal input", errx.TypeValidation)
	}
	refunds, err := t.payment.ListRefunds(ctx, t.tenantID, kernel.PaymentID(in.PaymentID))
	if err != nil {
		return "", errx.Wrap(err, "list_refunds", errx.TypeInternal)
	}
	if len(refunds) == 0 {
		return "No refunds found for this payment.", nil
	}
	out := fmt.Sprintf("Refunds (%d):\n\n", len(refunds))
	for _, r := range refunds {
		out += fmt.Sprintf("- ID: %s | Amount: %d %s | Status: %s | Reason: %s | Created: %s\n",
			r.ID, r.Amount.Amount, r.Amount.Currency, r.Status, r.Reason,
			r.CreatedAt.Format("2006-01-02 15:04"))
	}
	return out, nil
}

// ─── SearchProductsTool ───────────────────────────────────────────────────────

type SearchProductsTool struct {
	search   *searchsrv.Service
	tenantID kernel.TenantID
}

type searchProductsInput struct {
	Query      string   `json:"query"`
	CategoryID string   `json:"category_id"`
	Tags       []string `json:"tags"`
	MinPrice   *int64   `json:"min_price"`
	MaxPrice   *int64   `json:"max_price"`
	SortBy     string   `json:"sort_by"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
}

func (t *SearchProductsTool) Name() string { return "search_products" }

func (t *SearchProductsTool) Description() string {
	return "Full-text search across products with optional filters (category, tags, price range)."
}

func (t *SearchProductsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"query":       map[string]any{"type": "string", "description": "Search query text"},
			"category_id": map[string]any{"type": "string", "description": "Filter by category ID (optional)"},
			"tags":        map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Filter by tags (optional)"},
			"min_price":   map[string]any{"type": "integer", "description": "Minimum price in cents (optional)"},
			"max_price":   map[string]any{"type": "integer", "description": "Maximum price in cents (optional)"},
			"sort_by":     map[string]any{"type": "string", "description": "Sort: relevance, price_asc, price_desc, name, created_at"},
			"page":        map[string]any{"type": "integer", "description": "Page number (default 1)"},
			"page_size":   map[string]any{"type": "integer", "description": "Results per page (default 20)"},
		},
		"required": []string{"query"},
	}
}

func (t *SearchProductsTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in searchProductsInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "search_products: unmarshal input", errx.TypeValidation)
	}
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 20
	}
	result, err := t.search.Search(ctx, t.tenantID, search.SearchQuery{
		Query:      in.Query,
		CategoryID: in.CategoryID,
		Tags:       in.Tags,
		MinPrice:   in.MinPrice,
		MaxPrice:   in.MaxPrice,
		SortBy:     in.SortBy,
		Page:       in.Page,
		PageSize:   in.PageSize,
	})
	if err != nil {
		return "", errx.Wrap(err, "search_products", errx.TypeInternal)
	}
	if len(result.Products) == 0 {
		return fmt.Sprintf("No products found for query '%s'.", in.Query), nil
	}
	out := fmt.Sprintf("Search results for '%s' (%d total, page %d/%d):\n\n",
		result.Query, result.Total, result.Page, result.TotalPages)
	for _, p := range result.Products {
		out += fmt.Sprintf("- %s | %d %s | SKU: %s | ID: %s\n",
			p.Name, p.Price.Amount, p.Price.Currency, p.SKU, p.ID)
	}
	return out, nil
}

// ─── SearchSuggestionsTool ────────────────────────────────────────────────────

type SearchSuggestionsTool struct {
	search   *searchsrv.Service
	tenantID kernel.TenantID
}

type searchSuggestionsInput struct {
	Prefix string `json:"prefix"`
	Limit  int    `json:"limit"`
}

func (t *SearchSuggestionsTool) Name() string { return "search_suggestions" }

func (t *SearchSuggestionsTool) Description() string {
	return "Get autocomplete suggestions for a partial search query."
}

func (t *SearchSuggestionsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"prefix": map[string]any{"type": "string", "description": "Partial search text to autocomplete"},
			"limit":  map[string]any{"type": "integer", "description": "Max suggestions to return (default 5)"},
		},
		"required": []string{"prefix"},
	}
}

func (t *SearchSuggestionsTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in searchSuggestionsInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "search_suggestions: unmarshal input", errx.TypeValidation)
	}
	if in.Limit <= 0 {
		in.Limit = 5
	}
	suggestions, err := t.search.Suggest(ctx, t.tenantID, in.Prefix, in.Limit)
	if err != nil {
		return "", errx.Wrap(err, "search_suggestions", errx.TypeInternal)
	}
	if len(suggestions) == 0 {
		return "No suggestions found.", nil
	}
	out := fmt.Sprintf("Suggestions for '%s':\n\n", in.Prefix)
	for _, s := range suggestions {
		out += fmt.Sprintf("- %s (%d results)\n", s.Term, s.Count)
	}
	return out, nil
}

// ─── CreateProductOptionTool ──────────────────────────────────────────────────

type CreateProductOptionTool struct {
	products *productsrv.Service
	tenantID kernel.TenantID
}

type createProductOptionInput struct {
	ProductID string   `json:"product_id"`
	Name      string   `json:"name"`
	Position  int      `json:"position"`
	Values    []string `json:"values"`
}

func (t *CreateProductOptionTool) Name() string { return "create_product_option" }

func (t *CreateProductOptionTool) Description() string {
	return "Create a product option (e.g. Size, Color) with its possible values."
}

func (t *CreateProductOptionTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"product_id": map[string]any{"type": "string", "description": "The product ID"},
			"name":       map[string]any{"type": "string", "description": "Option name, e.g. 'Size' or 'Color'"},
			"position":   map[string]any{"type": "integer", "description": "Display position (0-based)"},
			"values":     map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Possible values, e.g. ['S','M','L','XL']"},
		},
		"required": []string{"product_id", "name", "values"},
	}
}

func (t *CreateProductOptionTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in createProductOptionInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "create_product_option: unmarshal input", errx.TypeValidation)
	}
	opt, err := t.products.CreateOption(ctx, t.tenantID, kernel.ProductID(in.ProductID), in.Name, in.Position, in.Values)
	if err != nil {
		return "", errx.Wrap(err, "create_product_option", errx.TypeInternal)
	}
	return fmt.Sprintf("Product option created.\nID: %s\nName: %s\nValues: %v",
		opt.ID, opt.Name, opt.Values), nil
}

// ─── ListProductOptionsTool ───────────────────────────────────────────────────

type ListProductOptionsTool struct {
	products *productsrv.Service
	tenantID kernel.TenantID
}

type listProductOptionsInput struct {
	ProductID string `json:"product_id"`
}

func (t *ListProductOptionsTool) Name() string { return "list_product_options" }

func (t *ListProductOptionsTool) Description() string {
	return "List all options (e.g. Size, Color) for a product."
}

func (t *ListProductOptionsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"product_id": map[string]any{"type": "string", "description": "The product ID"},
		},
		"required": []string{"product_id"},
	}
}

func (t *ListProductOptionsTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in listProductOptionsInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "list_product_options: unmarshal input", errx.TypeValidation)
	}
	opts, err := t.products.ListOptions(ctx, t.tenantID, kernel.ProductID(in.ProductID))
	if err != nil {
		return "", errx.Wrap(err, "list_product_options", errx.TypeInternal)
	}
	if len(opts) == 0 {
		return "No options found for this product.", nil
	}
	out := fmt.Sprintf("Product options (%d):\n\n", len(opts))
	for _, o := range opts {
		out += fmt.Sprintf("- ID: %s | %s: %v (position %d)\n", o.ID, o.Name, o.Values, o.Position)
	}
	return out, nil
}

// ─── CreateProductVariantTool ─────────────────────────────────────────────────

type CreateProductVariantTool struct {
	products *productsrv.Service
	tenantID kernel.TenantID
}

type createProductVariantInput struct {
	ProductID string            `json:"product_id"`
	SKU       string            `json:"sku"`
	Price     int64             `json:"price"`
	Currency  string            `json:"currency"`
	Stock     int               `json:"stock"`
	Options   map[string]string `json:"options"`
}

func (t *CreateProductVariantTool) Name() string { return "create_product_variant" }

func (t *CreateProductVariantTool) Description() string {
	return "Create a product variant with specific option values, SKU, price, and stock."
}

func (t *CreateProductVariantTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"product_id": map[string]any{"type": "string", "description": "The product ID"},
			"sku":        map[string]any{"type": "string", "description": "Variant SKU"},
			"price":      map[string]any{"type": "integer", "description": "Variant price in cents"},
			"currency":   map[string]any{"type": "string", "description": "Currency code, e.g. 'USD'"},
			"stock":      map[string]any{"type": "integer", "description": "Stock quantity"},
			"options":    map[string]any{"type": "object", "description": "Option name→value map, e.g. {\"Size\":\"L\",\"Color\":\"Red\"}"},
		},
		"required": []string{"product_id", "sku", "price", "currency", "stock", "options"},
	}
}

func (t *CreateProductVariantTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in createProductVariantInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "create_product_variant: unmarshal input", errx.TypeValidation)
	}
	v, err := t.products.CreateVariant(ctx, t.tenantID, kernel.ProductID(in.ProductID),
		in.SKU, kernel.Money{Amount: in.Price, Currency: in.Currency}, in.Stock, in.Options)
	if err != nil {
		return "", errx.Wrap(err, "create_product_variant", errx.TypeInternal)
	}
	return fmt.Sprintf("Variant created.\nID: %s\nSKU: %s\nPrice: %d %s\nStock: %d\nOptions: %v",
		v.ID, v.SKU, v.Price.Amount, v.Price.Currency, v.Stock, v.Options), nil
}

// ─── ListProductVariantsTool ──────────────────────────────────────────────────

type ListProductVariantsTool struct {
	products *productsrv.Service
	tenantID kernel.TenantID
}

type listProductVariantsInput struct {
	ProductID string `json:"product_id"`
}

func (t *ListProductVariantsTool) Name() string { return "list_product_variants" }

func (t *ListProductVariantsTool) Description() string {
	return "List all variants for a product, showing SKU, price, stock, and option values."
}

func (t *ListProductVariantsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"product_id": map[string]any{"type": "string", "description": "The product ID"},
		},
		"required": []string{"product_id"},
	}
}

func (t *ListProductVariantsTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in listProductVariantsInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "list_product_variants: unmarshal input", errx.TypeValidation)
	}
	variants, err := t.products.ListVariants(ctx, t.tenantID, kernel.ProductID(in.ProductID))
	if err != nil {
		return "", errx.Wrap(err, "list_product_variants", errx.TypeInternal)
	}
	if len(variants) == 0 {
		return "No variants found for this product.", nil
	}
	out := fmt.Sprintf("Product variants (%d):\n\n", len(variants))
	for _, v := range variants {
		active := ""
		if !v.Active {
			active = " [inactive]"
		}
		out += fmt.Sprintf("- ID: %s | SKU: %s | %d %s | Stock: %d | %v%s\n",
			v.ID, v.SKU, v.Price.Amount, v.Price.Currency, v.Stock, v.Options, active)
	}
	return out, nil
}

// ─── ListCustomerGroupsTool ───────────────────────────────────────────────────

type ListCustomerGroupsTool struct {
	groups   *customergroupsrv.Service
	tenantID kernel.TenantID
}

func (t *ListCustomerGroupsTool) Name() string { return "list_customer_groups" }

func (t *ListCustomerGroupsTool) Description() string {
	return "List all customer groups/segments configured for the store."
}

func (t *ListCustomerGroupsTool) InputSchema() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

func (t *ListCustomerGroupsTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	groups, err := t.groups.List(ctx, t.tenantID)
	if err != nil {
		return "", errx.Wrap(err, "list_customer_groups", errx.TypeInternal)
	}
	if len(groups) == 0 {
		return "No customer groups found.", nil
	}
	out := fmt.Sprintf("Customer groups (%d):\n\n", len(groups))
	for _, g := range groups {
		auto := ""
		if g.AutoAssign {
			auto = " [auto-assign]"
		}
		out += fmt.Sprintf("- ID: %s | %s: %s%s\n", g.ID, g.Name, g.Description, auto)
	}
	return out, nil
}

// ─── CreateCustomerGroupTool ──────────────────────────────────────────────────

type CreateCustomerGroupTool struct {
	groups   *customergroupsrv.Service
	tenantID kernel.TenantID
}

type createCustomerGroupInput struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Rules       customergroup.GroupRules `json:"rules"`
	AutoAssign  bool                    `json:"auto_assign"`
}

func (t *CreateCustomerGroupTool) Name() string { return "create_customer_group" }

func (t *CreateCustomerGroupTool) Description() string {
	return "Create a customer group for segmentation, promo targeting, or tiered pricing."
}

func (t *CreateCustomerGroupTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":        map[string]any{"type": "string", "description": "Group name, e.g. 'VIP Customers'"},
			"description": map[string]any{"type": "string", "description": "Group description"},
			"rules": map[string]any{
				"type":        "object",
				"description": "Auto-qualification rules (optional): min_orders, min_spent, tags, etc.",
			},
			"auto_assign": map[string]any{"type": "boolean", "description": "If true, customers matching rules are auto-assigned"},
		},
		"required": []string{"name"},
	}
}

func (t *CreateCustomerGroupTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in createCustomerGroupInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "create_customer_group: unmarshal input", errx.TypeValidation)
	}
	group, err := t.groups.Create(ctx, t.tenantID, customergroup.CreateGroupRequest{
		Name:        in.Name,
		Description: in.Description,
		Rules:       in.Rules,
		AutoAssign:  in.AutoAssign,
	})
	if err != nil {
		return "", errx.Wrap(err, "create_customer_group", errx.TypeInternal)
	}
	return fmt.Sprintf("Customer group created.\nID: %s\nName: %s\nDescription: %s\nAuto-assign: %v",
		group.ID, group.Name, group.Description, group.AutoAssign), nil
}

// ─── AddGroupMemberTool ───────────────────────────────────────────────────────

type AddGroupMemberTool struct {
	groups   *customergroupsrv.Service
	tenantID kernel.TenantID
}

type addGroupMemberInput struct {
	GroupID    string `json:"group_id"`
	CustomerID string `json:"customer_id"`
}

func (t *AddGroupMemberTool) Name() string { return "add_group_member" }

func (t *AddGroupMemberTool) Description() string {
	return "Add a customer to a customer group."
}

func (t *AddGroupMemberTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"group_id":    map[string]any{"type": "string", "description": "Customer group ID"},
			"customer_id": map[string]any{"type": "string", "description": "Customer ID to add"},
		},
		"required": []string{"group_id", "customer_id"},
	}
}

func (t *AddGroupMemberTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in addGroupMemberInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "add_group_member: unmarshal input", errx.TypeValidation)
	}
	membership, err := t.groups.AddMember(ctx, t.tenantID,
		kernel.CustomerGroupID(in.GroupID), kernel.CustomerID(in.CustomerID))
	if err != nil {
		return "", errx.Wrap(err, "add_group_member", errx.TypeInternal)
	}
	return fmt.Sprintf("Member added.\nGroup: %s\nCustomer: %s\nJoined: %s",
		membership.GroupID, membership.CustomerID, membership.AssignedAt.Format("2006-01-02 15:04")), nil
}

// ─── compile-time guards ──────────────────────────────────────────────────────

var (
	_ Tool = (*ListShippingZonesTool)(nil)
	_ Tool = (*CreateShippingZoneTool)(nil)
	_ Tool = (*CalculateShippingTool)(nil)
	_ Tool = (*ListTaxRatesTool)(nil)
	_ Tool = (*CreateTaxRateTool)(nil)
	_ Tool = (*CalculateTaxTool)(nil)
	_ Tool = (*GetOrderPaymentTool)(nil)
	_ Tool = (*ListRefundsTool)(nil)
	_ Tool = (*SearchProductsTool)(nil)
	_ Tool = (*SearchSuggestionsTool)(nil)
	_ Tool = (*CreateProductOptionTool)(nil)
	_ Tool = (*ListProductOptionsTool)(nil)
	_ Tool = (*CreateProductVariantTool)(nil)
	_ Tool = (*ListProductVariantsTool)(nil)
	_ Tool = (*ListCustomerGroupsTool)(nil)
	_ Tool = (*CreateCustomerGroupTool)(nil)
	_ Tool = (*AddGroupMemberTool)(nil)

	// P3 tools
	_ Tool = (*ListGiftCardsTool)(nil)
	_ Tool = (*CreateGiftCardTool)(nil)
	_ Tool = (*CheckGiftCardBalanceTool)(nil)
	_ Tool = (*RedeemGiftCardTool)(nil)
	_ Tool = (*ListRecoveryEmailsTool)(nil)
	_ Tool = (*GetRecoveryStatsTool)(nil)
	_ Tool = (*ListCurrencyRatesTool)(nil)
	_ Tool = (*SetCurrencyRateTool)(nil)
	_ Tool = (*ConvertCurrencyTool)(nil)
	_ Tool = (*SetTranslationsTool)(nil)
	_ Tool = (*GetTranslationsTool)(nil)
	_ Tool = (*ListSupportedLocalesTool)(nil)
	_ Tool = (*ListSubscriptionsTool)(nil)
	_ Tool = (*CreateSubscriptionTool)(nil)
	_ Tool = (*CancelSubscriptionTool)(nil)
)

// ─── ListGiftCardsTool ──────────────────────────────────────────────────────

type ListGiftCardsTool struct {
	giftcards *giftcardsrv.Service
	tenantID  kernel.TenantID
}

func (t *ListGiftCardsTool) Name() string        { return "list_gift_cards" }
func (t *ListGiftCardsTool) Description() string  { return "List all gift cards for the store" }
func (t *ListGiftCardsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"page":      map[string]any{"type": "integer", "description": "Page number (default 1)"},
			"page_size": map[string]any{"type": "integer", "description": "Items per page (default 20)"},
		},
	}
}

func (t *ListGiftCardsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "list_gift_cards: unmarshal input", errx.TypeValidation)
	}
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 20
	}
	result, err := t.giftcards.List(ctx, t.tenantID, kernel.PaginationOptions{Page: in.Page, PageSize: in.PageSize})
	if err != nil {
		return "", errx.Wrap(err, "list_gift_cards", errx.TypeInternal)
	}
	out := fmt.Sprintf("Found %d gift cards (page %d/%d):\n", result.Total, result.Page, result.TotalPages)
	for _, gc := range result.Items {
		out += fmt.Sprintf("- [%s] Code: %s, Balance: %d %s / %d %s, Active: %v\n",
			gc.ID, gc.Code, gc.Balance.Amount, gc.Balance.Currency,
			gc.InitialAmount.Amount, gc.InitialAmount.Currency, gc.Active)
	}
	return out, nil
}

// ─── CreateGiftCardTool ─────────────────────────────────────────────────────

type CreateGiftCardTool struct {
	giftcards *giftcardsrv.Service
	tenantID  kernel.TenantID
}

func (t *CreateGiftCardTool) Name() string        { return "create_gift_card" }
func (t *CreateGiftCardTool) Description() string  { return "Create a new gift card with an initial balance" }
func (t *CreateGiftCardTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"code":        map[string]any{"type": "string", "description": "Unique gift card code"},
			"amount":      map[string]any{"type": "integer", "description": "Initial amount in cents"},
			"currency":    map[string]any{"type": "string", "description": "Currency code (e.g. USD)"},
			"created_by":  map[string]any{"type": "string", "description": "Who created this card"},
		},
		"required": []string{"code", "amount", "currency"},
	}
}

func (t *CreateGiftCardTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Code      string `json:"code"`
		Amount    int64  `json:"amount"`
		Currency  string `json:"currency"`
		CreatedBy string `json:"created_by"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "create_gift_card: unmarshal input", errx.TypeValidation)
	}
	gc, err := t.giftcards.Create(ctx, t.tenantID, giftcard.CreateInput{
		Code:          in.Code,
		InitialAmount: kernel.Money{Amount: in.Amount, Currency: in.Currency},
		CreatedBy:     in.CreatedBy,
	})
	if err != nil {
		return "", errx.Wrap(err, "create_gift_card", errx.TypeInternal)
	}
	return fmt.Sprintf("Gift card created: ID=%s, Code=%s, Balance=%d %s", gc.ID, gc.Code, gc.Balance.Amount, gc.Balance.Currency), nil
}

// ─── CheckGiftCardBalanceTool ───────────────────────────────────────────────

type CheckGiftCardBalanceTool struct {
	giftcards *giftcardsrv.Service
	tenantID  kernel.TenantID
}

func (t *CheckGiftCardBalanceTool) Name() string        { return "check_gift_card_balance" }
func (t *CheckGiftCardBalanceTool) Description() string  { return "Check the balance of a gift card by its code" }
func (t *CheckGiftCardBalanceTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"code": map[string]any{"type": "string", "description": "Gift card code"},
		},
		"required": []string{"code"},
	}
}

func (t *CheckGiftCardBalanceTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "check_gift_card_balance: unmarshal input", errx.TypeValidation)
	}
	gc, err := t.giftcards.CheckBalance(ctx, t.tenantID, in.Code)
	if err != nil {
		return "", errx.Wrap(err, "check_gift_card_balance", errx.TypeInternal)
	}
	return fmt.Sprintf("Gift card %s: Balance=%d %s, Active=%v", gc.Code, gc.Balance.Amount, gc.Balance.Currency, gc.Active), nil
}

// ─── RedeemGiftCardTool ─────────────────────────────────────────────────────

type RedeemGiftCardTool struct {
	giftcards *giftcardsrv.Service
	tenantID  kernel.TenantID
}

func (t *RedeemGiftCardTool) Name() string        { return "redeem_gift_card" }
func (t *RedeemGiftCardTool) Description() string  { return "Redeem (deduct) an amount from a gift card" }
func (t *RedeemGiftCardTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"code":     map[string]any{"type": "string", "description": "Gift card code"},
			"amount":   map[string]any{"type": "integer", "description": "Amount to redeem in cents"},
			"currency": map[string]any{"type": "string", "description": "Currency code"},
			"note":     map[string]any{"type": "string", "description": "Optional note for the transaction"},
		},
		"required": []string{"code", "amount", "currency"},
	}
}

func (t *RedeemGiftCardTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Code     string `json:"code"`
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
		Note     string `json:"note"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "redeem_gift_card: unmarshal input", errx.TypeValidation)
	}
	gc, err := t.giftcards.Redeem(ctx, t.tenantID, giftcardsrv.RedeemInput{
		Code:   in.Code,
		Amount: kernel.Money{Amount: in.Amount, Currency: in.Currency},
		Note:   in.Note,
	})
	if err != nil {
		return "", errx.Wrap(err, "redeem_gift_card", errx.TypeInternal)
	}
	return fmt.Sprintf("Redeemed %d %s from gift card %s. Remaining balance: %d %s",
		in.Amount, in.Currency, gc.Code, gc.Balance.Amount, gc.Balance.Currency), nil
}

// ─── ListRecoveryEmailsTool ─────────────────────────────────────────────────

type ListRecoveryEmailsTool struct {
	recovery *cartrecoverysrv.Service
	tenantID kernel.TenantID
}

func (t *ListRecoveryEmailsTool) Name() string        { return "list_recovery_emails" }
func (t *ListRecoveryEmailsTool) Description() string  { return "List abandoned cart recovery emails" }
func (t *ListRecoveryEmailsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"page":      map[string]any{"type": "integer", "description": "Page number (default 1)"},
			"page_size": map[string]any{"type": "integer", "description": "Items per page (default 20)"},
		},
	}
}

func (t *ListRecoveryEmailsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "list_recovery_emails: unmarshal input", errx.TypeValidation)
	}
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 20
	}
	result, err := t.recovery.List(ctx, t.tenantID, in.Page, in.PageSize)
	if err != nil {
		return "", errx.Wrap(err, "list_recovery_emails", errx.TypeInternal)
	}
	out := fmt.Sprintf("Found %d recovery emails (page %d/%d):\n", result.Total, result.Page, result.TotalPages)
	for _, r := range result.Items {
		out += fmt.Sprintf("- [%s] Cart: %s, Email: %s, Step: %d, Status: %s\n",
			r.ID, r.CartID, r.Email, r.Step, r.Status)
	}
	return out, nil
}

// ─── GetRecoveryStatsTool ───────────────────────────────────────────────────

type GetRecoveryStatsTool struct {
	recovery *cartrecoverysrv.Service
	tenantID kernel.TenantID
}

func (t *GetRecoveryStatsTool) Name() string        { return "get_recovery_stats" }
func (t *GetRecoveryStatsTool) Description() string  { return "Get abandoned cart recovery statistics" }
func (t *GetRecoveryStatsTool) InputSchema() map[string]any {
	return map[string]any{"type": "object", "properties": map[string]any{}}
}

func (t *GetRecoveryStatsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	stats, err := t.recovery.GetStats(ctx, t.tenantID)
	if err != nil {
		return "", errx.Wrap(err, "get_recovery_stats", errx.TypeInternal)
	}
	return fmt.Sprintf("Recovery Stats: Total=%d, Sent=%d, Clicked=%d, Converted=%d, Rate=%.1f%%",
		stats.Total, stats.Sent, stats.Clicked, stats.Converted, stats.ConversionRate), nil
}

// ─── ListCurrencyRatesTool ──────────────────────────────────────────────────

type ListCurrencyRatesTool struct {
	currency *currencysrv.Service
	tenantID kernel.TenantID
}

func (t *ListCurrencyRatesTool) Name() string        { return "list_currency_rates" }
func (t *ListCurrencyRatesTool) Description() string  { return "List all configured currency exchange rates" }
func (t *ListCurrencyRatesTool) InputSchema() map[string]any {
	return map[string]any{"type": "object", "properties": map[string]any{}}
}

func (t *ListCurrencyRatesTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	rates, err := t.currency.ListRates(ctx, t.tenantID)
	if err != nil {
		return "", errx.Wrap(err, "list_currency_rates", errx.TypeInternal)
	}
	out := fmt.Sprintf("Found %d currency rates:\n", len(rates))
	for _, r := range rates {
		out += fmt.Sprintf("- %s → %s: %.6f\n", r.BaseCurrency, r.TargetCurrency, r.Rate)
	}
	return out, nil
}

// ─── SetCurrencyRateTool ────────────────────────────────────────────────────

type SetCurrencyRateTool struct {
	currency *currencysrv.Service
	tenantID kernel.TenantID
}

func (t *SetCurrencyRateTool) Name() string        { return "set_currency_rate" }
func (t *SetCurrencyRateTool) Description() string  { return "Set or update a currency exchange rate" }
func (t *SetCurrencyRateTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"base_currency":   map[string]any{"type": "string", "description": "Base currency code (e.g. USD)"},
			"target_currency": map[string]any{"type": "string", "description": "Target currency code (e.g. EUR)"},
			"rate":            map[string]any{"type": "number", "description": "Exchange rate"},
			"auto_update":     map[string]any{"type": "boolean", "description": "Enable auto-update"},
		},
		"required": []string{"base_currency", "target_currency", "rate"},
	}
}

func (t *SetCurrencyRateTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		BaseCurrency   string  `json:"base_currency"`
		TargetCurrency string  `json:"target_currency"`
		Rate           float64 `json:"rate"`
		AutoUpdate     bool    `json:"auto_update"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "set_currency_rate: unmarshal input", errx.TypeValidation)
	}
	r, err := t.currency.SetRate(ctx, t.tenantID, currencysrv.SetRateInput{
		BaseCurrency:   in.BaseCurrency,
		TargetCurrency: in.TargetCurrency,
		Rate:           in.Rate,
		AutoUpdate:     in.AutoUpdate,
	})
	if err != nil {
		return "", errx.Wrap(err, "set_currency_rate", errx.TypeInternal)
	}
	return fmt.Sprintf("Rate set: %s → %s = %.6f (ID: %s)", r.BaseCurrency, r.TargetCurrency, r.Rate, r.ID), nil
}

// ─── ConvertCurrencyTool ────────────────────────────────────────────────────

type ConvertCurrencyTool struct {
	currency *currencysrv.Service
	tenantID kernel.TenantID
}

func (t *ConvertCurrencyTool) Name() string        { return "convert_currency" }
func (t *ConvertCurrencyTool) Description() string  { return "Convert an amount from one currency to another" }
func (t *ConvertCurrencyTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"amount":          map[string]any{"type": "integer", "description": "Amount in cents"},
			"from_currency":   map[string]any{"type": "string", "description": "Source currency code"},
			"target_currency": map[string]any{"type": "string", "description": "Target currency code"},
		},
		"required": []string{"amount", "from_currency", "target_currency"},
	}
}

func (t *ConvertCurrencyTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Amount         int64  `json:"amount"`
		FromCurrency   string `json:"from_currency"`
		TargetCurrency string `json:"target_currency"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "convert_currency: unmarshal input", errx.TypeValidation)
	}
	result, err := t.currency.Convert(ctx, t.tenantID, kernel.Money{Amount: in.Amount, Currency: in.FromCurrency}, in.TargetCurrency)
	if err != nil {
		return "", errx.Wrap(err, "convert_currency", errx.TypeInternal)
	}
	return fmt.Sprintf("Converted %d %s → %d %s (rate: %.6f)",
		result.OriginalAmount.Amount, result.OriginalAmount.Currency,
		result.ConvertedAmount.Amount, result.ConvertedAmount.Currency, result.Rate), nil
}

// ─── SetTranslationsTool ────────────────────────────────────────────────────

type SetTranslationsTool struct {
	i18n     *i18nsrv.Service
	tenantID kernel.TenantID
}

func (t *SetTranslationsTool) Name() string        { return "set_translations" }
func (t *SetTranslationsTool) Description() string  { return "Set translations for an entity in a specific locale" }
func (t *SetTranslationsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"entity_type": map[string]any{"type": "string", "description": "Entity type (product, category, collection, page)"},
			"entity_id":   map[string]any{"type": "string", "description": "Entity ID"},
			"locale":      map[string]any{"type": "string", "description": "Locale code (e.g. es, fr, de)"},
			"fields":      map[string]any{"type": "object", "description": "Map of field name to translated value"},
		},
		"required": []string{"entity_type", "entity_id", "locale", "fields"},
	}
}

func (t *SetTranslationsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		EntityType string            `json:"entity_type"`
		EntityID   string            `json:"entity_id"`
		Locale     string            `json:"locale"`
		Fields     map[string]string `json:"fields"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "set_translations: unmarshal input", errx.TypeValidation)
	}
	if err := t.i18n.SetTranslations(ctx, t.tenantID, in.EntityType, in.EntityID, in.Locale, in.Fields); err != nil {
		return "", errx.Wrap(err, "set_translations", errx.TypeInternal)
	}
	return fmt.Sprintf("Set %d translations for %s/%s in locale %s", len(in.Fields), in.EntityType, in.EntityID, in.Locale), nil
}

// ─── GetTranslationsTool ────────────────────────────────────────────────────

type GetTranslationsTool struct {
	i18n     *i18nsrv.Service
	tenantID kernel.TenantID
}

func (t *GetTranslationsTool) Name() string        { return "get_translations" }
func (t *GetTranslationsTool) Description() string  { return "Get translations for an entity in a specific locale" }
func (t *GetTranslationsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"entity_type": map[string]any{"type": "string", "description": "Entity type (product, category, collection, page)"},
			"entity_id":   map[string]any{"type": "string", "description": "Entity ID"},
			"locale":      map[string]any{"type": "string", "description": "Locale code (e.g. es, fr, de)"},
		},
		"required": []string{"entity_type", "entity_id", "locale"},
	}
}

func (t *GetTranslationsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		EntityType string `json:"entity_type"`
		EntityID   string `json:"entity_id"`
		Locale     string `json:"locale"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "get_translations: unmarshal input", errx.TypeValidation)
	}
	bundle, err := t.i18n.GetTranslations(ctx, t.tenantID, in.EntityType, in.EntityID, in.Locale)
	if err != nil {
		return "", errx.Wrap(err, "get_translations", errx.TypeInternal)
	}
	out := fmt.Sprintf("Translations for %s/%s [%s]:\n", bundle.EntityType, bundle.EntityID, bundle.Locale)
	for field, value := range bundle.Fields {
		out += fmt.Sprintf("- %s: %s\n", field, value)
	}
	return out, nil
}

// ─── ListSupportedLocalesTool ───────────────────────────────────────────────

type ListSupportedLocalesTool struct {
	i18n *i18nsrv.Service
}

func (t *ListSupportedLocalesTool) Name() string        { return "list_supported_locales" }
func (t *ListSupportedLocalesTool) Description() string  { return "List all supported locales for translations" }
func (t *ListSupportedLocalesTool) InputSchema() map[string]any {
	return map[string]any{"type": "object", "properties": map[string]any{}}
}

func (t *ListSupportedLocalesTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	locales := t.i18n.ListSupportedLocales()
	out := fmt.Sprintf("Supported locales (%d):\n", len(locales))
	for code, info := range locales {
		out += fmt.Sprintf("- %s: %s (%s)\n", code, info.Name, info.Name)
	}
	return out, nil
}

// ─── ListSubscriptionsTool ──────────────────────────────────────────────────

type ListSubscriptionsTool struct {
	subscriptions *subscriptionsrv.Service
	tenantID      kernel.TenantID
}

func (t *ListSubscriptionsTool) Name() string        { return "list_subscriptions" }
func (t *ListSubscriptionsTool) Description() string  { return "List all subscriptions" }
func (t *ListSubscriptionsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"page":      map[string]any{"type": "integer", "description": "Page number (default 1)"},
			"page_size": map[string]any{"type": "integer", "description": "Items per page (default 20)"},
		},
	}
}

func (t *ListSubscriptionsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "list_subscriptions: unmarshal input", errx.TypeValidation)
	}
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 20
	}
	result, err := t.subscriptions.List(ctx, t.tenantID, in.Page, in.PageSize)
	if err != nil {
		return "", errx.Wrap(err, "list_subscriptions", errx.TypeInternal)
	}
	out := fmt.Sprintf("Found %d subscriptions (page %d/%d):\n", result.Total, result.Page, result.TotalPages)
	for _, s := range result.Items {
		out += fmt.Sprintf("- [%s] Customer: %s, Product: %s, Status: %s, Interval: %s, Price: %d %s\n",
			s.ID, s.CustomerID, s.ProductID, s.Status, s.Interval, s.Price.Amount, s.Price.Currency)
	}
	return out, nil
}

// ─── CreateSubscriptionTool ─────────────────────────────────────────────────

type CreateSubscriptionTool struct {
	subscriptions *subscriptionsrv.Service
	tenantID      kernel.TenantID
}

func (t *CreateSubscriptionTool) Name() string        { return "create_subscription" }
func (t *CreateSubscriptionTool) Description() string  { return "Create a new subscription for a customer" }
func (t *CreateSubscriptionTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"customer_id": map[string]any{"type": "string", "description": "Customer ID"},
			"product_id":  map[string]any{"type": "string", "description": "Product ID"},
			"amount":      map[string]any{"type": "integer", "description": "Price in cents"},
			"currency":    map[string]any{"type": "string", "description": "Currency code"},
			"interval":    map[string]any{"type": "string", "description": "Billing interval (weekly, monthly, quarterly, yearly)"},
		},
		"required": []string{"customer_id", "product_id", "amount", "currency", "interval"},
	}
}

func (t *CreateSubscriptionTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		CustomerID string `json:"customer_id"`
		ProductID  string `json:"product_id"`
		Amount     int64  `json:"amount"`
		Currency   string `json:"currency"`
		Interval   string `json:"interval"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "create_subscription: unmarshal input", errx.TypeValidation)
	}
	sub, err := t.subscriptions.Create(ctx, t.tenantID, subscriptionsrv.CreateInput{
		CustomerID: kernel.CustomerID(in.CustomerID),
		ProductID:  kernel.ProductID(in.ProductID),
		Price:      kernel.Money{Amount: in.Amount, Currency: in.Currency},
		Interval:   subscription.BillingInterval(in.Interval),
	})
	if err != nil {
		return "", errx.Wrap(err, "create_subscription", errx.TypeInternal)
	}
	return fmt.Sprintf("Subscription created: ID=%s, Customer=%s, Product=%s, Interval=%s, NextBilling=%s",
		sub.ID, sub.CustomerID, sub.ProductID, sub.Interval, sub.NextBillingDate.Format("2006-01-02")), nil
}

// ─── CancelSubscriptionTool ─────────────────────────────────────────────────

type CancelSubscriptionTool struct {
	subscriptions *subscriptionsrv.Service
	tenantID      kernel.TenantID
}

func (t *CancelSubscriptionTool) Name() string        { return "cancel_subscription" }
func (t *CancelSubscriptionTool) Description() string  { return "Cancel an active subscription" }
func (t *CancelSubscriptionTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"subscription_id": map[string]any{"type": "string", "description": "Subscription ID to cancel"},
		},
		"required": []string{"subscription_id"},
	}
}

func (t *CancelSubscriptionTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		SubscriptionID string `json:"subscription_id"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "cancel_subscription: unmarshal input", errx.TypeValidation)
	}
	sub, err := t.subscriptions.Cancel(ctx, t.tenantID, kernel.SubscriptionID(in.SubscriptionID))
	if err != nil {
		return "", errx.Wrap(err, "cancel_subscription", errx.TypeInternal)
	}
	return fmt.Sprintf("Subscription %s cancelled. Status: %s", sub.ID, sub.Status), nil
}
