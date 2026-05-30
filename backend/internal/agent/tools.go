package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/catalog/catalogsrv"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/order/ordersrv"
	"github.com/Abraxas-365/hada-commerce/internal/product/productsrv"
	"github.com/Abraxas-365/hada-commerce/internal/promo"
	"github.com/Abraxas-365/hada-commerce/internal/promo/promosrv"
	"github.com/Abraxas-365/hada-commerce/internal/storefront"
	"github.com/Abraxas-365/hada-commerce/internal/storefront/storefrontsrv"
	"github.com/Abraxas-365/hada-commerce/internal/theme"
	"github.com/Abraxas-365/hada-commerce/internal/theme/themesrv"
)

// ─── CreatePageTool ──────────────────────────────────────────────────────────

// CreatePageTool lets the agent create a new storefront page.
// Agent-created pages always land in pending_review status.
type CreatePageTool struct {
	sf       *storefrontsrv.Service
	tenantID kernel.TenantID
}

type createPageInput struct {
	Slug        string          `json:"slug"`
	Title       string          `json:"title"`
	HTML        string          `json:"html"`
	CSS         string          `json:"css"`
	Description string          `json:"description"`
	OGTitle     string          `json:"og_title"`
	OGImage     string          `json:"og_image"`
	Keywords    []string        `json:"keywords"`
}

func (t *CreatePageTool) Name() string { return "create_page" }

func (t *CreatePageTool) Description() string {
	return "Create a new storefront page. The page will be placed in pending_review status for admin approval before going live."
}

func (t *CreatePageTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"slug":        map[string]any{"type": "string", "description": "URL-friendly slug, e.g. 'about-us'"},
			"title":       map[string]any{"type": "string", "description": "Page title"},
			"html":        map[string]any{"type": "string", "description": "Full HTML body content"},
			"css":         map[string]any{"type": "string", "description": "Page-scoped CSS styles"},
			"description": map[string]any{"type": "string", "description": "SEO meta description"},
			"og_title":    map[string]any{"type": "string", "description": "Open Graph title (optional)"},
			"og_image":    map[string]any{"type": "string", "description": "Open Graph image URL (optional)"},
			"keywords":    map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "SEO keywords"},
		},
		"required": []string{"slug", "title", "html"},
	}
}

func (t *CreatePageTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in createPageInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "create_page: unmarshal input", errx.TypeValidation)
	}

	page, err := t.sf.CreatePage(ctx, storefrontsrv.CreatePageInput{
		TenantID: t.tenantID,
		Slug:     in.Slug,
		Title:    in.Title,
		HTML:     in.HTML,
		CSS:      in.CSS,
		Meta: storefront.PageMeta{
			Description: in.Description,
			OGTitle:     in.OGTitle,
			OGImage:     in.OGImage,
			Keywords:    in.Keywords,
		},
		CreatedBy: "agent",
		ByAgent:   true,
	})
	if err != nil {
		return "", errx.Wrap(err, "create_page", errx.TypeInternal)
	}

	return fmt.Sprintf(
		"Page created successfully.\nID: %s\nSlug: %s\nTitle: %s\nStatus: %s\n\nThe page is now pending_review. An admin must approve it before it goes live.",
		page.ID, page.Slug, page.Title, page.Status,
	), nil
}

// ─── UpdatePageTool ───────────────────────────────────────────────────────────

// UpdatePageTool lets the agent revise an existing page's content.
type UpdatePageTool struct {
	sf       *storefrontsrv.Service
	tenantID kernel.TenantID
}

type updatePageInput struct {
	PageID  string  `json:"page_id"`
	Title   *string `json:"title,omitempty"`
	HTML    *string `json:"html,omitempty"`
	CSS     *string `json:"css,omitempty"`
	Comment string  `json:"comment"`
}

func (t *UpdatePageTool) Name() string { return "update_page" }

func (t *UpdatePageTool) Description() string {
	return "Update the HTML, CSS, or title of an existing storefront page. Creates a new version snapshot."
}

