package cartrecoverysrv

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Abraxas-365/hada-commerce/internal/cartrecovery"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Service implements business logic for the cart recovery domain.
type Service struct {
	repo cartrecovery.Repository
}

// New creates a new cart recovery service.
func New(repo cartrecovery.Repository) *Service {
	return &Service{repo: repo}
}

// ScheduleRecovery creates a step-1 pending recovery email for an abandoned cart.
// Returns ErrAlreadyScheduled if a recovery email already exists for the cart.
func (s *Service) ScheduleRecovery(
	ctx context.Context,
	tenantID kernel.TenantID,
	cartID kernel.CartID,
	customerID kernel.CustomerID,
	email string,
) (*cartrecovery.RecoveryEmail, error) {
	// Guard: do not schedule twice for the same cart.
	existing, err := s.repo.GetByCartID(ctx, tenantID, cartID)
	if err != nil {
		return nil, err
	}
	if len(existing) > 0 {
		return nil, cartrecovery.ErrAlreadyScheduled
	}

	rec := &cartrecovery.RecoveryEmail{
		ID:         kernel.NewRecoveryID(uuid.NewString()),
		TenantID:   tenantID,
		CartID:     cartID,
		CustomerID: customerID,
		Email:      email,
		Step:       1,
		Status:     cartrecovery.StatusPending,
		CreatedAt:  time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, rec); err != nil {
		return nil, err
	}
	return rec, nil
}

// MarkSent transitions a recovery email to "sent" and records the sent timestamp.
func (s *Service) MarkSent(ctx context.Context, tenantID kernel.TenantID, id kernel.RecoveryID) (*cartrecovery.RecoveryEmail, error) {
	rec, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if rec.Status != cartrecovery.StatusPending {
		return nil, cartrecovery.ErrInvalidStatus
	}

	now := time.Now().UTC()
	rec.Status = cartrecovery.StatusSent
	rec.SentAt = &now

	if err := s.repo.Update(ctx, rec); err != nil {
		return nil, err
	}
	return rec, nil
}

// MarkClicked transitions a recovery email to "clicked" and records the clicked timestamp.
func (s *Service) MarkClicked(ctx context.Context, tenantID kernel.TenantID, id kernel.RecoveryID) (*cartrecovery.RecoveryEmail, error) {
	rec, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if rec.Status != cartrecovery.StatusSent {
		return nil, cartrecovery.ErrInvalidStatus
	}

	now := time.Now().UTC()
	rec.Status = cartrecovery.StatusClicked
	rec.ClickedAt = &now

	if err := s.repo.Update(ctx, rec); err != nil {
		return nil, err
	}
	return rec, nil
}

// MarkConverted transitions a recovery email to "converted" and records the converted timestamp.
func (s *Service) MarkConverted(ctx context.Context, tenantID kernel.TenantID, id kernel.RecoveryID) (*cartrecovery.RecoveryEmail, error) {
	rec, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if rec.Status != cartrecovery.StatusClicked && rec.Status != cartrecovery.StatusSent {
		return nil, cartrecovery.ErrInvalidStatus
	}

	now := time.Now().UTC()
	rec.Status = cartrecovery.StatusConverted
	rec.ConvertedAt = &now

	if err := s.repo.Update(ctx, rec); err != nil {
		return nil, err
	}
	return rec, nil
}

// ListPending returns all pending recovery emails for the given tenant.
func (s *Service) ListPending(ctx context.Context, tenantID kernel.TenantID) ([]cartrecovery.RecoveryEmail, error) {
	return s.repo.ListPending(ctx, tenantID)
}

// List returns a paginated list of all recovery emails for the given tenant.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, page, pageSize int) (kernel.Paginated[cartrecovery.RecoveryEmail], error) {
	return s.repo.List(ctx, tenantID, page, pageSize)
}

// GetStats returns aggregate recovery statistics for the given tenant.
func (s *Service) GetStats(ctx context.Context, tenantID kernel.TenantID) (*cartrecovery.RecoveryStats, error) {
	return s.repo.GetStats(ctx, tenantID)
}
