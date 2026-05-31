package currency

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Repository defines the persistence interface for the currency domain.
type Repository interface {
	Create(ctx context.Context, rate *CurrencyRate) error
	GetByPair(ctx context.Context, tenantID kernel.TenantID, base, target string) (*CurrencyRate, error)
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CurrencyRateID) (*CurrencyRate, error)
	List(ctx context.Context, tenantID kernel.TenantID) ([]CurrencyRate, error)
	Update(ctx context.Context, rate *CurrencyRate) error
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CurrencyRateID) error
}
