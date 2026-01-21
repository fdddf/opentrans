package database

import (
	"fmt"

	"gorm.io/gorm"
)

// DBInterface defines the interface for database operations
type DBInterface interface {
	GetDB() *gorm.DB
	Close() error
}

// DatabaseService wraps the database functionality with dependency injection
type DatabaseService struct {
	DB *gorm.DB
}

// NewDatabaseService creates a new DatabaseService instance
func NewDatabaseService(db *gorm.DB) *DatabaseService {
	return &DatabaseService{
		DB: db,
	}
}

// GetDB returns the database instance
func (ds *DatabaseService) GetDB() *gorm.DB {
	return ds.DB
}

// Connect initializes the database connection (placeholder)
func Connect() error {
	// This function now returns an error as database connections
	// should be handled exclusively through dependency injection
	return fmt.Errorf("database should be initialized through fx dependency injection - direct connection not supported")
}

// Close closes the database connection using direct gorm.DB instance
func Close(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	sql, err := db.DB()
	if err != nil {
		return err
	}
	return sql.Close()
}

// CloseWithService closes the database connection using the service
func (ds *DatabaseService) Close() error {
	if ds.DB == nil {
		return fmt.Errorf("database not initialized")
	}

	sql, err := ds.DB.DB()
	if err != nil {
		return err
	}
	return sql.Close()
}
