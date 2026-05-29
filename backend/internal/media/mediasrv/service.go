package mediasrv

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/media"
)

// Service implements media business logic.
type Service struct {
	repo    media.MediaRepository
	storage media.StorageProvider
}

// New creates a new media Service.
func New(repo media.MediaRepository, storage media.StorageProvider) *Service {
	return &Service{repo: repo, storage: storage}
}

// UploadInput holds everything needed to store a media file.
type UploadInput struct {
	TenantID    kernel.TenantID
	Filename    string
	ContentType string
	Size        int64
	Alt         string
	UploadedBy  string
	Data        io.Reader
}

// Upload stores the file via the StorageProvider and persists metadata.
func (s *Service) Upload(ctx context.Context, input UploadInput) (*media.Media, error) {
	id := generateUUID()
	// Build a storage key that is unique and tenant-scoped.
	key := fmt.Sprintf("media/%s/%s/%s", string(input.TenantID), id, sanitizeFilename(input.Filename))

	url, err := s.storage.Upload(ctx, key, input.ContentType, input.Data)
	if err != nil {
		return nil, fmt.Errorf("upload file: %w", err)
	}

	m := &media.Media{
		ID:          kernel.MediaID(id),
		TenantID:    input.TenantID,
		Filename:    input.Filename,
		ContentType: input.ContentType,
		Size:        input.Size,
		URL:         url,
		Alt:         input.Alt,
		UploadedBy:  input.UploadedBy,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, m); err != nil {
		// Best-effort cleanup of orphaned file.
		_ = s.storage.Delete(ctx, key)
		return nil, fmt.Errorf("persist media metadata: %w", err)
	}
	return m, nil
}

// Delete removes the media metadata record. Callers should handle storage
// cleanup separately if the URL / storage key is needed (the URL is returned).
func (s *Service) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.MediaID) error {
	m, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	// Derive the storage key from the URL — the local provider stores relative keys,
	// so we use the path component of the URL as the key.
	key := path.Base(m.URL)
	if err := s.repo.Delete(ctx, tenantID, id); err != nil {
		return fmt.Errorf("delete media metadata: %w", err)
	}
	// Best-effort; failure here does not roll back the metadata deletion.
	_ = s.storage.Delete(ctx, key)
	return nil
}

// GetByID retrieves a media record by ID.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.MediaID) (*media.Media, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// List returns paginated media records for a tenant.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[media.Media], error) {
	return s.repo.List(ctx, tenantID, p)
}

// generateUUID produces a random UUID v4 string.
func generateUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// sanitizeFilename strips directory components from an uploaded filename.
func sanitizeFilename(name string) string {
	return path.Base(name)
}
