package marketplacecontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/marketplace/marketplaceapi"
	"github.com/Abraxas-365/vendex/internal/marketplace/marketplaceinfra"
	"github.com/Abraxas-365/vendex/internal/marketplace/marketplacesrv"
)

// Container wires together all marketplace domain dependencies.
type Container struct {
	Service       *marketplacesrv.VendorService
	PresetService *marketplacesrv.PresetService
	Handler       *marketplaceapi.Handler
	PresetHandler *marketplaceapi.PresetHandler
}

// New creates a fully-wired marketplace container.
func New(db *sqlx.DB) *Container {
	vendorRepo  := marketplaceinfra.NewPostgresVendorRepo(db)
	productRepo := marketplaceinfra.NewPostgresVendorProductRepo(db)
	orderRepo   := marketplaceinfra.NewPostgresVendorOrderRepo(db)
	svc     := marketplacesrv.New(vendorRepo, productRepo, orderRepo)
	handler := marketplaceapi.NewHandler(svc)

	presetRepo  := marketplaceinfra.NewPostgresPresetRepo(db)
	installRepo := marketplaceinfra.NewPostgresPresetInstallRepo(db)
	presetSvc     := marketplacesrv.NewPresetService(presetRepo, installRepo)
	presetHandler := marketplaceapi.NewPresetHandler(presetSvc)

	return &Container{
		Service:       svc,
		PresetService: presetSvc,
		Handler:       handler,
		PresetHandler: presetHandler,
	}
}

// RegisterRoutes registers marketplace HTTP routes on the given Fiber router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
	c.PresetHandler.RegisterPresetRoutes(router)
}
