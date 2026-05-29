# Module Structure — How to Scaffold a Domain

Use the CLI to generate new domains — never create them manually:
```bash
manifesto add internal/recruitment/candidate
```

The last path segment (`candidate`) drives all naming conventions.

---

## Generated Directory Structure

```
internal/<context>/<entity>/
├── <entity>.go                          # Entity + DTOs
├── port.go                              # Repository interface
├── errors.go                            # Error registry
├── <entity>srv/service.go               # Business logic
├── <entity>infra/postgres.go            # Data layer (Postgres)
├── <entity>api/handler.go               # HTTP handlers (Fiber)
└── <entity>container/container.go       # DI wiring
```

**Example** — `manifesto add internal/recruitment/candidate`:
```
internal/recruitment/candidate/
├── candidate.go
├── port.go
├── errors.go
├── candidatesrv/service.go
├── candidateinfra/postgres.go
├── candidateapi/handler.go
└── candidatecontainer/container.go
```

---

## Naming Conventions

All derived from the last path segment (e.g. `candidate`):

| Element | Pattern | Example |
|---------|---------|---------|
| Entity type | `PascalCase` | `Candidate` |
| ID type | `{Entity}ID` | `CandidateID` |
| DB table | plural snake_case | `candidates` |
| Service pkg | `{entity}srv` | `candidatesrv` |
| Infra pkg | `{entity}infra` | `candidateinfra` |
| API pkg | `{entity}api` | `candidateapi` |
| Container pkg | `{entity}container` | `candidatecontainer` |
| Error code prefix | `UPPER_SNAKE` | `CANDIDATE_NOT_FOUND` |

---

## Layer Responsibilities

### `<entity>.go` — Entity + DTOs
- Aggregate root struct with `ID`, `TenantID`, and domain fields
- `Create<Entity>Request` and `Update<Entity>Request` DTOs
- No business logic — pure data shapes

```go
type Candidate struct {
    ID       CandidateID
    TenantID kernel.TenantID
    Name     string
    Email    kernel.Email
    // ...
}

type CreateCandidateRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

### `port.go` — Repository Interface
- Context-first signatures
- Pagination-aware `List`
- Only the operations the service actually needs

```go
type Repository interface {
    Create(ctx context.Context, c Candidate) (Candidate, error)
    GetByID(ctx context.Context, tenantID kernel.TenantID, id CandidateID) (Candidate, error)
    List(ctx context.Context, tenantID kernel.TenantID, page, pageSize int) ([]Candidate, int, error)
    Update(ctx context.Context, c Candidate) (Candidate, error)
    Delete(ctx context.Context, tenantID kernel.TenantID, id CandidateID) error
}
```

### `errors.go` — Error Registry
- Package-level `errx` vars — never inline errors
- HTTP status codes baked in

```go
var (
    ErrNotFound  = errx.New("CANDIDATE_NOT_FOUND", "candidate not found", http.StatusNotFound)
    ErrConflict  = errx.New("CANDIDATE_CONFLICT", "candidate already exists", http.StatusConflict)
)
```

### `<entity>srv/service.go` — Business Logic
- Depends on `Repository` interface (injected via constructor)
- Generates IDs (`uuid.New()`)
- Enforces tenant scoping
- Contains all business rules

```go
type Service struct {
    repo candidate.Repository
}

func NewService(repo candidate.Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, req candidate.CreateCandidateRequest) (candidate.Candidate, error) {
    c := candidate.Candidate{
        ID:       candidate.CandidateID(uuid.New().String()),
        TenantID: tenantID,
        // map req fields...
    }
    return s.repo.Create(ctx, c)
}
```

### `<entity>infra/postgres.go` — Data Layer
- Implements `candidate.Repository`
- Raw SQL with `LIMIT`/`OFFSET` pagination
- Maps DB constraint errors to domain errors

```go
type PostgresRepository struct {
    db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
    return &PostgresRepository{db: db}
}
```

### `<entity>api/handler.go` — HTTP Handlers (Fiber)
- Routes: `POST /`, `GET /:id`, `GET /`, `PUT /:id`, `DELETE /:id`
- Extracts `TenantID` from JWT context
- Delegates all logic to service

```go
type Handler struct {
    svc *candidatesrv.Service
}

func NewHandler(svc *candidatesrv.Service) *Handler {
    return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r fiber.Router) {
    r.Post("/", h.Create)
    r.Get("/:id", h.GetByID)
    r.Get("/", h.List)
    r.Put("/:id", h.Update)
    r.Delete("/:id", h.Delete)
}
```

### `<entity>container/container.go` — DI Wiring
- Assembles: `DB → repo → service → handler`
- Exposes `RegisterRoutes(router)` to the app container

```go
type Container struct {
    Handler *candidateapi.Handler
}

func New(db *sql.DB) *Container {
    repo    := candidateinfra.NewPostgresRepository(db)
    svc     := candidatesrv.NewService(repo)
    handler := candidateapi.NewHandler(svc)
    return &Container{Handler: handler}
}

func (c *Container) RegisterRoutes(r fiber.Router) {
    c.Handler.RegisterRoutes(r.Group("/candidates"))
}
```

---

## Auto-injected into Root (marker-based, idempotent)

The CLI also modifies these root files automatically:

- **`cmd/container.go`** — adds import, struct field, and `New<Entity>Container()` call
- **`cmd/server.go`** — registers routes under `/api/v1` group
- **`internal/kernel/proj_ids.go`** — appends `type <Entity>ID string`

---

## What You Add Manually After Scaffolding

- Business rules inside service methods
- SQL schema and migrations
- Request validation in DTOs
- Authorization checks in handlers
- Advanced filtering beyond basic pagination
