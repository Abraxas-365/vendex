package cartsrv

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Abraxas-365/vendex/internal/cart"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Service handles cart business logic.
type Service struct {
	repo cart.Repository
	bus  eventbus.Bus
}

// New creates a new cart service.
func New(repo cart.Repository, bus eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// CreateCart creates a new empty cart for the given tenant.
func (s *Service) CreateCart(ctx context.Context, tenantID kernel.TenantID, sessionID string, customerID kernel.CustomerID, currency string) (*cart.Cart, error) {
	if currency == "" {
		currency = "USD"
	}
	now := time.Now()
	c := &cart.Cart{
		ID:         kernel.CartID(uuid.NewString()),
		TenantID:   tenantID,
		CustomerID: customerID,
		SessionID:  sessionID,
		Items:      []cart.CartItem{},
		Currency:   currency,
		CreatedAt:  now,
		UpdatedAt:  now,
		ExpiresAt:  now.AddDate(0, 0, 30),
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, errx.Wrap(err, "creating cart", errx.TypeInternal)
	}

	if evt, err := eventbus.NewEvent(eventbus.CartCreated, tenantID, eventbus.CartPayload{
		CartID:     string(c.ID),
		CustomerID: string(c.CustomerID),
		ItemCount:  0,
		Total:      0,
		Currency:   c.Currency,
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return c, nil
}

// GetOrCreateCart retrieves an existing cart by session or customer, or creates a new one.
func (s *Service) GetOrCreateCart(ctx context.Context, tenantID kernel.TenantID, sessionID string, customerID kernel.CustomerID, currency string) (*cart.Cart, error) {
	// Try by customer first.
	if !customerID.IsEmpty() {
		c, err := s.repo.GetByCustomer(ctx, tenantID, customerID)
		if err == nil {
			return c, nil
		}
		if !errx.IsNotFound(err) {
			return nil, err
		}
	}

	// Try by session.
	if sessionID != "" {
		c, err := s.repo.GetBySession(ctx, tenantID, sessionID)
		if err == nil {
			return c, nil
		}
		if !errx.IsNotFound(err) {
			return nil, err
		}
	}

	// Create a new cart.
	return s.CreateCart(ctx, tenantID, sessionID, customerID, currency)
}

// GetCart retrieves a cart by ID, scoped to the tenant.
func (s *Service) GetCart(ctx context.Context, tenantID kernel.TenantID, cartID kernel.CartID) (*cart.Cart, error) {
	return s.repo.GetByID(ctx, tenantID, cartID)
}

// AddItem adds (or merges) an item into the cart.
func (s *Service) AddItem(ctx context.Context, tenantID kernel.TenantID, cartID kernel.CartID, productID kernel.ProductID, variantID string, quantity int, unitPrice kernel.Money) (*cart.Cart, error) {
	if quantity <= 0 {
		return nil, cart.ErrInvalidQty
	}

	c, err := s.repo.GetByID(ctx, tenantID, cartID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	item := cart.CartItem{
		ID:        kernel.CartItemID(uuid.NewString()),
		CartID:    cartID,
		TenantID:  tenantID,
		ProductID: productID,
		VariantID: variantID,
		Quantity:  quantity,
		UnitPrice: unitPrice,
		CreatedAt: now,
		UpdatedAt: now,
	}
	c.AddItem(item)
	c.UpdatedAt = now

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, errx.Wrap(err, "updating cart", errx.TypeInternal)
	}
	if err := s.repo.SaveItems(ctx, cartID, c.Items); err != nil {
		return nil, errx.Wrap(err, "saving cart items", errx.TypeInternal)
	}

	s.publishCartUpdated(ctx, c)
	return c, nil
}

// UpdateItemQuantity changes the quantity of a specific cart item.
func (s *Service) UpdateItemQuantity(ctx context.Context, tenantID kernel.TenantID, cartID kernel.CartID, itemID kernel.CartItemID, quantity int) (*cart.Cart, error) {
	c, err := s.repo.GetByID(ctx, tenantID, cartID)
	if err != nil {
		return nil, err
	}

	if err := c.UpdateItemQuantity(itemID, quantity); err != nil {
		return nil, err
	}
	c.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, errx.Wrap(err, "updating cart", errx.TypeInternal)
	}
	if err := s.repo.SaveItems(ctx, cartID, c.Items); err != nil {
		return nil, errx.Wrap(err, "saving cart items", errx.TypeInternal)
	}

	s.publishCartUpdated(ctx, c)
	return c, nil
}

// RemoveItem removes an item from the cart.
func (s *Service) RemoveItem(ctx context.Context, tenantID kernel.TenantID, cartID kernel.CartID, itemID kernel.CartItemID) (*cart.Cart, error) {
	c, err := s.repo.GetByID(ctx, tenantID, cartID)
	if err != nil {
		return nil, err
	}

	c.RemoveItem(itemID)
	c.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, errx.Wrap(err, "updating cart", errx.TypeInternal)
	}
	if err := s.repo.SaveItems(ctx, cartID, c.Items); err != nil {
		return nil, errx.Wrap(err, "saving cart items", errx.TypeInternal)
	}

	s.publishCartUpdated(ctx, c)
	return c, nil
}

// ClearCart removes all items from the cart.
func (s *Service) ClearCart(ctx context.Context, tenantID kernel.TenantID, cartID kernel.CartID) (*cart.Cart, error) {
	c, err := s.repo.GetByID(ctx, tenantID, cartID)
	if err != nil {
		return nil, err
	}

	c.Clear()
	c.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, errx.Wrap(err, "updating cart", errx.TypeInternal)
	}
	if err := s.repo.SaveItems(ctx, cartID, c.Items); err != nil {
		return nil, errx.Wrap(err, "saving cart items", errx.TypeInternal)
	}

	return c, nil
}

// DeleteCart removes the cart and all its items.
func (s *Service) DeleteCart(ctx context.Context, tenantID kernel.TenantID, cartID kernel.CartID) error {
	return s.repo.Delete(ctx, tenantID, cartID)
}

// publishCartUpdated is a helper to fire the CartUpdated event.
func (s *Service) publishCartUpdated(ctx context.Context, c *cart.Cart) {
	subtotal := c.Subtotal()
	evt, err := eventbus.NewEvent(eventbus.CartUpdated, c.TenantID, eventbus.CartPayload{
		CartID:     string(c.ID),
		CustomerID: string(c.CustomerID),
		ItemCount:  c.ItemCount(),
		Total:      int(subtotal.Amount),
		Currency:   c.Currency,
	})
	if err == nil {
		_ = s.bus.Publish(ctx, evt)
	}
}
