package models

import "time"

// NewBillItem creates a new BillItem instance
func NewBillItem(name string, price float64, currency string) *BillItem {
	now := time.Now()
	if currency == "" || !IsSupportedCurrency(currency) {
		currency = DefaultCurrency()
	}
	return &BillItem{
		Name:      name,
		Price:     price,
		Currency:  currency,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
