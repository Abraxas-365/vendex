package order

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines persistence operations for orders.
type Repository interface {
	Create(ctx context.Context, o *Order) error
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.OrderID) (*Order, error)
	Update(ctx context.Context, o *Order) error
	UpdateCheckoutFields(ctx context.Context, o *Order) error
	List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[Order], error)
	ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, pg kernel.PaginationOptions) (kernel.Paginated[Order], error)
}
