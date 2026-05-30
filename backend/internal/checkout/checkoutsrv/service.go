package checkoutsrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/cart"
	"github.com/Abraxas-365/hada-commerce/internal/checkout"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/order"
	"github.com/Abraxas-365/hada-commerce/internal/order/ordersrv"
	"github.com/Abraxas-365/hada-commerce/internal/payment"
	"github.com/Abraxas-365/hada-commerce/internal/shipping"
	"github.com/Abraxas-365/hada-commerce/internal/tax"
)

// ---------------------------------------------------------------------------
// Dependency interfaces
// ---------------------------------------------------------------------------

// CartGetter abstracts cart retrieval and clearing.
type CartGetter interface {
	GetCart(ctx context.Context, tenantID kernel.TenantID, cartID kernel.CartID) (*cart.Cart, error)
	ClearCart(ctx context.Context, tenantID kernel.TenantID, cartID kernel.CartID) (*cart.Cart, error)
}

// OrderCreator abstracts order creation and status transitions.
type OrderCreator interface {
	Create(ctx context.Context, tenantID kernel.TenantID, in ordersrv.CreateInput) (*order.Order, error)
	UpdateStatus(ctx context.Context, tenantID kernel.TenantID, id kernel.OrderID, status order.OrderStatus) (*order.Order, error)
}

// OrderCheckoutUpdater abstracts the repo method for writing checkout fields.
type OrderCheckoutUpdater interface {
	UpdateCheckoutFields(ctx context.Context, o *order.Order) error
}

// ShippingGetter abstracts fetching a shipping rate.
type ShippingGetter interface {
	GetRate(ctx context.Context, tenantID kernel.TenantID, id kernel.ShippingRateID) (*shipping.ShippingRate, error)
}

// TaxCalculator abstracts tax computation.
type TaxCalculator interface {
	CalculateTax(ctx context.Context, tenantID kernel.TenantID, subtotalCents int64, shippingCents int64, country, state, city, zipCode string) (*tax.TaxResult, error)
}

// PaymentProcessor abstracts payment creation and processing.
type PaymentProcessor interface {
	CreatePayment(ctx context.Context, tenantID kernel.TenantID, orderID kernel.OrderID, amount int64, currency string, providerName string, method string) (*payment.Payment, error)
	ProcessPayment(ctx context.Context, tenantID kernel.TenantID, paymentID kernel.PaymentID, token string) (*payment.Payment, error)
}

// PromoApplier abstracts promo validation and application.
type PromoApplier interface {
	Validate(ctx context.Context, tenantID kernel.TenantID, code string, orderTotalCents int64) (ValidationResult, error)
	Apply(ctx context.Context, tenantID kernel.TenantID, code string, orderTotalCents int64) (int64, error)
}

// ValidationResult mirrors promosrv.ValidationResult to avoid an import cycle.
type ValidationResult struct {
	Valid          bool
	DiscountCents  int64
	IsFreeShipping bool
	Reason         string
}

// ---------------------------------------------------------------------------
// Service
// ---------------------------------------------------------------------------

// Service orchestrates the checkout flow: cart → shipping → tax → order → payment.
type Service struct {
	cartSvc     CartGetter
	orderSvc    OrderCreator
	orderRepo   OrderCheckoutUpdater
	shippingSvc ShippingGetter
	taxSvc      TaxCalculator
	paymentSvc  PaymentProcessor
	promoSvc    PromoApplier
	bus         eventbus.Bus
}

// New constructs a checkout Service.
func New(
	cartSvc CartGetter,
	orderSvc OrderCreator,
	orderRepo OrderCheckoutUpdater,
	shippingSvc ShippingGetter,
	taxSvc TaxCalculator,
	paymentSvc PaymentProcessor,
	promoSvc PromoApplier,
	bus eventbus.Bus,
) *Service {
	return &Service{
		cartSvc:     cartSvc,
		orderSvc:    orderSvc,
		orderRepo:   orderRepo,
		shippingSvc: shippingSvc,
		taxSvc:      taxSvc,
		paymentSvc:  paymentSvc,
		promoSvc:    promoSvc,
		bus:         bus,
	}
}

// ---------------------------------------------------------------------------
// Preview — read-only summary (no mutations)
// ---------------------------------------------------------------------------

