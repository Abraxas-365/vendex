// Package agenttrigger defines the domain entity for event-triggered agent actions.
package agenttrigger

import (
	"encoding/json"
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Trigger represents a rule that fires an agent action when a store event occurs.
type Trigger struct {
	ID          kernel.AgentTriggerID `json:"id"`
	TenantID    kernel.TenantID       `json:"tenant_id"`
	Name        string                `json:"name"`
	EventType   string                `json:"event_type"`
	Prompt      string                `json:"prompt"`      // supports {{.EventPayload}} template
	PresetID    string                `json:"preset_id"`   // empty = default store manager
	Enabled     bool                  `json:"enabled"`
	Cooldown    int                   `json:"cooldown"`    // minimum seconds between triggers
	LastFiredAt *time.Time            `json:"last_fired_at,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// TriggerLog records each execution of a trigger.
type TriggerLog struct {
	ID            kernel.TriggerLogID    `json:"id"`
	TriggerID     kernel.AgentTriggerID  `json:"trigger_id"`
	TenantID      kernel.TenantID        `json:"tenant_id"`
	EventType     string                 `json:"event_type"`
	EventPayload  json.RawMessage        `json:"event_payload"`
	AgentResponse string                 `json:"agent_response"`
	Status        string                 `json:"status"` // "success", "error", "skipped_cooldown"
	CreatedAt     time.Time              `json:"created_at"`
}

// CreateTriggerRequest holds input for creating a new trigger.
type CreateTriggerRequest struct {
	Name      string `json:"name"`
	EventType string `json:"event_type"`
	Prompt    string `json:"prompt"`
	PresetID  string `json:"preset_id,omitempty"`
	Cooldown  int    `json:"cooldown"` // seconds, default 300
}

// UpdateTriggerRequest holds fields for updating a trigger (all optional).
type UpdateTriggerRequest struct {
	Name     *string `json:"name,omitempty"`
	Prompt   *string `json:"prompt,omitempty"`
	PresetID *string `json:"preset_id,omitempty"`
	Enabled  *bool   `json:"enabled,omitempty"`
	Cooldown *int    `json:"cooldown,omitempty"`
}

// TriggerStatus values for TriggerLog.Status.
const (
	StatusSuccess         = "success"
	StatusError           = "error"
	StatusSkippedCooldown = "skipped_cooldown"
)
