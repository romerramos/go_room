package repository_test

import (
	"bills/db"
	"bills/internal/models"
	"bills/internal/repository"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Use a file-based database for testing
	testDB, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Enable foreign keys
	_, err = testDB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		t.Fatalf("Failed to enable foreign keys: %v", err)
	}

	// Run migrations
	if err := db.MigrateDB(testDB, "file::memory:?cache=shared"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return testDB
}

func createTestBill(t *testing.T, db *sql.DB) int64 {
	// Create test issuer
	_, err := db.Exec(`
		INSERT INTO issuers (name, vat_number, street, city, state, zip_code, country, created_at, updated_at)
		VALUES ('Test Issuer', '123456', '123 Street', 'City', 'State', '12345', 'Country', ?, ?)
	`, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to create test issuer: %v", err)
	}

	// Create test receiver
	_, err = db.Exec(`
		INSERT INTO receivers (name, vat_number, street, city, state, zip_code, country, created_at, updated_at)
		VALUES ('Test Receiver', '654321', '321 Street', 'City', 'State', '54321', 'Country', ?, ?)
	`, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to create test receiver: %v", err)
	}

	// Create test bill
	result, err := db.Exec(`
		INSERT INTO bills (amount, due_date, paid, issuer_id, receiver_id, created_at, updated_at)
		VALUES (0, ?, false, 1, 1, ?, ?)
	`, time.Now(), time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to create test bill: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get bill ID: %v", err)
	}

	return id
}

func createTestBillItem(t *testing.T, db *sql.DB) int64 {
	result, err := db.Exec(`
		INSERT INTO bill_items (description, default_price, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, "Test Item", 100.00, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to create test bill item: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get bill item ID: %v", err)
	}

	return id
}

func TestBillItemAssignmentRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewSQLiteBillItemAssignmentRepository(db)

	// Create test data
	billID := createTestBill(t, db)
	itemID := createTestBillItem(t, db)

	// Test Create
	t.Run("Create", func(t *testing.T) {
		assignment := models.NewBillItemAssignment(billID, itemID, 2, 100.00)
		err := repo.Create(assignment)
		if err != nil {
			t.Fatalf("Failed to create assignment: %v", err)
		}
		if assignment.ID == 0 {
			t.Error("Expected assignment ID to be set")
		}
		if assignment.Subtotal != 200.00 {
			t.Errorf("Expected subtotal to be 200.00, got %.2f", assignment.Subtotal)
		}
	})

	// Test GetByBillID
	t.Run("GetByBillID", func(t *testing.T) {
		assignments, err := repo.GetByBillID(billID)
		if err != nil {
			t.Fatalf("Failed to get assignments: %v", err)
		}
		if len(assignments) != 1 {
			t.Fatalf("Expected 1 assignment, got %d", len(assignments))
		}

		assignment := assignments[0]
		if assignment.BillID != billID {
			t.Errorf("Expected bill ID %d, got %d", billID, assignment.BillID)
		}
		if assignment.ItemID != itemID {
			t.Errorf("Expected item ID %d, got %d", itemID, assignment.ItemID)
		}
		if assignment.Quantity != 2 {
			t.Errorf("Expected quantity 2, got %d", assignment.Quantity)
		}
		if assignment.UnitPrice != 100.00 {
			t.Errorf("Expected unit price 100.00, got %.2f", assignment.UnitPrice)
		}
		if assignment.Subtotal != 200.00 {
			t.Errorf("Expected subtotal 200.00, got %.2f", assignment.Subtotal)
		}
		if assignment.BillItem == nil {
			t.Error("Expected bill item to be loaded")
		} else {
			if assignment.BillItem.Description != "Test Item" {
				t.Errorf("Expected item description 'Test Item', got '%s'", assignment.BillItem.Description)
			}
		}
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		assignments, err := repo.GetByBillID(billID)
		if err != nil {
			t.Fatalf("Failed to get assignments: %v", err)
		}

		assignment := assignments[0]
		assignment.Quantity = 3
		assignment.UpdateSubtotal()

		err = repo.Update(assignment)
		if err != nil {
			t.Fatalf("Failed to update assignment: %v", err)
		}

		// Verify update
		assignments, err = repo.GetByBillID(billID)
		if err != nil {
			t.Fatalf("Failed to get assignments after update: %v", err)
		}

		updated := assignments[0]
		if updated.Quantity != 3 {
			t.Errorf("Expected quantity 3, got %d", updated.Quantity)
		}
		if updated.Subtotal != 300.00 {
			t.Errorf("Expected subtotal 300.00, got %.2f", updated.Subtotal)
		}
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		assignments, err := repo.GetByBillID(billID)
		if err != nil {
			t.Fatalf("Failed to get assignments: %v", err)
		}

		err = repo.Delete(assignments[0].ID)
		if err != nil {
			t.Fatalf("Failed to delete assignment: %v", err)
		}

		// Verify deletion
		assignments, err = repo.GetByBillID(billID)
		if err != nil {
			t.Fatalf("Failed to get assignments after deletion: %v", err)
		}
		if len(assignments) != 0 {
			t.Errorf("Expected 0 assignments after deletion, got %d", len(assignments))
		}
	})
}
