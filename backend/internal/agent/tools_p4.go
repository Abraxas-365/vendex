package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Abraxas-365/vendex/internal/audit"
	"github.com/Abraxas-365/vendex/internal/audit/auditsrv"
	"github.com/Abraxas-365/vendex/internal/bundle"
	"github.com/Abraxas-365/vendex/internal/bundle/bundlesrv"
	"github.com/Abraxas-365/vendex/internal/dashboard"
	"github.com/Abraxas-365/vendex/internal/dashboard/dashboardsrv"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/inventory"
	"github.com/Abraxas-365/vendex/internal/inventory/inventorysrv"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/loyalty"
	"github.com/Abraxas-365/vendex/internal/loyalty/loyaltysrv"
	"github.com/Abraxas-365/vendex/internal/notification/notificationsrv"
	"github.com/Abraxas-365/vendex/internal/returns"
	"github.com/Abraxas-365/vendex/internal/returns/returnssrv"
	"github.com/Abraxas-365/vendex/internal/review/reviewsrv"
	"github.com/Abraxas-365/vendex/internal/webhook"
	"github.com/Abraxas-365/vendex/internal/webhook/webhooksrv"
)

// ─── Inventory Tools ────────────────────────────────────────────────────────

type ListWarehousesTool struct {
	inventory *inventorysrv.Service
	tenantID  kernel.TenantID
}

func (t *ListWarehousesTool) Name() string        { return "list_warehouses" }
func (t *ListWarehousesTool) Description() string  { return "List all warehouses for the tenant" }
func (t *ListWarehousesTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"page":      map[string]any{"type": "integer", "description": "Page number (default 1)"},
			"page_size": map[string]any{"type": "integer", "description": "Items per page (default 20)"},
		},
	}
}
func (t *ListWarehousesTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "list_warehouses: unmarshal input", errx.TypeValidation)
	}
	if in.Page == 0 { in.Page = 1 }
	if in.PageSize == 0 { in.PageSize = 20 }
	result, err := t.inventory.ListWarehouses(ctx, t.tenantID, kernel.PaginationOptions{Page: in.Page, PageSize: in.PageSize})
	if err != nil {
		return "", errx.Wrap(err, "list_warehouses", errx.TypeInternal)
	}
	var lines []string
	for _, w := range result.Items {
		lines = append(lines, fmt.Sprintf("- %s: %s (address: %s, default: %v)", w.ID, w.Name, w.Address, w.IsDefault))
	}
	return fmt.Sprintf("Found %d warehouses (page %d/%d):\n%s", result.Total, result.Page, result.TotalPages, strings.Join(lines, "\n")), nil
}

type CreateWarehouseTool struct {
	inventory *inventorysrv.Service
	tenantID  kernel.TenantID
}

func (t *CreateWarehouseTool) Name() string        { return "create_warehouse" }
func (t *CreateWarehouseTool) Description() string  { return "Create a new warehouse" }
func (t *CreateWarehouseTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":       map[string]any{"type": "string", "description": "Warehouse name"},
			"address":    map[string]any{"type": "string", "description": "Warehouse address"},
			"is_default": map[string]any{"type": "boolean", "description": "Whether this is the default warehouse"},
		},
		"required": []string{"name"},
	}
}
func (t *CreateWarehouseTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in inventory.CreateWarehouseInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "create_warehouse: unmarshal input", errx.TypeValidation)
	}
	w, err := t.inventory.CreateWarehouse(ctx, t.tenantID, in)
	if err != nil {
		return "", errx.Wrap(err, "create_warehouse", errx.TypeInternal)
	}
	return fmt.Sprintf("Created warehouse %s: %s", w.ID, w.Name), nil
}

type AdjustStockTool struct {
	inventory *inventorysrv.Service
	tenantID  kernel.TenantID
}

