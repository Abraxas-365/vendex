package recommendation

import (
	"context"
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Repository defines the persistence contract for the recommendation domain.
type Repository interface {
	// TrackView persists a product view event.
	TrackView(ctx context.Context, view ProductView) error

	// TrackInteraction persists a product interaction event.
	TrackInteraction(ctx context.Context, interaction ProductInteraction) error

	// GetFrequentlyBoughtTogether returns products frequently co-purchased with the given product.
	GetFrequentlyBoughtTogether(ctx context.Context, tenantID kernel.TenantID, productID string, limit int) ([]RecommendedProduct, error)

	// GetTrending returns products with the highest interaction count within the given duration.
	GetTrending(ctx context.Context, tenantID kernel.TenantID, limit int, since time.Duration) ([]RecommendedProduct, error)

	// GetRecentlyViewed returns the most recently viewed products for a visitor.
	GetRecentlyViewed(ctx context.Context, tenantID kernel.TenantID, visitorID string, limit int) ([]RecommendedProduct, error)

	// GetPersonalized returns recommended products based on the visitor's interaction history.
	GetPersonalized(ctx context.Context, tenantID kernel.TenantID, visitorID string, limit int) ([]RecommendedProduct, error)

	// CreateRule persists a new recommendation rule.
	CreateRule(ctx context.Context, rule RecommendationRule) (RecommendationRule, error)

	// GetRuleByID returns a recommendation rule scoped to a tenant.
	GetRuleByID(ctx context.Context, tenantID kernel.TenantID, id kernel.RecommendationRuleID) (RecommendationRule, error)

	// ListRules returns all recommendation rules for a tenant.
	ListRules(ctx context.Context, tenantID kernel.TenantID) ([]RecommendationRule, error)

	// UpdateRule replaces a recommendation rule's mutable fields.
	UpdateRule(ctx context.Context, rule RecommendationRule) (RecommendationRule, error)

	// DeleteRule removes a recommendation rule.
	DeleteRule(ctx context.Context, tenantID kernel.TenantID, id kernel.RecommendationRuleID) error
}
