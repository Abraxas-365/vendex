# Hada Commerce — Session Context

> Use this document to bootstrap a new AI coding session with full project context.

## Project Overview

**hada-commerce** is a fully customizable e-commerce platform comparable to Shopify/VTEX/Magento/WooCommerce. It features 57+ DDD bounded contexts, 42 admin pages, 80+ AI agent tools, and a harness-powered AI assistant.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.25, Fiber v2, PostgreSQL 16, sqlx, raw SQL (lib/pq) |
| Frontend | React 19, TypeScript, Vite 8, TanStack Router, TanStack Query, Radix UI Themes, Tailwind v4 |
| AI Agent | github.com/Abraxas-365/harness (local dep via `replace` directive) |
| Infra | docker-compose (postgres:16, redis:7), Makefile |
| Module | `github.com/Abraxas-365/hada-commerce` |

## Directory Layout

```
backend/
  cmd/                        — main.go, container.go (DI), server.go (routing)
  migrations/                 — 001-042 sequential .up.sql/.down.sql files
  internal/
    kernel/                   — ID types (69 types), Money, AuthContext
    errx/                     — Domain error handling (Wrap, NotFound, Validation, Business)
    eventbus/                 — In-memory event bus (68 event types + typed payloads)
    config/                   — App config from env vars
    logx/, ptrx/, asyncx/    — Utility packages
    fsx/                      — File storage abstraction
    iam/                      — Admin authentication + authorization
    notifx/                   — Email/notification abstraction (SES)
    agent/                    — AI agent tools + harness integration
      adapter.go              — agent.Tool -> harness tools.Tool adapter
      handler.go              — EventHandler implementing query.EventHandler
      setup.go                — Services struct (31 fields) + Setup() returns []Tool
      tools.go                — 15 core tools (storefront, products, promos, orders, catalog, themes)
      tools_commerce.go       — 32 tools (shipping, tax, payment, search, variants, groups, gifts, cart recovery, currency, i18n, subscriptions)
      tools_p4.go             — 24 tools (inventory, reviews, returns, webhooks, audit, loyalty, bundles, dashboard, social auth, notifications)
      tools_p5.go             — 18 tools (multistore, bulkops, blog, collections, abtest, recommendations)
      agentapi/handler.go     — SSE streaming POST /agent/chat endpoint
    <domain>/                 — 57+ DDD bounded contexts (see list below)

frontend/
  src/
    pages/admin/              — 42 admin pages (Dashboard, Products, Orders, etc.)
    types/index.ts            — All TypeScript types (1306 lines)
    lib/api.ts                — API client functions (1559 lines)
    lib/hooks.ts              — TanStack Query hooks (2494 lines)
    routeTree.tsx             — Route definitions (679 lines)
    components/               — Shared components (layouts, sidebar, etc.)
  e2e/screenshots.spec.ts    — Playwright spec for 45 pages
  screenshots/index.html      — HTML gallery of all screenshots
```

## DDD Pattern (per domain)

Every domain follows this 5-layer pattern:
```
internal/<domain>/
  entity.go         — Domain entities, value objects
  errors.go         — errx-based domain errors
  port.go           — Repository interface (Port)
  <domain>srv/
    service.go      — Business logic service
  <domain>infra/
    postgres.go     — PostgreSQL repository implementation
  <domain>api/
    handler.go      — Fiber HTTP handlers
  <domain>container/
    container.go    — DI container: New(db, bus) -> {Service, Handler}
```

## All 57+ Domains

