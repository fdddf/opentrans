package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"

	"github.com/fdddf/xcstrings-translator/internal/database"
)

// generateDAOCmd represents the generate DAO command
var generateDAOCmd = &cobra.Command{
	Use:   "generate-dao",
	Short: "Generate DAO code using GORM gen",
	Long:  `Generate DAO code using GORM gen for all models`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize the generator
		g := gen.NewGenerator(gen.Config{
			OutPath:       "./internal/dao/query",                                             // Output path for generated code
			Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // Generate with specific modes
			FieldNullable: true,                                                               // Generate pointer fields for nullable columns
		})

		// Connect to database (for parsing table structures)
		// This is for code generation only, not for production use
		db, err := gorm.Open(postgres.Open("host=localhost user=i18n password=change_this_password dbname=i18n port=5432 sslmode=disable"))
		if err != nil {
			log.Fatal("failed to connect to database for code generation:", err)
		}

		// Use the database connection
		g.UseDB(db)

		// Generate code for all models defined in database package
		g.ApplyBasic(
			new(database.User),
			new(database.Project),
			new(database.Translation),
			new(database.UserActivity),
			new(database.ProviderConfig),
			new(database.App),
			new(database.AppLocalization),
			new(database.Subscription),
			new(database.AppUser),
			new(database.AppProviderConfig),
			new(database.TranslationQueue),
		)

		// Execute the generation
		g.Execute()

		log.Println("DAO code generation completed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(generateDAOCmd)
}
