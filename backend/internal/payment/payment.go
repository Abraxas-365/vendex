package payment

import (
	"encoding/json"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// PaymentStatus represents the lifecycle of a payment.
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted  PaymentStatus = "completed"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusRefunded   PaymentStatus = "refunded"
)

// Payment represents a payment transaction for an order.
type Payment struct {
	ID                kernel.PaymentID `json:"id" db:"id"`
	TenantID          kernel.TenantID  `json:"tenant_id" db:"tenant_id"`
	OrderID           kernel.OrderID   `json:"order_id" db:"order_id"`
	Amount            kernel.Money     `json:"amount"`
	Status            PaymentStatus    `json:"status" db:"status"`
	Provider          string           `json:"provider" db:"provider"`                     // "stripe", "paypal", "manual"
	ProviderPaymentID string           `json:"provider_payment_id,omitempty" db:"provider_payment_id"`
	ProviderData      json.RawMessage  `json:"provider_data,omitempty" db:"provider_data"` // JSONB
	Method            string           `json:"method,omitempty" db:"method"`               // "card", "bank_transfer", "cash"
	ErrorMessage      string           `json:"error_message,omitempty" db:"error_message"`
	PaidAt            *time.Time       `json:"paid_at,omitempty" db:"paid_at"`
	CreatedAt         time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at" db:"updated_at"`
}

// RefundStatus represents the lifecycle of a refund.
type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "pending"
	RefundStatusCompleted RefundStatus = "completed"
	RefundStatusFailed    RefundStatus = "failed"
)

// Refund represents a refund against a payment.
type Refund struct {
	ID               kernel.RefundID  `json:"id" db:"id"`
	TenantID         kernel.TenantID  `json:"tenant_id" db:"tenant_id"`
	PaymentID        kernel.PaymentID `json:"payment_id" db:"payment_id"`
	OrderID          kernel.OrderID   `json:"order_id" db:"order_id"`
	Amount           kernel.Money     `json:"amount"`
	Reason           string           `json:"reason,omitempty" db:"reason"`
	Status           RefundStatus     `json:"status" db:"status"`
	ProviderRefundID string           `json:"provider_refund_id,omitempty" db:"provider_refund_id"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at" db:"updated_at"`
}