// Preview calculates a cost breakdown for the given checkout input without committing anything.
func (s *Service) Preview(ctx context.Context, tenantID kernel.TenantID, input checkout.CheckoutInput) (*checkout.CheckoutSummary, error) {
	// 1. Get cart and validate it has items.
	c, err := s.cartSvc.GetCart(ctx, tenantID, input.CartID)
	if err != nil {
		return nil, errx.Wrap(err, "getting cart for preview", errx.TypeInternal)
	}
	if len(c.Items) == 0 {
		return nil, checkout.ErrEmptyCart
	}

	// 2. Calculate subtotal from cart.
	subtotal := c.Subtotal()

	// 3. Get shipping rate.
	rate, err := s.shippingSvc.GetRate(ctx, tenantID, input.ShippingRateID)
	if err != nil {
		return nil, errx.Wrap(err, "getting shipping rate", errx.TypeInternal)
	}
	shippingMoney := rate.Price

	// 4. Calculate tax.
	addr := input.ShippingAddress
	taxResult, err := s.taxSvc.CalculateTax(
		ctx, tenantID,
		subtotal.Amount, shippingMoney.Amount,
		addr.Country, addr.State, addr.City, addr.PostalCode,
	)
	if err != nil {
		// Tax errors are non-fatal in preview — treat as zero tax.
		taxResult = &tax.TaxResult{TotalTax: 0}
	}
	taxMoney := kernel.NewMoney(taxResult.TotalTax, subtotal.Currency)

	// 5. Validate promo code (read-only — does NOT increment usage).
	discountMoney := kernel.NewMoney(0, subtotal.Currency)
	freeShipping := false
	if input.PromoCode != "" {
		subtotalForPromo := subtotal.Amount + shippingMoney.Amount + taxMoney.Amount
		result, err := s.promoSvc.Validate(ctx, tenantID, input.PromoCode, subtotalForPromo)
		if err == nil && result.Valid {
			discountMoney = kernel.NewMoney(result.DiscountCents, subtotal.Currency)
			freeShipping = result.IsFreeShipping
		}
	}

	if freeShipping {
		shippingMoney = kernel.NewMoney(0, subtotal.Currency)
	}

	// 6. Compute total.
	total := subtotal.Amount + shippingMoney.Amount + taxMoney.Amount - discountMoney.Amount
	if total < 0 {
		total = 0
	}

	return &checkout.CheckoutSummary{
		Subtotal: subtotal,
		Shipping: shippingMoney,
		Tax:      taxMoney,
		Discount: discountMoney,
		Total:    kernel.NewMoney(total, subtotal.Currency),
		Currency: subtotal.Currency,
	}, nil
}

// ---------------------------------------------------------------------------
// Process — full checkout orchestration
// ---------------------------------------------------------------------------

