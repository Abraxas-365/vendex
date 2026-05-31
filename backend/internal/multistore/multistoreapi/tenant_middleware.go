package multistoreapi

import (
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/multistore"
	"github.com/gofiber/fiber/v2"
)

// TenantMiddleware resolves the tenant from the request and stores it in Locals.
// It extracts tenant from: X-Tenant-ID header, subdomain, or custom domain.
type TenantMiddleware struct {
	resolver *multistore.TenantResolver
}

// NewTenantMiddleware creates a new tenant resolution middleware.
func NewTenantMiddleware(resolver *multistore.TenantResolver) *TenantMiddleware {
	return &TenantMiddleware{resolver: resolver}
}

// Resolve returns a Fiber handler that resolves and sets the tenant ID.
func (m *TenantMiddleware) Resolve() fiber.Handler {
	return func(c *fiber.Ctx) error {
		host := c.Hostname()
		headerTenant := c.Get("X-Tenant-ID")

		tenantID, err := m.resolver.Resolve(c.Context(), host, headerTenant)
		if err != nil {
			return err
		}

		// Store resolved tenant ID in locals for downstream handlers
		c.Locals("tenant_id", tenantID)
		// Also set it as a header so existing handlers that read X-Tenant-ID still work
		c.Request().Header.Set("X-Tenant-ID", string(tenantID))
		return c.Next()
	}
}

// TenantFromLocals extracts the resolved tenant ID from Fiber context.
func TenantFromLocals(c *fiber.Ctx) kernel.TenantID {
	if id, ok := c.Locals("tenant_id").(kernel.TenantID); ok {
		return id
	}
	return kernel.TenantID(c.Get("X-Tenant-ID"))
}
