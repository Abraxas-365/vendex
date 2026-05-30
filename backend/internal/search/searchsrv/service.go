package searchsrv

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/search"
)

// Service implements business logic for the search domain.
type Service struct {
	repo search.Repository
}

// New creates a new search Service.
func New(repo search.Repository) *Service {
	return &Service{repo: repo}
}

// Search validates and executes a product search.
func (s *Service) Search(ctx context.Context, tenantID kernel.TenantID, q search.SearchQuery) (*search.SearchResult, error) {
	// Apply pagination defaults.
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = 20
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}

	// Default sort: relevance when a query is present, newest otherwise.
	if q.SortBy == "" {
		if q.Query != "" {
			q.SortBy = "relevance"
		} else {
			q.SortBy = "created_at"
		}
	}

	return s.repo.Search(ctx, tenantID, q)
}

// Suggest returns autocomplete suggestions for the given prefix.
func (s *Service) Suggest(ctx context.Context, tenantID kernel.TenantID, prefix string, limit int) ([]search.SearchSuggestion, error) {
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	return s.repo.Suggest(ctx, tenantID, prefix, limit)
}