// Process converts a cart into a confirmed, paid order.
func (s *Service) Process(ctx context.Context, tenantID kernel.TenantID, input checkout.CheckoutInput) (*checkout.CheckoutResult, error) {
	// 1. Publish CheckoutStarted.
	if evt, err := eventbus.NewEvent(eventbus.CheckoutStarted, tenantID, map[string]string{
		"cart_id":     string(input.CartID),
		"customer_id": string(input.CustomerID),
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	// 2. Get cart and validate.
	c, err := s.cartSvc.GetCart(ctx, tenantID, input.CartID)
	if err != nil {
		s.publishFailed(ctx, tenantID, input, "cart not found")
		return nil, errx.Wrap(err, "getting cart", errx.TypeInternal)
	}
	if len(c.Items) == 0 {
		s.publishFailed(ctx, tenantID, input, "cart is empty")
		return nil, checkout.ErrEmptyCart
	}

	// 3. Get shipping rate.
	rate, err := s.shippingSvc.GetRate(ctx, tenantID, input.ShippingRateID)
	if err != nil {
		s.publishFailed(ctx, tenantID, input, "shipping rate not found")
		return nil, errx.Wrap(err, "getting shipping rate", errx.TypeInternal)
	}
	shippingMoney := rate.Price

	// 4. Calculate subtotal from cart items.
	subtotal := c.Subtotal()

	// 5. Calculate tax.
	addr := input.ShippingAddress
	taxResult, err := s.taxSvc.CalculateTax(
		ctx, tenantID,
		subtotal.Amount, shippingMoney.Amount,
		addr.Country, addr.State, addr.City, addr.PostalCode,
	)
	if err != nil {
		// Tax errors are non-fatal — fall back to zero tax.
		taxResult = &tax.TaxResult{TotalTax: 0}
	}
	taxMoney := kernel.NewMoney(taxResult.TotalTax, subtotal.Currency)

	// 6. Apply promo if provided (increments usage counter).
	discountMoney := kernel.NewMoney(0, subtotal.Currency)
	freeShipping := false
	appliedPromo := ""
	if input.PromoCode != "" {
		subtotalForPromo := subtotal.Amount + shippingMoney.Amount + taxMoney.Amount
		discountCents, err := s.promoSvc.Apply(ctx, tenantID, input.PromoCode, subtotalForPromo)
		if err != nil {
			// Promo errors are non-fatal for process — skip discount.
			discountMoney = kernel.NewMoney(0, subtotal.Currency)
		} else {
			discountMoney = kernel.NewMoney(discountCents, subtotal.Currency)
			appliedPromo = input.PromoCode

			// Check if the promo grants free shipping via a fresh Validate call.
			result, verr := s.promoSvc.Validate(ctx, tenantID, input.PromoCode, subtotalForPromo)
			if verr == nil && result.IsFreeShipping {
				freeShipping = true
			}
		}
	}

	if freeShipping {
		shippingMoney = kernel.NewMoney(0, subtotal.Currency)
	}

	// 7. Calculate total.
	totalCents := subtotal.Amount + shippingMoney.Amount + taxMoney.Amount - discountMoney.Amount
	if totalCents < 0 {
		totalCents = 0
	}

	// 8. Build order items from cart items.
	orderItems := make([]ordersrv.CreateItemInput, len(c.Items))
	for i, item := range c.Items {
		orderItems[i] = ordersrv.CreateItemInput{
			ProductID:   item.ProductID,
			ProductName: string(item.ProductID), // fallback; enrichment handled upstream
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		}
	}

	// 9. Create order.
	newOrder, err := s.orderSvc.Create(ctx, tenantID, ordersrv.CreateInput{
		CustomerID:      input.CustomerID,
		Items:           orderItems,
		ShippingAddress: input.ShippingAddress,
	})
	if err != nil {
		s.publishFailed(ctx, tenantID, input, "order creation failed")
		return nil, errx.Wrap(err, "creating order", errx.TypeInternal)
	}

	// 10. Enrich order with checkout fields.
	billing := input.BillingAddress
	if billing == nil {
		billing = &input.ShippingAddress
	}
	now := time.Now()
	newOrder.SubtotalAmount = subtotal
	newOrder.ShippingAmount = shippingMoney
	newOrder.TaxAmount = taxMoney
	newOrder.DiscountAmount = discountMoney
	newOrder.TotalAmount = kernel.NewMoney(totalCents, subtotal.Currency)
	newOrder.ShippingMethod = rate.Name
	newOrder.BillingAddress = billing
	newOrder.PaymentStatus = "pending"
	newOrder.PaymentMethod = input.PaymentMethod
	newOrder.PromoCode = appliedPromo
	newOrder.CartID = string(input.CartID)
	newOrder.UpdatedAt = now

	if err := s.orderRepo.UpdateCheckoutFields(ctx, newOrder); err != nil {
		s.publishFailed(ctx, tenantID, input, "updating checkout fields failed")
		return nil, errx.Wrap(err, "persisting checkout fields", errx.TypeInternal)
	}

	// 11. Create payment record.
	pay, err := s.paymentSvc.CreatePayment(
		ctx, tenantID, newOrder.ID,
		totalCents, subtotal.Currency,
		input.PaymentProvider, input.PaymentMethod,
	)
	if err != nil {
		// Cancel order on payment creation failure.
		_, _ = s.orderSvc.UpdateStatus(ctx, tenantID, newOrder.ID, order.StatusCancelled)
		s.publishFailed(ctx, tenantID, input, "payment creation failed")
		return nil, errx.Wrap(err, "creating payment", errx.TypeInternal)
	}

	// 12. Process (charge) payment.
	pay, err = s.paymentSvc.ProcessPayment(ctx, tenantID, pay.ID, input.PaymentToken)
	if err != nil {
		// Cancel order on payment failure.
		_, _ = s.orderSvc.UpdateStatus(ctx, tenantID, newOrder.ID, order.StatusCancelled)
		s.publishFailed(ctx, tenantID, input, "payment processing failed")
		return nil, errx.Wrap(err, "processing payment", errx.TypeExternal)
	}

	// 13. Confirm order and clear cart.
	confirmedOrder, err := s.orderSvc.UpdateStatus(ctx, tenantID, newOrder.ID, order.StatusConfirmed)
	if err != nil {
		// Payment succeeded but status update failed — log and continue.
		confirmedOrder = newOrder
	}

	// Update payment_status on the order record (best-effort).
	confirmedOrder.PaymentStatus = "paid"
	confirmedOrder.UpdatedAt = time.Now()
	_ = s.orderRepo.UpdateCheckoutFields(ctx, confirmedOrder)

	// Clear the cart (best-effort — don't fail checkout if this errors).
	_, _ = s.cartSvc.ClearCart(ctx, tenantID, input.CartID)

	// 14. Publish CheckoutCompleted.
	if evt, err := eventbus.NewEvent(eventbus.CheckoutCompleted, tenantID, map[string]interface{}{
		"order_id":    string(confirmedOrder.ID),
		"cart_id":     string(input.CartID),
		"customer_id": string(input.CustomerID),
		"total":       totalCents,
		"currency":    subtotal.Currency,
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return &checkout.CheckoutResult{
		Order:   confirmedOrder,
		Payment: pay,
	}, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func (s *Service) publishFailed(ctx context.Context, tenantID kernel.TenantID, input checkout.CheckoutInput, reason string) {
	if evt, err := eventbus.NewEvent(eventbus.CheckoutFailed, tenantID, map[string]string{
		"cart_id":     string(input.CartID),
		"customer_id": string(input.CustomerID),
		"reason":      reason,
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}
}
