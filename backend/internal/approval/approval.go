// Package approval defines the approval request domain entity and DTOs.
package approval

import (
	"encoding/json"
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Status values for an approval request.
const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"
)

// ApprovalRequest represents a queued agent tool action awaiting human review.
type ApprovalRequest struct {
	ID          kernel.ApprovalRequestID `json:"id"`
	TenantID    kernel.TenantID          `json:"tenant_id"`
	SessionID   string                   `json:"session_id"`
	ToolName    string                   `json:"tool_name"`
	ToolInput   json.RawMessage          `json:"tool_input"`
	Status      string                   `json:"status"`
	Reason      string                   `json:"reason"`
	RequestedBy string                   `json:"requested_by"`
	ReviewedBy  string                   `json:"reviewed_by"`
	CreatedAt   time.Time                `json:"created_at"`
	ReviewedAt  *time.Time               `json:"reviewed_at,omitempty"`
}

// CreateApprovalRequest holds the input for creating a new approval request.
type CreateApprovalRequest struct {
	SessionID   string          `json:"session_id"`
	ToolName    string          `json:"tool_name"`
	ToolInput   json.RawMessage `json:"tool_input"`
	RequestedBy string          `json:"requested_by"`
}

// ReviewApprovalRequest holds the input for approving or rejecting a request.
type ReviewApprovalRequest struct {
	Reason string `json:"reason"`
}
