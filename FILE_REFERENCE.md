# File Reference Guide

Quick reference for all analyzed files with their exact contents and purposes.

---

## Customer Domain Files

### `backend/internal/customer/customer.go`
**Purpose**: Define Customer entity and value objects  
**Key Types**:
- `Customer`: Entity with ID, TenantID, Email, Name, Phone, Addresses[], CreatedAt, UpdatedAt
- `Address`: Nested value object (Street, City, State, Country, PostalCode, IsDefault)
**Key Methods**:
- `DefaultAddress() *Address` - Returns the address marked as default
- `AddAddress(addr Address)` - Adds address, clears default from others if needed
- `SetDefaultAddress(idx int) bool` - Sets address at index as default

### `backend/internal/customer/port.go`
**Purpose**: Define Repository interface contract  
**Key Interface**:
- `Repository`: CRUD + List operations scoped by tenant
  - `Create(ctx, c *Customer) error`
  - `GetByID(ctx, tenantID, id) (*Customer, error)`
  - `GetByEmail(ctx, tenantID, email) (*Customer, error)`
  - `Update(ctx, c *Customer) error`
  - `Delete(ctx, tenantID, id) error`
  - `List(ctx, tenantID, pg) (Paginated[Customer], error)`

### `backend/internal/customer/errors.go`
**Purpose**: Domain-specific error definitions  
**Errors**:
- `ErrNotFound` (TypeNotFound)
- `ErrDuplicateEmail` (TypeConflict)
- `ErrInvalidEmail` (TypeValidation)

### `backend/internal/customer/customersrv/service.go`
**Purpose**: Business logic for customer operations  
**Key Type**:
- `Service`: Holds Repository + EventBus
- `CreateInput`: DTO for creating customers
**Key Methods**:
- `Create(ctx, tenantID, input) → (*Customer, error)` - Creates customer, checks email uniqueness, publishes CustomerRegistered event
- `GetByID(ctx, tenantID, id) → (*Customer, error)`
- `GetByEmail(ctx, tenantID, email) → (*Customer, error)`
- `Update(ctx, customer) error` - Updates and publishes CustomerUpdated event
- `Delete(ctx, tenantID, id) error`
- `List(ctx, tenantID, pg) → (Paginated[Customer], error)`
**Events Published**:
- `eventbus.CustomerRegistered` (on Create)
- `eventbus.CustomerUpdated` (on Update)

### `backend/internal/customer/customerinfra/postgres.go`
**Purpose**: PostgreSQL implementation of Customer Repository  
**Key Type**:
- `PostgresRepo`: Wraps *sqlx.DB
**Implementation Details**:
- Addresses stored as JSON string in DB, marshaled/unmarshaled
- All queries scoped by `tenant_id` (WHERE tenant_id = $N)
- Uses scanner pattern for Row/Rows polymorphism
- Error handling via `errx.Wrap()`
**Tables Used**:
- `customers` (id, tenant_id, email, name, phone, addresses, created_at, updated_at)

### `backend/internal/customer/customercontainer/container.go`
**Purpose**: Dependency injection container for Customer domain  
**Key Type**:
- `Container`: Holds Service + Handler
**Wiring**:
1. Create PostgresRepo(db)
2. Create Service(repo, bus)
3. Create Handler(service)
4. Return Container
**Public Method**:
- `RegisterRoutes(router fiber.Router)` - Registers HTTP routes

---

## Order Domain Files

### `backend/internal/order/ordersrv/service.go`
**Purpose**: Business logic for order operations  
**Key Types**:
- `Service`: Holds Repository + EventBus
- `CreateInput`: DTO with CustomerID, Items[], ShippingAddress
- `CreateItemInput`: Product info for order line item
**Key Methods**:
- `Create(ctx, tenantID, input) → (*Order, error)` - Creates order with items, calculates total, publishes OrderPlaced
- `GetByID(ctx, tenantID, id) → (*Order, error)`
- `UpdateStatus(ctx, tenantID, id, newStatus) → (*Order, error)` - Status transitions, publishes OrderConfirmed/Shipped/Delivered
- `Cancel(ctx, tenantID, id) → (*Order, error)` - Cancels order, publishes OrderCancelled
- `List(ctx, tenantID, pg) → (Paginated[Order], error)` - Lists all orders for tenant
- `ListByCustomer(ctx, tenantID, customerID, pg) → (Paginated[Order], error)` - Lists customer's orders
**Events Published**:
- `eventbus.OrderPlaced`, `OrderConfirmed`, `OrderShipped`, `OrderDelivered`, `OrderCancelled`

