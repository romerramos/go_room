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

func setupBillTestDB(t *testing.T) *sql.DB {
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

func createTestData(t *testing.T, db *sql.DB) (int64, int64, int64) {
	// Create test issuer
	result, err := db.Exec(`
		INSERT INTO issuers (name, vat_number, street, city, state, zip_code, country, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, "Test Issuer", "123456", "123 Street", "City", "State", "12345", "Country", time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to create test issuer: %v", err)
	}
	issuerID, _ := result.LastInsertId()

	// Create test receiver
	result, err = db.Exec(`
		INSERT INTO receivers (name, vat_number, street, city, state, zip_code, country, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, "Test Receiver", "654321", "321 Street", "City", "State", "54321", "Country", time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to create test receiver: %v", err)
	}
	receiverID, _ := result.LastInsertId()

	// Create test bill item
	result, err = db.Exec(`
		INSERT INTO bill_items (description, default_price, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, "Test Item", 100.00, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to create test bill item: %v", err)
	}
	itemID, _ := result.LastInsertId()

	return issuerID, receiverID, itemID
}

func TestBillRepository(t *testing.T) {
	db := setupBillTestDB(t)
	defer db.Close()

	repo := repository.NewSQLiteBillRepository(db)
	issuerID, receiverID, itemID := createTestData(t, db)

	// Test Create
	t.Run("Create", func(t *testing.T) {
		bill := models.NewBill(time.Now(), issuerID, receiverID)
		err := repo.Create(bill)
		if err != nil {
			t.Fatalf("Failed to create bill: %v", err)
		}
		if bill.ID == 0 {
			t.Error("Expected bill ID to be set")
		}
	})

	// Test GetByID with items
	t.Run("GetByID with items", func(t *testing.T) {
		// Create a bill with an item
		bill := models.NewBill(time.Now(), issuerID, receiverID)
		err := repo.Create(bill)
		if err != nil {
			t.Fatalf("Failed to create bill: %v", err)
		}

		// Add a bill item assignment
		_, err = db.Exec(`
			INSERT INTO bill_item_assignments (bill_id, item_id, quantity, unit_price, subtotal, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, bill.ID, itemID, 2, 100.00, 200.00, time.Now(), time.Now())
		if err != nil {
			t.Fatalf("Failed to create bill item assignment: %v", err)
		}

		// Get the bill and verify
		retrieved, err := repo.GetByID(bill.ID)
		if err != nil {
			t.Fatalf("Failed to get bill: %v", err)
		}

		if retrieved == nil {
			t.Fatal("Expected bill to be found")
		}

		if len(retrieved.Items) != 1 {
			t.Errorf("Expected 1 item, got %d", len(retrieved.Items))
		}

		item := retrieved.Items[0]
		if item.Quantity != 2 {
			t.Errorf("Expected quantity 2, got %d", item.Quantity)
		}
		if item.UnitPrice != 100.00 {
			t.Errorf("Expected unit price 100.00, got %.2f", item.UnitPrice)
		}
		if item.Subtotal != 200.00 {
			t.Errorf("Expected subtotal 200.00, got %.2f", item.Subtotal)
		}
		if item.BillItem == nil {
			t.Error("Expected bill item to be loaded")
		} else {
			if item.BillItem.Description != "Test Item" {
				t.Errorf("Expected item description 'Test Item', got '%s'", item.BillItem.Description)
			}
		}
	})

	// Test GetAll with items
	t.Run("GetAll with items", func(t *testing.T) {
		bills, err := repo.GetAll()
		if err != nil {
			t.Fatalf("Failed to get bills: %v", err)
		}

		if len(bills) == 0 {
			t.Fatal("Expected bills to be found")
		}

		// Check the last bill (the one we just created)
		bill := bills[len(bills)-1]
		if len(bill.Items) != 1 {
			t.Errorf("Expected 1 item, got %d", len(bill.Items))
		}

		if len(bill.Items) > 0 {
			item := bill.Items[0]
			if item.Quantity != 2 {
				t.Errorf("Expected quantity 2, got %d", item.Quantity)
			}
			if item.UnitPrice != 100.00 {
				t.Errorf("Expected unit price 100.00, got %.2f", item.UnitPrice)
			}
			if item.Subtotal != 200.00 {
				t.Errorf("Expected subtotal 200.00, got %.2f", item.Subtotal)
			}
			if item.BillItem == nil {
				t.Error("Expected bill item to be loaded")
			} else {
				if item.BillItem.Description != "Test Item" {
					t.Errorf("Expected item description 'Test Item', got '%s'", item.BillItem.Description)
				}
			}
		}
	})
}
