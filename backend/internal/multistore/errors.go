package multistore

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	ErrNotFound          = errx.NotFound("storefront not found")
	ErrSlugConflict      = errx.Conflict("storefront slug already exists")
	ErrDomainConflict    = errx.Conflict("storefront domain already exists")
	ErrSlugRequired      = errx.Validation("slug is required")
	ErrNameRequired      = errx.Validation("name is required")
	ErrDeleteDefault     = errx.Business("cannot delete the default storefront")
	ErrCatalogNotFound   = errx.NotFound("storefront catalog entry not found")
	ErrCatalogConflict   = errx.Conflict("catalog already linked to this storefront")
)
