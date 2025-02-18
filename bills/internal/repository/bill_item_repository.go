package repository

import (
	"database/sql"
	"time"

	"bills/internal/models"
)

// BillItemRepository defines the interface for bill item storage operations
type BillItemRepository interface {
	Create(item *models.BillItem) error
	GetByID(id int64) (*models.BillItem, error)
	GetAll() ([]*models.BillItem, error)
	Update(item *models.BillItem) error
	Delete(id int64) error
}

// SQLiteBillItemRepository implements BillItemRepository using SQLite
type SQLiteBillItemRepository struct {
	db *sql.DB
}

// NewSQLiteBillItemRepository creates a new SQLite repository instance
func NewSQLiteBillItemRepository(db *sql.DB) *SQLiteBillItemRepository {
	return &SQLiteBillItemRepository{db: db}
}

// InitBillItemDB initializes the bill_items table
func InitBillItemDB(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS bill_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			price REAL NOT NULL,
			currency TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create bill_item_assignments table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS bill_item_assignments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bill_id INTEGER NOT NULL,
			item_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			unit_price REAL NOT NULL,
			subtotal REAL NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (bill_id) REFERENCES bills(id) ON DELETE CASCADE,
			FOREIGN KEY (item_id) REFERENCES bill_items(id) ON DELETE RESTRICT
		)
	`)
	return err
}

func (r *SQLiteBillItemRepository) Create(item *models.BillItem) error {
	query := `
		INSERT INTO bill_items (name, price, currency, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		item.Name,
		item.Price,
		item.Currency,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	item.ID = id
	return nil
}

func (r *SQLiteBillItemRepository) GetByID(id int64) (*models.BillItem, error) {
	item := &models.BillItem{}
	err := r.db.QueryRow(`
		SELECT id, name, price, currency, created_at, updated_at
		FROM bill_items WHERE id = ?
	`, id).Scan(
		&item.ID,
		&item.Name,
		&item.Price,
		&item.Currency,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return item, err
}

func (r *SQLiteBillItemRepository) GetAll() ([]*models.BillItem, error) {
	rows, err := r.db.Query(`
		SELECT id, name, price, currency, created_at, updated_at
		FROM bill_items ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.BillItem
	for rows.Next() {
		item := &models.BillItem{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Price,
			&item.Currency,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *SQLiteBillItemRepository) Update(item *models.BillItem) error {
	item.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE bill_items
		SET name = ?, price = ?, currency = ?, updated_at = ?
		WHERE id = ?
	`,
		item.Name,
		item.Price,
		item.Currency,
		item.UpdatedAt,
		item.ID,
	)
	return err
}

func (r *SQLiteBillItemRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM bill_items WHERE id = ?", id)
	return err
}