func (t *UpdatePageTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"page_id": map[string]any{"type": "string", "description": "The ID of the page to update"},
			"title":   map[string]any{"type": "string", "description": "New title (optional)"},
			"html":    map[string]any{"type": "string", "description": "New HTML content (optional)"},
			"css":     map[string]any{"type": "string", "description": "New CSS (optional)"},
			"comment": map[string]any{"type": "string", "description": "Brief description of the changes made"},
		},
		"required": []string{"page_id"},
	}
}

func (t *UpdatePageTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in updatePageInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "update_page: unmarshal input", errx.TypeValidation)
	}

	page, err := t.sf.UpdatePage(ctx, storefrontsrv.UpdatePageInput{
		TenantID: t.tenantID,
		ID:       kernel.PageID(in.PageID),
		Title:    in.Title,
		HTML:     in.HTML,
		CSS:      in.CSS,
		EditedBy: "agent",
		Comment:  in.Comment,
	})
	if err != nil {
		return "", errx.Wrap(err, "update_page", errx.TypeInternal)
	}

	return fmt.Sprintf(
		"Page updated.\nID: %s\nTitle: %s\nVersion: %d\nStatus: %s",
		page.ID, page.Title, page.Version, page.Status,
	), nil
}

// ─── ListPagesTool ────────────────────────────────────────────────────────────

// ListPagesTool lets the agent browse all storefront pages.
type ListPagesTool struct {
	sf       *storefrontsrv.Service
	tenantID kernel.TenantID
}

type listPagesInput struct {
	Status   string `json:"status,omitempty"`
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
}

func (t *ListPagesTool) Name() string { return "list_pages" }

func (t *ListPagesTool) Description() string {
	return "List storefront pages. Optionally filter by status: draft, pending_review, published, archived."
}

func (t *ListPagesTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"status":    map[string]any{"type": "string", "description": "Filter by status: draft | pending_review | published | archived (optional)"},
			"page":      map[string]any{"type": "integer", "description": "Page number (default 1)"},
			"page_size": map[string]any{"type": "integer", "description": "Items per page (default 20, max 100)"},
		},
	}
}

func (t *ListPagesTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in listPagesInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "list_pages: unmarshal input", errx.TypeValidation)
	}

	pg := kernel.NewPaginationOptions(in.Page, in.PageSize)

	var statusFilter *storefront.PageStatus
	if in.Status != "" {
		s := storefront.PageStatus(in.Status)
		statusFilter = &s
	}

	result, err := t.sf.ListPages(ctx, t.tenantID, statusFilter, pg)
	if err != nil {
		return "", errx.Wrap(err, "list_pages", errx.TypeInternal)
	}

	if len(result.Items) == 0 {
		return "No pages found.", nil
	}

	out := fmt.Sprintf("Pages (page %d of %d, total %d):\n\n", result.Page, result.TotalPages, result.Total)
	for _, p := range result.Items {
		out += fmt.Sprintf("- ID: %s | Slug: %s | Title: %s | Status: %s | Version: %d\n",
			p.ID, p.Slug, p.Title, p.Status, p.Version)
	}
	return out, nil
}

// ─── CreateProductTool ────────────────────────────────────────────────────────

// CreateProductTool lets the agent create a new product.
type CreateProductTool struct {
	products *productsrv.Service
	tenantID kernel.TenantID
}

type createProductInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceCents  int64  `json:"price_cents"`
	Currency    string `json:"currency"`
	SKU         string `json:"sku,omitempty"`
	CategoryID  string `json:"category_id,omitempty"`
	Stock       int    `json:"stock"`
	Tags        []string `json:"tags,omitempty"`
}

func (t *CreateProductTool) Name() string { return "create_product" }

func (t *CreateProductTool) Description() string {
	return "Create a new product in the catalog. The product starts as a draft."
}

func (t *CreateProductTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":        map[string]any{"type": "string", "description": "Product name"},
			"description": map[string]any{"type": "string", "description": "Product description"},
			"price_cents": map[string]any{"type": "integer", "description": "Price in the smallest currency unit (e.g. cents for USD)"},
			"currency":    map[string]any{"type": "string", "description": "ISO 4217 currency code, e.g. USD"},
			"sku":         map[string]any{"type": "string", "description": "Stock-keeping unit (optional, must be unique)"},
			"category_id": map[string]any{"type": "string", "description": "Category ID (optional)"},
			"stock":       map[string]any{"type": "integer", "description": "Initial stock quantity"},
			"tags":        map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Product tags (optional)"},
		},
		"required": []string{"name", "price_cents", "currency"},
	}
}