---

## Authentication (IAM) Domain

### `backend/internal/iam/auth/jwt_service.go`
**Purpose**: Generate and validate JWT tokens  
**Key Type**:
- `JWTService`: Holds secretKey, tokenTTLs, issuer, audience
- `JWTClaims`: Custom JWT claims (UserID, TenantID, Email, Name, Scopes, RegisteredClaims)
**Key Methods**:
- `NewJWTServiceFromConfig(cfg *JWTConfig) *JWTService`
- `GenerateAccessToken(userID, tenantID, claims map[string]any) (string, error)` - Creates access token with 15m TTL by default
- `ValidateAccessToken(tokenString) (*TokenClaims, error)` - Parses and validates, returns TokenClaims
- `GenerateRefreshToken(userID) (string, error)` - Creates refresh token with 7d TTL by default
**Token Config** (from JWTConfig):
- `SecretKey`: HMAC secret (min 32 chars, required)
- `AccessTokenTTL`: Default 15 minutes
- `RefreshTokenTTL`: Default 7 days
- `Issuer`: e.g., "manifesto"
- `Audience`: e.g., ["manifesto-api"]

### `backend/internal/iam/auth/unified_middleware.go`
**Purpose**: Authenticate requests using JWT or API Keys  
**Key Type**:
- `UnifiedAuthMiddleware`: Holds APIKeyService + TokenService
**Key Methods**:
- `Authenticate() fiber.Handler` - Middleware that sets c.Locals("auth") = AuthContext
  - Tries API Key first (extractAPIKey), falls back to JWT
- `authenticateAPIKey(c, keyString) error` - Validates key, sets auth context with IsAPIKey=true
- `authenticateJWT(c) error` - Validates JWT, sets auth context with UserID pointer
- `RequireScope(scope) fiber.Handler` - Middleware to enforce single scope
- `RequireAnyScope(scopes...) fiber.Handler` - Middleware to enforce ANY of scopes
- `RequireAllScopes(scopes...) fiber.Handler` - Middleware to enforce ALL scopes
- `RequireAdmin() fiber.Handler` - Requires "*" or "admin:*" scope
- `extractAPIKey(c *fiber.Ctx) string` - Extracts from: Authorization header, X-API-Key header, api_key query
**Token Extraction Sources**:
1. Authorization: Bearer <token> header
2. Authorization: X-API-Key <key> header
3. X-API-Key header
4. api_key query parameter
5. access_token cookie

### `backend/internal/iam/auth/authinfra/bcrypt_password_service.go`
**Purpose**: Hash and verify passwords using bcrypt  
**Key Type**:
- `BcryptPasswordService`: Holds cost factor (default 10)
**Key Methods**:
- `NewBcryptPasswordService(cost int) user.PasswordService`
- `HashPassword(password) (string, error)` - Returns bcrypt hash
- `VerifyPassword(hashedPassword, password) bool` - Returns true if valid

### `backend/internal/iam/auth/port.go`
**Purpose**: Define token/session/password reset interfaces  
**Key Types**:
- `RefreshToken`: Stored token with UserID, TenantID, ExpiresAt, IsRevoked
- `UserSession`: Session with SessionToken, IPAddress, UserAgent, ExpiresAt, LastActivity
- `PasswordResetToken`: Reset token with UserID, ExpiresAt, IsUsed flag
- `TokenClaims`: Claims returned after validation (UserID, TenantID, Email, Name, Scopes, IssuedAt, ExpiresAt)
**Key Interfaces**:
- `TokenRepository`: SaveRefreshToken, FindRefreshToken, RevokeRefreshToken, RevokeAllUserTokens, CleanExpiredTokens
- `SessionRepository`: SaveSession, FindSession, FindUserSessions, UpdateSessionActivity, RevokeSession, RevokeAllUserSessions, CleanExpiredSessions
- `PasswordResetRepository`: SaveResetToken, FindResetToken, ConsumeResetToken, CleanExpiredResetTokens
- `TokenService`: GenerateAccessToken, ValidateAccessToken, GenerateRefreshToken (implemented by JWTService)
- `AuditService`: LogLoginAttempt, LogLogout, LogTokenRefresh, LogOTPVerification, LogAccountCreated, LogAccountLinked

