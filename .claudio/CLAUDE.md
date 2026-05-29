# hada-commerce — Project Context

## What Is This?

**hada-commerce** is an AI-extensible e-commerce platform. It exposes AI agents that can:
- Create and manage CMS pages (HTML + CSS, with a `pending_review → published` workflow and versioning)
- Manage products, categories, collections, and catalogs
- Create and apply promotional codes and discount rules
- Query orders and customer data
- Upload and organize media assets

The agent integration is built on **[harness](https://github.com/Abraxas-365/harness)** — each domain service is wrapped as a harness `Tool`, making it callable from the AI loop.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.25 |
| Web framework | Fiber v2 |
| Database | PostgreSQL 16 |
| Cache / Queue | Redis 7 |
| DDD scaffold | Manifesto patterns (local, not library import) |
| AI agent | github.com/Abraxas-365/harness |
| Module | `github.com/Abraxas-365/hada-commerce` |

Infrastructure is wired via `docker-compose.yml` (postgres:16-alpine, redis:7-alpine).

---

## Directory Structure

```
hada-commerce/
├── cmd/                          # Entry point + DI container + Fiber server
├── deploy/                       # Docker / k8s deployment files
├── frontend/                     # (static or SSR frontend assets)
├── migrations/                   # SQL migration files (numbered, sequential)
├── internal/
│   ├── kernel/                   # Shared value objects — import everywhere
│   │   ├── ids.go                # All domain ID types
│   │   ├── money.go              # Money (Amount int64 cents + Currency string)
│   │   ├── pagination.go         # Pagination + PaginatedResult[T]
│   │   ├── email.go              # Email value object with validation
│   │   └── errx/errx.go          # Project-local errx (see below)
│   ├── config/                   # App config loading
│   ├── agent/                    # Harness tool wrappers (AI integration layer)
│   │
│   ├── product/                  # Bounded context: products
│   │   ├── product/              # Entity + DTOs + port + errors
│   │   ├── productsrv/           # Business logic
│   │   ├── productinfra/         # Postgres repository
│   │   ├── productapi/           # Fiber HTTP handlers
│   │   └── productcontainer/     # DI wiring
│   │
│   ├── order/                    # Bounded context: orders
│   │   ├── order/
│   │   ├── ordersrv/
│   │   ├── orderinfra/
│   │   ├── orderapi/
│   │   └── ordercontainer/
│   │
│   ├── customer/                 # Bounded context: customers
│   │   ├── customer/
│   │   ├── customersrv/
│   │   ├── customerinfra/
│   │   ├── customerapi/
│   │   └── customercontainer/
│   │
│   ├── catalog/                  # Bounded context: categories + collections
│   │   ├── catalog/
│   │   ├── catalogsrv/
│   │   ├── cataloginfra/
│   │   ├── catalogapi/
│   │   └── catalogcontainer/
│   │
│   ├── storefront/               # Bounded context: CMS pages
│   │   ├── storefront/
│   │   ├── storefrontsrv/
│   │   ├── storefrontinfra/
│   │   ├── storefrontapi/
│   │   └── storefrontcontainer/
│   │
│   ├── promo/                    # Bounded context: promotions + discount codes
│   │   ├── promo/
│   │   ├── promosrv/
│   │   ├── promoinfra/
│   │   ├── promoapi/
│   │   └── promocontainer/
│   │
│   ├── media/                    # Bounded context: media assets
│   │   ├── media/
│   │   ├── mediasrv/
│   │   ├── mediainfra/
│   │   ├── mediaapi/
│   │   └── mediacontainer/
│   │
│   └── analytics/                # Bounded context: analytics / reporting
│       ├── analytics/
│       ├── analyticssrv/
│       ├── analyticsinfra/
│       ├── analyticsapi/
│       └── analyticscontainer/
└── docker-compose.yml
```

---

## Kernel Value Objects

Import path: `github.com/Abraxas-365/hada-commerce/internal/kernel`

These are the canonical types — **never use raw `string` for IDs**.

```go
// Identity types
kernel.TenantID     // string alias — multi-tenant scope key
kernel.UserID       // string alias

// Domain ID types
kernel.ProductID
kernel.OrderID
kernel.OrderItemID
kernel.CustomerID
kernel.CategoryID
kernel.CollectionID
kernel.PageID
kernel.PageVersionID
kernel.PromoID
kernel.MediaID

// Value objects
kernel.Money{Amount int64, Currency string}   // cents + ISO 4217
kernel.NewMoney(cents, "USD")
m.Add(other), m.Multiply(qty)

kernel.Email                                   // validated email
kernel.NewEmail("user@example.com")            // returns (Email, error)

kernel.Pagination{Page, PageSize}             // page ≥ 1, pageSize 1–100 (default 20)
kernel.NewPagination(page, pageSize)
p.Offset(), p.Limit()

kernel.PaginatedResult[T]{Items, Total, Page, PageSize, TotalPages}
kernel.NewPaginatedResult(items, total, pagination)
```

---

## Error Handling (errx)

**Project-local errx** at `internal/kernel/errx` — **not** the manifesto library import.
Usage is identical to manifesto errx:

```go
import "github.com/Abraxas-365/hada-commerce/internal/kernel/errx"

// Define at package level in errors.go
var (
    ErrNotFound     = errx.New("PRODUCT_NOT_FOUND", "product not found", http.StatusNotFound)
    ErrConflict     = errx.New("PRODUCT_CONFLICT", "product already exists", http.StatusConflict)
    ErrInvalidInput = errx.New("PRODUCT_INVALID_INPUT", "invalid input", http.StatusBadRequest)
)

// Return
return errx.Wrap(ErrNotFound, "id: "+id.String())

// Check
if errx.Is(err, ErrNotFound) { ... }

// HTTP handler
status := errx.HTTPStatus(err)
code   := errx.Code(err)
msg    := errx.Message(err)
```

**Code convention:** `<DOMAIN>_<DESCRIPTION>` in SCREAMING_SNAKE_CASE.

---

## Domain Module Pattern

Every bounded context follows the same layout. Example — `product`:

```
internal/product/
├── product/
│   ├── product.go       # Entity + DTOs
│   ├── port.go          # Repository interface
│   └── errors.go        # Error registry (errx vars)
├── productsrv/
│   └── service.go       # Business logic (depends on Repository interface)
├── productinfra/
│   └── postgres.go      # Implements product.Repository via raw SQL
├── productapi/
│   └── handler.go       # Fiber handlers — extracts TenantID from JWT, delegates to svc
└── productcontainer/
    └── container.go     # Wires db → repo → svc → handler; exposes RegisterRoutes()
```

### Naming rules (derived from entity name, e.g. `product`):

| Element | Pattern | Example |
|---------|---------|---------|
| Entity type | PascalCase | `Product` |
| ID type | `{Entity}ID` | `ProductID` |
| DB table | plural snake_case | `products` |
| Service pkg | `{entity}srv` | `productsrv` |
| Infra pkg | `{entity}infra` | `productinfra` |
| API pkg | `{entity}api` | `productapi` |
| Container pkg | `{entity}container` | `productcontainer` |
| Error codes | `DOMAIN_DESCRIPTION` | `PRODUCT_NOT_FOUND` |

### Entity structure (entity.go):
```go
type Product struct {
    ID       kernel.ProductID
    TenantID kernel.TenantID
    // domain fields...
}

type CreateProductRequest struct { /* ... */ }
type UpdateProductRequest struct { /* ... */ }
```

### Repository interface (port.go):
```go
type Repository interface {
    Create(ctx context.Context, p Product) (Product, error)
    GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ProductID) (Product, error)
    List(ctx context.Context, tenantID kernel.TenantID, p kernel.Pagination) (kernel.PaginatedResult[Product], error)
    Update(ctx context.Context, p Product) (Product, error)
    Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.ProductID) error
}
```

### Service (productsrv/service.go):
```go
type Service struct { repo product.Repository }

func NewService(repo product.Repository) *Service { return &Service{repo: repo} }

func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, req product.CreateProductRequest) (product.Product, error) {
    p := product.Product{
        ID:       kernel.ProductID(uuid.New().String()),
        TenantID: tenantID,
        // map req...
    }
    return s.repo.Create(ctx, p)
}
```

### Container (productcontainer/container.go):
```go
type Container struct { Handler *productapi.Handler }

func New(db *sql.DB) *Container {
    repo    := productinfra.NewPostgresRepository(db)
    svc     := productsrv.NewService(repo)
    handler := productapi.NewHandler(svc)
    return &Container{Handler: handler}
}

func (c *Container) RegisterRoutes(r fiber.Router) {
    c.Handler.RegisterRoutes(r.Group("/products"))
}
```

---

## CMS Page System (storefront domain)

The `storefront` domain manages CMS pages (HTML + CSS) with versioning and workflow:

### Workflow states:
```
draft → pending_review → published
                       → rejected → draft (revise and resubmit)
```

### Key concepts:
- **Page** (`kernel.PageID`) — top-level content unit with slug, title, and current published version
- **PageVersion** (`kernel.PageVersionID`) — immutable snapshot of HTML + CSS content
- Agents create new versions; humans (or auto-approval) publish them
- Published versions are served at the storefront; older versions are archived but queryable

### Agent actions on storefront:
1. `CreatePage(tenantID, slug, title)` → creates page in draft state
2. `CreatePageVersion(tenantID, pageID, html, css)` → creates a new version in `pending_review`
3. `PublishPageVersion(tenantID, pageVersionID)` → moves to `published`
4. `GetPage(tenantID, slug)` → returns published page HTML + CSS

---

## Agent Integration Layer (internal/agent/)

The `internal/agent/` package wraps domain services as **harness tools**.

### Pattern:
```go
// Each tool implements tools.Tool from github.com/Abraxas-365/harness
type CreateProductTool struct {
    svc *productsrv.Service
}

func (t *CreateProductTool) Name() string        { return "create_product" }
func (t *CreateProductTool) Description() string { return "Create a new product in the catalog" }
func (t *CreateProductTool) InputSchema() json.RawMessage { /* JSON schema */ }
func (t *CreateProductTool) IsReadOnly() bool    { return false }
func (t *CreateProductTool) RequiresApproval(input json.RawMessage) bool { return false }

func (t *CreateProductTool) Execute(ctx context.Context, input json.RawMessage) (*tools.Result, error) {
    var req CreateProductInput
    if err := json.Unmarshal(input, &req); err != nil {
        return &tools.Result{Content: "invalid input: " + err.Error(), IsError: true}, nil
    }
    product, err := t.svc.Create(ctx, kernel.TenantID(req.TenantID), req.toServiceRequest())
    if err != nil {
        return &tools.Result{Content: errx.Message(err), IsError: true}, nil
    }
    result, _ := json.Marshal(product)
    return &tools.Result{Content: string(result)}, nil
}
```

### Tool registration in harness:
```go
h := harness.New(apiKey,
    harness.WithTools(
        agent.NewCreateProductTool(productSvc),
        agent.NewCreatePageTool(storefrontSvc),
        agent.NewCreatePromoTool(promoSvc),
        // ...
    ),
)
```

---

## Key Conventions

1. **Always multi-tenant**: every DB query must filter by `kernel.TenantID`
2. **errx for all domain errors**: never `errors.New` or `fmt.Errorf` for domain-level errors
3. **Context first**: `ctx context.Context` is always the first parameter
4. **IDs from kernel**: never use raw `string` for any domain ID
5. **Money in cents**: `kernel.Money.Amount` is always in smallest currency unit
6. **Pagination everywhere**: list endpoints always return `kernel.PaginatedResult[T]`
7. **Tool errors vs Go errors**: harness tools return `&tools.Result{IsError: true}` for user-visible errors; Go errors only for framework failures
8. **No global state**: all dependencies injected via constructors

---

## Manifesto Skill

The `manifesto` skill (`.claudio/skills/manifesto/`) documents the DDD patterns used in this project. Load it when scaffolding new domains, adding services, or implementing repository methods.
