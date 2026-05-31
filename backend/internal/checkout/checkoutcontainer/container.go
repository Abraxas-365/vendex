package checkoutcontainer

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/checkout/checkoutapi"
	"github.com/Abraxas-365/vendex/internal/checkout/checkoutsrv"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/order"
	"github.com/Abraxas-365/vendex/internal/order/orderinfra"
	"github.com/Abraxas-365/vendex/internal/order/ordersrv"
	"github.com/Abraxas-365/vendex/internal/payment"
	"github.com/Abraxas-365/vendex/internal/payment/paymentsrv"
	"github.com/Abraxas-365/vendex/internal/promo/promosrv"
	"github.com/Abraxas-365/vendex/internal/shipping/shippingsrv"
	"github.com/Abraxas-365/vendex/internal/tax/taxsrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all checkout domain dependencies.
type Container struct {
	Service *checkoutsrv.Service
	Handler *checkoutapi.Handler
}

// New creates a fully-wired checkout container.
func New(
	db *sqlx.DB,
	bus eventbus.Bus,
	cartSvc checkoutsrv.CartGetter,
	orderSvc *ordersrv.Service,
	shippingSvc *shippingsrv.Service,
	taxSvc *taxsrv.Service,
	paymentSvc *paymentsrv.Service,
	promoSvc *promosrv.Service,
) *Container {
	// Create a fresh order repo for checkout field updates.
	orderRepo := orderinfra.NewPostgresRepo(db)

	svc := checkoutsrv.New(
		cartSvc,
		orderSvcAdapter{orderSvc},
		orderRepo,
		shippingSvc,
		taxSvc,
		paymentSvcAdapter{paymentSvc},
		promoSvcAdapter{promoSvc},
		bus,
	)
	handler := checkoutapi.NewHandler(svc)

	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterPublicRoutes registers public checkout HTTP routes.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.Handler.RegisterPublicRoutes(router)
}

// ---------------------------------------------------------------------------
// Adapters — bridge concrete service types to checkoutsrv interfaces
// ---------------------------------------------------------------------------

// orderSvcAdapter adapts *ordersrv.Service to checkoutsrv.OrderCreator.
type orderSvcAdapter struct{ svc *ordersrv.Service }

func (a orderSvcAdapter) Create(ctx context.Context, tenantID kernel.TenantID, in ordersrv.CreateInput) (*order.Order, error) {
	return a.svc.Create(ctx, tenantID, in)
}

func (a orderSvcAdapter) UpdateStatus(ctx context.Context, tenantID kernel.TenantID, id kernel.OrderID, status order.OrderStatus) (*order.Order, error) {
	return a.svc.UpdateStatus(ctx, tenantID, id, status)
}

// paymentSvcAdapter adapts *paymentsrv.Service to checkoutsrv.PaymentProcessor.
type paymentSvcAdapter struct{ svc *paymentsrv.Service }

func (a paymentSvcAdapter) CreatePayment(ctx context.Context, tenantID kernel.TenantID, orderID kernel.OrderID, amount int64, currency string, providerName string, method string) (*payment.Payment, error) {
	return a.svc.CreatePayment(ctx, tenantID, orderID, amount, currency, providerName, method)
}

func (a paymentSvcAdapter) ProcessPayment(ctx context.Context, tenantID kernel.TenantID, paymentID kernel.PaymentID, token string) (*payment.Payment, error) {
	return a.svc.ProcessPayment(ctx, tenantID, paymentID, token)
}

// promoSvcAdapter adapts *promosrv.Service to checkoutsrv.PromoApplier.
type promoSvcAdapter struct{ svc *promosrv.Service }

func (a promoSvcAdapter) Validate(ctx context.Context, tenantID kernel.TenantID, code string, orderTotalCents int64) (checkoutsrv.ValidationResult, error) {
	r, err := a.svc.Validate(ctx, tenantID, code, orderTotalCents)
	if err != nil {
		return checkoutsrv.ValidationResult{}, err
	}
	return checkoutsrv.ValidationResult{
		Valid:          r.Valid,
		DiscountCents:  r.DiscountCents,
		IsFreeShipping: r.IsFreeShipping,
		Reason:         r.Reason,
	}, nil
}

func (a promoSvcAdapter) Apply(ctx context.Context, tenantID kernel.TenantID, code string, orderTotalCents int64) (int64, error) {
	return a.svc.Apply(ctx, tenantID, code, orderTotalCents)
}
