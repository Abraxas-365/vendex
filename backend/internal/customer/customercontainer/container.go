package customercontainer

import (
	"database/sql"
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/customer/customerapi"
	"github.com/Abraxas-365/hada-commerce/internal/customer/customerinfra"
	"github.com/Abraxas-365/hada-commerce/internal/customer/customersrv"
)

// Container wires together all customer domain dependencies.
type Container struct {
	Service *customersrv.Service
	Handler *customerapi.Handler
}

// New creates a fully-wired customer container.
func New(db *sql.DB) *Container {
	repo := customerinfra.NewPostgresRepo(db)
	svc := customersrv.New(repo)
	handler := customerapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers customer HTTP routes on the given mux.
func (c *Container) RegisterRoutes(mux *http.ServeMux) {
	c.Handler.RegisterRoutes(mux)
}
