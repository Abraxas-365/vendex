package notificationcontainer

import (
	"context"
	"encoding/json"

	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/logx"
	"github.com/Abraxas-365/vendex/internal/notification"
	"github.com/Abraxas-365/vendex/internal/notification/notificationapi"
	"github.com/Abraxas-365/vendex/internal/notification/notificationinfra"
	"github.com/Abraxas-365/vendex/internal/notification/notificationsrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all notification domain dependencies.
type Container struct {
	Service *notificationsrv.Service
	Handler *notificationapi.Handler
}

// New creates a fully-wired notification container and subscribes to relevant domain events.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := notificationinfra.NewPostgresRepo(db)
	svc := notificationsrv.New(repo)
	handler := notificationapi.NewHandler(svc)

	c := &Container{
		Service: svc,
		Handler: handler,
	}

	// Subscribe to domain events and create in-app notifications.
	// Notifications are broadcast to a sentinel user "system" for now —
	// real implementations should fan-out to all admin users of the tenant.
	// Passing an empty UserID here is intentional: consumers can query by tenantID alone.

	bus.Subscribe(eventbus.OrderPlaced, func(ctx context.Context, event eventbus.Event) error {
		var p eventbus.OrderPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			logx.Errorf("notification: unmarshal order payload: %v", err)
			return nil
		}
		_, err := svc.Create(ctx, event.TenantID, notification.CreateNotificationInput{
			UserID:       kernel.UserID(""),
			Title:        "New order received",
			Body:         "Order " + p.OrderID + " has been placed",
			Type:         notification.TypeInfo,
			ResourceType: "order",
			ResourceID:   p.OrderID,
		})
		if err != nil {
			logx.Errorf("notification: create for order.placed: %v", err)
		}
		return nil
	})

	bus.Subscribe(eventbus.ReturnRequested, func(ctx context.Context, event eventbus.Event) error {
		var p eventbus.ReturnPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			logx.Errorf("notification: unmarshal return payload: %v", err)
			return nil
		}
		_, err := svc.Create(ctx, event.TenantID, notification.CreateNotificationInput{
			UserID:       kernel.UserID(""),
			Title:        "New return request",
			Body:         "Return " + p.ReturnID + " has been requested for order " + p.OrderID,
			Type:         notification.TypeWarning,
			ResourceType: "return",
			ResourceID:   p.ReturnID,
		})
		if err != nil {
			logx.Errorf("notification: create for return.requested: %v", err)
		}
		return nil
	})

	bus.Subscribe(eventbus.ReviewCreated, func(ctx context.Context, event eventbus.Event) error {
		var p eventbus.ReviewPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			logx.Errorf("notification: unmarshal review payload: %v", err)
			return nil
		}
		_, err := svc.Create(ctx, event.TenantID, notification.CreateNotificationInput{
			UserID:       kernel.UserID(""),
			Title:        "New product review",
			Body:         "A review has been submitted for product " + p.ProductID,
			Type:         notification.TypeInfo,
			ResourceType: "review",
			ResourceID:   p.ReviewID,
		})
		if err != nil {
			logx.Errorf("notification: create for review.created: %v", err)
		}
		return nil
	})

	bus.Subscribe(eventbus.StockLowAlert, func(ctx context.Context, event eventbus.Event) error {
		_, err := svc.Create(ctx, event.TenantID, notification.CreateNotificationInput{
			UserID:       kernel.UserID(""),
			Title:        "Low stock alert",
			Body:         "One or more products are running low on stock",
			Type:         notification.TypeWarning,
			ResourceType: "inventory",
			ResourceID:   "",
		})
		if err != nil {
			logx.Errorf("notification: create for stock.low_alert: %v", err)
		}
		return nil
	})

	return c
}

// RegisterRoutes registers protected notification routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