func (t *CreateProductTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in createProductInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "create_product: unmarshal input", errx.TypeValidation)
	}

	p, err := t.products.Create(ctx, t.tenantID, productsrv.CreateInput{
		Name:        in.Name,
		Description: in.Description,
		Price:       kernel.NewMoney(in.PriceCents, in.Currency),
		SKU:         in.SKU,
		CategoryID:  kernel.CategoryID(in.CategoryID),
		Tags:        in.Tags,
		Stock:       in.Stock,
	})
	if err != nil {
		return "", errx.Wrap(err, "create_product", errx.TypeInternal)
	}

	return fmt.Sprintf(
		"Product created.\nID: %s\nName: %s\nSKU: %s\nPrice: %d %s\nStock: %d\nStatus: %s",
		p.ID, p.Name, p.SKU, p.Price.Amount, p.Price.Currency, p.Stock, p.Status,
	), nil
}

// ─── ListProductsTool ─────────────────────────────────────────────────────────

// ListProductsTool lets the agent browse the product catalog.
type ListProductsTool struct {
	products *productsrv.Service
	tenantID kernel.TenantID
}

type listProductsInput struct {
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}

func (t *ListProductsTool) Name() string { return "list_products" }

func (t *ListProductsTool) Description() string {
	return "List products in the catalog with pagination."
}

func (t *ListProductsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"page":      map[string]any{"type": "integer", "description": "Page number (default 1)"},
			"page_size": map[string]any{"type": "integer", "description": "Items per page (default 20, max 100)"},
		},
	}
}

func (t *ListProductsTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in listProductsInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "list_products: unmarshal input", errx.TypeValidation)
	}

	pg := kernel.NewPaginationOptions(in.Page, in.PageSize)
	result, err := t.products.List(ctx, t.tenantID, pg)
	if err != nil {
		return "", errx.Wrap(err, "list_products", errx.TypeInternal)
	}

	if len(result.Items) == 0 {
		return "No products found.", nil
	}

	out := fmt.Sprintf("Products (page %d of %d, total %d):\n\n", result.Page, result.TotalPages, result.Total)
	for _, p := range result.Items {
		out += fmt.Sprintf("- ID: %s | Name: %s | SKU: %s | Price: %d %s | Stock: %d | Status: %s\n",
			p.ID, p.Name, p.SKU, p.Price.Amount, p.Price.Currency, p.Stock, p.Status)
	}
	return out, nil
}

// ─── CreatePromoTool ──────────────────────────────────────────────────────────

// CreatePromoTool lets the agent create discount/promo codes.
type CreatePromoTool struct {
	promos   *promosrv.Service
	tenantID kernel.TenantID
}

type createPromoInput struct {
	Code           string  `json:"code"`
	Type           string  `json:"type"`
	Value          int64   `json:"value"`
	MinOrderCents  *int64  `json:"min_order_cents,omitempty"`
	MaxUses        *int    `json:"max_uses,omitempty"`
	StartsAtUnix   *int64  `json:"starts_at_unix,omitempty"`
	EndsAtUnix     *int64  `json:"ends_at_unix,omitempty"`
}

func (t *CreatePromoTool) Name() string { return "create_promo" }

func (t *CreatePromoTool) Description() string {
	return "Create a promotional discount code. Types: percentage (0-100), fixed_amount (cents), free_shipping."
}

