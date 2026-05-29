package settings

import "github.com/Abraxas-365/hada-commerce/internal/kernel/errx"

// ErrNotFound is returned when no settings row exists for the requested tenant.
var ErrNotFound = errx.New("SETTINGS_NOT_FOUND", "store settings not found", 404)
