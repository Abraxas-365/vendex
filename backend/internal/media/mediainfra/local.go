package mediainfra

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
		return nil, fmt.Errorf("create storage dir: %w", err)
	}
	return &LocalStorageProvider{BaseDir: baseDir, BaseURL: baseURL}, nil
}

// Upload writes the file to the local filesystem and returns its public URL.
func (p *LocalStorageProvider) Upload(_ context.Context, key string, _ string, r io.Reader) (string, error) {
	dest := filepath.Join(p.BaseDir, filepath.FromSlash(key))

	// Ensure parent directories exist.
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return "", fmt.Errorf("create upload dir: %w", err)
	}

	f, err := os.Create(dest)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	return p.BaseURL + "/" + key, nil
}

// Delete removes a file from the local filesystem.
func (p *LocalStorageProvider) Delete(_ context.Context, key string) error {
	dest := filepath.Join(p.BaseDir, filepath.FromSlash(key))
	if err := os.Remove(dest); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete file: %w", err)
	}
	return nil
}

// GetURL returns the public URL for a stored key without re-reading the file.
func (p *LocalStorageProvider) GetURL(_ context.Context, key string) (string, error) {
	return p.BaseURL + "/" + key, nil
}
