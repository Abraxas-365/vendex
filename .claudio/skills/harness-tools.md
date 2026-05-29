---
name: harness-tools
description: "How to write harness Tool wrappers in backend/internal/agent/. Load when adding new tools to the AI agent loop, or when modifying existing tools in tools.go."
---

# Harness Tools — hada-commerce

Tools in `backend/internal/agent/` expose domain services to the AI agent loop.
Read `backend/internal/agent/tools.go` before writing any new tool — follow its patterns exactly.

## Local Tool interface
The project defines its own Tool interface (not the upstream harness library):

```go
// backend/internal/agent/ — see handler.go for the actual definition
type Tool interface {
    Name()        string
    Description() string
    InputSchema() map[string]any
    Execute(ctx context.Context, raw json.RawMessage) (string, error)
}
```

## Full tool template
```go
// ─── CreateWidgetTool ────────────────────────────────────────────────────────

type CreateWidgetTool struct {
    svc      *widgetsrv.Service
    tenantID kernel.TenantID
}

type createWidgetInput struct {
    Name       string `json:"name"`
    PriceCents int64  `json:"price_cents"`
    Currency   string `json:"currency"`
}

func (t *CreateWidgetTool) Name() string { return "create_widget" }

func (t *CreateWidgetTool) Description() string {
    return "Create a new widget in the catalog. Returns the created widget's ID and status."
}

func (t *CreateWidgetTool) InputSchema() map[string]any {
    return map[string]any{
        "type": "object",
        "properties": map[string]any{
            "name":        map[string]any{"type": "string",  "description": "Widget name"},
            "price_cents": map[string]any{"type": "integer", "description": "Price in smallest currency unit (e.g. cents)"},
            "currency":    map[string]any{"type": "string",  "description": "ISO 4217 currency code, e.g. USD"},
        },
        "required": []string{"name", "price_cents", "currency"},
    }
}

func (t *CreateWidgetTool) Execute(ctx context.Context, raw json.RawMessage) (string, error) {
    var in createWidgetInput
    if err := json.Unmarshal(raw, &in); err != nil {
        return "", fmt.Errorf("create_widget: unmarshal input: %w", err)
    }

    w, err := t.svc.Create(ctx, t.tenantID, widgetsrv.CreateInput{
        Name:       in.Name,
        PriceCents: in.PriceCents,
        Currency:   in.Currency,
    })
    if err != nil {
        return "", fmt.Errorf("create_widget: %w", err)
    }

    return fmt.Sprintf("Widget created.\nID: %s\nName: %s\nPrice: %d %s",
        w.ID, w.Name, w.Price.Amount, w.Price.Currency), nil
}

var _ Tool = (*CreateWidgetTool)(nil)  // compile-time interface guard
```

## Naming conventions
| Element | Pattern | Example |
|---------|---------|---------|
| Tool struct | `<Action><Domain>Tool` | `CreateProductTool` |
| Tool name | `verb_noun` snake_case | `"create_product"` |
| Input struct | `<action><Domain>Input` | `createProductInput` (unexported) |

## Error handling
```go
// Unmarshal errors = framework failure → return Go error
if err := json.Unmarshal(raw, &in); err != nil {
    return "", fmt.Errorf("tool_name: unmarshal: %w", err)
}

// Service errors = user-visible → wrap and return as Go error
// The harness agent will surface the error.Message to the LLM
if err != nil {
    return "", fmt.Errorf("tool_name: %w", err)
}
```

## InputSchema rules
- Use `map[string]any` for the schema (not json.RawMessage)
- Every property must have a `"description"` key — the LLM reads these
- Always set `"required"` for mandatory fields
- Optional fields: omit from `"required"`, add `omitempty` to Go struct tag
- Use `"integer"` not `"number"` for int64/int fields
- Use `"array"` with `"items"` for slices

## Output format
Return a human-readable string the LLM can parse:
```go
// Good — structured and readable
return fmt.Sprintf("Widget created.\nID: %s\nName: %s\nStatus: %s", w.ID, w.Name, w.Status), nil

// For lists — header + one item per line
out := fmt.Sprintf("Widgets (page %d of %d, total %d):\n\n", result.Page, result.TotalPages, result.Total)
for _, w := range result.Items {
    out += fmt.Sprintf("- ID: %s | Name: %s | Price: %d %s\n", w.ID, w.Name, w.Price.Amount, w.Price.Currency)
}
return out, nil
```

## Tool registration
Find where the harness instance is built (check `backend/internal/agent/setup.go`).
Add the new tool to the slice of tools passed to the harness constructor.

## Existing tools (don't duplicate)
- `create_page` — CreatePageTool
- `update_page` — UpdatePageTool
- `list_pages` — ListPagesTool
- `create_product` — CreateProductTool
- `list_products` — ListProductsTool
- `create_promo` — CreatePromoTool
- `query_orders` — QueryOrdersTool
- `search_catalog` — SearchCatalogTool

## Gotchas
- Always add `var _ Tool = (*YourTool)(nil)` — if the interface changes, you'll catch it at compile time
- The `tenantID` field is pre-set at construction time (from config or request context) — tools don't receive it from the LLM
- Never call other tools from within a tool's Execute — call the service directly
- Time values from LLM: accept as Unix timestamps (int64), convert with `time.Unix(ts, 0).UTC()`
