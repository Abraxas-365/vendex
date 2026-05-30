package subscription

import (
	"context"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines the persistence contract for subscriptions.
type Repository interface {
	// Subscription CRUD
	Create(ctx context.Context, sub *Subscription) error
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.SubscriptionID) (*Subscription, error)
	Update(ctx context.Context, sub *Subscription) error
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.SubscriptionID) error

	// Listing
	List(ctx context.Context, tenantID kernel.TenantID, page, pageSize int) (kernel.Paginated[Subscription], error)
	ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) ([]Subscription, error)
	ListDueBilling(ctx context.Context, tenantID kernel.TenantID, before time.Time) ([]Subscription, error)

	// Billing records
	CreateBillingRecord(ctx context.Context, record *BillingRecord) error
	ListBillingRecords(ctx context.Context, tenantID kernel.TenantID, subID kernel.SubscriptionID, page, pageSize int) (kernel.Paginated[BillingRecord], error)
}
