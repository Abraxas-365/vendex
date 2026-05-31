package approval

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines persistence operations for ApprovalRequest records.
type Repository interface {
	Create(ctx context.Context, req ApprovalRequest) (ApprovalRequest, error)
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ApprovalRequestID) (ApprovalRequest, error)
	List(ctx context.Context, tenantID kernel.TenantID, status string, p kernel.PaginationOptions) (kernel.Paginated[ApprovalRequest], error)
	UpdateStatus(ctx context.Context, tenantID kernel.TenantID, id kernel.ApprovalRequestID, status, reason, reviewedBy string) (ApprovalRequest, error)
	CountPending(ctx context.Context, tenantID kernel.TenantID) (int, error)
}
