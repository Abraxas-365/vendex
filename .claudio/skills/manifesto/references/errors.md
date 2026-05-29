# Error Handling (errx)

Import: `github.com/Abraxas-365/hada-commerce/internal/kernel/errx`

**Rule:** Never use `errors.New`, `fmt.Errorf`, or `errors.Wrap` for domain errors. Always use `errx`.

---

## Defining errors

Declare at package level — one var block per module:

```go
var (
    ErrNotFound      = errx.New("USER_NOT_FOUND", "user not found", http.StatusNotFound)
    ErrConflict      = errx.New("USER_CONFLICT", "user already exists", http.StatusConflict)
    ErrInvalidInput  = errx.New("USER_INVALID_INPUT", "invalid input", http.StatusBadRequest)
    ErrUnauthorized  = errx.New("USER_UNAUTHORIZED", "unauthorized", http.StatusUnauthorized)
)
```

Code convention: `<MODULE>_<DESCRIPTION>` in SCREAMING_SNAKE_CASE.

---

## Returning errors

```go
// Plain
return ErrNotFound

// With added context
return errx.Wrap(ErrNotFound, "id: "+id)
```

---

## Checking errors

```go
if errx.Is(err, ErrNotFound) {
    // handle specifically
}
```

---

## HTTP handlers

```go
status := errx.HTTPStatus(err) // extracts registered status code
code   := errx.Code(err)       // extracts the error code string
msg    := errx.Message(err)    // extracts the human-readable message
```

Use these to build consistent JSON error responses in HTTP handlers.

---

## Pattern: module error registry

Each bounded context owns its error vars. Never share errors across packages — define equivalent errors per package with their own codes. This keeps HTTP status mapping predictable and decoupled.

---

## errx API (project-local implementation)

```go
// Constructor
errx.New(code string, message string, httpStatus int) *Error

// Wrap with detail (preserves code + status)
errx.Wrap(err *Error, detail string) *Error

// Identity check (by error code)
errx.Is(err error, target *Error) bool

// HTTP helpers
errx.HTTPStatus(err error) int   // defaults to 500
errx.Code(err error) string      // defaults to "INTERNAL_ERROR"
errx.Message(err error) string   // defaults to "internal error"
```

---

## Harness tool errors

In `internal/agent/` tool implementations, translate domain errors to tool results:

```go
product, err := svc.Create(ctx, tenantID, req)
if err != nil {
    return &tools.Result{Content: errx.Message(err), IsError: true}, nil
}
```

Return Go errors only for framework-level failures (serialization, nil pointer, etc.).
