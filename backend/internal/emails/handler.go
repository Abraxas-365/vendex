// Package emails wires transactional email notifications to the event bus.
// It subscribes to domain events (order placed, shipped, customer registered, etc.)
// and sends templated emails via notifx.Client.
package emails

import (
	"fmt"

	"github.com/Abraxas-365/hada-commerce/internal/notifx"
)

// Handler handles transactional email delivery by subscribing to domain events.
type Handler struct {
	client    *notifx.Client
	fromEmail string
	storeName string
}

// New creates a new email Handler, registers all HTML templates, and returns it
// ready for subscription wiring via RegisterSubscriptions.
func New(client *notifx.Client, fromEmail, storeName string) *Handler {
	h := &Handler{
		client:    client,
		fromEmail: fromEmail,
		storeName: storeName,
	}
	h.registerTemplates()
	return h
}

// registerTemplates parses and stores all transactional email templates.
func (h *Handler) registerTemplates() {
	templates := map[string]string{
		"order_confirmation": orderConfirmationTmpl,
		"order_shipped":      orderShippedTmpl,
		"order_delivered":    orderDeliveredTmpl,
		"customer_welcome":   customerWelcomeTmpl,
		"payment_completed":  paymentCompletedTmpl,
		"refund_completed":   refundCompletedTmpl,
	}

	for name, tmpl := range templates {
		if err := h.client.RegisterTemplate(name, tmpl); err != nil {
			// Template parse errors are programming mistakes — log and panic early
			// rather than silently failing at send time.
			panic(fmt.Sprintf("emails: failed to register template %q: %v", name, err))
		}
	}
}

// formatMoney converts integer cents to a human-readable currency string.
// e.g. formatMoney(1999, "USD") → "USD 19.99"
func formatMoney(cents int64, currency string) string {
	return fmt.Sprintf("%s %.2f", currency, float64(cents)/100)
}
