package handlers_test

import (
	"bills/db"
	"bills/internal/handlers"
	"bills/internal/models"
	"bills/internal/repository"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
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
	billItem := models.NewBillItem("Test Item", 100.00)
	if err := billItemRepo.Create(billItem); err != nil {
		t.Fatalf("Failed to create test bill item: %v", err)
	}

	return issuer.ID, receiver.ID, billItem.ID
}

func TestCreateBill(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	issuerID, receiverID, itemID := createTestData(t, db)

	// Initialize repositories
	billRepo := repository.NewSQLiteBillRepository(db)
	receiverRepo := repository.NewSQLiteReceiverRepository(db)
	issuerRepo := repository.NewSQLiteIssuerRepository(db)
	billItemRepo := repository.NewSQLiteBillItemRepository(db)
	billItemAssignRepo := repository.NewSQLiteBillItemAssignmentRepository(db)

	// Initialize handler
	tmpl := template.Must(template.New("test").Parse("{{.}}"))
	handler := handlers.NewBillHandler(billRepo, receiverRepo, issuerRepo, billItemRepo, billItemAssignRepo, tmpl)

	// Create Echo instance
	e := echo.New()

	// Create form values
	form := url.Values{}
	form.Set("due_date", time.Now().Format("2006-01-02"))
	form.Set("issuer_id", fmt.Sprintf("%d", issuerID))
	form.Set("receiver_id", fmt.Sprintf("%d", receiverID))
	form.Add("item_ids[]", fmt.Sprintf("%d", itemID))
	form.Add("quantities[]", "2")
	form.Add("unit_prices[]", "100.00")

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/bills", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test CreateBill
	if err := handler.CreateBill(c); err != nil {
		t.Fatalf("Failed to create bill: %v", err)
	}

	// Verify bill was created
	bills, err := billRepo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get bills: %v", err)
	}

	if len(bills) != 1 {
		t.Fatalf("Expected 1 bill, got %d", len(bills))
	}

	bill := bills[0]
	if bill.Amount != 200.00 {
		t.Errorf("Expected amount 200.00, got %.2f", bill.Amount)
	}
	if bill.IssuerID != issuerID {
		t.Errorf("Expected issuer ID %d, got %d", issuerID, bill.IssuerID)
	}
	if bill.ReceiverID != receiverID {
		t.Errorf("Expected receiver ID %d, got %d", receiverID, bill.ReceiverID)
	}

	if len(bill.Items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(bill.Items))
	}

	item := bill.Items[0]
	if item.ItemID != itemID {
		t.Errorf("Expected item ID %d, got %d", itemID, item.ItemID)
	}
	if item.Quantity != 2 {
		t.Errorf("Expected quantity 2, got %d", item.Quantity)
	}
	if item.UnitPrice != 100.00 {
		t.Errorf("Expected unit price 100.00, got %.2f", item.UnitPrice)
	}
	if item.Subtotal != 200.00 {
		t.Errorf("Expected subtotal 200.00, got %.2f", item.Subtotal)
	}
}
