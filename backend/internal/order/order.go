package order

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// OrderStatus represents the lifecycle of an order.
type OrderStatus string

const (
	StatusPending    OrderStatus = "pending"
	StatusConfirmed  OrderStatus = "confirmed"
	StatusProcessing OrderStatus = "processing"
	StatusShipped    OrderStatus = "shipped"
	StatusDelivered  OrderStatus = "delivered"
	StatusCancelled  OrderStatus = "cancelled"
)

// validTransitions defines allowed status transitions.
var validTransitions = map[OrderStatus][]OrderStatus{
	StatusPending:    {StatusConfirmed, StatusCancelled},
	StatusConfirmed:  {StatusProcessing, StatusCancelled},
	StatusProcessing: {StatusShipped, StatusCancelled},
	StatusShipped:    {StatusDelivered},
	StatusDelivered:  {},
	StatusCancelled:  {},
}

// Address represents a shipping address.
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
}

// OrderItem represents a line item in an order.
type OrderItem struct {
	ID          kernel.OrderItemID `json:"id" db:"id"`
	ProductID   kernel.ProductID   `json:"product_id" db:"product_id"`
	ProductName string             `json:"product_name" db:"product_name"`
	Quantity    int                `json:"quantity" db:"quantity"`
	UnitPrice   kernel.Money       `json:"unit_price"`
	Total       kernel.Money       `json:"total"`
}

// Order is the aggregate root for a purchase.
type Order struct {
	ID              kernel.OrderID    `json:"id" db:"id"`
	TenantID        kernel.TenantID   `json:"tenant_id" db:"tenant_id"`
	CustomerID      kernel.CustomerID `json:"customer_id" db:"customer_id"`
	Items           []OrderItem       `json:"items"`
	Status          OrderStatus       `json:"status" db:"status"`
	TotalAmount     kernel.Money      `json:"total_amount"`
	ShippingAddress Address           `json:"shipping_address"`
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" db:"updated_at"`

	// Checkout-enriched fields (populated after checkout orchestration)
	SubtotalAmount kernel.Money `json:"subtotal_amount"`
	ShippingAmount kernel.Money `json:"shipping_amount"`
	TaxAmount      kernel.Money `json:"tax_amount"`
	DiscountAmount kernel.Money `json:"discount_amount"`
	ShippingMethod string       `json:"shipping_method,omitempty" db:"shipping_method"`
	BillingAddress *Address     `json:"billing_address,omitempty"`
	PaymentStatus  string       `json:"payment_status,omitempty" db:"payment_status"`
	PaymentMethod  string       `json:"payment_method,omitempty" db:"payment_method"`
	TrackingNumber string       `json:"tracking_number,omitempty" db:"tracking_number"`
	Carrier        string       `json:"carrier,omitempty" db:"carrier"`
	PromoCode      string       `json:"promo_code,omitempty" db:"promo_code"`
	CartID         string       `json:"cart_id,omitempty" db:"cart_id"`
}

// CanTransitionTo checks if the status transition is valid.
func (o *Order) CanTransitionTo(next OrderStatus) bool {
	allowed, ok := validTransitions[o.Status]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == next {
			return true
		}
	}
	return false
}

// TransitionTo updates the order status if the transition is valid.
// Returns false if the transition is not allowed.
func (o *Order) TransitionTo(next OrderStatus) bool {
	if !o.CanTransitionTo(next) {
		return false
	}
	o.Status = next
	o.UpdatedAt = time.Now()
	return true
}

// Cancel transitions the order to cancelled status.
func (o *Order) Cancel() bool {
	return o.TransitionTo(StatusCancelled)
}

// CalculateTotal recomputes the total amount from line items.
func (o *Order) CalculateTotal() {
	var total int64
	currency := ""
	for i := range o.Items {
		o.Items[i].Total = o.Items[i].UnitPrice.Multiply(o.Items[i].Quantity)
		total += o.Items[i].Total.Amount
		if currency == "" {
			currency = o.Items[i].UnitPrice.Currency
		}
	}
	o.TotalAmount = kernel.NewMoney(total, currency)
}

// ItemCount returns the total number of units across all items.
func (o *Order) ItemCount() int {
	count := 0
	for _, item := range o.Items {
		count += item.Quantity
	}
	return count
}
