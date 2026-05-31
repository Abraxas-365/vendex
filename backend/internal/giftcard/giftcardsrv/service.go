package giftcardsrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/giftcard"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/google/uuid"
)

// Service implements gift card business logic.
type Service struct {
	repo giftcard.Repository
}

// New creates a new gift card Service.
func New(repo giftcard.Repository) *Service {
	return &Service{repo: repo}
}

// Create persists a new gift card with the given initial balance.
func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, input giftcard.CreateInput) (*giftcard.GiftCard, error) {
	if input.Code == "" {
		return nil, giftcard.ErrInvalidCode
	}

	now := time.Now().UTC()
	gc := &giftcard.GiftCard{
		ID:            kernel.GiftCardID(uuid.NewString()),
		TenantID:      tenantID,
		Code:          input.Code,
		InitialAmount: input.InitialAmount,
		Balance:       input.InitialAmount,
		ExpiresAt:     input.ExpiresAt,
		Active:        true,
		CreatedBy:     input.CreatedBy,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.Create(ctx, gc); err != nil {
		return nil, err
	}

	// Record the initial credit transaction.
	tx := &giftcard.GiftCardTransaction{
		ID:         kernel.GiftCardTransactionID(uuid.NewString()),
		GiftCardID: gc.ID,
		TenantID:   tenantID,
		Type:       giftcard.TransactionTypeCredit,
		Amount:     input.InitialAmount,
		Note:       "initial credit",
		CreatedAt:  now,
	}
	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, errx.Wrap(err, "record initial credit transaction", errx.TypeInternal)
	}

	return gc, nil
}

// GetByID retrieves a gift card by its primary key.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.GiftCardID) (*giftcard.GiftCard, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// GetByCode retrieves a gift card by its code.
func (s *Service) GetByCode(ctx context.Context, tenantID kernel.TenantID, code string) (*giftcard.GiftCard, error) {
	return s.repo.GetByCode(ctx, tenantID, code)
}

// CheckBalance returns the gift card for a given code so the caller can inspect the balance.
// This is a public operation — it does not require authentication beyond tenant scoping.
func (s *Service) CheckBalance(ctx context.Context, tenantID kernel.TenantID, code string) (*giftcard.GiftCard, error) {
	return s.repo.GetByCode(ctx, tenantID, code)
}

// RedeemInput holds the data needed to redeem a gift card.
type RedeemInput struct {
	Code    string
	Amount  kernel.Money
	OrderID *kernel.OrderID
	Note    string
}

// Redeem validates a gift card and deducts the given amount from its balance.
// It records a debit transaction. Returns the updated gift card.
func (s *Service) Redeem(ctx context.Context, tenantID kernel.TenantID, input RedeemInput) (*giftcard.GiftCard, error) {
	gc, err := s.repo.GetByCode(ctx, tenantID, input.Code)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	if !gc.Active {
		return nil, giftcard.ErrInactive
	}
	if gc.IsExpired(now) {
		return nil, giftcard.ErrExpired
	}
	if !gc.HasSufficientBalance(input.Amount) {
		return nil, giftcard.ErrInsufficientBalance
	}

	// Deduct balance.
	gc.Balance.Amount -= input.Amount.Amount
	gc.UpdatedAt = now

	if err := s.repo.Update(ctx, gc); err != nil {
		return nil, errx.Wrap(err, "update gift card balance", errx.TypeInternal)
	}

	// Record debit transaction.
	note := input.Note
	if note == "" {
		note = "redemption"
	}
	tx := &giftcard.GiftCardTransaction{
		ID:         kernel.GiftCardTransactionID(uuid.NewString()),
		GiftCardID: gc.ID,
		TenantID:   tenantID,
		Type:       giftcard.TransactionTypeDebit,
		Amount:     input.Amount,
		OrderID:    input.OrderID,
		Note:       note,
		CreatedAt:  now,
	}
	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, errx.Wrap(err, "record debit transaction", errx.TypeInternal)
	}

	return gc, nil
}

// Disable sets a gift card's Active flag to false.
func (s *Service) Disable(ctx context.Context, tenantID kernel.TenantID, id kernel.GiftCardID) (*giftcard.GiftCard, error) {
	gc, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	gc.Active = false
	gc.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, gc); err != nil {
		return nil, errx.Wrap(err, "disable gift card", errx.TypeInternal)
	}
	return gc, nil
}

// List returns all gift cards for a tenant with pagination.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[giftcard.GiftCard], error) {
	return s.repo.List(ctx, tenantID, p)
}

// ListTransactions returns all transactions for a given gift card.
func (s *Service) ListTransactions(ctx context.Context, tenantID kernel.TenantID, giftCardID kernel.GiftCardID) ([]giftcard.GiftCardTransaction, error) {
	// Verify the gift card belongs to the tenant.
	if _, err := s.repo.GetByID(ctx, tenantID, giftCardID); err != nil {
		return nil, err
	}
	return s.repo.ListTransactions(ctx, tenantID, giftCardID)
}

// Delete removes a gift card permanently.
func (s *Service) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.GiftCardID) error {
	// Verify exists before deleting.
	if _, err := s.repo.GetByID(ctx, tenantID, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, tenantID, id)
}

// UpdateCard persists a mutated gift card (admin operation).
// The caller is responsible for setting the UpdatedAt timestamp.
func (s *Service) UpdateCard(ctx context.Context, tenantID kernel.TenantID, gc *giftcard.GiftCard) error {
	if gc.TenantID != tenantID {
		return errx.Validation("tenant mismatch")
	}
	gc.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, gc)
}
