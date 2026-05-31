// Package agenttriggersrv implements business logic for event-triggered agent actions.
package agenttriggersrv

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/Abraxas-365/vendex/internal/agenttrigger/agenttrigger"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// ValidEventTypes lists all supported eventbus event types available for trigger configuration.
var ValidEventTypes = []string{
	string(eventbus.OrderPlaced),
	string(eventbus.OrderConfirmed),
	string(eventbus.OrderShipped),
	string(eventbus.OrderDelivered),
	string(eventbus.OrderCancelled),
	string(eventbus.CustomerRegistered),
	string(eventbus.CustomerUpdated),
	string(eventbus.ProductCreated),
	string(eventbus.ProductUpdated),
	string(eventbus.ProductDeleted),
	string(eventbus.CategoryCreated),
	string(eventbus.CollectionCreated),
	string(eventbus.CollectionUpdated),
	string(eventbus.PagePublished),
	string(eventbus.PageUnpublished),
	string(eventbus.PluginInstalled),
	string(eventbus.PluginUninstalled),
	string(eventbus.ThemeActivated),
	string(eventbus.ThemeUpdated),
	string(eventbus.SettingsUpdated),
	string(eventbus.CartCreated),
	string(eventbus.CartUpdated),
	string(eventbus.CartAbandoned),
	string(eventbus.CheckoutStarted),
	string(eventbus.CheckoutCompleted),
	string(eventbus.CheckoutFailed),
	string(eventbus.PaymentCreated),
	string(eventbus.PaymentCompleted),
	string(eventbus.PaymentFailed),
	string(eventbus.RefundCreated),
	string(eventbus.RefundCompleted),
	string(eventbus.ShippingZoneCreated),
	string(eventbus.ShippingRateCreated),
	string(eventbus.TaxRateCreated),
	string(eventbus.GiftCardCreated),
	string(eventbus.GiftCardRedeemed),
	string(eventbus.SubscriptionCreated),
	string(eventbus.SubscriptionCancelled),
	string(eventbus.SubscriptionBilled),
	string(eventbus.StockUpdated),
	string(eventbus.StockLowAlert),
	string(eventbus.ReviewCreated),
	string(eventbus.ReviewApproved),
	string(eventbus.ReviewRejected),
	string(eventbus.ReturnRequested),
	string(eventbus.ReturnApproved),
	string(eventbus.ReturnCompleted),
	string(eventbus.LoyaltyPointsEarned),
	string(eventbus.LoyaltyPointsRedeemed),
	string(eventbus.LoyaltyRewardCreated),
	string(eventbus.BundleCreated),
	string(eventbus.BundleUpdated),
	string(eventbus.StorefrontCreated),
	string(eventbus.StorefrontUpdated),
	string(eventbus.StorefrontDeleted),
	string(eventbus.BulkOperationStarted),
	string(eventbus.BulkOperationCompleted),
	string(eventbus.BlogPostPublished),
	string(eventbus.ExperimentStarted),
	string(eventbus.ExperimentCompleted),
}

// validEventTypeSet is a fast lookup set built from ValidEventTypes.
var validEventTypeSet map[string]struct{}

func init() {
	validEventTypeSet = make(map[string]struct{}, len(ValidEventTypes))
	for _, et := range ValidEventTypes {
		validEventTypeSet[et] = struct{}{}
	}
}

// Service manages event trigger lifecycle and execution recording.
type Service struct {
	triggerRepo agenttrigger.TriggerRepository
	logRepo     agenttrigger.TriggerLogRepository
}

// NewService creates a new agenttrigger Service.
func NewService(
	triggerRepo agenttrigger.TriggerRepository,
	logRepo agenttrigger.TriggerLogRepository,
) *Service {
	return &Service{
		triggerRepo: triggerRepo,
		logRepo:     logRepo,
	}
}

