package audit

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// AuditEntry records a single admin action performed within a tenant.
type AuditEntry struct {
	ID           kernel.AuditEntryID `json:"id"`
	TenantID     kernel.TenantID     `json:"tenant_id"`
	UserID       string              `json:"user_id"`
	UserEmail    string              `json:"user_email,omitempty"`
	Action       string              `json:"action"`
	ResourceType string              `json:"resource_type"`
	ResourceID   string              `json:"resource_id,omitempty"`
	Changes      map[string]any      `json:"changes,omitempty"`
	Metadata     map[string]any      `json:"metadata,omitempty"`
	IPAddress    string              `json:"ip_address,omitempty"`
	UserAgent    string              `json:"user_agent,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
}

// CreateAuditInput carries the data required to log an audit event.
type CreateAuditInput struct {
	TenantID     kernel.TenantID
	UserID       string
	UserEmail    string
	Action       string
	ResourceType string
	ResourceID   string
	Changes      map[string]any
	Metadata     map[string]any
	IPAddress    string
	UserAgent    string
}

// AuditFilter supports filtering the audit log list query.
type AuditFilter struct {
	UserID       string
	Action       string
	ResourceType string
	ResourceID   string
	From         *time.Time
	To           *time.Time
}

// ActionStats holds a count keyed by action name.
type ActionStats struct {
	Action string `json:"action"`
	Count  int    `json:"count"`
}
