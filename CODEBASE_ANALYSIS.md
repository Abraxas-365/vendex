# Hada-Commerce Backend Codebase Analysis

## Project Overview
- **Module**: `github.com/Abraxas-365/hada-commerce`
- **Language**: Go 1.25.4
- **Architecture**: Modular DDD-inspired with separate domains (Customer, Order, Product, IAM, Payment, etc.)
- **Web Framework**: Fiber v2 (lightweight HTTP framework)
- **Database**: PostgreSQL with sqlx
- **Cache/Queue**: Redis
- **Authentication**: JWT tokens + API Keys + OAuth (Google, Microsoft)
- **Notifications**: AWS SES or console (development)
- **File Storage**: Local or AWS S3

---

## Architecture Overview

### High-Level Structure
```
backend/
├── cmd/                           # Entry point
│   ├── server.go                 # Fiber app setup, routes, middleware
│   └── container.go              # Dependency injection container
├── internal/
│   ├── kernel/                   # Shared domain types
│   ├── config/                   # Configuration loading
│   ├── iam/                       # Authentication & Authorization
│   ├── customer/                 # Customer domain
│   ├── order/                    # Order domain
│   ├── product/                  # Product domain
│   ├── payment/                  # Payment domain
│   ├── cart/                     # Shopping cart domain
│   ├── eventbus/                 # Domain events
│   └── [... other domains ...]
├── migrations/                   # Database migration files
└── go.mod                        # Module dependencies
```

### Module/Domain Pattern (Manifesto DDD)
Each domain follows a consistent 3-layer pattern:

```
internal/{domain}/
├── {domain}.go                   # Entity definitions + domain methods
├── errors.go                     # Domain-specific errors
├── port.go                       # Repository/Service interfaces
├── {domain}srv/
│   └── service.go               # Business logic service
├── {domain}infra/
│   └── postgres.go              # PostgreSQL repository implementation
└── {domain}container/
    └── container.go             # Dependency injection & wiring
```

---

## Key Types & Kernel

### ID Types (string-based)
**Location**: `internal/kernel/`

```go
// Common types
type UserID string      // User in the IAM system
type TenantID string    // Multi-tenant organization
type CustomerID string  // Customer in commerce domain

// Commerce domain types
type OrderID string
type OrderItemID string
type ProductID string
type CartID string
type CartItemID string
type CategoryID string
type PaymentID string
type RefundID string
type ShippingZoneID string
type TaxRateID string

// Each type has helper methods:
func NewCustomerID(id string) CustomerID
func (c CustomerID) String() string
func (c CustomerID) IsEmpty() bool
```

### Auth Context Type
**Location**: `internal/kernel/context.go`

```go
type AuthContext struct {
    UserID   *UserID  // Pointer because JWT auth fills it, API key might not
    TenantID TenantID // Always required
    Email    string   // User's email from JWT claims
    Name     string   // User's name from JWT claims
    Scopes   []string // Authorization scopes ["orders:read", "products:*", etc.]
    IsAPIKey bool     // True if authenticated via API key
}

// Scope validation methods:
func (ac *AuthContext) HasScope(scope string) bool
func (ac *AuthContext) HasAnyScope(scopes ...string) bool
func (ac *AuthContext) HasAllScopes(scopes ...string) bool
```

### Money Type
**Location**: `internal/kernel/money.go`

```go
type Money struct {
    Amount   int64  // in cents
    Currency string // ISO 4217 code (USD, EUR, etc.)
}
```

### Pagination
**Location**: `internal/kernel/pagination.go`

```go
// Two types (old and new pattern):
type Pagination struct { Page, PageSize int }
type PaginationOptions struct { Page, PageSize int }

type Paginated[T any] struct {
    Items      []T  `json:"items"`
    Total      int  `json:"total"`
    Page       int  `json:"page"`
    PageSize   int  `json:"page_size"`
    TotalPages int  `json:"total_pages"`
}
```

