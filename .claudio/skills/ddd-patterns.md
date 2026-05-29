---
name: ddd-patterns
description: "Project-specific DDD patterns for hada-commerce: 5-layer bounded context structure, naming conventions, multi-tenancy, error handling, and the net/http handler pattern. Load when scaffolding domains, writing handlers, or adding repositories."
---

# DDD Patterns — hada-commerce

Every bounded context follows the same 5-layer layout. Deviate from this and the project breaks.

## Quick rules
- HTTP framework: `net/http` stdlib only — NO Fiber, NO Gin, NO Chi
- DB access: raw SQL via `github.com/lib/pq` — NO ORM, NO query builder
- Errors: `internal/kernel/errx` only — NO `errors.New`, NO `fmt.Errorf` for domain errors
- IDs: `kernel.<Entity>ID` always — NO raw `string`
- Multi-tenancy: every query must `WHERE tenant_id = $N`
- Money: `kernel.Money{Amount int64, Currency string}` — always in cents
- Pagination: all list endpoints return `kernel.PaginatedResult[T]`

## 5-Layer structure

```
backend/internal/<entity>/
├── <entity>/
│   ├── entity.go         ← struct + DTOs
│   ├── port.go           ← Repository interface
│   └── errors.go         ← errx var declarations
├── <entity>srv/
│   └── service.go        ← business logic
├── <entity>infra/
│   └── postgres.go       ← implements Repository
├── <entity>api/
│   └── handler.go        ← net/http handlers + RegisterRoutes
└── <entity>container/
    └── container.go      ← DI wiring
```

## entity.go pattern
```go
package <entity>

import "github.com/Abraxas-365/hada-commerce/internal/kernel"

type Widget struct {
    ID          kernel.WidgetID
    TenantID    kernel.TenantID
    Name        string
    Price       kernel.Money
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type CreateWidgetRequest struct {
    Name       string
    PriceCents int64
    Currency   string
}

type UpdateWidgetRequest struct {
    Name  *string
    Price *kernel.Money
}
```

## port.go pattern
```go
package <entity>

import (
    "context"
    "github.com/Abraxas-365/hada-commerce/internal/kernel"
)

type Repository interface {
    Create(ctx context.Context, w Widget) (Widget, error)
    GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.WidgetID) (Widget, error)
    List(ctx context.Context, tenantID kernel.TenantID, p kernel.Pagination) (kernel.PaginatedResult[Widget], error)
    Update(ctx context.Context, w Widget) (Widget, error)
    Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.WidgetID) error
}
```

## errors.go pattern
```go
package <entity>

import (
    "net/http"
    "github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
)

var (
    ErrNotFound     = errx.New("WIDGET_NOT_FOUND",    "widget not found",    http.StatusNotFound)
    ErrConflict     = errx.New("WIDGET_CONFLICT",     "widget already exists", http.StatusConflict)
    ErrInvalidInput = errx.New("WIDGET_INVALID_INPUT","invalid input",        http.StatusBadRequest)
)
```

Code pattern: `DOMAIN_DESCRIPTION` in SCREAMING_SNAKE_CASE.

## service.go pattern
```go
package widgetsrv

import (
    "context"
    "github.com/Abraxas-365/hada-commerce/internal/kernel"
    "github.com/Abraxas-365/hada-commerce/internal/widget/widget"
    "github.com/google/uuid"
)

type Service struct { repo widget.Repository }

func NewService(repo widget.Repository) *Service { return &Service{repo: repo} }

func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, req widget.CreateWidgetRequest) (widget.Widget, error) {
    w := widget.Widget{
        ID:       kernel.WidgetID(uuid.New().String()),
        TenantID: tenantID,
        Name:     req.Name,
        Price:    kernel.NewMoney(req.PriceCents, req.Currency),
    }
    return s.repo.Create(ctx, w)
}
```

