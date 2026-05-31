package wishlistsrv

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/wishlist"
)

// Service handles wishlist business logic.
type Service struct {
	repo wishlist.Repository
}

// New creates a new wishlist service.
func New(repo wishlist.Repository) *Service {
	return &Service{repo: repo}
}

// GetOrCreateWishlist finds an existing wishlist for a customer, or creates an empty one.
func (s *Service) GetOrCreateWishlist(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) (*wishlist.Wishlist, error) {
	w, err := s.repo.GetByCustomer(ctx, tenantID, customerID)
	if err == nil {
		return w, nil
	}
	if !errx.IsNotFound(err) {
		return nil, err
	}

	// Not found — create a new one.
	now := time.Now()
	w = &wishlist.Wishlist{
		ID:         kernel.WishlistID(uuid.NewString()),
		TenantID:   tenantID,
		CustomerID: customerID,
		Items:      []wishlist.WishlistItem{},
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.repo.Create(ctx, w); err != nil {
		return nil, errx.Wrap(err, "creating wishlist", errx.TypeInternal)
	}
	return w, nil
}

// GetWishlist retrieves the wishlist for a customer with all items.
func (s *Service) GetWishlist(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) (*wishlist.Wishlist, error) {
	return s.repo.GetByCustomer(ctx, tenantID, customerID)
}

// AddItem adds a product to the customer's wishlist, creating the wishlist if needed.
func (s *Service) AddItem(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, productID kernel.ProductID, variantID string) (*wishlist.Wishlist, error) {
	w, err := s.GetOrCreateWishlist(ctx, tenantID, customerID)
	if err != nil {
		return nil, err
	}

	if w.HasProduct(productID, variantID) {
		return nil, wishlist.ErrAlreadyInWishlist
	}

	item := &wishlist.WishlistItem{
		ID:         kernel.WishlistItemID(uuid.NewString()),
		WishlistID: w.ID,
		ProductID:  productID,
		VariantID:  variantID,
		AddedAt:    time.Now(),
	}

	if err := s.repo.AddItem(ctx, w.ID, item); err != nil {
		return nil, errx.Wrap(err, "adding wishlist item", errx.TypeInternal)
	}

	// Append to in-memory slice for response.
	w.Items = append(w.Items, *item)
	return w, nil
}

// RemoveItem removes an item from the customer's wishlist.
func (s *Service) RemoveItem(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, itemID kernel.WishlistItemID) (*wishlist.Wishlist, error) {
	w, err := s.repo.GetByCustomer(ctx, tenantID, customerID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.RemoveItem(ctx, w.ID, itemID); err != nil {
		return nil, err
	}

	// Update in-memory slice.
	w.RemoveItem(itemID)
	return w, nil
}

// ClearWishlist deletes the customer's entire wishlist.
func (s *Service) ClearWishlist(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) error {
	w, err := s.repo.GetByCustomer(ctx, tenantID, customerID)
	if err != nil {
		if errx.IsNotFound(err) {
			return nil // Already gone — idempotent.
		}
		return err
	}
	return s.repo.Delete(ctx, tenantID, w.ID)
}
