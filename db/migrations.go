package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Migration represents a databse migrations
type Migration struct {
	ID   string
	File string
	SQL  string
}

// RunMigragtions runs the migrations files in order
func RunMigrations() error {
	log.Println("Running migrations...")

	// Get the migrations path directory
	migrationsDir := "db/migrations"
	if os.Getenv("MIGRATIONS_DIR") != "" {
		migrationsDir = os.Getenv("MIGRATIONS_DIR")
	}

	// Get the list of migration files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Filter and sort migrations files
	var migrations []Migration
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			filepath := filepath.Join(migrationsDir, file.Name())
			content, err := os.ReadFile(filepath)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", filepath, err)
			}

			migrations = append(migrations, Migration{
				ID:   strings.TrimSuffix(file.Name(), ".sql"),
				File: file.Name(),
				SQL:  string(content),
			})
		}
	}

	// Sort migrations by ID
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].ID < migrations[j].ID
	})

	// Create migrations table if it doesn't exist
	createMigrationsTable := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT NOW()
		);
	`
	_, err = DB.Exec(createMigrationsTable)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Execute migrations
	for _, migration := range migrations {
		// Check if migration is already applied
		var count int
		err := DB.QueryRow("SELECT COUNT(*) FROM migrations WHERE version = $1", migration.ID).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check if migration %s is applied: %w", migration.ID, err)
		}

		if count > 0 {
			log.Printf("Skipping migration %s (already applied)", migration.ID)
			continue
		}

		// Execute migration
		log.Printf("Applying migration %s", migration.ID)
		_, err = DB.Exec(migration.SQL)
		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.ID, err)
		}

		log.Printf("Migration applied successfully: %s", migration.ID)
	}

	log.Println("All migrations applied successfully")
	return nil
}
