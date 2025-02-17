package models

import "time"

// BillItemAssignment represents the assignment of a BillItem to a Bill
type BillItemAssignment struct {
	ID        int64     `json:"id"`
	BillID    int64     `json:"bill_id"`
	BillItem  *BillItem `json:"bill_item,omitempty"`
	ItemID    int64     `json:"item_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
	Subtotal  float64   `json:"subtotal"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewBillItemAssignment creates a new BillItemAssignment instance
func NewBillItemAssignment(billID, itemID int64, quantity int, unitPrice float64) *BillItemAssignment {
	now := time.Now()
	return &BillItemAssignment{
		BillID:    billID,
		ItemID:    itemID,
		Quantity:  quantity,
		UnitPrice: unitPrice,
		Subtotal:  float64(quantity) * unitPrice,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateSubtotal recalculates the subtotal based on quantity and unit price
func (a *BillItemAssignment) UpdateSubtotal() {
	a.Subtotal = float64(a.Quantity) * a.UnitPrice
}