### `backend/internal/iam/auth/auth.go`
**Purpose**: Error definitions for auth domain  
**Error Types** (via ErrRegistry):
- `INVALID_REFRESH_TOKEN`, `EXPIRED_REFRESH_TOKEN`
- `INVALID_OAUTH_PROVIDER`, `OAUTH_AUTHORIZATION_FAILED`, `INVALID_STATE`
- `TOKEN_GENERATION_FAILED`, `TOKEN_VALIDATION_FAILED`
- `OAUTH_CALLBACK_ERROR`

---

## Kernel (Core Types)

### `backend/internal/kernel/common_ids.go`
**Purpose**: Define UserID and TenantID types  
**Key Types**:
- `UserID string` - User in IAM system (methods: NewUserID, String, IsEmpty)
- `TenantID string` - Organization (methods: NewTenantID, String, IsEmpty)

### `backend/internal/kernel/proj_ids.go`
**Purpose**: Define commerce domain IDs  
**Key Types** (all string-based, with NewX, String, IsEmpty methods):
- `ProductID`, `OrderID`, `OrderItemID`
- `CustomerID`, `CategoryID`, `CollectionID`
- `PageID`, `PageVersionID`
- `PromoID`, `MediaID`
- `PluginID`, `PluginVersionID`, `InstallationID`
- `SettingID`, `VendorID`
- `BlockTypeID`, `BlockID`, `ThemeID`
- `CartID`, `CartItemID`
- `ShippingZoneID`, `ShippingRateID`
- `TaxRateID`, `PaymentID`, `RefundID`

### `backend/internal/kernel/context.go`
**Purpose**: Auth context and context keys  
**Key Type**:
- `AuthContext`:
  - `UserID *UserID` - Pointer (null for API keys)
  - `TenantID TenantID` - Always required
  - `Email string` - From JWT claims
  - `Name string` - From JWT claims
  - `Scopes []string` - Authorization scopes
  - `IsAPIKey bool` - True if API key auth
  - Methods: `IsValid()`, `HasScope(scope)`, `HasAnyScope(scopes...)`, `HasAllScopes(scopes...)`, `IsAdmin()`
**Context Keys**:
- `AuthContextKey` - For storing AuthContext in context.Context
- `TenantContextKey`, `UserContextKey`, `RequestIDKey`

### `backend/internal/kernel/money.go`
**Purpose**: Money type for prices and amounts  
**Key Type**:
- `Money`:
  - `Amount int64` - In smallest currency unit (cents)
  - `Currency string` - ISO 4217 code (USD, EUR, etc.)
  - Methods: `NewMoney(cents, currency)`, `Add(other)`, `Multiply(qty)`

### `backend/internal/kernel/pagination.go`
**Purpose**: Pagination types  
**Key Types**:
- `PaginationOptions`: Standard pagination params (Page, PageSize)
  - Methods: `NewPaginationOptions(page, pageSize)`, `Offset()`, `Limit()`
- `Paginated[T]`: Response wrapper
  - Fields: `Items []T`, `Total int`, `Page int`, `PageSize int`, `TotalPages int`
  - Constructor: `NewPaginated[T](items, page, pageSize, total)`

### `backend/internal/kernel/common_objvalue.go`
**Purpose**: Common value objects  
**Key Types**:
- `Email string` - Email type with methods: NewEmail, String, IsEmpty
- `Phone string` - Phone type
- `FirstName string` - First name type
- `LastName string` - Last name type

---

## Configuration