| # | Domain | Package | Migration | Key Features |
|---|--------|---------|-----------|-------------|
| 1 | Product | `product` | 001 | CRUD, has_variants flag |
| 2 | Order | `order` | 001 | CRUD, status workflow, subtotal/shipping/tax/discount amounts |
| 3 | Customer | `customer` | 001 | CRUD, addresses, customer auth (register/login/JWT) |
| 4 | Catalog | `catalog` | 001 | Categories, collections |
| 5 | Storefront | `storefront` | 001, 005, 008, 009 | Pages, blocks, sections, navigation menus, template overrides |
| 6 | Promo | `promo` | 001, 020 | Codes, validation, usage tracking, advanced targeting, buy-X-get-Y |
| 7 | Media | `media` | 001 | File uploads, S3 storage |
| 8 | Marketplace | `marketplace` | 002 | Plugin registry |
| 9 | Analytics | `analytics` | - | Read-only dashboard stats |
| 10 | Settings | `settings` | 003 | Key-value store settings |
| 11 | IAM | `iam` | 004 | Admin auth, roles, permissions |
| 12 | Theme | `theme` | 006 | Design tokens, theme CRUD, activation |
| 13 | Cart | `cart` | 010 | Cart + items, public routes |
| 14 | Shipping | `shipping` | 011 | Zones, rates, calculate endpoint |
| 15 | Tax | `tax` | 012 | Rates, compound tax, calculate endpoint |
| 16 | Payment | `payment` | 013 | Create/process/refund, manual provider |
| 17 | Checkout | `checkout` | 014 | Preview + process orchestrator (wraps cart, order, shipping, tax, promo, payment) |
| 18 | Product Variants | Extended `product` | 015 | Options, variants with own SKU/price/stock |
| 19 | Search | `search` | 016 | PostgreSQL tsvector + GIN index, autocomplete |
| 20 | Customer Auth | `customer/auth` | 017 | Register, login, JWT, profile, order history |
| 21 | Emails | `emails` | - | Event-driven transactional emails via notifx |
| 22 | Wishlist | `wishlist` | 018 | Customer wishlists + items |
| 23 | SEO | Extended `product` | 019 | Product SEO fields (meta title, description, slug) |
| 24 | Advanced Promos | Extended `promo` | 020 | Targeting rules, buy-X-get-Y, max discount cap |
| 25 | Customer Groups | `customergroup` | 021 | Groups, memberships, bulk assignment |
| 26 | Gift Cards | `giftcard` | 022 | Create, redeem, check balance, transaction history |
| 27 | Cart Recovery | `cartrecovery` | 023 | Abandoned cart detection, recovery emails |
| 28 | Currency | `currency` | 024 | Exchange rates, conversion |
| 29 | i18n | `i18n` | 025 | Multi-language translations, supported locales |
| 30 | Subscriptions | `subscription` | 026 | Recurring billing, plan management |
| 31 | Inventory | `inventory` | 027 | Warehouses, stock levels, adjustments, low-stock alerts |
| 32 | Reviews | `review` | 028 | Ratings, moderation (approve/reject) |
| 33 | Returns | `returns` | 029 | Return requests, approval workflow |
| 34 | Webhooks | `webhook` | 030 | Event subscriptions, delivery tracking |
| 35 | Audit | `audit` | 031 | Audit logs with actor/action/resource tracking |
| 36 | Loyalty | `loyalty` | 032 | Points, rewards, earn/redeem |
| 37 | Bundles | `bundle` | 033 | Product bundles with discount |
| 38 | Dashboard | `dashboard` | 034 | Sales overview, top products, revenue by day |
| 39 | Social Auth | `socialauth` | 035 | OAuth providers (Google, Facebook, etc.) |
| 40 | Notifications | `notification` | 036 | In-app notification center |
| 41 | Multi-Storefront | `multistore` | 037 | Multiple storefronts per tenant |
| 42 | Bulk Operations | `bulkops` | 038 | Async bulk import/update/delete |
| 43 | Blog | `blog` | 039 | Blog posts, CMS |
| 44 | Collections | `collection` | 040 | Product collection management |
| 45 | A/B Testing | `abtest` | 041 | Experiments, variants, results |
| 46 | Recommendations | `recommendation` | 042 | Rules-based product recommendations, trending |
| + | Plugin Runtime | `pluginrt` | 002 | Plugin sandboxing (scaffold) |
| + | Import/Export | `importexport` | - | CSV import/export |
| + | Sitemap | `sitemap` | - | XML sitemap generation |

## Key Architecture Patterns

### Error Handling
```go
// Always use errx, never errors.New or fmt.Errorf for domain errors
errx.NotFound("product", productID)
errx.Validation("price must be positive")
errx.Business("insufficient stock")
errx.Wrap(err, "failed to create product")
```

### ID Types
All IDs are typed via `kernel` package (69 types in `kernel/proj_ids.go` + 2 in `kernel/common_ids.go`):
```go
type ProductID string
type OrderID string
type TenantID string
// etc.
```

### Money
```go
kernel.Money{Amount: 1999, Currency: "USD"} // $19.99 — always cents
```

### Auth
- **Admin**: `c.Locals("auth").(*kernel.AuthContext)` — IAM JWT middleware
- **Public**: `c.Get("X-Tenant-ID")` header
- **Customer**: `c.Locals("customer_auth").(*auth.CustomerAuthContext)` — separate JWT, routes at `/api/v1/storefront/auth/*` (public) and `/api/v1/storefront/account/*` (protected)

### Multi-tenancy
Every DB query MUST be scoped by `tenant_id`. All repositories take `tenantID` as first parameter.

### Event Bus
In-memory pub/sub with 68 event types defined in `eventbus/eventbus.go`. Typed payloads in `eventbus/payloads.go`. Used by emails, cart recovery, webhooks, analytics.

