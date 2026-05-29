package main

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/Abraxas-365/hada-commerce/internal/analytics/analyticscontainer"
	"github.com/Abraxas-365/hada-commerce/internal/catalog/catalogcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/config"
	"github.com/Abraxas-365/hada-commerce/internal/customer/customercontainer"
	"github.com/Abraxas-365/hada-commerce/internal/marketplace/marketplacecontainer"
	"github.com/Abraxas-365/hada-commerce/internal/media/mediacontainer"
	"github.com/Abraxas-365/hada-commerce/internal/order/ordercontainer"
	"github.com/Abraxas-365/hada-commerce/internal/product/productcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/promo/promocontainer"
	"github.com/Abraxas-365/hada-commerce/internal/settings/settingscontainer"
	"github.com/Abraxas-365/hada-commerce/internal/storefront/storefrontcontainer"
)

// Container is the top-level DI container that wires all domain containers together.
// Each domain container owns its own repository, service, and handler graph.
type Container struct {
	Product    *productcontainer.Container
	Order      *ordercontainer.Container
	Customer   *customercontainer.Container
	Catalog    *catalogcontainer.Container
	Storefront *storefrontcontainer.Container
	Promo      *promocontainer.Container
	Media       *mediacontainer.Container
	Marketplace *marketplacecontainer.Container
	Analytics   *analyticscontainer.Container
	Settings    *settingscontainer.Container
}

// NewContainer builds all domain containers in dependency order.
// Domains that are independent of each other are built first; domains with
// cross-domain dependencies (if any are added later) come after.
func NewContainer(db *sql.DB, cfg *config.Config) (*Container, error) {
	// Independent domain containers — order doesn't matter today,
	// but we group them logically: core commerce, then CMS.
	product := productcontainer.New(db)
	order := ordercontainer.New(db)
	customer := customercontainer.New(db)
	catalog := catalogcontainer.New(db)

	// CMS domains
	storefront := storefrontcontainer.New(db)
	promo := promocontainer.New(db)
	marketplace := marketplacecontainer.New(db)

	// Analytics domain — read-only, queries across tables.
	analyticsCtr := analyticscontainer.New(db)

	// Settings domain — per-tenant store configuration.
	settingsCtr := settingscontainer.New(db)

	// Media needs a storage provider; choose based on config.
	mediaCtr, err := newMediaContainer(db, cfg)
	if err != nil {
		return nil, fmt.Errorf("media container: %w", err)
	}

	return &Container{
		Product:    product,
		Order:      order,
		Customer:   customer,
		Catalog:    catalog,
		Storefront: storefront,
		Promo:      promo,
		Media:       mediaCtr,
		Marketplace: marketplace,
		Analytics:   analyticsCtr,
		Settings:    settingsCtr,
	}, nil
}

// newMediaContainer selects the right storage backend based on cfg.MediaStorage.
func newMediaContainer(db *sql.DB, cfg *config.Config) (*mediacontainer.Container, error) {
	switch cfg.MediaStorage {
	case "local", "":
		baseDir := filepath.Join(".", "uploads")
		baseURL := "/uploads"
		return mediacontainer.NewWithLocalStorage(db, baseDir, baseURL)
	// TODO: add "s3" case when S3 storage provider is implemented
	default:
		return nil, fmt.Errorf("unsupported media storage type: %q", cfg.MediaStorage)
	}
}
