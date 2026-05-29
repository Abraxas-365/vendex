package storefront

import (
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
)

var (
	ErrPageNotFound      = errx.New("PAGE_NOT_FOUND", "page not found", http.StatusNotFound)
	ErrSlugAlreadyExists = errx.New("SLUG_ALREADY_EXISTS", "a page with this slug already exists", http.StatusConflict)
	ErrPageNotPublished  = errx.New("PAGE_NOT_PUBLISHED", "page is not published", http.StatusNotFound)
	ErrPageArchived      = errx.New("PAGE_ARCHIVED", "archived pages cannot be edited", http.StatusUnprocessableEntity)
	ErrInvalidStatus     = errx.New("INVALID_STATUS", "invalid page status transition", http.StatusUnprocessableEntity)
	ErrVersionNotFound   = errx.New("VERSION_NOT_FOUND", "page version not found", http.StatusNotFound)
)
