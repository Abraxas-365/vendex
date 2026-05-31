package returnssrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/returns"
	"github.com/google/uuid"
)

// Service handles return request business logic.
type Service struct {
	repo returns.Repository
	bus  eventbus.Bus
}

// New creates a new returns service.
func New(repo returns.Repository, bus eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// CreateReturn creates a new return request.
func (s *Service) CreateReturn(ctx context.Context, tenantID kernel.TenantID, in returns.CreateReturnInput) (*returns.ReturnRequest, error) {
	if len(in.Items) == 0 {
		return nil, returns.ErrNoItems
	}
	if in.Reason == "" {
		return nil, returns.ErrInvalidInput
	}

	now := time.Now()
	items := make([]returns.ReturnItem, len(in.Items))
	for i, it := range in.Items {
		cond := it.Condition
		if cond == "" {
			cond = returns.ConditionUnopened
		}
		items[i] = returns.ReturnItem{
			ID:        kernel.ReturnItemID(uuid.NewString()),
			TenantID:  tenantID,
			ProductID: it.ProductID,
			VariantID: it.VariantID,
			Quantity:  it.Quantity,
			Reason:    it.Reason,
			Condition: cond,
			CreatedAt: now,
		}
	}

	r := &returns.ReturnRequest{
		ID:           kernel.ReturnID(uuid.NewString()),
		TenantID:     tenantID,
		OrderID:      in.OrderID,
		CustomerID:   in.CustomerID,
		Status:       returns.StatusRequested,
		Reason:       in.Reason,
		Notes:        in.Notes,
		RefundAmount: kernel.NewMoney(0, "USD"),
		Items:        items,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, r); err != nil {
		return nil, errx.Wrap(err, "creating return request", errx.TypeInternal)
	}

	if evt, err := eventbus.NewEvent(eventbus.ReturnRequested, tenantID, eventbus.ReturnPayload{
		ReturnID:   string(r.ID),
		OrderID:    string(r.OrderID),
		CustomerID: string(r.CustomerID),
		Status:     string(r.Status),
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return r, nil
}

// GetByID retrieves a return request by ID, scoped to tenant.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ReturnID) (*returns.ReturnRequest, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// List returns a paginated list of return requests for a tenant, optionally filtered by status.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, status string, page, pageSize int) (kernel.Paginated[returns.ReturnRequest], error) {
	pg := kernel.NewPaginationOptions(page, pageSize)
	return s.repo.List(ctx, tenantID, status, pg)
}

// ListByOrder returns return requests for a specific order.
func (s *Service) ListByOrder(ctx context.Context, tenantID kernel.TenantID, orderID kernel.OrderID, page, pageSize int) (kernel.Paginated[returns.ReturnRequest], error) {
	pg := kernel.NewPaginationOptions(page, pageSize)
	return s.repo.ListByOrder(ctx, tenantID, orderID, pg)
}

// ListByCustomer returns return requests for a specific customer.
func (s *Service) ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, page, pageSize int) (kernel.Paginated[returns.ReturnRequest], error) {
	pg := kernel.NewPaginationOptions(page, pageSize)
	return s.repo.ListByCustomer(ctx, tenantID, customerID, pg)
}

// Approve approves a return request.
func (s *Service) Approve(ctx context.Context, tenantID kernel.TenantID, id kernel.ReturnID, in returns.ApproveInput) (*returns.ReturnRequest, error) {
	r, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if r.Status != returns.StatusRequested {
		return nil, returns.ErrInvalidStatus
	}

	currency := in.RefundCurrency
	if currency == "" {
		currency = "USD"
	}

	r.Status = returns.StatusApproved
	r.AdminNotes = in.AdminNotes
	r.Resolution = in.Resolution
	r.RefundAmount = kernel.NewMoney(in.RefundCents, currency)
	r.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, r); err != nil {
		return nil, errx.Wrap(err, "approving return request", errx.TypeInternal)
	}

	if evt, err := eventbus.NewEvent(eventbus.ReturnApproved, tenantID, eventbus.ReturnPayload{
		ReturnID:   string(r.ID),
		OrderID:    string(r.OrderID),
		CustomerID: string(r.CustomerID),
		Status:     string(r.Status),
		Resolution: string(r.Resolution),
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return r, nil
}

// Reject rejects a return request.
func (s *Service) Reject(ctx context.Context, tenantID kernel.TenantID, id kernel.ReturnID, adminNotes string) (*returns.ReturnRequest, error) {
	r, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if r.Status != returns.StatusRequested {
		return nil, returns.ErrInvalidStatus
	}

	r.Status = returns.StatusRejected
	r.AdminNotes = adminNotes
	r.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, r); err != nil {
		return nil, errx.Wrap(err, "rejecting return request", errx.TypeInternal)
	}

	return r, nil
}

// MarkReceived marks the return as received.
func (s *Service) MarkReceived(ctx context.Context, tenantID kernel.TenantID, id kernel.ReturnID) (*returns.ReturnRequest, error) {
	r, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if r.Status != returns.StatusApproved {
		return nil, returns.ErrInvalidStatus
	}

	r.Status = returns.StatusReceived
	r.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, r); err != nil {
		return nil, errx.Wrap(err, "marking return as received", errx.TypeInternal)
	}

	return r, nil
}

// MarkRefunded marks the return as refunded.
func (s *Service) MarkRefunded(ctx context.Context, tenantID kernel.TenantID, id kernel.ReturnID) (*returns.ReturnRequest, error) {
	r, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if r.Status != returns.StatusReceived {
		return nil, returns.ErrInvalidStatus
	}

	r.Status = returns.StatusRefunded
	r.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, r); err != nil {
		return nil, errx.Wrap(err, "marking return as refunded", errx.TypeInternal)
	}

	if evt, err := eventbus.NewEvent(eventbus.ReturnCompleted, tenantID, eventbus.ReturnPayload{
		ReturnID:   string(r.ID),
		OrderID:    string(r.OrderID),
		CustomerID: string(r.CustomerID),
		Status:     string(r.Status),
		Resolution: string(r.Resolution),
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return r, nil
}

// Close closes a return request.
func (s *Service) Close(ctx context.Context, tenantID kernel.TenantID, id kernel.ReturnID) (*returns.ReturnRequest, error) {
	r, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if r.Status == returns.StatusClosed {
		return nil, returns.ErrAlreadyClosed
	}
	if r.Status == returns.StatusRejected {
		return nil, returns.ErrAlreadyRejected
	}

	r.Status = returns.StatusClosed
	r.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, r); err != nil {
		return nil, errx.Wrap(err, "closing return request", errx.TypeInternal)
	}

	return r, nil
}
