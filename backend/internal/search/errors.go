package search

import "github.com/Abraxas-365/hada-commerce/internal/errx"

var (
	ErrInvalidQuery = errx.New("invalid search query", errx.TypeValidation)
	ErrSearchFailed = errx.New("search operation failed", errx.TypeInternal)
)