func (t *CreatePromoTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"code":            map[string]any{"type": "string", "description": "Promo code string (e.g. SUMMER20)"},
			"type":            map[string]any{"type": "string", "description": "One of: percentage, fixed_amount, free_shipping"},
			"value":           map[string]any{"type": "integer", "description": "Discount value: percentage (0-100) or fixed cents. Ignored for free_shipping."},
			"min_order_cents": map[string]any{"type": "integer", "description": "Minimum order total in cents to use this promo (optional)"},
			"max_uses":        map[string]any{"type": "integer", "description": "Maximum number of times this code can be used (optional, nil = unlimited)"},
			"starts_at_unix":  map[string]any{"type": "integer", "description": "Unix timestamp when this promo becomes active (optional)"},
			"ends_at_unix":    map[string]any{"type": "integer", "description": "Unix timestamp when this promo expires (optional)"},
		},
		"required": []string{"code", "type", "value"},
	}
}

func (t *CreatePromoTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in createPromoInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "create_promo: unmarshal input", errx.TypeValidation)
	}

	ci := promosrv.CreateInput{
		TenantID:       t.tenantID,
		Code:           in.Code,
		Type:           promo.PromoType(in.Type),
		Value:          in.Value,
		MinOrderAmount: in.MinOrderCents,
		MaxUses:        in.MaxUses,
	}

	if in.StartsAtUnix != nil {
		ts := time.Unix(*in.StartsAtUnix, 0).UTC()
		ci.StartsAt = &ts
	}
	if in.EndsAtUnix != nil {
		ts := time.Unix(*in.EndsAtUnix, 0).UTC()
		ci.EndsAt = &ts
	}

	p, err := t.promos.Create(ctx, ci)
	if err != nil {
		return "", errx.Wrap(err, "create_promo", errx.TypeInternal)
	}

	return fmt.Sprintf(
		"Promo code created.\nID: %s\nCode: %s\nType: %s\nValue: %d\nActive: %v",
		p.ID, p.Code, p.Type, p.Value, p.Active,
	), nil
}

// ─── QueryOrdersTool ──────────────────────────────────────────────────────────

// QueryOrdersTool lets the agent inspect recent orders.
type QueryOrdersTool struct {
	orders   *ordersrv.Service
	tenantID kernel.TenantID
}

type queryOrdersInput struct {
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}

func (t *QueryOrdersTool) Name() string { return "query_orders" }

func (t *QueryOrdersTool) Description() string {
	return "List recent orders for the store. Returns order IDs, statuses, totals, and item counts."
}

func (t *QueryOrdersTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"page":      map[string]any{"type": "integer", "description": "Page number (default 1)"},
			"page_size": map[string]any{"type": "integer", "description": "Items per page (default 20, max 100)"},
		},
	}
}

func (t *QueryOrdersTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in queryOrdersInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "query_orders: unmarshal input", errx.TypeValidation)
	}

	pg := kernel.NewPaginationOptions(in.Page, in.PageSize)
	result, err := t.orders.List(ctx, t.tenantID, pg)
	if err != nil {
		return "", errx.Wrap(err, "query_orders", errx.TypeInternal)
	}

	if len(result.Items) == 0 {
		return "No orders found.", nil
	}

	out := fmt.Sprintf("Orders (page %d of %d, total %d):\n\n", result.Page, result.TotalPages, result.Total)
	for _, o := range result.Items {
		out += fmt.Sprintf("- ID: %s | Status: %s | Total: %d %s | Items: %d | Created: %s\n",
			o.ID, o.Status, o.TotalAmount.Amount, o.TotalAmount.Currency, o.ItemCount(), o.CreatedAt.Format("2006-01-02 15:04"))
	}
	return out, nil
}

// ─── SearchCatalogTool ────────────────────────────────────────────────────────

// SearchCatalogTool lets the agent browse categories and collections.
type SearchCatalogTool struct {
	catalog  *catalogsrv.Service
	tenantID kernel.TenantID
}

type searchCatalogInput struct {
	Target   string `json:"target"` // "categories" | "collections"
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
}

func (t *SearchCatalogTool) Name() string { return "search_catalog" }

func (t *SearchCatalogTool) Description() string {
	return "Browse the catalog structure: list categories or collections."
}

func (t *SearchCatalogTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"target":    map[string]any{"type": "string", "description": "What to search: 'categories' or 'collections'"},
			"page":      map[string]any{"type": "integer", "description": "Page number (default 1)"},
			"page_size": map[string]any{"type": "integer", "description": "Items per page (default 20, max 100)"},
		},
		"required": []string{"target"},
	}
}

