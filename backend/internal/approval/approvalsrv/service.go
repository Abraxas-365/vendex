// Package approvalsrv implements business logic for the approval workflow domain.
package approvalsrv

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/Abraxas-365/vendex/internal/approval"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Service handles approval request lifecycle.
type Service struct {
	repo approval.Repository
}

// NewService creates a new approval Service.
func NewService(repo approval.Repository) *Service {
	return &Service{repo: repo}
}

// RequestApproval creates a new pending approval request for a tool execution.
func (s *Service) RequestApproval(
	ctx context.Context,
	tenantID kernel.TenantID,
	toolName string,
	toolInput json.RawMessage,
	requestedBy string,
	sessionID string,
) (approval.ApprovalRequest, error) {
	if len(toolInput) == 0 {
		toolInput = json.RawMessage("{}")
	}

	req := approval.ApprovalRequest{
		ID:          kernel.ApprovalRequestID(uuid.New().String()),
		TenantID:    tenantID,
		SessionID:   sessionID,
		ToolName:    toolName,
		ToolInput:   toolInput,
		Status:      approval.StatusPending,
		Reason:      "",
		RequestedBy: requestedBy,
		ReviewedBy:  "",
		CreatedAt:   time.Now(),
		ReviewedAt:  nil,
	}

	return s.repo.Create(ctx, req)
}

// GetByID retrieves a single approval request scoped to the tenant.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ApprovalRequestID) (approval.ApprovalRequest, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// List returns paginated approval requests, optionally filtered by status.
// Pass an empty status to return all statuses.
func (s *Service) List(
	ctx context.Context,
	tenantID kernel.TenantID,
	status string,
	p kernel.PaginationOptions,
) (kernel.Paginated[approval.ApprovalRequest], error) {
	return s.repo.List(ctx, tenantID, status, p)
}

// ListPending returns paginated pending approval requests.
func (s *Service) ListPending(
	ctx context.Context,
	tenantID kernel.TenantID,
	p kernel.PaginationOptions,
) (kernel.Paginated[approval.ApprovalRequest], error) {
	return s.repo.List(ctx, tenantID, approval.StatusPending, p)
}

// Approve marks an approval request as approved.
// Returns ErrAlreadyReviewed if the request is not in pending state.
func (s *Service) Approve(
	ctx context.Context,
	tenantID kernel.TenantID,
	id kernel.ApprovalRequestID,
	reviewedBy string,
	reason string,
) (approval.ApprovalRequest, error) {
	existing, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return approval.ApprovalRequest{}, err
	}
	if existing.Status != approval.StatusPending {
		return approval.ApprovalRequest{}, approval.ErrAlreadyReviewed
	}

	return s.repo.UpdateStatus(ctx, tenantID, id, approval.StatusApproved, reason, reviewedBy)
}

// Reject marks an approval request as rejected.
// Returns ErrAlreadyReviewed if the request is not in pending state.
func (s *Service) Reject(
	ctx context.Context,
	tenantID kernel.TenantID,
	id kernel.ApprovalRequestID,
	reviewedBy string,
	reason string,
) (approval.ApprovalRequest, error) {
	existing, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return approval.ApprovalRequest{}, err
	}
	if existing.Status != approval.StatusPending {
		return approval.ApprovalRequest{}, approval.ErrAlreadyReviewed
	}

	return s.repo.UpdateStatus(ctx, tenantID, id, approval.StatusRejected, reason, reviewedBy)
}

// CountPending returns the number of pending approval requests for a tenant.
func (s *Service) CountPending(ctx context.Context, tenantID kernel.TenantID) (int, error) {
	return s.repo.CountPending(ctx, tenantID)
}
