package notification

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Type represents the severity/kind of a notification.
type Type string

const (
	TypeInfo    Type = "info"
	TypeWarning Type = "warning"
	TypeSuccess Type = "success"
	TypeError   Type = "error"
)

// Notification is the aggregate root for the notifications domain.
type Notification struct {
	ID           kernel.NotificationID `json:"id"`
	TenantID     kernel.TenantID       `json:"tenant_id"`
	UserID       kernel.UserID         `json:"user_id"`
	Title        string                `json:"title"`
	Body         string                `json:"body,omitempty"`
	Type         Type                  `json:"type"`
	ResourceType string                `json:"resource_type,omitempty"`
	ResourceID   string                `json:"resource_id,omitempty"`
	Read         bool                  `json:"read"`
	ReadAt       *time.Time            `json:"read_at,omitempty"`
	CreatedAt    time.Time             `json:"created_at"`
}

// CreateNotificationInput holds the data required to create a new notification.
type CreateNotificationInput struct {
	UserID       kernel.UserID
	Title        string
	Body         string
	Type         Type
	ResourceType string
	ResourceID   string
}
