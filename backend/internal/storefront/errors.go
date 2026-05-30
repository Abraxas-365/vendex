package storefront

import (
	"github.com/Abraxas-365/hada-commerce/internal/errx"
)

var (
	ErrPageNotFound      = errx.New("page not found", errx.TypeNotFound)
	ErrSlugAlreadyExists = errx.New("a page with this slug already exists", errx.TypeConflict)
	ErrPageNotPublished  = errx.New("page is not published", errx.TypeNotFound)
	ErrPageArchived      = errx.New("archived pages cannot be edited", errx.TypeBusiness)
	ErrInvalidStatus     = errx.New("invalid page status transition", errx.TypeBusiness)
	ErrVersionNotFound   = errx.New("page version not found", errx.TypeNotFound)

	ErrBlockTypeNotFound = errx.New("block type not found", errx.TypeNotFound)
	ErrBlockTypeExists   = errx.New("block type with this name already exists", errx.TypeConflict)
)
