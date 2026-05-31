package emails

import (
	"context"
	"encoding/json"

	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/logx"
	"github.com/Abraxas-365/hada-commerce/internal/notifx"
)

// RegisterSubscriptions wires all transactional email handlers to the event bus.
// All handlers are fault-tolerant: they log errors and always return nil so that
// a failed email send never blocks other event subscribers.
func (h *Handler) RegisterSubscriptions(bus eventbus.Bus) {
	bus.Subscribe(eventbus.OrderPlaced, h.onOrderPlaced)
	bus.Subscribe(eventbus.OrderShipped, h.onOrderShipped)
	bus.Subscribe(eventbus.OrderDelivered, h.onOrderDelivered)
	bus.Subscribe(eventbus.CustomerRegistered, h.onCustomerRegistered)
	bus.Subscribe(eventbus.PaymentCompleted, h.onPaymentCompleted)
	bus.Subscribe(eventbus.RefundCompleted, h.onRefundCompleted)
}

// onOrderPlaced fires when an order is placed.
// Resolves the customer email via the EmailResolver and sends an order confirmation.
func (h *Handler) onOrderPlaced(ctx context.Context, event eventbus.Event) error {
	var payload eventbus.OrderPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		logx.Errorf("emails: failed to unmarshal order payload for order_confirmation: %v", err)
		return nil // fault-tolerant
	}

	email, err := h.resolver.ResolveEmail(ctx, string(event.TenantID), payload.CustomerID)
	if err != nil {
		logx.Errorf("emails: failed to resolve email for customer %s (order_confirmation): %v", payload.CustomerID, err)
		return nil
	}

	data := map[string]any{
		"StoreName": h.storeName,
		"OrderID":   payload.OrderID,
		"Total":     formatMoney(int64(payload.Total), payload.Currency),
		"ItemCount": payload.ItemCount,
	}

	msg := notifx.EmailMessage{
		From:    h.fromEmail,
		To:      []string{email},
		Subject: "Order Confirmation — " + payload.OrderID,
	}

	if err := h.client.SendTemplatedEmail(ctx, "order_confirmation", data, msg); err != nil {
		logx.Errorf("emails: failed to send order_confirmation to %s: %v", email, err)
		return nil // fault-tolerant
	}

	logx.Infof("emails: sent order_confirmation to %s (order_id: %s)", email, payload.OrderID)
	return nil
}

// onOrderShipped fires when an order status transitions to shipped.
// Resolves the customer email via the EmailResolver and sends a shipped notification.
func (h *Handler) onOrderShipped(ctx context.Context, event eventbus.Event) error {
	var payload eventbus.OrderPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		logx.Errorf("emails: failed to unmarshal order payload for order_shipped: %v", err)
		return nil
	}

	email, err := h.resolver.ResolveEmail(ctx, string(event.TenantID), payload.CustomerID)
	if err != nil {
		logx.Errorf("emails: failed to resolve email for customer %s (order_shipped): %v", payload.CustomerID, err)
		return nil
	}

	data := map[string]any{
		"StoreName": h.storeName,
		"OrderID":   payload.OrderID,
		"Status":    payload.Status,
	}

	msg := notifx.EmailMessage{
		From:    h.fromEmail,
		To:      []string{email},
		Subject: "Your order has shipped — " + payload.OrderID,
	}

	if err := h.client.SendTemplatedEmail(ctx, "order_shipped", data, msg); err != nil {
		logx.Errorf("emails: failed to send order_shipped to %s: %v", email, err)
		return nil // fault-tolerant
	}

	logx.Infof("emails: sent order_shipped to %s (order_id: %s)", email, payload.OrderID)
	return nil
}

// onOrderDelivered fires when an order status transitions to delivered.
// Resolves the customer email via the EmailResolver and sends a delivered notification.
func (h *Handler) onOrderDelivered(ctx context.Context, event eventbus.Event) error {
	var payload eventbus.OrderPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		logx.Errorf("emails: failed to unmarshal order payload for order_delivered: %v", err)
		return nil
	}

	email, err := h.resolver.ResolveEmail(ctx, string(event.TenantID), payload.CustomerID)
	if err != nil {
		logx.Errorf("emails: failed to resolve email for customer %s (order_delivered): %v", payload.CustomerID, err)
		return nil
	}

	data := map[string]any{
		"StoreName": h.storeName,
		"OrderID":   payload.OrderID,
	}

	msg := notifx.EmailMessage{
		From:    h.fromEmail,
		To:      []string{email},
		Subject: "Your order has been delivered — " + payload.OrderID,
	}

	if err := h.client.SendTemplatedEmail(ctx, "order_delivered", data, msg); err != nil {
		logx.Errorf("emails: failed to send order_delivered to %s: %v", email, err)
		return nil // fault-tolerant
	}

	logx.Infof("emails: sent order_delivered to %s (order_id: %s)", email, payload.OrderID)
	return nil
}

