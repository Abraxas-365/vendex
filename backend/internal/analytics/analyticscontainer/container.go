package analyticscontainer

import (
	"database/sql"
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/analytics/analyticsapi"
	"github.com/Abraxas-365/hada-commerce/internal/analytics/analyticsinfra"
	"github.com/Abraxas-365/hada-commerce/internal/analytics/analyticssrv"
)

// Container wires together all analytics domain dependencies.
type Container struct {
	Service *analyticssrv.Service
	Handler *analyticsapi.Handler
}

// New creates a fully-wired analytics container.
func New(db *sql.DB) *Container {
	repo := analyticsinfra.NewPostgresRepo(db)
	svc := analyticssrv.New(repo)
	handler := analyticsapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers analytics HTTP routes on the given mux.
func (c *Container) RegisterRoutes(mux *http.ServeMux) {
	c.Handler.RegisterRoutes(mux)
}
