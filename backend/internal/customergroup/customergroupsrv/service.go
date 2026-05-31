package customergroupsrv

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Abraxas-365/vendex/internal/customergroup"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Service handles business logic for the customer group domain.
type Service struct {
	repo customergroup.Repository
}

// New creates a new customer group service.
func New(repo customergroup.Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new customer group for the given tenant.
func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, req customergroup.CreateGroupRequest) (*customergroup.CustomerGroup, error) {
	if req.Name == "" {
		return nil, customergroup.ErrGroupNotFound
	}

	now := time.Now()
	g := &customergroup.CustomerGroup{
		ID:          kernel.CustomerGroupID(uuid.NewString()),
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		Rules:       req.Rules,
		AutoAssign:  req.AutoAssign,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, g); err != nil {
		return nil, err
	}
	return g, nil
}

// GetByID retrieves a customer group by ID, scoped to tenant.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerGroupID) (*customergroup.CustomerGroup, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// List returns all customer groups for a tenant.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID) ([]customergroup.CustomerGroup, error) {
	return s.repo.List(ctx, tenantID)
}

// Update applies changes to a customer group.
func (s *Service) Update(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerGroupID, req customergroup.UpdateGroupRequest) (*customergroup.CustomerGroup, error) {
	g, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		g.Name = *req.Name
	}
	if req.Description != nil {
		g.Description = *req.Description
	}
	if req.Rules != nil {
		g.Rules = *req.Rules
	}
	if req.AutoAssign != nil {
		g.AutoAssign = *req.AutoAssign
	}
	g.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, g); err != nil {
		return nil, err
	}
	return g, nil
}

// Delete removes a customer group by ID, scoped to tenant.
func (s *Service) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerGroupID) error {
	return s.repo.Delete(ctx, tenantID, id)
}

// AddMember adds a customer to a group.
func (s *Service) AddMember(ctx context.Context, tenantID kernel.TenantID, groupID kernel.CustomerGroupID, customerID kernel.CustomerID) (*customergroup.GroupMembership, error) {
	// Verify the group exists for this tenant.
	if _, err := s.repo.GetByID(ctx, tenantID, groupID); err != nil {
		return nil, err
	}

	m := &customergroup.GroupMembership{
		ID:         kernel.CustomerGroupMembershipID(uuid.NewString()),
		GroupID:    groupID,
		CustomerID: customerID,
		TenantID:   tenantID,
		AssignedAt: time.Now(),
	}

	if err := s.repo.AddMember(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

// RemoveMember removes a customer from a group.
func (s *Service) RemoveMember(ctx context.Context, tenantID kernel.TenantID, groupID kernel.CustomerGroupID, customerID kernel.CustomerID) error {
	// Verify group exists.
	if _, err := s.repo.GetByID(ctx, tenantID, groupID); err != nil {
		return err
	}
	return s.repo.RemoveMember(ctx, tenantID, groupID, customerID)
}

// ListMembers returns all members of a customer group.
func (s *Service) ListMembers(ctx context.Context, tenantID kernel.TenantID, groupID kernel.CustomerGroupID) ([]customergroup.GroupMembership, error) {
	// Verify group exists.
	if _, err := s.repo.GetByID(ctx, tenantID, groupID); err != nil {
		return nil, err
	}
	return s.repo.ListMembers(ctx, tenantID, groupID)
}

// GetCustomerGroups returns all groups a customer belongs to.
func (s *Service) GetCustomerGroups(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) ([]customergroup.CustomerGroup, error) {
	return s.repo.GetCustomerGroups(ctx, tenantID, customerID)
}
