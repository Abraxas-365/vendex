package media

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Media represents a stored media asset (image, document, etc.).
// The actual file bytes live in a StorageProvider; this entity holds metadata.
type Media struct {
	ID          kernel.MediaID  `json:"id"`
	TenantID    kernel.TenantID `json:"tenant_id"`
	Filename    string          `json:"filename"`
	ContentType string          `json:"content_type"`
	Size        int64           `json:"size"` // bytes
	URL         string          `json:"url"`
	Alt         string          `json:"alt"`
	UploadedBy  string          `json:"uploaded_by"`
	CreatedAt   time.Time       `json:"created_at"`
}