### `backend/internal/config/config.go`
**Purpose**: Main configuration loader  
**Key Type**:
- `Config`:
  - `Server`, `Database`, `Redis`: Infrastructure config
  - `Environment`: development/staging/production
  - `Auth`: Authentication config
  - `OAuth`, `TenantConfig`, `Jobx`, `Notifx`: Other configs
**Key Function**:
- `Load() (*Config, error)` - Loads all config from environment
**Environment Types**:
- `EnvironmentDevelopment`, `EnvironmentStaging`, `EnvironmentProduction`
**Methods**:
- `IsDevelopment()`, `IsStaging()`, `IsProd()`

### `backend/internal/config/auth.go`
**Purpose**: Authentication-specific configuration  
**Key Types**:
- `AuthConfig`: Top-level auth config
  - Fields: JWT, APIKey, Session, OTP, Invitation, PasswordReset, Cookie, Password
- `JWTConfig`:
  - `SecretKey string` (required, min 32 chars)
  - `AccessTokenTTL` (default: 15 min, env: JWT_ACCESS_TOKEN_TTL)
  - `RefreshTokenTTL` (default: 7 days, env: JWT_REFRESH_TOKEN_TTL)
  - `Issuer` (default: "manifesto", env: JWT_ISSUER)
  - `Audience` (default: ["manifesto-api"], env: JWT_AUDIENCE)
- `APIKeyConfig`: LivePrefix, TestPrefix, TokenLength
- `SessionConfig`: ExpirationTime, CleanupInterval, MaxSessions
- `OTPConfig`: CodeLength, ExpirationTime, MaxAttempts, RateLimitWindow, TokenByteLength
- `InvitationConfig`: DefaultExpirationDays, TokenByteLength, MaxPendingPerTenant
- `PasswordResetConfig`: TokenByteLength, ExpirationTime, RateLimitWindow, MaxAttemptsPerWindow
- `CookieConfig`: AccessTokenName, RefreshTokenName, Domain, Path, Secure, HTTPOnly, SameSite
- `PasswordConfig`: BcryptCost (default: 10)
**Function**:
- `loadAuthConfig() AuthConfig` - Loads from environment variables

---

## Main Server

### `backend/cmd/container.go`
**Purpose**: Dependency injection container for entire application  
**Key Type**:
- `Container`:
  - **Infrastructure**: Config, DB (*sqlx.DB), Redis, FileSystem, S3Client, JobClient, NotifxClient, EventBus
  - **Modules**: IAM, Cart, Product, Order, Payment, Customer, Catalog, Storefront, Promo, Media, Marketplace, Analytics, Settings, Theme, Plugin, Shipping, Tax (each a *{domain}container.Container)
**Key Functions**:
- `NewContainer(cfg) *Container` - Initializes all infrastructure and modules
  - `initInfrastructure()` - DB, Redis
  - `initModules()` - Calls initFileStorage, initJobx, initNotifx, creates IAM & all domain containers
  - `initFileStorage()` - Local or S3 file storage
  - `initJobx()` - Redis job queue
  - `initNotifx()` - Email provider (SES or console)
- `StartBackgroundServices(ctx)` - Starts Job processor and IAM background services
- `Cleanup()` - Closes DB and Redis connections
**Adapters**:
- `NotifxOTPNotifier` - Sends OTP via email
- `NotifxInvitationNotifier` - Sends invitations via email

### `backend/cmd/server.go`
**Purpose**: HTTP server setup and route registration  
**Main Flow**:
1. Load configuration
2. Initialize logger
3. Create DI container
4. Start background services
5. Create Fiber app with error handler & body limit (10MB)
6. Setup middleware (recover, request ID, CORS, logging)
7. Register health check & info endpoints
8. Register all routes (public & protected)
9. Start server with graceful shutdown
**Public Routes** (no auth):
- `/api/v1/*` - Storefront, Theme, Product, Catalog, Settings, Cart, Shipping, Tax
**Protected Routes** (require UnifiedAuthMiddleware.Authenticate()):
- `/api/v1/*` - All commerce domains (Cart, Product, Order, Payment, Customer, etc.)
- `/api/v1/*` - API keys, Invitations
**Middleware Stack**:
1. recover.New() - Panic recovery
2. requestid.New() - X-Request-ID header
3. cors.New() - CORS headers
4. logger.New() - Request logging
**Key Handlers**:
- `healthCheckHandler()` - `/health` - Checks DB and Redis
- `infoHandler()` - `/` - API info
- `notFoundHandler()` - 404 responses
- Global error handler - Centralized error response formatting