func (t *SearchCatalogTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in searchCatalogInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "search_catalog: unmarshal input", errx.TypeValidation)
	}

	pg := kernel.NewPaginationOptions(in.Page, in.PageSize)

	switch in.Target {
	case "categories":
		return t.listCategories(ctx, pg)
	case "collections":
		return t.listCollections(ctx, pg)
	default:
		return "", errx.Validation(fmt.Sprintf("search_catalog: unknown target %q; use 'categories' or 'collections'", in.Target))
	}
}

func (t *SearchCatalogTool) listCategories(ctx context.Context, pg kernel.PaginationOptions) (string, error) {
	result, err := t.catalog.ListCategories(ctx, t.tenantID, pg)
	if err != nil {
		return "", errx.Wrap(err, "list categories", errx.TypeInternal)
	}
	if len(result.Items) == 0 {
		return "No categories found.", nil
	}
	out := fmt.Sprintf("Categories (page %d of %d, total %d):\n\n", result.Page, result.TotalPages, result.Total)
	for _, c := range result.Items {
		parent := "root"
		if c.ParentID != nil {
			parent = string(*c.ParentID)
		}
		out += fmt.Sprintf("- ID: %s | Name: %s | Slug: %s | Parent: %s\n", c.ID, c.Name, c.Slug, parent)
	}
	return out, nil
}

func (t *SearchCatalogTool) listCollections(ctx context.Context, pg kernel.PaginationOptions) (string, error) {
	result, err := t.catalog.ListCollections(ctx, t.tenantID, pg)
	if err != nil {
		return "", errx.Wrap(err, "list collections", errx.TypeInternal)
	}
	if len(result.Items) == 0 {
		return "No collections found.", nil
	}
	out := fmt.Sprintf("Collections (page %d of %d, total %d):\n\n", result.Page, result.TotalPages, result.Total)
	for _, c := range result.Items {
		out += fmt.Sprintf("- ID: %s | Name: %s | Slug: %s | Products: %d | Automatic: %v\n",
			c.ID, c.Name, c.Slug, len(c.ProductIDs), c.IsAutomatic)
	}
	return out, nil
}

// ─── ListBlockTypesTool ───────────────────────────────────────────────────────

// ListBlockTypesTool lets the agent browse registered block types, optionally
// filtered by category.
type ListBlockTypesTool struct {
	sf *storefrontsrv.Service
}

type listBlockTypesInput struct {
	Category string `json:"category,omitempty"`
}

func (t *ListBlockTypesTool) Name() string { return "list_block_types" }

func (t *ListBlockTypesTool) Description() string {
	return "List registered storefront block types. Optionally filter by category: content, commerce, media, layout."
}

func (t *ListBlockTypesTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"category": map[string]any{
				"type":        "string",
				"description": "Optional category filter: content | commerce | media | layout",
			},
		},
	}
}

func (t *ListBlockTypesTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in listBlockTypesInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "list_block_types: unmarshal input", errx.TypeValidation)
	}

	bts, err := t.sf.ListBlockTypes(ctx, in.Category)
	if err != nil {
		return "", errx.Wrap(err, "list_block_types", errx.TypeInternal)
	}

	if len(bts) == 0 {
		return "No block types found.", nil
	}

	out := fmt.Sprintf("Block types (%d):\n\n", len(bts))
	for _, bt := range bts {
		out += fmt.Sprintf("- ID: %s | Name: %s | Display: %s | Category: %s | Icon: %s\n",
			bt.ID, bt.Name, bt.DisplayName, bt.Category, bt.Icon)
	}
	return out, nil
}

var _ Tool = (*ListBlockTypesTool)(nil)

// ─── CreateBlockTypeTool ──────────────────────────────────────────────────────

// CreateBlockTypeTool lets the agent register a new block type in the storefront.
type CreateBlockTypeTool struct {
	sf *storefrontsrv.Service
}

