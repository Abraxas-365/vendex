package webhooksrv

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/webhook"
	"github.com/google/uuid"
)

// Service implements webhook business logic.
type Service struct {
	repo webhook.Repository
}

// New creates a webhook Service.
func New(repo webhook.Repository) *Service {
	return &Service{repo: repo}
}

// ---------------------------------------------------------------------------
// Webhook CRUD
// ---------------------------------------------------------------------------

// Create registers a new webhook endpoint.
func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, input webhook.CreateWebhookInput) (*webhook.Webhook, error) {
	if input.URL == "" {
		return nil, webhook.ErrInvalidURL
	}
	if len(input.Events) == 0 {
		return nil, webhook.ErrNoEvents
	}

	now := time.Now().UTC()
	wh := &webhook.Webhook{
		ID:          kernel.NewWebhookID(uuid.NewString()),
		TenantID:    tenantID,
		URL:         input.URL,
		Secret:      input.Secret,
		Events:      input.Events,
		Active:      true,
		Description: input.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, wh); err != nil {
		return nil, err
	}

	return wh, nil
}

// GetByID returns a webhook by ID scoped to the tenant.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.WebhookID) (*webhook.Webhook, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// List returns a paginated list of webhooks for a tenant.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, page, pageSize int) (kernel.Paginated[webhook.Webhook], error) {
	return s.repo.List(ctx, tenantID, page, pageSize)
}

// Update modifies a webhook's mutable fields.
func (s *Service) Update(ctx context.Context, tenantID kernel.TenantID, id kernel.WebhookID, input webhook.UpdateWebhookInput) (*webhook.Webhook, error) {
	wh, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if input.URL != nil {
		if *input.URL == "" {
			return nil, webhook.ErrInvalidURL
		}
		wh.URL = *input.URL
	}
	if input.Secret != nil {
		wh.Secret = *input.Secret
	}
	if input.Events != nil {
		if len(input.Events) == 0 {
			return nil, webhook.ErrNoEvents
		}
		wh.Events = input.Events
	}
	if input.Description != nil {
		wh.Description = *input.Description
	}

	wh.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, wh); err != nil {
		return nil, err
	}

	return wh, nil
}

// Delete removes a webhook and all its deliveries (cascade).
func (s *Service) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.WebhookID) error {
	_, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, tenantID, id)
}

// Toggle sets the active state of a webhook.
func (s *Service) Toggle(ctx context.Context, tenantID kernel.TenantID, id kernel.WebhookID, active bool) (*webhook.Webhook, error) {
	wh, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if wh.Active == active {
		if active {
			return nil, webhook.ErrAlreadyActive
		}
		return nil, webhook.ErrAlreadyInactive
	}

	wh.Active = active
	wh.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, wh); err != nil {
		return nil, err
	}

	return wh, nil
}

// ---------------------------------------------------------------------------
// Event delivery
// ---------------------------------------------------------------------------

// Deliver finds all active webhooks for the tenant subscribing to eventType,
// creates a pending delivery record for each, and marks them for async dispatch.
// Real HTTP dispatch is intentionally omitted — a background worker would pick
// up "pending" deliveries and execute them.
func (s *Service) Deliver(ctx context.Context, tenantID kernel.TenantID, eventType string, payload json.RawMessage) error {
	webhooks, err := s.repo.ListActiveByEvent(ctx, tenantID, eventType)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	for _, wh := range webhooks {
		delivery := &webhook.WebhookDelivery{
			ID:          kernel.NewWebhookDeliveryID(uuid.NewString()),
			TenantID:    tenantID,
			WebhookID:   wh.ID,
			EventType:   eventType,
			Payload:     payload,
			Status:      webhook.DeliveryPending,
			Attempts:    0,
			MaxAttempts: 3,
			CreatedAt:   now,
		}
		if err := s.repo.CreateDelivery(ctx, delivery); err != nil {
			// Log but continue — don't fail for one bad delivery
			_ = err
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Delivery management
// ---------------------------------------------------------------------------

// GetDelivery returns a single delivery by ID.
func (s *Service) GetDelivery(ctx context.Context, tenantID kernel.TenantID, id kernel.WebhookDeliveryID) (*webhook.WebhookDelivery, error) {
	return s.repo.GetDelivery(ctx, tenantID, id)
}

// ListDeliveries returns paginated deliveries for a specific webhook.
func (s *Service) ListDeliveries(ctx context.Context, tenantID kernel.TenantID, webhookID kernel.WebhookID, page, pageSize int) (kernel.Paginated[webhook.WebhookDelivery], error) {
	// Verify the webhook belongs to the tenant before listing
	_, err := s.repo.GetByID(ctx, tenantID, webhookID)
	if err != nil {
		return kernel.Paginated[webhook.WebhookDelivery]{}, err
	}
	return s.repo.ListDeliveries(ctx, tenantID, webhookID, page, pageSize)
}

// RetryDelivery resets a failed delivery so it will be re-attempted.
func (s *Service) RetryDelivery(ctx context.Context, tenantID kernel.TenantID, id kernel.WebhookDeliveryID) (*webhook.WebhookDelivery, error) {
	delivery, err := s.repo.GetDelivery(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if delivery.Status != webhook.DeliveryFailed {
		return nil, errx.Business("only failed deliveries can be retried")
	}

	now := time.Now().UTC()
	delivery.Status = webhook.DeliveryPending
	delivery.NextRetryAt = &now

	if err := s.repo.UpdateDelivery(ctx, delivery); err != nil {
		return nil, err
	}

	return delivery, nil
}