func (t *AdjustStockTool) Name() string        { return "adjust_stock" }
func (t *AdjustStockTool) Description() string  { return "Adjust stock level for a product in a warehouse (receive, sell, adjust, transfer, return)" }
func (t *AdjustStockTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"product_id":   map[string]any{"type": "string"},
			"warehouse_id": map[string]any{"type": "string"},
			"quantity":     map[string]any{"type": "integer", "description": "Positive to add, negative to subtract"},
			"type":         map[string]any{"type": "string", "enum": []string{"receive", "sell", "adjust", "transfer", "return"}},
			"reference":    map[string]any{"type": "string"},
			"note":         map[string]any{"type": "string"},
		},
		"required": []string{"product_id", "warehouse_id", "quantity", "type"},
	}
}
func (t *AdjustStockTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in inventory.AdjustStockInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "adjust_stock: unmarshal input", errx.TypeValidation)
	}
	stock, err := t.inventory.AdjustStock(ctx, t.tenantID, in)
	if err != nil {
		return "", errx.Wrap(err, "adjust_stock", errx.TypeInternal)
	}
	return fmt.Sprintf("Stock adjusted. Product %s at warehouse %s: available=%d, reserved=%d", stock.ProductID, stock.WarehouseID, stock.Available(), stock.Reserved), nil
}

type GetLowStockTool struct {
	inventory *inventorysrv.Service
	tenantID  kernel.TenantID
}

func (t *GetLowStockTool) Name() string        { return "get_low_stock" }
func (t *GetLowStockTool) Description() string  { return "Get all products with low stock levels" }
func (t *GetLowStockTool) InputSchema() map[string]any {
	return map[string]any{"type": "object", "properties": map[string]any{}}
}
func (t *GetLowStockTool) Execute(ctx context.Context, _ json.RawMessage) (string, error) {
	items, err := t.inventory.GetLowStockItems(ctx, t.tenantID)
	if err != nil {
		return "", errx.Wrap(err, "get_low_stock", errx.TypeInternal)
	}
	if len(items) == 0 {
		return "No low stock items found.", nil
	}
	var lines []string
	for _, s := range items {
		lines = append(lines, fmt.Sprintf("- Product %s at warehouse %s: available=%d, threshold=%d", s.ProductID, s.WarehouseID, s.Available(), s.LowStockThreshold))
	}
	return fmt.Sprintf("%d low stock items:\n%s", len(items), strings.Join(lines, "\n")), nil
}

// ─── Reviews Tools ──────────────────────────────────────────────────────────

type ListReviewsTool struct {
	reviews  *reviewsrv.Service
	tenantID kernel.TenantID
}

func (t *ListReviewsTool) Name() string        { return "list_reviews" }
func (t *ListReviewsTool) Description() string  { return "List product reviews, optionally filtered by status (pending, approved, rejected)" }
func (t *ListReviewsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"status":    map[string]any{"type": "string", "enum": []string{"pending", "approved", "rejected"}},
			"page":      map[string]any{"type": "integer"},
			"page_size": map[string]any{"type": "integer"},
		},
	}
}
func (t *ListReviewsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Status   string `json:"status"`
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "list_reviews: unmarshal input", errx.TypeValidation)
	}
	if in.Page == 0 { in.Page = 1 }
	if in.PageSize == 0 { in.PageSize = 20 }
	result, err := t.reviews.List(ctx, t.tenantID, in.Status, kernel.PaginationOptions{Page: in.Page, PageSize: in.PageSize})
	if err != nil {
		return "", errx.Wrap(err, "list_reviews", errx.TypeInternal)
	}
	var lines []string
	for _, r := range result.Items {
		lines = append(lines, fmt.Sprintf("- [%s] %s: %d stars by customer %s — %s", r.Status, r.ID, r.Rating, r.CustomerID, r.Title))
	}
	return fmt.Sprintf("Found %d reviews (page %d/%d):\n%s", result.Total, result.Page, result.TotalPages, strings.Join(lines, "\n")), nil
}

type ApproveReviewTool struct {
	reviews  *reviewsrv.Service
	tenantID kernel.TenantID
}

func (t *ApproveReviewTool) Name() string        { return "approve_review" }
func (t *ApproveReviewTool) Description() string  { return "Approve a pending product review" }
func (t *ApproveReviewTool) InputSchema() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{"review_id": map[string]any{"type": "string"}},
		"required":   []string{"review_id"},
	}
}
func (t *ApproveReviewTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct{ ReviewID string `json:"review_id"` }
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "approve_review: unmarshal input", errx.TypeValidation)
	}
	r, err := t.reviews.Approve(ctx, t.tenantID, kernel.ReviewID(in.ReviewID))
	if err != nil {
		return "", errx.Wrap(err, "approve_review", errx.TypeInternal)
	}
	return fmt.Sprintf("Review %s approved (product %s, %d stars)", r.ID, r.ProductID, r.Rating), nil
}

