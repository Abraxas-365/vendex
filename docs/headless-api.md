# Headless API

hada-commerce exposes a public REST API that lets you build a fully custom storefront using any frontend framework — Next.js, Remix, Astro, plain HTML, or a native mobile app. The backend handles products, catalog, orders, pages, themes, and settings; your frontend handles rendering.

---

## Authentication

### Public endpoints (no auth required)

These endpoints are safe to call from a browser or edge function. They require only the `X-Tenant-ID` header to scope the response to the correct store.

```
X-Tenant-ID: my-store
```

### Admin/protected endpoints

Mutations (create, update, delete) require a Bearer token in addition to the tenant header:

```
X-Tenant-ID: my-store
Authorization: Bearer <token>
```

Obtain a token by authenticating through the IAM endpoints.

---

## Public Endpoint Reference

### Products

#### List products

```
GET /api/v1/products
```

Query parameters:

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | `1` | Page number |
| `page_size` | integer | `20` | Items per page (max 100) |

```bash
curl http://localhost:8080/api/v1/products \
  -H "X-Tenant-ID: my-store"
```

Response:

```json
{
  "items": [
    {
      "id": "prd_abc123",
      "tenant_id": "my-store",
      "name": "Classic T-Shirt",
      "description": "A comfortable everyday tee.",
      "price": { "amount": 2999, "currency": "USD" },
      "sku": "TSHIRT-001",
      "images": ["https://cdn.example.com/tshirt.jpg"],
      "category_id": "cat_xyz",
      "tags": ["clothing", "basics"],
      "status": "active",
      "stock": 87,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-10T00:00:00Z"
    }
  ],
  "total": 142,
  "page": 1,
  "page_size": 20,
  "total_pages": 8
}
```

> **Note:** `price.amount` is always in the smallest currency unit (cents for USD). `2999` = $29.99.

#### Get a single product

```
GET /api/v1/products/{id}
```

```bash
curl http://localhost:8080/api/v1/products/prd_abc123 \
  -H "X-Tenant-ID: my-store"
```

---

### Categories

#### List categories

```
GET /api/v1/categories
```

```bash
curl http://localhost:8080/api/v1/categories \
  -H "X-Tenant-ID: my-store"
```

Response:

```json
{
  "items": [
    {
      "id": "cat_xyz",
      "tenant_id": "my-store",
      "name": "Clothing",
      "slug": "clothing",
      "parent_id": null,
      "description": "All clothing items",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 12,
  "page": 1,
  "page_size": 20,
  "total_pages": 1
}
```

Categories support hierarchies via `parent_id`. Root categories have `parent_id: null`.

---

### Collections

#### List collections

```
GET /api/v1/collections
```

```bash
curl http://localhost:8080/api/v1/collections \
  -H "X-Tenant-ID: my-store"
```

#### Get a single collection

```
GET /api/v1/collections/{id}
```

```bash
curl http://localhost:8080/api/v1/collections/col_featured \
  -H "X-Tenant-ID: my-store"
```

Response:

```json
{
  "id": "col_featured",
  "tenant_id": "my-store",
  "name": "Featured",
  "slug": "featured",
  "description": "Hand-picked bestsellers",
  "type": "manual",
  "products": ["prd_abc123", "prd_def456"],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-15T00:00:00Z"
}
```

---

### Store Settings

#### Get store settings

```
GET /api/v1/settings
```

```bash
curl http://localhost:8080/api/v1/settings \
  -H "X-Tenant-ID: my-store"
```

Response:

```json
{
  "tenant_id": "my-store",
  "store_name": "My Store",
  "store_email": "hello@mystore.com",
  "store_phone": "+1 555-0100",
  "currency": "USD",
  "timezone": "America/New_York",
  "address": {
    "street": "123 Main St",
    "city": "New York",
    "state": "NY",
    "country": "US",
    "zip": "10001"
  },
  "logo_url": "https://cdn.mystore.com/logo.png",
  "favicon_url": "https://cdn.mystore.com/favicon.ico",
  "social_links": {
    "instagram": "https://instagram.com/mystore",
    "twitter": "https://twitter.com/mystore",
    "facebook": "https://facebook.com/mystore"
  },
  "checkout_config": {
    "guest_checkout": true,
    "require_phone": false
  },
  "updated_at": "2024-01-15T10:00:00Z"
}
```

