package database

import (
	"fmt"

	"github.com/fdddf/xcstrings-translator/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB holds the database connection
var DB *gorm.DB

// Connect initializes the database connection
func Connect() error {
	dbConfig := config.GetDatabaseConfig()

	var err error
	DB, err = gorm.Open(postgres.Open(dbConfig.URL), &gorm.Config{
		Logger: logger.Default,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Run migrations
	err = DB.AutoMigrate(&User{}, &Project{}, &Translation{}, &ProviderConfig{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	return nil
}

// Close closes the database connection
func Close() error {
	sql, err := DB.DB()
	if err != nil {
		return err
	}
	return sql.Close()
}
