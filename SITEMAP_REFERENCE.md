# Sitemap.xml Endpoint - Complete Reference

## Module Information
- **Go Module**: `github.com/Abraxas-365/hada-commerce`
- **Go Version**: 1.25.4
- **Framework**: Fiber v2
- **Database**: PostgreSQL (via sqlx)
- **Port File Path**: `backend/internal/*/port.go`

---

## 1. Container Structure

### Main Container (`backend/cmd/container.go`)
```go
type Container struct {
    Config *config.Config

    // Infrastructure
    DB    *sqlx.DB
    Redis *redis.Client

    // File storage
    FileSystem fsx.FileSystem
    S3Client   *s3.Client

    // Background services
    JobClient    *jobx.Client
    NotifxClient *notifx.Client

    // Event bus
    EventBus eventbus.Bus

    // IAM
    IAM *iamcontainer.Container

    // Commerce domains
    Cart        *cartcontainer.Container
    Product     *productcontainer.Container
    Order       *ordercontainer.Container
    Payment     *paymentcontainer.Container
    Customer    *customercontainer.Container
    Catalog     *catalogcontainer.Container
    Storefront  *storefrontcontainer.Container
    Promo       *promocontainer.Container
    Media       *mediacontainer.Container
    Marketplace *marketplacecontainer.Container
    Analytics   *analyticscontainer.Container
    Settings    *settingscontainer.Container
    Theme       *themecontainer.Container
    Plugin      *plugincontainer.Container
    Search      *searchcontainer.Container
    Shipping    *shippingcontainer.Container
    Tax         *taxcontainer.Container
    Checkout    *checkoutcontainer.Container
}
```

### Example Domain Container (`backend/internal/product/productcontainer/container.go`)
```go
type Container struct {
    Service *productsrv.Service
    Handler *productapi.Handler
}

// New creates a fully-wired product container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
    repo := productinfra.NewPostgresRepo(db)
    variantRepo := productinfra.NewVariantPostgresRepo(db)
    svc := productsrv.New(repo, variantRepo, bus)
    handler := productapi.NewHandler(svc)
    return &Container{
        Service: svc,
        Handler: handler,
    }
}

// RegisterRoutes registers product HTTP routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
    c.Handler.RegisterRoutes(router)
}
```

---

## 2. Route Registration Pattern (`backend/cmd/server.go`)

### Public Routes (No Authentication)
```go
// Public storefront routes (no auth)
public := app.Group("/api/v1")
container.Storefront.Handler.RegisterPublicRoutes(public)
container.Theme.Handler.RegisterPublicRoutes(public)
container.Product.Handler.RegisterPublicRoutes(public)
container.Catalog.Handler.RegisterPublicRoutes(public)
container.Settings.Handler.RegisterPublicRoutes(public)
container.Plugin.Handler.RegisterPublicRoutes(public)
container.Cart.Handler.RegisterPublicRoutes(public)
container.Search.Handler.RegisterPublicRoutes(public)
container.Shipping.Handler.RegisterPublicRoutes(public)
container.Tax.Handler.RegisterPublicRoutes(public)
container.Checkout.RegisterPublicRoutes(public)
container.Customer.RegisterPublicRoutes(public)
```

### Protected Routes (Authentication Required)
```go
protected := app.Group("/api/v1",
    container.IAM.UnifiedAuthMiddleware.Authenticate(),
)
container.Product.RegisterRoutes(protected)
// ... other protected routes
```

### Tenant Identification
- **For protected routes**: Extracted from `AuthContext` in Fiber locals
- **For public routes**: From `X-Tenant-ID` header

---

## 3. Handler Pattern (`backend/internal/product/productapi/handler.go`)

### Handler Structure
```go
type Handler struct {
    svc *productsrv.Service
}

func NewHandler(svc *productsrv.Service) *Handler {
    return &Handler{svc: svc}
}

// RegisterPublicRoutes registers unauthenticated, read-only product routes.
// Tenant is identified via the X-Tenant-ID header.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
    g := router.Group("/products")
    g.Get("/", h.ListProductsPublic)
    g.Get("/:id", h.GetProductPublic)
    // ...
}
```

