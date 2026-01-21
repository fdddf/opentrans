package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// AppliedMigration represents a database migration
type AppliedMigration struct {
	ID        uint   `gorm:"column:id;primaryKey" json:"id"`
	Name      string `gorm:"column:name;type:varchar(255);uniqueIndex;not null" json:"name"`
	AppliedAt string `gorm:"column:applied_at" json:"applied_at"` // Changed to string to avoid time field issues
}

// TableName returns the table name for the AppliedMigration model
func (AppliedMigration) TableName() string {
	return "schema_migrations"
}

// RunMigrations runs all pending migrations
func RunMigrations(db *Database) error {
	// Create migrations table manually with SQL to avoid GORM issues
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL,
		applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	)`).Error; err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Get all applied migrations using raw SQL
	rows, err := db.Raw("SELECT name FROM schema_migrations ORDER BY id ASC").Rows()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %v", err)
	}
	defer rows.Close()

	appliedMap := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan migration name: %v", err)
		}
		appliedMap[name] = true
	}

	// Run migrations in order
	for _, migration := range migrations {
		if !appliedMap[migration.name] {
			fmt.Printf("Running migration: %s\n", migration.name)
			if err := migration.up(db.DB); err != nil {
				return fmt.Errorf("failed to run migration %s: %v", migration.name, err)
			}

			// Record the migration using raw SQL
			if err := db.Exec("INSERT INTO schema_migrations (name, applied_at) VALUES (?, ?)",
				migration.name, time.Now().Format(time.RFC3339)).Error; err != nil {
				return fmt.Errorf("failed to record migration %s: %v", migration.name, err)
			}

			fmt.Printf("Migration %s completed\n", migration.name)
		}
	}

	// Ensure role column exists for legacy databases
	if err := addUserRoleColumn(db.DB); err != nil {
		return fmt.Errorf("failed to ensure user role column: %v", err)
	}

	// Ensure app origin column exists for legacy databases
	if err := addAppOriginColumn(db.DB); err != nil {
		return fmt.Errorf("failed to ensure app origin column: %v", err)
	}

	// Skip auto-migration since we're running explicit migrations above
	// Run AutoMigrate for all models to ensure schema is up to date
	// if err := autoMigrateModels(db); err != nil {
	// 	return fmt.Errorf("failed to auto migrate models: %v", err)
	// }

	return nil
}

// autoMigrateModels runs AutoMigrate for all models
func autoMigrateModels(db *Database) error {
	return db.AutoMigrate(
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
}

// migration represents a single migration
type migration struct {
	name string
	up   func(*gorm.DB) error
}

// MigrationsList exports the migrations list for testing
var MigrationsList = []string{
	"001_create_users_table",
	"002_create_projects_table",
	"003_create_translations_table",
	"004_create_provider_configs_table",
	"005_create_user_activities_table",
	"006_create_apps_table",
	"007_create_app_localizations_table",
	"008_create_subscriptions_table",
	"009_create_app_users_table",
	"010_create_translation_queue_table",
}

// migrations is the list of all migrations in order
var migrations = []migration{
	{
		name: "001_create_users_table",
		up:   createUsersTable,
	},
	{
		name: "002_create_projects_table",
		up:   createProjectsTable,
	},
	{
		name: "003_create_translations_table",
		up:   createTranslationsTable,
	},
	{
		name: "004_create_provider_configs_table",
		up:   createProviderConfigsTable,
	},
	{
		name: "005_create_user_activities_table",
		up:   createUserActivitiesTable,
	},
	{
		name: "006_create_apps_table",
		up:   createAppsTable,
	},
	{
		name: "007_create_app_localizations_table",
		up:   createAppLocalizationsTable,
	},
	{
		name: "008_create_subscriptions_table",
		up:   createSubscriptionsTable,
	},
	{
		name: "009_create_app_users_table",
		up:   createAppUsersTable,
	},
	{
		name: "010_create_translation_queue_table",
		up:   createTranslationQueueTable,
	},
}

// addAppOriginColumn ensures origin column exists for legacy DBs
func addAppOriginColumn(db *gorm.DB) error {
	return db.Exec(`ALTER TABLE apps ADD COLUMN IF NOT EXISTS origin VARCHAR(20) NOT NULL DEFAULT 'manual'`).Error
}

// Migration helper for adding role column if missing (idempotent for existing DBs)
func addUserRoleColumn(db *gorm.DB) error {
	return db.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) NOT NULL DEFAULT 'user'`).Error
}

// Migration functions
func createUsersTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE,
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			is_activated BOOLEAN NOT NULL DEFAULT FALSE,
			activation_code VARCHAR(255),
			role VARCHAR(20) NOT NULL DEFAULT 'user',
			is_subscribed BOOLEAN NOT NULL DEFAULT FALSE,
			subscription_type VARCHAR(50) NOT NULL DEFAULT 'free',
			subscription_end TIMESTAMP WITH TIME ZONE,
			max_apps INTEGER NOT NULL DEFAULT 1,
			max_translations INTEGER NOT NULL DEFAULT 1000,
			current_usage INTEGER NOT NULL DEFAULT 0,
			current_app_count INTEGER NOT NULL DEFAULT 0
		)
	`).Error
}

func createProjectsTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			file_content TEXT,
			file_name VARCHAR(255),
			source_language VARCHAR(50),
			content_structure JSONB,
			settings JSONB
		)
	`).Error
}

func createTranslationsTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS translations (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE,
			project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			key VARCHAR(255) NOT NULL,
			source_text TEXT NOT NULL,
			target_text TEXT,
			target_language VARCHAR(50) NOT NULL,
			state VARCHAR(50),
			translation_provider VARCHAR(50)
		)
	`).Error
}

func createProviderConfigsTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS provider_configs (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			provider_type VARCHAR(50) NOT NULL,
			config_data JSONB NOT NULL,
			is_default BOOLEAN NOT NULL DEFAULT FALSE
		)
	`).Error
}

func createUserActivitiesTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS user_activities (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			action VARCHAR(255) NOT NULL,
			details TEXT,
			ip_address VARCHAR(255),
			user_agent TEXT
		)
	`).Error
}

func createAppsTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS apps (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			bundle_id VARCHAR(255) UNIQUE NOT NULL,
			apple_id VARCHAR(255),
			primary_locale VARCHAR(50),
			apple_connect_token TEXT,
			origin VARCHAR(20) NOT NULL DEFAULT 'manual',
			short_description TEXT,
			long_description TEXT,
			keywords TEXT,
			support_url VARCHAR(255),
			marketing_url VARCHAR(255),
			privacy_url VARCHAR(255),
			version VARCHAR(50),
			app_category VARCHAR(100),
			is_ready_for_review BOOLEAN NOT NULL DEFAULT FALSE
		)
	`).Error
}

func createAppLocalizationsTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS app_localizations (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE,
			app_id INTEGER NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
			language_code VARCHAR(50) NOT NULL,
			name VARCHAR(255),
			subtitle VARCHAR(255),
			privacy_url VARCHAR(255),
			marketing_url VARCHAR(255),
			support_url VARCHAR(255),
			download_description TEXT,
			short_description TEXT,
			long_description TEXT,
			keywords TEXT,
			release_notes TEXT,
			UNIQUE(app_id, language_code)
		)
	`).Error
}

func createSubscriptionsTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS subscriptions (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			stripe_subscription_id VARCHAR(255) UNIQUE,
			stripe_customer_id VARCHAR(255),
			subscription_type VARCHAR(50) NOT NULL,
			subscription_status VARCHAR(50) NOT NULL,
			current_period_start TIMESTAMP WITH TIME ZONE NOT NULL,
			current_period_end TIMESTAMP WITH TIME ZONE NOT NULL,
			trial_end TIMESTAMP WITH TIME ZONE,
			cancel_at_period_end BOOLEAN NOT NULL DEFAULT FALSE
		)
	`).Error
}

func createAppUsersTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS app_users (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE,
			app_id INTEGER NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			role VARCHAR(50) NOT NULL DEFAULT 'viewer',
			UNIQUE(app_id, user_id)
		)
	`).Error
}

func createTranslationQueueTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS translation_queue (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
			app_id INTEGER REFERENCES apps(id) ON DELETE CASCADE,
			type VARCHAR(50) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			priority INTEGER NOT NULL DEFAULT 0,
			progress INTEGER NOT NULL DEFAULT 0,
			total INTEGER NOT NULL DEFAULT 0,
			done INTEGER NOT NULL DEFAULT 0,
			error TEXT,
			provider_type VARCHAR(50),
			source_language VARCHAR(50),
			target_languages JSONB,
			config_data JSONB,
			result_data JSONB
		)
	`).Error
}

// Rollback rolls back the last migration (for development only)
func Rollback(db *Database) error {
	// Get the last applied migration
	var lastMigration AppliedMigration
	if err := db.Order("id DESC").First(&lastMigration).Error; err != nil {
		return fmt.Errorf("no migration to rollback: %v", err)
	}

	// Find the migration function
	var migration *migration
	for i := len(migrations) - 1; i >= 0; i-- {
		if migrations[i].name == lastMigration.Name {
			migration = &migrations[i]
			break
		}
	}

	if migration == nil {
		return fmt.Errorf("migration %s not found", lastMigration.Name)
	}

	fmt.Printf("Rolling back migration: %s\n", migration.name)

	// Note: Implementing rollback functions would require additional work
	// For now, this is a placeholder
	// In production, you should implement proper rollback functions for each migration

	// Delete the migration record
	if err := db.Delete(&lastMigration).Error; err != nil {
		return fmt.Errorf("failed to delete migration record: %v", err)
	}

	fmt.Printf("Migration %s rolled back\n", migration.name)
	return nil
}

// GetMigrationStatus returns the status of all migrations
func GetMigrationStatus(db *Database) ([]AppliedMigration, error) {
	var migrations []AppliedMigration
	if err := db.Order("id ASC").Find(&migrations).Error; err != nil {
		return nil, err
	}
	return migrations, nil
}
