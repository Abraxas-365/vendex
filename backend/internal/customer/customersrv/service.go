package customersrv

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Abraxas-365/vendex/internal/customer"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Service handles customer business logic.
type Service struct {
	repo customer.Repository
	bus  eventbus.Bus
}

// New creates a new customer service.
func New(repo customer.Repository, bus eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// CreateInput holds the data needed to create a customer.
type CreateInput struct {
	Email     string
	Name      string
	Phone     string
	Addresses []customer.Address
}

// Create creates a new customer for the given tenant.
func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, in CreateInput) (*customer.Customer, error) {
	email := kernel.NewEmail(in.Email)
	if email.IsEmpty() {
		return nil, customer.ErrInvalidEmail
	}

	// Check for duplicate email.
	existing, err := s.repo.GetByEmail(ctx, tenantID, email)
	if err != nil && !errx.Is(err, customer.ErrNotFound) {
		return nil, errx.Wrap(err, "checking email uniqueness", errx.TypeInternal)
	}
	if existing != nil {
		return nil, customer.ErrDuplicateEmail
	}

	now := time.Now()
	c := &customer.Customer{
		ID:        kernel.CustomerID(uuid.NewString()),
		TenantID:  tenantID,
		Email:     email,
		Name:      in.Name,
		Phone:     in.Phone,
		Addresses: in.Addresses,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if c.Addresses == nil {
		c.Addresses = []customer.Address{}
	}

	if err := s.repo.Create(ctx, c); err != nil {
		return nil, errx.Wrap(err, "creating customer", errx.TypeInternal)
	}

	if evt, err := eventbus.NewEvent(eventbus.CustomerRegistered, tenantID, eventbus.CustomerPayload{
		CustomerID: string(c.ID),
		Email:      string(c.Email),
		Name:       c.Name,
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return c, nil
}

// GetByID retrieves a customer by ID, scoped to tenant.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerID) (*customer.Customer, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// GetByEmail retrieves a customer by email, scoped to tenant.
func (s *Service) GetByEmail(ctx context.Context, tenantID kernel.TenantID, email kernel.Email) (*customer.Customer, error) {
	return s.repo.GetByEmail(ctx, tenantID, email)
}

// Update persists changes to a customer.
func (s *Service) Update(ctx context.Context, c *customer.Customer) error {
	c.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, c); err != nil {
		return err
	}

	if evt, err := eventbus.NewEvent(eventbus.CustomerUpdated, c.TenantID, eventbus.CustomerPayload{
		CustomerID: string(c.ID),
		Email:      string(c.Email),
		Name:       c.Name,
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return nil
}

// Delete removes a customer by ID, scoped to tenant.
func (s *Service) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerID) error {
	return s.repo.Delete(ctx, tenantID, id)
}

// List returns a paginated list of customers for a tenant.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[customer.Customer], error) {
	return s.repo.List(ctx, tenantID, pg)
}

