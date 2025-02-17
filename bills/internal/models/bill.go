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
	ID         int64                 `json:"id"`
	Amount     float64               `json:"amount"`
	DueDate    time.Time             `json:"due_date"`
	Paid       bool                  `json:"paid"`
	IssuerID   int64                 `json:"issuer_id"`
	Issuer     *Issuer               `json:"issuer,omitempty"`
	ReceiverID int64                 `json:"receiver_id"`
	Receiver   *Receiver             `json:"receiver,omitempty"`
	Items      []*BillItemAssignment `json:"items,omitempty"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
}

// NewBill creates a new Bill instance with default values
func NewBill(dueDate time.Time, issuerID, receiverID int64) *Bill {
	now := time.Now()
	return &Bill{
		Amount:     0, // Will be calculated from items
		DueDate:    dueDate,
		Paid:       false,
		IssuerID:   issuerID,
		ReceiverID: receiverID,
		Items:      make([]*BillItemAssignment, 0),
		Issuer:     &Issuer{},
		Receiver:   &Receiver{},
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// CalculateAmount calculates the total amount from bill items
func (b *Bill) CalculateAmount() {
	var total float64
	for _, item := range b.Items {
		item.Subtotal = float64(item.Quantity) * item.UnitPrice
		total += item.Subtotal
	}
	b.Amount = total
}
