# Architecture

This document explains how hada-commerce is structured internally — for developers who want to contribute to the platform, build deep integrations, or understand how the pieces fit together.

---

## System Overview

```
┌──────────────────────────────────────────────────────────────────────┐
│                           Clients                                    │
│                                                                      │
│  ┌─────────────────┐  ┌──────────────────┐  ┌────────────────────┐  │
│  │  React Admin    │  │  Custom          │  │  Plugin backends   │  │
│  │  Panel (SPA)    │  │  Storefront      │  │  (webhook recv.)   │  │
│  └────────┬────────┘  └────────┬─────────┘  └─────────┬──────────┘  │
└───────────┼────────────────────┼─────────────────────┼─────────────┘
            │ Bearer token       │ X-Tenant-ID          │ POST events
            ▼                    ▼                       ▼
┌──────────────────────────────────────────────────────────────────────┐
│                    Go Backend  (net/http, :8080)                     │
│                                                                      │
│   Middleware stack:                                                  │
│   ┌─────────┐ → ┌──────────────────┐ → ┌────────────────────────┐  │
│   │  CORS   │   │ Request logging  │   │ Tenant extraction      │  │
│   └─────────┘   └──────────────────┘   │ (X-Tenant-ID → ctx)   │  │
│                                         └────────────────────────┘  │
│                                                                      │
│  ┌───────────────────── Domain Containers ──────────────────────┐   │
│  │ product  order  customer  catalog  storefront  theme  plugin  │   │
│  │ settings  promo  media  analytics  marketplace  iam           │   │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                      │
│  ┌────────────────────── Event Bus ────────────────────────────┐    │
│  │  In-process pub/sub · 21 event types · goroutine fanout     │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                      │
│  ┌──────────────────────── Infra ──────────────────────────────┐    │
│  │  fsx (file storage)  ·  jobx (queue)  ·  notifx (email)    │    │
│  │  logx (logging)      ·  errx (typed errors)                 │    │
│  └─────────────────────────────────────────────────────────────┘    │
└──────────────────────────────┬───────────────────────────────────────┘
                               │
               ┌───────────────┴───────────────┐
               ▼                               ▼
   ┌───────────────────────┐     ┌─────────────────────────┐
   │    PostgreSQL 16       │     │       Redis 7            │
   │  (primary data store)  │     │  (job queue, cache)      │
   │  Multi-tenant schemas  │     │  (background workers)    │
   └───────────────────────┘     └─────────────────────────┘
```

---

## Domain-Driven Design

Each bounded context lives in `backend/internal/<domain>/` and follows a strict 5-layer structure:

```
backend/internal/<domain>/
├── <domain>/              ← Core domain (entities, interfaces, errors)
│   ├── entity.go          ← Structs with typed IDs and kernel types
│   ├── port.go            ← Repository interface (context-first signatures)
│   └── errors.go          ← errx.New() typed error vars
├── <domain>srv/           ← Application service (business logic)
│   └── service.go
├── <domain>infra/         ← Infrastructure (Postgres implementation)
│   └── postgres.go
├── <domain>api/           ← HTTP handlers (net/http)
│   └── handler.go
└── <domain>container/     ← Dependency injection wiring
    └── container.go
```

### Dependency flow

```
HTTP Handler
     │ calls
     ▼
  Service (business logic)
     │ calls (via interface)
     ▼
  Repository interface   ←── defined in domain package
     ▲
     │ implemented by
  PostgresRepository      ←── in infra package
```

The `infra` package depends on the `domain` package (for entity types and the repository interface), not the other way around. The `service` package only knows about the repository interface — it never imports `infra` directly.

---

## Bounded Contexts

### product
Manages the product catalog.

| Layer | Key types |
|-------|-----------|
| Entity | `Product` with `ProductID`, `TenantID`, `Money`, `Status` (draft/active/archived) |
| Repository | `GetByID`, `List`, `Create`, `Update`, `Delete`, `UpdateStock` |
| Service | Validates SKU uniqueness, emits `product.created/updated/deleted` events |
| Handler | `POST /api/v1/products`, `GET /api/v1/products/{id}`, etc. |

