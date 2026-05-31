package abtest

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines persistence operations for the A/B testing domain.
type Repository interface {
	// Experiment CRUD
	CreateExperiment(ctx context.Context, e *Experiment) error
	GetExperimentByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ExperimentID) (*Experiment, error)
	UpdateExperiment(ctx context.Context, e *Experiment) error
	DeleteExperiment(ctx context.Context, tenantID kernel.TenantID, id kernel.ExperimentID) error
	ListExperiments(ctx context.Context, tenantID kernel.TenantID, status string, pg kernel.PaginationOptions) (kernel.Paginated[Experiment], error)

	// Variant operations
	CreateVariant(ctx context.Context, v *ExperimentVariant) error
	GetVariantByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ExperimentVariantID) (*ExperimentVariant, error)
	ListVariants(ctx context.Context, tenantID kernel.TenantID, experimentID kernel.ExperimentID) ([]ExperimentVariant, error)
	DeleteVariant(ctx context.Context, tenantID kernel.TenantID, id kernel.ExperimentVariantID) error
	IncrementVariantVisitors(ctx context.Context, variantID kernel.ExperimentVariantID) error
	IncrementVariantConversions(ctx context.Context, variantID kernel.ExperimentVariantID, revenueCents int64) error

	// Assignment operations
	CreateAssignment(ctx context.Context, a *ExperimentAssignment) error
	GetAssignment(ctx context.Context, tenantID kernel.TenantID, experimentID kernel.ExperimentID, visitorID string) (*ExperimentAssignment, error)
	RecordConversion(ctx context.Context, tenantID kernel.TenantID, experimentID kernel.ExperimentID, visitorID string, revenueCents int64) error
}
