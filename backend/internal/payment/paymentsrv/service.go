package paymentsrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/payment"
	"github.com/google/uuid"
)

// Service handles payment business logic.
type Service struct {
	repo      payment.Repository
	bus       eventbus.Bus
	providers map[string]payment.PaymentProvider
}

// New creates a new payment service.
func New(repo payment.Repository, bus eventbus.Bus, providers map[string]payment.PaymentProvider) *Service {
	return &Service{
		repo:      repo,
		bus:       bus,
		providers: providers,
	}
}

// CreatePayment creates a new pending payment for an order.
func (s *Service) CreatePayment(
	ctx context.Context,
	tenantID kernel.TenantID,
	orderID kernel.OrderID,
	amount int64,
	currency string,
	providerName string,
	method string,
) (*payment.Payment, error) {
	if amount <= 0 {
		return nil, payment.ErrInvalidAmount
	}

	now := time.Now().UTC()
	p := &payment.Payment{
		ID:        kernel.PaymentID(uuid.NewString()),
		TenantID:  tenantID,
		OrderID:   orderID,
		Amount:    kernel.NewMoney(amount, currency),
		Status:    payment.PaymentStatusPending,
		Provider:  providerName,
		Method:    method,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.CreatePayment(ctx, p); err != nil {
		return nil, errx.Wrap(err, "creating payment", errx.TypeInternal)
	}

	if evt, err := eventbus.NewEvent(eventbus.PaymentCreated, tenantID, eventbus.PaymentPayload{
		PaymentID: string(p.ID),
		OrderID:   string(p.OrderID),
		Amount:    p.Amount.Amount,
		Currency:  p.Amount.Currency,
		Provider:  p.Provider,
		Status:    string(p.Status),
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return p, nil
}

// ProcessPayment charges via the configured provider and updates payment status.
func (s *Service) ProcessPayment(
	ctx context.Context,
	tenantID kernel.TenantID,
	paymentID kernel.PaymentID,
	token string,
) (*payment.Payment, error) {
	p, err := s.repo.GetPaymentByID(ctx, tenantID, paymentID)
	if err != nil {
		return nil, err
	}

	if p.Status == payment.PaymentStatusCompleted {
		return nil, payment.ErrAlreadyPaid
	}

	prov, ok := s.providers[p.Provider]
	if !ok {
		return nil, errx.New("unknown payment provider: "+p.Provider, errx.TypeBusiness)
	}

	// Mark as processing
	p.Status = payment.PaymentStatusProcessing
	p.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdatePayment(ctx, p); err != nil {
		return nil, errx.Wrap(err, "updating payment to processing", errx.TypeInternal)
	}

	result, chargeErr := prov.Charge(ctx, p.Amount, token, map[string]string{
		"payment_id": string(p.ID),
		"order_id":   string(p.OrderID),
		"tenant_id":  string(p.TenantID),
	})

	now := time.Now().UTC()
	if chargeErr != nil {
		p.Status = payment.PaymentStatusFailed
		p.ErrorMessage = chargeErr.Error()
		p.UpdatedAt = now
		_ = s.repo.UpdatePayment(ctx, p)

		if evt, err := eventbus.NewEvent(eventbus.PaymentFailed, tenantID, eventbus.PaymentPayload{
			PaymentID: string(p.ID),
			OrderID:   string(p.OrderID),
			Amount:    p.Amount.Amount,
			Currency:  p.Amount.Currency,
			Provider:  p.Provider,
			Status:    string(p.Status),
		}); err == nil {
			_ = s.bus.Publish(ctx, evt)
		}

		return p, errx.Wrap(chargeErr, "payment provider charge failed", errx.TypeExternal)
	}

	p.Status = payment.PaymentStatusCompleted
	p.ProviderPaymentID = result.ProviderID
	p.ProviderData = result.Data
	p.PaidAt = &now
	p.UpdatedAt = now
	if err := s.repo.UpdatePayment(ctx, p); err != nil {
		return nil, errx.Wrap(err, "updating payment after charge", errx.TypeInternal)
	}

	if evt, err := eventbus.NewEvent(eventbus.PaymentCompleted, tenantID, eventbus.PaymentPayload{
		PaymentID: string(p.ID),
		OrderID:   string(p.OrderID),
		Amount:    p.Amount.Amount,
		Currency:  p.Amount.Currency,
		Provider:  p.Provider,
		Status:    string(p.Status),
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return p, nil
}

// GetPayment retrieves a payment by ID, scoped to tenant.
func (s *Service) GetPayment(
	ctx context.Context,
	tenantID kernel.TenantID,
	paymentID kernel.PaymentID,
) (*payment.Payment, error) {
	return s.repo.GetPaymentByID(ctx, tenantID, paymentID)
}

// GetPaymentByOrder retrieves the latest payment for an order.
func (s *Service) GetPaymentByOrder(
	ctx context.Context,
	tenantID kernel.TenantID,
	orderID kernel.OrderID,
) (*payment.Payment, error) {
	return s.repo.GetPaymentByOrder(ctx, tenantID, orderID)
}

// CreateRefund creates and processes a refund against a completed payment.
func (s *Service) CreateRefund(
	ctx context.Context,
	tenantID kernel.TenantID,
	paymentID kernel.PaymentID,
	amount int64,
	reason string,
) (*payment.Refund, error) {
	p, err := s.repo.GetPaymentByID(ctx, tenantID, paymentID)
	if err != nil {
		return nil, err
	}

	if amount <= 0 {
		return nil, payment.ErrInvalidAmount
	}
	if amount > p.Amount.Amount {
		return nil, payment.ErrRefundExceedsPayment
	}

	prov, ok := s.providers[p.Provider]
	if !ok {
		return nil, errx.New("unknown payment provider: "+p.Provider, errx.TypeBusiness)
	}

	now := time.Now().UTC()
	r := &payment.Refund{
		ID:        kernel.RefundID(uuid.NewString()),
		TenantID:  tenantID,
		PaymentID: paymentID,
		OrderID:   p.OrderID,
		Amount:    kernel.NewMoney(amount, p.Amount.Currency),
		Reason:    reason,
		Status:    payment.RefundStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.CreateRefund(ctx, r); err != nil {
		return nil, errx.Wrap(err, "creating refund", errx.TypeInternal)
	}

	if evt, err := eventbus.NewEvent(eventbus.RefundCreated, tenantID, eventbus.RefundPayload{
		RefundID:  string(r.ID),
		PaymentID: string(r.PaymentID),
		OrderID:   string(r.OrderID),
		Amount:    r.Amount.Amount,
		Currency:  r.Amount.Currency,
		Status:    string(r.Status),
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	result, refundErr := prov.Refund(ctx, p.ProviderPaymentID, r.Amount)

	now = time.Now().UTC()
	if refundErr != nil {
		r.Status = payment.RefundStatusFailed
		r.UpdatedAt = now
		_ = s.repo.UpdateRefund(ctx, r)
		return r, errx.Wrap(refundErr, "payment provider refund failed", errx.TypeExternal)
	}

	r.Status = payment.RefundStatusCompleted
	r.ProviderRefundID = result.ProviderID
	r.UpdatedAt = now
	if err := s.repo.UpdateRefund(ctx, r); err != nil {
		return nil, errx.Wrap(err, "updating refund after provider response", errx.TypeInternal)
	}

	if evt, err := eventbus.NewEvent(eventbus.RefundCompleted, tenantID, eventbus.RefundPayload{
		RefundID:  string(r.ID),
		PaymentID: string(r.PaymentID),
		OrderID:   string(r.OrderID),
		Amount:    r.Amount.Amount,
		Currency:  r.Amount.Currency,
		Status:    string(r.Status),
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return r, nil
}

// ListRefunds returns all refunds for a given payment.
func (s *Service) ListRefunds(
	ctx context.Context,
	tenantID kernel.TenantID,
	paymentID kernel.PaymentID,
) ([]payment.Refund, error) {
	return s.repo.ListRefundsByPayment(ctx, tenantID, paymentID)
}

// ListPaymentsByOrder returns all payments for a given order.
func (s *Service) ListPaymentsByOrder(
	ctx context.Context,
	tenantID kernel.TenantID,
	orderID kernel.OrderID,
) ([]payment.Payment, error) {
	return s.repo.ListPaymentsByOrder(ctx, tenantID, orderID)
}
