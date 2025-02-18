package currency

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// ExchangeRate represents a currency exchange rate
type ExchangeRate struct {
	ID        int64     `json:"id"`
	From      string    `json:"currency_from"`
	To        string    `json:"currency_to"`
	Rate      float64   `json:"rate"`
	CreatedAt time.Time `json:"created_at"`
}

// ExchangeService handles currency exchange operations
type ExchangeService struct {
	apiKey string
	db     *sql.DB
}

// NewExchangeService creates a new exchange service
func NewExchangeService(db *sql.DB) *ExchangeService {
	return &ExchangeService{
		apiKey: os.Getenv("CURRENCY_API_KEY"),
		db:     db,
	}
}

// GetRate gets the exchange rate from the API and stores it
func (s *ExchangeService) GetRate(from, to string) (*ExchangeRate, error) {
	if from == to {
		return &ExchangeRate{
			From:      from,
			To:        to,
			Rate:      1.0,
			CreatedAt: time.Now(),
		}, nil
	}

	url := fmt.Sprintf("https://api.freecurrencyapi.com/v1/latest?apikey=%s&base_currency=%s&currencies=%s",
		s.apiKey, from, to)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rate: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data map[string]float64 `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	rate := result.Data[to]
	if rate == 0 {
		return nil, fmt.Errorf("failed to get rate for %s to %s", from, to)
	}

	exchangeRate := &ExchangeRate{
		From:      from,
		To:        to,
		Rate:      rate,
		CreatedAt: time.Now(),
	}

	// Store the rate
	err = s.storeRate(exchangeRate)
	if err != nil {
		return nil, fmt.Errorf("failed to store rate: %w", err)
	}

	return exchangeRate, nil
}

// storeRate stores the exchange rate in the database
func (s *ExchangeService) storeRate(rate *ExchangeRate) error {
	result, err := s.db.Exec(`
        INSERT INTO exchange_rates (currency_from, currency_to, rate, created_at)
        VALUES (?, ?, ?, ?)
    `,
		rate.From,
		rate.To,
		rate.Rate,
		rate.CreatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	rate.ID = id

	return nil
}

// Convert converts an amount from one currency to another
func (s *ExchangeService) Convert(amount float64, from, to string) (float64, float64, error) {
	rate, err := s.GetRate(from, to)
	if err != nil {
		return 0, 0, err
	}

	return amount * rate.Rate, rate.Rate, nil
}

// SupportedCurrencies returns a list of supported currencies
func (s *ExchangeService) SupportedCurrencies() []string {
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
