package customersrv

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/customer"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
)

// Service handles customer business logic.
type Service struct {
	repo customer.Repository
}

// New creates a new customer service.
func New(repo customer.Repository) *Service {
	return &Service{repo: repo}
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
	email, err := kernel.NewEmail(in.Email)
	if err != nil {
		return nil, customer.ErrInvalidEmail
	}

	// Check for duplicate email.
	existing, err := s.repo.GetByEmail(ctx, tenantID, email)
	if err != nil && !errx.Is(err, customer.ErrNotFound) {
		return nil, fmt.Errorf("checking email uniqueness: %w", err)
	}
	if existing != nil {
		return nil, customer.ErrDuplicateEmail
	}

	now := time.Now()
	c := &customer.Customer{
		ID:        kernel.CustomerID(generateID()),
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
		return nil, fmt.Errorf("creating customer: %w", err)
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
	return s.repo.Update(ctx, c)
}

// Delete removes a customer by ID, scoped to tenant.
func (s *Service) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerID) error {
	return s.repo.Delete(ctx, tenantID, id)
}

// List returns a paginated list of customers for a tenant.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, pg kernel.Pagination) (kernel.PaginatedResult[customer.Customer], error) {
	return s.repo.List(ctx, tenantID, pg)
}

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
