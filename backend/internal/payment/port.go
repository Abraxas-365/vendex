package payment

import (
	"context"
	"encoding/json"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Repository defines persistence operations for payments and refunds.
type Repository interface {
	CreatePayment(ctx context.Context, p *Payment) error
	GetPaymentByID(ctx context.Context, tenantID kernel.TenantID, id kernel.PaymentID) (*Payment, error)
	GetPaymentByOrder(ctx context.Context, tenantID kernel.TenantID, orderID kernel.OrderID) (*Payment, error)
	UpdatePayment(ctx context.Context, p *Payment) error
	ListPaymentsByOrder(ctx context.Context, tenantID kernel.TenantID, orderID kernel.OrderID) ([]Payment, error)
	CreateRefund(ctx context.Context, r *Refund) error
	GetRefundByID(ctx context.Context, tenantID kernel.TenantID, id kernel.RefundID) (*Refund, error)
	ListRefundsByPayment(ctx context.Context, tenantID kernel.TenantID, paymentID kernel.PaymentID) ([]Refund, error)
	UpdateRefund(ctx context.Context, r *Refund) error
}

// PaymentProvider is the interface for payment gateway integrations.
type PaymentProvider interface {
	Name() string
	Charge(ctx context.Context, amount kernel.Money, token string, metadata map[string]string) (*ProviderResult, error)
	Refund(ctx context.Context, providerPaymentID string, amount kernel.Money) (*ProviderResult, error)
}

// ProviderResult holds the result from a payment provider operation.
type ProviderResult struct {
	ProviderID string          `json:"provider_id"`
	Status     string          `json:"status"`
	Data       json.RawMessage `json:"data,omitempty"`
}