type createBlockTypeInput struct {
	Name            string          `json:"name"`
	DisplayName     string          `json:"display_name"`
	Category        string          `json:"category"`
	Schema          json.RawMessage `json:"schema,omitempty"`
	DefaultSettings json.RawMessage `json:"default_settings,omitempty"`
	Icon            string          `json:"icon,omitempty"`
}

func (t *CreateBlockTypeTool) Name() string { return "create_block_type" }

func (t *CreateBlockTypeTool) Description() string {
	return "Register a new storefront block type with its JSON schema and default settings."
}

func (t *CreateBlockTypeTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":             map[string]any{"type": "string", "description": "Unique machine-readable name, e.g. 'hero_banner'"},
			"display_name":     map[string]any{"type": "string", "description": "Human-readable label shown in the editor"},
			"category":         map[string]any{"type": "string", "description": "One of: content, commerce, media, layout"},
			"schema":           map[string]any{"type": "object", "description": "JSON Schema describing the block's settings (optional)"},
			"default_settings": map[string]any{"type": "object", "description": "Default values matching the schema (optional)"},
			"icon":             map[string]any{"type": "string", "description": "Icon identifier or SVG string (optional)"},
		},
		"required": []string{"name", "display_name", "category"},
	}
}

func (t *CreateBlockTypeTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in createBlockTypeInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "create_block_type: unmarshal input", errx.TypeValidation)
	}

	// Convert optional JSON objects to []byte (nil if not provided).
	var schemaBytes, defaultSettingsBytes []byte
	if len(in.Schema) > 0 && string(in.Schema) != "null" {
		schemaBytes = []byte(in.Schema)
	}
	if len(in.DefaultSettings) > 0 && string(in.DefaultSettings) != "null" {
		defaultSettingsBytes = []byte(in.DefaultSettings)
	}

	bt, err := t.sf.CreateBlockType(ctx, storefrontsrv.CreateBlockTypeInput{
		Name:            in.Name,
		DisplayName:     in.DisplayName,
		Category:        in.Category,
		Schema:          schemaBytes,
		DefaultSettings: defaultSettingsBytes,
		Icon:            in.Icon,
	})
	if err != nil {
		return "", errx.Wrap(err, "create_block_type", errx.TypeInternal)
	}

	return fmt.Sprintf(
		"Block type created.\nID: %s\nName: %s\nDisplay: %s\nCategory: %s",
		bt.ID, bt.Name, bt.DisplayName, bt.Category,
	), nil
}

var _ Tool = (*CreateBlockTypeTool)(nil)

// ─── ListThemesTool ───────────────────────────────────────────────────────────

// ListThemesTool lets the agent browse all themes for a tenant.
type ListThemesTool struct {
	themes   *themesrv.Service
	tenantID kernel.TenantID
}

func (t *ListThemesTool) Name() string { return "list_themes" }

func (t *ListThemesTool) Description() string {
	return "List all themes for the tenant, showing which one is currently active."
}

func (t *ListThemesTool) InputSchema() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

func (t *ListThemesTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	themes, err := t.themes.ListThemes(ctx, t.tenantID)
	if err != nil {
		return "", errx.Wrap(err, "list_themes", errx.TypeInternal)
	}

	if len(themes) == 0 {
		return "No themes found.", nil
	}

	out := fmt.Sprintf("Themes (%d):\n\n", len(themes))
	for _, th := range themes {
		active := ""
		if th.IsActive {
			active = " [ACTIVE]"
		}
		out += fmt.Sprintf("- ID: %s | Name: %s%s | Updated: %s\n",
			th.ID, th.Name, active, th.UpdatedAt.Format("2006-01-02"))
	}
	return out, nil
}

var _ Tool = (*ListThemesTool)(nil)

// ─── GetActiveThemeTool ───────────────────────────────────────────────────────

// GetActiveThemeTool retrieves the currently active theme and its design tokens.
type GetActiveThemeTool struct {
	themes   *themesrv.Service
	tenantID kernel.TenantID
}

func (t *GetActiveThemeTool) Name() string { return "get_active_theme" }

func (t *GetActiveThemeTool) Description() string {
	return "Get the currently active theme and all its design tokens (colors, typography, spacing, borders, shadows)."
}

func (t *GetActiveThemeTool) InputSchema() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

