package mediainfra

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/Abraxas-365/vendex/internal/errx"
)

// LocalStorageProvider implements media.StorageProvider using the local filesystem.
// Files are stored under BaseDir and served at BaseURL/<key>.
// This is suitable for development; use an S3 provider in production.
type LocalStorageProvider struct {
	// BaseDir is the root directory where files are written.
	BaseDir string
	// BaseURL is the public URL prefix (e.g. "http://localhost:8080/uploads").
	BaseURL string
}

// NewLocalStorageProvider creates a new LocalStorageProvider.
// It creates BaseDir if it does not exist.
func NewLocalStorageProvider(baseDir, baseURL string) (*LocalStorageProvider, error) {
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, errx.Wrap(err, "create storage dir", errx.TypeExternal)
	}
	return &LocalStorageProvider{BaseDir: baseDir, BaseURL: baseURL}, nil
}

// Upload writes the file to the local filesystem and returns its public URL.
func (p *LocalStorageProvider) Upload(_ context.Context, key string, _ string, r io.Reader) (string, error) {
	dest := filepath.Join(p.BaseDir, filepath.FromSlash(key))

	// Ensure parent directories exist.
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return "", errx.Wrap(err, "create upload dir", errx.TypeExternal)
	}

	f, err := os.Create(dest)
	if err != nil {
		return "", errx.Wrap(err, "create file", errx.TypeExternal)
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return "", errx.Wrap(err, "write file", errx.TypeExternal)
	}

	return p.BaseURL + "/" + key, nil
}

// Delete removes a file from the local filesystem.
func (p *LocalStorageProvider) Delete(_ context.Context, key string) error {
	dest := filepath.Join(p.BaseDir, filepath.FromSlash(key))
	if err := os.Remove(dest); err != nil && !os.IsNotExist(err) {
		return errx.Wrap(err, "delete file", errx.TypeExternal)
	}
	return nil
}

// GetURL returns the public URL for a stored key without re-reading the file.
func (p *LocalStorageProvider) GetURL(_ context.Context, key string) (string, error) {
	return p.BaseURL + "/" + key, nil
}
