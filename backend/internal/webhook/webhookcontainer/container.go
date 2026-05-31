package webhookcontainer

import (
	"context"
	"encoding/json"

	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/logx"
	"github.com/Abraxas-365/vendex/internal/webhook/webhookapi"
	"github.com/Abraxas-365/vendex/internal/webhook/webhookinfra"
	"github.com/Abraxas-365/vendex/internal/webhook/webhooksrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all webhook domain dependencies.
type Container struct {
	Service *webhooksrv.Service
	Handler *webhookapi.Handler
}

// New creates a fully-wired webhook container and subscribes to all events on the bus.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := webhookinfra.NewPostgresRepo(db)
	svc := webhooksrv.New(repo)
	handler := webhookapi.NewHandler(svc)

	c := &Container{
		Service: svc,
		Handler: handler,
	}

	// Subscribe to all domain events and fan-out to registered webhooks.
	bus.SubscribeAll(func(ctx context.Context, event eventbus.Event) error {
		payload, err := json.Marshal(event)
		if err != nil {
			logx.Errorf("webhook container: failed to marshal event %s: %v", event.Type, err)
			return nil // non-fatal
		}

		if err := svc.Deliver(ctx, event.TenantID, string(event.Type), json.RawMessage(payload)); err != nil {
			logx.Errorf("webhook container: deliver event %s to tenant %s: %v", event.Type, event.TenantID, err)
		}
		return nil
	})

	return c
}

// RegisterRoutes registers protected webhook routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
