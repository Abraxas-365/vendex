# Plugins Guide

The hada-commerce plugin system lets you extend the platform with custom business logic, third-party integrations, and UI widgets — without modifying the core codebase. Plugins receive webhook events for every domain action and can inject JavaScript into the admin panel.

---

## Plugin Architecture

```
Plugin Registry (global)
   │
   ├── Plugin          ← metadata: name, author, category
   │    └── PluginVersion ← versioned bundle: manifest, frontend URL, backend entry
   │
   └── PluginInstallation (per-tenant)
        ├── status: active | inactive | failed
        ├── settings: tenant-specific config JSON
        └── version: the installed PluginVersion
```

### Three Core Entities

| Entity | Scope | Description |
|--------|-------|-------------|
| `Plugin` | Global | The plugin catalog entry. Created once by the plugin developer. |
| `PluginVersion` | Global | A versioned release of a plugin with a manifest and code bundles. |
| `PluginInstallation` | Per-tenant | When a tenant installs a plugin, an installation record is created with per-tenant settings. |

---

## Plugin Lifecycle

```
Developer                         Merchant
    │                                 │
    │  POST /api/v1/plugins           │
    │  (create Plugin)                │
    │                                 │
    │  POST /api/v1/plugins/{id}/versions
    │  (publish PluginVersion)        │
    │                                 │
    │                    POST /api/v1/plugins/{id}/install
    │                    (create PluginInstallation)
    │                                 │
    │                    PUT /api/v1/plugins/installations/{id}/settings
    │                    (configure tenant settings)
    │                                 │
    │                    POST /api/v1/plugins/installations/{id}/activate
    │                    (set status = active)
    │                                 │
    │                    DELETE /api/v1/plugins/{id}/uninstall
    │                    (remove installation)
```

---

## Plugin Manifest

Every `PluginVersion` carries a **manifest** that describes the plugin's capabilities:

```json
{
  "name": "my-analytics",
  "display_name": "My Analytics",
  "version": "1.2.0",
  "description": "Track store events and build custom dashboards.",
  "author": "Acme Corp",
  "permissions": [
    "orders:read",
    "products:read",
    "customers:read"
  ],
  "ui": {
    "tabs": [
      {
        "label": "Analytics",
        "icon": "bar-chart",
        "entry": "/tabs/dashboard.html"
      }
    ],
    "widgets": [
      {
        "slot": "dashboard",
        "component": "RevenueWidget",
        "entry": "/widgets/revenue.js"
      },
      {
        "slot": "product-detail",
        "component": "ViewsWidget",
        "entry": "/widgets/views.js"
      }
    ]
  }
}
```

### UI Widget Slots

The admin panel exposes these injection points:

| Slot | Location |
|------|---------|
| `dashboard` | Main dashboard page |
| `product-detail` | Product edit page sidebar |
| `checkout` | Checkout flow |
| `order-detail` | Order detail page |

---

## Creating a Plugin (Developer Guide)

### Step 1 — Register the Plugin

```bash
curl -X POST http://localhost:8080/api/v1/plugins \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-analytics",
    "display_name": "My Analytics",
    "description": "Track store events and build dashboards.",
    "author": "Acme Corp",
    "category": "analytics",
    "tags": ["analytics", "reporting"]
  }'
```

### Step 2 — Publish a Version

```bash
curl -X POST http://localhost:8080/api/v1/plugins/{plugin_id}/versions \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "version": "1.0.0",
    "changelog": "Initial release",
    "permissions": ["orders:read", "products:read"],
    "frontend_url": "https://cdn.acmecorp.com/plugins/my-analytics/1.0.0",
    "backend_entry": "https://api.acmecorp.com/webhooks/hada",
    "min_platform_ver": "1.0.0",
    "manifest": {
      "name": "my-analytics",
      "display_name": "My Analytics",
      "version": "1.0.0",
      "ui": {
        "tabs": [
          { "label": "Analytics", "icon": "bar-chart", "entry": "/tabs/dashboard.html" }
        ]
      }
    }
  }'
```

The `backend_entry` URL is where the platform will POST webhook events. It must be publicly reachable from the server.

### Step 3 — Implement the Webhook Receiver

Your backend must accept POST requests at the `backend_entry` URL. The platform sends a JSON body for every subscribed event:

