package models

import (
	"testing"
	"time"
)

func TestNewBill(t *testing.T) {
	dueDate := time.Now()
	bill := NewBill(dueDate, 1, 2)

	if bill.Amount != 0 {
		t.Errorf("Expected initial amount to be 0, got %.2f", bill.Amount)
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
}

func TestCalculateAmount(t *testing.T) {
	tests := []struct {
		name  string
		items []*BillItemAssignment
		want  float64
	}{
		{
			name: "Single item",
			items: []*BillItemAssignment{
				NewBillItemAssignment(1, 1, 2, 100.00),
			},
			want: 200.00,
		},
		{
			name: "Multiple items",
			items: []*BillItemAssignment{
				NewBillItemAssignment(1, 1, 2, 100.00),
				NewBillItemAssignment(1, 2, 1, 50.00),
			},
			want: 250.00,
		},
		{
			name:  "No items",
			items: []*BillItemAssignment{},
			want:  0.00,
		},
		{
			name: "Decimal prices",
			items: []*BillItemAssignment{
				NewBillItemAssignment(1, 1, 3, 19.99),
				NewBillItemAssignment(1, 2, 2, 24.50),
			},
			want: 108.97,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bill := NewBill(time.Now(), 1, 1)
			bill.Items = tt.items
			bill.CalculateAmount()

			if bill.Amount != tt.want {
				t.Errorf("Amount = %.2f, want %.2f", bill.Amount, tt.want)
			}

			// Verify each item's subtotal is calculated correctly
			for _, item := range bill.Items {
				expectedSubtotal := float64(item.Quantity) * item.UnitPrice
				if item.Subtotal != expectedSubtotal {
					t.Errorf("Item subtotal = %.2f, want %.2f", item.Subtotal, expectedSubtotal)
				}
			}
		})
	}
}
