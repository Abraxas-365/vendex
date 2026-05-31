package recommendationsrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/recommendation"
	"github.com/google/uuid"
)

// Service implements the recommendation business logic.
type Service struct {
	repo recommendation.Repository
}

// New creates a recommendation Service.
func New(repo recommendation.Repository) *Service {
	return &Service{repo: repo}
}

// ---------------------------------------------------------------------------
// Tracking
// ---------------------------------------------------------------------------

// TrackView records that a visitor viewed a product.
func (s *Service) TrackView(ctx context.Context, tenantID kernel.TenantID, input recommendation.TrackViewInput) error {
	if input.ProductID == "" {
		return recommendation.ErrProductIDReq
	}
	if input.VisitorID == "" {
		return recommendation.ErrVisitorIDReq
	}

	view := recommendation.ProductView{
		ID:         kernel.NewProductViewID(uuid.NewString()),
		TenantID:   tenantID,
		VisitorID:  input.VisitorID,
		CustomerID: input.CustomerID,
		ProductID:  input.ProductID,
		Source:     input.Source,
		ViewedAt:   time.Now().UTC(),
	}

	return s.repo.TrackView(ctx, view)
}

// TrackInteraction records a product interaction (add_to_cart, purchase, wishlist, etc.).
func (s *Service) TrackInteraction(ctx context.Context, tenantID kernel.TenantID, input recommendation.TrackInteractionInput) error {
	if input.ProductID == "" {
		return recommendation.ErrProductIDReq
	}
	if input.VisitorID == "" {
		return recommendation.ErrVisitorIDReq
	}
	if input.InteractionType == "" {
		return recommendation.ErrInvalidInput
	}

	interaction := recommendation.ProductInteraction{
		ID:              kernel.NewProductInteractionID(uuid.NewString()),
		TenantID:        tenantID,
		VisitorID:       input.VisitorID,
		CustomerID:      input.CustomerID,
		ProductID:       input.ProductID,
		InteractionType: input.InteractionType,
		Metadata:        input.Metadata,
		CreatedAt:       time.Now().UTC(),
	}

	return s.repo.TrackInteraction(ctx, interaction)
}

// ---------------------------------------------------------------------------
// Recommendations
// ---------------------------------------------------------------------------

// GetRecommendations dispatches to the appropriate algorithm based on input.Type.
func (s *Service) GetRecommendations(ctx context.Context, tenantID kernel.TenantID, input recommendation.GetRecommendationsInput) ([]recommendation.RecommendedProduct, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 10
	}

	switch input.Type {
	case "frequently_bought_together":
		if input.ProductID == "" {
			return nil, recommendation.ErrProductIDReq
		}
		return s.repo.GetFrequentlyBoughtTogether(ctx, tenantID, input.ProductID, limit)

	case "trending":
		return s.repo.GetTrending(ctx, tenantID, limit, 7*24*time.Hour)

	case "recently_viewed":
		if input.VisitorID == "" {
			return nil, recommendation.ErrVisitorIDReq
		}
		return s.repo.GetRecentlyViewed(ctx, tenantID, input.VisitorID, limit)

	case "personalized":
		if input.VisitorID == "" {
			return nil, recommendation.ErrVisitorIDReq
		}
		return s.repo.GetPersonalized(ctx, tenantID, input.VisitorID, limit)

	default:
		// Default: frequently bought together if product provided, else trending
		if input.ProductID != "" {
			return s.repo.GetFrequentlyBoughtTogether(ctx, tenantID, input.ProductID, limit)
		}
		return s.repo.GetTrending(ctx, tenantID, limit, 7*24*time.Hour)
	}
}

// GetFrequentlyBoughtTogether returns products co-purchased with the given product.
func (s *Service) GetFrequentlyBoughtTogether(ctx context.Context, tenantID kernel.TenantID, productID string, limit int) ([]recommendation.RecommendedProduct, error) {
	if productID == "" {
		return nil, recommendation.ErrProductIDReq
	}
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetFrequentlyBoughtTogether(ctx, tenantID, productID, limit)
}

// GetTrending returns products with the most interactions recently.
func (s *Service) GetTrending(ctx context.Context, tenantID kernel.TenantID, limit int, since time.Duration) ([]recommendation.RecommendedProduct, error) {
	if limit <= 0 {
		limit = 10
	}
	if since <= 0 {
		since = 7 * 24 * time.Hour
	}
	return s.repo.GetTrending(ctx, tenantID, limit, since)
}

// GetRecentlyViewed returns the most recently viewed products for a visitor.
func (s *Service) GetRecentlyViewed(ctx context.Context, tenantID kernel.TenantID, visitorID string, limit int) ([]recommendation.RecommendedProduct, error) {
	if visitorID == "" {
		return nil, recommendation.ErrVisitorIDReq
	}
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetRecentlyViewed(ctx, tenantID, visitorID, limit)
}

// GetPersonalized returns personalised recommendations based on a visitor's history.
func (s *Service) GetPersonalized(ctx context.Context, tenantID kernel.TenantID, visitorID string, limit int) ([]recommendation.RecommendedProduct, error) {
	if visitorID == "" {
		return nil, recommendation.ErrVisitorIDReq
	}
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetPersonalized(ctx, tenantID, visitorID, limit)
}

// ---------------------------------------------------------------------------
// Rules
// ---------------------------------------------------------------------------

// CreateRule persists a new recommendation rule.
func (s *Service) CreateRule(ctx context.Context, tenantID kernel.TenantID, input recommendation.CreateRuleInput) (*recommendation.RecommendationRule, error) {
	if input.Name == "" {
		return nil, recommendation.ErrRuleNameReq
	}
	if input.Type == "" {
		return nil, recommendation.ErrRuleTypeReq
	}

	now := time.Now().UTC()
	rule := recommendation.RecommendationRule{
		ID:        kernel.NewRecommendationRuleID(uuid.NewString()),
		TenantID:  tenantID,
		Name:      input.Name,
		Type:      input.Type,
		Config:    input.Config,
		IsActive:  input.IsActive,
		Priority:  input.Priority,
		CreatedAt: now,
		UpdatedAt: now,
	}

	created, err := s.repo.CreateRule(ctx, rule)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

// ListRules returns all recommendation rules for a tenant.
func (s *Service) ListRules(ctx context.Context, tenantID kernel.TenantID) ([]recommendation.RecommendationRule, error) {
	return s.repo.ListRules(ctx, tenantID)
}

// UpdateRule applies changes to an existing rule.
func (s *Service) UpdateRule(ctx context.Context, tenantID kernel.TenantID, id kernel.RecommendationRuleID, input recommendation.UpdateRuleInput) (*recommendation.RecommendationRule, error) {
	existing, err := s.repo.GetRuleByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Type != nil {
		existing.Type = *input.Type
	}
	if input.Config != nil {
		existing.Config = input.Config
	}
	if input.IsActive != nil {
		existing.IsActive = *input.IsActive
	}
	if input.Priority != nil {
		existing.Priority = *input.Priority
	}
	existing.UpdatedAt = time.Now().UTC()

	updated, err := s.repo.UpdateRule(ctx, existing)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

// DeleteRule removes a recommendation rule.
func (s *Service) DeleteRule(ctx context.Context, tenantID kernel.TenantID, id kernel.RecommendationRuleID) error {
	_, err := s.repo.GetRuleByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	return s.repo.DeleteRule(ctx, tenantID, id)
}