type RejectReviewTool struct {
	reviews  *reviewsrv.Service
	tenantID kernel.TenantID
}

func (t *RejectReviewTool) Name() string        { return "reject_review" }
func (t *RejectReviewTool) Description() string  { return "Reject a product review" }
func (t *RejectReviewTool) InputSchema() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{"review_id": map[string]any{"type": "string"}},
		"required":   []string{"review_id"},
	}
}
func (t *RejectReviewTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct{ ReviewID string `json:"review_id"` }
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "reject_review: unmarshal input", errx.TypeValidation)
	}
	r, err := t.reviews.Reject(ctx, t.tenantID, kernel.ReviewID(in.ReviewID))
	if err != nil {
		return "", errx.Wrap(err, "reject_review", errx.TypeInternal)
	}
	return fmt.Sprintf("Review %s rejected", r.ID), nil
}

// ─── Returns Tools ──────────────────────────────────────────────────────────

type ListReturnsTool struct {
	returns  *returnssrv.Service
	tenantID kernel.TenantID
}

func (t *ListReturnsTool) Name() string        { return "list_returns" }
func (t *ListReturnsTool) Description() string  { return "List return requests, optionally filtered by status" }
func (t *ListReturnsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"status":    map[string]any{"type": "string", "enum": []string{"pending", "approved", "rejected", "received", "refunded", "closed"}},
			"page":      map[string]any{"type": "integer"},
			"page_size": map[string]any{"type": "integer"},
		},
	}
}
func (t *ListReturnsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Status   string `json:"status"`
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "list_returns: unmarshal input", errx.TypeValidation)
	}
	if in.Page == 0 { in.Page = 1 }
	if in.PageSize == 0 { in.PageSize = 20 }
	result, err := t.returns.List(ctx, t.tenantID, in.Status, in.Page, in.PageSize)
	if err != nil {
		return "", errx.Wrap(err, "list_returns", errx.TypeInternal)
	}
	var lines []string
	for _, r := range result.Items {
		lines = append(lines, fmt.Sprintf("- %s: order %s, status=%s, reason=%s", r.ID, r.OrderID, r.Status, r.Reason))
	}
	return fmt.Sprintf("Found %d returns (page %d/%d):\n%s", result.Total, result.Page, result.TotalPages, strings.Join(lines, "\n")), nil
}

type ApproveReturnTool struct {
	returns  *returnssrv.Service
	tenantID kernel.TenantID
}

