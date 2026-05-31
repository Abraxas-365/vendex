package provider

import (
	"context"
	"encoding/json"

	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/payment"
	"github.com/google/uuid"
)

// ManualProvider is a payment provider that marks payments as completed immediately.
// It is the default provider for development and testing.
type ManualProvider struct{}

// NewManualProvider creates a new ManualProvider.
func NewManualProvider() *ManualProvider {
	return &ManualProvider{}
}

// Name returns the provider identifier.
func (p *ManualProvider) Name() string {
	return "manual"
}

// Charge marks the payment as completed without contacting any external gateway.
func (p *ManualProvider) Charge(_ context.Context, amount kernel.Money, _ string, metadata map[string]string) (*payment.ProviderResult, error) {
	data, _ := json.Marshal(map[string]interface{}{
		"provider": "manual",
		"metadata": metadata,
		"amount":   amount.Amount,
		"currency": amount.Currency,
	})

	return &payment.ProviderResult{
		ProviderID: "manual_" + uuid.NewString(),
		Status:     "completed",
		Data:       json.RawMessage(data),
	}, nil
}

// Refund marks a refund as completed without contacting any external gateway.
func (p *ManualProvider) Refund(_ context.Context, providerPaymentID string, amount kernel.Money) (*payment.ProviderResult, error) {
	data, _ := json.Marshal(map[string]interface{}{
		"provider":            "manual",
		"provider_payment_id": providerPaymentID,
		"amount":              amount.Amount,
		"currency":            amount.Currency,
	})

	return &payment.ProviderResult{
		ProviderID: "manual_refund_" + uuid.NewString(),
		Status:     "completed",
		Data:       json.RawMessage(data),
	}, nil
}

// Ensure ManualProvider implements PaymentProvider.
var _ payment.PaymentProvider = (*ManualProvider)(nil)