// onCustomerRegistered fires when a new customer account is created.
// CustomerPayload includes the email address, so we can send the welcome email directly.
func (h *Handler) onCustomerRegistered(ctx context.Context, event eventbus.Event) error {
	var payload eventbus.CustomerPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		logx.Errorf("emails: failed to unmarshal customer payload for customer_welcome: %v", err)
		return nil
	}

	if payload.Email == "" {
		logx.Errorf("emails: customer_welcome skipped — payload has no email (customer_id: %s)", payload.CustomerID)
		return nil
	}

	data := map[string]any{
		"StoreName": h.storeName,
		"Name":      payload.Name,
		"Email":     payload.Email,
	}

	msg := notifx.EmailMessage{
		From:    h.fromEmail,
		To:      []string{payload.Email},
		Subject: "Welcome to " + h.storeName + "!",
	}

	if err := h.client.SendTemplatedEmail(ctx, "customer_welcome", data, msg); err != nil {
		logx.Errorf("emails: failed to send customer_welcome to %s: %v", payload.Email, err)
		return nil // fault-tolerant
	}

	logx.Infof("emails: sent customer_welcome to %s (customer_id: %s)", payload.Email, payload.CustomerID)
	return nil
}

// onPaymentCompleted fires when a payment is successfully charged.
// Resolves the customer email via order lookup then EmailResolver.
func (h *Handler) onPaymentCompleted(ctx context.Context, event eventbus.Event) error {
	var payload eventbus.PaymentPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		logx.Errorf("emails: failed to unmarshal payment payload for payment_completed: %v", err)
		return nil
	}

	customerID, err := h.orderResolver.ResolveOrderCustomerID(ctx, string(event.TenantID), payload.OrderID)
	if err != nil {
		logx.Errorf("emails: failed to resolve customer for order %s (payment_completed): %v", payload.OrderID, err)
		return nil
	}

	email, err := h.resolver.ResolveEmail(ctx, string(event.TenantID), customerID)
	if err != nil {
		logx.Errorf("emails: failed to resolve email for customer %s (payment_completed): %v", customerID, err)
		return nil
	}

	data := map[string]any{
		"StoreName": h.storeName,
		"OrderID":   payload.OrderID,
		"Amount":    formatMoney(payload.Amount, payload.Currency),
		"PaymentID": payload.PaymentID,
	}

	msg := notifx.EmailMessage{
		From:    h.fromEmail,
		To:      []string{email},
		Subject: "Payment confirmed — order " + payload.OrderID,
	}

	if err := h.client.SendTemplatedEmail(ctx, "payment_completed", data, msg); err != nil {
		logx.Errorf("emails: failed to send payment_completed to %s: %v", email, err)
		return nil // fault-tolerant
	}

	logx.Infof("emails: sent payment_completed to %s (order_id: %s, payment_id: %s)", email, payload.OrderID, payload.PaymentID)
	return nil
}

// onRefundCompleted fires when a refund is successfully processed.
// Resolves the customer email via order lookup then EmailResolver.
func (h *Handler) onRefundCompleted(ctx context.Context, event eventbus.Event) error {
	var payload eventbus.RefundPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		logx.Errorf("emails: failed to unmarshal refund payload for refund_completed: %v", err)
		return nil
	}

	customerID, err := h.orderResolver.ResolveOrderCustomerID(ctx, string(event.TenantID), payload.OrderID)
	if err != nil {
		logx.Errorf("emails: failed to resolve customer for order %s (refund_completed): %v", payload.OrderID, err)
		return nil
	}

	email, err := h.resolver.ResolveEmail(ctx, string(event.TenantID), customerID)
	if err != nil {
		logx.Errorf("emails: failed to resolve email for customer %s (refund_completed): %v", customerID, err)
		return nil
	}

	data := map[string]any{
		"StoreName": h.storeName,
		"OrderID":   payload.OrderID,
		"Amount":    formatMoney(payload.Amount, payload.Currency),
		"RefundID":  payload.RefundID,
	}

	msg := notifx.EmailMessage{
		From:    h.fromEmail,
		To:      []string{email},
		Subject: "Refund processed — order " + payload.OrderID,
	}

	if err := h.client.SendTemplatedEmail(ctx, "refund_completed", data, msg); err != nil {
		logx.Errorf("emails: failed to send refund_completed to %s: %v", email, err)
		return nil // fault-tolerant
	}

	logx.Infof("emails: sent refund_completed to %s (order_id: %s, refund_id: %s)", email, payload.OrderID, payload.RefundID)
	return nil
}
