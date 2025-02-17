package repository

import (
	"database/sql"
	"time"

	"bills/internal/models"
)

// ReceiverRepository defines the interface for receiver storage operations
type ReceiverRepository interface {
	Create(receiver *models.Receiver) error
	GetByID(id int64) (*models.Receiver, error)
	GetAll() ([]*models.Receiver, error)
	Update(receiver *models.Receiver) error
	Delete(id int64) error
}

// SQLiteReceiverRepository implements ReceiverRepository using SQLite
type SQLiteReceiverRepository struct {
	db *sql.DB
}

// NewSQLiteReceiverRepository creates a new SQLite repository instance
func NewSQLiteReceiverRepository(db *sql.DB) *SQLiteReceiverRepository {
	return &SQLiteReceiverRepository{db: db}
}

// InitReceiverDB initializes the receivers table
func InitReceiverDB(db *sql.DB) error {
	_, err := db.Exec(`
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
		)
	`)
	return err
}

func (r *SQLiteReceiverRepository) Create(receiver *models.Receiver) error {
	query := `
		INSERT INTO receivers (name, vat_number, street, city, state, zip_code, country, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		receiver.Name,
		receiver.VATNumber,
		receiver.Street,
		receiver.City,
		receiver.State,
		receiver.ZipCode,
		receiver.Country,
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
	receiver.ID = id
	return nil
}

func (r *SQLiteReceiverRepository) GetByID(id int64) (*models.Receiver, error) {
	receiver := &models.Receiver{}
	err := r.db.QueryRow(`
		SELECT id, name, vat_number, street, city, state, zip_code, country, created_at, updated_at
		FROM receivers WHERE id = ?
	`, id).Scan(
		&receiver.ID,
		&receiver.Name,
		&receiver.VATNumber,
		&receiver.Street,
		&receiver.City,
		&receiver.State,
		&receiver.ZipCode,
		&receiver.Country,
		&receiver.CreatedAt,
		&receiver.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return receiver, err
}

func (r *SQLiteReceiverRepository) GetAll() ([]*models.Receiver, error) {
	rows, err := r.db.Query(`
		SELECT id, name, vat_number, street, city, state, zip_code, country, created_at, updated_at
		FROM receivers ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receivers []*models.Receiver
	for rows.Next() {
		receiver := &models.Receiver{}
		err := rows.Scan(
			&receiver.ID,
			&receiver.Name,
			&receiver.VATNumber,
			&receiver.Street,
			&receiver.City,
			&receiver.State,
			&receiver.ZipCode,
			&receiver.Country,
			&receiver.CreatedAt,
			&receiver.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		receivers = append(receivers, receiver)
	}
	return receivers, rows.Err()
}

func (r *SQLiteReceiverRepository) Update(receiver *models.Receiver) error {
	receiver.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE receivers
		SET name = ?, vat_number = ?, street = ?, city = ?, state = ?, zip_code = ?, country = ?, updated_at = ?
		WHERE id = ?
	`,
		receiver.Name,
		receiver.VATNumber,
		receiver.Street,
		receiver.City,
		receiver.State,
		receiver.ZipCode,
		receiver.Country,
		receiver.UpdatedAt,
		receiver.ID,
	)
	return err
}

func (r *SQLiteReceiverRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM receivers WHERE id = ?", id)
	return err
}
