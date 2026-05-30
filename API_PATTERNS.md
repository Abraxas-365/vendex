# API Implementation Patterns

## HTTP Handler Pattern

### Basic Handler Structure
```go
package customerapi

import (
    "github.com/gofiber/fiber/v2"
    "github.com/Abraxas-365/hada-commerce/internal/auth"
    "github.com/Abraxas-365/hada-commerce/internal/customer"
)

type Handler struct {
    svc *customersrv.Service
}

func NewHandler(svc *customersrv.Service) *Handler {
    return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
    router.Post("/customers", h.CreateCustomer)
    router.Get("/customers/:id", h.GetCustomer)
    router.Put("/customers/:id", h.UpdateCustomer)
    router.Delete("/customers/:id", h.DeleteCustomer)
    router.Get("/customers", h.ListCustomers)
}
```

## Request/Response Patterns

### Extract Auth Context
```go
func (h *Handler) CreateCustomer(c *fiber.Ctx) error {
    // Get auth context (set by UnifiedAuthMiddleware)
    authCtx, ok := auth.GetAuthContext(c)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Authentication required",
        })
    }
    
    tenantID := authCtx.TenantID  // Use for scoping
    // ... rest of handler
}
```

### Parse Request Body
```go
type CreateCustomerRequest struct {
    Email string            `json:"email"`
    Name  string            `json:"name"`
    Phone string            `json:"phone"`
    Addresses []Address      `json:"addresses,omitempty"`
}

func (h *Handler) CreateCustomer(c *fiber.Ctx) error {
    authCtx, _ := auth.GetAuthContext(c)
    
    var req CreateCustomerRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    // Convert to service input type
    input := customersrv.CreateInput{
        Email:     req.Email,
        Name:      req.Name,
        Phone:     req.Phone,
        Addresses: req.Addresses,
    }
    
    cust, err := h.svc.Create(c.Context(), authCtx.TenantID, input)
    if err != nil {
        return h.handleError(c, err)
    }
    
    return c.Status(fiber.StatusCreated).JSON(cust)
}
```

### Extract Path Parameters
```go
func (h *Handler) GetCustomer(c *fiber.Ctx) error {
    authCtx, _ := auth.GetAuthContext(c)
    id := kernel.CustomerID(c.Params("id"))
    
    cust, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
    if err != nil {
        return h.handleError(c, err)
    }
    
    return c.JSON(cust)
}
```

### Extract Query Parameters (Pagination)
```go
func (h *Handler) ListCustomers(c *fiber.Ctx) error {
    authCtx, _ := auth.GetAuthContext(c)
    
    page := c.QueryInt("page", 1)
    pageSize := c.QueryInt("page_size", 20)
    
    pg := kernel.NewPaginationOptions(page, pageSize)
    
    result, err := h.svc.List(c.Context(), authCtx.TenantID, pg)
    if err != nil {
        return h.handleError(c, err)
    }
    
    return c.JSON(result)
}
```

### Error Handling in Handlers
```go
func (h *Handler) handleError(c *fiber.Ctx, err error) error {
    switch {
    case errx.Is(err, customer.ErrNotFound):
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Customer not found",
        })
    case errx.Is(err, customer.ErrDuplicateEmail):
        return c.Status(fiber.StatusConflict).JSON(fiber.Map{
            "error": "Email already in use",
        })
    case errx.Is(err, customer.ErrInvalidEmail):
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid email format",
        })
    default:
        // Internal error
        logx.Errorf("Unexpected error: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Internal server error",
        })
    }
}
```

## Protected Route Pattern

### Require Specific Scope
```go
protected := app.Group("/api/v1",
    container.IAM.UnifiedAuthMiddleware.Authenticate(),
)

protected.Delete("/customers/:id",
    container.IAM.UnifiedAuthMiddleware.RequireScope("customers:delete"),
    handler.DeleteCustomer,
)
```

### Require Admin
```go
protected.Post("/admin/settings",
    container.IAM.UnifiedAuthMiddleware.RequireAdmin(),
    handler.UpdateSettings,
)
```

### Require Multiple Scopes (AND)
```go
protected.Post("/orders/:id/refund",
    container.IAM.UnifiedAuthMiddleware.RequireAllScopes("orders:read", "payments:write"),
    handler.RefundOrder,
)
```

## Pagination Response Pattern

```go
// Service returns:
type Paginated[T] struct {
    Items      []T `json:"items"`
    Total      int `json:"total"`
    Page       int `json:"page"`
    PageSize   int `json:"page_size"`
    TotalPages int `json:"total_pages"`
}

// HTTP Response (same structure):
{
  "items": [
    { "id": "cust_123", "email": "john@example.com", ... },
    ...
  ],
  "total": 150,
  "page": 1,
  "page_size": 20,
  "total_pages": 8
}
```

