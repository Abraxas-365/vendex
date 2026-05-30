# Sitemap.xml Endpoint - Implementation Guide

## Quick Start Overview

The sitemap endpoint should:
1. Accept `X-Tenant-ID` header (public, no auth required)
2. Fetch active products, published pages, and categories from the database
3. Generate an XML sitemap following the [sitemaps.org](https://www.sitemaps.org) protocol
4. Return `application/xml` content type with proper `<urlset>` structure

---

## Step 1: Create the Domain Module Structure

### Directory Layout
```
backend/internal/sitemap/
├── port.go              # Interface definitions
├── sitemap.go           # Domain model
├── error.go             # Error types
├── sitemapapi/
│   └── handler.go       # HTTP handler
├── sitemapinfra/
│   └── (no repo needed - uses existing Product/Catalog/Storefront repos)
├── sitemapsrv/
│   └── service.go       # Business logic
└── sitemapcontainer/
    └── container.go     # Dependency wiring
```

---

## Step 2: Define the Domain Model (`sitemap/sitemap.go`)

```go
package sitemap

import "time"

// Entry represents a single URL entry in the sitemap
type Entry struct {
    URL        string
    LastMod    time.Time
    ChangeFreq string // always, hourly, daily, weekly, monthly, yearly, never
    Priority   float64 // 0.0 to 1.0
}

// Sitemap holds all entries for a tenant
type Sitemap struct {
    TenantID string
    Entries  []Entry
}

// Common change frequencies
const (
    ChangeFreqAlways  = "always"
    ChangeFreqHourly  = "hourly"
    ChangeFreqDaily   = "daily"
    ChangeFreqWeekly  = "weekly"
    ChangeFreqMonthly = "monthly"
    ChangeFreqYearly  = "yearly"
    ChangeFreqNever   = "never"
)

// Default priorities
const (
    PriorityProduct   = 0.8
    PriorityCategory  = 0.7
    PriorityPage      = 0.6
    PriorityHome      = 1.0
)
```

---

## Step 3: Define the Port Interface (`sitemap/port.go`)

```go
package sitemap

import (
    "context"
    "github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Generator interface for building sitemap entries
// Different implementations could generate sitemap entries from different sources
type Generator interface {
    // GenerateEntries returns all sitemap entries for a tenant
    GenerateEntries(ctx context.Context, tenantID kernel.TenantID, baseURL string) ([]Entry, error)
}
```

---

## Step 4: Create the Service (`sitemap/sitemapsrv/service.go`)

```go
package sitemapsrv

import (
    "context"
    "fmt"
    "time"

    "github.com/Abraxas-365/hada-commerce/internal/kernel"
    "github.com/Abraxas-365/hada-commerce/internal/product"
    "github.com/Abraxas-365/hada-commerce/internal/catalog"
    "github.com/Abraxas-365/hada-commerce/internal/storefront"
    "github.com/Abraxas-365/hada-commerce/internal/sitemap"
)

// Service generates sitemaps by aggregating entries from multiple domains
type Service struct {
    productRepo   product.Repository
    categoryRepo  catalog.CategoryRepository
    collectionRepo catalog.CollectionRepository
    pageRepo      storefront.PageRepository
}

// New creates a new sitemap service
func New(
    productRepo product.Repository,
    categoryRepo catalog.CategoryRepository,
    collectionRepo catalog.CollectionRepository,
    pageRepo storefront.PageRepository,
) *Service {
    return &Service{
        productRepo:    productRepo,
        categoryRepo:   categoryRepo,
        collectionRepo: collectionRepo,
        pageRepo:       pageRepo,
    }
}

// Generate creates a complete sitemap for a tenant
func (s *Service) Generate(ctx context.Context, tenantID kernel.TenantID, baseURL string) (*sitemap.Sitemap, error) {
    entries := []sitemap.Entry{}

    // 1. Add homepage
    entries = append(entries, sitemap.Entry{
        URL:        baseURL,
        LastMod:    time.Now(),
        ChangeFreq: sitemap.ChangeFreqDaily,
        Priority:   sitemap.PriorityHome,
    })

    // 2. Add products
    productEntries, err := s.generateProductEntries(ctx, tenantID, baseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to generate product entries: %w", err)
    }
    entries = append(entries, productEntries...)

    // 3. Add categories
    categoryEntries, err := s.generateCategoryEntries(ctx, tenantID, baseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to generate category entries: %w", err)
    }
    entries = append(entries, categoryEntries...)

    // 4. Add collections
    collectionEntries, err := s.generateCollectionEntries(ctx, tenantID, baseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to generate collection entries: %w", err)
    }
    entries = append(entries, collectionEntries...)

    // 5. Add published pages
    pageEntries, err := s.generatePageEntries(ctx, tenantID, baseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to generate page entries: %w", err)
    }
    entries = append(entries, pageEntries...)

    return &sitemap.Sitemap{
        TenantID: tenantID.String(),
        Entries:  entries,
    }, nil
}

// generateProductEntries creates sitemap entries for all active products
func (s *Service) generateProductEntries(ctx context.Context, tenantID kernel.TenantID, baseURL string) ([]sitemap.Entry, error) {
    var entries []sitemap.Entry
    page := 1
    pageSize := 100
    
    for {
        pg := kernel.NewPaginationOptions(page, pageSize)
        result, err := s.productRepo.List(ctx, tenantID, pg)
        if err != nil {
            return nil, err
        }

        for _, p := range result.Items {
            // Only include active products
            if p.Status != product.StatusActive {
                continue
            }

            entries = append(entries, sitemap.Entry{
                URL:        fmt.Sprintf("%s/products/%s", baseURL, p.ID.String()),
                LastMod:    p.UpdatedAt,
                ChangeFreq: sitemap.ChangeFreqWeekly,
                Priority:   sitemap.PriorityProduct,
            })
        }

        if len(result.Items) < pageSize {
            break
        }
        page++
    }

    return entries, nil
}

// generateCategoryEntries creates sitemap entries for all categories
func (s *Service) generateCategoryEntries(ctx context.Context, tenantID kernel.TenantID, baseURL string) ([]sitemap.Entry, error) {
    var entries []sitemap.Entry
    page := 1
    pageSize := 100

    for {
        pg := kernel.NewPaginationOptions(page, pageSize)
        result, err := s.categoryRepo.List(ctx, tenantID, pg)
        if err != nil {
            return nil, err
        }

        for _, c := range result.Items {
            entries = append(entries, sitemap.Entry{
                URL:        fmt.Sprintf("%s/categories/%s", baseURL, c.Slug),
                LastMod:    c.UpdatedAt,
                ChangeFreq: sitemap.ChangeFreqWeekly,
                Priority:   sitemap.PriorityCategory,
            })
        }

        if len(result.Items) < pageSize {
            break
        }
        page++
    }

    return entries, nil
}

// generateCollectionEntries creates sitemap entries for all collections
func (s *Service) generateCollectionEntries(ctx context.Context, tenantID kernel.TenantID, baseURL string) ([]sitemap.Entry, error) {
    var entries []sitemap.Entry
    page := 1
    pageSize := 100

    for {
        pg := kernel.NewPaginationOptions(page, pageSize)
        result, err := s.collectionRepo.List(ctx, tenantID, pg)
        if err != nil {
            return nil, err
        }

        for _, col := range result.Items {
            entries = append(entries, sitemap.Entry{
                URL:        fmt.Sprintf("%s/collections/%s", baseURL, col.Slug),
                LastMod:    col.UpdatedAt,
                ChangeFreq: sitemap.ChangeFreqWeekly,
                Priority:   sitemap.PriorityCategory,
            })
        }

        if len(result.Items) < pageSize {
            break
        }
        page++
    }

    return entries, nil
}

// generatePageEntries creates sitemap entries for all published pages
func (s *Service) generatePageEntries(ctx context.Context, tenantID kernel.TenantID, baseURL string) ([]sitemap.Entry, error) {
    var entries []sitemap.Entry
    page := 1
    pageSize := 100

    for {
        pg := kernel.NewPaginationOptions(page, pageSize)
        result, err := s.pageRepo.ListByStatus(ctx, tenantID, storefront.PageStatusPublished, pg)
        if err != nil {
            return nil, err
        }

        for _, p := range result.Items {
            entries = append(entries, sitemap.Entry{
                URL:        fmt.Sprintf("%s/pages/%s", baseURL, p.Slug),
                LastMod:    p.UpdatedAt,
                ChangeFreq: sitemap.ChangeFreqMonthly,
                Priority:   sitemap.PriorityPage,
            })
        }

        if len(result.Items) < pageSize {
            break
        }
        page++
    }

    return entries, nil
}
```

---

## Step 5: Create the HTTP Handler (`sitemap/sitemapapi/handler.go`)

```go
package sitemapapi

import (
    "encoding/xml"
    "fmt"

    "github.com/gofiber/fiber/v2"
    "github.com/Abraxas-365/hada-commerce/internal/errx"
    "github.com/Abraxas-365/hada-commerce/internal/kernel"
    "github.com/Abraxas-365/hada-commerce/internal/sitemap/sitemapsrv"
)

// Handler exposes HTTP endpoints for sitemap generation
type Handler struct {
    svc *sitemapsrv.Service
}

// NewHandler creates a new sitemap handler
func NewHandler(svc *sitemapsrv.Service) *Handler {
    return &Handler{svc: svc}
}

// RegisterPublicRoutes registers unauthenticated sitemap routes
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
    router.Get("/sitemap.xml", h.GetSitemap)
}

// GetSitemap handles GET /sitemap.xml
// Returns an XML sitemap for the tenant identified by X-Tenant-ID header
func (h *Handler) GetSitemap(c *fiber.Ctx) error {
    // 1. Extract tenant from header
    tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
    if tenantID == "" {
        return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
    }

    // 2. Get base URL (construct from request or use from config)
    // Option A: Use from host header
    baseURL := fmt.Sprintf("https://%s", c.Hostname())
    
    // Option B: Could also use query parameter
    if customURL := c.Query("base_url"); customURL != "" {
        baseURL = customURL
    }

    // 3. Generate sitemap
    sitemap, err := h.svc.Generate(c.Context(), tenantID, baseURL)
    if err != nil {
        return err
    }

    // 4. Convert to XML format
    xmlSitemap := toXMLSitemap(sitemap)

    // 5. Serialize to XML
    xmlBytes, err := xml.MarshalIndent(xmlSitemap, "", "  ")
    if err != nil {
        return errx.New("failed to marshal sitemap", errx.TypeInternal)
    }

    // 6. Add XML declaration
    xmlWithDeclaration := append(
        []byte(`<?xml version="1.0" encoding="UTF-8"?>`+"\n"),
        xmlBytes...,
    )

    // 7. Return response
    c.Set("Content-Type", "application/xml; charset=utf-8")
    c.Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
    return c.Send(xmlWithDeclaration)
}

// ------- XML Types -------

// URLSet represents the root element of a sitemap XML
type URLSet struct {
    XMLName string    `xml:"urlset"`
    Xmlns   string    `xml:"xmlns,attr"`
    URLs    []URLEntry `xml:"url"`
}

// URLEntry represents a single URL entry in the sitemap
type URLEntry struct {
    Loc        string `xml:"loc"`
    LastMod    string `xml:"lastmod,omitempty"`
    ChangeFreq string `xml:"changefreq,omitempty"`
    Priority   string `xml:"priority,omitempty"`
}

// toXMLSitemap converts domain Sitemap to XML-serializable URLSet
func toXMLSitemap(s *sitemap.Sitemap) *URLSet {
    urls := make([]URLEntry, len(s.Entries))
    
    for i, entry := range s.Entries {
        urls[i] = URLEntry{
            Loc:        entry.URL,
            LastMod:    entry.LastMod.Format("2006-01-02"),
            ChangeFreq: entry.ChangeFreq,
            Priority:   fmt.Sprintf("%.1f", entry.Priority),
        }
    }

    return &URLSet{
        XMLName: "urlset",
        Xmlns:   "http://www.sitemaps.org/schemas/sitemap/0.9",
        URLs:    urls,
    }
}
```

---

## Step 6: Wire the Container (`sitemap/sitemapcontainer/container.go`)

```go
package sitemapcontainer

import (
    "github.com/gofiber/fiber/v2"
    "github.com/jmoiron/sqlx"

    "github.com/Abraxas-365/hada-commerce/internal/catalog"
    "github.com/Abraxas-365/hada-commerce/internal/product"
    "github.com/Abraxas-365/hada-commerce/internal/sitemap/sitemapapi"
    "github.com/Abraxas-365/hada-commerce/internal/sitemap/sitemapsrv"
    "github.com/Abraxas-365/hada-commerce/internal/storefront"
)

// Container wires together all sitemap domain dependencies
type Container struct {
    Service *sitemapsrv.Service
    Handler *sitemapapi.Handler
}

// New creates a fully-wired sitemap container
// Takes repositories from other domains as parameters
func New(
    productRepo product.Repository,
    categoryRepo catalog.CategoryRepository,
    collectionRepo catalog.CollectionRepository,
    pageRepo storefront.PageRepository,
) *Container {
    svc := sitemapsrv.New(productRepo, categoryRepo, collectionRepo, pageRepo)
    handler := sitemapapi.NewHandler(svc)
    return &Container{
        Service: svc,
        Handler: handler,
    }
}

// RegisterPublicRoutes registers sitemap HTTP routes on the given router
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
    c.Handler.RegisterPublicRoutes(router)
}
```

---

## Step 7: Integrate with Main Container (`backend/cmd/container.go`)

Add to Container struct:
```go
type Container struct {
    // ... existing fields ...
    Sitemap *sitemapcontainer.Container
}
```

Add to initModules():
```go
// Sitemap — uses existing domain repos
c.Sitemap = sitemapcontainer.New(
    c.Product.Service,      // Implements product.Repository
    c.Catalog.Service,      // Implements catalog.CategoryRepository
    c.Catalog.Service,      // Implements catalog.CollectionRepository
    c.Storefront.Service,   // Implements storefront.PageRepository
)
```

---

## Step 8: Register Routes (`backend/cmd/server.go`)

In `registerRoutes()` function, add to public routes:
```go
public := app.Group("/api/v1")
// ... existing routes ...
container.Sitemap.RegisterPublicRoutes(public)
```

---

## Step 9: XML Output Example

Request:
```
GET /api/v1/sitemap.xml
X-Tenant-ID: tenant-123
```

Response:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://store.example.com</loc>
    <lastmod>2024-05-30</lastmod>
    <changefreq>daily</changefreq>
    <priority>1.0</priority>
  </url>
  <url>
    <loc>https://store.example.com/products/prod-123</loc>
    <lastmod>2024-05-29</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
  <url>
    <loc>https://store.example.com/categories/electronics</loc>
    <lastmod>2024-05-25</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.7</priority>
  </url>
  <url>
    <loc>https://store.example.com/pages/about-us</loc>
    <lastmod>2024-05-20</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.6</priority>
  </url>
</urlset>
```

---

## Testing the Endpoint

### curl Example
```bash
curl -H "X-Tenant-ID: tenant-123" \
  http://localhost:3000/api/v1/sitemap.xml
```

### With custom base URL
```bash
curl -H "X-Tenant-ID: tenant-123" \
  "http://localhost:3000/api/v1/sitemap.xml?base_url=https://mystore.com"
```

---

## Performance Considerations

1. **Pagination**: Service uses pagination internally (100 items per page) to avoid memory issues
2. **Caching**: Handler sets `Cache-Control: max-age=3600` to cache for 1 hour
3. **Database**: All queries are indexed by TenantID and Status
4. **Lazy Loading**: Only fetches published/active items

For large stores (100k+ products):
- Consider implementing async sitemap generation
- Use a background job to cache sitemap_index.xml
- Store pre-generated sitemaps in Redis or file storage

---

## Deployment Notes

1. **No Auth Required**: Public endpoint accessible without API key/token
2. **Multi-tenant Safe**: Filters by X-Tenant-ID header
3. **Search Engine Compatible**: Valid sitemaps.org XML format
4. **Robots.txt Integration**: Reference in robots.txt:
   ```
   Sitemap: https://store.example.com/api/v1/sitemap.xml
   ```

---

## Next Steps

1. ✅ Create the domain module directory structure
2. ✅ Implement all 5 files (model, port, service, handler, container)
3. ✅ Integrate with main container and register routes
4. ✅ Test with curl/Postman
5. ✅ Optimize queries if needed for large stores
6. ⚠️ Consider SEO fields (meta description, keywords) if task #47 is related
