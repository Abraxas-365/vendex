package review

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// ReviewStatus represents the moderation state of a review.
type ReviewStatus string

const (
	StatusPending  ReviewStatus = "pending"
	StatusApproved ReviewStatus = "approved"
	StatusRejected ReviewStatus = "rejected"
)

// Review is the aggregate root for a product review.
type Review struct {
	ID               kernel.ReviewID   `json:"id"`
	TenantID         kernel.TenantID   `json:"tenant_id"`
	ProductID        kernel.ProductID  `json:"product_id"`
	CustomerID       kernel.CustomerID `json:"customer_id"`
	Rating           int               `json:"rating"`
	Title            string            `json:"title,omitempty"`
	Body             string            `json:"body,omitempty"`
	Status           ReviewStatus      `json:"status"`
	VerifiedPurchase bool              `json:"verified_purchase"`
	HelpfulCount     int               `json:"helpful_count"`
	Images           []string          `json:"images,omitempty"`
	AdminResponse    *string           `json:"admin_response,omitempty"`
	AdminRespondedAt *time.Time        `json:"admin_responded_at,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

// CreateReviewInput holds the data needed to create a review.
type CreateReviewInput struct {
	ProductID  kernel.ProductID
	CustomerID kernel.CustomerID
	Rating     int
	Title      string
	Body       string
	Images     []string
}

// ReviewStats holds aggregated rating statistics for a product.
type ReviewStats struct {
	AverageRating float64      `json:"average_rating"`
	TotalReviews  int          `json:"total_reviews"`
	Distribution  map[int]int  `json:"distribution"` // rating value -> count
}
