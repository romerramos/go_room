package models

import "testing"

func TestNewBillItemAssignment(t *testing.T) {
	tests := []struct {
		name         string
		billID       int64
		itemID       int64
		quantity     int
		price        float64
		currency     string
		exchangeRate float64
		wantOriginal float64
		wantEUR      float64
	}{
		{
			name:         "Simple calculation EUR",
			billID:       1,
			itemID:       1,
			quantity:     2,
			price:        100.00,
			currency:     "EUR",
			exchangeRate: 1.0,
			wantOriginal: 200.00,
			wantEUR:      200.00,
		},
		{
			name:         "USD to EUR conversion",
			billID:       1,
			itemID:       1,
			quantity:     2,
			price:        100.00,
			currency:     "USD",
			exchangeRate: 0.85,
			wantOriginal: 200.00,
			wantEUR:      170.00,
		},
		{
			name:         "Zero quantity",
			billID:       1,
			itemID:       1,
			quantity:     0,
			price:        100.00,
			currency:     "EUR",
			exchangeRate: 1.0,
			wantOriginal: 0.00,
			wantEUR:      0.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assignment := NewBillItemAssignment(tt.billID, tt.itemID, tt.quantity, tt.price, tt.currency, tt.exchangeRate)

			if assignment.BillID != tt.billID {
				t.Errorf("BillID = %v, want %v", assignment.BillID, tt.billID)
			}
			if assignment.ItemID != tt.itemID {
				t.Errorf("ItemID = %v, want %v", assignment.ItemID, tt.itemID)
			}
			if assignment.Quantity != tt.quantity {
				t.Errorf("Quantity = %v, want %v", assignment.Quantity, tt.quantity)
			}
			if assignment.Price != tt.price {
				t.Errorf("Price = %v, want %v", assignment.Price, tt.price)
			}
			if assignment.Currency != tt.currency {
				t.Errorf("Currency = %v, want %v", assignment.Currency, tt.currency)
			}
			if assignment.ExchangeRate != tt.exchangeRate {
				t.Errorf("ExchangeRate = %v, want %v", assignment.ExchangeRate, tt.exchangeRate)
			}
			if assignment.OriginalAmount != tt.wantOriginal {
				t.Errorf("OriginalAmount = %v, want %v", assignment.OriginalAmount, tt.wantOriginal)
			}
			if assignment.EURAmount != tt.wantEUR {
				t.Errorf("EURAmount = %v, want %v", assignment.EURAmount, tt.wantEUR)
			}
		})
	}
}

func TestCalculateAmounts(t *testing.T) {
	tests := []struct {
		name         string
		quantity     int
		price        float64
		currency     string
		exchangeRate float64
		wantOriginal float64
		wantEUR      float64
	}{
		{
			name:         "EUR calculation",
			quantity:     5,
			price:        10.00,
			currency:     "EUR",
			exchangeRate: 1.0,
			wantOriginal: 50.00,
			wantEUR:      50.00,
		},
		{
			name:         "USD calculation",
			quantity:     2,
			price:        25.50,
			currency:     "USD",
			exchangeRate: 0.85,
			wantOriginal: 51.00,
			wantEUR:      43.35,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assignment := &BillItemAssignment{
				Quantity:     tt.quantity,
				Price:        tt.price,
				Currency:     tt.currency,
				ExchangeRate: tt.exchangeRate,
			}
			assignment.CalculateAmounts()

			if assignment.OriginalAmount != tt.wantOriginal {
				t.Errorf("OriginalAmount = %v, want %v", assignment.OriginalAmount, tt.wantOriginal)
			}
			if assignment.EURAmount != tt.wantEUR {
				t.Errorf("EURAmount = %v, want %v", assignment.EURAmount, tt.wantEUR)
			}
		})
	}
}