### order
Manages order lifecycle.

| Layer | Key types |
|-------|-----------|
| Entity | `Order` with `OrderItem[]`, `OrderStatus` (pending→confirmed→processing→shipped→delivered/cancelled) |
| Service | Validates status transitions, deducts product stock, emits order events |
| Handler | `PUT /api/v1/orders/{id}/status`, `POST /api/v1/orders/{id}/cancel` |

Valid status transitions:

```
pending ──► confirmed ──► processing ──► shipped ──► delivered
   │              │              │
   └──────────────┴──────────────┴────────────────► cancelled
```

### customer
Customer profiles and addresses.

### catalog
- `Category` — hierarchical via `parent_id` (self-referencing FK)
- `Collection` — manual (explicit product list) or automatic (rule-based)

### storefront
Page builder, SSR renderer, block type registry.

| Entity | Description |
|--------|-------------|
| `Page` | Slug-addressed content. `ContentType` = `html` (legacy) or `blocks` |
| `Section` | Full-width page region containing blocks |
| `Block` | Atomic content unit with a type and JSON settings |
| `BlockType` | Schema definition for a block type (built-in or custom) |
| `PageVersion` | Immutable snapshot of page content, created on every save |

Page status lifecycle: `draft` → `pending_review` → `published` → `archived`

### theme
Design token presets. One active theme per tenant. See [themes-guide.md](./themes-guide.md) for token reference.

### plugin
Plugin registry, versioned bundles, per-tenant installations, webhook delivery. See [plugins-guide.md](./plugins-guide.md).

### settings
Per-tenant store configuration (name, currency, checkout rules, social links, branding). Single record per tenant (upserted).

### promo
Discount codes: `percentage` (0–100), `fixed_amount` (cents), `free_shipping`.

### media
File upload and asset management. Files stored via `fsx` (local or S3).

### marketplace
Multi-vendor support: `Vendor` (with approval workflow), `VendorProduct`, `VendorOrder`.

### iam
Identity and access: `User`, `Tenant`, `APIKey`, `OTP`, `Invitation`. Scope-based authorization with wildcard support.

---

## Multi-Tenancy

Every resource is scoped to a tenant. The tenant ID is a string (e.g. `my-store`) provided by the client in the `X-Tenant-ID` header.

```
Request → Middleware reads X-Tenant-ID → stores in context.Context
Handler → kernel.TenantID(TenantFromContext(r.Context()))
Service → passes TenantID to repository
Repository → every SQL query includes WHERE tenant_id = $N
```

There is no cross-tenant data access. Each query is enforced at the database level.

```sql
-- Every query looks like this
SELECT id, name, price_cents, ...
FROM products
WHERE tenant_id = $1   -- always scoped
  AND id = $2
```

Tenant IDs are stored as `TEXT` (not UUIDs) to allow human-readable slugs.

---

## Event-Driven Architecture

### Event Bus

An in-process publish/subscribe system (`backend/internal/eventbus`). Handlers subscribe to event types and run in goroutines.

```go
// Publishing an event (done by service layer)
bus.Publish(eventbus.Event{
    ID:       uuid.New().String(),
    Type:     eventbus.ProductCreated,
    TenantID: tenantID,
    Payload:  productPayload,
    Timestamp: time.Now().UTC(),
})

// Subscribing (done by containers at startup)
bus.Subscribe(eventbus.ProductCreated, func(e eventbus.Event) {
    // handle event
})

// Subscribe to all events (used by webhook dispatcher)
bus.SubscribeAll(func(e eventbus.Event) {
    // handle all events
})
```

### Webhook Dispatcher

The plugin container starts a `WebhookDispatcher` that subscribes to all events and delivers them to active plugin installations:

```
Domain event published
        │
        ▼
  WebhookDispatcher.handleEvent()
        │
        ├── find all active PluginInstallations for event.TenantID
        │
        └── for each installation:
              POST event JSON → installation.backend_entry
              timeout: 5 seconds
              no retries (fire-and-forget)
```

