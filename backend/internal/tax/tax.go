package tax

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// TaxRate represents a tax rate configuration for a jurisdiction.
type TaxRate struct {
	ID               kernel.TaxRateID `json:"id" db:"id"`
	TenantID         kernel.TenantID  `json:"tenant_id" db:"tenant_id"`
	Name             string           `json:"name" db:"name"`
	Rate             float64          `json:"rate" db:"rate"`               // e.g. 0.0825 for 8.25%
	Country          string           `json:"country" db:"country"`         // ISO 3166-1 alpha-2
	State            string           `json:"state,omitempty" db:"state"`   // nullable
	City             string           `json:"city,omitempty" db:"city"`     // nullable
	ZipCode          string           `json:"zip_code,omitempty" db:"zip_code"` // nullable
	Priority         int              `json:"priority" db:"priority"`       // for ordering compound taxes
	Compound         bool             `json:"compound" db:"compound"`       // applied on top of previous taxes
	IncludesShipping bool             `json:"includes_shipping" db:"includes_shipping"`
	Active           bool             `json:"active" db:"active"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at" db:"updated_at"`
}

// TaxResult holds the result of a tax calculation.
type TaxResult struct {
	TotalTax     int64         `json:"total_tax"`      // cents
	TaxBreakdown []TaxLineItem `json:"tax_breakdown"`
}

// TaxLineItem represents a single tax rate's contribution to the total.
type TaxLineItem struct {
	RateID   kernel.TaxRateID `json:"rate_id"`
	Name     string           `json:"name"`
	Rate     float64          `json:"rate"`
	Amount   int64            `json:"amount"` // cents
	Compound bool             `json:"compound"`
}

// Domain errors.
var (
	ErrNotFound          = errx.New("tax rate not found", errx.TypeNotFound)
	ErrInvalidRate       = errx.New("tax rate must be between 0 and 1", errx.TypeValidation)
	ErrNoRatesConfigured = errx.New("no tax rates configured for this location", errx.TypeBusiness)
)
