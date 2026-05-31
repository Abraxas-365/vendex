package shippingsrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/shipping"
	"github.com/google/uuid"
)

// AvailableRate represents a computed shipping option returned to the caller.
type AvailableRate struct {
	RateID     kernel.ShippingRateID `json:"rate_id"`
	Name       string                `json:"name"`
	Price      kernel.Money          `json:"price"`
	EstDaysMin *int                  `json:"est_days_min,omitempty"`
	EstDaysMax *int                  `json:"est_days_max,omitempty"`
}

// Service handles shipping business logic.
type Service struct {
	zones ZoneRepo
	rates RateRepo
	bus   eventbus.Bus
}

// ZoneRepo is a local alias kept internal to the service package.
type ZoneRepo = shipping.ZoneRepository

// RateRepo is a local alias kept internal to the service package.
type RateRepo = shipping.RateRepository

// New creates a new shipping service.
func New(zones shipping.ZoneRepository, rates shipping.RateRepository, bus eventbus.Bus) *Service {
	return &Service{zones: zones, rates: rates, bus: bus}
}

// ---------------------------------------------------------------------------
// Zone CRUD
// ---------------------------------------------------------------------------

// CreateZone persists a new shipping zone.
func (s *Service) CreateZone(ctx context.Context, tenantID kernel.TenantID, name string, countries, states []string) (*shipping.ShippingZone, error) {
	if name == "" {
		return nil, errx.New("zone name is required", errx.TypeValidation)
	}
	now := time.Now().UTC()
	zone := &shipping.ShippingZone{
		ID:        kernel.ShippingZoneID(uuid.NewString()),
		TenantID:  tenantID,
		Name:      name,
		Countries: countries,
		States:    states,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if zone.Countries == nil {
		zone.Countries = []string{}
	}
	if zone.States == nil {
		zone.States = []string{}
	}

	if err := s.zones.Create(ctx, zone); err != nil {
		return nil, err
	}

	if evt, err := eventbus.NewEvent(eventbus.ShippingZoneCreated, tenantID, eventbus.ShippingZonePayload{
		ZoneID:    string(zone.ID),
		Name:      zone.Name,
		Countries: zone.Countries,
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return zone, nil
}

// GetZone retrieves a shipping zone by ID.
func (s *Service) GetZone(ctx context.Context, tenantID kernel.TenantID, id kernel.ShippingZoneID) (*shipping.ShippingZone, error) {
	return s.zones.GetByID(ctx, tenantID, id)
}

// ListZones returns all shipping zones for a tenant.
func (s *Service) ListZones(ctx context.Context, tenantID kernel.TenantID) ([]shipping.ShippingZone, error) {
	return s.zones.List(ctx, tenantID)
}

// UpdateZone modifies an existing shipping zone.
func (s *Service) UpdateZone(ctx context.Context, tenantID kernel.TenantID, id kernel.ShippingZoneID, name string, countries, states []string) (*shipping.ShippingZone, error) {
	zone, err := s.zones.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		zone.Name = name
	}
	if countries != nil {
		zone.Countries = countries
	}
	if states != nil {
		zone.States = states
	}
	zone.UpdatedAt = time.Now().UTC()

	if err := s.zones.Update(ctx, zone); err != nil {
		return nil, err
	}
	return zone, nil
}

// DeleteZone removes a shipping zone.
func (s *Service) DeleteZone(ctx context.Context, tenantID kernel.TenantID, id kernel.ShippingZoneID) error {
	return s.zones.Delete(ctx, tenantID, id)
}

// ---------------------------------------------------------------------------
// Rate CRUD
// ---------------------------------------------------------------------------

// CreateRateInput holds data for creating a shipping rate.
type CreateRateInput struct {
	ZoneID         kernel.ShippingZoneID
	Name           string
	Type           shipping.RateType
	PriceAmount    int64
	PriceCurrency  string
	MinWeight      *float64
	MaxWeight      *float64
	MinOrderAmount *int64
	MaxOrderAmount *int64
	EstDaysMin     *int
	EstDaysMax     *int
}

// CreateRate persists a new shipping rate within a zone.
func (s *Service) CreateRate(ctx context.Context, tenantID kernel.TenantID, input CreateRateInput) (*shipping.ShippingRate, error) {
	if input.Name == "" {
		return nil, errx.New("rate name is required", errx.TypeValidation)
	}
	now := time.Now().UTC()
	currency := input.PriceCurrency
	if currency == "" {
		currency = "USD"
	}
	rate := &shipping.ShippingRate{
		ID:             kernel.ShippingRateID(uuid.NewString()),
		ZoneID:         input.ZoneID,
		TenantID:       tenantID,
		Name:           input.Name,
		Type:           input.Type,
		Price:          kernel.NewMoney(input.PriceAmount, currency),
		MinWeight:      input.MinWeight,
		MaxWeight:      input.MaxWeight,
		MinOrderAmount: input.MinOrderAmount,
		MaxOrderAmount: input.MaxOrderAmount,
		EstDaysMin:     input.EstDaysMin,
		EstDaysMax:     input.EstDaysMax,
		Active:         true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.rates.Create(ctx, rate); err != nil {
		return nil, err
	}

	if evt, err := eventbus.NewEvent(eventbus.ShippingRateCreated, tenantID, eventbus.ShippingRatePayload{
		RateID: string(rate.ID),
		ZoneID: string(rate.ZoneID),
		Name:   rate.Name,
		Type:   string(rate.Type),
		Price:  rate.Price.Amount,
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return rate, nil
}

// GetRate retrieves a shipping rate by ID.
func (s *Service) GetRate(ctx context.Context, tenantID kernel.TenantID, id kernel.ShippingRateID) (*shipping.ShippingRate, error) {
	return s.rates.GetByID(ctx, tenantID, id)
}

// ListRates returns all rates for a given zone.
func (s *Service) ListRates(ctx context.Context, tenantID kernel.TenantID, zoneID kernel.ShippingZoneID) ([]shipping.ShippingRate, error) {
	return s.rates.ListByZone(ctx, tenantID, zoneID)
}

// UpdateRateInput holds partial update data for a shipping rate.
type UpdateRateInput struct {
	Name           *string
	Type           *shipping.RateType
	PriceAmount    *int64
	PriceCurrency  *string
	MinWeight      *float64
	MaxWeight      *float64
	MinOrderAmount *int64
	MaxOrderAmount *int64
	EstDaysMin     *int
	EstDaysMax     *int
	Active         *bool
}

// UpdateRate modifies an existing shipping rate.
func (s *Service) UpdateRate(ctx context.Context, tenantID kernel.TenantID, id kernel.ShippingRateID, input UpdateRateInput) (*shipping.ShippingRate, error) {
	rate, err := s.rates.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		rate.Name = *input.Name
	}
	if input.Type != nil {
		rate.Type = *input.Type
	}
	if input.PriceAmount != nil {
		rate.Price.Amount = *input.PriceAmount
	}
	if input.PriceCurrency != nil {
		rate.Price.Currency = *input.PriceCurrency
	}
	if input.MinWeight != nil {
		rate.MinWeight = input.MinWeight
	}
	if input.MaxWeight != nil {
		rate.MaxWeight = input.MaxWeight
	}
	if input.MinOrderAmount != nil {
		rate.MinOrderAmount = input.MinOrderAmount
	}
	if input.MaxOrderAmount != nil {
		rate.MaxOrderAmount = input.MaxOrderAmount
	}
	if input.EstDaysMin != nil {
		rate.EstDaysMin = input.EstDaysMin
	}
	if input.EstDaysMax != nil {
		rate.EstDaysMax = input.EstDaysMax
	}
	if input.Active != nil {
		rate.Active = *input.Active
	}
	rate.UpdatedAt = time.Now().UTC()

	if err := s.rates.Update(ctx, rate); err != nil {
		return nil, err
	}
	return rate, nil
}

// DeleteRate removes a shipping rate.
func (s *Service) DeleteRate(ctx context.Context, tenantID kernel.TenantID, id kernel.ShippingRateID) error {
	return s.rates.Delete(ctx, tenantID, id)
}

// ---------------------------------------------------------------------------
// Rate calculation
// ---------------------------------------------------------------------------

// CalculateShipping finds matching zones for the address and returns all active
// rates that satisfy the given weight and order amount constraints.
func (s *Service) CalculateShipping(ctx context.Context, tenantID kernel.TenantID, country, state string, orderAmount int64, weight float64) ([]AvailableRate, error) {
	zones, err := s.zones.FindByAddress(ctx, tenantID, country, state)
	if err != nil {
		return nil, err
	}

	var available []AvailableRate
	seen := map[kernel.ShippingRateID]bool{}

	for _, zone := range zones {
		rates, err := s.rates.ListByZone(ctx, tenantID, zone.ID)
		if err != nil {
			return nil, err
		}

		for _, rate := range rates {
			if !rate.Active {
				continue
			}
			if seen[rate.ID] {
				continue
			}

			// Weight constraints
			if rate.MinWeight != nil && weight < *rate.MinWeight {
				continue
			}
			if rate.MaxWeight != nil && weight > *rate.MaxWeight {
				continue
			}

			// Order amount constraints
			if rate.MinOrderAmount != nil && orderAmount < *rate.MinOrderAmount {
				continue
			}
			if rate.MaxOrderAmount != nil && orderAmount > *rate.MaxOrderAmount {
				continue
			}

			seen[rate.ID] = true
			available = append(available, AvailableRate{
				RateID:     rate.ID,
				Name:       rate.Name,
				Price:      rate.Price,
				EstDaysMin: rate.EstDaysMin,
				EstDaysMax: rate.EstDaysMax,
			})
		}
	}

	if len(available) == 0 {
		return nil, shipping.ErrNoRatesAvailable
	}

	return available, nil
}
