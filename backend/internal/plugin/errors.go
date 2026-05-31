package plugin

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	ErrPluginNotFound       = errx.New("plugin not found", errx.TypeNotFound)
	ErrVersionNotFound      = errx.New("plugin version not found", errx.TypeNotFound)
	ErrInstallationNotFound = errx.New("plugin installation not found", errx.TypeNotFound)
	ErrAlreadyInstalled     = errx.New("plugin is already installed for this tenant", errx.TypeConflict)
	ErrNotInstalled         = errx.New("plugin is not installed for this tenant", errx.TypeBusiness)
	ErrInvalidInput         = errx.New("invalid plugin input", errx.TypeValidation)
)