func (t *ApproveReturnTool) Name() string        { return "approve_return" }
func (t *ApproveReturnTool) Description() string  { return "Approve a pending return request" }
func (t *ApproveReturnTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"return_id":   map[string]any{"type": "string"},
			"admin_notes": map[string]any{"type": "string"},
		},
		"required": []string{"return_id"},
	}
}
func (t *ApproveReturnTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		ReturnID   string `json:"return_id"`
		AdminNotes string `json:"admin_notes"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "approve_return: unmarshal input", errx.TypeValidation)
	}
	r, err := t.returns.Approve(ctx, t.tenantID, kernel.ReturnID(in.ReturnID), returns.ApproveInput{AdminNotes: in.AdminNotes})
	if err != nil {
		return "", errx.Wrap(err, "approve_return", errx.TypeInternal)
	}
	return fmt.Sprintf("Return %s approved", r.ID), nil
}

// ─── Webhooks Tools ─────────────────────────────────────────────────────────

type ListWebhooksTool struct {
	webhooks *webhooksrv.Service
	tenantID kernel.TenantID
}

func (t *ListWebhooksTool) Name() string        { return "list_webhooks" }
func (t *ListWebhooksTool) Description() string  { return "List registered webhooks" }
func (t *ListWebhooksTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"page":      map[string]any{"type": "integer"},
			"page_size": map[string]any{"type": "integer"},
		},
	}
}
func (t *ListWebhooksTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "list_webhooks: unmarshal input", errx.TypeValidation)
	}
	if in.Page == 0 { in.Page = 1 }
	if in.PageSize == 0 { in.PageSize = 20 }
	result, err := t.webhooks.List(ctx, t.tenantID, in.Page, in.PageSize)
	if err != nil {
		return "", errx.Wrap(err, "list_webhooks", errx.TypeInternal)
	}
	var lines []string
	for _, w := range result.Items {
		lines = append(lines, fmt.Sprintf("- %s: %s (active=%v, events=%v)", w.ID, w.URL, w.Active, w.Events))
	}
	return fmt.Sprintf("Found %d webhooks (page %d/%d):\n%s", result.Total, result.Page, result.TotalPages, strings.Join(lines, "\n")), nil
}

type CreateWebhookTool struct {
	webhooks *webhooksrv.Service
	tenantID kernel.TenantID
}

func (t *CreateWebhookTool) Name() string        { return "create_webhook" }
func (t *CreateWebhookTool) Description() string  { return "Register a new webhook to receive event notifications" }
func (t *CreateWebhookTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"url":         map[string]any{"type": "string", "description": "The URL to POST events to"},
			"events":      map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Event types to subscribe to"},
			"secret":      map[string]any{"type": "string", "description": "HMAC secret for payload signing"},
			"description": map[string]any{"type": "string"},
		},
		"required": []string{"url", "events"},
	}
}
func (t *CreateWebhookTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in webhook.CreateWebhookInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "create_webhook: unmarshal input", errx.TypeValidation)
	}
	w, err := t.webhooks.Create(ctx, t.tenantID, in)
	if err != nil {
		return "", errx.Wrap(err, "create_webhook", errx.TypeInternal)
	}
	return fmt.Sprintf("Created webhook %s: %s", w.ID, w.URL), nil
}

type ToggleWebhookTool struct {
	webhooks *webhooksrv.Service
	tenantID kernel.TenantID
}

func (t *ToggleWebhookTool) Name() string        { return "toggle_webhook" }
func (t *ToggleWebhookTool) Description() string  { return "Activate or deactivate a webhook" }
func (t *ToggleWebhookTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"webhook_id": map[string]any{"type": "string"},
			"active":     map[string]any{"type": "boolean"},
		},
		"required": []string{"webhook_id", "active"},
	}
}
func (t *ToggleWebhookTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		WebhookID string `json:"webhook_id"`
		Active    bool   `json:"active"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "toggle_webhook: unmarshal input", errx.TypeValidation)
	}
	w, err := t.webhooks.Toggle(ctx, t.tenantID, kernel.WebhookID(in.WebhookID), in.Active)
	if err != nil {
		return "", errx.Wrap(err, "toggle_webhook", errx.TypeInternal)
	}
	return fmt.Sprintf("Webhook %s is now active=%v", w.ID, w.Active), nil
}

// ─── Audit Tools ────────────────────────────────────────────────────────────

type ListAuditLogsTool struct {
	audit    *auditsrv.Service
	tenantID kernel.TenantID
}

