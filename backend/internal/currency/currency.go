package currency

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// CurrencyRate represents an exchange rate between two currencies for a tenant.
type CurrencyRate struct {
	ID             kernel.CurrencyRateID `json:"id" db:"id"`
	TenantID       kernel.TenantID       `json:"tenant_id" db:"tenant_id"`
	BaseCurrency   string                `json:"base_currency" db:"base_currency"`
	TargetCurrency string                `json:"target_currency" db:"target_currency"`
	Rate           float64               `json:"rate" db:"rate"`
	AutoUpdate     bool                  `json:"auto_update" db:"auto_update"`
	UpdatedAt      time.Time             `json:"updated_at" db:"updated_at"`
	CreatedAt      time.Time             `json:"created_at" db:"created_at"`
}

// ConvertResult holds the result of a currency conversion operation.
type ConvertResult struct {
	OriginalAmount  kernel.Money `json:"original_amount"`
	ConvertedAmount kernel.Money `json:"converted_amount"`
	Rate            float64      `json:"rate"`
}

// SupportedCurrency describes a currency that this system supports.
type SupportedCurrency struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	DecimalPlaces int    `json:"decimal_places"`
}

// SupportedCurrencies is the static list of currencies supported by the system.
var SupportedCurrencies = map[string]SupportedCurrency{
	"USD": {Code: "USD", Name: "US Dollar", Symbol: "$", DecimalPlaces: 2},
	"EUR": {Code: "EUR", Name: "Euro", Symbol: "€", DecimalPlaces: 2},
	"GBP": {Code: "GBP", Name: "British Pound", Symbol: "£", DecimalPlaces: 2},
	"JPY": {Code: "JPY", Name: "Japanese Yen", Symbol: "¥", DecimalPlaces: 0},
	"CAD": {Code: "CAD", Name: "Canadian Dollar", Symbol: "CA$", DecimalPlaces: 2},
	"AUD": {Code: "AUD", Name: "Australian Dollar", Symbol: "A$", DecimalPlaces: 2},
	"BRL": {Code: "BRL", Name: "Brazilian Real", Symbol: "R$", DecimalPlaces: 2},
	"MXN": {Code: "MXN", Name: "Mexican Peso", Symbol: "MX$", DecimalPlaces: 2},
	"COP": {Code: "COP", Name: "Colombian Peso", Symbol: "COL$", DecimalPlaces: 2},
}
