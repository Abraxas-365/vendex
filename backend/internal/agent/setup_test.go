package agent

import (
	"encoding/json"
	"testing"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// TestSetup_ReturnsAllTools verifies that Setup() with nil services returns
// all tool structs (>= 89) without panicking. Setup only constructs structs,
// it never calls service methods, so nil service pointers are safe.
func TestSetup_ReturnsAllTools(t *testing.T) {
	tools := Setup(kernel.TenantID("test-tenant"), Services{})

	if len(tools) < 89 {
		t.Errorf("expected at least 89 tools, got %d", len(tools))
	}

	for i, tool := range tools {
		if tool == nil {
			t.Errorf("tool at index %d is nil", i)
		}
	}
}

// TestSetup_UniqueToolNames verifies that no two tools share the same name.
func TestSetup_UniqueToolNames(t *testing.T) {
	tools := Setup(kernel.TenantID("test-tenant"), Services{})

	seen := make(map[string]int)
	for i, tool := range tools {
		name := tool.Name()
		if prev, exists := seen[name]; exists {
			t.Errorf("duplicate tool name %q at indices %d and %d", name, prev, i)
		}
		seen[name] = i
	}
}

// TestSetup_AllToolNamesNonEmpty verifies every tool has a non-empty Name().
func TestSetup_AllToolNamesNonEmpty(t *testing.T) {
	tools := Setup(kernel.TenantID("test-tenant"), Services{})

	for i, tool := range tools {
		if tool.Name() == "" {
			t.Errorf("tool at index %d has empty Name()", i)
		}
	}
}

// TestSetup_AllToolDescriptionsNonEmpty verifies every tool has a non-empty Description().
func TestSetup_AllToolDescriptionsNonEmpty(t *testing.T) {
	tools := Setup(kernel.TenantID("test-tenant"), Services{})

	for i, tool := range tools {
		if tool.Description() == "" {
			t.Errorf("tool %q (index %d) has empty Description()", tool.Name(), i)
		}
	}
}

// TestSetup_AllInputSchemasValid verifies every tool's InputSchema() contains
// "type": "object" and a "properties" key, as required by the JSON Schema spec
// for tool inputs.
func TestSetup_AllInputSchemasValid(t *testing.T) {
	tools := Setup(kernel.TenantID("test-tenant"), Services{})

	for i, tool := range tools {
		schema := tool.InputSchema()
		if schema == nil {
			t.Errorf("tool %q (index %d) returned nil InputSchema()", tool.Name(), i)
			continue
		}

		// Marshal to JSON and back so we can inspect keys uniformly.
		raw, err := json.Marshal(schema)
		if err != nil {
			t.Errorf("tool %q (index %d): InputSchema() could not be marshalled: %v", tool.Name(), i, err)
			continue
		}

		var parsed map[string]any
		if err := json.Unmarshal(raw, &parsed); err != nil {
			t.Errorf("tool %q (index %d): InputSchema() is not valid JSON: %v", tool.Name(), i, err)
			continue
		}

		// Must have "type": "object"
		typeVal, ok := parsed["type"]
		if !ok {
			t.Errorf("tool %q (index %d): InputSchema() missing 'type' key", tool.Name(), i)
		} else if typeVal != "object" {
			t.Errorf("tool %q (index %d): InputSchema() 'type' is %q, want 'object'", tool.Name(), i, typeVal)
		}

		// Must have "properties" key
		if _, ok := parsed["properties"]; !ok {
			t.Errorf("tool %q (index %d): InputSchema() missing 'properties' key", tool.Name(), i)
		}
	}
}

// TestSetup_KnownToolNames verifies that all expected tool names are present
// in the Setup() output. This acts as a regression guard: if a tool is
// accidentally removed or renamed, this test will catch it.
func TestSetup_KnownToolNames(t *testing.T) {
	tools := Setup(kernel.TenantID("test-tenant"), Services{})

	// Build a lookup map.
	nameSet := make(map[string]bool, len(tools))
	for _, tool := range tools {
		nameSet[tool.Name()] = true
	}

	required := []string{
		// Storefront / CMS
		"create_page",
		"update_page",
		"list_pages",
		"list_block_types",
		"create_block_type",

		// Products
		"create_product",
		"list_products",

		// Promos
		"create_promo",

		// Orders
		"query_orders",

		// Catalog
		"search_catalog",

		// Themes
		"list_themes",
		"get_active_theme",
		"create_theme",
		"update_theme",
		"activate_theme",

		// Shipping
		"list_shipping_zones",
		"create_shipping_zone",
		"calculate_shipping",

		// Tax
		"list_tax_rates",
		"create_tax_rate",
		"calculate_tax",

		// Payments
		"get_order_payment",
		"list_refunds",

		// Search
		"search_products",
		"search_suggestions",

		// Product variants
		"create_product_option",
		"list_product_options",
		"create_product_variant",
		"list_product_variants",

		// Customer groups
		"list_customer_groups",
		"create_customer_group",
		"add_group_member",

		// Gift cards
		"list_gift_cards",
		"create_gift_card",
		"check_gift_card_balance",
		"redeem_gift_card",

		// Cart recovery
		"list_recovery_emails",
		"get_recovery_stats",

		// Currency
		"list_currency_rates",
		"set_currency_rate",
		"convert_currency",

		// I18n
		"set_translations",
		"get_translations",
		"list_supported_locales",

		// Subscriptions
		"list_subscriptions",
		"create_subscription",
		"cancel_subscription",

		// Inventory
		"list_warehouses",
		"create_warehouse",
		"adjust_stock",
		"get_low_stock",

		// Reviews
		"list_reviews",
		"approve_review",
		"reject_review",

		// Returns
		"list_returns",
		"approve_return",

		// Webhooks
		"list_webhooks",
		"create_webhook",
		"toggle_webhook",

		// Audit
		"list_audit_logs",

		// Loyalty
		"list_loyalty_accounts",
		"earn_loyalty_points",
		"list_loyalty_rewards",
		"create_loyalty_reward",

		// Bundles
		"list_bundles",
		"create_bundle",

		// Dashboard
		"get_sales_overview",
		"get_top_products",
		"get_revenue_by_day",

		// Notifications
		"get_unread_notification_count",
		"mark_all_notifications_read",

		// Multi-storefront
		"list_storefronts",
		"create_storefront",

		// Bulk operations
		"list_bulk_operations",
		"create_bulk_operation",

		// Blog
		"list_blog_posts",
		"create_blog_post",
		"publish_blog_post",

		// Collections
		"list_collections",
		"create_collection",
		"add_collection_product",

		// A/B Testing
		"list_experiments",
		"create_experiment",
		"get_experiment_results",

		// Recommendations
		"list_recommendation_rules",
		"get_trending_products",

		// Agent Memory
		"search_memory",
		"save_memory",
		"get_memory_context",
	}

	for _, name := range required {
		if !nameSet[name] {
			t.Errorf("expected tool %q to be in Setup() output, but it was not found", name)
		}
	}
}
