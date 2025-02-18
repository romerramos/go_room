package repository_test

import (
	"bills/internal/models"
	"bills/internal/repository"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
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
			updated_at DATETIME NOT NULL
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

func createTestBill(t *testing.T, db *sql.DB) int64 {
	// Create test bill
	result, err := db.Exec(`
		INSERT INTO bills (
			due_date, currency, original_total, eur_total, paid,
			issuer_id, receiver_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		time.Now(),
		models.DefaultCurrency(),
		0.0,
		0.0,
		false,
		1, // Dummy issuer ID
		1, // Dummy receiver ID
		time.Now(),
		time.Now(),
	)
	if err != nil {
		t.Fatalf("Failed to create test bill: %v", err)
	}

	billID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get bill ID: %v", err)
	}

	return billID
}

func createTestBillItem(t *testing.T, db *sql.DB) int64 {
	// Create test bill item
	billItemRepo := repository.NewSQLiteBillItemRepository(db)
	billItem := models.NewBillItem("Test Item", 100.00, models.DefaultCurrency())
	if err := billItemRepo.Create(billItem); err != nil {
		t.Fatalf("Failed to create test bill item: %v", err)
	}

	return billItem.ID
}

func TestBillItemAssignmentRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewSQLiteBillItemAssignmentRepository(db)
	billID := createTestBill(t, db)
	itemID := createTestBillItem(t, db)

	// Test Create
	t.Run("Create", func(t *testing.T) {
		assignment := models.NewBillItemAssignment(billID, itemID, 2, 100.00, models.DefaultCurrency(), 1.0)
		err := repo.Create(assignment)
		if err != nil {
			t.Fatalf("Failed to create assignment: %v", err)
		}

		if assignment.ID == 0 {
			t.Error("Expected assignment ID to be set")
		}
		if assignment.OriginalAmount != 200.00 {
			t.Errorf("Expected original amount to be 200.00, got %.2f", assignment.OriginalAmount)
		}
	})

	// Test GetByID
	t.Run("GetByID", func(t *testing.T) {
		assignment := models.NewBillItemAssignment(billID, itemID, 2, 100.00, models.DefaultCurrency(), 1.0)
		err := repo.Create(assignment)
		if err != nil {
			t.Fatalf("Failed to create assignment: %v", err)
		}

		retrieved, err := repo.GetByID(assignment.ID)
		if err != nil {
			t.Fatalf("Failed to get assignment: %v", err)
		}

		if retrieved == nil {
			t.Fatal("Expected assignment to be found")
		}

		if retrieved.Quantity != 2 {
			t.Errorf("Expected quantity 2, got %d", retrieved.Quantity)
		}
		if retrieved.Price != 100.00 {
			t.Errorf("Expected price 100.00, got %.2f", retrieved.Price)
		}
		if retrieved.OriginalAmount != 200.00 {
			t.Errorf("Expected original amount 200.00, got %.2f", retrieved.OriginalAmount)
		}
		if retrieved.BillItem == nil {
			t.Error("Expected bill item to be loaded")
		} else {
			if retrieved.BillItem.Name != "Test Item" {
				t.Errorf("Expected item name 'Test Item', got '%s'", retrieved.BillItem.Name)
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
		assignment.CalculateAmounts()

		err = repo.Update(assignment)
		if err != nil {
			t.Fatalf("Failed to update assignment: %v", err)
		}

		updated, err := repo.GetByID(assignment.ID)
		if err != nil {
			t.Fatalf("Failed to get updated assignment: %v", err)
		}

		if updated.Quantity != 3 {
			t.Errorf("Expected quantity 3, got %d", updated.Quantity)
		}
		if updated.OriginalAmount != 300.00 {
			t.Errorf("Expected original amount 300.00, got %.2f", updated.OriginalAmount)
		}
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		assignment := models.NewBillItemAssignment(billID, itemID, 2, 100.00, models.DefaultCurrency(), 1.0)
		err := repo.Create(assignment)
		if err != nil {
			t.Fatalf("Failed to create assignment: %v", err)
		}

		err = repo.Delete(assignment.ID)
		if err != nil {
			t.Fatalf("Failed to delete assignment: %v", err)
		}

		deleted, err := repo.GetByID(assignment.ID)
		if err != nil {
			t.Fatalf("Failed to check deleted assignment: %v", err)
		}
		if deleted != nil {
			t.Error("Expected assignment to be deleted")
		}
	})

	// Test DeleteByBillID
	t.Run("DeleteByBillID", func(t *testing.T) {
		assignment := models.NewBillItemAssignment(billID, itemID, 2, 100.00, models.DefaultCurrency(), 1.0)
		err := repo.Create(assignment)
		if err != nil {
			t.Fatalf("Failed to create assignment: %v", err)
		}

		err = repo.DeleteByBillID(billID)
		if err != nil {
			t.Fatalf("Failed to delete assignments by bill ID: %v", err)
		}

		assignments, err := repo.GetByBillID(billID)
		if err != nil {
			t.Fatalf("Failed to check deleted assignments: %v", err)
		}
		if len(assignments) > 0 {
			t.Error("Expected all assignments to be deleted")
		}
	})
}