---

### Storefront Pages

#### Get a page by slug (HTML or JSON)

```
GET /api/v1/storefront/pages/by-slug/{slug}
```

Use the `Accept` header to control the response format:

| Accept header | Response |
|--------------|---------|
| `text/html` | Server-rendered HTML page with theme CSS variables |
| `application/json` | Raw page JSON with sections and blocks |

```bash
# Get rendered HTML (for SSR proxy or iframe embed)
curl http://localhost:8080/api/v1/storefront/pages/by-slug/home \
  -H "X-Tenant-ID: my-store" \
  -H "Accept: text/html"

# Get raw JSON (for client-side rendering)
curl http://localhost:8080/api/v1/storefront/pages/by-slug/home \
  -H "X-Tenant-ID: my-store" \
  -H "Accept: application/json"
```

**JSON response:**

```json
{
  "id": "pg_home",
  "tenant_id": "my-store",
  "slug": "home",
  "title": "Home",
  "content_type": "blocks",
  "status": "published",
  "sections": [
    {
      "id": "sec_1",
      "type": "hero",
      "settings": {},
      "blocks": [
        {
          "id": "blk_1",
          "type": "bt-hero",
          "settings": {
            "heading": "Welcome to My Store",
            "cta_text": "Shop Now",
            "cta_link": "/products"
          }
        }
      ]
    }
  ],
  "meta": {
    "description": "Official store for My Store",
    "og_title": "My Store",
    "og_image": "https://cdn.mystore.com/og.jpg",
    "keywords": ["fashion", "clothing"]
  },
  "published_at": "2024-01-10T00:00:00Z"
}
```

**HTML response** is a complete HTML document ready to embed or proxy:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Home | My Store</title>
  <meta property="og:title" content="My Store">
  <meta property="og:image" content="https://cdn.mystore.com/og.jpg">
  <style>
    :root {
      --color-primary: #6366f1;
      --color-background: #ffffff;
      /* ... all theme tokens ... */
    }
    /* base reset + utility styles */
  </style>
</head>
<body>
  <div class="section-wrapper" data-section-id="sec_1" data-section-type="hero">
    <!-- rendered block HTML -->
  </div>
</body>
</html>
```

---

### Active Theme

#### Get the active theme tokens

```
GET /api/v1/themes/active
```

```bash
curl http://localhost:8080/api/v1/themes/active \
  -H "X-Tenant-ID: my-store"
```

Response:

```json
{
  "id": "thm_abc",
  "name": "Ocean Blue",
  "is_active": true,
  "tokens": {
    "colors": {
      "primary": "#0EA5E9",
      "secondary": "#38BDF8",
      "background": "#F0F9FF"
    },
    "typography": {
      "font_heading": "Poppins, sans-serif",
      "font_body": "Inter, sans-serif",
      "base_size": "16px",
      "scale_ratio": 1.25
    },
    "spacing": {
      "unit": "4px",
      "section_padding": "64px"
    },
    "borders": {
      "radius_sm": "4px",
      "radius_md": "8px",
      "radius_lg": "16px",
      "radius_full": "9999px"
    },
    "shadows": {
      "sm": "0 1px 2px 0 rgb(0 0 0 / 0.05)",
      "md": "0 4px 6px -1px rgb(0 0 0 / 0.1)",
      "lg": "0 10px 15px -3px rgb(0 0 0 / 0.1)"
    }
  }
}
```

Use these tokens in your custom frontend to match the merchant's design system.

---

### Plugin JS Manifest

#### Get frontend JS scripts for active plugins

```
GET /api/v1/plugins/js-manifest
```

```bash
curl http://localhost:8080/api/v1/plugins/js-manifest \
  -H "X-Tenant-ID: my-store"
```

Useful if you want to load plugin widgets in your custom storefront.

---

## Building a Custom Storefront

### Next.js Example

Install the fetch utility:

```bash
npm install
```

Create a shared API client (`lib/api.ts`):

```typescript
const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080';

