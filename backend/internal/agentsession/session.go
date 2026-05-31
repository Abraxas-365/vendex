package agentsession

import (
	"encoding/json"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// SessionStatus represents the lifecycle of an agent session.
type SessionStatus string

const (
	SessionStatusCreating SessionStatus = "creating"
	SessionStatusRunning  SessionStatus = "running"
	SessionStatusStopped  SessionStatus = "stopped"
	SessionStatusFailed   SessionStatus = "failed"
)

// Session represents a running agent workspace linking a tenant to a preset.
type Session struct {
	ID          kernel.AgentSessionID `json:"id"`
	TenantID    kernel.TenantID       `json:"tenant_id"`
	PresetID    kernel.PresetID       `json:"preset_id"`
	ContainerID string                `json:"container_id"` // Docker container ID
	Status      SessionStatus         `json:"status"`
	FrontendURL string                `json:"frontend_url"` // proxied URL to container frontend
	Metadata    json.RawMessage       `json:"metadata"`     // arbitrary session metadata
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	StoppedAt   *time.Time            `json:"stopped_at"`
}

// ChatMessage holds a single message in a session's chat history.
type ChatMessage struct {
	ID        string                `json:"id"`
	SessionID kernel.AgentSessionID `json:"session_id"`
	Role      string                `json:"role"` // "user", "assistant", "tool"
	Content   string                `json:"content"`
	ToolName  string                `json:"tool_name,omitempty"`
	CreatedAt time.Time             `json:"created_at"`
}

// CreateSessionRequest holds input for creating a new agent session.
type CreateSessionRequest struct {
	PresetID kernel.PresetID `json:"preset_id"`
	Config   json.RawMessage `json:"config"` // optional overrides
}
