package media

import (
	"context"
	"io"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// MediaRepository defines persistence operations for Media metadata.
// All operations are scoped by TenantID.
type MediaRepository interface {
	// Create persists new media metadata.
	Create(ctx context.Context, m *Media) error
	// GetByID retrieves a media record by primary key.
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.MediaID) (*Media, error)
	// Delete removes a media metadata record.
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.MediaID) error
	// List returns all media records for a tenant with pagination.
	List(ctx context.Context, tenantID kernel.TenantID, p kernel.Pagination) (kernel.PaginatedResult[Media], error)
}

// StorageProvider abstracts file storage backends (local filesystem, S3, etc.).
// The service calls Upload to persist bytes and Delete to remove them.
type StorageProvider interface {
	// Upload stores the file and returns the public URL.
	Upload(ctx context.Context, key string, contentType string, r io.Reader) (url string, err error)
	// Delete removes a file from storage by its key.
	Delete(ctx context.Context, key string) error
	// GetURL returns the public URL for a storage key without re-uploading.
	GetURL(ctx context.Context, key string) (string, error)
}
