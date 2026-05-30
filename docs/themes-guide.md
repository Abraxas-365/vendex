# Themes Guide

A **theme** in hada-commerce is a named collection of **design tokens** — values for colors, typography, spacing, borders, and shadows. At render time, the active theme is converted to CSS custom properties and injected into every storefront page.

You can have multiple themes per tenant (e.g. "Light", "Dark", "Holiday") but only one can be active at a time.

---

## How Themes Work

```
Theme (stored in DB as JSONB)
        │
        ▼
  TokensToCSS()
        │
        ▼
  :root { --color-primary: #6366f1; ... }
        │
        ▼
  Injected into <head> of every rendered page
        │
        ▼
  Block templates reference CSS vars (var(--color-primary))
```

The renderer resolves the active theme for the current tenant before rendering any page. If no theme is active, platform defaults are used.

---

## Design Token Reference

### Colors — 11 tokens

| Token | CSS Variable | Default | Description |
|-------|-------------|---------|-------------|
| `primary` | `--color-primary` | `#6366f1` | Main brand/action color |
| `secondary` | `--color-secondary` | `#8b5cf6` | Accent color |
| `background` | `--color-background` | `#ffffff` | Page background |
| `surface` | `--color-surface` | `#f9fafb` | Card / elevated surfaces |
| `text` | `--color-text` | `#111827` | Body text |
| `text_muted` | `--color-text-muted` | `#6b7280` | Secondary / de-emphasized text |
| `border` | `--color-border` | `#e5e7eb` | Borders and dividers |
| `error` | `--color-error` | `#ef4444` | Error states |
| `success` | `--color-success` | `#22c55e` | Success states |
| `warning` | `--color-warning` | `#f59e0b` | Warning states |
| `info` | `--color-info` | `#3b82f6` | Informational states |

### Typography — 4 tokens

| Token | CSS Variable | Default | Description |
|-------|-------------|---------|-------------|
| `font_heading` | `--font-heading` | `Inter, sans-serif` | Heading font stack |
| `font_body` | `--font-body` | `Inter, sans-serif` | Body font stack |
| `base_size` | `--font-base-size` | `16px` | Root font size |
| `scale_ratio` | `--font-scale-ratio` | `1.25` | Modular scale multiplier |

### Spacing — 2 tokens

| Token | CSS Variable | Default | Description |
|-------|-------------|---------|-------------|
| `unit` | `--spacing-unit` | `4px` | Base spacing unit |
| `section_padding` | `--spacing-section-padding` | `64px` | Vertical padding for sections |

### Borders — 4 tokens

| Token | CSS Variable | Default | Description |
|-------|-------------|---------|-------------|
| `radius_sm` | `--border-radius-sm` | `4px` | Small corners (inputs, badges) |
| `radius_md` | `--border-radius-md` | `8px` | Medium corners (cards, buttons) |
| `radius_lg` | `--border-radius-lg` | `16px` | Large corners (modals, panels) |
| `radius_full` | `--border-radius-full` | `9999px` | Fully rounded (pills, avatars) |

### Shadows — 3 tokens

| Token | CSS Variable | Default | Description |
|-------|-------------|---------|-------------|
| `sm` | `--shadow-sm` | `0 1px 2px 0 rgb(0 0 0 / 0.05)` | Subtle elevation |
| `md` | `--shadow-md` | `0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)` | Medium elevation |
| `lg` | `--shadow-lg` | `0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)` | High elevation |

---

## Generated CSS Output

When the active theme is applied, the renderer injects a `<style>` block into every page:

```css
:root {
  --color-primary: #6366f1;
  --color-secondary: #8b5cf6;
  --color-background: #ffffff;
  --color-surface: #f9fafb;
  --color-text: #111827;
  --color-text-muted: #6b7280;
  --color-border: #e5e7eb;
  --color-error: #ef4444;
  --color-success: #22c55e;
  --color-warning: #f59e0b;
  --color-info: #3b82f6;

  --font-heading: Inter, sans-serif;
  --font-body: Inter, sans-serif;
  --font-base-size: 16px;
  --font-scale-ratio: 1.25;

  --spacing-unit: 4px;
  --spacing-section-padding: 64px;

  --border-radius-sm: 4px;
  --border-radius-md: 8px;
  --border-radius-lg: 16px;
  --border-radius-full: 9999px;

  --shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.05);
  --shadow-md: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);
  --shadow-lg: 0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1);
}
```

Block templates use these variables directly in inline styles or CSS classes, so changing the active theme instantly reskins the entire storefront.

---

## Managing Themes via API

### Create a Theme

