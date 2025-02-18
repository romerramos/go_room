package models

import "time"

// NewBill creates a new Bill instance with default values
func NewBill(dueDate time.Time, issuerID, receiverID int64) *Bill {
	now := time.Now()
	return &Bill{
		DueDate:       dueDate,
		IssuerID:      issuerID,
		ReceiverID:    receiverID,
		Currency:      DefaultCurrency(),
		OriginalTotal: 0,
		EURTotal:      0,
		Paid:          false,
		Items:         make([]*BillItemAssignment, 0),
		Issuer:        &Issuer{},
		Receiver:      &Receiver{},
		IssuerName:    "",
		ReceiverName:  "",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// CalculateTotals calculates both original and EUR totals for a bill
func (b *Bill) CalculateTotals() {
	var originalTotal, eurTotal float64
	for _, item := range b.Items {
		if item.Currency == b.Currency {
			// If item currency matches bill currency, add to original total directly
			originalTotal += item.OriginalAmount
		} else {
			// If item currency is different, convert to bill currency
			// For now, we'll use the EUR amount since we don't have direct conversion rates
			if b.Currency == DefaultCurrency() {
				originalTotal += item.EURAmount
			} else {
				// If bill currency is not EUR, we should convert from EUR to bill currency
				// TODO: Use proper exchange rate service
				originalTotal += item.EURAmount
			}
		}
		eurTotal += item.EURAmount
	}
	b.OriginalTotal = originalTotal
	b.EURTotal = eurTotal
}
