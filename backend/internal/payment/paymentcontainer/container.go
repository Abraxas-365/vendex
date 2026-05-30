package paymentcontainer

import (
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/payment"
	"github.com/Abraxas-365/hada-commerce/internal/payment/paymentapi"
	"github.com/Abraxas-365/hada-commerce/internal/payment/paymentinfra"
	"github.com/Abraxas-365/hada-commerce/internal/payment/paymentsrv"
	"github.com/Abraxas-365/hada-commerce/internal/payment/provider"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all payment domain dependencies.
type Container struct {
	Service *paymentsrv.Service
	Handler *paymentapi.Handler
}

// New creates a fully-wired payment container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := paymentinfra.NewPostgresRepo(db)
	providers := map[string]payment.PaymentProvider{
		"manual": provider.NewManualProvider(),
	}
	svc := paymentsrv.New(repo, bus, providers)
	handler := paymentapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers payment HTTP routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