async function hadaFetch<T>(
  path: string,
  tenantId: string,
  init?: RequestInit
): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers: {
      'X-Tenant-ID': tenantId,
      'Content-Type': 'application/json',
      ...init?.headers,
    },
  });

  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.message ?? 'API error');
  }

  return res.json();
}

export const api = {
  products: {
    list: (tenantId: string, page = 1, pageSize = 20) =>
      hadaFetch(`/api/v1/products?page=${page}&page_size=${pageSize}`, tenantId),
    get: (tenantId: string, id: string) =>
      hadaFetch(`/api/v1/products/${id}`, tenantId),
  },
  categories: {
    list: (tenantId: string) =>
      hadaFetch('/api/v1/categories', tenantId),
  },
  collections: {
    list: (tenantId: string) =>
      hadaFetch('/api/v1/collections', tenantId),
    get: (tenantId: string, id: string) =>
      hadaFetch(`/api/v1/collections/${id}`, tenantId),
  },
  pages: {
    bySlug: (tenantId: string, slug: string) =>
      hadaFetch(`/api/v1/storefront/pages/by-slug/${slug}`, tenantId),
  },
  theme: {
    active: (tenantId: string) =>
      hadaFetch('/api/v1/themes/active', tenantId),
  },
  settings: {
    get: (tenantId: string) =>
      hadaFetch('/api/v1/settings', tenantId),
  },
};
```

### Product listing page (`app/products/page.tsx`):

```typescript
import { api } from '@/lib/api';

const TENANT_ID = process.env.HADA_TENANT_ID!;

export default async function ProductsPage() {
  const { items: products, total } = await api.products.list(TENANT_ID);
  const settings = await api.settings.get(TENANT_ID);

  return (
    <main>
      <h1>Products — {settings.store_name}</h1>
      <p>{total} products available</p>
      <div className="product-grid">
        {products.map((product) => (
          <div key={product.id} className="product-card">
            <img src={product.images[0]} alt={product.name} />
            <h2>{product.name}</h2>
            <p>${(product.price.amount / 100).toFixed(2)}</p>
          </div>
        ))}
      </div>
    </main>
  );
}
```

### Apply theme tokens (`app/layout.tsx`):

```typescript
import { api } from '@/lib/api';

const TENANT_ID = process.env.HADA_TENANT_ID!;

export default async function RootLayout({ children }) {
  const theme = await api.theme.active(TENANT_ID);
  const { colors, typography, spacing, borders, shadows } = theme.tokens;

  const cssVars = `
    :root {
      --color-primary: ${colors.primary};
      --color-background: ${colors.background};
      --color-text: ${colors.text};
      --color-surface: ${colors.surface};
      --font-heading: ${typography.font_heading};
      --font-body: ${typography.font_body};
      --font-base-size: ${typography.base_size};
      --spacing-section-padding: ${spacing.section_padding};
      --border-radius-md: ${borders.radius_md};
      --shadow-md: ${shadows.md};
    }
  `;

  return (
    <html lang="en">
      <head>
        <style dangerouslySetInnerHTML={{ __html: cssVars }} />
      </head>
      <body style={{ fontFamily: 'var(--font-body)' }}>
        {children}
      </body>
    </html>
  );
}
```

### Block page renderer (`app/[slug]/page.tsx`):

```typescript
import { api } from '@/lib/api';
import { notFound } from 'next/navigation';

const TENANT_ID = process.env.HADA_TENANT_ID!;

export default async function StorefrontPage({ params }) {
  const page = await api.pages.bySlug(TENANT_ID, params.slug).catch(() => null);

  if (!page || page.status !== 'published') {
    return notFound();
  }

  return (
    <main>
      {page.sections.map((section) => (
        <section key={section.id} data-section-type={section.type}>
          {section.blocks.map((block) => (
            <Block key={block.id} block={block} />
          ))}
        </section>
      ))}
    </main>
  );
}