### Event Types (21 total)

```
order.placed         order.confirmed      order.shipped
order.delivered      order.cancelled
customer.registered  customer.updated
product.created      product.updated      product.deleted
category.created     collection.updated
page.published       page.unpublished
theme.activated      theme.updated
settings.updated
plugin.installed     plugin.uninstalled
cart.updated
```

---

## SSR Renderer Pipeline

```
GET /api/v1/storefront/pages/by-slug/{slug}
  Accept: text/html
        │
        ▼
1. Load Page from DB (status must be "published")
        │
        ▼
2. Resolve active Theme for tenant
   (falls back to default tokens if no active theme)
        │
        ▼
3. TokensToCSS(theme.tokens)
   → ":root { --color-primary: ...; ... }"
        │
        ▼
4. Render page body:
   if ContentType == "blocks":
     for each section in page.Sections:
       for each block in section.Blocks:
         render block template (go template or HTML string)
         → wrap in <div class="section-wrapper" ...>
   else (ContentType == "html"):
     use page.HTML directly
        │
        ▼
5. Assemble HTML5 document:
   <html>
     <head>
       <title>{page.Title}</title>
       <meta property="og:..." />
       <style>:root { ... CSS vars ... }</style>
       <style>{ page.CSS }</style>    ← per-page custom CSS
     </head>
     <body>
       { rendered page body }
     </body>
   </html>
        │
        ▼
6. Return HTML response (Content-Type: text/html)
```

---

## HTTP Handler Pattern

All handlers follow this exact pattern:

```go
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    // 1. Extract tenant (returns "" if not set)
    tenantID := kernel.TenantID(TenantFromContext(r.Context()))
    if tenantID == "" {
        writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing tenant")
        return
    }

    // 2. Decode request body
    var req widget.CreateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
        return
    }

    // 3. Call service
    result, err := h.svc.Create(r.Context(), tenantID, req)
    if err != nil {
        writeErrx(w, err)   // maps errx typed errors → HTTP status
        return
    }

    // 4. Write response
    writeJSON(w, http.StatusCreated, result)
}
```

### JSON helpers (in every handler package)

```go
func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
    writeJSON(w, status, map[string]string{"error": code, "message": message})
}

func writeErrx(w http.ResponseWriter, err error) {
    writeError(w, errx.HTTPStatus(err), errx.Code(err), errx.Message(err))
}
```

---

## Error Handling

Errors are typed using `errx`:

```go
// Defined at package level in errors.go
var (
    ErrNotFound     = errx.New("PRODUCT_NOT_FOUND",    "product not found",    http.StatusNotFound)
    ErrConflict     = errx.New("PRODUCT_CONFLICT",     "product already exists", http.StatusConflict)
    ErrInvalidInput = errx.New("PRODUCT_INVALID_INPUT","invalid input",        http.StatusBadRequest)
)
```

In the infra layer, `sql.ErrNoRows` is wrapped into the domain error:

```go
if err == sql.ErrNoRows {
    return Product{}, errx.Wrap(ErrNotFound, string(id))
}
```

In the handler layer, `writeErrx` maps the error to an HTTP response:

```go
// errx.HTTPStatus(err) → 404
// errx.Code(err)       → "PRODUCT_NOT_FOUND"
// errx.Message(err)    → "product not found"
```

JSON response body:
```json
{ "error": "PRODUCT_NOT_FOUND", "message": "product not found" }
```

---

## Kernel Types

All business types are in `backend/internal/kernel`. Never use raw strings for IDs.

```go
// IDs (all string-backed, distinct Go types)
kernel.ProductID     kernel.OrderID       kernel.CustomerID
kernel.CategoryID    kernel.CollectionID  kernel.PageID
kernel.PageVersionID kernel.ThemeID       kernel.BlockTypeID
kernel.PluginID      kernel.PluginVersionID  kernel.InstallationID
kernel.MediaID       kernel.PromoID       kernel.VendorID
kernel.TenantID      kernel.UserID

// Money (always in smallest currency unit — cents for USD)
kernel.Money{Amount: 2999, Currency: "USD"}
kernel.NewMoney(2999, "USD")

// Pagination
kernel.NewPagination(page, pageSize)      // defaults: page=1, size=20, max=100
kernel.PaginatedResult[T]
kernel.NewPaginatedResult(items, total, p)
```