---

## Domain: Customer

### Customer Entity
**File**: `internal/customer/customer.go`

```go
type Customer struct {
    ID        kernel.CustomerID // UUID string
    TenantID  kernel.TenantID   // Multi-tenant scoping
    Email     kernel.Email      // Validated email
    Name      string            // Full name
    Phone     string            // Phone number
    Addresses []Address         // Nested array of addresses
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Address struct {
    Street     string `json:"street"`
    City       string `json:"city"`
    State      string `json:"state"`
    Country    string `json:"country"`
    PostalCode string `json:"postal_code"`
    IsDefault  bool   `json:"is_default"`
}

// Domain methods:
func (c *Customer) DefaultAddress() *Address
func (c *Customer) AddAddress(addr Address)
func (c *Customer) SetDefaultAddress(idx int) bool
```

### Repository Interface
**File**: `internal/customer/port.go`

```go
type Repository interface {
    Create(ctx context.Context, c *Customer) error
    GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerID) (*Customer, error)
    GetByEmail(ctx context.Context, tenantID kernel.TenantID, email kernel.Email) (*Customer, error)
    Update(ctx context.Context, c *Customer) error
    Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerID) error
    List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[Customer], error)
}
```

### Service Layer
**File**: `internal/customer/customersrv/service.go`

```go
type Service struct {
    repo customer.Repository
    bus  eventbus.Bus  // Event publishing
}

type CreateInput struct {
    Email     string
    Name      string
    Phone     string
    Addresses []customer.Address
}

// Main methods:
func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, in CreateInput) (*customer.Customer, error)
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerID) (*customer.Customer, error)
func (s *Service) GetByEmail(ctx context.Context, tenantID kernel.TenantID, email kernel.Email) (*customer.Customer, error)
func (s *Service) Update(ctx context.Context, c *customer.Customer) error
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[customer.Customer], error)

// Publishes events:
// - eventbus.CustomerRegistered (on Create)
// - eventbus.CustomerUpdated (on Update)
```

### Repository Implementation (PostgreSQL)
**File**: `internal/customer/customerinfra/postgres.go`

```go
type PostgresRepo struct {
    db *sqlx.DB
}

// Key implementation details:
// - Addresses stored as JSON in DB
// - Uses json.Marshal/Unmarshal for nested address array
// - All operations scoped by tenant_id (multi-tenant)
// - Uses prepared statements with $1, $2, etc. (sqlx placeholders)
// - Error handling via errx package

// Scanner pattern for both Row and Rows:
func scanCustomer(row *sql.Row) (*customer.Customer, error)
func scanCustomerRow(rows *sql.Rows) (*customer.Customer, error)
func scanFields(s scanner) (*customer.Customer, error)
```

### Container/Wiring
**File**: `internal/customer/customercontainer/container.go`

```go
type Container struct {
    Service *customersrv.Service
    Handler *customerapi.Handler
}

func New(db *sqlx.DB, bus eventbus.Bus) *Container {
    repo := customerinfra.NewPostgresRepo(db)
    svc := customersrv.New(repo, bus)
    handler := customerapi.NewHandler(svc)
    return &Container{Service: svc, Handler: handler}
}

func (c *Container) RegisterRoutes(router fiber.Router) {
    c.Handler.RegisterRoutes(router)
}
```

### Error Handling
**File**: `internal/customer/errors.go`

```go
var (
    ErrNotFound       = errx.New("customer not found", errx.TypeNotFound)
    ErrDuplicateEmail = errx.New("customer with this email already exists", errx.TypeConflict)
    ErrInvalidEmail   = errx.New("invalid email address", errx.TypeValidation)
)
```

---

## Domain: Order

### Order Entity
**File**: `internal/order/order.go` (not fully shown, but used by service)

