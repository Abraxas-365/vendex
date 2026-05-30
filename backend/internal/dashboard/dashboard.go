package dashboard

import "time"

// DateRange defines the time window for reporting queries.
type DateRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// SalesOverview holds high-level sales KPIs for a date range.
type SalesOverview struct {
	TotalRevenue     int64  `json:"total_revenue"`      // cents
	OrderCount       int    `json:"order_count"`
	AverageOrderValue int64 `json:"average_order_value"` // cents
	RefundTotal      int64  `json:"refund_total"`        // cents (cancelled orders)
	Currency         string `json:"currency"`
}

// TopProduct is a product ranked by revenue in a date range.
type TopProduct struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Revenue   int64  `json:"revenue"`   // cents
	Quantity  int    `json:"quantity"`
}

// DailyRevenue is the revenue aggregated for a single calendar day.
type DailyRevenue struct {
	Date       string `json:"date"`        // "2024-01-15"
	Revenue    int64  `json:"revenue"`     // cents
	OrderCount int    `json:"order_count"`
}

// CustomerStats holds customer metrics for a date range.
type CustomerStats struct {
	TotalCustomers     int `json:"total_customers"`
	NewCustomers       int `json:"new_customers"`
	ReturningCustomers int `json:"returning_customers"`
}

// ConversionFunnel represents a simple checkout conversion funnel.
// Visitors and cart/checkout stages are tracked by external analytics;
// we provide the order-completion count from the DB.
type ConversionFunnel struct {
	Visitors          int `json:"visitors"`           // placeholder — not stored in DB
	CartsCreated      int `json:"carts_created"`
	CheckoutsStarted  int `json:"checkouts_started"`  // placeholder — not stored in DB
	OrdersCompleted   int `json:"orders_completed"`
}
