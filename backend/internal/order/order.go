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
	ID          kernel.OrderItemID `json:"id"`
	ProductID   kernel.ProductID   `json:"product_id"`
	ProductName string             `json:"product_name"`
	Quantity    int                `json:"quantity"`
	UnitPrice   kernel.Money       `json:"unit_price"`
	Total       kernel.Money       `json:"total"`
}

// Order is the aggregate root for a purchase.
type Order struct {
	ID              kernel.OrderID    `json:"id"`
	TenantID        kernel.TenantID   `json:"tenant_id"`
	CustomerID      kernel.CustomerID `json:"customer_id"`
	Items           []OrderItem       `json:"items"`
	Status          OrderStatus       `json:"status"`
	TotalAmount     kernel.Money      `json:"total_amount"`
	ShippingAddress Address           `json:"shipping_address"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
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
