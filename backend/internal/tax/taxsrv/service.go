package taxsrv

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/tax"
	"github.com/google/uuid"
)

// Service handles tax domain business logic.
type Service struct {
	repo tax.Repository
	bus  eventbus.Bus
}

// New creates a new tax service.
func New(repo tax.Repository, bus eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// CreateRateInput holds the data needed to create a tax rate.
type CreateRateInput struct {
	Name             string
	Rate             float64
	Country          string
	State            string
	City             string
	ZipCode          string
	Priority         int
	Compound         bool
	IncludesShipping bool
	Active           bool
}

// UpdateRateInput holds the data needed to update a tax rate.
type UpdateRateInput struct {
	Name             string
	Rate             float64
	Country          string
	State            string
	City             string
	ZipCode          string
	Priority         int
	Compound         bool
	IncludesShipping bool
	Active           bool
}

// CreateRate creates a new tax rate for the given tenant.
func (s *Service) CreateRate(ctx context.Context, tenantID kernel.TenantID, in CreateRateInput) (*tax.TaxRate, error) {
	if in.Rate < 0 || in.Rate > 1 {
		return nil, tax.ErrInvalidRate
	}

	now := time.Now()
	rate := &tax.TaxRate{
		ID:               kernel.TaxRateID(uuid.NewString()),
		TenantID:         tenantID,
		Name:             in.Name,
		Rate:             in.Rate,
		Country:          in.Country,
		State:            in.State,
		City:             in.City,
		ZipCode:          in.ZipCode,
		Priority:         in.Priority,
		Compound:         in.Compound,
		IncludesShipping: in.IncludesShipping,
		Active:           in.Active,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.repo.Create(ctx, rate); err != nil {
		return nil, errx.Wrap(err, "creating tax rate", errx.TypeInternal)
	}

	if evt, err := eventbus.NewEvent(eventbus.TaxRateCreated, tenantID, eventbus.TaxRatePayload{
		RateID:  string(rate.ID),
		Name:    rate.Name,
		Rate:    rate.Rate,
		Country: rate.Country,
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return rate, nil
}

// GetRate retrieves a tax rate by ID, scoped to tenant.
func (s *Service) GetRate(ctx context.Context, tenantID kernel.TenantID, id kernel.TaxRateID) (*tax.TaxRate, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// ListRates returns all tax rates for a tenant.
func (s *Service) ListRates(ctx context.Context, tenantID kernel.TenantID) ([]tax.TaxRate, error) {
	return s.repo.List(ctx, tenantID)
}

// UpdateRate updates an existing tax rate.
func (s *Service) UpdateRate(ctx context.Context, tenantID kernel.TenantID, id kernel.TaxRateID, in UpdateRateInput) (*tax.TaxRate, error) {
	if in.Rate < 0 || in.Rate > 1 {
		return nil, tax.ErrInvalidRate
	}

	existing, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	existing.Name = in.Name
	existing.Rate = in.Rate
	existing.Country = in.Country
	existing.State = in.State
	existing.City = in.City
	existing.ZipCode = in.ZipCode
	existing.Priority = in.Priority
	existing.Compound = in.Compound
	existing.IncludesShipping = in.IncludesShipping
	existing.Active = in.Active
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, errx.Wrap(err, "updating tax rate", errx.TypeInternal)
	}

	return existing, nil
}

// DeleteRate removes a tax rate by ID, scoped to tenant.
func (s *Service) DeleteRate(ctx context.Context, tenantID kernel.TenantID, id kernel.TaxRateID) error {
	return s.repo.Delete(ctx, tenantID, id)
}

// CalculateTax calculates tax for a given subtotal and shipping amount at a location.
// It finds active rates matching the location, sorts by priority, applies non-compound
// rates first (to subtotal + optionally shipping), then compound rates on top.
func (s *Service) CalculateTax(
	ctx context.Context,
	tenantID kernel.TenantID,
	subtotalCents int64,
	shippingCents int64,
	country, state, city, zipCode string,
) (*tax.TaxResult, error) {
	rates, err := s.repo.FindByLocation(ctx, tenantID, country, state, city, zipCode)
	if err != nil {
		return nil, err
	}

	// Sort by priority (ascending)
	sort.Slice(rates, func(i, j int) bool {
		return rates[i].Priority < rates[j].Priority
	})

	result := &tax.TaxResult{
		TaxBreakdown: []tax.TaxLineItem{},
	}

	var nonCompoundTotal int64 // accumulated tax from non-compound rates

	// First pass: apply non-compound rates
	for _, r := range rates {
		if r.Compound {
			continue
		}
		taxable := subtotalCents
		if r.IncludesShipping {
			taxable += shippingCents
		}
		amount := roundToCents(float64(taxable) * r.Rate)
		nonCompoundTotal += amount

		result.TaxBreakdown = append(result.TaxBreakdown, tax.TaxLineItem{
			RateID:   r.ID,
			Name:     r.Name,
			Rate:     r.Rate,
			Amount:   amount,
			Compound: false,
		})
	}

	// Second pass: apply compound rates (on top of subtotal + non-compound taxes)
	for _, r := range rates {
		if !r.Compound {
			continue
		}
		taxable := subtotalCents + nonCompoundTotal
		if r.IncludesShipping {
			taxable += shippingCents
		}
		amount := roundToCents(float64(taxable) * r.Rate)

		result.TaxBreakdown = append(result.TaxBreakdown, tax.TaxLineItem{
			RateID:   r.ID,
			Name:     r.Name,
			Rate:     r.Rate,
			Amount:   amount,
			Compound: true,
		})
		result.TotalTax += amount
	}

	// Add non-compound total to grand total
	result.TotalTax += nonCompoundTotal

	return result, nil
}

// roundToCents rounds a float amount to the nearest cent (int64).
func roundToCents(amount float64) int64 {
	return int64(math.Round(amount))
}
