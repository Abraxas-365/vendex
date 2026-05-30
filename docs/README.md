# hada-commerce

A multi-tenant, headless-ready e-commerce platform built with Go and React. Ships with a block-based page builder, a design token theme system, a plugin framework, and a fully documented public API — so you can use the built-in admin panel _or_ build a completely custom storefront.

---

## Key Features

| Feature | Description |
|---------|-------------|
| **Block-based page builder** | Compose pages from 10 built-in block types. Each block has a typed JSON settings schema. |
| **Design token themes** | 24 CSS custom properties (colors, typography, spacing, borders, shadows) stored per tenant, applied at render time. |
| **Plugin system** | Install plugins that receive webhook events and inject frontend widgets into the admin UI. |
| **Event-driven core** | An in-process event bus emits 21 domain events (orders, products, customers, pages, themes, …). |
| **Headless API** | Public REST endpoints (no auth required) for building custom storefronts with Next.js, Remix, Astro, etc. |
| **SSR renderer** | Server-side HTML rendering of block pages with theme CSS variables injected — no client JS required. |
| **Multi-tenancy** | Every resource is fully scoped to a tenant. A single deployment serves unlimited stores. |
| **Admin panel** | React SPA for managing products, orders, customers, catalog, pages, themes, plugins, and settings. |

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         Clients                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────────────────┐ │
│  │  Admin Panel │  │  Custom      │  │  Plugin Webhooks      │ │
│  │  (React SPA) │  │  Storefront  │  │  (external services)  │ │
│  └──────┬───────┘  └──────┬───────┘  └──────────┬────────────┘ │
└─────────┼─────────────────┼──────────────────────┼─────────────┘
          │  Bearer token   │  X-Tenant-ID          │  POST events
          ▼                 ▼                        ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Go Backend (net/http)                     │
│                                                                 │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────────┐  │
│  │ product  │  │  order   │  │ customer │  │   catalog     │  │
│  └──────────┘  └──────────┘  └──────────┘  └───────────────┘  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────────┐  │
│  │storefront│  │  theme   │  │  plugin  │  │   settings    │  │
│  └──────────┘  └──────────┘  └──────────┘  └───────────────┘  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────────┐  │
│  │  promo   │  │  media   │  │analytics │  │  marketplace  │  │
│  └──────────┘  └──────────┘  └──────────┘  └───────────────┘  │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    Event Bus                            │   │
│  │  order.*  product.*  customer.*  page.*  theme.*  …    │   │
│  └─────────────────────────────────────────────────────────┘   │
└────────────────────┬─────────────────────────────┬─────────────┘
                     │                             │
          ┌──────────▼──────────┐      ┌───────────▼──────────┐
          │   PostgreSQL 16     │      │      Redis 7          │
          │   (primary store)   │      │  (jobs, cache)        │
          └─────────────────────┘      └──────────────────────┘
```

### Bounded Contexts

| Context | Responsibility |
|---------|---------------|
| `product` | Product catalog: name, price, SKU, stock, images, tags, status |
| `order` | Order lifecycle: pending → confirmed → shipped → delivered |
| `customer` | Customer profiles and addresses |
| `catalog` | Categories (hierarchical) and collections |
| `storefront` | Block-based pages, sections, block types, SSR renderer |
| `theme` | Design token presets, CSS custom property generation |
| `plugin` | Plugin registry, versioning, per-tenant installation, webhook delivery |
| `settings` | Per-tenant store configuration (name, currency, checkout rules) |
| `promo` | Discount codes (percentage, fixed, free-shipping) |
| `media` | File uploads and asset management |
| `analytics` | Event tracking |
| `marketplace` | Multi-vendor support (vendors, vendor products, vendor orders) |

---

## Documentation

| Guide | Audience |
|-------|---------|
| [Getting Started](./getting-started.md) | Developers setting up a local environment |
| [Blocks Guide](./blocks-guide.md) | Merchants and developers building pages |
| [Themes Guide](./themes-guide.md) | Designers customizing store appearance |
| [Plugins Guide](./plugins-guide.md) | Developers building integrations |
| [Headless API](./headless-api.md) | Frontend developers building custom storefronts |
| [Architecture](./architecture.md) | Developers contributing to the platform |

---

## Quick Links

- **Admin panel**: `http://localhost:3000` (after running locally)
- **API base URL**: `http://localhost:8080/api/v1`
- **Module**: `github.com/Abraxas-365/hada-commerce`
- **Go version**: 1.24+
- **PostgreSQL**: 16
- **Redis**: 7
