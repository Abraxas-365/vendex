package review

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines persistence operations for the review domain.
type Repository interface {
	// Create persists a new review.
	Create(ctx context.Context, r Review) (Review, error)

	// GetByID retrieves a review by ID, scoped to tenant.
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ReviewID) (Review, error)

	// ListByProduct returns paginated reviews for a product, optionally filtered by status.
	ListByProduct(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID, status string, pg kernel.PaginationOptions) (kernel.Paginated[Review], error)

	// ListByCustomer returns paginated reviews submitted by a customer.
	ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, pg kernel.PaginationOptions) (kernel.Paginated[Review], error)

	// List returns all reviews for a tenant (admin use), optionally filtered by status.
	List(ctx context.Context, tenantID kernel.TenantID, status string, pg kernel.PaginationOptions) (kernel.Paginated[Review], error)

	// UpdateStatus changes the moderation status of a review.
	UpdateStatus(ctx context.Context, tenantID kernel.TenantID, id kernel.ReviewID, status ReviewStatus) (Review, error)

	// GetStats returns aggregated rating statistics for a product.
	GetStats(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID) (ReviewStats, error)

	// IncrementHelpful increments the helpful_count for a review.
	IncrementHelpful(ctx context.Context, tenantID kernel.TenantID, id kernel.ReviewID) error

	// SetAdminResponse records an admin response to a review.
	SetAdminResponse(ctx context.Context, tenantID kernel.TenantID, id kernel.ReviewID, response string) (Review, error)
}
