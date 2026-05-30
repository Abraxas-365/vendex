package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Abraxas-365/hada-commerce/internal/customergroup"
	"github.com/Abraxas-365/hada-commerce/internal/customergroup/customergroupsrv"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/payment/paymentsrv"
	"github.com/Abraxas-365/hada-commerce/internal/product/productsrv"
	"github.com/Abraxas-365/hada-commerce/internal/search"
	"github.com/Abraxas-365/hada-commerce/internal/search/searchsrv"
	"github.com/Abraxas-365/hada-commerce/internal/shipping/shippingsrv"
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
)