### Public Handler Pattern
```go
// ListProductsPublic handles GET /products (public, no auth).
// Accepts optional query params: page, page_size, category_id.
func (h *Handler) ListProductsPublic(c *fiber.Ctx) error {
    tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
    if tenantID == "" {
        return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
    }

    pg := paginationFromQuery(c)
    categoryIDStr := c.Query("category_id")

    if categoryIDStr != "" {
        result, err := h.svc.ListByCategory(c.Context(), tenantID, kernel.CategoryID(categoryIDStr), pg)
        if err != nil {
            return err
        }
        return c.JSON(result)
    }

    result, err := h.svc.List(c.Context(), tenantID, pg)
    if err != nil {
        return err
    }
    return c.JSON(result)
}
```

### Pagination Helper
```go
func paginationFromQuery(c *fiber.Ctx) kernel.PaginationOptions {
    page, _ := strconv.Atoi(c.Query("page"))
    pageSize, _ := strconv.Atoi(c.Query("page_size"))
    opts := kernel.PaginationOptions{Page: page, PageSize: pageSize}
    if opts.Page < 1 {
        opts.Page = 1
    }
    if opts.PageSize < 1 {
        opts.PageSize = 20
    }
    if opts.PageSize > 100 {
        opts.PageSize = 100
    }
    return opts
}
```

---

## 4. Kernel Types

### Authentication Context (`backend/internal/kernel/context.go`)
```go
type AuthContext struct {
    UserID   *UserID  `json:"user_id"`
    TenantID TenantID `json:"tenant_id"`
    Email    string   `json:"email"`
    Name     string   `json:"name"`
    Scopes   []string `json:"scopes"`
    IsAPIKey bool     `json:"is_api_key"`
}

// Extract from Fiber locals in protected routes:
authCtx := c.Locals("auth").(*kernel.AuthContext)
tenantID := authCtx.TenantID
```

### ID Types (`backend/internal/kernel/common_ids.go`, `proj_ids.go`)
```go
type UserID string
type TenantID string
type ProductID string
type OrderID string
type CustomerID string
type CategoryID string
type CollectionID string
type PageID string
type VariantID string
type OptionID string

// All follow same pattern:
func NewProductID(id string) ProductID { return ProductID(id) }
func (p ProductID) String() string     { return string(p) }
func (p ProductID) IsEmpty() bool      { return string(p) == "" }
```

### Pagination Types (`backend/internal/kernel/pagination.go`)
```go
type PaginationOptions struct {
    Page     int
    PageSize int
}

func (p PaginationOptions) Offset() int {
    return (p.Page - 1) * p.PageSize
}

func (p PaginationOptions) Limit() int {
    return p.PageSize
}

type Paginated[T any] struct {
    Items      []T `json:"items"`
    Total      int `json:"total"`
    Page       int `json:"page"`
    PageSize   int `json:"page_size"`
    TotalPages int `json:"total_pages"`
}

func NewPaginated[T any](items []T, page, pageSize, total int) Paginated[T] {
    // ... implementation
}
```

---

## 5. Domain Models