```go
type Order struct {
    ID              kernel.OrderID
    TenantID        kernel.TenantID
    CustomerID      kernel.CustomerID
    Items           []OrderItem
    Status          OrderStatus  // Pending, Confirmed, Shipped, Delivered, Cancelled
    ShippingAddress Address
    TotalAmount     kernel.Money
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

type OrderItem struct {
    ID          kernel.OrderItemID
    ProductID   kernel.ProductID
    ProductName string
    Quantity    int
    UnitPrice   kernel.Money
}
```

### Service Layer
**File**: `internal/order/ordersrv/service.go`

```go
type Service struct {
    repo order.Repository
    bus  eventbus.Bus
}

type CreateInput struct {
    CustomerID      kernel.CustomerID
    Items           []CreateItemInput
    ShippingAddress order.Address
}

// Core methods:
func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, in CreateInput) (*order.Order, error)
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.OrderID) (*order.Order, error)
func (s *Service) UpdateStatus(ctx context.Context, tenantID kernel.TenantID, id kernel.OrderID, newStatus order.OrderStatus) (*order.Order, error)
func (s *Service) Cancel(ctx context.Context, tenantID kernel.TenantID, id kernel.OrderID) (*order.Order, error)
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[order.Order], error)
func (s *Service) ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, pg kernel.PaginationOptions) (kernel.Paginated[order.Order], error)

// Publishes events:
// - eventbus.OrderPlaced
// - eventbus.OrderConfirmed
// - eventbus.OrderShipped
// - eventbus.OrderDelivered
// - eventbus.OrderCancelled
```

---

## Authentication & Authorization (IAM)

### JWT Service
**File**: `internal/iam/auth/jwt_service.go`

```go
type JWTService struct {
    secretKey       []byte
    accessTokenTTL  time.Duration  // Default: 15 minutes
    refreshTokenTTL time.Duration  // Default: 7 days
    issuer          string         // e.g., "manifesto"
    audience        []string       // e.g., ["manifesto-api"]
}

type JWTClaims struct {
    UserID   kernel.UserID   `json:"user_id"`
    TenantID kernel.TenantID `json:"tenant_id"`
    Email    string          `json:"email"`
    Name     string          `json:"name"`
    Scopes   []string        `json:"scopes"`
    jwt.RegisteredClaims
}

// Token interface (TokenService):
func (j *JWTService) GenerateAccessToken(userID kernel.UserID, tenantID kernel.TenantID, claims map[string]any) (string, error)
func (j *JWTService) ValidateAccessToken(tokenString string) (*TokenClaims, error)
func (j *JWTService) GenerateRefreshToken(userID kernel.UserID) (string, error)

// Returns TokenClaims type with IssuedAt and ExpiresAt
```

### Password Hashing
**File**: `internal/iam/auth/authinfra/bcrypt_password_service.go`

```go
type BcryptPasswordService struct {
    cost int  // bcrypt cost factor (default: 10)
}

// Implements user.PasswordService interface:
func (s *BcryptPasswordService) HashPassword(password string) (string, error)
func (s *BcryptPasswordService) VerifyPassword(hashedPassword, password string) bool
```

### Unified Auth Middleware
**File**: `internal/iam/auth/unified_middleware.go`

```go
type UnifiedAuthMiddleware struct {
    apiKeyService *apikeysrv.APIKeyService
    tokenService  TokenService  // JWT service
}

// Main handler:
func (am *UnifiedAuthMiddleware) Authenticate() fiber.Handler
    // Tries API Key first, then JWT token

// Extracts tokens from:
// - Authorization: Bearer <token> header
// - Authorization: X-API-Key <key> header
// - X-API-Key header
// - api_key query parameter
// - access_token cookie

// Sets c.Locals("auth") = *kernel.AuthContext

// Scope validation methods:
func (am *UnifiedAuthMiddleware) RequireScope(scope string) fiber.Handler
func (am *UnifiedAuthMiddleware) RequireAnyScope(scopes ...string) fiber.Handler
func (am *UnifiedAuthMiddleware) RequireAllScopes(scopes ...string) fiber.Handler
func (am *UnifiedAuthMiddleware) RequireAdmin() fiber.Handler
```

