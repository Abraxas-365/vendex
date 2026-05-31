package eventbus

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/google/uuid"
)

// EventType identifies the kind of domain event.
type EventType string

// Standard commerce event types.
const (
	// Order events
	OrderPlaced    EventType = "order.placed"
	OrderConfirmed EventType = "order.confirmed"
	OrderShipped   EventType = "order.shipped"
	OrderDelivered EventType = "order.delivered"
	OrderCancelled EventType = "order.cancelled"

	// Customer events
	CustomerRegistered EventType = "customer.registered"
	CustomerUpdated    EventType = "customer.updated"

	// Product events
	ProductCreated EventType = "product.created"
	ProductUpdated EventType = "product.updated"
	ProductDeleted EventType = "product.deleted"

	// Catalog events
	CategoryCreated   EventType = "category.created"
	CollectionUpdated EventType = "collection.updated"

	// Storefront events
	PagePublished   EventType = "page.published"
	PageUnpublished EventType = "page.unpublished"

	// Plugin events
	PluginInstalled   EventType = "plugin.installed"
	PluginUninstalled EventType = "plugin.uninstalled"

	// Theme events
	ThemeActivated EventType = "theme.activated"
	ThemeUpdated   EventType = "theme.updated"

	// Settings events
	SettingsUpdated EventType = "settings.updated"

	// Cart events
	CartCreated   EventType = "cart.created"
	CartUpdated   EventType = "cart.updated"
	CartAbandoned EventType = "cart.abandoned"

	// Checkout events
	CheckoutStarted   EventType = "checkout.started"
	CheckoutCompleted EventType = "checkout.completed"
	CheckoutFailed    EventType = "checkout.failed"

	// Payment events
	PaymentCreated   EventType = "payment.created"
	PaymentCompleted EventType = "payment.completed"
	PaymentFailed    EventType = "payment.failed"
	RefundCreated    EventType = "refund.created"
	RefundCompleted  EventType = "refund.completed"

	// Shipping events
	ShippingZoneCreated EventType = "shipping_zone.created"
	ShippingRateCreated EventType = "shipping_rate.created"

	// Tax events
	TaxRateCreated EventType = "tax_rate.created"

	// Gift card events
	GiftCardCreated  EventType = "gift_card.created"
	GiftCardRedeemed EventType = "gift_card.redeemed"

	// Subscription events
	SubscriptionCreated   EventType = "subscription.created"
	SubscriptionCancelled EventType = "subscription.cancelled"
	SubscriptionBilled    EventType = "subscription.billed"

	// Inventory events
	StockUpdated  EventType = "stock.updated"
	StockLowAlert EventType = "stock.low_alert"

	// Review events
	ReviewCreated  EventType = "review.created"
	ReviewApproved EventType = "review.approved"
	ReviewRejected EventType = "review.rejected"
	// Return events
	ReturnRequested EventType = "return.requested"
	ReturnApproved  EventType = "return.approved"
	ReturnCompleted EventType = "return.completed"

	// Loyalty events
	LoyaltyPointsEarned  EventType = "loyalty.points_earned"
	LoyaltyPointsRedeemed EventType = "loyalty.points_redeemed"
	LoyaltyRewardCreated EventType = "loyalty.reward_created"
	// Bundle events
	BundleCreated EventType = "bundle.created"
	BundleUpdated EventType = "bundle.updated"

	// Multi-storefront events
	StorefrontCreated EventType = "storefront_entry.created"
	StorefrontUpdated EventType = "storefront_entry.updated"
	StorefrontDeleted EventType = "storefront_entry.deleted"
	// Bulk operation events
	BulkOperationStarted   EventType = "bulk_operation.started"
	BulkOperationCompleted EventType = "bulk_operation.completed"
)

// Event is a domain event that has occurred in the system.
type Event struct {
	ID        string          `json:"id"`
	Type      EventType       `json:"type"`
	TenantID  kernel.TenantID `json:"tenant_id"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp time.Time       `json:"timestamp"`
}

// Handler processes a domain event. Handlers should be idempotent.
// Returning an error logs the failure but does not prevent other handlers from running.
type Handler func(ctx context.Context, event Event) error

// Bus defines the event bus interface for publishing and subscribing to domain events.
type Bus interface {
	// Publish dispatches an event to all registered handlers for the event type.
	// Handlers are called synchronously in the order they were registered.
	Publish(ctx context.Context, event Event) error

	// Subscribe registers a handler for a specific event type.
	// Multiple handlers can be registered for the same event type.
	Subscribe(eventType EventType, handler Handler)

	// SubscribeAll registers a handler for ALL event types (useful for logging/webhooks).
	SubscribeAll(handler Handler)
}

// NewEvent is a helper to construct an Event with a generated ID and timestamp.
func NewEvent(eventType EventType, tenantID kernel.TenantID, payload any) (Event, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return Event{}, err
	}
	return Event{
		ID:        uuid.NewString(),
		Type:      eventType,
		TenantID:  tenantID,
		Payload:   data,
		Timestamp: time.Now().UTC(),
	}, nil
}