## postgres.go pattern
```go
package widgetinfra

import (
    "context"
    "database/sql"
    "github.com/Abraxas-365/hada-commerce/internal/kernel"
    "github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
    "github.com/Abraxas-365/hada-commerce/internal/widget/widget"
)

type PostgresRepository struct { db *sql.DB }

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
    return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.WidgetID) (widget.Widget, error) {
    const q = `SELECT id, tenant_id, name, price_cents, currency, created_at, updated_at
               FROM widgets WHERE id = $1 AND tenant_id = $2`
    var w widget.Widget
    err := r.db.QueryRowContext(ctx, q, id, tenantID).Scan(
        &w.ID, &w.TenantID, &w.Name, &w.Price.Amount, &w.Price.Currency, &w.CreatedAt, &w.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return widget.Widget{}, errx.Wrap(widget.ErrNotFound, string(id))
    }
    return w, err
}

func (r *PostgresRepository) List(ctx context.Context, tenantID kernel.TenantID, p kernel.Pagination) (kernel.PaginatedResult[widget.Widget], error) {
    const count = `SELECT COUNT(*) FROM widgets WHERE tenant_id = $1`
    var total int
    if err := r.db.QueryRowContext(ctx, count, tenantID).Scan(&total); err != nil {
        return kernel.PaginatedResult[widget.Widget]{}, err
    }

    const q = `SELECT id, tenant_id, name, price_cents, currency, created_at, updated_at
               FROM widgets WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
    rows, err := r.db.QueryContext(ctx, q, tenantID, p.Limit(), p.Offset())
    if err != nil {
        return kernel.PaginatedResult[widget.Widget]{}, err
    }
    defer rows.Close()

    var items []widget.Widget
    for rows.Next() {
        var w widget.Widget
        if err := rows.Scan(&w.ID, &w.TenantID, &w.Name, &w.Price.Amount, &w.Price.Currency, &w.CreatedAt, &w.UpdatedAt); err != nil {
            return kernel.PaginatedResult[widget.Widget]{}, err
        }
        items = append(items, w)
    }
    return kernel.NewPaginatedResult(items, total, p), nil
}
```

## handler.go pattern
```go
package widgetapi

import (
    "encoding/json"
    "net/http"
    "github.com/Abraxas-365/hada-commerce/internal/kernel"
    "github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
    "github.com/Abraxas-365/hada-commerce/internal/widget/widget"
    "github.com/Abraxas-365/hada-commerce/internal/widget/widgetsrv"
)

type Handler struct { svc *widgetsrv.Service }

func NewHandler(svc *widgetsrv.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("GET /api/v1/widgets/",       h.List)
    mux.HandleFunc("POST /api/v1/widgets/",      h.Create)
    mux.HandleFunc("GET /api/v1/widgets/{id}",   h.GetByID)
    mux.HandleFunc("DELETE /api/v1/widgets/{id}", h.Delete)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    tenantID := kernel.TenantID(tenantFromContext(r.Context()))
    var req widget.CreateWidgetRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
        return
    }
    result, err := h.svc.Create(r.Context(), tenantID, req)
    if err != nil {
        writeErrx(w, err)
        return
    }
    writeJSON(w, http.StatusCreated, result)
}

// helper — copy the writeJSON/writeError/writeErrx pattern from an existing handler
```

## container.go pattern
```go
package widgetcontainer

import (
    "database/sql"
    "net/http"
    "github.com/Abraxas-365/hada-commerce/internal/widget/widgetapi"
    "github.com/Abraxas-365/hada-commerce/internal/widget/widgetinfra"
    "github.com/Abraxas-365/hada-commerce/internal/widget/widgetsrv"
)

type Container struct { Handler *widgetapi.Handler }

func New(db *sql.DB) *Container {
    repo    := widgetinfra.NewPostgresRepository(db)
    svc     := widgetsrv.NewService(repo)
    handler := widgetapi.NewHandler(svc)
    return &Container{Handler: handler}
}

func (c *Container) RegisterRoutes(mux *http.ServeMux) {
    c.Handler.RegisterRoutes(mux)
}
```

## Migration pattern
File: `backend/migrations/NNN_<feature>.up.sql`
```sql
CREATE TABLE widgets (
    id          TEXT        NOT NULL,
    tenant_id   TEXT        NOT NULL,
    name        TEXT        NOT NULL,
    price_cents BIGINT      NOT NULL DEFAULT 0,
    currency    TEXT        NOT NULL DEFAULT 'USD',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    UNIQUE (tenant_id, name)
);

CREATE INDEX idx_widgets_tenant_id ON widgets(tenant_id);
```

Rules:
- Always index `tenant_id`
- Use `TEXT` for IDs (UUIDs stored as strings)
- Use `BIGINT` for money amounts (cents)
- Never edit existing .up.sql files — add a new one

## Gotchas
- `TenantFromContext` lives in `cmd/server.go` — import it as a package-level helper or copy the pattern to the handler package
- `sql.ErrNoRows` must be wrapped with `errx.Wrap(ErrNotFound, ...)` — never returned raw
- `rows.Close()` must be `defer`red immediately after `QueryContext`
- `kernel.NewPagination(0, 0)` returns defaults (page=1, pageSize=20) — pass user-provided values
- Don't forget to register the new container in `backend/cmd/container.go` and `server.go`
