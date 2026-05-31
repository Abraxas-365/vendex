package agenttrigger

import (
	"context"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// TriggerRepository defines persistence operations for Trigger records.
type TriggerRepository interface {
	Create(ctx context.Context, t Trigger) (Trigger, error)
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentTriggerID) (Trigger, error)
	Update(ctx context.Context, t Trigger) (Trigger, error)
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentTriggerID) error
	List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[Trigger], error)
	// ListByEventType is cross-tenant — returns all ENABLED triggers for an event type.
	ListByEventType(ctx context.Context, eventType string) ([]Trigger, error)
	UpdateLastFired(ctx context.Context, id kernel.AgentTriggerID, firedAt time.Time) error
}

// TriggerLogRepository defines persistence operations for TriggerLog records.
type TriggerLogRepository interface {
	Create(ctx context.Context, log TriggerLog) (TriggerLog, error)
	ListByTrigger(ctx context.Context, tenantID kernel.TenantID, triggerID kernel.AgentTriggerID, p kernel.PaginationOptions) (kernel.Paginated[TriggerLog], error)
}
