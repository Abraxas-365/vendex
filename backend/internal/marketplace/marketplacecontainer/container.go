package marketplacecontainer

import (
	"database/sql"
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/marketplace/marketplaceapi"
	"github.com/Abraxas-365/hada-commerce/internal/marketplace/marketplaceinfra"
	"github.com/Abraxas-365/hada-commerce/internal/marketplace/marketplacesrv"
)

// Container wires together all marketplace domain dependencies.
type Container struct {
	Service *marketplacesrv.Service
	Handler *marketplaceapi.Handler
}

// New creates a fully-wired marketplace container.
func New(db *sql.DB) *Container {
	repo := marketplaceinfra.NewPostgresRepo(db)
	svc := marketplacesrv.New(repo)
	handler := marketplaceapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers marketplace HTTP routes on the given mux.
func (c *Container) RegisterRoutes(mux *http.ServeMux) {
	c.Handler.RegisterRoutes(mux)
}
