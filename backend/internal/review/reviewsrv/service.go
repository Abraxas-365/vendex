package reviewsrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/review"
	"github.com/google/uuid"
)

// Service handles review business logic.
type Service struct {
	repo review.Repository
	bus  eventbus.Bus
}

// New creates a new review service.
func New(repo review.Repository, bus eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// Create validates and persists a new review, then fires review.created.
func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, input review.CreateReviewInput) (review.Review, error) {
	if input.ProductID.IsEmpty() {
		return review.Review{}, review.ErrProductIDRequired
	}
	if input.CustomerID.IsEmpty() {
		return review.Review{}, review.ErrCustomerIDRequired
	}
	if input.Rating < 1 || input.Rating > 5 {
		return review.Review{}, review.ErrInvalidRating
	}

	now := time.Now().UTC()
	images := input.Images
	if images == nil {
		images = []string{}
	}

	r := review.Review{
		ID:               kernel.ReviewID(uuid.NewString()),
		TenantID:         tenantID,
		ProductID:        input.ProductID,
		CustomerID:       input.CustomerID,
		Rating:           input.Rating,
		Title:            input.Title,
		Body:             input.Body,
		Status:           review.StatusPending,
		VerifiedPurchase: false,
		HelpfulCount:     0,
		Images:           images,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	created, err := s.repo.Create(ctx, r)
	if err != nil {
		return review.Review{}, errx.Wrap(err, "creating review", errx.TypeInternal)
	}

	if evt, err := eventbus.NewEvent(eventbus.ReviewCreated, tenantID, eventbus.ReviewPayload{
		ReviewID:   string(created.ID),
		ProductID:  string(created.ProductID),
		CustomerID: string(created.CustomerID),
		Rating:     created.Rating,
		Status:     string(created.Status),
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return created, nil
}

// ListByProduct returns paginated, optionally status-filtered reviews for a product.
func (s *Service) ListByProduct(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID, status string, pg kernel.PaginationOptions) (kernel.Paginated[review.Review], error) {
	return s.repo.ListByProduct(ctx, tenantID, productID, status, pg)
}

// ListByCustomer returns paginated reviews for a specific customer.
func (s *Service) ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, pg kernel.PaginationOptions) (kernel.Paginated[review.Review], error) {
	return s.repo.ListByCustomer(ctx, tenantID, customerID, pg)
}

// List returns all reviews for a tenant, optionally filtered by status (admin use).
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, status string, pg kernel.PaginationOptions) (kernel.Paginated[review.Review], error) {
	return s.repo.List(ctx, tenantID, status, pg)
}

// GetByID retrieves a review by its ID.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ReviewID) (review.Review, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// Approve approves a pending review and fires review.approved.
func (s *Service) Approve(ctx context.Context, tenantID kernel.TenantID, id kernel.ReviewID) (review.Review, error) {
	r, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return review.Review{}, err
	}
	if r.Status != review.StatusPending {
		return review.Review{}, review.ErrAlreadyModerated
	}

	updated, err := s.repo.UpdateStatus(ctx, tenantID, id, review.StatusApproved)
	if err != nil {
		return review.Review{}, err
	}

	if evt, err := eventbus.NewEvent(eventbus.ReviewApproved, tenantID, eventbus.ReviewPayload{
		ReviewID:   string(updated.ID),
		ProductID:  string(updated.ProductID),
		CustomerID: string(updated.CustomerID),
		Rating:     updated.Rating,
		Status:     string(updated.Status),
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return updated, nil
}

// Reject rejects a pending review and fires review.rejected.
func (s *Service) Reject(ctx context.Context, tenantID kernel.TenantID, id kernel.ReviewID) (review.Review, error) {
	r, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return review.Review{}, err
	}
	if r.Status != review.StatusPending {
		return review.Review{}, review.ErrAlreadyModerated
	}

	updated, err := s.repo.UpdateStatus(ctx, tenantID, id, review.StatusRejected)
	if err != nil {
		return review.Review{}, err
	}

	if evt, err := eventbus.NewEvent(eventbus.ReviewRejected, tenantID, eventbus.ReviewPayload{
		ReviewID:   string(updated.ID),
		ProductID:  string(updated.ProductID),
		CustomerID: string(updated.CustomerID),
		Rating:     updated.Rating,
		Status:     string(updated.Status),
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return updated, nil
}

// MarkHelpful increments the helpful count for a review.
func (s *Service) MarkHelpful(ctx context.Context, tenantID kernel.TenantID, id kernel.ReviewID) error {
	// Verify the review exists and belongs to tenant.
	if _, err := s.repo.GetByID(ctx, tenantID, id); err != nil {
		return err
	}
	return s.repo.IncrementHelpful(ctx, tenantID, id)
}

// RespondAsAdmin records an admin response on a review.
func (s *Service) RespondAsAdmin(ctx context.Context, tenantID kernel.TenantID, id kernel.ReviewID, response string) (review.Review, error) {
	if response == "" {
		return review.Review{}, errx.New("response cannot be empty", errx.TypeValidation)
	}
	// Verify the review exists and belongs to tenant.
	if _, err := s.repo.GetByID(ctx, tenantID, id); err != nil {
		return review.Review{}, err
	}
	return s.repo.SetAdminResponse(ctx, tenantID, id, response)
}

// GetStats returns aggregated rating statistics for a product.
func (s *Service) GetStats(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID) (review.ReviewStats, error) {
	return s.repo.GetStats(ctx, tenantID, productID)
}
