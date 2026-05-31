package bulkops

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// OperationType identifies the kind of bulk operation.
type OperationType string

const (
	OpPriceUpdate  OperationType = "price_update"
	OpStatusChange OperationType = "status_change"
	OpTagAdd       OperationType = "tag_add"
	OpTagRemove    OperationType = "tag_remove"
	OpDelete       OperationType = "delete"
)

// OperationStatus is the lifecycle status of a bulk operation.
type OperationStatus string

const (
	StatusPending    OperationStatus = "pending"
	StatusProcessing OperationStatus = "processing"
	StatusCompleted  OperationStatus = "completed"
	StatusFailed     OperationStatus = "failed"
	StatusCancelled  OperationStatus = "cancelled"
)

// ItemStatus is the per-item processing status.
type ItemStatus string

const (
	ItemPending ItemStatus = "pending"
	ItemSuccess ItemStatus = "success"
	ItemFailed  ItemStatus = "failed"
)

// OperationError holds a per-resource error recorded during processing.
type OperationError struct {
	ResourceID string `json:"resource_id"`
	Message    string `json:"message"`
}

// BulkOperation is the aggregate root for a bulk operation.
type BulkOperation struct {
	ID             kernel.BulkOperationID `json:"id"`
	TenantID       kernel.TenantID        `json:"tenant_id"`
	Type           OperationType          `json:"type"`
	ResourceType   string                 `json:"resource_type"` // "product" or "order"
	Status         OperationStatus        `json:"status"`
	TotalItems     int                    `json:"total_items"`
	ProcessedItems int                    `json:"processed_items"`
	FailedItems    int                    `json:"failed_items"`
	Parameters     map[string]interface{} `json:"parameters"`
	Errors         []OperationError       `json:"errors"`
	CreatedBy      string                 `json:"created_by"`
	StartedAt      *time.Time             `json:"started_at,omitempty"`
	CompletedAt    *time.Time             `json:"completed_at,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
}

// BulkOperationItem tracks the per-resource result of a bulk operation.
type BulkOperationItem struct {
	ID           kernel.BulkOperationItemID `json:"id"`
	TenantID     kernel.TenantID            `json:"tenant_id"`
	OperationID  kernel.BulkOperationID     `json:"operation_id"`
	ResourceID   string                     `json:"resource_id"`
	Status       ItemStatus                 `json:"status"`
	ErrorMessage string                     `json:"error_message,omitempty"`
	ProcessedAt  *time.Time                 `json:"processed_at,omitempty"`
	CreatedAt    time.Time                  `json:"created_at"`
}

// CreateInput holds the data required to create a new bulk operation.
type CreateInput struct {
	Type         OperationType
	ResourceType string
	ResourceIDs  []string
	Parameters   map[string]interface{}
	CreatedBy    string
}
