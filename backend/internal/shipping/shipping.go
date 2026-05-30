package shipping

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// RateType represents how the shipping rate price is determined.
type RateType string

const (
	RateFlat        RateType = "flat"
	RateWeightBased RateType = "weight_based"
	RatePriceBased  RateType = "price_based"
	RateFree        RateType = "free"
)

// ShippingZone defines geographic coverage for shipping.
type ShippingZone struct {
	ID        kernel.ShippingZoneID `json:"id" db:"id"`
	TenantID  kernel.TenantID       `json:"tenant_id" db:"tenant_id"`
	Name      string                `json:"name" db:"name"`
	Countries []string              `json:"countries"`  // JSONB
	States    []string              `json:"states"`     // JSONB
	CreatedAt time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt time.Time             `json:"updated_at" db:"updated_at"`
}

// MatchesAddress returns true if the zone covers the given country and state.
// A zone with an empty Countries list matches any country.
// A zone with an empty States list matches any state within the matched country.
func (z *ShippingZone) MatchesAddress(country, state string) bool {
	if len(z.Countries) == 0 {
		return true
	}
	countryMatch := false
	for _, c := range z.Countries {
		if c == country {
			countryMatch = true
			break
		}
	}
	if !countryMatch {
		return false
	}

	if len(z.States) == 0 {
		return true
	}
	for _, s := range z.States {
		if s == state {
			return true
		}
	}
	return false
}

// ShippingRate defines cost and conditions for a shipping option within a zone.
type ShippingRate struct {
	ID             kernel.ShippingRateID `json:"id" db:"id"`
	ZoneID         kernel.ShippingZoneID `json:"zone_id" db:"zone_id"`
	TenantID       kernel.TenantID       `json:"tenant_id" db:"tenant_id"`
	Name           string                `json:"name" db:"name"`
	Type           RateType              `json:"type" db:"type"`
	Price          kernel.Money          `json:"price"`
	MinWeight      *float64              `json:"min_weight,omitempty" db:"min_weight"`
	MaxWeight      *float64              `json:"max_weight,omitempty" db:"max_weight"`
	MinOrderAmount *int64                `json:"min_order_amount,omitempty" db:"min_order_amount"`
	MaxOrderAmount *int64                `json:"max_order_amount,omitempty" db:"max_order_amount"`
	EstDaysMin     *int                  `json:"est_days_min,omitempty" db:"est_days_min"`
	EstDaysMax     *int                  `json:"est_days_max,omitempty" db:"est_days_max"`
	Active         bool                  `json:"active" db:"active"`
	CreatedAt      time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at" db:"updated_at"`
}
