package giftcardcontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/giftcard/giftcardapi"
	"github.com/Abraxas-365/vendex/internal/giftcard/giftcardinfra"
	"github.com/Abraxas-365/vendex/internal/giftcard/giftcardsrv"
)

// Container wires together the gift card domain's repository, service, and handler.
type Container struct {
	Service *giftcardsrv.Service
	Handler *giftcardapi.Handler
}

// New builds the full gift card dependency graph.
func New(db *sqlx.DB) *Container {
	repo := giftcardinfra.NewPostgresRepository(db)
	svc := giftcardsrv.New(repo)
	handler := giftcardapi.New(svc)
	return &Container{Service: svc, Handler: handler}
}

// RegisterRoutes wires all protected gift card routes onto the provided Fiber router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}

// RegisterPublicRoutes wires all public gift card routes onto the provided Fiber router.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.Handler.RegisterPublicRoutes(router)
}
