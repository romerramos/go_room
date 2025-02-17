package models

import "time"

// BillItem represents a service or product that can be added to bills
type BillItem struct {
	ID           int64     `json:"id"`
	Description  string    `json:"description"`
	DefaultPrice float64   `json:"default_price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewBillItem creates a new BillItem instance
func NewBillItem(description string, defaultPrice float64) *BillItem {
	now := time.Now()
	return &BillItem{
		Description:  description,
		DefaultPrice: defaultPrice,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
