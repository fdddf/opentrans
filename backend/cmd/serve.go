package cmd

import (
	"fmt"
	"time"

	"github.com/fdddf/xcstrings-translator/internal/config"
	"github.com/fdddf/xcstrings-translator/internal/database"
	"github.com/fdddf/xcstrings-translator/internal/server"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// serveCmd launches the Fiber web UI server.
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web UI for visual localisation",
	RunE: func(cmd *cobra.Command, args []string) error {
		addr, _ := cmd.Flags().GetString("addr")
		fmt.Printf("Serving web UI on %s\n", addr)

		// Get database configuration
		dbConfig := config.GetDatabaseConfig()

		// Initialize database connection
		gormDB, err := gorm.Open(postgres.Open(dbConfig.URL), &gorm.Config{
			Logger: logger.Default,
		})
		if err != nil {
			return fmt.Errorf("failed to connect to database: %v", err)
		}

		sqlDB, err := gormDB.DB()
		if err != nil {
			return fmt.Errorf("failed to get database instance: %v", err)
		}

		// Set connection pool settings
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
		defer sqlDB.Close()

		db := &database.Database{DB: gormDB}

		// Run migrations
		if err := database.RunMigrations(db); err != nil {
			return fmt.Errorf("failed to run migrations: %v", err)
		}

		return server.Serve(addr, db)
	},
}

func init() {
	serveCmd.Flags().String("addr", ":8080", "listen address for the web UI")
	rootCmd.AddCommand(serveCmd)
}
