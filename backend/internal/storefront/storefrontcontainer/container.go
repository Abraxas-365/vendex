package storefrontcontainer

import (
	"database/sql"
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/storefront/storefrontapi"
	"github.com/Abraxas-365/hada-commerce/internal/storefront/storefrontinfra"
	"github.com/Abraxas-365/hada-commerce/internal/storefront/storefrontsrv"
)

// Container wires together the storefront domain's repository, service, and handler.
type Container struct {
	handler *storefrontapi.Handler
}

// New builds the full storefront dependency graph.
func New(db *sql.DB) *Container {
	pageRepo := storefrontinfra.NewPostgresPageRepository(db)
	versionRepo := storefrontinfra.NewPostgresPageVersionRepository(db)
	svc := storefrontsrv.New(pageRepo, versionRepo)
	handler := storefrontapi.New(svc)
	return &Container{handler: handler}
}

// RegisterRoutes wires all storefront routes onto the provided ServeMux.
func (c *Container) RegisterRoutes(mux *http.ServeMux) {
	c.handler.RegisterRoutes(mux)
}
