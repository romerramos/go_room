package models

import (
	"testing"
	"time"
)

func TestNewBill(t *testing.T) {
	dueDate := time.Now()
	bill := NewBill(dueDate, 1, 2)

	if bill.OriginalTotal != 0 {
		t.Errorf("Expected initial original total to be 0, got %.2f", bill.OriginalTotal)
	}
	if bill.EURTotal != 0 {
		t.Errorf("Expected initial EUR total to be 0, got %.2f", bill.EURTotal)
	}
	if bill.DueDate != dueDate {
		t.Errorf("Expected due date to be %v, got %v", dueDate, bill.DueDate)
	}
	if bill.Paid {
		t.Error("Expected new bill to be unpaid")
	}
	if bill.IssuerID != 1 {
		t.Errorf("Expected issuer ID to be 1, got %d", bill.IssuerID)
	}
	if bill.ReceiverID != 2 {
		t.Errorf("Expected receiver ID to be 2, got %d", bill.ReceiverID)
	}
	if len(bill.Items) != 0 {
		t.Errorf("Expected items to be empty, got %d items", len(bill.Items))
	}
	if bill.Currency != DefaultCurrency() {
		t.Errorf("Expected currency to be %s, got %s", DefaultCurrency(), bill.Currency)
	}
}

func TestCalculateTotals(t *testing.T) {
	tests := []struct {
		name         string
		items        []*BillItemAssignment
		wantOriginal float64
		wantEUR      float64
	}{
		{
			name: "Single EUR item",
			items: []*BillItemAssignment{
				NewBillItemAssignment(1, 1, 2, 100.00, "EUR", 1.0),
			},
			wantOriginal: 200.00,
			wantEUR:      200.00,
		},
		{
			name: "Multiple items with different currencies",
			items: []*BillItemAssignment{
				NewBillItemAssignment(1, 1, 2, 100.00, "EUR", 1.0),
				NewBillItemAssignment(1, 2, 1, 50.00, "USD", 0.85),
			},
			wantOriginal: 250.00,
			wantEUR:      242.50,
		},
		{
			name:         "No items",
			items:        []*BillItemAssignment{},
			wantOriginal: 0.00,
			wantEUR:      0.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bill := NewBill(time.Now(), 1, 1)
			bill.Items = tt.items
			bill.CalculateTotals()

			if bill.OriginalTotal != tt.wantOriginal {
				t.Errorf("OriginalTotal = %.2f, want %.2f", bill.OriginalTotal, tt.wantOriginal)
			}
			if bill.EURTotal != tt.wantEUR {
				t.Errorf("EURTotal = %.2f, want %.2f", bill.EURTotal, tt.wantEUR)
			}
		})
	}
}