func (t *GetActiveThemeTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	th, err := t.themes.GetActiveTheme(ctx, t.tenantID)
	if err != nil {
		return "", errx.Wrap(err, "get_active_theme", errx.TypeInternal)
	}

	return formatTheme(th), nil
}

var _ Tool = (*GetActiveThemeTool)(nil)

// ─── CreateThemeTool ──────────────────────────────────────────────────────────

// CreateThemeTool lets the agent create a new theme with design tokens.
type CreateThemeTool struct {
	themes   *themesrv.Service
	tenantID kernel.TenantID
}

type createThemeInput struct {
	Name   string            `json:"name"`
	Tokens *theme.ThemeTokens `json:"tokens,omitempty"`
}

func (t *CreateThemeTool) Name() string { return "create_theme" }

func (t *CreateThemeTool) Description() string {
	return "Create a new theme with optional design tokens. Unset tokens use sensible defaults."
}

func (t *CreateThemeTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{"type": "string", "description": "Theme name, e.g. 'Dark Mode' or 'Holiday Special'"},
			"tokens": map[string]any{
				"type":        "object",
				"description": "Optional design tokens to override defaults (colors, typography, spacing, borders, shadows)",
			},
		},
		"required": []string{"name"},
	}
}

func (t *CreateThemeTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in createThemeInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "create_theme: unmarshal input", errx.TypeValidation)
	}

	th, err := t.themes.CreateTheme(ctx, themesrv.CreateThemeInput{
		TenantID: t.tenantID,
		Name:     in.Name,
		Tokens:   in.Tokens,
	})
	if err != nil {
		return "", errx.Wrap(err, "create_theme", errx.TypeInternal)
	}

	return fmt.Sprintf(
		"Theme created.\nID: %s\nName: %s\nActive: %v\nPrimary color: %s\nFont heading: %s",
		th.ID, th.Name, th.IsActive, th.Tokens.Colors.Primary, th.Tokens.Typography.FontHeading,
	), nil
}

var _ Tool = (*CreateThemeTool)(nil)

// ─── UpdateThemeTool ──────────────────────────────────────────────────────────

// UpdateThemeTool lets the agent update a theme's name or design tokens.
type UpdateThemeTool struct {
	themes   *themesrv.Service
	tenantID kernel.TenantID
}

type updateThemeInput struct {
	ID     string             `json:"id"`
	Name   *string            `json:"name,omitempty"`
	Tokens *theme.ThemeTokens `json:"tokens,omitempty"`
}

func (t *UpdateThemeTool) Name() string { return "update_theme" }

func (t *UpdateThemeTool) Description() string {
	return "Update a theme's name or design tokens. Only the provided fields are changed."
}

func (t *UpdateThemeTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"id":   map[string]any{"type": "string", "description": "Theme ID to update"},
			"name": map[string]any{"type": "string", "description": "New name for the theme (optional)"},
			"tokens": map[string]any{
				"type":        "object",
				"description": "Design tokens to apply; merges with existing values (optional)",
			},
		},
		"required": []string{"id"},
	}
}

func (t *UpdateThemeTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in updateThemeInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "update_theme: unmarshal input", errx.TypeValidation)
	}

	th, err := t.themes.UpdateTheme(ctx, themesrv.UpdateThemeInput{
		TenantID: t.tenantID,
		ID:       kernel.ThemeID(in.ID),
		Name:     in.Name,
		Tokens:   in.Tokens,
	})
	if err != nil {
		return "", errx.Wrap(err, "update_theme", errx.TypeInternal)
	}

	return fmt.Sprintf(
		"Theme updated.\nID: %s\nName: %s\nActive: %v\nPrimary color: %s",
		th.ID, th.Name, th.IsActive, th.Tokens.Colors.Primary,
	), nil
}

var _ Tool = (*UpdateThemeTool)(nil)

// ─── ActivateThemeTool ────────────────────────────────────────────────────────

// ActivateThemeTool switches the active theme for the tenant.
type ActivateThemeTool struct {
	themes   *themesrv.Service
	tenantID kernel.TenantID
}

