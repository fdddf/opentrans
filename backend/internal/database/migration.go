package database

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs all pending migrations
func RunMigrations(db *Database) error {
	m, err := newMigrator(db)
	if err != nil {
		return err
	}
	// Don't close the migrator as it shares the database connection with GORM
	// defer closeMigrator(m)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %v", err)
	}

	return nil
}

// autoMigrateModels runs AutoMigrate for all models
// MigrationStatus represents the current migration status.
type MigrationStatus struct {
	Version uint
	Dirty   bool
}

func newMigrator(db *Database) (*migrate.Migrate, error) {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to access database: %v", err)
	}

	path, err := migrationsPath()
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(path, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migrator: %v", err)
	}

	return m, nil
}

func closeMigrator(m *migrate.Migrate) {
	_, _ = m.Close()
}

func migrationsPath() (string, error) {
	if custom := os.Getenv("MIGRATIONS_PATH"); custom != "" {
		return ensureFileURL(custom)
	}

	candidates := []string{
		filepath.Join("backend", "migrations"),
		"migrations",
	}

	for _, candidate := range candidates {
		if dirExists(candidate) {
			return ensureFileURL(candidate)
		}
	}

	if exe, err := os.Executable(); err == nil {
		candidate := filepath.Join(filepath.Dir(exe), "migrations")
		if dirExists(candidate) {
			return ensureFileURL(candidate)
		}
	}

	return "", fmt.Errorf("migrations directory not found; set MIGRATIONS_PATH")
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func ensureFileURL(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve migrations path: %v", err)
	}
	return "file://" + absPath, nil
}

// Rollback rolls back the last migration (for development only)
func Rollback(db *Database) error {
	m, err := newMigrator(db)
	if err != nil {
		return err
	}
	// Don't close the migrator as it shares the database connection with GORM
	// defer closeMigrator(m)

	if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to rollback migration: %v", err)
	}

	return nil
}

// GetMigrationStatus returns the status of all migrations
func GetMigrationStatus(db *Database) (*MigrationStatus, error) {
	m, err := newMigrator(db)
	if err != nil {
		return nil, err
	}
	// Don't close the migrator as it shares the database connection with GORM
	// defer closeMigrator(m)

	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			return &MigrationStatus{Version: 0, Dirty: false}, nil
		}
		return nil, fmt.Errorf("failed to read migration version: %v", err)
	}

	return &MigrationStatus{Version: version, Dirty: dirty}, nil
}
