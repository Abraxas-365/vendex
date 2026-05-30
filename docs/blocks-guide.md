# Blocks Guide

hada-commerce uses a **block-based page builder**. Instead of writing raw HTML, you compose pages from reusable block components. Each block has a well-defined JSON settings schema, which the renderer uses to produce HTML.

---

## Concepts

### Block

A **block** is the smallest content unit. It has a `type` (e.g. `bt-hero`) and a `settings` object that conforms to that block type's schema.

### Section

A **section** is a full-width page region. It contains one or more blocks and its own settings (background color, padding, etc.). Sections are rendered top-to-bottom.

### Page

A **page** is identified by a `slug` (URL path) and contains an ordered list of sections. Pages go through a lifecycle: **draft → published → archived**.

---

## Page Lifecycle

```
                ┌─────────┐
                │  draft  │◄──────────────────────┐
                └────┬────┘                       │
                     │ submit for review           │
                     ▼                             │
            ┌────────────────┐                    │
            │ pending_review │                    │ edit
            └───────┬────────┘                    │
                    │ approve                     │
                    ▼                             │
             ┌───────────┐                       │
             │ published │──────── archive ──────►┤
             └───────────┘               ┌───────┴──────┐
                                         │   archived   │
                                         └──────────────┘
```

| Status | Description |
|--------|-------------|
| `draft` | Work in progress. Not visible to shoppers. |
| `pending_review` | Submitted for approval. Edits locked. |
| `published` | Live and publicly accessible. |
| `archived` | Read-only. Hidden from the public. |

> **Note:** Only `published` pages are returned by the public storefront API.

---

## The 10 Built-in Block Types

### 1. `bt-hero`

Full-width hero section with a heading, subheading, and call-to-action button.

**Settings schema:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `heading` | string | `"Welcome"` | Main headline |
| `subheading` | string | `""` | Supporting text below the headline |
| `background_image` | string | `""` | URL of background image |
| `background_color` | string | `"#4F46E5"` | Fallback background color |
| `cta_text` | string | `""` | Button label |
| `cta_link` | string | `""` | Button URL |
| `alignment` | `"left"` \| `"center"` \| `"right"` | `"center"` | Text alignment |

**Example:**

```json
{
  "type": "bt-hero",
  "settings": {
    "heading": "Summer Collection",
    "subheading": "Up to 40% off selected styles",
    "background_image": "https://cdn.example.com/hero.jpg",
    "cta_text": "Shop Now",
    "cta_link": "/collections/summer",
    "alignment": "center"
  }
}
```

---

### 2. `bt-rich-text`

Rich HTML content block for body copy, FAQs, policy text, etc.

**Settings schema:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `content` | string (HTML) | `""` | Raw HTML content |

**Example:**

```json
{
  "type": "bt-rich-text",
  "settings": {
    "content": "<h2>Our Story</h2><p>Founded in 2020, we believe in sustainable fashion...</p>"
  }
}
```

---

### 3. `bt-product-grid`

Displays a grid of products from a collection or category.

**Settings schema:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `title` | string | `"Featured Products"` | Section heading |
| `collection_id` | string | `""` | Filter by collection ID (optional) |
| `columns` | integer (1–6) | `4` | Number of grid columns |
| `limit` | integer (1–24) | `8` | Maximum products to show |
| `show_price` | boolean | `true` | Whether to display the price |

**Example:**

```json
{
  "type": "bt-product-grid",
  "settings": {
    "title": "Best Sellers",
    "collection_id": "col_abc123",
    "columns": 4,
    "limit": 8,
    "show_price": true
  }
}
```

---

### 4. `bt-image-banner`

Full-width image with optional link overlay.

**Settings schema:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `image_url` | string | `""` | Image source URL |
| `alt_text` | string | `""` | Accessibility alt text |
| `link` | string | `""` | Optional click-through URL |
| `height` | string (CSS) | `"400px"` | Banner height |

**Example:**

```json
{
  "type": "bt-image-banner",
  "settings": {
    "image_url": "https://cdn.example.com/banner.jpg",
    "alt_text": "New arrivals banner",
    "link": "/collections/new",
    "height": "500px"
  }
}
```

---

### 5. `bt-newsletter`

Email newsletter sign-up form.

**Settings schema:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `heading` | string | `"Subscribe to our newsletter"` | Form heading |
| `description` | string | `""` | Supporting text |
| `button_text` | string | `"Subscribe"` | Submit button label |
| `placeholder` | string | `"Enter your email"` | Input placeholder |

**Example:**

```json
{
  "type": "bt-newsletter",
  "settings": {
    "heading": "Stay in the loop",
    "description": "Get exclusive deals delivered to your inbox.",
    "button_text": "Sign me up",
    "placeholder": "your@email.com"
  }
}
```

---

### 6. `bt-category-grid`

Displays a grid of product categories as card links.

**Settings schema:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `title` | string | `"Shop by Category"` | Section heading |
| `columns` | integer (1–6) | `3` | Number of grid columns |
| `show_description` | boolean | `false` | Show category description text |

**Example:**

```json
{
  "type": "bt-category-grid",
  "settings": {
    "title": "Browse Categories",
    "columns": 4,
    "show_description": true
  }
}
```

---

### 7. `bt-cta`

Full-width call-to-action banner with heading, description, and a button.

