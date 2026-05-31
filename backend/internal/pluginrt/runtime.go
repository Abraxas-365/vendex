package pluginrt

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sync"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// LoadedPlugin represents a plugin that has been loaded and initialized.
type LoadedPlugin struct {
	Name     string
	Manifest *Manifest
	Sandbox  *Sandbox
	Status   string // "loaded", "active", "error"
	Error    error  // last error if status is "error"
}

// Runtime manages the lifecycle of installed plugins.
// It handles mounting HTTP routes, serving frontend bundles, and registering agent tools.
type Runtime struct {
	mu      sync.RWMutex
	plugins map[string]*LoadedPlugin // name → loaded plugin
	mux     *http.ServeMux           // shared router for mounting plugin routes
	db      *sql.DB
}

// New creates a new plugin Runtime.
// mux is the shared HTTP router onto which plugin static routes are mounted.
// db is the shared database connection passed into each plugin's Sandbox.
func New(mux *http.ServeMux, db *sql.DB) *Runtime {
	return &Runtime{
		plugins: make(map[string]*LoadedPlugin),
		mux:     mux,
		db:      db,
	}
}

// LoadPlugin loads a plugin from its manifest and creates a sandbox for the given tenant.
// It mounts static file serving at /plugins/{name}/ and registers the plugin in memory.
// Calling LoadPlugin again for an already-loaded plugin returns an error.
func (r *Runtime) LoadPlugin(ctx context.Context, tenantID kernel.TenantID, manifest *Manifest) error {
	if manifest.Name == "" {
		return errx.Validation("loading plugin: manifest name must not be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plugins[manifest.Name]; exists {
		return errx.Conflict(fmt.Sprintf("plugin %q already loaded", manifest.Name))
	}

	sandbox := newSandbox(manifest.Name, tenantID, manifest.Permissions, r.db)

	lp := &LoadedPlugin{
		Name:     manifest.Name,
		Manifest: manifest,
		Sandbox:  sandbox,
		Status:   "loaded",
	}

	// Mount static file serving for the plugin's frontend bundle.
	// Plugins are expected to serve their assets from /plugins/{name}/.
	// NOTE: We register a no-op handler here; the actual static files are
	// served by the cmd layer wiring a file server to this pattern.
	pattern := "/plugins/" + manifest.Name + "/"
	r.mux.Handle(pattern, http.StripPrefix(pattern, http.NotFoundHandler()))

	r.plugins[manifest.Name] = lp
	return nil
}

// UnloadPlugin removes a loaded plugin by name.
// Returns an error if the plugin is not currently loaded.
func (r *Runtime) UnloadPlugin(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plugins[name]; !exists {
		return errx.NotFound(fmt.Sprintf("plugin %q not loaded", name))
	}

	delete(r.plugins, name)
	return nil
}

// GetLoaded returns info about a loaded plugin.
// Returns false if no plugin with that name is loaded.
func (r *Runtime) GetLoaded(name string) (*LoadedPlugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	lp, ok := r.plugins[name]
	return lp, ok
}

// ListLoaded returns a snapshot of all currently loaded plugins.
func (r *Runtime) ListLoaded() []*LoadedPlugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]*LoadedPlugin, 0, len(r.plugins))
	for _, lp := range r.plugins {
		out = append(out, lp)
	}
	return out
}

// GetManifest returns the manifest for a loaded plugin.
// Returns an error if no plugin with that name is loaded.
func (r *Runtime) GetManifest(name string) (*Manifest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	lp, ok := r.plugins[name]
	if !ok {
		return nil, errx.NotFound(fmt.Sprintf("plugin %q not loaded", name))
	}
	return lp.Manifest, nil
}

// ListManifests returns manifests for all loaded plugins.
// The frontend calls this to build the sidebar navigation from plugin tabs/widgets.
func (r *Runtime) ListManifests() []*Manifest {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]*Manifest, 0, len(r.plugins))
	for _, lp := range r.plugins {
		out = append(out, lp.Manifest)
	}
	return out
}
