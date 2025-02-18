package repository_test

import (
	"bills/internal/models"
	"bills/internal/repository"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupBillTestDB(t *testing.T) *sql.DB {
	// Use an in-memory database for testing
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		t.Fatalf("Failed to enable foreign keys: %v", err)
	}

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS issuers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			vat_number TEXT NOT NULL,
			street TEXT NOT NULL,
			city TEXT NOT NULL,
			state TEXT NOT NULL,
			zip_code TEXT NOT NULL,
			country TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS receivers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			vat_number TEXT NOT NULL,
			street TEXT NOT NULL,
			city TEXT NOT NULL,
			state TEXT NOT NULL,
			zip_code TEXT NOT NULL,
			country TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS bill_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			price REAL NOT NULL,
			currency TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS bills (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			due_date DATETIME NOT NULL,
			currency TEXT NOT NULL,
			original_total REAL NOT NULL,
			eur_total REAL NOT NULL,
			paid BOOLEAN DEFAULT FALSE,
			issuer_id INTEGER NOT NULL,
			receiver_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (issuer_id) REFERENCES issuers(id),
			FOREIGN KEY (receiver_id) REFERENCES receivers(id)
		);

		CREATE TABLE IF NOT EXISTS bill_item_assignments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bill_id INTEGER NOT NULL,
			item_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			price REAL NOT NULL,
			currency TEXT NOT NULL,
			exchange_rate REAL NOT NULL,
			original_amount REAL NOT NULL,
			eur_amount REAL NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (bill_id) REFERENCES bills(id) ON DELETE CASCADE,
			FOREIGN KEY (item_id) REFERENCES bill_items(id) ON DELETE RESTRICT
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	return db
}

func createTestData(t *testing.T, db *sql.DB) (int64, int64, int64) {
	// Create test issuer
	issuerRepo := repository.NewSQLiteIssuerRepository(db)
	issuer := models.NewIssuer("Test Issuer", "123456", "123 Street", "City", "State", "12345", "Country")
	if err := issuerRepo.Create(issuer); err != nil {
		t.Fatalf("Failed to create test issuer: %v", err)
	}

	// Create test receiver
	receiverRepo := repository.NewSQLiteReceiverRepository(db)
	receiver := models.NewReceiver("Test Receiver", "654321", "321 Street", "City", "State", "54321", "Country")
	if err := receiverRepo.Create(receiver); err != nil {
		t.Fatalf("Failed to create test receiver: %v", err)
	}

	// Create test bill item
	billItemRepo := repository.NewSQLiteBillItemRepository(db)
	billItem := models.NewBillItem("Test Item", 100.00, models.DefaultCurrency())
	if err := billItemRepo.Create(billItem); err != nil {
		t.Fatalf("Failed to create test bill item: %v", err)
	}

	return issuer.ID, receiver.ID, billItem.ID
}

func TestBillRepository(t *testing.T) {
	db := setupBillTestDB(t)
	defer db.Close()

	repo := repository.NewSQLiteBillRepository(db)
	issuerID, receiverID, itemID := createTestData(t, db)

	// Test Create
	t.Run("Create", func(t *testing.T) {
		bill := models.NewBill(time.Now(), issuerID, receiverID)
		assignment := models.NewBillItemAssignment(0, itemID, 2, 100.00, models.DefaultCurrency(), 1.0)
		bill.Items = append(bill.Items, assignment)
		bill.CalculateTotals()

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
		assignment := models.NewBillItemAssignment(0, itemID, 2, 100.00, models.DefaultCurrency(), 1.0)
		bill.Items = append(bill.Items, assignment)
		bill.CalculateTotals()

		err := repo.Create(bill)
		if err != nil {
			t.Fatalf("Failed to create bill: %v", err)
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
		if item.Price != 100.00 {
			t.Errorf("Expected price 100.00, got %.2f", item.Price)
		}
		if item.OriginalAmount != 200.00 {
			t.Errorf("Expected original amount 200.00, got %.2f", item.OriginalAmount)
		}
		if item.BillItem == nil {
			t.Error("Expected bill item to be loaded")
		} else {
			if item.BillItem.Name != "Test Item" {
				t.Errorf("Expected item name 'Test Item', got '%s'", item.BillItem.Name)
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
			t.Error("Expected to find bills")
		}

		for _, bill := range bills {
			if len(bill.Items) == 0 {
				t.Error("Expected bill to have items")
			}
			for _, item := range bill.Items {
				if item.Quantity != 2 {
					t.Errorf("Expected quantity 2, got %d", item.Quantity)
				}
				if item.Price != 100.00 {
					t.Errorf("Expected price 100.00, got %.2f", item.Price)
				}
				if item.OriginalAmount != 200.00 {
					t.Errorf("Expected original amount 200.00, got %.2f", item.OriginalAmount)
				}
				if item.BillItem == nil {
					t.Error("Expected bill item to be loaded")
				} else {
					if item.BillItem.Name != "Test Item" {
						t.Errorf("Expected item name 'Test Item', got '%s'", item.BillItem.Name)
					}
				}
			}
		}
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		bill := models.NewBill(time.Now(), issuerID, receiverID)
		assignment := models.NewBillItemAssignment(0, itemID, 2, 100.00, models.DefaultCurrency(), 1.0)
		bill.Items = append(bill.Items, assignment)
		bill.CalculateTotals()

		err := repo.Create(bill)
		if err != nil {
			t.Fatalf("Failed to create bill: %v", err)
		}

		bill.Paid = true
		err = repo.Update(bill)
		if err != nil {
			t.Fatalf("Failed to update bill: %v", err)
		}

		updated, err := repo.GetByID(bill.ID)
		if err != nil {
			t.Fatalf("Failed to get updated bill: %v", err)
		}

		if !updated.Paid {
			t.Error("Expected bill to be marked as paid")
		}
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		bill := models.NewBill(time.Now(), issuerID, receiverID)
		err := repo.Create(bill)
		if err != nil {
			t.Fatalf("Failed to create bill: %v", err)
		}

		err = repo.Delete(bill.ID)
		if err != nil {
			t.Fatalf("Failed to delete bill: %v", err)
		}

		deleted, err := repo.GetByID(bill.ID)
		if err != nil {
			t.Fatalf("Failed to check deleted bill: %v", err)
		}
		if deleted != nil {
			t.Error("Expected bill to be deleted")
		}
	})
}