---

## Database Migrations

**Location**: `backend/migrations/`

**Available migrations** (sequential by number):

1. **001_initial.up.sql** - Core tables
   - Tenants, Users, Customers, Products, Categories, etc.

2. **002_marketplace.up.sql** - Marketplace feature
   - Vendors, Marketplace settings

3. **003_settings.up.sql** - Store settings
   - Global and tenant-level settings

4. **004_iam.up.sql** - Identity & Access Management (13KB)
   - Users, Roles, Permissions, API Keys, Sessions, Refresh Tokens, Password Reset Tokens

5. **005_blocks.up.sql** - Page blocks
   - Block content storage

6. **006_themes.up.sql** - Theme system
   - Theme configurations

7. **007_seed_block_types.up.sql** - Block type seeds (8KB)
   - Bootstrap common block types

8. **008_navigation_menus.up.sql** - Navigation menus
   - Menu structures

9. **009_template_overrides.up.sql** - Custom templates
   - Store-specific template overrides

10. **010_cart.up.sql** - Shopping cart
    - Cart and cart items tables

11. **011_shipping.up.sql** - Shipping
    - Shipping zones and rates

12. **012_tax.up.sql** - Taxes
    - Tax rates and rules

13. **013_payment.up.sql** - Payments
    - Payment records and transaction tracking

**Key tables**:
- `customers` (id, tenant_id, email, name, phone, addresses JSON, created_at, updated_at)
- `orders` (id, tenant_id, customer_id, items JSON, status, total_amount, shipping_address, created_at, updated_at)
- `users` (id, tenant_id, email, password_hash, created_at, updated_at)
- `api_keys` (id, key, user_id, tenant_id, scopes, created_at, expires_at)
- `sessions` (id, user_id, tenant_id, session_token, ip_address, user_agent, created_at, expires_at, last_activity)
- `refresh_tokens` (id, token, user_id, tenant_id, created_at, expires_at, is_revoked)
- `password_reset_tokens` (id, token, user_id, created_at, expires_at, is_used)

---

## Dependencies Summary

