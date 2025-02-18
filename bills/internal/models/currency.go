package models

// SupportedCurrencies returns a list of supported currencies
func SupportedCurrencies() []string {
	return []string{
		"EUR", // Base currency
		"USD", // US Dollar
		"CAD", // Canadian Dollar
		"GBP", // British Pound
		"AUD", // Australian Dollar
		"JPY", // Japanese Yen
		"CHF", // Swiss Franc
		"CNY", // Chinese Yuan
		"NZD", // New Zealand Dollar
		"MXN", // Mexican Peso
	}
}

// IsSupportedCurrency checks if a currency is supported
func IsSupportedCurrency(currency string) bool {
	for _, c := range SupportedCurrencies() {
		if c == currency {
			return true
		}
	}
	return false
}

// DefaultCurrency returns the default currency (EUR)
func DefaultCurrency() string {
	return "EUR"
}

// IsDefaultCurrency checks if a currency is the default currency (EUR)
func IsDefaultCurrency(currency string) bool {
	return currency == DefaultCurrency()
}