**Settings schema:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `heading` | string | `"Ready to get started?"` | Primary heading |
| `description` | string | `""` | Supporting description |
| `button_text` | string | `"Get Started"` | CTA button label |
| `button_link` | string | `""` | CTA button URL |
| `background_color` | string | `"#111827"` | Block background color |

**Example:**

```json
{
  "type": "bt-cta",
  "settings": {
    "heading": "Join 10,000+ happy customers",
    "description": "Start your free trial today. No credit card required.",
    "button_text": "Start Free Trial",
    "button_link": "/signup",
    "background_color": "#1D4ED8"
  }
}
```

---

### 8. `bt-spacer`

Empty vertical space for layout control.

**Settings schema:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `height` | string (CSS) | `"48px"` | Height of the spacer |

**Example:**

```json
{
  "type": "bt-spacer",
  "settings": {
    "height": "80px"
  }
}
```

---

### 9. `bt-divider`

Horizontal rule / visual separator.

**Settings schema:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `color` | string | `"#E5E7EB"` | Line color |
| `thickness` | string (CSS) | `"1px"` | Line thickness |
| `width` | string (CSS) | `"100%"` | Line width |

**Example:**

```json
{
  "type": "bt-divider",
  "settings": {
    "color": "#CBD5E1",
    "thickness": "2px",
    "width": "80%"
  }
}
```

---

### 10. `bt-testimonials`

Customer testimonial carousel or grid.

**Settings schema:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `title` | string | `"What our customers say"` | Section heading |
| `items` | array | `[]` | List of testimonial objects |

Each item in `items`:

| Field | Type | Description |
|-------|------|-------------|
| `quote` | string | The testimonial text |
| `author` | string | Customer name |
| `role` | string | Customer title or location |
| `avatar` | string | URL to avatar image |

**Example:**

```json
{
  "type": "bt-testimonials",
  "settings": {
    "title": "Loved by thousands",
    "items": [
      {
        "quote": "Best purchase I've made all year. Quality is outstanding.",
        "author": "Sarah M.",
        "role": "Verified Buyer",
        "avatar": "https://cdn.example.com/avatars/sarah.jpg"
      },
      {
        "quote": "Fast shipping and exactly as described.",
        "author": "James K.",
        "role": "Repeat Customer",
        "avatar": ""
      }
    ]
  }
}
```

---

## How Sections Compose into Pages

A page is a JSON structure like this:

```json
{
  "slug": "home",
  "title": "Home",
  "content_type": "blocks",
  "sections": [
    {
      "id": "section-1",
      "type": "hero",
      "settings": {},
      "blocks": [
        {
          "id": "block-1",
          "type": "bt-hero",
          "settings": {
            "heading": "Welcome to our store",
            "cta_text": "Shop Now",
            "cta_link": "/products"
          }
        }
      ]
    },
    {
      "id": "section-2",
      "type": "products",
      "settings": {},
      "blocks": [
        {
          "id": "block-2",
          "type": "bt-product-grid",
          "settings": {
            "title": "New Arrivals",
            "limit": 8
          }
        }
      ]
    }
  ]
}
```

### Creating a Page via API

```bash
curl -X POST http://localhost:8080/api/v1/storefront/pages \
  -H "X-Tenant-ID: my-store" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "slug": "home",
    "title": "Home",
    "content_type": "blocks",
    "sections": [
      {
        "type": "hero",
        "settings": {},
        "blocks": [
          {
            "type": "bt-hero",
            "settings": {
              "heading": "Welcome!",
              "cta_text": "Shop Now",
              "cta_link": "/products"
            }
          }
        ]
      }
    ]
  }'
```

### Publishing a Page

After creation, pages start as `draft`. To make a page live:

```bash
curl -X POST http://localhost:8080/api/v1/storefront/pages/{id}/publish \
  -H "X-Tenant-ID: my-store" \
  -H "Authorization: Bearer <token>"
```

### Fetching a Rendered Page

Any published page is available at the public storefront endpoint. Use the `Accept` header to control the response format:

```bash
# Get rendered HTML (for SSR integration)
curl http://localhost:8080/api/v1/storefront/pages/by-slug/home \
  -H "X-Tenant-ID: my-store" \
  -H "Accept: text/html"

# Get raw JSON (for client-side rendering)
curl http://localhost:8080/api/v1/storefront/pages/by-slug/home \
  -H "X-Tenant-ID: my-store" \
  -H "Accept: application/json"
```

---

## Creating Custom Block Types

You can register custom block types via the API. Custom blocks appear alongside built-in blocks in the page editor.

```bash
curl -X POST http://localhost:8080/api/v1/storefront/block-types \
  -H "X-Tenant-ID: my-store" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "bt-countdown",
    "name": "Countdown Timer",
    "category": "content",
    "schema": {
      "target_date": {
        "type": "string",
        "label": "End date (ISO 8601)",
        "default": ""
      },
      "heading": {
        "type": "string",
        "label": "Heading",
        "default": "Sale ends in:"
      }
    },
    "defaults": {
      "target_date": "",
      "heading": "Sale ends in:"
    }
  }'
```

The `schema` object describes each setting field for the admin UI form. Your custom block template is rendered using the same block rendering pipeline.

---

## Version History

Every time you save a page, a `PageVersion` is created. You can retrieve the version history for a page to audit or restore a previous state:

```bash
curl http://localhost:8080/api/v1/storefront/pages/{id}/versions \
  -H "X-Tenant-ID: my-store" \
  -H "Authorization: Bearer <token>"
```