### Product (`backend/internal/product/product.go`)
```go
type Status string
const (
    StatusDraft    Status = "draft"
    StatusActive   Status = "active"
    StatusArchived Status = "archived"
)

type Product struct {
    ID          kernel.ProductID  `json:"id" db:"id"`
    TenantID    kernel.TenantID   `json:"tenant_id" db:"tenant_id"`
    Name        string            `json:"name" db:"name"`
    Description string            `json:"description" db:"description"`
    Price       kernel.Money      `json:"price"`
    SKU         string            `json:"sku" db:"sku"`
    Images      []string          `json:"images"`
    CategoryID  kernel.CategoryID `json:"category_id" db:"category_id"`
    Tags        []string          `json:"tags"`
    Status      Status            `json:"status" db:"status"`
    Stock       int               `json:"stock" db:"stock"`
    HasVariants bool              `json:"has_variants" db:"has_variants"`
    Options     []ProductOption   `json:"options,omitempty"`
    Variants    []ProductVariant  `json:"variants,omitempty"`
    CreatedAt   time.Time         `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
}
```

### Category (`backend/internal/catalog/catalog.go`)
```go
type Category struct {
    ID          kernel.CategoryID  `json:"id" db:"id"`
    TenantID    kernel.TenantID    `json:"tenant_id" db:"tenant_id"`
    Name        string             `json:"name" db:"name"`
    Slug        string             `json:"slug" db:"slug"`
    ParentID    *kernel.CategoryID `json:"parent_id" db:"parent_id"`
    Description string             `json:"description" db:"description"`
    CreatedAt   time.Time          `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time          `json:"updated_at" db:"updated_at"`
}

func (c *Category) IsRoot() bool {
    return c.ParentID == nil
}
```

### Collection (`backend/internal/catalog/catalog.go`)
```go
type Collection struct {
    ID          kernel.CollectionID `json:"id" db:"id"`
    TenantID    kernel.TenantID     `json:"tenant_id" db:"tenant_id"`
    Name        string              `json:"name" db:"name"`
    Slug        string              `json:"slug" db:"slug"`
    Description string              `json:"description" db:"description"`
    ProductIDs  []kernel.ProductID  `json:"product_ids" db:"product_ids"`
    IsAutomatic bool                `json:"is_automatic" db:"is_automatic"`
    Rules       map[string]any      `json:"rules" db:"rules"`
    CreatedAt   time.Time           `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time           `json:"updated_at" db:"updated_at"`
}
```

### Page (Storefront) (`backend/internal/storefront/storefront.go`)
```go
type PageStatus string
const (
    PageStatusDraft         PageStatus = "draft"
    PageStatusPendingReview PageStatus = "pending_review"
    PageStatusPublished     PageStatus = "published"
    PageStatusArchived      PageStatus = "archived"
)

type PageMeta struct {
    Description string   `json:"description"`
    OGTitle     string   `json:"og_title"`
    OGImage     string   `json:"og_image"`
    Keywords    []string `json:"keywords"`
}

type Page struct {
    ID          kernel.PageID   `json:"id" db:"id"`
    TenantID    kernel.TenantID `json:"tenant_id" db:"tenant_id"`
    Slug        string          `json:"slug" db:"slug"`
    Title       string          `json:"title" db:"title"`
    HTML        string          `json:"html" db:"html"`
    CSS         string          `json:"css" db:"css"`
    Meta        PageMeta        `json:"meta" db:"meta"`
    ContentType ContentType     `json:"content_type" db:"content_type"`
    Sections    []Section       `json:"sections"`
    Status      PageStatus      `json:"status" db:"status"`
    Version     int             `json:"version" db:"version"`
    CreatedBy   string          `json:"created_by" db:"created_by"`
    PublishedAt *time.Time      `json:"published_at,omitempty" db:"published_at"`
    CreatedAt   time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

func (p *Page) IsPublished() bool {
    return p.Status == PageStatusPublished
}
```

### Store Settings (`backend/internal/settings/settings.go`)
```go
type StoreSettings struct {
    TenantID       kernel.TenantID `json:"tenant_id" db:"tenant_id"`
    StoreName      string          `json:"store_name" db:"store_name"`
    StoreEmail     string          `json:"store_email" db:"store_email"`
    StorePhone     string          `json:"store_phone" db:"store_phone"`
    Currency       string          `json:"currency" db:"currency"`
    Timezone       string          `json:"timezone" db:"timezone"`
    Address        StoreAddress    `json:"address" db:"address"`
    LogoURL        string          `json:"logo_url" db:"logo_url"`
    FaviconURL     string          `json:"favicon_url" db:"favicon_url"`
    SocialLinks    SocialLinks     `json:"social_links" db:"social_links"`
    CheckoutConfig CheckoutConfig  `json:"checkout_config" db:"checkout_config"`
    UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
}
```

---

## 6. Port (Repository) Interfaces

### Product Repository (`backend/internal/product/port.go`)
```go
type Repository interface {
    Create(ctx context.Context, p *Product) error
    GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ProductID) (*Product, error)
    Update(ctx context.Context, p *Product) error
    Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.ProductID) error
    List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[Product], error)
    ListByCategory(ctx context.Context, tenantID kernel.TenantID, categoryID kernel.CategoryID, pg kernel.PaginationOptions) (kernel.Paginated[Product], error)
    GetBySKU(ctx context.Context, tenantID kernel.TenantID, sku string) (*Product, error)
}
```

### Catalog Repository (`backend/internal/catalog/port.go`)
```go
type CategoryRepository interface {
    Create(ctx context.Context, c *Category) error
    GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CategoryID) (*Category, error)
    GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*Category, error)
    Update(ctx context.Context, c *Category) error
    Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CategoryID) error
    List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[Category], error)
    ListByParent(ctx context.Context, tenantID kernel.TenantID, parentID *kernel.CategoryID, pg kernel.PaginationOptions) (kernel.Paginated[Category], error)
}

