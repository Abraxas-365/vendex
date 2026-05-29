package productcontainer

import (
	"database/sql"
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/product/productapi"
	"github.com/Abraxas-365/hada-commerce/internal/product/productinfra"
	"github.com/Abraxas-365/hada-commerce/internal/product/productsrv"
)

// Container wires together all product domain dependencies.
type Container struct {
	Service *productsrv.Service
	Handler *productapi.Handler
}

// New creates a fully-wired product container.
func New(db *sql.DB) *Container {
	repo := productinfra.NewPostgresRepo(db)
	svc := productsrv.New(repo)
	handler := productapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers product HTTP routes on the given mux.
func (c *Container) RegisterRoutes(mux *http.ServeMux) {
	c.Handler.RegisterRoutes(mux)
}