func (t *ListAuditLogsTool) Name() string        { return "list_audit_logs" }
func (t *ListAuditLogsTool) Description() string  { return "Search audit logs by user, action, resource type, or date range" }
func (t *ListAuditLogsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"user_id":       map[string]any{"type": "string"},
			"action":        map[string]any{"type": "string"},
			"resource_type": map[string]any{"type": "string"},
			"page":          map[string]any{"type": "integer"},
			"page_size":     map[string]any{"type": "integer"},
		},
	}
}
func (t *ListAuditLogsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		UserID       string `json:"user_id"`
		Action       string `json:"action"`
		ResourceType string `json:"resource_type"`
		Page         int    `json:"page"`
		PageSize     int    `json:"page_size"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "list_audit_logs: unmarshal input", errx.TypeValidation)
	}
	if in.Page == 0 { in.Page = 1 }
	if in.PageSize == 0 { in.PageSize = 20 }
	filter := audit.AuditFilter{
		UserID:       in.UserID,
		Action:       in.Action,
		ResourceType: in.ResourceType,
	}
	result, err := t.audit.List(ctx, t.tenantID, filter, in.Page, in.PageSize)
	if err != nil {
		return "", errx.Wrap(err, "list_audit_logs", errx.TypeInternal)
	}
	var lines []string
	for _, e := range result.Items {
		lines = append(lines, fmt.Sprintf("- %s: [%s] %s %s/%s by %s", e.CreatedAt.Format(time.RFC3339), e.Action, e.ResourceType, e.ResourceID, e.ID, e.UserID))
	}
	return fmt.Sprintf("Found %d audit entries (page %d/%d):\n%s", result.Total, result.Page, result.TotalPages, strings.Join(lines, "\n")), nil
}

// ─── Loyalty Tools ──────────────────────────────────────────────────────────

type ListLoyaltyAccountsTool struct {
	loyalty  *loyaltysrv.Service
	tenantID kernel.TenantID
}

func (t *ListLoyaltyAccountsTool) Name() string        { return "list_loyalty_accounts" }
func (t *ListLoyaltyAccountsTool) Description() string  { return "List customer loyalty accounts with points and tiers" }
func (t *ListLoyaltyAccountsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"page":      map[string]any{"type": "integer"},
			"page_size": map[string]any{"type": "integer"},
		},
	}
}
func (t *ListLoyaltyAccountsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "list_loyalty_accounts: unmarshal input", errx.TypeValidation)
	}
	if in.Page == 0 { in.Page = 1 }
	if in.PageSize == 0 { in.PageSize = 20 }
	result, err := t.loyalty.ListAccounts(ctx, t.tenantID, kernel.PaginationOptions{Page: in.Page, PageSize: in.PageSize})
	if err != nil {
		return "", errx.Wrap(err, "list_loyalty_accounts", errx.TypeInternal)
	}
	var lines []string
	for _, a := range result.Items {
		lines = append(lines, fmt.Sprintf("- %s: customer=%s, points=%d, tier=%s, lifetime=%d", a.ID, a.CustomerID, a.PointsBalance, a.Tier, a.LifetimePoints))
	}
	return fmt.Sprintf("Found %d loyalty accounts (page %d/%d):\n%s", result.Total, result.Page, result.TotalPages, strings.Join(lines, "\n")), nil
}

type EarnLoyaltyPointsTool struct {
	loyalty  *loyaltysrv.Service
	tenantID kernel.TenantID
}

func (t *EarnLoyaltyPointsTool) Name() string        { return "earn_loyalty_points" }
func (t *EarnLoyaltyPointsTool) Description() string  { return "Award loyalty points to a customer" }
func (t *EarnLoyaltyPointsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"customer_id": map[string]any{"type": "string"},
			"points":      map[string]any{"type": "integer"},
			"reference":   map[string]any{"type": "string", "description": "Order ID or reason"},
			"note":        map[string]any{"type": "string"},
		},
		"required": []string{"customer_id", "points"},
	}
}
func (t *EarnLoyaltyPointsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in loyalty.EarnPointsInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "earn_loyalty_points: unmarshal input", errx.TypeValidation)
	}
	acct, err := t.loyalty.EarnPoints(ctx, t.tenantID, in)
	if err != nil {
		return "", errx.Wrap(err, "earn_loyalty_points", errx.TypeInternal)
	}
	return fmt.Sprintf("Awarded %d points to customer %s. Balance: %d, Tier: %s", in.Points, acct.CustomerID, acct.PointsBalance, acct.Tier), nil
}

type ListLoyaltyRewardsTool struct {
	loyalty  *loyaltysrv.Service
	tenantID kernel.TenantID
}

func (t *ListLoyaltyRewardsTool) Name() string        { return "list_loyalty_rewards" }
func (t *ListLoyaltyRewardsTool) Description() string  { return "List available loyalty rewards" }
func (t *ListLoyaltyRewardsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"page":      map[string]any{"type": "integer"},
			"page_size": map[string]any{"type": "integer"},
		},
	}
}
func (t *ListLoyaltyRewardsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "list_loyalty_rewards: unmarshal input", errx.TypeValidation)
	}
	if in.Page == 0 { in.Page = 1 }
	if in.PageSize == 0 { in.PageSize = 20 }
	result, err := t.loyalty.ListRewards(ctx, t.tenantID, kernel.PaginationOptions{Page: in.Page, PageSize: in.PageSize})
	if err != nil {
		return "", errx.Wrap(err, "list_loyalty_rewards", errx.TypeInternal)
	}
	var lines []string
	for _, r := range result.Items {
		lines = append(lines, fmt.Sprintf("- %s: %s (%d pts, type=%s, value=%d cents)", r.ID, r.Name, r.PointsCost, r.RewardType, r.ValueCents))
	}
	return fmt.Sprintf("Found %d rewards (page %d/%d):\n%s", result.Total, result.Page, result.TotalPages, strings.Join(lines, "\n")), nil
}

type CreateLoyaltyRewardTool struct {
	loyalty  *loyaltysrv.Service
	tenantID kernel.TenantID
}

func (t *CreateLoyaltyRewardTool) Name() string        { return "create_loyalty_reward" }
func (t *CreateLoyaltyRewardTool) Description() string  { return "Create a new loyalty reward" }
func (t *CreateLoyaltyRewardTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":        map[string]any{"type": "string"},
			"description": map[string]any{"type": "string"},
			"points_cost": map[string]any{"type": "integer"},
			"reward_type": map[string]any{"type": "string", "enum": []string{"discount", "free_shipping", "gift_card"}},
			"value_cents": map[string]any{"type": "integer"},
		},
		"required": []string{"name", "points_cost", "reward_type"},
	}
}
func (t *CreateLoyaltyRewardTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in loyalty.CreateRewardInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "create_loyalty_reward: unmarshal input", errx.TypeValidation)
	}
	r, err := t.loyalty.CreateReward(ctx, t.tenantID, in)
	if err != nil {
		return "", errx.Wrap(err, "create_loyalty_reward", errx.TypeInternal)
	}
	return fmt.Sprintf("Created reward %s: %s (%d pts)", r.ID, r.Name, r.PointsCost), nil
}

// ─── Bundle Tools ───────────────────────────────────────────────────────────

type ListBundlesTool struct {
	bundles  *bundlesrv.Service
	tenantID kernel.TenantID
}

func (t *ListBundlesTool) Name() string        { return "list_bundles" }
func (t *ListBundlesTool) Description() string  { return "List product bundles" }
func (t *ListBundlesTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"active_only": map[string]any{"type": "boolean"},
			"page":        map[string]any{"type": "integer"},
			"page_size":   map[string]any{"type": "integer"},
		},
	}
}
func (t *ListBundlesTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		ActiveOnly bool `json:"active_only"`
		Page       int  `json:"page"`
		PageSize   int  `json:"page_size"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "list_bundles: unmarshal input", errx.TypeValidation)
	}
	if in.Page == 0 { in.Page = 1 }
	if in.PageSize == 0 { in.PageSize = 20 }
	result, err := t.bundles.List(ctx, t.tenantID, in.ActiveOnly, kernel.PaginationOptions{Page: in.Page, PageSize: in.PageSize})
	if err != nil {
		return "", errx.Wrap(err, "list_bundles", errx.TypeInternal)
	}
	var lines []string
	for _, b := range result.Items {
		lines = append(lines, fmt.Sprintf("- %s: %s (slug=%s, discount=%s %d, active=%v)", b.ID, b.Name, b.Slug, b.DiscountType, b.DiscountValue, b.Active))
	}
	return fmt.Sprintf("Found %d bundles (page %d/%d):\n%s", result.Total, result.Page, result.TotalPages, strings.Join(lines, "\n")), nil
}

