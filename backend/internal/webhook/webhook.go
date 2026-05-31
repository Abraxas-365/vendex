package webhook

import (
	"encoding/json"
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// DeliveryStatus represents the state of a webhook delivery attempt.
type DeliveryStatus string

const (
	DeliveryPending DeliveryStatus = "pending"
	DeliverySuccess DeliveryStatus = "success"
	DeliveryFailed  DeliveryStatus = "failed"
)

// Webhook is the aggregate root for registered HTTP callbacks.
type Webhook struct {
	ID          kernel.WebhookID   `json:"id"`
	TenantID    kernel.TenantID    `json:"tenant_id"`
	URL         string             `json:"url"`
	Secret      string             `json:"secret"`
	Events      []string           `json:"events"`
	Active      bool               `json:"active"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// WebhookDelivery captures a single delivery attempt for a webhook event.
type WebhookDelivery struct {
	ID             kernel.WebhookDeliveryID `json:"id"`
	TenantID       kernel.TenantID          `json:"tenant_id"`
	WebhookID      kernel.WebhookID         `json:"webhook_id"`
	EventType      string                   `json:"event_type"`
	Payload        json.RawMessage          `json:"payload"`
	ResponseStatus *int                     `json:"response_status,omitempty"`
	ResponseBody   string                   `json:"response_body,omitempty"`
	Status         DeliveryStatus           `json:"status"`
	Attempts       int                      `json:"attempts"`
	MaxAttempts    int                      `json:"max_attempts"`
	NextRetryAt    *time.Time               `json:"next_retry_at,omitempty"`
	DeliveredAt    *time.Time               `json:"delivered_at,omitempty"`
	CreatedAt      time.Time                `json:"created_at"`
}

// CreateWebhookInput holds the data required to register a new webhook.
type CreateWebhookInput struct {
	URL         string   `json:"url"`
	Secret      string   `json:"secret"`
	Events      []string `json:"events"`
	Description string   `json:"description"`
}

// UpdateWebhookInput holds the fields that can be updated on a webhook.
type UpdateWebhookInput struct {
	URL         *string  `json:"url,omitempty"`
	Secret      *string  `json:"secret,omitempty"`
	Events      []string `json:"events,omitempty"`
	Description *string  `json:"description,omitempty"`
}