// Create creates a new trigger for a tenant.
func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, req agenttrigger.CreateTriggerRequest) (agenttrigger.Trigger, error) {
	if req.Name == "" {
		return agenttrigger.Trigger{}, agenttrigger.ErrInvalidInput
	}
	if req.EventType == "" {
		return agenttrigger.Trigger{}, agenttrigger.ErrInvalidInput
	}
	if _, ok := validEventTypeSet[req.EventType]; !ok {
		return agenttrigger.Trigger{}, agenttrigger.ErrInvalidEventType
	}
	if req.Prompt == "" {
		return agenttrigger.Trigger{}, agenttrigger.ErrInvalidInput
	}
	cooldown := req.Cooldown
	if cooldown <= 0 {
		cooldown = 300
	}

	now := time.Now()
	t := agenttrigger.Trigger{
		ID:        kernel.AgentTriggerID(uuid.New().String()),
		TenantID:  tenantID,
		Name:      req.Name,
		EventType: req.EventType,
		Prompt:    req.Prompt,
		PresetID:  req.PresetID,
		Enabled:   true,
		Cooldown:  cooldown,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.triggerRepo.Create(ctx, t)
}

// Get retrieves a single trigger scoped to a tenant.
func (s *Service) Get(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentTriggerID) (agenttrigger.Trigger, error) {
	return s.triggerRepo.GetByID(ctx, tenantID, id)
}

// Update updates mutable fields on a trigger.
func (s *Service) Update(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentTriggerID, req agenttrigger.UpdateTriggerRequest) (agenttrigger.Trigger, error) {
	t, err := s.triggerRepo.GetByID(ctx, tenantID, id)
	if err != nil {
		return agenttrigger.Trigger{}, err
	}

	if req.Name != nil {
		t.Name = *req.Name
	}
	if req.Prompt != nil {
		t.Prompt = *req.Prompt
	}
	if req.PresetID != nil {
		t.PresetID = *req.PresetID
	}
	if req.Enabled != nil {
		t.Enabled = *req.Enabled
	}
	if req.Cooldown != nil {
		if *req.Cooldown > 0 {
			t.Cooldown = *req.Cooldown
		}
	}
	t.UpdatedAt = time.Now()

	return s.triggerRepo.Update(ctx, t)
}

// Delete removes a trigger.
func (s *Service) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentTriggerID) error {
	return s.triggerRepo.Delete(ctx, tenantID, id)
}

// List returns paginated triggers for a tenant.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[agenttrigger.Trigger], error) {
	return s.triggerRepo.List(ctx, tenantID, p)
}

// Enable enables a trigger.
func (s *Service) Enable(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentTriggerID) (agenttrigger.Trigger, error) {
	t, err := s.triggerRepo.GetByID(ctx, tenantID, id)
	if err != nil {
		return agenttrigger.Trigger{}, err
	}
	t.Enabled = true
	t.UpdatedAt = time.Now()
	return s.triggerRepo.Update(ctx, t)
}

// Disable disables a trigger.
func (s *Service) Disable(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentTriggerID) (agenttrigger.Trigger, error) {
	t, err := s.triggerRepo.GetByID(ctx, tenantID, id)
	if err != nil {
		return agenttrigger.Trigger{}, err
	}
	t.Enabled = false
	t.UpdatedAt = time.Now()
	return s.triggerRepo.Update(ctx, t)
}

// GetTriggersForEvent returns all enabled triggers for an event type (cross-tenant).
// Used by the event handler to fan out to all tenant triggers.
func (s *Service) GetTriggersForEvent(ctx context.Context, eventType string) ([]agenttrigger.Trigger, error) {
	return s.triggerRepo.ListByEventType(ctx, eventType)
}

// RecordExecution logs a trigger execution.
func (s *Service) RecordExecution(
	ctx context.Context,
	triggerID kernel.AgentTriggerID,
	tenantID kernel.TenantID,
	eventType string,
	payload json.RawMessage,
	response string,
	status string,
) (agenttrigger.TriggerLog, error) {
	if len(payload) == 0 {
		payload = json.RawMessage("{}")
	}
	log := agenttrigger.TriggerLog{
		ID:            kernel.TriggerLogID(uuid.New().String()),
		TriggerID:     triggerID,
		TenantID:      tenantID,
		EventType:     eventType,
		EventPayload:  payload,
		AgentResponse: response,
		Status:        status,
		CreatedAt:     time.Now(),
	}
	return s.logRepo.Create(ctx, log)
}

// GetLogs returns paginated execution logs for a specific trigger.
func (s *Service) GetLogs(ctx context.Context, tenantID kernel.TenantID, triggerID kernel.AgentTriggerID, p kernel.PaginationOptions) (kernel.Paginated[agenttrigger.TriggerLog], error) {
	// Verify the trigger belongs to the tenant before returning logs.
	if _, err := s.triggerRepo.GetByID(ctx, tenantID, triggerID); err != nil {
		return kernel.Paginated[agenttrigger.TriggerLog]{}, err
	}
	return s.logRepo.ListByTrigger(ctx, tenantID, triggerID, p)
}

// GetValidEventTypes returns the list of event types that can be used in triggers.
func (s *Service) GetValidEventTypes() []string {
	return ValidEventTypes
}

// UpdateLastFired records the last-fired timestamp for cooldown tracking.
func (s *Service) UpdateLastFired(ctx context.Context, id kernel.AgentTriggerID, firedAt time.Time) error {
	return s.triggerRepo.UpdateLastFired(ctx, id, firedAt)
}
