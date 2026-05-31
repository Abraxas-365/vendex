package abtest

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// ExperimentStatus represents the lifecycle state of an experiment.
type ExperimentStatus string

const (
	StatusDraft     ExperimentStatus = "draft"
	StatusRunning   ExperimentStatus = "running"
	StatusPaused    ExperimentStatus = "paused"
	StatusCompleted ExperimentStatus = "completed"
)

// ExperimentType represents what is being tested.
type ExperimentType string

const (
	TypePage   ExperimentType = "page"
	TypePrice  ExperimentType = "price"
	TypeLayout ExperimentType = "layout"
	TypeCopy   ExperimentType = "copy"
)

// Experiment is the core entity for an A/B test.
type Experiment struct {
	ID              kernel.ExperimentID         `json:"id" db:"id"`
	TenantID        kernel.TenantID             `json:"tenant_id" db:"tenant_id"`
	Name            string                      `json:"name" db:"name"`
	Description     string                      `json:"description" db:"description"`
	Type            ExperimentType              `json:"type" db:"type"`
	Status          ExperimentStatus            `json:"status" db:"status"`
	TrafficPercent  int                         `json:"traffic_percent" db:"traffic_percent"`
	StartedAt       *time.Time                  `json:"started_at,omitempty" db:"started_at"`
	EndedAt         *time.Time                  `json:"ended_at,omitempty" db:"ended_at"`
	WinnerVariantID *kernel.ExperimentVariantID `json:"winner_variant_id,omitempty" db:"winner_variant_id"`
	Variants        []ExperimentVariant         `json:"variants,omitempty"`
	CreatedAt       time.Time                   `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time                   `json:"updated_at" db:"updated_at"`
}

// ExperimentVariant represents one branch of an experiment.
type ExperimentVariant struct {
	ID           kernel.ExperimentVariantID `json:"id" db:"id"`
	TenantID     kernel.TenantID            `json:"tenant_id" db:"tenant_id"`
	ExperimentID kernel.ExperimentID        `json:"experiment_id" db:"experiment_id"`
	Name         string                     `json:"name" db:"name"`
	Description  string                     `json:"description" db:"description"`
	Weight       int                        `json:"weight" db:"weight"`
	IsControl    bool                       `json:"is_control" db:"is_control"`
	Config       map[string]interface{}     `json:"config"`
	Visitors     int                        `json:"visitors" db:"visitors"`
	Conversions  int                        `json:"conversions" db:"conversions"`
	RevenueCents int64                      `json:"revenue_cents" db:"revenue_cents"`
	CreatedAt    time.Time                  `json:"created_at" db:"created_at"`
}

// ExperimentAssignment records which variant a visitor was assigned to.
type ExperimentAssignment struct {
	ID           kernel.ExperimentAssignmentID `json:"id" db:"id"`
	TenantID     kernel.TenantID               `json:"tenant_id" db:"tenant_id"`
	ExperimentID kernel.ExperimentID           `json:"experiment_id" db:"experiment_id"`
	VariantID    kernel.ExperimentVariantID    `json:"variant_id" db:"variant_id"`
	VisitorID    string                        `json:"visitor_id" db:"visitor_id"`
	Converted    bool                          `json:"converted" db:"converted"`
	RevenueCents int64                         `json:"revenue_cents" db:"revenue_cents"`
	AssignedAt   time.Time                     `json:"assigned_at" db:"assigned_at"`
	ConvertedAt  *time.Time                    `json:"converted_at,omitempty" db:"converted_at"`
}

// ExperimentResults holds calculated statistics for an experiment.
type ExperimentResults struct {
	ExperimentID kernel.ExperimentID `json:"experiment_id"`
	Variants     []VariantResult     `json:"variants"`
}

// VariantResult holds statistics for a single variant.
type VariantResult struct {
	VariantID      kernel.ExperimentVariantID `json:"variant_id"`
	Name           string                     `json:"name"`
	IsControl      bool                       `json:"is_control"`
	Visitors       int                        `json:"visitors"`
	Conversions    int                        `json:"conversions"`
	ConversionRate float64                    `json:"conversion_rate"`
	Revenue        int64                      `json:"revenue"`
	IsWinner       bool                       `json:"is_winner"`
}

// CreateExperimentInput holds the data needed to create an experiment.
type CreateExperimentInput struct {
	Name           string
	Description    string
	Type           ExperimentType
	TrafficPercent int
}

// UpdateExperimentInput holds updatable fields.
type UpdateExperimentInput struct {
	Name           *string
	Description    *string
	TrafficPercent *int
}

// CreateVariantInput holds the data needed to add a variant.
type CreateVariantInput struct {
	Name        string
	Description string
	Weight      int
	IsControl   bool
	Config      map[string]interface{}
}
