package settings

import "github.com/Abraxas-365/hada-commerce/internal/errx"

// ErrNotFound is returned when no settings row exists for the requested tenant.
var ErrNotFound = errx.New("store settings not found", errx.TypeNotFound)
