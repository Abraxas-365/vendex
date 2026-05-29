package mediacontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/media"
	"github.com/Abraxas-365/hada-commerce/internal/media/mediaapi"
	"github.com/Abraxas-365/hada-commerce/internal/media/mediainfra"
	"github.com/Abraxas-365/hada-commerce/internal/media/mediasrv"
)

// Container wires together all media domain dependencies.
type Container struct {
	Service *mediasrv.Service
	Handler *mediaapi.Handler
}

// New creates a fully-wired media container with a local filesystem storage provider.
// uploadDir is the directory where uploaded files will be stored.
// baseURL is the public URL prefix for serving uploaded files (e.g. "http://localhost:3000/uploads").
func New(db *sqlx.DB, uploadDir string, baseURL string) (*Container, error) {
	storage, err := mediainfra.NewLocalStorageProvider(uploadDir, baseURL)
	if err != nil {
		return nil, err
	}
	repo := mediainfra.NewPostgresMediaRepository(db)
	svc := mediasrv.New(repo, storage)
	handler := mediaapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}, nil
}

// NewWithStorage creates a media container with a custom StorageProvider.
// Useful for production deployments (e.g., S3-backed storage).
func NewWithStorage(db *sqlx.DB, storage media.StorageProvider) *Container {
	repo := mediainfra.NewPostgresMediaRepository(db)
	svc := mediasrv.New(repo, storage)
	handler := mediaapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers media HTTP routes on the given Fiber router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
