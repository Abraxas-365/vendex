package ordersrv

import (
	"context"
	"fmt"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/google/uuid"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/order"
)

// Service handles order business logic.
type Service struct {
	repo order.Repository
}

// New creates a new order service.
func New(repo order.Repository) *Service {
	return &Service{repo: repo}
}

// CreateItemInput represents one line item when creating an order.
type CreateItemInput struct {
	ProductID   kernel.ProductID
	ProductName string
	Quantity    int
	UnitPrice   kernel.Money
}

// CreateInput holds all data needed to create an order.
type CreateInput struct {
	CustomerID      kernel.CustomerID
	Items           []CreateItemInput
	ShippingAddress order.Address
}

// Create builds a new order from the given items and persists it.
func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, in CreateInput) (*order.Order, error) {
	if len(in.Items) == 0 {
		return nil, order.ErrEmptyOrder
	}

	now := time.Now()
	items := make([]order.OrderItem, len(in.Items))
	for i, it := range in.Items {
		items[i] = order.OrderItem{
			ID:          kernel.OrderItemID(uuid.NewString()),
			ProductID:   it.ProductID,
			ProductName: it.ProductName,
			Quantity:    it.Quantity,
			UnitPrice:   it.UnitPrice,
		}
	}

	o := &order.Order{
		ID:              kernel.OrderID(uuid.NewString()),
		TenantID:        tenantID,
		CustomerID:      in.CustomerID,
		Items:           items,
		Status:          order.StatusPending,
		ShippingAddress: in.ShippingAddress,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	o.CalculateTotal()

	if err := s.repo.Create(ctx, o); err != nil {
		return nil, errx.Wrap(err, "creating order", errx.TypeInternal)
	}
	return o, nil
}

// GetByID retrieves an order by ID, scoped to tenant.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.OrderID) (*order.Order, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// UpdateStatus transitions the order to a new status.
func (s *Service) UpdateStatus(ctx context.Context, tenantID kernel.TenantID, id kernel.OrderID, newStatus order.OrderStatus) (*order.Order, error) {
	o, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if !o.TransitionTo(newStatus) {
		return nil, order.ErrInvalidTransition
	}

	if err := s.repo.Update(ctx, o); err != nil {
		return nil, errx.Wrap(err, "updating order status", errx.TypeInternal)
	}
	return o, nil
}

// Cancel cancels an order if possible.
func (s *Service) Cancel(ctx context.Context, tenantID kernel.TenantID, id kernel.OrderID) (*order.Order, error) {
	o, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if o.Status == order.StatusCancelled {
		return nil, order.ErrAlreadyCancelled
	}

	if !o.Cancel() {
		return nil, order.ErrInvalidTransition
	}

	if err := s.repo.Update(ctx, o); err != nil {
		return nil, errx.Wrap(err, "cancelling order", errx.TypeInternal)
	}
	return o, nil
}

// List returns a paginated list of orders for a tenant.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[order.Order], error) {
	return s.repo.List(ctx, tenantID, pg)
}

// ListByCustomer returns orders for a specific customer.
func (s *Service) ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, pg kernel.PaginationOptions) (kernel.Paginated[order.Order], error) {
	return s.repo.ListByCustomer(ctx, tenantID, customerID, pg)
}