### Auth Types
**File**: `internal/iam/auth/port.go` (auth.go)

```go
// Token types:
type RefreshToken struct {
    ID        string
    Token     string
    UserID    kernel.UserID
    TenantID  kernel.TenantID
    ExpiresAt time.Time
    CreatedAt time.Time
    IsRevoked bool
}

type UserSession struct {
    ID           string
    UserID       kernel.UserID
    TenantID     kernel.TenantID
    SessionToken string
    IPAddress    string
    UserAgent    string
    ExpiresAt    time.Time
    CreatedAt    time.Time
    LastActivity time.Time
}

type PasswordResetToken struct {
    ID        string
    Token     string
    UserID    kernel.UserID
    ExpiresAt time.Time
    CreatedAt time.Time
    IsUsed    bool
}

type TokenClaims struct {
    UserID    kernel.UserID
    TenantID  kernel.TenantID
    Email     string
    Name      string
    Scopes    []string
    IssuedAt  time.Time
    ExpiresAt time.Time
}

// Repository interfaces:
type TokenRepository interface
type SessionRepository interface
type PasswordResetRepository interface
type TokenService interface       // Implemented by JWTService
type AuditService interface
type Invitation interface
```

---

## Configuration

### Main Config
**File**: `internal/config/config.go`

```go
type Config struct {
    Server       ServerConfig
    Database     DatabaseConfig
    Redis        RedisConfig
    Environment  Environment  // development, staging, production
    Auth         AuthConfig
    OAuth        OAuthConfig
    TenantConfig TenantConfig
    Jobx         JobxConfig
    Notifx       NotifxConfig
}

func Load() (*Config, error)  // Loads from environment variables
```

### Auth Config
**File**: `internal/config/auth.go`

```go
type AuthConfig struct {
    JWT           JWTConfig
    APIKey        APIKeyConfig
    Session       SessionConfig
    OTP           OTPConfig
    Invitation    InvitationConfig
    PasswordReset PasswordResetConfig
    Cookie        CookieConfig
    Password      PasswordConfig
}

type JWTConfig struct {
    SecretKey       string        // Required, min 32 chars
    AccessTokenTTL  time.Duration // Default: 15 min (JWT_ACCESS_TOKEN_TTL)
    RefreshTokenTTL time.Duration // Default: 7 days (JWT_REFRESH_TOKEN_TTL)
    Issuer          string        // Default: "manifesto" (JWT_ISSUER)
    Audience        []string      // Default: ["manifesto-api"] (JWT_AUDIENCE)
}

type SessionConfig struct {
    ExpirationTime  time.Duration // Default: 24h (SESSION_EXPIRATION_TIME)
    CleanupInterval time.Duration // Default: 1h (SESSION_CLEANUP_INTERVAL)
    MaxSessions     int            // Default: 10 (SESSION_MAX_PER_USER)
}

type OTPConfig struct {
    CodeLength      int            // Default: 6
    ExpirationTime  time.Duration // Default: 10m
    MaxAttempts     int            // Default: 5
    RateLimitWindow time.Duration // Default: 1m
    TokenByteLength int            // Default: 3
}

type PasswordResetConfig struct {
    TokenByteLength      int            // Default: 32
    ExpirationTime       time.Duration // Default: 1h
    RateLimitWindow      time.Duration // Default: 15m
    MaxAttemptsPerWindow int            // Default: 3
}

type CookieConfig struct {
    AccessTokenName  string // Default: "access_token"
    RefreshTokenName string // Default: "refresh_token"
    Domain           string
    Path             string // Default: "/"
    Secure           bool   // Default: false
    HTTPOnly         bool   // Default: true
    SameSite         string // Default: "Lax"
}

type PasswordConfig struct {
    BcryptCost int  // Default: 10 (BCRYPT_COST)
}
```

