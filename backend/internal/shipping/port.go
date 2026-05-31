package shipping

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// ZoneRepository defines persistence operations for shipping zones.
type ZoneRepository interface {
	Create(ctx context.Context, zone *ShippingZone) error
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ShippingZoneID) (*ShippingZone, error)
	List(ctx context.Context, tenantID kernel.TenantID) ([]ShippingZone, error)
	Update(ctx context.Context, zone *ShippingZone) error
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.ShippingZoneID) error
	FindByAddress(ctx context.Context, tenantID kernel.TenantID, country, state string) ([]ShippingZone, error)
}

// RateRepository defines persistence operations for shipping rates.
type RateRepository interface {
	Create(ctx context.Context, rate *ShippingRate) error
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ShippingRateID) (*ShippingRate, error)
	ListByZone(ctx context.Context, tenantID kernel.TenantID, zoneID kernel.ShippingZoneID) ([]ShippingRate, error)
	Update(ctx context.Context, rate *ShippingRate) error
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.ShippingRateID) error
}