type CollectionRepository interface {
    Create(ctx context.Context, c *Collection) error
    GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CollectionID) (*Collection, error)
    GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*Collection, error)
    Update(ctx context.Context, c *Collection) error
    Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CollectionID) error
    List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[Collection], error)
}
```

### Storefront Repository (`backend/internal/storefront/port.go`)
```go
type PageRepository interface {
    Create(ctx context.Context, page *Page) error
    GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.PageID) (*Page, error)
    GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*Page, error)
    GetPublished(ctx context.Context, tenantID kernel.TenantID, slug string) (*Page, error)
    Update(ctx context.Context, page *Page) error
    ListByStatus(ctx context.Context, tenantID kernel.TenantID, status PageStatus, p kernel.PaginationOptions) (kernel.Paginated[Page], error)
    List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[Page], error)
}
```

---

## 7. Fiber Context Access Patterns

### Get Tenant from Protected Route (with auth)
```go
authCtx := c.Locals("auth").(*kernel.AuthContext)
tenantID := authCtx.TenantID
```

### Get Tenant from Public Route (header-based)
```go
tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
if tenantID == "" {
    return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
}
```

### Extract Query Parameters
```go
page, _ := strconv.Atoi(c.Query("page"))
pageSize, _ := strconv.Atoi(c.Query("page_size"))
```

### Return JSON Response
```go
return c.JSON(result)
return c.Status(fiber.StatusCreated).JSON(result)
return c.SendStatus(fiber.StatusNoContent)
```

---

## 8. Available Domain Services (via Container)

For your sitemap endpoint, these containers and their services are available:

- **container.Product.Service** - ProductRepository interface
- **container.Catalog.Service** - CategoryRepository & CollectionRepository interfaces
- **container.Storefront.Service** - PageRepository interface
- **container.Settings.Service** - StoreSettings access
- **container.DB** - Direct database access (*sqlx.DB) if needed

---

## 9. Key Notes for Sitemap Implementation

### Multi-tenancy
- Every request requires tenant identification
- Filter all data by TenantID before including in sitemap
- Use X-Tenant-ID header for public endpoints

### Status Filtering
- Only include **published** pages (PageStatus == PageStatusPublished)
- Only include **active** products (Status == StatusActive)
- Only include **published** pages from storefront

### URL Slug Generation
- Products: Use ProductID or slug if available
- Categories: Use slug field from Category
- Collections: Use slug field from Collection
- Pages: Use slug field from Page

### Pagination
- Use kernel.PaginationOptions for large result sets
- Max PageSize is 100 (enforced in pagination helper)

### Error Handling
- Use errx.New() for validation errors (returns JSON with proper HTTP status)
- Return proper HTTP status codes:
  - 200 OK for successful GET
  - 400 Bad Request for missing X-Tenant-ID
  - 404 Not Found if tenant has no data

---

## 10. Example: Complete Handler Flow

```go
func (h *Handler) GetSitemap(c *fiber.Ctx) error {
    // 1. Get tenant from header (public endpoint)
    tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
    if tenantID == "" {
        return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
    }

    // 2. Get pagination (optional)
    pg := kernel.NewPaginationOptions(1, 100)

    // 3. Fetch data from service
    products, err := h.productSvc.List(c.Context(), tenantID, pg)
    if err != nil {
        return err
    }

    // 4. Transform to XML/JSON
    // ... transform logic ...

    // 5. Return response
    c.Set("Content-Type", "application/xml")
    return c.Send(xmlBytes)
}
```

---

## Summary

This reference provides all the patterns, types, and interfaces needed to:
- Create a sitemap endpoint handler
- Access tenant-scoped data via repositories
- Handle pagination and public access
- Return proper error responses
- Follow the hada-commerce patterns and conventions