type CreateBundleTool struct {
	bundles  *bundlesrv.Service
	tenantID kernel.TenantID
}

func (t *CreateBundleTool) Name() string        { return "create_bundle" }
func (t *CreateBundleTool) Description() string  { return "Create a new product bundle with discount" }
func (t *CreateBundleTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":           map[string]any{"type": "string"},
			"slug":           map[string]any{"type": "string"},
			"description":    map[string]any{"type": "string"},
			"discount_type":  map[string]any{"type": "string", "enum": []string{"percentage", "fixed"}},
			"discount_value": map[string]any{"type": "integer"},
			"active":         map[string]any{"type": "boolean"},
		},
		"required": []string{"name", "discount_type", "discount_value"},
	}
}
func (t *CreateBundleTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in bundle.CreateBundleInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "create_bundle: unmarshal input", errx.TypeValidation)
	}
	b, err := t.bundles.Create(ctx, t.tenantID, in)
	if err != nil {
		return "", errx.Wrap(err, "create_bundle", errx.TypeInternal)
	}
	return fmt.Sprintf("Created bundle %s: %s (slug=%s)", b.ID, b.Name, b.Slug), nil
}

// ─── Dashboard Tools ────────────────────────────────────────────────────────

type GetSalesOverviewTool struct {
	dashboard *dashboardsrv.Service
	tenantID  kernel.TenantID
}