```
POST https://api.acmecorp.com/webhooks/hada
Content-Type: application/json

{
  "id": "evt_abc123",
  "type": "order.placed",
  "tenant_id": "my-store",
  "timestamp": "2024-01-15T10:30:00Z",
  "payload": {
    "order_id": "ord_xyz789",
    "customer_id": "cus_def456",
    "status": "pending",
    "total": 5999,
    "currency": "USD",
    "item_count": 2
  }
}
```

Return any 2xx status to acknowledge the event. Events are delivered **fire-and-forget** with a 5-second timeout — there are no automatic retries. Plan for idempotency in your webhook handler.

---

## Webhook Event Reference

The platform emits 21 event types. Your webhook endpoint receives all events for tenants that have your plugin installed and active.

### Order Events

| Event Type | Trigger | Payload fields |
|-----------|---------|----------------|
| `order.placed` | Customer submits an order | `order_id`, `customer_id`, `status`, `total`, `currency`, `item_count` |
| `order.confirmed` | Order status → confirmed | `order_id`, `customer_id`, `status`, `total`, `currency`, `item_count` |
| `order.shipped` | Order status → shipped | `order_id`, `customer_id`, `status`, `total`, `currency`, `item_count` |
| `order.delivered` | Order status → delivered | `order_id`, `customer_id`, `status`, `total`, `currency`, `item_count` |
| `order.cancelled` | Order cancelled | `order_id`, `customer_id`, `status`, `total`, `currency`, `item_count` |

**Order payload example:**
```json
{
  "order_id": "ord_xyz789",
  "customer_id": "cus_def456",
  "status": "placed",
  "total": 5999,
  "currency": "USD",
  "item_count": 2
}
```

### Customer Events

| Event Type | Trigger | Payload fields |
|-----------|---------|----------------|
| `customer.registered` | New customer created | `customer_id`, `email`, `name` |
| `customer.updated` | Customer profile changed | `customer_id`, `email`, `name` |

**Customer payload example:**
```json
{
  "customer_id": "cus_def456",
  "email": "jane@example.com",
  "name": "Jane Smith"
}
```

### Product Events

| Event Type | Trigger | Payload fields |
|-----------|---------|----------------|
| `product.created` | New product created | `product_id`, `name`, `sku`, `price`, `currency` |
| `product.updated` | Product edited | `product_id`, `name`, `sku`, `price`, `currency` |
| `product.deleted` | Product deleted | `product_id`, `name`, `sku`, `price`, `currency` |

**Product payload example:**
```json
{
  "product_id": "prd_abc123",
  "name": "Classic T-Shirt",
  "sku": "TSHIRT-001",
  "price": 2999,
  "currency": "USD"
}
```

### Catalog Events

| Event Type | Trigger | Payload fields |
|-----------|---------|----------------|
| `category.created` | New category created | `category_id`, `name`, `slug` |
| `collection.updated` | Collection changed | `collection_id`, `name`, `slug` |

### Storefront Events

| Event Type | Trigger | Payload fields |
|-----------|---------|----------------|
| `page.published` | Page status → published | `page_id`, `slug`, `title` |
| `page.unpublished` | Page status → archived | `page_id`, `slug`, `title` |

### Theme Events

| Event Type | Trigger | Payload fields |
|-----------|---------|----------------|
| `theme.activated` | A theme is activated | `theme_id`, `name` |
| `theme.updated` | Theme tokens changed | `theme_id`, `name` |

### Settings Events

| Event Type | Trigger | Payload fields |
|-----------|---------|----------------|
| `settings.updated` | Store settings changed | `fields` (array of changed field names) |

**Settings payload example:**
```json
{
  "fields": ["store_name", "currency"]
}
```

### Plugin Events

| Event Type | Trigger | Payload fields |
|-----------|---------|----------------|
| `plugin.installed` | Plugin installed by a tenant | `plugin_id`, `plugin_name`, `version` |
| `plugin.uninstalled` | Plugin removed by a tenant | `plugin_id`, `plugin_name`, `version` |

### Cart Events

| Event Type | Trigger | Payload fields |
|-----------|---------|----------------|
| `cart.updated` | Cart contents changed | `cart_id`, `customer_id`, `item_count`, `total`, `currency` |

