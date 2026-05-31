package recommendation

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// ProductView records that a visitor viewed a product.
type ProductView struct {
	ID         kernel.ProductViewID `json:"id"`
	TenantID   kernel.TenantID      `json:"tenant_id"`
	VisitorID  string               `json:"visitor_id"`
	CustomerID string               `json:"customer_id,omitempty"`
	ProductID  string               `json:"product_id"`
	Source     string               `json:"source,omitempty"`
	ViewedAt   time.Time            `json:"viewed_at"`
}

// ProductInteraction records a behavioural event (view, cart, purchase, wishlist).
type ProductInteraction struct {
	ID              kernel.ProductInteractionID `json:"id"`
	TenantID        kernel.TenantID             `json:"tenant_id"`
	VisitorID       string                      `json:"visitor_id"`
	CustomerID      string                      `json:"customer_id,omitempty"`
	ProductID       string                      `json:"product_id"`
	InteractionType string                      `json:"interaction_type"`
	Metadata        map[string]interface{}      `json:"metadata,omitempty"`
	CreatedAt       time.Time                   `json:"created_at"`
}

// RecommendationRule configures how a given recommendation algorithm is applied.
type RecommendationRule struct {
	ID        kernel.RecommendationRuleID `json:"id"`
	TenantID  kernel.TenantID             `json:"tenant_id"`
	Name      string                      `json:"name"`
	Type      string                      `json:"type"`
	Config    map[string]interface{}      `json:"config,omitempty"`
	IsActive  bool                        `json:"is_active"`
	Priority  int                         `json:"priority"`
	CreatedAt time.Time                   `json:"created_at"`
	UpdatedAt time.Time                   `json:"updated_at"`
}

// RecommendedProduct is a lightweight result item returned to callers.
type RecommendedProduct struct {
	ProductID string  `json:"product_id"`
	Score     float64 `json:"score"`
	Reason    string  `json:"reason"`
}

// ---------------------------------------------------------------------------
// Input types
// ---------------------------------------------------------------------------

// TrackViewInput is the payload for recording a product view event.
type TrackViewInput struct {
	VisitorID  string
	CustomerID string
	ProductID  string
	Source     string
}

// TrackInteractionInput is the payload for recording a product interaction.
type TrackInteractionInput struct {
	VisitorID       string
	CustomerID      string
	ProductID       string
	InteractionType string
	Metadata        map[string]interface{}
}

// GetRecommendationsInput drives the generic recommendation lookup.
type GetRecommendationsInput struct {
	ProductID  string
	VisitorID  string
	CustomerID string
	Type       string
	Limit      int
}

// CreateRuleInput holds the data required to create a recommendation rule.
type CreateRuleInput struct {
	Name     string
	Type     string
	Config   map[string]interface{}
	IsActive bool
	Priority int
}

// UpdateRuleInput holds the mutable fields of a recommendation rule.
type UpdateRuleInput struct {
	Name     *string
	Type     *string
	Config   map[string]interface{}
	IsActive *bool
	Priority *int
}