**Key Libraries**:
- `github.com/gofiber/fiber/v2` - Web framework
- `github.com/golang-jwt/jwt/v5` - JWT tokens
- `github.com/jmoiron/sqlx` - Database ORM/driver
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/redis/go-redis/v9` - Redis client
- `golang.org/x/crypto` - Bcrypt & crypto utilities
- `github.com/aws/aws-sdk-go-v2` - AWS SDK (S3, SES)
- `github.com/google/uuid` - UUID generation

**Go Version**: 1.25.4

---

## Environment Variables Checklist

### Critical (Required)
- [ ] `JWT_SECRET_KEY` - Must be ≥32 chars

### Database
- [ ] `DATABASE_HOST` - Default: localhost
- [ ] `DATABASE_PORT` - Default: 5432
- [ ] `DATABASE_USER` - Default: hada
- [ ] `DATABASE_PASSWORD`
- [ ] `DATABASE_NAME` - Default: hada
- [ ] `DATABASE_SSL_MODE` - Default: disable

### JWT/Auth
- [ ] `JWT_ISSUER` - Default: manifesto
- [ ] `JWT_AUDIENCE` - Default: manifesto-api (comma-separated)
- [ ] `JWT_ACCESS_TOKEN_TTL` - Default: 15m
- [ ] `JWT_REFRESH_TOKEN_TTL` - Default: 168h (7 days)

### Session
- [ ] `SESSION_EXPIRATION_TIME` - Default: 24h
- [ ] `SESSION_MAX_PER_USER` - Default: 10

### Cookies
- [ ] `COOKIE_SECURE` - Default: false
- [ ] `COOKIE_HTTP_ONLY` - Default: true
- [ ] `COOKIE_SAME_SITE` - Default: Lax

### Server
- [ ] `SERVER_PORT` - Default: 3000
- [ ] `SERVER_ENVIRONMENT` - development/staging/production
- [ ] `SERVER_LOG_LEVEL` - Default: info
- [ ] `SERVER_CORS_ORIGINS` - Default: *

### Storage
- [ ] `STORAGE_MODE` - local or s3
- [ ] `UPLOAD_DIR` - Default: ./uploads (for local)
- [ ] `AWS_REGION`, `AWS_BUCKET` - (for S3)

### Email
- [ ] `NOTIFX_PROVIDER` - ses or console
- [ ] `NOTIFX_FROM_ADDRESS` - (if SES)
- [ ] `NOTIFX_AWS_REGION` - (if SES)

### Redis
- [ ] `REDIS_HOST` - Default: localhost
- [ ] `REDIS_PORT` - Default: 6379
- [ ] `REDIS_DB` - Default: 0

---

## Patterns & Conventions

### Domain Pattern
Each domain follows:
1. **Entity** (`{domain}.go`) - Core business object with value objects
2. **Port** (`port.go`) - Repository interface + other contracts
3. **Service** (`{domain}srv/service.go`) - Business logic, event publishing
4. **Infra** (`{domain}infra/postgres.go`) - Repository implementation
5. **Container** (`{domain}container/container.go`) - DI wiring
6. **API** (`{domain}api/handler.go`) - HTTP handlers (not fully shown)

### Error Pattern
```go
// Domain errors
var ErrNotFound = errx.New("customer not found", errx.TypeNotFound)

// Usage
if err := repo.GetByID(ctx, tenantID, id); err != nil {
    if errx.Is(err, customer.ErrNotFound) {
        return nil, customer.ErrNotFound
    }
    return nil, errx.Wrap(err, "fetching customer", errx.TypeInternal)
}
```

### Multi-Tenancy Pattern
```go
// All operations scoped by tenant
func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerID) (*Customer, error) {
    // WHERE id = $1 AND tenant_id = $2
}
```

### Event Publishing Pattern
```go
// In service after successful operation
if evt, err := eventbus.NewEvent(eventbus.CustomerRegistered, tenantID, eventbus.CustomerPayload{
    CustomerID: string(c.ID),
    Email:      string(c.Email),
    Name:       c.Name,
}); err == nil {
    _ = s.bus.Publish(ctx, evt)
}
```

### Repository Pattern
```go
// Service depends on interface
type Service struct {
    repo customer.Repository  // Interface, not implementation
    bus  eventbus.Bus
}

// Postgres implements it
type PostgresRepo struct {
    db *sqlx.DB
}

var _ customer.Repository = (*PostgresRepo)(nil)  // Compile-time check
```

### Fiber Middleware Pattern
```go
// Public routes
public := app.Group("/api/v1")
public.Get("/products", handler.ListProducts)

// Protected routes
protected := app.Group("/api/v1", 
    container.IAM.UnifiedAuthMiddleware.Authenticate(),
)
protected.Post("/orders", handler.CreateOrder)

// Route-specific scope enforcement
protected.Delete("/customers/:id",
    container.IAM.UnifiedAuthMiddleware.RequireScope("customers:delete"),
    handler.DeleteCustomer,
)
```

---

## Testing Approach

### Database Testing
- Unit tests should mock Repository interface
- Integration tests use real PostgreSQL (test database)
- Migrations auto-applied before tests

### Service Layer Testing
- Mock Repository implementations
- Verify event publishing via mock EventBus
- Test business logic and validation

### HTTP Handler Testing
- Use Fiber's test utilities
- Mock Service layer
- Verify status codes and response JSON

### Example Structure
```go
// customer/customersrv/service_test.go
func TestCreateCustomer(t *testing.T) {
    mockRepo := &mockRepository{}
    mockBus := &mockEventBus{}
    svc := customersrv.New(mockRepo, mockBus)
    
    // Test create with valid input
    // Test duplicate email error
    // Verify event was published
}
```