function Block({ block }) {
  switch (block.type) {
    case 'bt-hero':
      return (
        <div style={{
          backgroundColor: block.settings.background_color ?? 'var(--color-primary)',
          textAlign: block.settings.alignment ?? 'center',
          padding: 'var(--spacing-section-padding)',
        }}>
          <h1 style={{ fontFamily: 'var(--font-heading)' }}>
            {block.settings.heading}
          </h1>
          {block.settings.subheading && (
            <p>{block.settings.subheading}</p>
          )}
          {block.settings.cta_link && (
            <a href={block.settings.cta_link}>{block.settings.cta_text}</a>
          )}
        </div>
      );

    case 'bt-rich-text':
      return (
        <div
          dangerouslySetInnerHTML={{ __html: block.settings.content }}
          style={{ maxWidth: '65ch', margin: '0 auto' }}
        />
      );

    default:
      return null; // unknown block type — skip
  }
}
```

---

## Remix Example

Create a loader in `app/routes/products.tsx`:

```typescript
import { json, LoaderFunction } from '@remix-run/node';
import { useLoaderData } from '@remix-run/react';

const TENANT_ID = process.env.HADA_TENANT_ID!;
const API_BASE  = process.env.HADA_API_URL!;

export const loader: LoaderFunction = async () => {
  const res = await fetch(`${API_BASE}/api/v1/products?page_size=24`, {
    headers: { 'X-Tenant-ID': TENANT_ID },
  });
  const data = await res.json();
  return json(data);
};

export default function ProductsRoute() {
  const { items } = useLoaderData<typeof loader>();

  return (
    <ul>
      {items.map(p => (
        <li key={p.id}>
          {p.name} — ${(p.price.amount / 100).toFixed(2)}
        </li>
      ))}
    </ul>
  );
}
```

---

## Astro Example

In `src/pages/products.astro`:

```astro
---
const TENANT_ID = import.meta.env.HADA_TENANT_ID;
const API_BASE  = import.meta.env.HADA_API_URL;

const res  = await fetch(`${API_BASE}/api/v1/products`, {
  headers: { 'X-Tenant-ID': TENANT_ID },
});
const { items } = await res.json();
---

<html lang="en">
<head><title>Products</title></head>
<body>
  <ul>
    {items.map(p => (
      <li>
        <a href={`/products/${p.id}`}>{p.name}</a>
        — ${(p.price.amount / 100).toFixed(2)}
      </li>
    ))}
  </ul>
</body>
</html>
```

---

## Money Formatting

All `price.amount` values are integers in the **smallest currency unit**:

| Currency | Unit | Example |
|---------|------|---------|
| USD | cents | `2999` = $29.99 |
| EUR | cents | `2999` = €29.99 |
| GBP | pence | `2999` = £29.99 |
| JPY | yen (no subunit) | `2999` = ¥2,999 |

Use the `Intl.NumberFormat` API to format for display:

```typescript
function formatMoney(amount: number, currency: string, locale = 'en-US') {
  return new Intl.NumberFormat(locale, {
    style: 'currency',
    currency,
    minimumFractionDigits: currency === 'JPY' ? 0 : 2,
  }).format(amount / (currency === 'JPY' ? 1 : 100));
}

formatMoney(2999, 'USD'); // "$29.99"
formatMoney(2999, 'JPY'); // "¥2,999"
```

---

## Pagination

All list endpoints use consistent pagination:

```bash
curl "http://localhost:8080/api/v1/products?page=2&page_size=12" \
  -H "X-Tenant-ID: my-store"
```

Response always includes:

```json
{
  "items": [...],
  "total": 142,
  "page": 2,
  "page_size": 12,
  "total_pages": 12
}
```

| Parameter | Default | Max |
|-----------|---------|-----|
| `page` | `1` | — |
| `page_size` | `20` | `100` |

---

## Error Responses

All errors follow the same JSON structure:

```json
{
  "error": "PRODUCT_NOT_FOUND",
  "message": "product not found"
}
```

| HTTP Status | When |
|------------|------|
| `400` | Bad request / invalid parameters |
| `404` | Resource not found |
| `409` | Conflict (duplicate slug, SKU, etc.) |
| `422` | Validation failure |
| `500` | Internal server error |

---

## Environment Variables for Custom Frontends

```dotenv
HADA_TENANT_ID=my-store
HADA_API_URL=http://localhost:8080
```

In production, point `HADA_API_URL` to your deployed backend URL.
