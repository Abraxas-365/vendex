package tax

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Repository defines the persistence interface for the tax domain.
type Repository interface {
	Create(ctx context.Context, rate *TaxRate) error
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.TaxRateID) (*TaxRate, error)
	List(ctx context.Context, tenantID kernel.TenantID) ([]TaxRate, error)
	Update(ctx context.Context, rate *TaxRate) error
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.TaxRateID) error
	FindByLocation(ctx context.Context, tenantID kernel.TenantID, country, state, city, zipCode string) ([]TaxRate, error)
}
