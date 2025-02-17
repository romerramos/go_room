package repository

import (
	"database/sql"
	"time"

	"bills/internal/models"
)

// IssuerRepository defines the interface for issuer storage operations
type IssuerRepository interface {
	Create(issuer *models.Issuer) error
	GetByID(id int64) (*models.Issuer, error)
	GetAll() ([]*models.Issuer, error)
	Update(issuer *models.Issuer) error
	Delete(id int64) error
}

// SQLiteIssuerRepository implements IssuerRepository using SQLite
type SQLiteIssuerRepository struct {
	db *sql.DB
}

// NewSQLiteIssuerRepository creates a new SQLite repository instance
func NewSQLiteIssuerRepository(db *sql.DB) *SQLiteIssuerRepository {
	return &SQLiteIssuerRepository{db: db}
}

// InitIssuerDB initializes the issuers table
func InitIssuerDB(db *sql.DB) error {
	_, err := db.Exec(`
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
		)
	`)
	return err
}

func (r *SQLiteIssuerRepository) Create(issuer *models.Issuer) error {
	query := `
		INSERT INTO issuers (name, vat_number, street, city, state, zip_code, country, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		issuer.Name,
		issuer.VATNumber,
		issuer.Street,
		issuer.City,
		issuer.State,
		issuer.ZipCode,
		issuer.Country,
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
	issuer.ID = id
	return nil
}

func (r *SQLiteIssuerRepository) GetByID(id int64) (*models.Issuer, error) {
	issuer := &models.Issuer{}
	err := r.db.QueryRow(`
		SELECT id, name, vat_number, street, city, state, zip_code, country, created_at, updated_at
		FROM issuers WHERE id = ?
	`, id).Scan(
		&issuer.ID,
		&issuer.Name,
		&issuer.VATNumber,
		&issuer.Street,
		&issuer.City,
		&issuer.State,
		&issuer.ZipCode,
		&issuer.Country,
		&issuer.CreatedAt,
		&issuer.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return issuer, err
}

func (r *SQLiteIssuerRepository) GetAll() ([]*models.Issuer, error) {
	rows, err := r.db.Query(`
		SELECT id, name, vat_number, street, city, state, zip_code, country, created_at, updated_at
		FROM issuers ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issuers []*models.Issuer
	for rows.Next() {
		issuer := &models.Issuer{}
		err := rows.Scan(
			&issuer.ID,
			&issuer.Name,
			&issuer.VATNumber,
			&issuer.Street,
			&issuer.City,
			&issuer.State,
			&issuer.ZipCode,
			&issuer.Country,
			&issuer.CreatedAt,
			&issuer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		issuers = append(issuers, issuer)
	}
	return issuers, rows.Err()
}

func (r *SQLiteIssuerRepository) Update(issuer *models.Issuer) error {
	issuer.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE issuers
		SET name = ?, vat_number = ?, street = ?, city = ?, state = ?, zip_code = ?, country = ?, updated_at = ?
		WHERE id = ?
	`,
		issuer.Name,
		issuer.VATNumber,
		issuer.Street,
		issuer.City,
		issuer.State,
		issuer.ZipCode,
		issuer.Country,
		issuer.UpdatedAt,
		issuer.ID,
	)
	return err
}

func (r *SQLiteIssuerRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM issuers WHERE id = ?", id)
	return err
}
