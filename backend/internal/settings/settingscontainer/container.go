package settingscontainer

import (
	"database/sql"
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/settings/settingsapi"
	"github.com/Abraxas-365/hada-commerce/internal/settings/settingsinfra"
	"github.com/Abraxas-365/hada-commerce/internal/settings/settingssrv"
)

// Container wires together all settings domain dependencies.
type Container struct {
	Service *settingssrv.Service
	Handler *settingsapi.Handler
}

// New creates a fully-wired settings container.
func New(db *sql.DB) *Container {
	repo := settingsinfra.NewPostgresRepo(db)
	svc := settingssrv.New(repo)
	handler := settingsapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers settings HTTP routes on the given mux.
func (c *Container) RegisterRoutes(mux *http.ServeMux) {
	c.Handler.RegisterRoutes(mux)
}
