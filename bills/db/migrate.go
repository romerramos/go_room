package db

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// MigrateDB applies all pending migrations to the database
func MigrateDB(db *sql.DB, dbName string) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create sqlite driver: %w", err)
	}

	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create source driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Get current version before migration
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get current version: %w", err)
	}
	fmt.Printf("Current migration version: %d, dirty: %v\n", version, dirty)

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Get version after migration
	version, dirty, err = m.Version()
	if err != nil {
		return fmt.Errorf("failed to get version after migration: %w", err)
	}
	fmt.Printf("Migration version after Up(): %d, dirty: %v\n", version, dirty)

	return nil
}

// DropDB drops all tables in the database
func DropDB(db *sql.DB, dbName string) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create sqlite driver: %w", err)
	}

	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create source driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to drop database: %w", err)
	}

	return nil
}
