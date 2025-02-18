package models

import "time"

// NewBillItemAssignment creates a new BillItemAssignment instance
func NewBillItemAssignment(billID, itemID int64, quantity int, price float64, currency string, exchangeRate float64) *BillItemAssignment {
	now := time.Now()
	if currency == "" || !IsSupportedCurrency(currency) {
		currency = DefaultCurrency()
	}
	if IsDefaultCurrency(currency) {
		exchangeRate = 1.0
	}

	assignment := &BillItemAssignment{
		BillID:       billID,
		ItemID:       itemID,
		Quantity:     quantity,
		Price:        price,
		Currency:     currency,
		ExchangeRate: exchangeRate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	assignment.CalculateAmounts()
	return assignment
}

// CalculateAmounts calculates both original and EUR amounts
func (a *BillItemAssignment) CalculateAmounts() {
	a.OriginalAmount = float64(a.Quantity) * a.Price
	if IsDefaultCurrency(a.Currency) {
		a.EURAmount = a.OriginalAmount
	} else {
		a.EURAmount = a.OriginalAmount * a.ExchangeRate
	}
}
