package customergroup

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Repository defines the persistence interface for the customer group domain.
type Repository interface {
	Create(ctx context.Context, g *CustomerGroup) error
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerGroupID) (*CustomerGroup, error)
	List(ctx context.Context, tenantID kernel.TenantID) ([]CustomerGroup, error)
	Update(ctx context.Context, g *CustomerGroup) error
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerGroupID) error

	AddMember(ctx context.Context, m *GroupMembership) error
	RemoveMember(ctx context.Context, tenantID kernel.TenantID, groupID kernel.CustomerGroupID, customerID kernel.CustomerID) error
	ListMembers(ctx context.Context, tenantID kernel.TenantID, groupID kernel.CustomerGroupID) ([]GroupMembership, error)
	GetCustomerGroups(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) ([]CustomerGroup, error)
}
