package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

// newServeMux builds the top-level HTTP mux with all routes and middleware.
func newServeMux(ctr *Container) http.Handler {
	mux := http.NewServeMux()

	// Health check — always first.
	mux.HandleFunc("GET /health", handleHealth)

	// Domain routes — each container registers under /api/v1/<domain>/...
	ctr.Product.RegisterRoutes(mux)
	ctr.Order.RegisterRoutes(mux)
	ctr.Customer.RegisterRoutes(mux)
	ctr.Catalog.RegisterRoutes(mux)
	ctr.Storefront.RegisterRoutes(mux)
	ctr.Promo.RegisterRoutes(mux)
	ctr.Media.RegisterRoutes(mux)
	ctr.Marketplace.RegisterRoutes(mux)
	ctr.Settings.RegisterRoutes(mux)
	ctr.Analytics.RegisterRoutes(mux)

	// TODO: WebSocket endpoint for agent chat
	// mux.HandleFunc("/api/v1/agent/chat", handleAgentChat)

	// Wrap with middleware chain (innermost first).
	var handler http.Handler = mux
	handler = withTenantExtraction(handler)
	handler = withRequestLogging(handler)
	handler = withCORS(handler)

	return handler
}

// newHTTPServer creates a configured *http.Server ready to start.
func newHTTPServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
		BaseContext: func(_ net.Listener) context.Context {
			return context.Background()
		},
	}
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ok"}`)
}

// ---------------------------------------------------------------------------
// Middleware
// ---------------------------------------------------------------------------

// withCORS adds permissive CORS headers suitable for development.
// Tighten AllowOrigin for production.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Tenant-ID")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// withRequestLogging logs method, path, status code, and duration for every request.
func withRequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		log.Printf("%s %s %d %s",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			time.Since(start).Round(time.Millisecond),
		)
	})
}

// withTenantExtraction reads the tenant identifier from the X-Tenant-ID header
// and injects it into the request context. Downstream handlers retrieve it with
// tenantFromContext.
func withTenantExtraction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get("X-Tenant-ID")
		// TODO: In production, extract tenant from JWT claims instead of a
		// plain header. For now, we accept the header for development.
		if tenantID != "" {
			ctx := context.WithValue(r.Context(), ctxKeyTenant, tenantID)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

// ---------------------------------------------------------------------------
// Context keys & helpers
// ---------------------------------------------------------------------------

type contextKey string

const ctxKeyTenant contextKey = "tenant_id"

// TenantFromContext returns the tenant ID injected by the tenant middleware.
// Returns empty string if no tenant is set.
func TenantFromContext(ctx context.Context) string {
	v, _ := ctx.Value(ctxKeyTenant).(string)
	return v
}

// ---------------------------------------------------------------------------
// responseWriter wraps http.ResponseWriter to capture the status code.
// ---------------------------------------------------------------------------

type responseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.statusCode = code
		rw.wroteHeader = true
	}
	rw.ResponseWriter.WriteHeader(code)
}
