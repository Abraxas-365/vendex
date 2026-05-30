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
// NOTE: OrderPayload only carries CustomerID, not the customer email. To send
// this email to the customer we would need to resolve the email from the customer
// service. For now we log the intent and skip the actual send.
// TODO: inject a customer lookup service and resolve email before sending.
func (h *Handler) onOrderPlaced(ctx context.Context, event eventbus.Event) error {
	var payload eventbus.OrderPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		logx.Errorf("emails: failed to unmarshal order payload for order_confirmation: %v", err)
		return nil // fault-tolerant
	}

	logx.WithFields(logx.Fields{
		"order_id":    payload.OrderID,
		"customer_id": payload.CustomerID,
		"total":       payload.Total,
		"item_count":  payload.ItemCount,
	}).Info("emails: would send order_confirmation (customer email not in payload — TODO: look up via customer service)")

	return nil
}

// onOrderShipped fires when an order status transitions to shipped.
// NOTE: Same limitation as onOrderPlaced — customer email not in OrderPayload.
// TODO: inject a customer lookup service and resolve email before sending.
func (h *Handler) onOrderShipped(ctx context.Context, event eventbus.Event) error {
	var payload eventbus.OrderPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		logx.Errorf("emails: failed to unmarshal order payload for order_shipped: %v", err)
		return nil
	}

	logx.WithFields(logx.Fields{
		"order_id":    payload.OrderID,
		"customer_id": payload.CustomerID,
		"status":      payload.Status,
	}).Info("emails: would send order_shipped (customer email not in payload — TODO: look up via customer service)")

	return nil
}

// onOrderDelivered fires when an order status transitions to delivered.
// NOTE: Same limitation as onOrderPlaced — customer email not in OrderPayload.
// TODO: inject a customer lookup service and resolve email before sending.
func (h *Handler) onOrderDelivered(ctx context.Context, event eventbus.Event) error {
	var payload eventbus.OrderPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		logx.Errorf("emails: failed to unmarshal order payload for order_delivered: %v", err)
		return nil
	}

	logx.WithFields(logx.Fields{
		"order_id":    payload.OrderID,
		"customer_id": payload.CustomerID,
	}).Info("emails: would send order_delivered (customer email not in payload — TODO: look up via customer service)")

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
// NOTE: PaymentPayload carries an OrderID but not the customer email.
// TODO: inject a customer lookup service and resolve email before sending.
func (h *Handler) onPaymentCompleted(ctx context.Context, event eventbus.Event) error {
	var payload eventbus.PaymentPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		logx.Errorf("emails: failed to unmarshal payment payload for payment_completed: %v", err)
		return nil
	}

	logx.WithFields(logx.Fields{
		"payment_id": payload.PaymentID,
		"order_id":   payload.OrderID,
		"amount":     formatMoney(payload.Amount, payload.Currency),
	}).Info("emails: would send payment_completed (customer email not in payload — TODO: look up via customer service)")

	return nil
}

// onRefundCompleted fires when a refund is successfully processed.
// NOTE: RefundPayload carries an OrderID but not the customer email.
// TODO: inject a customer lookup service and resolve email before sending.
func (h *Handler) onRefundCompleted(ctx context.Context, event eventbus.Event) error {
	var payload eventbus.RefundPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		logx.Errorf("emails: failed to unmarshal refund payload for refund_completed: %v", err)
		return nil
	}

	logx.WithFields(logx.Fields{
		"refund_id": payload.RefundID,
		"order_id":  payload.OrderID,
		"amount":    formatMoney(payload.Amount, payload.Currency),
	}).Info("emails: would send refund_completed (customer email not in payload — TODO: look up via customer service)")

	return nil
}