---

## Plugin Settings Schema

When publishing a version, include a `config_schema` to tell the admin panel what settings to collect from the merchant:

```json
{
  "config_schema": {
    "api_key": {
      "type": "string",
      "label": "API Key",
      "description": "Your Acme Analytics API key",
      "required": true,
      "secret": true
    },
    "track_page_views": {
      "type": "boolean",
      "label": "Track page views",
      "default": true
    },
    "retention_days": {
      "type": "integer",
      "label": "Data retention (days)",
      "default": 90,
      "min": 7,
      "max": 365
    }
  }
}
```

The merchant fills out these fields in the admin panel when configuring the plugin. The values are stored in the `PluginInstallation.settings` JSON field and can be retrieved by your backend.

---

## Frontend Widget Injection

If your plugin's manifest declares `ui.widgets`, the admin panel loads the JavaScript bundle from `frontend_url + entry` and mounts the component into the specified slot.

**Example widget entry point** (`/widgets/revenue.js`):

```javascript
// This file is served from your CDN at frontend_url + entry
// The platform calls mount(container, context) when the slot is ready

export function mount(container, context) {
  const { tenantId, settings } = context;

  container.innerHTML = `
    <div class="revenue-widget">
      <h3>Today's Revenue</h3>
      <div id="revenue-amount">Loading...</div>
    </div>
  `;

  fetch(`https://api.acmecorp.com/revenue?tenant=${tenantId}`, {
    headers: { 'X-API-Key': settings.api_key }
  })
    .then(r => r.json())
    .then(data => {
      document.getElementById('revenue-amount').textContent =
        `$${(data.total_cents / 100).toFixed(2)}`;
    });
}

export function unmount(container) {
  container.innerHTML = '';
}
```

The `context` object provided by the platform includes:

| Field | Type | Description |
|-------|------|-------------|
| `tenantId` | string | The current tenant ID |
| `settings` | object | The plugin's installation settings |
| `locale` | string | Admin panel locale (e.g. `"en-US"`) |

---

## Installing a Plugin (Merchant Guide)

### Via API

```bash
# 1. Install the plugin
curl -X POST http://localhost:8080/api/v1/plugins/{plugin_id}/install \
  -H "X-Tenant-ID: my-store" \
  -H "Authorization: Bearer <token>"

# 2. Configure settings
curl -X PUT http://localhost:8080/api/v1/plugins/installations/{installation_id}/settings \
  -H "X-Tenant-ID: my-store" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "api_key": "sk_live_abc123",
    "track_page_views": true
  }'

# 3. Activate
curl -X POST http://localhost:8080/api/v1/plugins/installations/{installation_id}/activate \
  -H "X-Tenant-ID: my-store" \
  -H "Authorization: Bearer <token>"
```

### Via Admin Panel

Go to **Plugins** in the admin panel, browse the plugin catalog, click **Install**, fill in the settings form, and click **Activate**.

---

## Uninstalling a Plugin

```bash
curl -X DELETE http://localhost:8080/api/v1/plugins/{plugin_id}/uninstall \
  -H "X-Tenant-ID: my-store" \
  -H "Authorization: Bearer <token>"
```

This removes the `PluginInstallation` record and stops webhook delivery for that tenant. The `plugin.uninstalled` event is emitted so you can clean up any remote data.

---

## Webhook Security

**Tip:** To verify that webhook requests are genuinely from your hada-commerce instance (and not a spoofed request), validate a shared secret in the request headers. When registering your `backend_entry`, also store a secret; include it as a custom header in your receiver and reject requests that don't match.

This is a convention — the platform does not enforce webhook signing out-of-the-box in the current version.

---

## Plugin JS Manifest (Public Endpoint)

The admin panel fetches the list of JS bundles to inject from a public endpoint:

```bash
curl http://localhost:8080/api/v1/plugins/js-manifest \
  -H "X-Tenant-ID: my-store"
```

Response:

```json
{
  "scripts": [
    {
      "plugin_id": "plg_abc",
      "plugin_name": "my-analytics",
      "frontend_url": "https://cdn.acmecorp.com/plugins/my-analytics/1.0.0",
      "widgets": [
        { "slot": "dashboard", "entry": "/widgets/revenue.js" }
      ]
    }
  ]
}
```
