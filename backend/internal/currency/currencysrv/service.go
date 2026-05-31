package currencysrv

import (
	"context"
	"math"
	"time"

	"github.com/Abraxas-365/vendex/internal/currency"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/google/uuid"
)

// Service handles currency domain business logic.
type Service struct {
	repo currency.Repository
}

// New creates a new currency service.
func New(repo currency.Repository) *Service {
	return &Service{repo: repo}
}

// SetRateInput holds the data needed to create or update an exchange rate.
type SetRateInput struct {
	BaseCurrency   string
	TargetCurrency string
	Rate           float64
	AutoUpdate     bool
}

// SetRate upserts an exchange rate for a tenant.
// It validates that both currencies are supported, they are different, and rate > 0.
// If a rate already exists for the pair it is updated; otherwise it is created.
func (s *Service) SetRate(ctx context.Context, tenantID kernel.TenantID, in SetRateInput) (*currency.CurrencyRate, error) {
	if _, ok := currency.SupportedCurrencies[in.BaseCurrency]; !ok {
		return nil, currency.ErrUnsupportedCurrency
	}
	if _, ok := currency.SupportedCurrencies[in.TargetCurrency]; !ok {
		return nil, currency.ErrUnsupportedCurrency
	}
	if in.BaseCurrency == in.TargetCurrency {
		return nil, currency.ErrSameCurrency
	}
	if in.Rate <= 0 {
		return nil, currency.ErrInvalidRate
	}

	now := time.Now()

	// Try to find existing rate for the pair.
	existing, err := s.repo.GetByPair(ctx, tenantID, in.BaseCurrency, in.TargetCurrency)
	if err != nil && !errx.IsNotFound(err) {
		return nil, errx.Wrap(err, "checking existing exchange rate", errx.TypeInternal)
	}

	if existing != nil {
		// Update the existing rate.
		existing.Rate = in.Rate
		existing.AutoUpdate = in.AutoUpdate
		existing.UpdatedAt = now
		if err := s.repo.Update(ctx, existing); err != nil {
			return nil, errx.Wrap(err, "updating exchange rate", errx.TypeInternal)
		}
		return existing, nil
	}

	// Create a new rate.
	rate := &currency.CurrencyRate{
		ID:             kernel.CurrencyRateID(uuid.NewString()),
		TenantID:       tenantID,
		BaseCurrency:   in.BaseCurrency,
		TargetCurrency: in.TargetCurrency,
		Rate:           in.Rate,
		AutoUpdate:     in.AutoUpdate,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repo.Create(ctx, rate); err != nil {
		return nil, errx.Wrap(err, "creating exchange rate", errx.TypeInternal)
	}

	return rate, nil
}

// GetRate retrieves a specific exchange rate by currency pair.
func (s *Service) GetRate(ctx context.Context, tenantID kernel.TenantID, base, target string) (*currency.CurrencyRate, error) {
	return s.repo.GetByPair(ctx, tenantID, base, target)
}

// Convert converts an amount from one currency to another using the stored rate.
// The math: convertedCents = int64(math.Round(float64(amount.Amount) * rate))
func (s *Service) Convert(ctx context.Context, tenantID kernel.TenantID, amount kernel.Money, targetCurrency string) (*currency.ConvertResult, error) {
	if _, ok := currency.SupportedCurrencies[targetCurrency]; !ok {
		return nil, currency.ErrUnsupportedCurrency
	}
	if amount.Currency == targetCurrency {
		return &currency.ConvertResult{
			OriginalAmount:  amount,
			ConvertedAmount: amount,
			Rate:            1.0,
		}, nil
	}

	rate, err := s.repo.GetByPair(ctx, tenantID, amount.Currency, targetCurrency)
	if err != nil {
		return nil, err
	}

	convertedCents := int64(math.Round(float64(amount.Amount) * rate.Rate))

	return &currency.ConvertResult{
		OriginalAmount:  amount,
		ConvertedAmount: kernel.Money{Amount: convertedCents, Currency: targetCurrency},
		Rate:            rate.Rate,
	}, nil
}

// ListRates returns all exchange rates for a tenant.
func (s *Service) ListRates(ctx context.Context, tenantID kernel.TenantID) ([]currency.CurrencyRate, error) {
	return s.repo.List(ctx, tenantID)
}

// DeleteRate removes an exchange rate by ID, scoped to tenant.
func (s *Service) DeleteRate(ctx context.Context, tenantID kernel.TenantID, id kernel.CurrencyRateID) error {
	return s.repo.Delete(ctx, tenantID, id)
}

// ListSupportedCurrencies returns the static map of currencies supported by the system.
func (s *Service) ListSupportedCurrencies() map[string]currency.SupportedCurrency {
	return currency.SupportedCurrencies
}