func (t *GetSalesOverviewTool) Name() string        { return "get_sales_overview" }
func (t *GetSalesOverviewTool) Description() string  { return "Get sales KPIs (revenue, order count, AOV, refunds) for a date range" }
func (t *GetSalesOverviewTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"from": map[string]any{"type": "string", "description": "Start date (YYYY-MM-DD)"},
			"to":   map[string]any{"type": "string", "description": "End date (YYYY-MM-DD)"},
		},
	}
}
func (t *GetSalesOverviewTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		From string `json:"from"`
		To   string `json:"to"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "get_sales_overview: unmarshal input", errx.TypeValidation)
	}
	dr := dashboard.DateRange{}
	if in.From != "" {
		t, _ := time.Parse("2006-01-02", in.From)
		dr.From = t
	}
	if in.To != "" {
		t, _ := time.Parse("2006-01-02", in.To)
		dr.To = t
	}
	s, err := t.dashboard.GetSalesOverview(ctx, t.tenantID, dr)
	if err != nil {
		return "", errx.Wrap(err, "get_sales_overview", errx.TypeInternal)
	}
	return fmt.Sprintf("Sales Overview:\n- Revenue: %d cents (%s)\n- Orders: %d\n- AOV: %d cents\n- Refunds: %d cents", s.TotalRevenue, s.Currency, s.OrderCount, s.AverageOrderValue, s.RefundTotal), nil
}

type GetTopProductsTool struct {
	dashboard *dashboardsrv.Service
	tenantID  kernel.TenantID
}

func (t *GetTopProductsTool) Name() string        { return "get_top_products" }
func (t *GetTopProductsTool) Description() string  { return "Get top-selling products by revenue for a date range" }
func (t *GetTopProductsTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"from":  map[string]any{"type": "string", "description": "Start date (YYYY-MM-DD)"},
			"to":    map[string]any{"type": "string", "description": "End date (YYYY-MM-DD)"},
			"limit": map[string]any{"type": "integer", "description": "Number of products (default 10)"},
		},
	}
}
func (t *GetTopProductsTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		From  string `json:"from"`
		To    string `json:"to"`
		Limit int    `json:"limit"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "get_top_products: unmarshal input", errx.TypeValidation)
	}
	if in.Limit == 0 { in.Limit = 10 }
	dr := dashboard.DateRange{}
	if in.From != "" {
		t, _ := time.Parse("2006-01-02", in.From)
		dr.From = t
	}
	if in.To != "" {
		t, _ := time.Parse("2006-01-02", in.To)
		dr.To = t
	}
	products, err := t.dashboard.GetTopProducts(ctx, t.tenantID, dr, in.Limit)
	if err != nil {
		return "", errx.Wrap(err, "get_top_products", errx.TypeInternal)
	}
	var lines []string
	for i, p := range products {
		lines = append(lines, fmt.Sprintf("%d. %s — revenue: %d cents, qty: %d", i+1, p.Name, p.Revenue, p.Quantity))
	}
	return fmt.Sprintf("Top %d products:\n%s", len(products), strings.Join(lines, "\n")), nil
}

type GetRevenueByDayTool struct {
	dashboard *dashboardsrv.Service
	tenantID  kernel.TenantID
}

func (t *GetRevenueByDayTool) Name() string        { return "get_revenue_by_day" }
func (t *GetRevenueByDayTool) Description() string  { return "Get daily revenue breakdown for a date range" }
func (t *GetRevenueByDayTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"from": map[string]any{"type": "string", "description": "Start date (YYYY-MM-DD)"},
			"to":   map[string]any{"type": "string", "description": "End date (YYYY-MM-DD)"},
		},
	}
}
func (t *GetRevenueByDayTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct {
		From string `json:"from"`
		To   string `json:"to"`
	}
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "get_revenue_by_day: unmarshal input", errx.TypeValidation)
	}
	dr := dashboard.DateRange{}
	if in.From != "" {
		t, _ := time.Parse("2006-01-02", in.From)
		dr.From = t
	}
	if in.To != "" {
		t, _ := time.Parse("2006-01-02", in.To)
		dr.To = t
	}
	days, err := t.dashboard.GetRevenueByDay(ctx, t.tenantID, dr)
	if err != nil {
		return "", errx.Wrap(err, "get_revenue_by_day", errx.TypeInternal)
	}
	var lines []string
	for _, d := range days {
		lines = append(lines, fmt.Sprintf("- %s: %d cents (%d orders)", d.Date, d.Revenue, d.OrderCount))
	}
	return fmt.Sprintf("Daily revenue (%d days):\n%s", len(days), strings.Join(lines, "\n")), nil
}

