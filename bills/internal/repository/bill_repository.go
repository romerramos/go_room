package repository

import (
	"database/sql"
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
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert bill
	query := `
		INSERT INTO bills (
			due_date, currency, original_total, eur_total, paid,
			issuer_id, receiver_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := tx.Exec(query,
		bill.DueDate,
		bill.Currency,
		bill.OriginalTotal,
		bill.EURTotal,
		bill.Paid,
		bill.IssuerID,
		bill.ReceiverID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	billID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	bill.ID = billID

	// Insert bill items
	for _, item := range bill.Items {
		item.BillID = billID
		query = `
			INSERT INTO bill_item_assignments (
				bill_id, item_id, quantity, price, currency,
				exchange_rate, original_amount, eur_amount,
				created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		_, err = tx.Exec(query,
			item.BillID,
			item.ItemID,
			item.Quantity,
			item.Price,
			item.Currency,
			item.ExchangeRate,
			item.OriginalAmount,
			item.EURAmount,
			time.Now(),
			time.Now(),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLiteBillRepository) GetByID(id int64) (*models.Bill, error) {
	bill := &models.Bill{}
	err := r.db.QueryRow(`
		SELECT b.id, b.due_date, b.paid, b.issuer_id, b.receiver_id,
			   b.currency, b.original_total, b.eur_total,
			   b.created_at, b.updated_at,
			   i.name as issuer_name, r.name as receiver_name
		FROM bills b
		LEFT JOIN issuers i ON b.issuer_id = i.id
		LEFT JOIN receivers r ON b.receiver_id = r.id
		WHERE b.id = ?
	`, id).Scan(
		&bill.ID,
		&bill.DueDate,
		&bill.Paid,
		&bill.IssuerID,
		&bill.ReceiverID,
		&bill.Currency,
		&bill.OriginalTotal,
		&bill.EURTotal,
		&bill.CreatedAt,
		&bill.UpdatedAt,
		&bill.IssuerName,
		&bill.ReceiverName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Load items
	rows, err := r.db.Query(`
		SELECT a.id, a.bill_id, a.item_id, a.quantity, a.price,
			   a.currency, a.exchange_rate, a.original_amount, a.eur_amount,
			   a.created_at, a.updated_at,
			   i.id, i.name, i.price, i.currency, i.created_at, i.updated_at
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
			&assignment.Price,
			&assignment.Currency,
			&assignment.ExchangeRate,
			&assignment.OriginalAmount,
			&assignment.EURAmount,
			&assignment.CreatedAt,
			&assignment.UpdatedAt,
			&assignment.BillItem.ID,
			&assignment.BillItem.Name,
			&assignment.BillItem.Price,
			&assignment.BillItem.Currency,
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
		SELECT b.id, b.due_date, b.paid, b.issuer_id, b.receiver_id,
			   b.currency, b.original_total, b.eur_total,
			   b.created_at, b.updated_at,
			   i.name as issuer_name, r.name as receiver_name
		FROM bills b
		LEFT JOIN issuers i ON b.issuer_id = i.id
		LEFT JOIN receivers r ON b.receiver_id = r.id
		ORDER BY b.due_date DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []*models.Bill
	for rows.Next() {
		bill := &models.Bill{}
		err := rows.Scan(
			&bill.ID,
			&bill.DueDate,
			&bill.Paid,
			&bill.IssuerID,
			&bill.ReceiverID,
			&bill.Currency,
			&bill.OriginalTotal,
			&bill.EURTotal,
			&bill.CreatedAt,
			&bill.UpdatedAt,
			&bill.IssuerName,
			&bill.ReceiverName,
		)
		if err != nil {
			return nil, err
		}

		// Load items for each bill
		itemRows, err := r.db.Query(`
			SELECT a.id, a.bill_id, a.item_id, a.quantity, a.price,
				   a.currency, a.exchange_rate, a.original_amount, a.eur_amount,
				   a.created_at, a.updated_at,
				   i.id, i.name, i.price, i.currency, i.created_at, i.updated_at
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
				&assignment.Price,
				&assignment.Currency,
				&assignment.ExchangeRate,
				&assignment.OriginalAmount,
				&assignment.EURAmount,
				&assignment.CreatedAt,
				&assignment.UpdatedAt,
				&assignment.BillItem.ID,
				&assignment.BillItem.Name,
				&assignment.BillItem.Price,
				&assignment.BillItem.Currency,
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
		SET due_date = ?, currency = ?, original_total = ?, eur_total = ?,
			paid = ?, issuer_id = ?, receiver_id = ?, updated_at = ?
		WHERE id = ?
	`,
		bill.DueDate,
		bill.Currency,
		bill.OriginalTotal,
		bill.EURTotal,
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
