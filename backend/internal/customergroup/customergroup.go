package customergroup

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// GroupRules defines the automatic assignment rules for a customer group.
type GroupRules struct {
	MinOrders        *int       `json:"min_orders,omitempty"`
	MinSpent         *int64     `json:"min_spent,omitempty"`
	RegisteredBefore *time.Time `json:"registered_before,omitempty"`
	RegisteredAfter  *time.Time `json:"registered_after,omitempty"`
	Tags             []string   `json:"tags,omitempty"`
}

// CustomerGroup represents a segment of customers for promo targeting and tiered pricing.
type CustomerGroup struct {
	ID          kernel.CustomerGroupID `json:"id" db:"id"`
	TenantID    kernel.TenantID        `json:"tenant_id" db:"tenant_id"`
	Name        string                 `json:"name" db:"name"`
	Description string                 `json:"description" db:"description"`
	Rules       GroupRules             `json:"rules"`
	AutoAssign  bool                   `json:"auto_assign" db:"auto_assign"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// GroupMembership represents the assignment of a customer to a group.
type GroupMembership struct {
	ID         kernel.CustomerGroupMembershipID `json:"id" db:"id"`
	GroupID    kernel.CustomerGroupID           `json:"group_id" db:"group_id"`
	CustomerID kernel.CustomerID                `json:"customer_id" db:"customer_id"`
	TenantID   kernel.TenantID                  `json:"tenant_id" db:"tenant_id"`
	AssignedAt time.Time                        `json:"assigned_at" db:"assigned_at"`
}

// CreateGroupRequest holds input data for creating a customer group.
type CreateGroupRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Rules       GroupRules `json:"rules"`
	AutoAssign  bool       `json:"auto_assign"`
}

// UpdateGroupRequest holds input data for updating a customer group.
type UpdateGroupRequest struct {
	Name        *string     `json:"name,omitempty"`
	Description *string     `json:"description,omitempty"`
	Rules       *GroupRules `json:"rules,omitempty"`
	AutoAssign  *bool       `json:"auto_assign,omitempty"`
}