---

## Dependency Container (Main Wiring)

### Main Container
**File**: `cmd/container.go`

```go
type Container struct {
    Config *config.Config

    // Infrastructure
    DB             *sqlx.DB
    Redis          *redis.Client
    FileSystem     fsx.FileSystem
    S3Client       *s3.Client
    JobClient      *jobx.Client
    NotifxClient   *notifx.Client
    EventBus       eventbus.Bus

    // Modules (each is a *{domain}container.Container)
    IAM            *iamcontainer.Container
    Cart           *cartcontainer.Container
    Product        *productcontainer.Container
    Order          *ordercontainer.Container
    Payment        *paymentcontainer.Container
    Customer       *customercontainer.Container
    Catalog        *catalogcontainer.Container
    Storefront     *storefrontcontainer.Container
    Promo          *promocontainer.Container
    Media          *mediacontainer.Container
    Marketplace    *marketplacecontainer.Container
    Analytics      *analyticscontainer.Container
    Settings       *settingscontainer.Container
    Theme          *themecontainer.Container
    Plugin         *plugincontainer.Container
    Shipping       *shippingcontainer.Container
    Tax            *taxcontainer.Container
}

// Wiring:
func NewContainer(cfg *config.Config) *Container
    → initInfrastructure()    // DB, Redis
    → initModules()           // All domain containers
        → initFileStorage()   // Local or S3
        → initJobx()          // Redis job queue
        → initNotifx()        // Email provider
        → IAM.New()           // Auth system
        → eventbus.NewInMemoryBus()
        → All domain containers (Customer, Order, etc.)
```

### Example Domain Container Wiring
```go
// customer/customercontainer/container.go
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
    repo := customerinfra.NewPostgresRepo(db)
    svc := customersrv.New(repo, bus)
    handler := customerapi.NewHandler(svc)
    return &Container{Service: svc, Handler: handler}
}
```

---

## HTTP Server Setup

### Server Entry Point
**File**: `cmd/server.go`

```go
func main() {
    // 1. Load config from environment
    cfg, err := config.Load()

    // 2. Initialize logger based on log level

    // 3. Create container (DI, initialize all modules)
    container := NewContainer(cfg)

    // 4. Start background services (IAM, Job processor)
    container.StartBackgroundServices(ctx)

    // 5. Create Fiber app with error handler & body limit
    app := fiber.New(fiber.Config{
        AppName:           "hada-commerce API",
        ErrorHandler:      globalErrorHandler(cfg),
        BodyLimit:         10 * 1024 * 1024,  // 10MB
        IdleTimeout:       120,
    })

    // 6. Setup middleware (CORS, logging, request ID, panic recovery)
    setupMiddleware(app, cfg)

    // 7. Register health check & info endpoints
    app.Get("/health", healthCheckHandler)
    app.Get("/", infoHandler)

    // 8. Register all routes
    registerRoutes(app, container)

    // 9. Start server with graceful shutdown
    startServer(app, cfg, cancel)
}
```

### Route Registration
**File**: `cmd/server.go` → `registerRoutes()`

```go
// Public routes (no auth)
public := app.Group("/api/v1")
container.Storefront.Handler.RegisterPublicRoutes(public)
container.Theme.Handler.RegisterPublicRoutes(public)
container.Product.Handler.RegisterPublicRoutes(public)
container.Catalog.Handler.RegisterPublicRoutes(public)
container.Settings.Handler.RegisterPublicRoutes(public)
container.Cart.Handler.RegisterPublicRoutes(public)

// Protected routes (require auth via UnifiedAuthMiddleware)
protected := app.Group("/api/v1",
    container.IAM.UnifiedAuthMiddleware.Authenticate(),
)
container.Cart.RegisterRoutes(protected)
container.Product.RegisterRoutes(protected)
container.Order.RegisterRoutes(protected)
container.Payment.RegisterRoutes(protected)
container.Customer.RegisterRoutes(protected)
container.Catalog.RegisterRoutes(protected)
// ... all other protected domains
```

