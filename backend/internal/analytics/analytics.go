package analytics

// DashboardStats holds high-level aggregate metrics for a tenant's dashboard.
type DashboardStats struct {
	TotalProducts  int    `json:"total_products"`
	TotalOrders    int    `json:"total_orders"`
	TotalCustomers int    `json:"total_customers"`
	TotalRevenue   int64  `json:"total_revenue"` // cents
	Currency       string `json:"currency"`
	PendingOrders  int    `json:"pending_orders"`
	ActivePromos   int    `json:"active_promos"`
	PendingPages   int    `json:"pending_pages"`
}

// RevenuePoint represents revenue aggregated for a single day.
type RevenuePoint struct {
	Date   string `json:"date"`   // "2026-01-15"
	Amount int64  `json:"amount"` // cents
}

// TopProduct represents a product ranked by revenue or units sold.
type TopProduct struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	TotalSold   int    `json:"total_sold"`
	Revenue     int64  `json:"revenue"` // cents
}

// OrderStatusBreakdown holds the count of orders grouped by status.
type OrderStatusBreakdown struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// RecentOrder is a lightweight order summary for dashboard display.
type RecentOrder struct {
	ID          string `json:"id"`
	CustomerID  string `json:"customer_id"`
	Status      string `json:"status"`
	TotalAmount int64  `json:"total_amount"`
	Currency    string `json:"currency"`
	CreatedAt   string `json:"created_at"`
}
