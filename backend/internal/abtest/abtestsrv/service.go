package abtestsrv

import (
	"context"
	"math/rand"
	"time"

	"github.com/Abraxas-365/vendex/internal/abtest"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/google/uuid"
)

// Service handles A/B testing business logic.
type Service struct {
	repo abtest.Repository
	bus  eventbus.Bus
}

// New creates a new A/B testing service.
func New(repo abtest.Repository, bus eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// CreateExperiment creates a new experiment in draft status.
func (s *Service) CreateExperiment(ctx context.Context, tenantID kernel.TenantID, in abtest.CreateExperimentInput) (*abtest.Experiment, error) {
	if in.TrafficPercent < 1 || in.TrafficPercent > 100 {
		if in.TrafficPercent == 0 {
			in.TrafficPercent = 100
		} else {
			return nil, abtest.ErrInvalidTrafficPercent
		}
	}

	expType := in.Type
	if expType == "" {
		expType = abtest.TypePage
	}

	now := time.Now()
	e := &abtest.Experiment{
		ID:             kernel.ExperimentID(uuid.NewString()),
		TenantID:       tenantID,
		Name:           in.Name,
		Description:    in.Description,
		Type:           expType,
		Status:         abtest.StatusDraft,
		TrafficPercent: in.TrafficPercent,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repo.CreateExperiment(ctx, e); err != nil {
		return nil, err
	}

	return e, nil
}

// GetExperiment retrieves a single experiment with its variants.
func (s *Service) GetExperiment(ctx context.Context, tenantID kernel.TenantID, id kernel.ExperimentID) (*abtest.Experiment, error) {
	e, err := s.repo.GetExperimentByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	variants, err := s.repo.ListVariants(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	e.Variants = variants

	return e, nil
}

// ListExperiments returns paginated experiments, optionally filtered by status.
func (s *Service) ListExperiments(ctx context.Context, tenantID kernel.TenantID, status string, page, pageSize int) (kernel.Paginated[abtest.Experiment], error) {
	pg := kernel.NewPaginationOptions(page, pageSize)
	return s.repo.ListExperiments(ctx, tenantID, status, pg)
}

// UpdateExperiment updates mutable fields of a draft experiment.
func (s *Service) UpdateExperiment(ctx context.Context, tenantID kernel.TenantID, id kernel.ExperimentID, in abtest.UpdateExperimentInput) (*abtest.Experiment, error) {
	e, err := s.repo.GetExperimentByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if e.Status == abtest.StatusRunning {
		return nil, abtest.ErrCannotModifyRunning
	}

	if in.Name != nil {
		e.Name = *in.Name
	}
	if in.Description != nil {
		e.Description = *in.Description
	}
	if in.TrafficPercent != nil {
		if *in.TrafficPercent < 1 || *in.TrafficPercent > 100 {
			return nil, abtest.ErrInvalidTrafficPercent
		}
		e.TrafficPercent = *in.TrafficPercent
	}
	e.UpdatedAt = time.Now()

	if err := s.repo.UpdateExperiment(ctx, e); err != nil {
		return nil, err
	}

	return e, nil
}

// DeleteExperiment removes an experiment and its variants.
func (s *Service) DeleteExperiment(ctx context.Context, tenantID kernel.TenantID, id kernel.ExperimentID) error {
	_, err := s.repo.GetExperimentByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	return s.repo.DeleteExperiment(ctx, tenantID, id)
}

// StartExperiment transitions a draft/paused experiment to running.
func (s *Service) StartExperiment(ctx context.Context, tenantID kernel.TenantID, id kernel.ExperimentID) (*abtest.Experiment, error) {
	e, err := s.repo.GetExperimentByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if e.Status == abtest.StatusRunning {
		return nil, abtest.ErrAlreadyRunning
	}
	if e.Status == abtest.StatusCompleted {
		return nil, abtest.ErrAlreadyCompleted
	}

	variants, err := s.repo.ListVariants(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if len(variants) < 2 {
		return nil, abtest.ErrInsufficientVariants
	}

	now := time.Now()
	e.Status = abtest.StatusRunning
	e.StartedAt = &now
	e.UpdatedAt = now
	e.Variants = variants

	if err := s.repo.UpdateExperiment(ctx, e); err != nil {
		return nil, err
	}

	_ = s.publishEvent(ctx, eventbus.ExperimentStarted, tenantID, map[string]string{
		"experiment_id": string(e.ID),
	})

	return e, nil
}

// PauseExperiment pauses a running experiment.
func (s *Service) PauseExperiment(ctx context.Context, tenantID kernel.TenantID, id kernel.ExperimentID) (*abtest.Experiment, error) {
	e, err := s.repo.GetExperimentByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if e.Status != abtest.StatusRunning {
		return nil, abtest.ErrNotRunning
	}

	e.Status = abtest.StatusPaused
	e.UpdatedAt = time.Now()

	if err := s.repo.UpdateExperiment(ctx, e); err != nil {
		return nil, err
	}

	variants, err := s.repo.ListVariants(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	e.Variants = variants

	return e, nil
}

// CompleteExperiment marks an experiment as completed with a declared winner.
func (s *Service) CompleteExperiment(ctx context.Context, tenantID kernel.TenantID, id kernel.ExperimentID, winnerVariantID kernel.ExperimentVariantID) (*abtest.Experiment, error) {
	e, err := s.repo.GetExperimentByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if e.Status == abtest.StatusCompleted {
		return nil, abtest.ErrAlreadyCompleted
	}

	// Validate winner variant belongs to this experiment.
	_, err = s.repo.GetVariantByID(ctx, tenantID, winnerVariantID)
	if err != nil {
		return nil, abtest.ErrVariantNotFound
	}

	now := time.Now()
	e.Status = abtest.StatusCompleted
	e.WinnerVariantID = &winnerVariantID
	e.EndedAt = &now
	e.UpdatedAt = now

	if err := s.repo.UpdateExperiment(ctx, e); err != nil {
		return nil, err
	}

	variants, err := s.repo.ListVariants(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	e.Variants = variants

	_ = s.publishEvent(ctx, eventbus.ExperimentCompleted, tenantID, map[string]string{
		"experiment_id":     string(e.ID),
		"winner_variant_id": string(winnerVariantID),
	})

	return e, nil
}

// AddVariant adds a new variant to an experiment.
func (s *Service) AddVariant(ctx context.Context, tenantID kernel.TenantID, experimentID kernel.ExperimentID, in abtest.CreateVariantInput) (*abtest.ExperimentVariant, error) {
	e, err := s.repo.GetExperimentByID(ctx, tenantID, experimentID)
	if err != nil {
		return nil, err
	}

	if e.Status == abtest.StatusRunning {
		return nil, abtest.ErrCannotModifyRunning
	}

	weight := in.Weight
	if weight <= 0 {
		weight = 50
	}

	cfg := in.Config
	if cfg == nil {
		cfg = map[string]interface{}{}
	}

	v := &abtest.ExperimentVariant{
		ID:           kernel.ExperimentVariantID(uuid.NewString()),
		TenantID:     tenantID,
		ExperimentID: experimentID,
		Name:         in.Name,
		Description:  in.Description,
		Weight:       weight,
		IsControl:    in.IsControl,
		Config:       cfg,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.CreateVariant(ctx, v); err != nil {
		return nil, err
	}

	return v, nil
}

// RemoveVariant removes a variant from an experiment.
func (s *Service) RemoveVariant(ctx context.Context, tenantID kernel.TenantID, experimentID kernel.ExperimentID, variantID kernel.ExperimentVariantID) error {
	e, err := s.repo.GetExperimentByID(ctx, tenantID, experimentID)
	if err != nil {
		return err
	}

	if e.Status == abtest.StatusRunning {
		return abtest.ErrCannotModifyRunning
	}

	_, err = s.repo.GetVariantByID(ctx, tenantID, variantID)
	if err != nil {
		return err
	}

	return s.repo.DeleteVariant(ctx, tenantID, variantID)
}

// AssignVisitor assigns a visitor to a variant using weighted random selection.
// If the visitor is already assigned, returns the existing assignment.
func (s *Service) AssignVisitor(ctx context.Context, tenantID kernel.TenantID, experimentID kernel.ExperimentID, visitorID string) (*abtest.ExperimentAssignment, error) {
	// Check if already assigned.
	existing, err := s.repo.GetAssignment(ctx, tenantID, experimentID, visitorID)
	if err == nil && existing != nil {
		return existing, nil
	}
	if err != nil && !errx.Is(err, abtest.ErrAssignmentNotFound) {
		return nil, err
	}

	e, err := s.repo.GetExperimentByID(ctx, tenantID, experimentID)
	if err != nil {
		return nil, err
	}

	if e.Status != abtest.StatusRunning {
		return nil, abtest.ErrNotRunning
	}

	variants, err := s.repo.ListVariants(ctx, tenantID, experimentID)
	if err != nil {
		return nil, err
	}
	if len(variants) == 0 {
		return nil, abtest.ErrInsufficientVariants
	}

	// Weighted random selection.
	selected := weightedRandomVariant(variants)

	now := time.Now()
	a := &abtest.ExperimentAssignment{
		ID:           kernel.ExperimentAssignmentID(uuid.NewString()),
		TenantID:     tenantID,
		ExperimentID: experimentID,
		VariantID:    selected.ID,
		VisitorID:    visitorID,
		Converted:    false,
		RevenueCents: 0,
		AssignedAt:   now,
	}

	if err := s.repo.CreateAssignment(ctx, a); err != nil {
		return nil, err
	}

	// Increment variant visitor count.
	_ = s.repo.IncrementVariantVisitors(ctx, selected.ID)

	return a, nil
}

// RecordConversion marks a visitor's assignment as converted.
func (s *Service) RecordConversion(ctx context.Context, tenantID kernel.TenantID, experimentID kernel.ExperimentID, visitorID string, revenueCents int64) error {
	assignment, err := s.repo.GetAssignment(ctx, tenantID, experimentID, visitorID)
	if err != nil {
		return err
	}

	if assignment.Converted {
		// Idempotent — already converted.
		return nil
	}

	if err := s.repo.RecordConversion(ctx, tenantID, experimentID, visitorID, revenueCents); err != nil {
		return err
	}

	_ = s.repo.IncrementVariantConversions(ctx, assignment.VariantID, revenueCents)

	return nil
}

// GetResults returns calculated statistics for an experiment.
func (s *Service) GetResults(ctx context.Context, tenantID kernel.TenantID, experimentID kernel.ExperimentID) (*abtest.ExperimentResults, error) {
	e, err := s.repo.GetExperimentByID(ctx, tenantID, experimentID)
	if err != nil {
		return nil, err
	}

	variants, err := s.repo.ListVariants(ctx, tenantID, experimentID)
	if err != nil {
		return nil, err
	}

	results := &abtest.ExperimentResults{
		ExperimentID: experimentID,
		Variants:     make([]abtest.VariantResult, 0, len(variants)),
	}

	for _, v := range variants {
		var convRate float64
		if v.Visitors > 0 {
			convRate = float64(v.Conversions) / float64(v.Visitors)
		}

		isWinner := e.WinnerVariantID != nil && *e.WinnerVariantID == v.ID

		results.Variants = append(results.Variants, abtest.VariantResult{
			VariantID:      v.ID,
			Name:           v.Name,
			IsControl:      v.IsControl,
			Visitors:       v.Visitors,
			Conversions:    v.Conversions,
			ConversionRate: convRate,
			Revenue:        v.RevenueCents,
			IsWinner:       isWinner,
		})
	}

	return results, nil
}

// weightedRandomVariant selects a variant using weighted random selection.
func weightedRandomVariant(variants []abtest.ExperimentVariant) abtest.ExperimentVariant {
	totalWeight := 0
	for _, v := range variants {
		totalWeight += v.Weight
	}
	if totalWeight <= 0 {
		return variants[0]
	}

	r := rand.Intn(totalWeight)
	cumulative := 0
	for _, v := range variants {
		cumulative += v.Weight
		if r < cumulative {
			return v
		}
	}
	return variants[len(variants)-1]
}

// publishEvent is a fire-and-forget event publisher.
func (s *Service) publishEvent(ctx context.Context, eventType eventbus.EventType, tenantID kernel.TenantID, payload any) error {
	event, err := eventbus.NewEvent(eventType, tenantID, payload)
	if err != nil {
		return err
	}
	return s.bus.Publish(ctx, event)
}
