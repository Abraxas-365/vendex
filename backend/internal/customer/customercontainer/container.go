package customercontainer

import (
	"context"

	customerauth "github.com/Abraxas-365/hada-commerce/internal/customer/auth"
	"github.com/Abraxas-365/hada-commerce/internal/customer/customerapi"
	"github.com/Abraxas-365/hada-commerce/internal/customer/customerinfra"
	"github.com/Abraxas-365/hada-commerce/internal/customer/customersrv"
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/order"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// OrderService is the interface for order operations needed by customer auth.
type OrderService interface {
	ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, pg kernel.PaginationOptions) (kernel.Paginated[order.Order], error)
}

// Container wires together all customer domain dependencies.
type Container struct {
	Service       *customersrv.Service
	Handler       *customerapi.Handler
	AuthService   *customerauth.Service
	AuthMiddleware *customerauth.CustomerMiddleware
	AuthHandler   *customerauth.Handler
}

// New creates a fully-wired customer container.
func New(db *sqlx.DB, bus eventbus.Bus, jwtSecret string, orderSvc OrderService) *Container {
	repo := customerinfra.NewPostgresRepo(db)
	svc := customersrv.New(repo, bus)
	handler := customerapi.NewHandler(svc)

	// Auth sub-domain
	credRepo := customerauth.NewPostgresCredentialsRepo(db)
	authSvc := customerauth.NewService(credRepo, svc, jwtSecret)
	authMiddleware := customerauth.NewCustomerMiddleware(jwtSecret)
	authHandler := customerauth.NewHandler(authSvc, authMiddleware, orderSvc)

	return &Container{
		Service:        svc,
		Handler:        handler,
		AuthService:    authSvc,
		AuthMiddleware: authMiddleware,
		AuthHandler:    authHandler,
	}
}

// RegisterRoutes registers admin-facing customer HTTP routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}

// RegisterPublicRoutes registers public storefront auth routes.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.AuthHandler.RegisterPublicRoutes(router)
}

// RegisterCustomerProtectedRoutes registers customer-authenticated account routes.
func (c *Container) RegisterCustomerProtectedRoutes(router fiber.Router) {
	customerProtected := router.Group("", c.AuthMiddleware.Authenticate())
	c.AuthHandler.RegisterProtectedRoutes(customerProtected)
}