## Fiber Context Utilities

### Locals (Request-scoped storage)
```go
// Set by middleware
c.Locals("auth", authContext)

// Retrieve in handler
auth, ok := c.Locals("auth").(*kernel.AuthContext)
```

### Status Codes
```go
fiber.StatusOK                  // 200
fiber.StatusCreated             // 201
fiber.StatusBadRequest          // 400
fiber.StatusUnauthorized        // 401
fiber.StatusForbidden           // 403
fiber.StatusNotFound            // 404
fiber.StatusConflict            // 409
fiber.StatusInternalServerError // 500
```

### Response Types
```go
// JSON object
c.JSON(data)                    // Sets Content-Type: application/json

// With status
c.Status(fiber.StatusOK).JSON(data)

// Error map
c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
    "error": "Validation failed",
    "details": map[string]string{
        "email": "Invalid format",
    },
})
```

## Service Input/Output Pattern

### Input DTOs (Request-specific)
```go
type CreateInput struct {
    Email     string
    Name      string
    Phone     string
    Addresses []Address
}

// Handler converts HTTP request → Service Input
```

### Output is Entity
```go
func (s *Service) Create(...) (*customer.Customer, error)
func (s *Service) GetByID(...) (*customer.Customer, error)

// Entity JSON tags determine HTTP response shape
type Customer struct {
    ID        kernel.CustomerID `json:"id"`
    TenantID  kernel.TenantID   `json:"tenant_id"`
    Email     kernel.Email      `json:"email"`
    Name      string            `json:"name"`
    Addresses []Address         `json:"addresses"`
    CreatedAt time.Time         `json:"created_at"`
    UpdatedAt time.Time         `json:"updated_at"`
}
```

## Common Handler Methods

### Create (POST)
```go
func (h *Handler) Create(c *fiber.Ctx) error {
    authCtx, ok := auth.GetAuthContext(c)
    if !ok { return unauthorized(c) }
    
    var req CreateRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
    }
    
    entity, err := h.svc.Create(c.Context(), authCtx.TenantID, convertToInput(req))
    if err != nil {
        return h.handleError(c, err)
    }
    
    return c.Status(fiber.StatusCreated).JSON(entity)
}
```

### Get (GET)
```go
func (h *Handler) Get(c *fiber.Ctx) error {
    authCtx, ok := auth.GetAuthContext(c)
    if !ok { return unauthorized(c) }
    
    id := kernel.SomeID(c.Params("id"))
    
    entity, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
    if err != nil {
        return h.handleError(c, err)
    }
    
    return c.JSON(entity)
}
```

### Update (PUT)
```go
func (h *Handler) Update(c *fiber.Ctx) error {
    authCtx, ok := auth.GetAuthContext(c)
    if !ok { return unauthorized(c) }
    
    id := kernel.SomeID(c.Params("id"))
    
    var req UpdateRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
    }
    
    entity, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
    if err != nil {
        return h.handleError(c, err)
    }
    
    // Apply updates
    applyUpdates(entity, req)
    
    if err := h.svc.Update(c.Context(), entity); err != nil {
        return h.handleError(c, err)
    }
    
    return c.JSON(entity)
}
```

### Delete (DELETE)
```go
func (h *Handler) Delete(c *fiber.Ctx) error {
    authCtx, ok := auth.GetAuthContext(c)
    if !ok { return unauthorized(c) }
    
    id := kernel.SomeID(c.Params("id"))
    
    if err := h.svc.Delete(c.Context(), authCtx.TenantID, id); err != nil {
        return h.handleError(c, err)
    }
    
    return c.SendStatus(fiber.StatusNoContent)
}
```

### List (GET)
```go
func (h *Handler) List(c *fiber.Ctx) error {
    authCtx, ok := auth.GetAuthContext(c)
    if !ok { return unauthorized(c) }
    
    page := c.QueryInt("page", 1)
    pageSize := c.QueryInt("page_size", 20)
    pg := kernel.NewPaginationOptions(page, pageSize)
    
    result, err := h.svc.List(c.Context(), authCtx.TenantID, pg)
    if err != nil {
        return h.handleError(c, err)
    }
    
    return c.JSON(result)
}
```

## Helper Functions for Handlers

```go
func unauthorized(c *fiber.Ctx) error {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
        "error": "Authentication required",
    })
}

func forbidden(c *fiber.Ctx) error {
    return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
        "error": "Insufficient permissions",
    })
}

func notFound(c *fiber.Ctx) error {
    return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
        "error": "Not found",
    })
}

func badRequest(c *fiber.Ctx, message string) error {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "error": message,
    })
}

func serverError(c *fiber.Ctx) error {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error": "Internal server error",
    })
}
```
