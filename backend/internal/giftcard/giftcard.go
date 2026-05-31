package giftcard

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// GiftCard represents a monetary gift card issued to a customer.
// Business rules:
//   - A gift card is redeemable only when Active=true and not expired.
//   - The balance is decremented on each redemption, and a debit transaction is recorded.
//   - The balance can never go below zero.
type GiftCard struct {
	ID            kernel.GiftCardID `json:"id"             db:"id"`
	TenantID      kernel.TenantID   `json:"tenant_id"      db:"tenant_id"`
	Code          string            `json:"code"           db:"code"`
	InitialAmount kernel.Money      `json:"initial_amount" db:"initial_amount"`
	Balance       kernel.Money      `json:"balance"        db:"balance"`
	ExpiresAt     *time.Time        `json:"expires_at,omitempty" db:"expires_at"`
	Active        bool              `json:"active"         db:"active"`
	CreatedBy     string            `json:"created_by"     db:"created_by"`
	CreatedAt     time.Time         `json:"created_at"     db:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"     db:"updated_at"`
}

// IsExpired returns true when the gift card is past its ExpiresAt date.
func (g *GiftCard) IsExpired(now time.Time) bool {
	return g.ExpiresAt != nil && now.After(*g.ExpiresAt)
}

// HasSufficientBalance returns true when the gift card balance is >= the requested amount.
func (g *GiftCard) HasSufficientBalance(amount kernel.Money) bool {
	return g.Balance.Amount >= amount.Amount
}

// GiftCardTransaction records a credit or debit against a gift card.
type GiftCardTransaction struct {
	ID         kernel.GiftCardTransactionID `json:"id"                    db:"id"`
	GiftCardID kernel.GiftCardID            `json:"gift_card_id"          db:"gift_card_id"`
	TenantID   kernel.TenantID              `json:"tenant_id"             db:"tenant_id"`
	Type       string                       `json:"type"                  db:"type"` // "credit" or "debit"
	Amount     kernel.Money                 `json:"amount"                db:"amount"`
	OrderID    *kernel.OrderID              `json:"order_id,omitempty"    db:"order_id"`
	Note       string                       `json:"note"                  db:"note"`
	CreatedAt  time.Time                    `json:"created_at"            db:"created_at"`
}

// Transaction type constants.
const (
	TransactionTypeCredit = "credit"
	TransactionTypeDebit  = "debit"
)

// CreateInput holds the data needed to create a new gift card.
type CreateInput struct {
	Code          string
	InitialAmount kernel.Money
	ExpiresAt     *time.Time
	CreatedBy     string
}

// UpdateInput holds the data for updating a gift card (admin use).
type UpdateInput struct {
	Code      *string
	ExpiresAt *time.Time
	Active    *bool
}