---

## Database Conventions

| Convention | Rule |
|-----------|------|
| IDs | `TEXT` columns, UUID values stored as strings |
| Tenant scoping | Every table has `tenant_id TEXT NOT NULL` |
| Money | `BIGINT` for amounts (cents), `TEXT` for currency code |
| Timestamps | `TIMESTAMPTZ NOT NULL DEFAULT NOW()` |
| Soft JSON | `JSONB` for flexible settings, addresses, social links |
| Indexes | Always index `tenant_id`; add composite indexes for common queries |
| Migrations | Sequential numbered files: `001_init.up.sql`, `002_...up.sql` — never edited, only appended |

---

## Infrastructure Packages

| Package | Purpose |
|---------|---------|
| `fsx` | File storage abstraction (local FS or S3) |
| `jobx` | Background job queue (Redis-backed) |
| `notifx` | Notification dispatching (console or AWS SES) |
| `logx` | Structured logging |
| `asyncx` | Async primitives (Future, Pool, Retry, Race) |
| `errx` | Typed error system with HTTP status codes |

---

## Directory Layout

```
hada-commerce/
├── backend/
│   ├── cmd/
│   │   ├── main.go          ← Entry point
│   │   ├── server.go        ← HTTP server setup, middleware, route registration
│   │   └── container.go     ← Wires all domain containers
│   ├── internal/
│   │   ├── kernel/          ← Core types (IDs, Money, Pagination, AuthContext)
│   │   ├── eventbus/        ← In-process event bus
│   │   ├── errx/            ← Typed error system
│   │   ├── product/         ← Product domain (5 layers)
│   │   ├── order/           ← Order domain
│   │   ├── customer/        ← Customer domain
│   │   ├── catalog/         ← Category + Collection domain
│   │   ├── storefront/      ← Page builder + SSR renderer
│   │   ├── theme/           ← Theme + design tokens
│   │   ├── plugin/          ← Plugin system + webhook dispatcher
│   │   ├── settings/        ← Store settings
│   │   ├── promo/           ← Discount codes
│   │   ├── media/           ← File uploads
│   │   ├── analytics/       ← Event tracking
│   │   ├── marketplace/     ← Multi-vendor
│   │   └── iam/             ← Identity and access management
│   └── migrations/
│       ├── 001_init.up.sql
│       ├── 002_catalog.up.sql
│       └── ...
├── frontend/
│   └── src/
│       └── ...              ← React admin panel
└── docker-compose.yml
```

---

## Adding a New Domain

Follow the 5-layer pattern. Use existing domains as reference (e.g. `backend/internal/product/`).

1. Create the 5 packages: `<entity>`, `<entity>srv`, `<entity>infra`, `<entity>api`, `<entity>container`
2. Add a numbered migration in `backend/migrations/`
3. Register the container in `backend/cmd/container.go`
4. Call `RegisterRoutes(mux)` in `backend/cmd/server.go`
5. Run `go build ./... && go vet ./...` to verify

### Checklist

- [ ] Entity struct uses `kernel.<Entity>ID` for ID field
- [ ] Entity struct has `TenantID kernel.TenantID`
- [ ] Repository interface in `port.go` (context-first signatures)
- [ ] Errors in `errors.go` using `errx.New("DOMAIN_DESC", "msg", status)`
- [ ] Service depends only on the Repository **interface**
- [ ] `postgres.go` wraps `sql.ErrNoRows` → domain `ErrNotFound`
- [ ] Every SQL query filters by `tenant_id`
- [ ] List endpoints return `kernel.PaginatedResult[T]`
- [ ] Handler extracts tenantID and checks for empty string
- [ ] Migration file is the next sequential number, never edits existing files
