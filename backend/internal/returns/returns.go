package returns

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// ReturnStatus represents the lifecycle of a return request.
type ReturnStatus string

const (
	StatusRequested ReturnStatus = "requested"
	StatusApproved  ReturnStatus = "approved"
	StatusRejected  ReturnStatus = "rejected"
	StatusReceived  ReturnStatus = "received"
	StatusRefunded  ReturnStatus = "refunded"
	StatusExchanged ReturnStatus = "exchanged"
	StatusClosed    ReturnStatus = "closed"
)

// Resolution represents how the return will be resolved.
type Resolution string

const (
	ResolutionRefund      Resolution = "refund"
	ResolutionExchange    Resolution = "exchange"
	ResolutionStoreCredit Resolution = "store_credit"
)

// ItemCondition describes the condition of a returned item.
type ItemCondition string

const (
	ConditionUnopened ItemCondition = "unopened"
	ConditionLikeNew  ItemCondition = "like_new"
	ConditionUsed     ItemCondition = "used"
	ConditionDamaged  ItemCondition = "damaged"
)

// ReturnItem represents a single item in a return request.
type ReturnItem struct {
	ID        kernel.ReturnItemID `json:"id"`
	ReturnID  kernel.ReturnID     `json:"return_id"`
	TenantID  kernel.TenantID     `json:"tenant_id"`
	ProductID kernel.ProductID    `json:"product_id"`
	VariantID kernel.VariantID    `json:"variant_id,omitempty"`
	Quantity  int                 `json:"quantity"`
	Reason    string              `json:"reason,omitempty"`
	Condition ItemCondition       `json:"condition"`
	CreatedAt time.Time           `json:"created_at"`
}

// ReturnRequest is the aggregate root for a return/exchange request.
type ReturnRequest struct {
	ID             kernel.ReturnID   `json:"id"`
	TenantID       kernel.TenantID   `json:"tenant_id"`
	OrderID        kernel.OrderID    `json:"order_id"`
	CustomerID     kernel.CustomerID `json:"customer_id"`
	Status         ReturnStatus      `json:"status"`
	Reason         string            `json:"reason"`
	Notes          string            `json:"notes,omitempty"`
	AdminNotes     string            `json:"admin_notes,omitempty"`
	RefundAmount   kernel.Money      `json:"refund_amount"`
	Resolution     Resolution        `json:"resolution,omitempty"`
	TrackingNumber string            `json:"tracking_number,omitempty"`
	Carrier        string            `json:"carrier,omitempty"`
	Items          []ReturnItem      `json:"items"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

// CreateReturnInput holds all data needed to create a return request.
type CreateReturnInput struct {
	OrderID    kernel.OrderID    `json:"order_id"`
	CustomerID kernel.CustomerID `json:"customer_id"`
	Reason     string            `json:"reason"`
	Notes      string            `json:"notes,omitempty"`
	Items      []ReturnItemInput `json:"items"`
}

// ReturnItemInput represents a single item in a create return request.
type ReturnItemInput struct {
	ProductID kernel.ProductID `json:"product_id"`
	VariantID kernel.VariantID `json:"variant_id,omitempty"`
	Quantity  int              `json:"quantity"`
	Reason    string           `json:"reason,omitempty"`
	Condition ItemCondition    `json:"condition"`
}

// ApproveInput holds data needed to approve a return.
type ApproveInput struct {
	AdminNotes   string     `json:"admin_notes,omitempty"`
	Resolution   Resolution `json:"resolution"`
	RefundCents  int64      `json:"refund_amount_cents"`
	RefundCurrency string   `json:"refund_currency"`
}
