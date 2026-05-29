package customer

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines persistence operations for customers.
type Repository interface {
	Create(ctx context.Context, c *Customer) error
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerID) (*Customer, error)
	GetByEmail(ctx context.Context, tenantID kernel.TenantID, email kernel.Email) (*Customer, error)
	Update(ctx context.Context, c *Customer) error
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerID) error
	List(ctx context.Context, tenantID kernel.TenantID, pg kernel.Pagination) (kernel.PaginatedResult[Customer], error)
}
