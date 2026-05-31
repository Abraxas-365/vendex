package eventbus

// OrderPayload is the payload for order.* events.
type OrderPayload struct {
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
	Status     string `json:"status"`
	Total      int    `json:"total"`
	Currency   string `json:"currency"`
	ItemCount  int    `json:"item_count"`
}

// ProductPayload is the payload for product.* events.
type ProductPayload struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	SKU       string `json:"sku"`
	Price     int    `json:"price"`
	Currency  string `json:"currency"`
}

// CustomerPayload is the payload for customer.* events.
type CustomerPayload struct {
	CustomerID string `json:"customer_id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
}

// CartPayload is the payload for cart.* events.
type CartPayload struct {
	CartID     string `json:"cart_id"`
	CustomerID string `json:"customer_id"`
	ItemCount  int    `json:"item_count"`
	Total      int    `json:"total"`
	Currency   string `json:"currency"`
}

// PagePayload is the payload for page.* events.
type PagePayload struct {
	PageID string `json:"page_id"`
	Slug   string `json:"slug"`
	Title  string `json:"title"`
}

// PluginPayload is the payload for plugin.* events.
type PluginPayload struct {
	PluginID   string `json:"plugin_id"`
	PluginName string `json:"plugin_name"`
	Version    string `json:"version"`
}

// ThemePayload is the payload for theme.* events.
type ThemePayload struct {
	ThemeID string `json:"theme_id"`
	Name    string `json:"name"`
}

// SettingsPayload is the payload for settings.* events.
type SettingsPayload struct {
	Fields []string `json:"fields"` // which fields changed
}

// CategoryPayload is the payload for category.* events.
type CategoryPayload struct {
	CategoryID string `json:"category_id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
}

// CollectionPayload is the payload for collection.* events.
type CollectionPayload struct {
	CollectionID string `json:"collection_id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
}

// ShippingZonePayload is the payload for shipping_zone.* events.
type ShippingZonePayload struct {
	ZoneID    string   `json:"zone_id"`
	Name      string   `json:"name"`
	Countries []string `json:"countries"`
}

// ShippingRatePayload is the payload for shipping_rate.* events.
type ShippingRatePayload struct {
	RateID string `json:"rate_id"`
	ZoneID string `json:"zone_id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Price  int64  `json:"price"`
}

// TaxRatePayload is the payload for tax_rate.* events.
type TaxRatePayload struct {
	RateID  string  `json:"rate_id"`
	Name    string  `json:"name"`
	Rate    float64 `json:"rate"`
	Country string  `json:"country"`
}

// PaymentPayload is the payload for payment.* events.
type PaymentPayload struct {
	PaymentID string `json:"payment_id"`
	OrderID   string `json:"order_id"`
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
	Provider  string `json:"provider"`
	Status    string `json:"status"`
}

// RefundPayload is the payload for refund.* events.
type RefundPayload struct {
	RefundID  string `json:"refund_id"`
	PaymentID string `json:"payment_id"`
	OrderID   string `json:"order_id"`
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
	Status    string `json:"status"`
}

// CheckoutPayload is the payload for checkout.* events.
type CheckoutPayload struct {
	CartID     string `json:"cart_id"`
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
	Total      int64  `json:"total"`
	Currency   string `json:"currency"`
}

// ReviewPayload is the payload for review.* events.
type ReviewPayload struct {
	ReviewID   string `json:"review_id"`
	ProductID  string `json:"product_id"`
	CustomerID string `json:"customer_id"`
	Rating     int    `json:"rating"`
	Status     string `json:"status"`
}

// ReturnPayload is the payload for return.* events.
type ReturnPayload struct {
	ReturnID   string `json:"return_id"`
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
	Status     string `json:"status"`
	Resolution string `json:"resolution,omitempty"`
}

// BundlePayload is the payload for bundle.* events.
type BundlePayload struct {
	BundleID      string `json:"bundle_id"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	DiscountType  string `json:"discount_type"`
	DiscountValue int    `json:"discount_value"`
}

// StorefrontPayload is the payload for storefront_entry.* events.
type StorefrontPayload struct {
	StorefrontID string `json:"storefront_id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
// BlogPostPayload is the payload for blog_post.* events.
type BlogPostPayload struct {
	PostID string `json:"post_id"`
	Title  string `json:"title"`
	Slug   string `json:"slug"`
}
