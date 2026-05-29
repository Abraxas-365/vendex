package kernel

// Money represents a monetary amount in the smallest currency unit (cents).
type Money struct {
	Amount   int64  // in cents
	Currency string // ISO 4217, e.g. "USD", "EUR"
}

func NewMoney(cents int64, currency string) Money {
	return Money{Amount: cents, Currency: currency}
}

func (m Money) Add(other Money) Money {
	return Money{Amount: m.Amount + other.Amount, Currency: m.Currency}
}

func (m Money) Multiply(qty int) Money {
	return Money{Amount: m.Amount * int64(qty), Currency: m.Currency}
}
