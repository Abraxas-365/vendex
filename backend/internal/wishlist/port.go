package wishlist

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines the persistence interface for the wishlist domain.
type Repository interface {
	// Create inserts a new wishlist (without items).
	Create(ctx context.Context, w *Wishlist) error

	// GetByCustomer retrieves a wishlist by customer, scoped to tenant, with items loaded.
	GetByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) (*Wishlist, error)

	// GetByID retrieves a wishlist by ID, scoped to tenant, with items loaded.
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.WishlistID) (*Wishlist, error)

	// AddItem inserts a new item into the wishlist.
	AddItem(ctx context.Context, wishlistID kernel.WishlistID, item *WishlistItem) error

	// RemoveItem deletes an item from the wishlist.
	RemoveItem(ctx context.Context, wishlistID kernel.WishlistID, itemID kernel.WishlistItemID) error

	// Delete removes the wishlist and all its items (cascade).
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.WishlistID) error
}
