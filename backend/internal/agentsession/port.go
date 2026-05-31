package agentsession

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// SessionRepository defines persistence operations for Session records.
type SessionRepository interface {
	Create(ctx context.Context, s Session) (Session, error)
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentSessionID) (Session, error)
	Update(ctx context.Context, s Session) (Session, error)
	ListByTenant(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[Session], error)
	ListActive(ctx context.Context, tenantID kernel.TenantID) ([]Session, error)
}

// ChatRepository defines persistence operations for ChatMessage records.
type ChatRepository interface {
	SaveMessage(ctx context.Context, msg ChatMessage) error
	ListMessages(ctx context.Context, sessionID kernel.AgentSessionID, p kernel.PaginationOptions) (kernel.Paginated[ChatMessage], error)
}
