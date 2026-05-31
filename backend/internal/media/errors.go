package media

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	ErrMediaNotFound   = errx.New("media not found", errx.TypeNotFound)
	ErrUploadFailed    = errx.New("file upload failed", errx.TypeInternal)
	ErrInvalidFile     = errx.New("invalid or missing file", errx.TypeValidation)
	ErrFileTooLarge    = errx.New("file exceeds maximum allowed size", errx.TypeValidation)
	ErrUnsupportedType = errx.New("unsupported content type", errx.TypeValidation)
)
