package giftcard

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Repository defines persistence operations for GiftCard entities.
// All operations are scoped by TenantID.
type Repository interface {
	// Create persists a new gift card.
	Create(ctx context.Context, giftCard *GiftCard) error
	// GetByID retrieves a gift card by its primary key.
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.GiftCardID) (*GiftCard, error)
	// GetByCode retrieves a gift card by its unique code (case-insensitive).
	GetByCode(ctx context.Context, tenantID kernel.TenantID, code string) (*GiftCard, error)
	// List returns all gift cards for a tenant with pagination.
	List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[GiftCard], error)
	// Update persists mutations to an existing gift card (balance, active, etc.).
	Update(ctx context.Context, giftCard *GiftCard) error
	// Delete removes a gift card by ID.
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.GiftCardID) error
	// CreateTransaction records a credit or debit transaction.
	CreateTransaction(ctx context.Context, tx *GiftCardTransaction) error
	// ListTransactions returns all transactions for a gift card.
	ListTransactions(ctx context.Context, tenantID kernel.TenantID, giftCardID kernel.GiftCardID) ([]GiftCardTransaction, error)
}
