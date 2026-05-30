package cart

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines persistence operations for carts.
type Repository interface {
	Create(ctx context.Context, cart *Cart) error
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CartID) (*Cart, error)
	GetBySession(ctx context.Context, tenantID kernel.TenantID, sessionID string) (*Cart, error)
	GetByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) (*Cart, error)
	Update(ctx context.Context, cart *Cart) error
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CartID) error
	SaveItems(ctx context.Context, cartID kernel.CartID, items []CartItem) error
}
