package sitemap

import (
	"context"
	"fmt"
	"time"

	"github.com/Abraxas-365/vendex/internal/catalog"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/product"
	"github.com/Abraxas-365/vendex/internal/storefront"
)

// ProductLister is the subset of productsrv.Service needed by the sitemap service.
type ProductLister interface {
	List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[product.Product], error)
}

// CategoryLister is the subset of catalogsrv.Service needed by the sitemap service.
type CategoryLister interface {
	ListCategories(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[catalog.Category], error)
	ListCollections(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[catalog.Collection], error)
}

// PageLister is the subset of storefrontsrv.Service needed by the sitemap service.
type PageLister interface {
	ListPages(ctx context.Context, tenantID kernel.TenantID, status *storefront.PageStatus, p kernel.PaginationOptions) (kernel.Paginated[storefront.Page], error)
}

// Service generates sitemaps by querying products, categories, collections, and pages.
type Service struct {
	products   ProductLister
	catalog    CategoryLister
	storefront PageLister
}

// NewService creates a new sitemap Service.
func NewService(products ProductLister, catalog CategoryLister, storefront PageLister) *Service {
	return &Service{
		products:   products,
		catalog:    catalog,
		storefront: storefront,
	}
}

// Generate produces a Sitemap for the given tenant using baseURL as the URL prefix.
// It fetches all active products, categories, collections, and published pages.
func (s *Service) Generate(ctx context.Context, tenantID kernel.TenantID, baseURL string) (*Sitemap, error) {
	now := FormatTime(time.Now())
	sm := &Sitemap{}

	// Always include the homepage.
	sm.URLs = append(sm.URLs, URLEntry{
		Loc:        baseURL + "/",
		LastMod:    now,
		ChangeFreq: "daily",
		Priority:   "1.0",
	})

	// Products — all pages, status=active filter is applied at fetch time.
	if err := s.appendProducts(ctx, tenantID, baseURL, sm); err != nil {
		return nil, err
	}

	// Categories
	if err := s.appendCategories(ctx, tenantID, baseURL, sm); err != nil {
		return nil, err
	}

	// Collections
	if err := s.appendCollections(ctx, tenantID, baseURL, sm); err != nil {
		return nil, err
	}

	// Published pages
	if err := s.appendPages(ctx, tenantID, baseURL, sm); err != nil {
		return nil, err
	}

	return sm, nil
}

const maxPageSize = 100

func (s *Service) appendProducts(ctx context.Context, tenantID kernel.TenantID, baseURL string, sm *Sitemap) error {
	page := 1
	for {
		pg := kernel.PaginationOptions{Page: page, PageSize: maxPageSize}
		result, err := s.products.List(ctx, tenantID, pg)
		if err != nil {
			return fmt.Errorf("sitemap: listing products page %d: %w", page, err)
		}

		for _, p := range result.Items {
			if p.Status != product.StatusActive {
				continue
			}
			sm.URLs = append(sm.URLs, URLEntry{
				Loc:        fmt.Sprintf("%s/products/%s", baseURL, string(p.ID)),
				LastMod:    FormatTime(p.UpdatedAt),
				ChangeFreq: "weekly",
				Priority:   "0.8",
			})
		}

		if page >= result.TotalPages || len(result.Items) == 0 {
			break
		}
		page++
	}
	return nil
}

func (s *Service) appendCategories(ctx context.Context, tenantID kernel.TenantID, baseURL string, sm *Sitemap) error {
	page := 1
	for {
		pg := kernel.PaginationOptions{Page: page, PageSize: maxPageSize}
		result, err := s.catalog.ListCategories(ctx, tenantID, pg)
		if err != nil {
			return fmt.Errorf("sitemap: listing categories page %d: %w", page, err)
		}

		for _, cat := range result.Items {
			slug := cat.Slug
			if slug == "" {
				slug = string(cat.ID)
			}
			sm.URLs = append(sm.URLs, URLEntry{
				Loc:        fmt.Sprintf("%s/categories/%s", baseURL, slug),
				LastMod:    FormatTime(cat.UpdatedAt),
				ChangeFreq: "weekly",
				Priority:   "0.6",
			})
		}

		if page >= result.TotalPages || len(result.Items) == 0 {
			break
		}
		page++
	}
	return nil
}

func (s *Service) appendCollections(ctx context.Context, tenantID kernel.TenantID, baseURL string, sm *Sitemap) error {
	page := 1
	for {
		pg := kernel.PaginationOptions{Page: page, PageSize: maxPageSize}
		result, err := s.catalog.ListCollections(ctx, tenantID, pg)
		if err != nil {
			return fmt.Errorf("sitemap: listing collections page %d: %w", page, err)
		}

		for _, col := range result.Items {
			slug := col.Slug
			if slug == "" {
				slug = string(col.ID)
			}
			sm.URLs = append(sm.URLs, URLEntry{
				Loc:        fmt.Sprintf("%s/collections/%s", baseURL, slug),
				LastMod:    FormatTime(col.UpdatedAt),
				ChangeFreq: "weekly",
				Priority:   "0.6",
			})
		}

		if page >= result.TotalPages || len(result.Items) == 0 {
			break
		}
		page++
	}
	return nil
}

func (s *Service) appendPages(ctx context.Context, tenantID kernel.TenantID, baseURL string, sm *Sitemap) error {
	published := storefront.PageStatusPublished
	page := 1
	for {
		pg := kernel.PaginationOptions{Page: page, PageSize: maxPageSize}
		result, err := s.storefront.ListPages(ctx, tenantID, &published, pg)
		if err != nil {
			return fmt.Errorf("sitemap: listing pages page %d: %w", page, err)
		}

		for _, p := range result.Items {
			sm.URLs = append(sm.URLs, URLEntry{
				Loc:        fmt.Sprintf("%s/%s", baseURL, p.Slug),
				LastMod:    FormatTime(p.UpdatedAt),
				ChangeFreq: "monthly",
				Priority:   "0.5",
			})
		}

		if page >= result.TotalPages || len(result.Items) == 0 {
			break
		}
		page++
	}
	return nil
}