```bash
curl -X POST http://localhost:8080/api/v1/themes \
  -H "X-Tenant-ID: my-store" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ocean Blue",
    "tokens": {
      "colors": {
        "primary": "#0EA5E9",
        "secondary": "#38BDF8",
        "background": "#F0F9FF",
        "surface": "#E0F2FE",
        "text": "#0C4A6E",
        "text_muted": "#0369A1",
        "border": "#BAE6FD"
      },
      "typography": {
        "font_heading": "Poppins, sans-serif",
        "font_body": "Inter, sans-serif",
        "base_size": "16px",
        "scale_ratio": 1.25
      }
    }
  }'
```

You can provide only the tokens you want to override — missing tokens fall back to the platform defaults.

**Response:**

```json
{
  "id": "thm_abc123",
  "tenant_id": "my-store",
  "name": "Ocean Blue",
  "is_active": false,
  "tokens": { "..." },
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

### List Themes

```bash
curl http://localhost:8080/api/v1/themes \
  -H "X-Tenant-ID: my-store" \
  -H "Authorization: Bearer <token>"
```

### Get a Specific Theme

```bash
curl http://localhost:8080/api/v1/themes/{id} \
  -H "X-Tenant-ID: my-store" \
  -H "Authorization: Bearer <token>"
```

### Activate a Theme

Only one theme can be active at a time. Activating a theme automatically deactivates the previous one.

```bash
curl -X POST http://localhost:8080/api/v1/themes/{id}/activate \
  -H "X-Tenant-ID: my-store" \
  -H "Authorization: Bearer <token>"
```

This emits a `theme.activated` event on the event bus, which is delivered to any installed plugins that subscribe to it.

### Update Theme Tokens

You can update individual token groups without replacing the entire theme:

```bash
curl -X PUT http://localhost:8080/api/v1/themes/{id} \
  -H "X-Tenant-ID: my-store" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "tokens": {
      "colors": {
        "primary": "#DC2626"
      }
    }
  }'
```

### Delete a Theme

You cannot delete the currently active theme. Deactivate it first.

```bash
curl -X DELETE http://localhost:8080/api/v1/themes/{id} \
  -H "X-Tenant-ID: my-store" \
  -H "Authorization: Bearer <token>"
```

### Get the Active Theme (Public)

This endpoint requires no authentication — it is used by custom storefronts:

```bash
curl http://localhost:8080/api/v1/themes/active \
  -H "X-Tenant-ID: my-store"
```

---

## Theme Presets

Below are example token sets to get started quickly.

### Minimal Light

```json
{
  "colors": {
    "primary": "#111827",
    "secondary": "#374151",
    "background": "#ffffff",
    "surface": "#f9fafb",
    "text": "#111827",
    "text_muted": "#6b7280",
    "border": "#e5e7eb"
  },
  "typography": {
    "font_heading": "Georgia, serif",
    "font_body": "system-ui, sans-serif"
  },
  "borders": {
    "radius_sm": "2px",
    "radius_md": "4px",
    "radius_lg": "8px",
    "radius_full": "9999px"
  }
}
```

### Bold Dark

```json
{
  "colors": {
    "primary": "#F97316",
    "secondary": "#FB923C",
    "background": "#0F172A",
    "surface": "#1E293B",
    "text": "#F8FAFC",
    "text_muted": "#94A3B8",
    "border": "#334155"
  },
  "typography": {
    "font_heading": "\"Space Grotesk\", sans-serif",
    "font_body": "\"DM Sans\", sans-serif",
    "base_size": "17px"
  }
}
```

---

## Template Overrides

Beyond token customization, you can override the HTML templates that individual block types use for rendering. This is an advanced feature that lets you completely change how a block is rendered for your tenant.

Template overrides are stored per-tenant in the storefront domain and are resolved at render time — if a tenant-specific template exists for a block type, it takes precedence over the platform default.

Contact your platform administrator or refer to the architecture guide for details on the template resolution pipeline.

---

## Using Theme Tokens in Custom CSS

If you add per-page custom CSS (via the page editor), reference theme tokens using CSS custom properties so your styles automatically adapt when the active theme changes:

```css
/* Good — adapts to the active theme */
.my-hero {
  background-color: var(--color-primary);
  color: var(--color-background);
  font-family: var(--font-heading);
  border-radius: var(--border-radius-lg);
  box-shadow: var(--shadow-md);
}

/* Avoid — hardcoded color breaks when theme changes */
.my-hero {
  background-color: #6366f1;
}
```

---

## FAQ

**Can I use a custom web font?**
Yes. Set `font_heading` or `font_body` to a font stack that includes the font name, and load the font separately using a `@import` in your page's custom CSS field.

**Do token changes affect published pages immediately?**
Yes. Tokens are resolved at render time, not stored with the page. Activating a new theme or updating tokens takes effect on the next page request.

**How many themes can I have?**
There is no limit on the number of themes per tenant. Only one can be active.