type activateThemeInput struct {
	ID string `json:"id"`
}

func (t *ActivateThemeTool) Name() string { return "activate_theme" }

func (t *ActivateThemeTool) Description() string {
	return "Activate a theme for the tenant. Deactivates all other themes. The active theme is applied to the storefront."
}

func (t *ActivateThemeTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"id": map[string]any{"type": "string", "description": "ID of the theme to activate"},
		},
		"required": []string{"id"},
	}
}

func (t *ActivateThemeTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
	var in activateThemeInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return "", errx.Wrap(err, "activate_theme: unmarshal input", errx.TypeValidation)
	}

	th, err := t.themes.ActivateTheme(ctx, t.tenantID, kernel.ThemeID(in.ID))
	if err != nil {
		return "", errx.Wrap(err, "activate_theme", errx.TypeInternal)
	}

	return fmt.Sprintf(
		"Theme activated.\nID: %s\nName: %s\nActive: %v",
		th.ID, th.Name, th.IsActive,
	), nil
}

var _ Tool = (*ActivateThemeTool)(nil)

// ─── formatTheme ─────────────────────────────────────────────────────────────

// formatTheme returns a human-readable summary of a theme and its design tokens.
func formatTheme(th *theme.Theme) string {
	return fmt.Sprintf(
		"Theme: %s\nID: %s\nActive: %v\nUpdated: %s\n\nColors:\n  Primary: %s\n  Secondary: %s\n  Background: %s\n  Surface: %s\n  Text: %s\n  Text Muted: %s\n  Border: %s\n  Error: %s\n  Success: %s\n  Warning: %s\n  Info: %s\n\nTypography:\n  Font Heading: %s\n  Font Body: %s\n  Base Size: %s\n  Scale Ratio: %.2f\n\nSpacing:\n  Unit: %s\n  Section Padding: %s\n\nBorders:\n  Sm: %s\n  Md: %s\n  Lg: %s\n  Full: %s\n\nShadows:\n  Sm: %s\n  Md: %s\n  Lg: %s",
		th.Name, th.ID, th.IsActive, th.UpdatedAt.Format("2006-01-02 15:04"),
		th.Tokens.Colors.Primary, th.Tokens.Colors.Secondary, th.Tokens.Colors.Background,
		th.Tokens.Colors.Surface, th.Tokens.Colors.Text, th.Tokens.Colors.TextMuted,
		th.Tokens.Colors.Border, th.Tokens.Colors.Error, th.Tokens.Colors.Success,
		th.Tokens.Colors.Warning, th.Tokens.Colors.Info,
		th.Tokens.Typography.FontHeading, th.Tokens.Typography.FontBody,
		th.Tokens.Typography.BaseSize, th.Tokens.Typography.ScaleRatio,
		th.Tokens.Spacing.Unit, th.Tokens.Spacing.SectionPadding,
		th.Tokens.Borders.RadiusSm, th.Tokens.Borders.RadiusMd,
		th.Tokens.Borders.RadiusLg, th.Tokens.Borders.RadiusFull,
		th.Tokens.Shadows.Sm, th.Tokens.Shadows.Md, th.Tokens.Shadows.Lg,
	)
}

// ─── compile-time guards ──────────────────────────────────────────────────────

var (
	_ Tool = (*CreatePageTool)(nil)
	_ Tool = (*UpdatePageTool)(nil)
	_ Tool = (*ListPagesTool)(nil)
	_ Tool = (*CreateProductTool)(nil)
	_ Tool = (*ListProductsTool)(nil)
	_ Tool = (*CreatePromoTool)(nil)
	_ Tool = (*QueryOrdersTool)(nil)
	_ Tool = (*SearchCatalogTool)(nil)
	_ Tool = (*ListBlockTypesTool)(nil)
	_ Tool = (*CreateBlockTypeTool)(nil)
	_ Tool = (*ListThemesTool)(nil)
	_ Tool = (*GetActiveThemeTool)(nil)
	_ Tool = (*CreateThemeTool)(nil)
	_ Tool = (*UpdateThemeTool)(nil)
	_ Tool = (*ActivateThemeTool)(nil)
)
