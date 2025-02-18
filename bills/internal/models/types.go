package models

import "time"

// Address represents a business address
type Address struct {
	Name      string `json:"name"`
	VATNumber string `json:"vat_number"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zip_code"`
	Country   string `json:"country"`
}

// Bill represents a bill entity in our system
type Bill struct {
	ID            int64                 `json:"id"`
	IssuerID      int64                 `json:"issuer_id"`
	ReceiverID    int64                 `json:"receiver_id"`
	DueDate       time.Time             `json:"due_date"`
	Currency      string                `json:"currency"`
	OriginalTotal float64               `json:"original_total"`
	EURTotal      float64               `json:"eur_total"`
	Paid          bool                  `json:"paid"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
	Issuer        *Issuer               `json:"issuer,omitempty"`
	Receiver      *Receiver             `json:"receiver,omitempty"`
	Items         []*BillItemAssignment `json:"items,omitempty"`
	// Helper fields for templates
	IssuerName   string `json:"-"`
	ReceiverName string `json:"-"`
}

// BillItem represents a service or product that can be added to bills
type BillItem struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BillItemAssignment represents the assignment of a BillItem to a Bill
type BillItemAssignment struct {
	ID             int64     `json:"id"`
	BillID         int64     `json:"bill_id"`
	BillItem       *BillItem `json:"bill_item,omitempty"`
	ItemID         int64     `json:"item_id"`
	Quantity       int       `json:"quantity"`
	Price          float64   `json:"price"`
	Currency       string    `json:"currency"`
	ExchangeRate   float64   `json:"exchange_rate"`
	OriginalAmount float64   `json:"original_amount"`
	EURAmount      float64   `json:"eur_amount"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
