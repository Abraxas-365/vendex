package marketplace

import (
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
)

var (
	ErrPluginNotFound      = errx.New("MARKETPLACE_PLUGIN_NOT_FOUND", "plugin not found", http.StatusNotFound)
	ErrVersionNotFound     = errx.New("MARKETPLACE_VERSION_NOT_FOUND", "plugin version not found", http.StatusNotFound)
	ErrAlreadyInstalled    = errx.New("MARKETPLACE_ALREADY_INSTALLED", "plugin is already installed", http.StatusConflict)
	ErrNotInstalled        = errx.New("MARKETPLACE_NOT_INSTALLED", "plugin is not installed", http.StatusNotFound)
	ErrIncompatibleVersion = errx.New("MARKETPLACE_INCOMPATIBLE_VERSION", "plugin version is incompatible with this platform version", http.StatusBadRequest)
	ErrPluginNameTaken     = errx.New("MARKETPLACE_PLUGIN_NAME_TAKEN", "a plugin with this name already exists", http.StatusConflict)
)
