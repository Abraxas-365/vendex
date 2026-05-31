package checkout

import (
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/order"
	"github.com/Abraxas-365/vendex/internal/payment"
)

// CheckoutInput holds all data needed to process a checkout.
type CheckoutInput struct {
	CartID          kernel.CartID
	ShippingAddress order.Address
	BillingAddress  *order.Address // optional, defaults to ShippingAddress if nil
	ShippingRateID  kernel.ShippingRateID
	PromoCode       string // optional
	PaymentProvider string // e.g. "manual", "stripe"
	PaymentMethod   string // e.g. "card", "cash"
	PaymentToken    string // provider-specific token
	CustomerID      kernel.CustomerID
}

// CheckoutResult is returned from a successful Process call.
type CheckoutResult struct {
	Order   *order.Order    `json:"order"`
	Payment *payment.Payment `json:"payment"`
}

// CheckoutSummary is returned from Preview — shows cost breakdown before committing.
type CheckoutSummary struct {
	Subtotal kernel.Money `json:"subtotal"`
	Shipping kernel.Money `json:"shipping"`
	Tax      kernel.Money `json:"tax"`
	Discount kernel.Money `json:"discount"`
	Total    kernel.Money `json:"total"`
	Currency string       `json:"currency"`
}

// Domain errors.
var (
	ErrEmptyCart      = errx.New("cart is empty — add items before checkout", errx.TypeBusiness)
	ErrCheckoutFailed = errx.New("checkout could not be completed", errx.TypeBusiness)
	ErrPaymentFailed  = errx.New("payment processing failed", errx.TypeExternal)
)
