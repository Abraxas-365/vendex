package pluginrt

import (
	"encoding/json"
	"net/http"
)

// Handler exposes plugin runtime info to the frontend via HTTP.
type Handler struct {
	runtime *Runtime
}

// NewHandler creates a new Handler backed by the given Runtime.
func NewHandler(rt *Runtime) *Handler {
	return &Handler{runtime: rt}
}

// RegisterRoutes adds plugin runtime routes to the mux:
//
//	GET /plugins/manifests          — returns all loaded plugin manifests (frontend uses this to build tabs)
//	GET /plugins/{name}/manifest    — returns a specific plugin's manifest
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /plugins/manifests", h.listManifests)
	mux.HandleFunc("GET /plugins/{name}/manifest", h.getManifest)
}

// listManifests returns all loaded plugin manifests as a JSON array.
func (h *Handler) listManifests(w http.ResponseWriter, r *http.Request) {
	manifests := h.runtime.ListManifests()
	writeJSON(w, http.StatusOK, manifests)
}

// getManifest returns the manifest for a single plugin by name.
func (h *Handler) getManifest(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "plugin name is required", http.StatusBadRequest)
		return
	}

	m, err := h.runtime.GetManifest(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, m)
}

// writeJSON encodes v as JSON and writes it to w with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// Header already written; nothing we can do but log would go here.
		_ = err
	}
}
