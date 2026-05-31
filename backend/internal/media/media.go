package media

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Media represents a stored media asset (image, document, etc.).
// The actual file bytes live in a StorageProvider; this entity holds metadata.
type Media struct {
	ID          kernel.MediaID  `json:"id" db:"id"`
	TenantID    kernel.TenantID `json:"tenant_id" db:"tenant_id"`
	Filename    string          `json:"filename" db:"filename"`
	ContentType string          `json:"content_type" db:"content_type"`
	Size        int64           `json:"size" db:"size"` // bytes
	URL         string          `json:"url" db:"url"`
	Alt         string          `json:"alt" db:"alt"`
	UploadedBy  string          `json:"uploaded_by" db:"uploaded_by"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
}
