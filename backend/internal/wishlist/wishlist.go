package wishlist

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Wishlist is the core entity for a customer's saved product list.
type Wishlist struct {
	ID         kernel.WishlistID `json:"id" db:"id"`
	TenantID   kernel.TenantID   `json:"tenant_id" db:"tenant_id"`
	CustomerID kernel.CustomerID `json:"customer_id" db:"customer_id"`
	Items      []WishlistItem    `json:"items"`
	CreatedAt  time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at" db:"updated_at"`
}

// WishlistItem represents a product saved in the wishlist.
type WishlistItem struct {
	ID         kernel.WishlistItemID `json:"id" db:"id"`
	WishlistID kernel.WishlistID     `json:"wishlist_id" db:"wishlist_id"`
	ProductID  kernel.ProductID      `json:"product_id" db:"product_id"`
	VariantID  string                `json:"variant_id,omitempty" db:"variant_id"`
	AddedAt    time.Time             `json:"added_at" db:"added_at"`
}

// --- Domain errors ---

var (
	ErrNotFound         = errx.New("wishlist not found", errx.TypeNotFound)
	ErrItemNotFound     = errx.New("wishlist item not found", errx.TypeNotFound)
	ErrAlreadyInWishlist = errx.New("product already in wishlist", errx.TypeBusiness)
)

// --- Domain methods ---

// AddItem adds an item to the wishlist. Returns ErrAlreadyInWishlist if the product+variant already exists.
func (w *Wishlist) AddItem(item WishlistItem) error {
	for _, existing := range w.Items {
		if existing.ProductID == item.ProductID && existing.VariantID == item.VariantID {
			return ErrAlreadyInWishlist
		}
	}
	w.Items = append(w.Items, item)
	return nil
}

// RemoveItem removes the item with the given ID from the wishlist.
func (w *Wishlist) RemoveItem(itemID kernel.WishlistItemID) {
	filtered := w.Items[:0]
	for _, item := range w.Items {
		if item.ID != itemID {
			filtered = append(filtered, item)
		}
	}
	w.Items = filtered
}

// HasProduct returns true if the wishlist contains a product (with optional variant).
func (w *Wishlist) HasProduct(productID kernel.ProductID, variantID string) bool {
	for _, item := range w.Items {
		if item.ProductID == productID && item.VariantID == variantID {
			return true
		}
	}
	return false
}

// ItemCount returns the number of items in the wishlist.
func (w *Wishlist) ItemCount() int {
	return len(w.Items)
}