### Middleware Setup
**File**: `cmd/server.go` → `setupMiddleware()`

```
1. recover.New()           // Panic recovery
2. requestid.New()         // X-Request-ID header
3. cors.New()              // CORS headers
4. logger.New()            // Request logging
```

---

## Database Migrations

**Location**: `backend/migrations/`

Available migrations (sequential):
1. `001_initial.up.sql` - Core tables (customers, users, tenants, etc.)
2. `002_marketplace.up.sql` - Marketplace tables
3. `003_settings.up.sql` - Settings
4. `004_iam.up.sql` - IAM tables (users, api_keys, sessions, etc.)
5. `005_blocks.up.sql` - Page blocks
6. `006_themes.up.sql` - Theme system
7. `007_seed_block_types.up.sql` - Seed block types
8. `008_navigation_menus.up.sql` - Navigation
9. `009_template_overrides.up.sql` - Template overrides
10. `010_cart.up.sql` - Shopping cart tables
11. `011_shipping.up.sql` - Shipping zones and rates
12. `012_tax.up.sql` - Tax rates
13. `013_payment.up.sql` - Payment records

**Database tables** (partial list):
- `customers` (id, tenant_id, email, name, phone, addresses JSON, created_at, updated_at)
- `orders` (id, tenant_id, customer_id, status, items JSON, total_amount, etc.)
- `users` (IAM system users, not customers)
- `api_keys` (for programmatic access)
- `sessions` (user sessions)
- `refresh_tokens` (for token rotation)
- `password_reset_tokens` (for password resets)
- `products`, `categories`, `collections`, `cart_items`, etc.

---

## Error Handling Pattern

The codebase uses an `errx` package for structured error handling:

```go
// Definition (e.g., customer/errors.go)
var ErrNotFound = errx.New("customer not found", errx.TypeNotFound)

// Usage in service layer:
if err != nil && !errx.Is(err, customer.ErrNotFound) {
    return nil, errx.Wrap(err, "operation context", errx.TypeInternal)
}

// Error types:
// - TypeNotFound      (404)
// - TypeValidation    (400)
// - TypeConflict      (409)
// - TypeAuthorization (401)
// - TypeForbidden     (403)
// - TypeInternal      (500)
// - TypeExternal      (from external service)
```

---

## Event Bus Pattern

**Location**: `internal/eventbus/`

```go
// Event types published by domains:
eventbus.CustomerRegistered  // When customer created
eventbus.CustomerUpdated     // When customer updated
eventbus.OrderPlaced         // When order created
eventbus.OrderConfirmed      // Order status changed to Confirmed
eventbus.OrderShipped        // Order status changed to Shipped
eventbus.OrderDelivered      // Order status changed to Delivered
eventbus.OrderCancelled      // Order cancelled

// Usage in service layer:
if evt, err := eventbus.NewEvent(eventbus.CustomerRegistered, tenantID, eventbus.CustomerPayload{
    CustomerID: string(c.ID),
    Email:      string(c.Email),
    Name:       c.Name,
}); err == nil {
    _ = s.bus.Publish(ctx, evt)
}
```

---

## Key Dependencies (go.mod)

```
github.com/gofiber/fiber/v2              v2.52.10   // Web framework
github.com/golang-jwt/jwt/v5              v5.2.2     // JWT tokens
github.com/google/uuid                    v1.6.0     // UUID generation
github.com/jmoiron/sqlx                   v1.4.0     // Database driver
github.com/lib/pq                         v1.10.9    // PostgreSQL driver
github.com/redis/go-redis/v9              v9.17.2    // Redis client
golang.org/x/crypto                       v0.40.0    // bcrypt for passwords
github.com/aws/aws-sdk-go-v2              v1.41.8    // AWS SDK (S3, SES)
```

