package repository

import (
	"database/sql"
	"fmt"
	"time"

	"bills/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

// BillRepository defines the interface for bill storage operations
type BillRepository interface {
	Create(bill *models.Bill) error
	GetByID(id int64) (*models.Bill, error)
	GetAll() ([]*models.Bill, error)
	Update(bill *models.Bill) error
	Delete(id int64) error
}

// SQLiteBillRepository implements BillRepository using SQLite
type SQLiteBillRepository struct {
	db *sql.DB
}

// NewSQLiteBillRepository creates a new SQLite repository instance
func NewSQLiteBillRepository(db *sql.DB) *SQLiteBillRepository {
	return &SQLiteBillRepository{db: db}
}

// InitDB initializes the database and creates the bills table
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create bills table with foreign keys for issuers and receivers
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS bills (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			amount REAL NOT NULL,
			due_date DATETIME NOT NULL,
			paid BOOLEAN DEFAULT FALSE,
			issuer_id INTEGER NOT NULL,
			receiver_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (issuer_id) REFERENCES issuers(id),
			FOREIGN KEY (receiver_id) REFERENCES receivers(id)
		)
	`)
	return db, err
}

func (r *SQLiteBillRepository) Create(bill *models.Bill) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if we return with error

	// Insert the bill
	result, err := tx.Exec(`
		INSERT INTO bills (amount, due_date, paid, issuer_id, receiver_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		bill.Amount,
		bill.DueDate,
		bill.Paid,
		bill.IssuerID,
		bill.ReceiverID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert bill: %w", err)
	}

	// Get the bill ID
	billID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get bill ID: %w", err)
	}
	bill.ID = billID

	// Insert bill item assignments
	for _, item := range bill.Items {
		item.BillID = billID
		_, err = tx.Exec(`
			INSERT INTO bill_item_assignments (bill_id, item_id, quantity, unit_price, subtotal, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`,
			item.BillID,
			item.ItemID,
			item.Quantity,
			item.UnitPrice,
			item.Subtotal,
			time.Now(),
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("failed to insert bill item assignment: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *SQLiteBillRepository) GetByID(id int64) (*models.Bill, error) {
	bill := &models.Bill{
		Issuer:   &models.Issuer{},
		Receiver: &models.Receiver{},
		Items:    make([]*models.BillItemAssignment, 0),
	}
	err := r.db.QueryRow(`
		SELECT b.id, b.amount, b.due_date, b.paid, b.issuer_id, b.receiver_id, b.created_at, b.updated_at,
			   i.name, i.vat_number, i.street, i.city, i.state, i.zip_code, i.country,
			   r.name, r.vat_number, r.street, r.city, r.state, r.zip_code, r.country
		FROM bills b
		LEFT JOIN issuers i ON b.issuer_id = i.id
		LEFT JOIN receivers r ON b.receiver_id = r.id
		WHERE b.id = ?
	`, id).Scan(
		&bill.ID,
		&bill.Amount,
		&bill.DueDate,
		&bill.Paid,
		&bill.IssuerID,
		&bill.ReceiverID,
		&bill.CreatedAt,
		&bill.UpdatedAt,
		// Issuer fields
		&bill.Issuer.Name,
		&bill.Issuer.VATNumber,
		&bill.Issuer.Street,
		&bill.Issuer.City,
		&bill.Issuer.State,
		&bill.Issuer.ZipCode,
		&bill.Issuer.Country,
		// Receiver fields
		&bill.Receiver.Name,
		&bill.Receiver.VATNumber,
		&bill.Receiver.Street,
		&bill.Receiver.City,
		&bill.Receiver.State,
		&bill.Receiver.ZipCode,
		&bill.Receiver.Country,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Load bill items
	rows, err := r.db.Query(`
		SELECT a.id, a.bill_id, a.item_id, a.quantity, a.unit_price, a.subtotal, a.created_at, a.updated_at,
			   i.id, i.description, i.default_price, i.created_at, i.updated_at
		FROM bill_item_assignments a
		LEFT JOIN bill_items i ON a.item_id = i.id
		WHERE a.bill_id = ?
		ORDER BY a.id ASC
	`, bill.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		assignment := &models.BillItemAssignment{
			BillItem: &models.BillItem{},
		}
		err := rows.Scan(
			&assignment.ID,
			&assignment.BillID,
			&assignment.ItemID,
			&assignment.Quantity,
			&assignment.UnitPrice,
			&assignment.Subtotal,
			&assignment.CreatedAt,
			&assignment.UpdatedAt,
			&assignment.BillItem.ID,
			&assignment.BillItem.Description,
			&assignment.BillItem.DefaultPrice,
			&assignment.BillItem.CreatedAt,
			&assignment.BillItem.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bill.Items = append(bill.Items, assignment)
	}

	return bill, rows.Err()
}

func (r *SQLiteBillRepository) GetAll() ([]*models.Bill, error) {
	rows, err := r.db.Query(`
		SELECT b.id, b.amount, b.due_date, b.paid, b.issuer_id, b.receiver_id, b.created_at, b.updated_at,
			   i.name, i.vat_number, i.street, i.city, i.state, i.zip_code, i.country,
			   r.name, r.vat_number, r.street, r.city, r.state, r.zip_code, r.country
		FROM bills b
		LEFT JOIN issuers i ON b.issuer_id = i.id
		LEFT JOIN receivers r ON b.receiver_id = r.id
		ORDER BY b.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []*models.Bill
	for rows.Next() {
		bill := &models.Bill{
			Issuer:   &models.Issuer{},
			Receiver: &models.Receiver{},
			Items:    make([]*models.BillItemAssignment, 0),
		}
		err := rows.Scan(
			&bill.ID,
			&bill.Amount,
			&bill.DueDate,
			&bill.Paid,
			&bill.IssuerID,
			&bill.ReceiverID,
			&bill.CreatedAt,
			&bill.UpdatedAt,
			// Issuer fields
			&bill.Issuer.Name,
			&bill.Issuer.VATNumber,
			&bill.Issuer.Street,
			&bill.Issuer.City,
			&bill.Issuer.State,
			&bill.Issuer.ZipCode,
			&bill.Issuer.Country,
			// Receiver fields
			&bill.Receiver.Name,
			&bill.Receiver.VATNumber,
			&bill.Receiver.Street,
			&bill.Receiver.City,
			&bill.Receiver.State,
			&bill.Receiver.ZipCode,
			&bill.Receiver.Country,
		)
		if err != nil {
			return nil, err
		}

		// Load bill items for each bill
		itemRows, err := r.db.Query(`
			SELECT a.id, a.bill_id, a.item_id, a.quantity, a.unit_price, a.subtotal, a.created_at, a.updated_at,
				   i.id, i.description, i.default_price, i.created_at, i.updated_at
			FROM bill_item_assignments a
			LEFT JOIN bill_items i ON a.item_id = i.id
			WHERE a.bill_id = ?
			ORDER BY a.id ASC
		`, bill.ID)
		if err != nil {
			return nil, err
		}
		defer itemRows.Close()

		for itemRows.Next() {
			assignment := &models.BillItemAssignment{
				BillItem: &models.BillItem{},
			}
			err := itemRows.Scan(
				&assignment.ID,
				&assignment.BillID,
				&assignment.ItemID,
				&assignment.Quantity,
				&assignment.UnitPrice,
				&assignment.Subtotal,
				&assignment.CreatedAt,
				&assignment.UpdatedAt,
				&assignment.BillItem.ID,
				&assignment.BillItem.Description,
				&assignment.BillItem.DefaultPrice,
				&assignment.BillItem.CreatedAt,
				&assignment.BillItem.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
			bill.Items = append(bill.Items, assignment)
		}
		if err = itemRows.Err(); err != nil {
			return nil, err
		}

		bills = append(bills, bill)
	}
	return bills, rows.Err()
}

func (r *SQLiteBillRepository) Update(bill *models.Bill) error {
	bill.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE bills
		SET amount = ?, due_date = ?, paid = ?, issuer_id = ?, receiver_id = ?, updated_at = ?
		WHERE id = ?
	`,
		bill.Amount,
		bill.DueDate,
		bill.Paid,
		bill.IssuerID,
		bill.ReceiverID,
		bill.UpdatedAt,
		bill.ID,
	)
	return err
}

func (r *SQLiteBillRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM bills WHERE id = ?", id)
	return err
}
