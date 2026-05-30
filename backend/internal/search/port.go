package search

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines persistence operations for the search domain.
type Repository interface {
	Search(ctx context.Context, tenantID kernel.TenantID, query SearchQuery) (*SearchResult, error)
	Suggest(ctx context.Context, tenantID kernel.TenantID, prefix string, limit int) ([]SearchSuggestion, error)
}
