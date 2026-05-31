package pluginrt

import (
	"database/sql"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Sandbox provides a plugin with scoped access to platform resources.
// Each installed plugin gets its own Sandbox instance.
type Sandbox struct {
	PluginName  string
	TenantID    kernel.TenantID
	Permissions []string // from manifest
	db          *sql.DB  // shared DB connection
}

// newSandbox creates a new Sandbox for a plugin installation.
func newSandbox(pluginName string, tenantID kernel.TenantID, permissions []string, db *sql.DB) *Sandbox {
	return &Sandbox{
		PluginName:  pluginName,
		TenantID:    tenantID,
		Permissions: permissions,
		db:          db,
	}
}

// DB returns the shared database connection.
// Plugins use table prefix "plugin_{name}_" for their own tables.
func (s *Sandbox) DB() *sql.DB {
	return s.db
}

// TablePrefix returns the plugin's table name prefix.
// Plugin-owned tables must be named with this prefix to avoid collisions.
func (s *Sandbox) TablePrefix() string {
	return "plugin_" + s.PluginName + "_"
}

// HasPermission checks if the plugin declared a given permission in its manifest.
func (s *Sandbox) HasPermission(perm string) bool {
	for _, p := range s.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}