### Container Wiring
`backend/cmd/container.go` (572 lines):
- `Container` struct holds all 42+ domain containers
- `initModules()` initializes each domain container with `db`, `bus`, and any cross-domain dependencies
- `BuildAgentServices()` maps container services to `agent.Services` struct (31 fields)
- Agent init: `if c.Config.Agent.APIKey != "" { ... }` creates `agentapi.Handler`

### Route Registration
`backend/cmd/server.go` (432 lines):
- `public` group: no auth required (storefront, cart, shipping calc, tax calc, checkout, search, customer auth)
- `protected` group: admin IAM middleware (all CRUD endpoints, agent chat)
- `customerProtected` group: customer JWT middleware (account, wishlist, reviews)

## AI Agent Architecture

### Current State (Just Completed)

The agent uses the [harness](https://github.com/Abraxas-365/harness) library for LLM orchestration:

1. **Domain Tools** (80+ tools): Thin Go structs wrapping domain services. Each implements `agent.Tool` interface: `Name()`, `Description()`, `InputSchema()`, `Execute(ctx, json.RawMessage) (string, error)`

2. **Adapter Layer** (`agent/adapter.go`): `HarnessTool` wraps `agent.Tool` into `harness tools.Tool`. Converts `map[string]any` InputSchema to `json.RawMessage`. Converts `(string, error)` Execute result to `(*tools.Result, error)`.

3. **EventHandler** (`agent/handler.go`): Implements `query.EventHandler` interface with all 10 methods. Streams events (text_delta, tool_start, tool_end, etc.) via `agent.Event` channel. Auto-approves tool execution and cost confirmation.

4. **Chat Endpoint** (`agent/agentapi/handler.go`): `POST /agent/chat` with SSE streaming. Accepts `{session_id, message}` JSON body. Sessions cached in-memory by `tenantID:sessionID`. Creates per-tenant tool sets dynamically (multi-tenant safe). Uses `context.WithoutCancel` so agent finishes even if HTTP drops.

5. **Config**: `ANTHROPIC_API_KEY` and `AGENT_MODEL` env vars. Agent only initializes if API key is set.

6. **harness dependency**: `replace github.com/Abraxas-365/harness => ../../harness` in go.mod (local development).

### Agent Tool Files
| File | Tools | Domains Covered |
|------|-------|----------------|
| `tools.go` | 15 | storefront, products, promos, orders, catalog, themes |
| `tools_commerce.go` | 32 | shipping, tax, payment, search, variants, customer groups, gift cards, cart recovery, currency, i18n, subscriptions |
| `tools_p4.go` | 24 | inventory, reviews, returns, webhooks, audit, loyalty, bundles, dashboard, social auth, notifications |
| `tools_p5.go` | 18 | multistore, bulkops, blog, collections, abtest, recommendations |

## Frontend Admin Pages (42)

```
Dashboard, Products, Orders, OrderDetail, Customers, CustomerDetail,
Catalog, Pages, PageEditor, Media, Promos, Marketplace, PluginView,
Settings, ThemeEditor, Shipping, Tax, Payments, ImportExport,
CustomerGroups, GiftCards, CartRecovery, CurrencyRates, Translations,
Subscriptions, Inventory, Reviews, Returns, Webhooks, AuditLogs,
Loyalty, Bundles, Reporting, SocialAccounts, Notifications,
Multistores, BulkOperations, Blog, Collections, ABTesting,
Recommendations, AgentChat
```

## Migrations (42 files, sequential 001-042)

All in `backend/migrations/`. Append-only — never modify existing files. Next migration number: **043**.

## Git History (latest commits)

```
7843399 feat(agent): wire harness agent into cmd and add multi-tenant support
34b0751 merge: integrate harness agent layer
60bec07 feat(config): add AgentConfig for harness integration
c29beef feat(agent): integrate harness library — adapter, EventHandler, and chat SSE endpoint
43c6262 chore(e2e): update screenshot spec to cover all 45 admin pages
e327e35 feat(agent): wire BuildAgentServices in cmd/container.go
e64ce86 feat(frontend): add P5 admin pages
b983442 feat(agent): add P5 agent tools
a0c2600 feat: merge P5 features - multistore, bulkops, blog, collections, abtest, recommendations
```

## Build Status

- `go build ./...` — PASS
- `go vet ./...` — PASS
- `go test ./...` — PASS (all packages)
- Frontend builds clean with `npm run build`

## What Was Just Completed

### Harness Integration (this session)
1. Added `github.com/Abraxas-365/harness` as dependency with local `replace` directive
2. Created `agent/adapter.go` — `HarnessTool` wraps agent.Tool to harness tools.Tool
3. Updated `agent/handler.go` — `EventHandler` implements full `query.EventHandler` interface (10 methods)
4. Created `agent/agentapi/handler.go` — SSE streaming chat endpoint with per-tenant sessions
5. Added `AgentConfig` to `config.go` (reads `ANTHROPIC_API_KEY`, `AGENT_MODEL`)
6. Wired agent into `cmd/container.go` — creates tools per-tenant dynamically via `BuildAgentServices()`
7. Registered `POST /agent/chat` on protected routes in `cmd/server.go`

## Known Issues / Technical Debt

1. **In-memory session cache** in `agentapi.Handler` has no eviction — needs LRU/TTL for production
2. **Email resolver limitation**: Only `customer.registered` email has customer email in payload. Order/payment/refund events carry CustomerID only — needs `CustomerEmailResolver` to look up email by ID
3. **harness replace directive**: Uses relative path `../../harness` — needs publishing or CI adjustment
4. **Agent system prompt**: Currently empty string `""` — needs a proper system prompt describing store context and available tools
5. **No storefront (public) routes** for many domains — admin-only for now
6. **Plugin runtime** (`pluginrt`) is scaffolded but not fully implemented
7. **Agent chat frontend page** exists (`AgentChat.tsx`) but may need updates to work with the new SSE endpoint

## Planned Next Steps

### Immediate (High Value)
1. **Agent system prompt** — Write a comprehensive system prompt that tells the LLM about the store, available tools, and how to be a helpful store assistant
2. **Frontend chat UI update** — Update `AgentChat.tsx` to connect to the SSE `POST /agent/chat` endpoint, render streaming responses and tool-use events
3. **Add harness built-in tools** — Cherry-pick high-value harness tools (WebSearch, WebFetch) so the agent can research products, competitors, SEO

### Medium Priority
4. **Session eviction** — Add TTL-based session cleanup in `agentapi.Handler`
5. **Agent context enrichment** — Include store name, plan, recent stats in system prompt dynamically
6. **Docker container for agent** (Phase 2) — Isolated execution for heavy tasks (report generation, bulk data processing, content creation with file I/O)
7. **Publish harness** — Remove local `replace` directive, publish to a Go module proxy or private repo

### Future Features (Not Yet Built)
8. **Multi-currency checkout** — Integrate currency conversion into checkout flow
9. **Real payment providers** — Stripe, PayPal integration (currently manual-only)
10. **Real search engine** — Elasticsearch/Meilisearch integration (currently PostgreSQL tsvector)
11. **Storefront SSR** — Server-side rendered public storefront (currently admin-only SPA)
12. **Customer-facing storefront** — Public product pages, cart, checkout UI
13. **CI/CD pipeline** — GitHub Actions, Docker build, deployment
14. **API documentation** — OpenAPI/Swagger spec generation

## Useful Commands

```bash
# Backend
cd backend && go build ./...           # Build
cd backend && go vet ./...             # Lint
cd backend && go test ./...            # Test
cd backend && make migrate             # Run migrations

# Frontend
cd frontend && npm run dev             # Dev server
cd frontend && npm run build           # Production build
cd frontend && npx playwright test     # E2E screenshots

# Infrastructure
docker compose up -d                   # Start postgres + redis
docker compose ps                      # Check services
```

## Agent Development Notes

### Adding a New Domain
1. Create DDD package at `backend/internal/<domain>/` (entity, errors, port, srv, infra, api, container)
2. Add migration `backend/migrations/<next_number>_<domain>.up.sql`
3. Add ID types to `backend/kernel/proj_ids.go`
4. Add event types to `backend/internal/eventbus/eventbus.go` + payloads
5. Add container to `cmd/container.go` struct + `initModules()`
6. Register routes in `cmd/server.go`
7. Add agent tools in `backend/internal/agent/tools_<batch>.go`
8. Add service to `agent.Services` struct and `Setup()` in `setup.go`
9. Map in `BuildAgentServices()` in `container.go`
10. Add TypeScript types, API functions, hooks, route, and admin page in frontend

### Adding Agent Tools
1. Create tool struct in `tools_<batch>.go` implementing `agent.Tool`
2. Register in `Setup()` function in `setup.go`
3. Add service field to `Services` struct if new domain
4. Map service in `BuildAgentServices()` in `container.go`

### Merge Conflict Pattern
When merging multiple agent branches, conflicts always occur in:
- `cmd/container.go` — struct fields, imports, init lines → keep both sides
- `cmd/server.go` — route registrations → keep both sides
- `kernel/proj_ids.go` — ID type definitions → keep both sides
- `eventbus/eventbus.go` — event constants → keep both sides
- `eventbus/payloads.go` — payload structs → keep both sides

Quick resolution: `sed -i '' '/^<<<<<<< HEAD$/d; /^=======$/d; /^>>>>>>> /d' <file>` (keeps both sides)

### Migration Numbering
Sequential 001-042. Next: **043**. When spawning parallel agents, assign non-conflicting numbers. If conflicts occur, renumber with `mv`.
