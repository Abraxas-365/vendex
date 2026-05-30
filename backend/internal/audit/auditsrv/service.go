package auditsrv

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Abraxas-365/hada-commerce/internal/audit"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Service provides audit-log business operations.
type Service struct {
	repo audit.Repository
}

// New returns a new audit Service.
func New(repo audit.Repository) *Service {
	return &Service{repo: repo}
}

// Log persists an audit entry for the given input.
func (s *Service) Log(ctx context.Context, input audit.CreateAuditInput) (audit.AuditEntry, error) {
	if string(input.TenantID) == "" {
		return audit.AuditEntry{}, errx.Wrap(audit.ErrInvalidInput, "tenant_id is required", errx.TypeValidation)
	}
	if input.UserID == "" {
		return audit.AuditEntry{}, errx.Wrap(audit.ErrInvalidInput, "user_id is required", errx.TypeValidation)
	}
	if input.Action == "" {
		return audit.AuditEntry{}, errx.Wrap(audit.ErrInvalidInput, "action is required", errx.TypeValidation)
	}
	if input.ResourceType == "" {
		return audit.AuditEntry{}, errx.Wrap(audit.ErrInvalidInput, "resource_type is required", errx.TypeValidation)
	}

	entry := audit.AuditEntry{
		ID:           kernel.AuditEntryID(uuid.NewString()),
		TenantID:     input.TenantID,
		UserID:       input.UserID,
		UserEmail:    input.UserEmail,
		Action:       input.Action,
		ResourceType: input.ResourceType,
		ResourceID:   input.ResourceID,
		Changes:      input.Changes,
		Metadata:     input.Metadata,
		IPAddress:    input.IPAddress,
		UserAgent:    input.UserAgent,
		CreatedAt:    time.Now().UTC(),
	}

	return s.repo.Create(ctx, entry)
}

// GetByID retrieves a single audit entry by its ID within the tenant.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.AuditEntryID) (audit.AuditEntry, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// List returns a paginated, filtered audit log for the tenant.
func (s *Service) List(
	ctx context.Context,
	tenantID kernel.TenantID,
	filter audit.AuditFilter,
	page, pageSize int,
) (kernel.Paginated[audit.AuditEntry], error) {
	p := kernel.NewPaginationOptions(page, pageSize)
	return s.repo.List(ctx, tenantID, filter, p)
}

// GetStats returns counts of audit entries grouped by action for the given period.
func (s *Service) GetStats(ctx context.Context, tenantID kernel.TenantID, from, to time.Time) ([]audit.ActionStats, error) {
	return s.repo.CountByAction(ctx, tenantID, from, to)
}
