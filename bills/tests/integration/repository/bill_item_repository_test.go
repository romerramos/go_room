package repository_test

import (
	"bills/internal/models"
	"bills/internal/repository"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupBillItemTestDB(t *testing.T) *sql.DB {
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
	`)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	return db
}

func TestBillItemRepository(t *testing.T) {
	db := setupBillItemTestDB(t)
	defer db.Close()

	repo := repository.NewSQLiteBillItemRepository(db)

	// Test Create
	t.Run("Create", func(t *testing.T) {
		item := models.NewBillItem("Test Item", 100.00, models.DefaultCurrency())
		err := repo.Create(item)
		if err != nil {
			t.Fatalf("Failed to create item: %v", err)
		}

		if item.ID == 0 {
			t.Error("Expected item ID to be set")
		}
		if item.Name != "Test Item" {
			t.Errorf("Expected item name 'Test Item', got '%s'", item.Name)
		}
		if item.Price != 100.00 {
			t.Errorf("Expected price 100.00, got %.2f", item.Price)
		}
		if item.Currency != models.DefaultCurrency() {
			t.Errorf("Expected currency '%s', got '%s'", models.DefaultCurrency(), item.Currency)
		}
	})

	// Test GetByID
	t.Run("GetByID", func(t *testing.T) {
		item := models.NewBillItem("Test Item", 100.00, models.DefaultCurrency())
		err := repo.Create(item)
		if err != nil {
			t.Fatalf("Failed to create item: %v", err)
		}

		retrieved, err := repo.GetByID(item.ID)
		if err != nil {
			t.Fatalf("Failed to get item: %v", err)
		}

		if retrieved == nil {
			t.Fatal("Expected item to be found")
		}
		if retrieved.Name != "Test Item" {
			t.Errorf("Expected item name 'Test Item', got '%s'", retrieved.Name)
		}
		if retrieved.Price != 100.00 {
			t.Errorf("Expected price 100.00, got %.2f", retrieved.Price)
		}
		if retrieved.Currency != models.DefaultCurrency() {
			t.Errorf("Expected currency '%s', got '%s'", models.DefaultCurrency(), retrieved.Currency)
		}
	})

	// Test GetAll
	t.Run("GetAll", func(t *testing.T) {
		items, err := repo.GetAll()
		if err != nil {
			t.Fatalf("Failed to get items: %v", err)
		}

		if len(items) == 0 {
			t.Error("Expected to find items")
		}

		for _, item := range items {
			if item.Name == "" {
				t.Error("Expected item name to be set")
			}
			if item.Price <= 0 {
				t.Errorf("Expected positive price, got %.2f", item.Price)
			}
			if item.Currency == "" {
				t.Error("Expected currency to be set")
			}
			if !models.IsSupportedCurrency(item.Currency) {
				t.Errorf("Expected supported currency, got '%s'", item.Currency)
			}
		}
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		item := models.NewBillItem("Test Item", 100.00, models.DefaultCurrency())
		err := repo.Create(item)
		if err != nil {
			t.Fatalf("Failed to create item: %v", err)
		}

		item.Name = "Updated Item"
		item.Price = 150.00
		item.Currency = "USD"
		err = repo.Update(item)
		if err != nil {
			t.Fatalf("Failed to update item: %v", err)
		}

		updated, err := repo.GetByID(item.ID)
		if err != nil {
			t.Fatalf("Failed to get updated item: %v", err)
		}

		if updated.Name != "Updated Item" {
			t.Errorf("Expected item name 'Updated Item', got '%s'", updated.Name)
		}
		if updated.Price != 150.00 {
			t.Errorf("Expected price 150.00, got %.2f", updated.Price)
		}
		if updated.Currency != "USD" {
			t.Errorf("Expected currency 'USD', got '%s'", updated.Currency)
		}
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		item := models.NewBillItem("Test Item", 100.00, models.DefaultCurrency())
		err := repo.Create(item)
		if err != nil {
			t.Fatalf("Failed to create item: %v", err)
		}

		err = repo.Delete(item.ID)
		if err != nil {
			t.Fatalf("Failed to delete item: %v", err)
		}

		deleted, err := repo.GetByID(item.ID)
		if err != nil {
			t.Fatalf("Failed to check deleted item: %v", err)
		}
		if deleted != nil {
			t.Error("Expected item to be deleted")
		}
	})
}
