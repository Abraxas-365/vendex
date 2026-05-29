---
name: manifesto
description: "Guide for writing Go code following the Manifesto DDD patterns used in hada-commerce. Use when implementing new domains, services, repositories, HTTP handlers, or agent tools. Covers: module structure, errx error pattern, kernel types, async, logging, file storage, IAM, and LLM/AI integration."
---

# Manifesto DDD Patterns — hada-commerce

This project follows Manifesto conventions for domain-driven design in Go.
**All imports use the project module** `github.com/Abraxas-365/hada-commerce/internal/...`
(not the upstream manifesto library).

---

## Package Map

| Concern | Import path | Reference |
|---------|-------------|-----------|
| Error types | `internal/kernel/errx` | `references/errors.md` |
| Kernel value objects | `internal/kernel` | IDs, Money, Pagination, Email |
| Module/domain structure | — | `references/module-structure.md` |
| Async primitives | `asyncx` (if added) | `references/async.md` |
| Logging | `logx` (if added) | `references/infra.md` |
| File storage | `fsx` (if added) | `references/infra.md` |
| Job queue | `jobx` (if added) | `references/infra.md` |
| LLM / AI / embeddings | `ai/*` (if added) | `references/llm-ai.md` |
| IAM (auth, tenant, user) | `iam/*` (if added) | `references/iam.md` |

---

## Core Conventions

### errx — domain errors
Define at package level in `errors.go` — never use `errors.New` or `fmt.Errorf` for domain errors:
```go
import "github.com/Abraxas-365/hada-commerce/internal/kernel/errx"

var (
    ErrNotFound     = errx.New("PRODUCT_NOT_FOUND", "product not found", http.StatusNotFound)
    ErrConflict     = errx.New("PRODUCT_CONFLICT", "product already exists", http.StatusConflict)
    ErrInvalidInput = errx.New("PRODUCT_INVALID_INPUT", "invalid input", http.StatusBadRequest)
)
```

### Kernel types — always use, never raw strings
```go
import "github.com/Abraxas-365/hada-commerce/internal/kernel"

kernel.TenantID, kernel.UserID
kernel.ProductID, kernel.OrderID, kernel.CustomerID, kernel.CategoryID
kernel.CollectionID, kernel.PageID, kernel.PageVersionID, kernel.PromoID, kernel.MediaID

kernel.Money{Amount: 1999, Currency: "USD"}   // always cents
kernel.NewMoney(1999, "USD")

kernel.Email                                   // validated
kernel.NewEmail("user@example.com")

kernel.Pagination, kernel.NewPagination(page, pageSize)
kernel.PaginatedResult[T], kernel.NewPaginatedResult(items, total, p)
```

### Layered architecture
```
HTTP Handler → Service → Domain (entity + port) → Infrastructure (Postgres)
```
Dependencies point inward. Infrastructure implements interfaces defined in the domain package.

### Multi-tenancy
All domain operations are scoped by `kernel.TenantID`. Every DB query must filter by tenant.

### Options pattern
For configurable calls, use variadic `With*` options — never bare config structs.

---

## Reference Files

When working in one of these areas, read the relevant reference file first:

| Area | Read this file |
|------|----------------|
| Module/domain scaffolding, directory layout, layer patterns | `.claudio/skills/manifesto/references/module-structure.md` |
| Async: Future, Map, Pool, Retry, Race | `.claudio/skills/manifesto/references/async.md` |
| Jobs, file storage, logging, email, config | `.claudio/skills/manifesto/references/infra.md` |
| Error types, registries, HTTP mapping | `.claudio/skills/manifesto/references/errors.md` |
| Auth, tenant, user, apikey, OTP, invitation | `.claudio/skills/manifesto/references/iam.md` |
| LLM, embeddings, vector store, agentx, memoryx | `.claudio/skills/manifesto/references/llm-ai.md` |

---

## New Domain Checklist

- [ ] `internal/<context>/<entity>/entity.go` — struct with `ID kernel.<Entity>ID`, `TenantID kernel.TenantID`
- [ ] `internal/<context>/<entity>/port.go` — `Repository` interface, context-first signatures
- [ ] `internal/<context>/<entity>/errors.go` — `errx.New(...)` vars, `DOMAIN_DESCRIPTION` codes
- [ ] `internal/<context>/<entity>/<entity>srv/service.go` — depends on `Repository` interface
- [ ] `internal/<context>/<entity>/<entity>infra/postgres.go` — implements `Repository`
- [ ] `internal/<context>/<entity>/<entity>api/handler.go` — Fiber, extracts TenantID from JWT
- [ ] `internal/<context>/<entity>/<entity>container/container.go` — wires db→repo→svc→handler
- [ ] All DB queries filter by `TenantID`
- [ ] IDs generated with `uuid.New().String()` cast to the kernel type
- [ ] List endpoints return `kernel.PaginatedResult[T]`
- [ ] No `errors.New` or `fmt.Errorf` for domain errors

## New Harness Tool Checklist (internal/agent/)

- [ ] Implements `tools.Tool` interface (Name, Description, InputSchema, Execute, IsReadOnly, RequiresApproval)
- [ ] `Execute` returns `&tools.Result{IsError: true}` for user-visible errors (not Go errors)
- [ ] `Execute` is goroutine-safe (no shared mutable state)
- [ ] `RequiresApproval` returns `true` for destructive mutations
- [ ] Input JSON schema defined inline as `json.RawMessage`
- [ ] Registered in the harness instance in `cmd/`
