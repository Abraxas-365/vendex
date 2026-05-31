package subscriptionsrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/subscription"
	"github.com/google/uuid"
)

// CreateInput holds the data required to create a new subscription.
type CreateInput struct {
	CustomerID  kernel.CustomerID
	ProductID   kernel.ProductID
	VariantID   *kernel.VariantID
	Price       kernel.Money
	Interval    subscription.BillingInterval
	TrialEndsAt *time.Time
	Metadata    map[string]string
}

// Service implements subscription business logic.
type Service struct {
	repo subscription.Repository
	bus  eventbus.Bus
}

// New creates a subscription Service.
func New(repo subscription.Repository, bus eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// Create creates a new subscription and schedules the first billing date.
func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, input CreateInput) (*subscription.Subscription, error) {
	if !input.Interval.IsValid() {
		return nil, subscription.ErrInvalidInterval
	}

	now := time.Now().UTC()

	// First billing date is after trial period (if any), otherwise from now.
	var nextBilling time.Time
	if input.TrialEndsAt != nil && input.TrialEndsAt.After(now) {
		nextBilling = *input.TrialEndsAt
	} else {
		nextBilling = calculateNextBilling(now, input.Interval)
	}

	meta := input.Metadata
	if meta == nil {
		meta = map[string]string{}
	}

	sub := &subscription.Subscription{
		ID:              kernel.NewSubscriptionID(uuid.NewString()),
		TenantID:        tenantID,
		CustomerID:      input.CustomerID,
		ProductID:       input.ProductID,
		VariantID:       input.VariantID,
		Price:           input.Price,
		Interval:        input.Interval,
		Status:          subscription.StatusActive,
		NextBillingDate: nextBilling,
		TrialEndsAt:     input.TrialEndsAt,
		Metadata:        meta,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, err
	}

	_ = s.publishEvent(ctx, eventbus.SubscriptionCreated, tenantID, sub)

	return sub, nil
}

// GetByID returns a subscription by ID.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.SubscriptionID) (*subscription.Subscription, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// Cancel marks a subscription as cancelled.
func (s *Service) Cancel(ctx context.Context, tenantID kernel.TenantID, id kernel.SubscriptionID) (*subscription.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if sub.Status == subscription.StatusCancelled {
		return nil, subscription.ErrAlreadyCancelled
	}

	now := time.Now().UTC()
	sub.Status = subscription.StatusCancelled
	sub.CancelledAt = &now
	sub.UpdatedAt = now

	if err := s.repo.Update(ctx, sub); err != nil {
		return nil, err
	}

	_ = s.publishEvent(ctx, eventbus.SubscriptionCancelled, tenantID, sub)

	return sub, nil
}

// Pause pauses an active subscription.
func (s *Service) Pause(ctx context.Context, tenantID kernel.TenantID, id kernel.SubscriptionID) (*subscription.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if sub.Status == subscription.StatusPaused {
		return nil, subscription.ErrAlreadyPaused
	}
	if sub.Status != subscription.StatusActive {
		return nil, subscription.ErrNotActive
	}

	now := time.Now().UTC()
	sub.Status = subscription.StatusPaused
	sub.PausedAt = &now
	sub.UpdatedAt = now

	if err := s.repo.Update(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

// Resume resumes a paused subscription and recalculates the next billing date.
func (s *Service) Resume(ctx context.Context, tenantID kernel.TenantID, id kernel.SubscriptionID) (*subscription.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if sub.Status != subscription.StatusPaused {
		return nil, subscription.ErrNotActive
	}

	now := time.Now().UTC()
	sub.Status = subscription.StatusActive
	sub.PausedAt = nil
	sub.NextBillingDate = calculateNextBilling(now, sub.Interval)
	sub.UpdatedAt = now

	if err := s.repo.Update(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

// ListByCustomer returns all subscriptions for a given customer.
func (s *Service) ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) ([]subscription.Subscription, error) {
	return s.repo.ListByCustomer(ctx, tenantID, customerID)
}

// List returns a paginated list of all subscriptions for a tenant.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, page, pageSize int) (kernel.Paginated[subscription.Subscription], error) {
	return s.repo.List(ctx, tenantID, page, pageSize)
}

// ListDueBilling returns all active subscriptions with a billing date at or before now.
func (s *Service) ListDueBilling(ctx context.Context, tenantID kernel.TenantID) ([]subscription.Subscription, error) {
	return s.repo.ListDueBilling(ctx, tenantID, time.Now().UTC())
}

// RecordBilling creates a billing record for a subscription. On success it advances the next billing date.
func (s *Service) RecordBilling(
	ctx context.Context,
	tenantID kernel.TenantID,
	subID kernel.SubscriptionID,
	amount kernel.Money,
	status string,
	orderID *kernel.OrderID,
	failureReason *string,
) (*subscription.BillingRecord, error) {
	sub, err := s.repo.GetByID(ctx, tenantID, subID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	record := &subscription.BillingRecord{
		ID:             kernel.NewBillingRecordID(uuid.NewString()),
		SubscriptionID: subID,
		TenantID:       tenantID,
		Amount:         amount,
		Status:         status,
		OrderID:        orderID,
		FailureReason:  failureReason,
		BilledAt:       now,
		CreatedAt:      now,
	}

	if err := s.repo.CreateBillingRecord(ctx, record); err != nil {
		return nil, err
	}

	// Advance next billing date only on success.
	if status == subscription.BillingSuccess {
		sub.LastBilledAt = &now
		sub.NextBillingDate = calculateNextBilling(now, sub.Interval)
		sub.UpdatedAt = now
		if err := s.repo.Update(ctx, sub); err != nil {
			return nil, err
		}

		_ = s.publishEvent(ctx, eventbus.SubscriptionBilled, tenantID, map[string]any{
			"subscription_id": sub.ID,
			"billing_record":  record,
		})
	}

	return record, nil
}

// ListBillingRecords returns a paginated list of billing records for a subscription.
func (s *Service) ListBillingRecords(ctx context.Context, tenantID kernel.TenantID, subID kernel.SubscriptionID, page, pageSize int) (kernel.Paginated[subscription.BillingRecord], error) {
	return s.repo.ListBillingRecords(ctx, tenantID, subID, page, pageSize)
}

// calculateNextBilling returns the next billing timestamp after `from` for the given interval.
func calculateNextBilling(from time.Time, interval subscription.BillingInterval) time.Time {
	switch interval {
	case subscription.IntervalWeekly:
		return from.AddDate(0, 0, 7)
	case subscription.IntervalMonthly:
		return from.AddDate(0, 1, 0)
	case subscription.IntervalQuarterly:
		return from.AddDate(0, 3, 0)
	case subscription.IntervalYearly:
		return from.AddDate(1, 0, 0)
	default:
		return from.AddDate(0, 1, 0)
	}
}

// publishEvent fires a domain event, ignoring errors (non-critical path).
func (s *Service) publishEvent(ctx context.Context, eventType eventbus.EventType, tenantID kernel.TenantID, payload any) error {
	event, err := eventbus.NewEvent(eventType, tenantID, payload)
	if err != nil {
		return err
	}
	return s.bus.Publish(ctx, event)
}
