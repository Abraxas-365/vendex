package theme

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	// ErrThemeNotFound is returned when a theme cannot be found.
	ErrThemeNotFound = errx.New("theme not found", errx.TypeNotFound)

	// ErrThemeConflict is returned when a theme with the same name already exists.
	ErrThemeConflict = errx.New("theme with this name already exists", errx.TypeConflict)

	// ErrThemeActiveDelete is returned when attempting to delete the active theme.
	ErrThemeActiveDelete = errx.New("cannot delete the active theme", errx.TypeBusiness)

	// ErrThemeNoActive is returned when no active theme exists.
	ErrThemeNoActive = errx.New("no active theme found", errx.TypeNotFound)
)
