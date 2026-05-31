package notificationsrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/notification"
	"github.com/google/uuid"
)

// Service implements notification business logic.
type Service struct {
	repo notification.Repository
}

// New creates a notification Service.
func New(repo notification.Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new notification for a user within a tenant.
func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, input notification.CreateNotificationInput) (*notification.Notification, error) {
	if input.Title == "" {
		return nil, notification.ErrInvalidInput
	}

	n := &notification.Notification{
		ID:           kernel.NewNotificationID(uuid.NewString()),
		TenantID:     tenantID,
		UserID:       input.UserID,
		Title:        input.Title,
		Body:         input.Body,
		Type:         input.Type,
		ResourceType: input.ResourceType,
		ResourceID:   input.ResourceID,
		Read:         false,
		CreatedAt:    time.Now().UTC(),
	}

	if n.Type == "" {
		n.Type = notification.TypeInfo
	}

	if err := s.repo.Create(ctx, n); err != nil {
		return nil, err
	}

	return n, nil
}

// List returns paginated notifications for a user within a tenant.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, userID kernel.UserID, unreadOnly bool, page, pageSize int) (kernel.Paginated[notification.Notification], error) {
	return s.repo.List(ctx, tenantID, userID, unreadOnly, page, pageSize)
}

// MarkRead marks a single notification as read.
func (s *Service) MarkRead(ctx context.Context, tenantID kernel.TenantID, id kernel.NotificationID) error {
	_, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	return s.repo.MarkRead(ctx, tenantID, id)
}

// MarkAllRead marks all notifications for a user as read.
func (s *Service) MarkAllRead(ctx context.Context, tenantID kernel.TenantID, userID kernel.UserID) error {
	return s.repo.MarkAllRead(ctx, tenantID, userID)
}

// GetUnreadCount returns the number of unread notifications for a user.
func (s *Service) GetUnreadCount(ctx context.Context, tenantID kernel.TenantID, userID kernel.UserID) (int, error) {
	return s.repo.GetUnreadCount(ctx, tenantID, userID)
}

// Delete removes a notification scoped to the tenant.
func (s *Service) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.NotificationID) error {
	_, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, tenantID, id)
}
