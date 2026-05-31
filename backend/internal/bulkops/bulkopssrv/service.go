package bulkopssrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/bulkops"
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/google/uuid"
)

// Service implements bulk operations business logic.
type Service struct {
	repo bulkops.Repository
	bus  eventbus.Bus
}

// New creates a new bulk operations Service.
func New(repo bulkops.Repository, bus eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// Create creates a new bulk operation together with its item records.
func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, input bulkops.CreateInput) (*bulkops.BulkOperation, error) {
	if input.Type == "" {
		return nil, bulkops.ErrInvalidInput
	}
	if input.ResourceType == "" {
		return nil, bulkops.ErrInvalidInput
	}
	if len(input.ResourceIDs) == 0 {
		return nil, bulkops.ErrNoResourceIDs
	}

	params := input.Parameters
	if params == nil {
		params = map[string]interface{}{}
	}

	now := time.Now().UTC()
	op := &bulkops.BulkOperation{
		ID:             kernel.BulkOperationID(uuid.NewString()),
		TenantID:       tenantID,
		Type:           input.Type,
		ResourceType:   input.ResourceType,
		Status:         bulkops.StatusPending,
		TotalItems:     len(input.ResourceIDs),
		ProcessedItems: 0,
		FailedItems:    0,
		Parameters:     params,
		Errors:         []bulkops.OperationError{},
		CreatedBy:      input.CreatedBy,
		CreatedAt:      now,
	}

	items := make([]bulkops.BulkOperationItem, 0, len(input.ResourceIDs))
	for _, rid := range input.ResourceIDs {
		items = append(items, bulkops.BulkOperationItem{
			ID:          kernel.BulkOperationItemID(uuid.NewString()),
			TenantID:    tenantID,
			OperationID: op.ID,
			ResourceID:  rid,
			Status:      bulkops.ItemPending,
			CreatedAt:   now,
		})
	}

	if err := s.repo.Create(ctx, op, items); err != nil {
		return nil, err
	}

	return op, nil
}

// GetByID returns a bulk operation scoped to the tenant.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.BulkOperationID) (*bulkops.BulkOperation, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// List returns paginated bulk operations for a tenant.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, page, pageSize int) (kernel.Paginated[bulkops.BulkOperation], error) {
	return s.repo.List(ctx, tenantID, page, pageSize)
}

// ListItems returns paginated items for a specific bulk operation.
func (s *Service) ListItems(ctx context.Context, tenantID kernel.TenantID, operationID kernel.BulkOperationID, page, pageSize int) (kernel.Paginated[bulkops.BulkOperationItem], error) {
	// Verify operation exists and belongs to tenant.
	if _, err := s.repo.GetByID(ctx, tenantID, operationID); err != nil {
		return kernel.Paginated[bulkops.BulkOperationItem]{}, err
	}
	return s.repo.ListItems(ctx, tenantID, operationID, page, pageSize)
}

// Process marks an operation as processing and synchronously applies it to each item.
// In a production system this would be dispatched to a job queue; here we run inline.
func (s *Service) Process(ctx context.Context, tenantID kernel.TenantID, id kernel.BulkOperationID) (*bulkops.BulkOperation, error) {
	op, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	switch op.Status {
	case bulkops.StatusCompleted:
		return nil, bulkops.ErrAlreadyCompleted
	case bulkops.StatusProcessing:
		return nil, bulkops.ErrAlreadyProcessing
	case bulkops.StatusCancelled, bulkops.StatusFailed:
		return nil, bulkops.ErrCannotCancel
	}

	// Mark as processing.
	now := time.Now().UTC()
	op.Status = bulkops.StatusProcessing
	op.StartedAt = &now
	if err := s.repo.UpdateOperation(ctx, op); err != nil {
		return nil, err
	}

	// Publish started event.
	startedEvt, evtErr := eventbus.NewEvent(eventbus.BulkOperationStarted, tenantID, BulkOperationPayload{
		OperationID:  string(op.ID),
		OperationType: string(op.Type),
		ResourceType: op.ResourceType,
		TotalItems:   op.TotalItems,
	})
	if evtErr == nil {
		_ = s.bus.Publish(ctx, startedEvt)
	}

	// Fetch all pending items.
	itemsResult, err := s.repo.ListItems(ctx, tenantID, id, 1, 1000)
	if err != nil {
		return nil, err
	}

	processedCount := 0
	failedCount := 0
	opErrors := make([]bulkops.OperationError, 0)

	for i := range itemsResult.Items {
		item := &itemsResult.Items[i]
		if item.Status != bulkops.ItemPending {
			continue
		}

		// Apply the operation to the item.
		itemErr := applyOperation(op.Type, op.Parameters, item.ResourceID)
		itemNow := time.Now().UTC()
		item.ProcessedAt = &itemNow

		if itemErr != nil {
			item.Status = bulkops.ItemFailed
			item.ErrorMessage = itemErr.Error()
			failedCount++
			opErrors = append(opErrors, bulkops.OperationError{
				ResourceID: item.ResourceID,
				Message:    itemErr.Error(),
			})
		} else {
			item.Status = bulkops.ItemSuccess
			processedCount++
		}

		_ = s.repo.UpdateItem(ctx, item)
	}

	// Finalise the operation.
	finishedAt := time.Now().UTC()
	op.ProcessedItems = processedCount
	op.FailedItems = failedCount
	op.Errors = opErrors
	op.CompletedAt = &finishedAt

	if failedCount == op.TotalItems {
		op.Status = bulkops.StatusFailed
	} else {
		op.Status = bulkops.StatusCompleted
	}

	if err := s.repo.UpdateOperation(ctx, op); err != nil {
		return nil, err
	}

	// Publish completed event.
	completedEvt, evtErr := eventbus.NewEvent(eventbus.BulkOperationCompleted, tenantID, BulkOperationPayload{
		OperationID:   string(op.ID),
		OperationType: string(op.Type),
		ResourceType:  op.ResourceType,
		TotalItems:    op.TotalItems,
		FailedItems:   failedCount,
	})
	if evtErr == nil {
		_ = s.bus.Publish(ctx, completedEvt)
	}

	return op, nil
}

// Cancel marks a pending operation as cancelled.
func (s *Service) Cancel(ctx context.Context, tenantID kernel.TenantID, id kernel.BulkOperationID) error {
	op, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	switch op.Status {
	case bulkops.StatusCompleted:
		return bulkops.ErrAlreadyCompleted
	case bulkops.StatusProcessing:
		return bulkops.ErrAlreadyProcessing
	case bulkops.StatusCancelled:
		return nil // idempotent
	}

	return s.repo.UpdateStatus(ctx, tenantID, id, bulkops.StatusCancelled)
}

// ---------------------------------------------------------------------------
// Operation application
// ---------------------------------------------------------------------------

// applyOperation performs the actual work for a single resource item.
// This is intentionally a stub — real implementations would call the relevant
// domain services (product, order, etc.) via injected interfaces.
func applyOperation(opType bulkops.OperationType, params map[string]interface{}, resourceID string) error {
	// Stub: in a real implementation, inject domain service interfaces and
	// delegate based on opType (price_update → product.UpdatePrice, etc.).
	// For now we record success for all items to keep the domain buildable.
	_ = opType
	_ = params
	_ = resourceID
	return nil
}

// BulkOperationPayload is the event payload for bulk operation events.
type BulkOperationPayload struct {
	OperationID   string `json:"operation_id"`
	OperationType string `json:"operation_type"`
	ResourceType  string `json:"resource_type"`
	TotalItems    int    `json:"total_items"`
	FailedItems   int    `json:"failed_items,omitempty"`
}
