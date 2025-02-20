package repository

import (
	"database/sql"
	"time"

	"bills/internal/models"
)

// BillItemAssignmentRepository defines the interface for bill item assignment storage operations
type BillItemAssignmentRepository interface {
	Create(assignment *models.BillItemAssignment) error
	GetByID(id int64) (*models.BillItemAssignment, error)
	GetByBillID(billID int64) ([]*models.BillItemAssignment, error)
	Update(assignment *models.BillItemAssignment) error
	Delete(id int64) error
	DeleteByBillID(billID int64) error
}

// SQLiteBillItemAssignmentRepository implements BillItemAssignmentRepository using SQLite
type SQLiteBillItemAssignmentRepository struct {
	db *sql.DB
}

// NewSQLiteBillItemAssignmentRepository creates a new SQLite repository instance
func NewSQLiteBillItemAssignmentRepository(db *sql.DB) *SQLiteBillItemAssignmentRepository {
	return &SQLiteBillItemAssignmentRepository{db: db}
}

func (r *SQLiteBillItemAssignmentRepository) Create(assignment *models.BillItemAssignment) error {
	query := `
		INSERT INTO bill_item_assignments (
			bill_id, item_id, quantity, price, currency, exchange_rate,
			original_amount, eur_amount, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		assignment.BillID,
		assignment.ItemID,
		assignment.Quantity,
		assignment.Price,
		assignment.Currency,
		assignment.ExchangeRate,
		assignment.OriginalAmount,
		assignment.EURAmount,
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
	assignment.ID = id
	return nil
}

func (r *SQLiteBillItemAssignmentRepository) GetByID(id int64) (*models.BillItemAssignment, error) {
	assignment := &models.BillItemAssignment{
		BillItem: &models.BillItem{},
	}
	err := r.db.QueryRow(`
		SELECT a.id, a.bill_id, a.item_id, a.quantity, a.price, a.currency,
			   a.exchange_rate, a.original_amount, a.eur_amount, a.created_at, a.updated_at,
			   i.name, i.price, i.currency, i.created_at, i.updated_at
		FROM bill_item_assignments a
		LEFT JOIN bill_items i ON a.item_id = i.id
		WHERE a.id = ?
	`, id).Scan(
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
		&assignment.BillItem.Name,
		&assignment.BillItem.Price,
		&assignment.BillItem.Currency,
		&assignment.BillItem.CreatedAt,
		&assignment.BillItem.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return assignment, err
}

func (r *SQLiteBillItemAssignmentRepository) GetByBillID(billID int64) ([]*models.BillItemAssignment, error) {
	rows, err := r.db.Query(`
		SELECT a.id, a.bill_id, a.item_id, a.quantity, a.price, a.currency,
			   a.exchange_rate, a.original_amount, a.eur_amount, a.created_at, a.updated_at,
			   i.id, i.name, i.price, i.currency, i.created_at, i.updated_at
		FROM bill_item_assignments a
		LEFT JOIN bill_items i ON a.item_id = i.id
		WHERE a.bill_id = ?
		ORDER BY a.id ASC
	`, billID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []*models.BillItemAssignment
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
		assignments = append(assignments, assignment)
	}
	return assignments, rows.Err()
}

func (r *SQLiteBillItemAssignmentRepository) Update(assignment *models.BillItemAssignment) error {
	assignment.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE bill_item_assignments
		SET quantity = ?, price = ?, currency = ?, exchange_rate = ?,
			original_amount = ?, eur_amount = ?, updated_at = ?
		WHERE id = ?
	`,
		assignment.Quantity,
		assignment.Price,
		assignment.Currency,
		assignment.ExchangeRate,
		assignment.OriginalAmount,
		assignment.EURAmount,
		assignment.UpdatedAt,
		assignment.ID,
	)
	return err
}

func (r *SQLiteBillItemAssignmentRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM bill_item_assignments WHERE id = ?", id)
	return err
}

func (r *SQLiteBillItemAssignmentRepository) DeleteByBillID(billID int64) error {
	_, err := r.db.Exec("DELETE FROM bill_item_assignments WHERE bill_id = ?", billID)
	return err
}
