package ordercontainer

import (
	"database/sql"
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/order/orderapi"
	"github.com/Abraxas-365/hada-commerce/internal/order/orderinfra"
	"github.com/Abraxas-365/hada-commerce/internal/order/ordersrv"
)

// Container wires together all order domain dependencies.
type Container struct {
	Service *ordersrv.Service
	Handler *orderapi.Handler
}

// New creates a fully-wired order container.
func New(db *sql.DB) *Container {
	repo := orderinfra.NewPostgresRepo(db)
	svc := ordersrv.New(repo)
	handler := orderapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers order HTTP routes on the given mux.
func (c *Container) RegisterRoutes(mux *http.ServeMux) {
	c.Handler.RegisterRoutes(mux)
}
