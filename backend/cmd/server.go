package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// newApp builds the Fiber app with all routes registered.
// Migrated domains (product, order, customer) register directly on Fiber.
// Legacy domains (catalog, storefront, promo, media, marketplace, analytics, settings)
// are bridged from net/http via fasthttpadaptor until they are migrated.
func newApp(ctr *Container) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	// Global middleware.
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Content-Type, Authorization, X-Tenant-ID",
		MaxAge:       86400,
	}))

	// Tenant extraction middleware — reads X-Tenant-ID header and injects AuthContext.
	app.Use(tenantMiddleware)

	// Health check.
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Migrated domains — register directly on the Fiber router under /api/v1.
	v1 := app.Group("/api/v1")
	ctr.Product.RegisterRoutes(v1)
	ctr.Order.RegisterRoutes(v1)
	ctr.Customer.RegisterRoutes(v1)

	// Legacy domains — bridge from net/http until migration is complete.
	legacyMux := buildLegacyMux(ctr)
	app.Use(bridgeNetHTTP(legacyMux))

	return app
}

// buildLegacyMux creates a net/http ServeMux for domains not yet migrated to Fiber.
func buildLegacyMux(ctr *Container) *http.ServeMux {
	mux := http.NewServeMux()
	ctr.Catalog.RegisterRoutes(mux)
	ctr.Storefront.RegisterRoutes(mux)
	ctr.Promo.RegisterRoutes(mux)
	ctr.Media.RegisterRoutes(mux)
	ctr.Marketplace.RegisterRoutes(mux)
	ctr.Analytics.RegisterRoutes(mux)
	ctr.Settings.RegisterRoutes(mux)
	return mux
}

// bridgeNetHTTP wraps a net/http handler as a Fiber middleware using fasthttpadaptor.
// This allows legacy net/http handlers to coexist with Fiber handlers during migration.
func bridgeNetHTTP(h http.Handler) fiber.Handler {
	fh := fasthttpadaptor.NewFastHTTPHandler(h)
	return func(c *fiber.Ctx) error {
		fh(c.Context())
		return nil
	}
}

// tenantMiddleware reads X-Tenant-ID from the request header and injects a
// *kernel.AuthContext into Fiber locals. Handlers read it with:
//
//	authCtx := c.Locals("auth").(*kernel.AuthContext)
func tenantMiddleware(c *fiber.Ctx) error {
	tenantID := c.Get("X-Tenant-ID")
	if tenantID != "" {
		c.Locals("auth", &kernel.AuthContext{
			TenantID: kernel.TenantID(tenantID),
		})
	}
	return c.Next()
}

// errorHandler converts domain errors (errx.Error) and Fiber errors to JSON responses.
func errorHandler(c *fiber.Ctx, err error) error {
	// Default: internal server error.
	status := fiber.StatusInternalServerError
	code := "INTERNAL_ERROR"
	message := "internal error"

	if e, ok := err.(*fiber.Error); ok {
		// Fiber-native error (e.g. 404 Not Found from routing).
		status = e.Code
		message = e.Message
	} else if e, ok := err.(*errx.Error); ok {
		// Domain error from errx — use the HTTP status and message from the error.
		status = e.HTTPStatus()
		message = e.Message()
		code = string(e.Type())
	}

	return c.Status(status).JSON(fiber.Map{
		"error":   code,
		"message": message,
	})
}

// ---------------------------------------------------------------------------
// Server lifecycle — used by main.go
// ---------------------------------------------------------------------------

// runServer starts the Fiber app and handles graceful shutdown on context cancellation.
func runServer(ctx context.Context, app *fiber.App, addr string) error {
	errCh := make(chan error, 1)
	go func() {
		log.Printf("hada-commerce listening on %s", addr)
		if err := app.Listen(addr); err != nil {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		log.Println("shutting down server...")
		shutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return app.ShutdownWithContext(shutCtx)
	case err := <-errCh:
		return err
	}
}
