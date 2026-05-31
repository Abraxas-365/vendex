package subscription

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// BillingInterval represents how often a subscription is billed.
type BillingInterval string

const (
	IntervalWeekly    BillingInterval = "weekly"
	IntervalMonthly   BillingInterval = "monthly"
	IntervalQuarterly BillingInterval = "quarterly"
	IntervalYearly    BillingInterval = "yearly"
)

// IsValid checks if the interval is one of the known values.
func (i BillingInterval) IsValid() bool {
	switch i {
	case IntervalWeekly, IntervalMonthly, IntervalQuarterly, IntervalYearly:
		return true
	}
	return false
}

// SubscriptionStatus represents the lifecycle state of a subscription.
type SubscriptionStatus string

const (
	StatusActive    SubscriptionStatus = "active"
	StatusPaused    SubscriptionStatus = "paused"
	StatusCancelled SubscriptionStatus = "cancelled"
	StatusExpired   SubscriptionStatus = "expired"
)

// BillingRecordStatus represents the outcome of a billing attempt.
const (
	BillingSuccess = "success"
	BillingFailed  = "failed"
	BillingPending = "pending"
)

// Subscription is the aggregate root for recurring billing contracts.
type Subscription struct {
	ID              kernel.SubscriptionID
	TenantID        kernel.TenantID
	CustomerID      kernel.CustomerID
	ProductID       kernel.ProductID
	VariantID       *kernel.VariantID
	Price           kernel.Money
	Interval        BillingInterval
	Status          SubscriptionStatus
	NextBillingDate time.Time
	LastBilledAt    *time.Time
	CancelledAt     *time.Time
	PausedAt        *time.Time
	TrialEndsAt     *time.Time
	Metadata        map[string]string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// BillingRecord captures the result of a single billing attempt.
type BillingRecord struct {
	ID             kernel.BillingRecordID
	SubscriptionID kernel.SubscriptionID
	TenantID       kernel.TenantID
	Amount         kernel.Money
	Status         string
	OrderID        *kernel.OrderID
	FailureReason  *string
	BilledAt       time.Time
	CreatedAt      time.Time
}
