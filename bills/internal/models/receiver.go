package models

import "time"

// Receiver represents a business entity that can receive bills
type Receiver struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	VATNumber string    `json:"vat_number"`
	Street    string    `json:"street"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	ZipCode   string    `json:"zip_code"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewReceiver creates a new Receiver instance
func NewReceiver(name, vatNumber, street, city, state, zipCode, country string) *Receiver {
	now := time.Now()
	return &Receiver{
		Name:      name,
		VATNumber: vatNumber,
		Street:    street,
		City:      city,
		State:     state,
		ZipCode:   zipCode,
		Country:   country,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
