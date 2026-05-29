package mediacontainer

import (
	"database/sql"
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/media"
	"github.com/Abraxas-365/hada-commerce/internal/media/mediaapi"
	"github.com/Abraxas-365/hada-commerce/internal/media/mediainfra"
	"github.com/Abraxas-365/hada-commerce/internal/media/mediasrv"
)

// Container wires together the media domain's repository, storage, service, and handler.
type Container struct {
	handler *mediaapi.Handler
}

// New builds the full media dependency graph using the provided storage backend.
func New(db *sql.DB, storage media.StorageProvider) *Container {
	repo := mediainfra.NewPostgresMediaRepository(db)
	svc := mediasrv.New(repo, storage)
	handler := mediaapi.New(svc)
	return &Container{handler: handler}
}

// NewWithLocalStorage is a convenience constructor that wires up a local filesystem
// StorageProvider rooted at baseDir, serving files at baseURL.
func NewWithLocalStorage(db *sql.DB, baseDir, baseURL string) (*Container, error) {
	storage, err := mediainfra.NewLocalStorageProvider(baseDir, baseURL)
	if err != nil {
		return nil, err
	}
	return New(db, storage), nil
}

// RegisterRoutes wires all media routes onto the provided ServeMux.
func (c *Container) RegisterRoutes(mux *http.ServeMux) {
	c.handler.RegisterRoutes(mux)
}
