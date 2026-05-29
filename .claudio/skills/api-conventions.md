---
name: api-conventions
description: "Backend REST API conventions for hada-commerce: route registration, tenant extraction, JSON helpers, error format, and middleware. Load when writing HTTP handlers or reviewing API correctness."
---

# API Conventions — hada-commerce

The backend uses Go's stdlib `net/http` with pattern-based routing (Go 1.22+ method+path syntax).
No external router. No Fiber. No Gin.

## Route registration

Routes use the `"METHOD /path"` pattern syntax added in Go 1.22:
```go
mux.HandleFunc("GET /api/v1/products/",        h.List)
mux.HandleFunc("POST /api/v1/products/",       h.Create)
mux.HandleFunc("GET /api/v1/products/{id}",    h.GetByID)
mux.HandleFunc("PUT /api/v1/products/{id}",    h.Update)
mux.HandleFunc("DELETE /api/v1/products/{id}", h.Delete)
```

Extract path parameters:
```go
id := r.PathValue("id")  // Go 1.22+ stdlib
```

Query parameters:
```go
page, _ := strconv.Atoi(r.URL.Query().Get("page"))
pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
```

## Tenant extraction
Tenant ID is injected by `withTenantExtraction` middleware in `backend/cmd/server.go`.
Handlers extract it with:
```go
tenantID := kernel.TenantID(TenantFromContext(r.Context()))
```

`TenantFromContext` is defined in `backend/cmd/server.go`. It returns `""` if not set.
Check for empty tenantID on protected routes and return 401.

## JSON helpers
Every handler package defines these helpers (or imports them from a shared place):
```go
func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
    writeJSON(w, status, map[string]string{"error": code, "message": message})
}

func writeErrx(w http.ResponseWriter, err error) {
    writeError(w, errx.HTTPStatus(err), errx.Code(err), errx.Message(err))
}
```

## Error response format
All errors return JSON with `error` (code) and `message` (human-readable):
```json
{"error": "PRODUCT_NOT_FOUND", "message": "product not found"}
```

HTTP status comes from the errx definition. Common cases:
- 400 Bad Request — invalid input, missing fields
- 404 Not Found — resource doesn't exist
- 409 Conflict — duplicate key/slug
- 422 Unprocessable Entity — validation failure
- 500 Internal Server Error — unexpected DB or system error

## Pagination query parameters
List endpoints accept `page` (default 1) and `page_size` (default 20, max 100):
```go
page, _ := strconv.Atoi(r.URL.Query().Get("page"))
pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
p := kernel.NewPagination(page, pageSize)  // handles defaults and clamping
```

Pagination response (always return this structure for list endpoints):
```json
{
  "items": [...],
  "total": 47,
  "page": 1,
  "page_size": 20,
  "total_pages": 3
}
```

## Middleware chain (applied in server.go)
Outermost first (applied bottom-up in code):
1. `withCORS` — adds CORS headers, handles OPTIONS preflight
2. `withRequestLogging` — logs method, path, status, duration
3. `withTenantExtraction` — reads X-Tenant-ID header → context

Handlers see: CORS headers already set, tenant in context, request logged.

## File upload (media domain)
```go
// Parse multipart form
if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max
    writeError(w, http.StatusBadRequest, "INVALID_FORM", "invalid multipart form")
    return
}
file, header, err := r.FormFile("file")
```

Uploaded files are stored locally at `uploads/` with URL prefix `/uploads/`.

## Handler checklist
- [ ] Extract tenantID from context — return 401 if empty on protected routes
- [ ] Decode body with `json.NewDecoder(r.Body).Decode(&req)` — check error
- [ ] Validate required fields — return 400 with errx
- [ ] Call service with ctx — check and handle error with writeErrx
- [ ] Return correct HTTP status (201 Created, 200 OK, 204 No Content for DELETE)

## Gotchas
- Go 1.22 routing: trailing `/` on collection routes matters (`/products/` not `/products`)
- `r.PathValue("id")` is Go 1.22+ only — do not use `mux.Vars(r)` (that's Gorilla, not here)
- `json.NewEncoder(w).Encode(v)` adds a trailing newline — that's fine
- CORS preflight (OPTIONS) is handled by `withCORS` middleware — handlers don't need to check it
- `TenantFromContext` returns empty string (not an error) — always check before using
