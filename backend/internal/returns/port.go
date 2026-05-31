package returns

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Repository defines persistence operations for return requests.
type Repository interface {
	Create(ctx context.Context, r *ReturnRequest) error
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ReturnID) (*ReturnRequest, error)
	Update(ctx context.Context, r *ReturnRequest) error
	List(ctx context.Context, tenantID kernel.TenantID, status string, pg kernel.PaginationOptions) (kernel.Paginated[ReturnRequest], error)
	ListByOrder(ctx context.Context, tenantID kernel.TenantID, orderID kernel.OrderID, pg kernel.PaginationOptions) (kernel.Paginated[ReturnRequest], error)
	ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, pg kernel.PaginationOptions) (kernel.Paginated[ReturnRequest], error)
}
