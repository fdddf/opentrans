package database

import (
	"context"
	"fmt"
	"time"

	"github.com/fdddf/xcstrings-translator/internal/config"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Module is the FX module for database
var Module = fx.Module("database",
	fx.Provide(NewDatabase),
	fx.Provide(NewGormDB),
	fx.Provide(NewDatabaseServiceFromFX),
)

// Database holds the database connection
type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database connection
func NewDatabase(lc fx.Lifecycle, cfg *config.FXConfig) (*Database, error) {
	var err error
	gormDB, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{
		Logger: logger.Default,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// In a complete DI implementation, we don't set the global DB
	// All database access should be done through injected dependencies
	// DB = gormDB  // Commenting out the global assignment

	db := &Database{DB: gormDB}

	// Add lifecycle hooks
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Database connection established")

			// Run migrations when the service starts
			err := RunMigrations(db)
			if err != nil {
				fmt.Printf("Failed to run migrations: %v\n", err)
				return err
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Closing database connection")
			return sqlDB.Close()
		},
	})

	return db, nil
}

// Migrate runs database migrations
func Migrate(db *Database) error {
	err := db.AutoMigrate(
		&User{},
		&Project{},
		&Translation{},
		&ProviderConfig{},
		&UserActivity{},
		&App{},
		&AppLocalization{},
		&Subscription{},
		&AppUser{},
		&TranslationQueue{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	fmt.Println("Database migrations completed")
	return nil
}

// DatabaseServiceParams holds the dependencies for DatabaseService
type DatabaseServiceParams struct {
	fx.In

	DB *Database
}

// NewGormDB provides the raw gorm.DB instance for dependency injection
func NewGormDB(database *Database) *gorm.DB {
	return database.DB
}

// NewDatabaseServiceFromFX creates a new DatabaseService instance with fx injection
func NewDatabaseServiceFromFX(p DatabaseServiceParams) *DatabaseService {
	return NewDatabaseService(p.DB.DB)
}
