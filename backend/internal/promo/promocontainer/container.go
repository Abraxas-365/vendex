package promocontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/promo/promoapi"
	"github.com/Abraxas-365/vendex/internal/promo/promoinfra"
	"github.com/Abraxas-365/vendex/internal/promo/promosrv"
)

// Container wires together the promo domain's repository, service, and handler.
type Container struct {
	Service *promosrv.Service
	Handler *promoapi.Handler
}

// New builds the full promo dependency graph.
func New(db *sqlx.DB) *Container {
	repo := promoinfra.NewPostgresPromoRepository(db)
	svc := promosrv.New(repo)
	handler := promoapi.New(svc)
	return &Container{Service: svc, Handler: handler}
}

// RegisterRoutes wires all promo routes onto the provided Fiber router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