// ─── Notification Tools ─────────────────────────────────────────────────────

type GetUnreadNotificationCountTool struct {
	notifications *notificationsrv.Service
	tenantID      kernel.TenantID
}

func (t *GetUnreadNotificationCountTool) Name() string        { return "get_unread_notification_count" }
func (t *GetUnreadNotificationCountTool) Description() string  { return "Get the number of unread notifications for a user" }
func (t *GetUnreadNotificationCountTool) InputSchema() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{"user_id": map[string]any{"type": "string"}},
		"required":   []string{"user_id"},
	}
}
func (t *GetUnreadNotificationCountTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct{ UserID string `json:"user_id"` }
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "get_unread_notification_count: unmarshal input", errx.TypeValidation)
	}
	count, err := t.notifications.GetUnreadCount(ctx, t.tenantID, kernel.UserID(in.UserID))
	if err != nil {
		return "", errx.Wrap(err, "get_unread_notification_count", errx.TypeInternal)
	}
	return fmt.Sprintf("User %s has %d unread notifications", in.UserID, count), nil
}

type MarkAllNotificationsReadTool struct {
	notifications *notificationsrv.Service
	tenantID      kernel.TenantID
}

func (t *MarkAllNotificationsReadTool) Name() string        { return "mark_all_notifications_read" }
func (t *MarkAllNotificationsReadTool) Description() string  { return "Mark all notifications as read for a user" }
func (t *MarkAllNotificationsReadTool) InputSchema() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{"user_id": map[string]any{"type": "string"}},
		"required":   []string{"user_id"},
	}
}
func (t *MarkAllNotificationsReadTool) Execute(ctx context.Context, input json.RawMessage) (string, error) {
	var in struct{ UserID string `json:"user_id"` }
	if err := json.Unmarshal(input, &in); err != nil {
		return "", errx.Wrap(err, "mark_all_notifications_read: unmarshal input", errx.TypeValidation)
	}
	err := t.notifications.MarkAllRead(ctx, t.tenantID, kernel.UserID(in.UserID))
	if err != nil {
		return "", errx.Wrap(err, "mark_all_notifications_read", errx.TypeInternal)
	}
	return fmt.Sprintf("All notifications marked as read for user %s", in.UserID), nil
}

// Compile-time interface guards
var (
	_ Tool = (*ListWarehousesTool)(nil)
	_ Tool = (*CreateWarehouseTool)(nil)
	_ Tool = (*AdjustStockTool)(nil)
	_ Tool = (*GetLowStockTool)(nil)
	_ Tool = (*ListReviewsTool)(nil)
	_ Tool = (*ApproveReviewTool)(nil)
	_ Tool = (*RejectReviewTool)(nil)
	_ Tool = (*ListReturnsTool)(nil)
	_ Tool = (*ApproveReturnTool)(nil)
	_ Tool = (*ListWebhooksTool)(nil)
	_ Tool = (*CreateWebhookTool)(nil)
	_ Tool = (*ToggleWebhookTool)(nil)
	_ Tool = (*ListAuditLogsTool)(nil)
	_ Tool = (*ListLoyaltyAccountsTool)(nil)
	_ Tool = (*EarnLoyaltyPointsTool)(nil)
	_ Tool = (*ListLoyaltyRewardsTool)(nil)
	_ Tool = (*CreateLoyaltyRewardTool)(nil)
	_ Tool = (*ListBundlesTool)(nil)
	_ Tool = (*CreateBundleTool)(nil)
	_ Tool = (*GetSalesOverviewTool)(nil)
	_ Tool = (*GetTopProductsTool)(nil)
	_ Tool = (*GetRevenueByDayTool)(nil)
	_ Tool = (*GetUnreadNotificationCountTool)(nil)
	_ Tool = (*MarkAllNotificationsReadTool)(nil)
)
