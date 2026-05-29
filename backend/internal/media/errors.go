package media

import (
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
)

var (
	ErrMediaNotFound    = errx.New("MEDIA_NOT_FOUND", "media not found", http.StatusNotFound)
	ErrUploadFailed     = errx.New("MEDIA_UPLOAD_FAILED", "file upload failed", http.StatusInternalServerError)
	ErrInvalidFile      = errx.New("MEDIA_INVALID_FILE", "invalid or missing file", http.StatusBadRequest)
	ErrFileTooLarge     = errx.New("MEDIA_FILE_TOO_LARGE", "file exceeds maximum allowed size", http.StatusRequestEntityTooLarge)
	ErrUnsupportedType  = errx.New("MEDIA_UNSUPPORTED_TYPE", "unsupported content type", http.StatusUnsupportedMediaType)
)
