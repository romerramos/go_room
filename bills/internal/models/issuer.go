package models

import "time"

// Issuer represents a business entity that can issue bills
type Issuer struct {
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

// NewIssuer creates a new Issuer instance
func NewIssuer(name, vatNumber, street, city, state, zipCode, country string) *Issuer {
	now := time.Now()
	return &Issuer{
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
