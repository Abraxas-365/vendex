package promocontainer

import (
	"database/sql"
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/promo/promoapi"
	"github.com/Abraxas-365/hada-commerce/internal/promo/promoinfra"
	"github.com/Abraxas-365/hada-commerce/internal/promo/promosrv"
)

// Container wires together the promo domain's repository, service, and handler.
type Container struct {
	handler *promoapi.Handler
}

// New builds the full promo dependency graph.
func New(db *sql.DB) *Container {
	repo := promoinfra.NewPostgresPromoRepository(db)
	svc := promosrv.New(repo)
	handler := promoapi.New(svc)
	return &Container{handler: handler}
}

// RegisterRoutes wires all promo routes onto the provided ServeMux.
func (c *Container) RegisterRoutes(mux *http.ServeMux) {
	c.handler.RegisterRoutes(mux)
}