---

## Environment Variables (Critical)

### Database
- `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_USER`, `DATABASE_PASSWORD`, `DATABASE_NAME`
- `DATABASE_MAX_OPEN_CONNS`, `DATABASE_MAX_IDLE_CONNS`

### JWT/Auth
- `JWT_SECRET_KEY` (required, min 32 chars)
- `JWT_ACCESS_TOKEN_TTL` (default: 15m)
- `JWT_REFRESH_TOKEN_TTL` (default: 7 days)
- `JWT_ISSUER` (default: "manifesto")
- `JWT_AUDIENCE` (comma-separated, default: "manifesto-api")

### Session/Security
- `SESSION_EXPIRATION_TIME` (default: 24h)
- `SESSION_MAX_PER_USER` (default: 10)
- `BCRYPT_COST` (default: 10, valid: 4-31)

### Cookies
- `COOKIE_ACCESS_TOKEN_NAME`, `COOKIE_REFRESH_TOKEN_NAME`
- `COOKIE_SECURE`, `COOKIE_HTTP_ONLY`, `COOKIE_SAME_SITE`

### File Storage
- `STORAGE_MODE` ("local" or "s3")
- If S3: `AWS_REGION`, `AWS_BUCKET`

### Email Notifications
- `NOTIFX_PROVIDER` ("ses" or defaults to console)
- If SES: `NOTIFX_FROM_ADDRESS`, `NOTIFX_AWS_REGION`

### Server
- `SERVER_PORT` (default: 3000)
- `SERVER_ENVIRONMENT` ("development", "staging", "production")
- `SERVER_LOG_LEVEL` ("debug", "info", "warn", "error")
- `SERVER_CORS_ORIGINS` (comma-separated, default: "*")

---

## Testing & Validation

### Build & Compile
```bash
cd backend
go build ./...   # Compile all packages
go vet ./...     # Run vet checks
```

### Test
```bash
go test ./...    # Run all tests
```

### Database Migration
```bash
# Migrations are auto-applied or manual via migration tool
# Sequential order enforced by numbered filenames
```

---

## Summary of Key Patterns

1. **DDD Domain Structure**: Each domain (Customer, Order, Product) follows Entity → Repository Interface → Service → Container → HTTP Handler pattern

2. **Multi-Tenancy**: All entities include `TenantID` field; all queries filter by tenant_id

3. **Error Handling**: Centralized `errx` package with typed errors and HTTP status mapping

4. **Event Publishing**: Domain services publish events to EventBus; listeners can subscribe

5. **Configuration**: All config loaded from environment variables via `internal/config` package

6. **Authentication**: JWT tokens + API Keys, both validated via `UnifiedAuthMiddleware`

7. **Dependency Injection**: Main container wires all modules; each domain has its own container

8. **Fiber Web Framework**: Lightweight Go HTTP framework; routes grouped by auth level (public vs protected)

9. **PostgreSQL + JSON**: Nested data (e.g., Addresses) stored as JSON; unmarshaled into Go structs

10. **Request Context**: AuthContext injected via Fiber locals (`c.Locals("auth")`); contains user, tenant, scopes

---

## Next Steps for Implementation

When adding new features:

1. **Define Entity** in `internal/{domain}/{domain}.go`
2. **Define Repository Interface** in `internal/{domain}/port.go`
3. **Implement Service** in `internal/{domain}/{domain}srv/service.go`
4. **Implement Repository** in `internal/{domain}/{domain}infra/postgres.go`
5. **Create HTTP Handler** in `internal/{domain}/{domain}api/handler.go`
6. **Wire Container** in `internal/{domain}/{domain}container/container.go`
7. **Add Routes** in `cmd/server.go` (public or protected group)
8. **Add DB Migration** in `backend/migrations/`
9. **Publish Events** if domain events needed via `eventbus.Publish()`

