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
