package notification

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Repository defines the persistence contract for the notification domain.
type Repository interface {
	// Create persists a new notification.
	Create(ctx context.Context, n *Notification) error

	// GetByID returns a notification scoped to the tenant.
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.NotificationID) (*Notification, error)

	// List returns paginated notifications for a user within a tenant.
	// If unreadOnly is true, only unread notifications are returned.
	List(ctx context.Context, tenantID kernel.TenantID, userID kernel.UserID, unreadOnly bool, page, pageSize int) (kernel.Paginated[Notification], error)

	// MarkRead marks a single notification as read.
	MarkRead(ctx context.Context, tenantID kernel.TenantID, id kernel.NotificationID) error

	// MarkAllRead marks all notifications for the user as read.
	MarkAllRead(ctx context.Context, tenantID kernel.TenantID, userID kernel.UserID) error

	// GetUnreadCount returns the count of unread notifications for the user.
	GetUnreadCount(ctx context.Context, tenantID kernel.TenantID, userID kernel.UserID) (int, error)

	// Delete removes a notification.
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.NotificationID) error
}
