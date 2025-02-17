package models

import "testing"

func TestNewBillItemAssignment(t *testing.T) {
	tests := []struct {
		name      string
		billID    int64
		itemID    int64
		quantity  int
		unitPrice float64
		want      float64 // expected subtotal
	}{
		{
			name:      "Simple calculation",
			billID:    1,
			itemID:    1,
			quantity:  2,
			unitPrice: 100.00,
			want:      200.00,
		},
		{
			name:      "Zero quantity",
			billID:    1,
			itemID:    1,
			quantity:  0,
			unitPrice: 100.00,
			want:      0.00,
		},
		{
			name:      "Zero price",
			billID:    1,
			itemID:    1,
			quantity:  2,
			unitPrice: 0.00,
			want:      0.00,
		},
		{
			name:      "Decimal price",
			billID:    1,
			itemID:    1,
			quantity:  3,
			unitPrice: 19.99,
			want:      59.97,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assignment := NewBillItemAssignment(tt.billID, tt.itemID, tt.quantity, tt.unitPrice)

			if assignment.BillID != tt.billID {
				t.Errorf("BillID = %v, want %v", assignment.BillID, tt.billID)
			}
			if assignment.ItemID != tt.itemID {
				t.Errorf("ItemID = %v, want %v", assignment.ItemID, tt.itemID)
			}
			if assignment.Quantity != tt.quantity {
				t.Errorf("Quantity = %v, want %v", assignment.Quantity, tt.quantity)
			}
			if assignment.UnitPrice != tt.unitPrice {
				t.Errorf("UnitPrice = %v, want %v", assignment.UnitPrice, tt.unitPrice)
			}
			if assignment.Subtotal != tt.want {
				t.Errorf("Subtotal = %v, want %v", assignment.Subtotal, tt.want)
			}
		})
	}
}

func TestUpdateSubtotal(t *testing.T) {
	tests := []struct {
		name      string
		quantity  int
		unitPrice float64
		want      float64
	}{
		{
			name:      "Update after quantity change",
			quantity:  5,
			unitPrice: 10.00,
			want:      50.00,
		},
		{
			name:      "Update after price change",
			quantity:  2,
			unitPrice: 25.50,
			want:      51.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assignment := &BillItemAssignment{
				Quantity:  tt.quantity,
				UnitPrice: tt.unitPrice,
			}
			assignment.UpdateSubtotal()

			if assignment.Subtotal != tt.want {
				t.Errorf("Subtotal = %v, want %v", assignment.Subtotal, tt.want)
			}
		})
	}
}
