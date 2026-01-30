package cmd

import (
	"fmt"
	"os"

	"github.com/fdddf/xcstrings-translator/internal/config"
	"github.com/fdddf/xcstrings-translator/internal/database"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// connectToDatabase creates a direct database connection for migration commands
func connectToDatabase() (*database.Database, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	gormDB, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	db := &database.Database{DB: gormDB}
	return db, nil
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := connectToDatabase()
		if err != nil {
			fmt.Printf("Failed to connect to database: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Running database migrations...")
		if err := database.RunMigrations(db); err != nil {
			fmt.Printf("Migration failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations completed successfully!")
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := connectToDatabase()
		if err != nil {
			fmt.Printf("Failed to connect to database: %v\n", err)
			os.Exit(1)
		}

		status, err := database.GetMigrationStatus(db)
		if err != nil {
			fmt.Printf("Failed to get migration status: %v\n", err)
			os.Exit(1)
		}

		if status.Version == 0 {
			fmt.Println("No migrations have been applied yet.")
			return
		}

		fmt.Printf("Current version: %d\n", status.Version)
		fmt.Printf("Dirty: %t\n", status.Dirty)
	},
}

var migrateRollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback the last migration (development only)",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := connectToDatabase()
		if err != nil {
			fmt.Printf("Failed to connect to database: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Rolling back last migration...")
		if err := database.Rollback(db); err != nil {
			fmt.Printf("Rollback failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Rollback completed successfully!")
	},
}

func init() {
	migrateCmd.AddCommand(migrateStatusCmd)
	migrateCmd.AddCommand(migrateRollbackCmd)
	rootCmd.AddCommand(migrateCmd)
}
