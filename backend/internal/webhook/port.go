package webhook

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Repository defines the persistence contract for the webhook domain.
type Repository interface {
	// Webhook CRUD
	Create(ctx context.Context, wh *Webhook) error
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.WebhookID) (*Webhook, error)
	List(ctx context.Context, tenantID kernel.TenantID, page, pageSize int) (kernel.Paginated[Webhook], error)
	Update(ctx context.Context, wh *Webhook) error
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.WebhookID) error

	// Active webhook lookup for event delivery
	ListActiveByEvent(ctx context.Context, tenantID kernel.TenantID, eventType string) ([]Webhook, error)

	// Delivery tracking
	CreateDelivery(ctx context.Context, d *WebhookDelivery) error
	GetDelivery(ctx context.Context, tenantID kernel.TenantID, id kernel.WebhookDeliveryID) (*WebhookDelivery, error)
	UpdateDelivery(ctx context.Context, d *WebhookDelivery) error
	ListDeliveries(ctx context.Context, tenantID kernel.TenantID, webhookID kernel.WebhookID, page, pageSize int) (kernel.Paginated[WebhookDelivery], error)
}
