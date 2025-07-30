package db

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Migration represents a database migration
type Migration struct {
	Version   int
	Name      string
	UpSQL     string
	DownSQL   string
	UpFile    string
	DownFile  string
	Checksum  string
	AppliedAt *time.Time
}

// MigrationStatus represents the status of migrations
type MigrationStatus struct {
	Version   int
	Name      string
	Applied   bool
	AppliedAt *time.Time
}

// ensureMigrationsTable creates the schema_migrations table if it doesn't exist
func ensureMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT NOW(),
			checksum VARCHAR(64)
		);`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}
	return nil
}

// calculateChecksum calculates MD5 checksum of migration content
func calculateChecksum(content string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(content)))
}

// readMigrationFiles reads all migration files from the migrations directory
func readMigrationFiles() ([]Migration, error) {
	migrationsPath := "db/migrations"
	migrations := make(map[int]*Migration)

	err := filepath.WalkDir(migrationsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		filename := d.Name()

		// Parse filename: 000001_create_reports_table.up.sql
		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			return fmt.Errorf("invalid migration filename format: %s", filename)
		}

		versionStr := parts[0]
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			return fmt.Errorf("invalid version number in filename %s: %w", filename, err)
		}

		// Skip version 0 (our schema_migrations table creation)
		if version == 0 {
			return nil
		}

		// Determine if this is an up or down migration
		isUp := strings.Contains(filename, ".up.sql")
		isDown := strings.Contains(filename, ".down.sql")

		if !isUp && !isDown {
			return fmt.Errorf("migration file must be .up.sql or .down.sql: %s", filename)
		}

		// Get or create migration entry
		migration, exists := migrations[version]
		if !exists {
			// Extract name from filename
			nameWithExt := strings.Join(parts[1:], "_")
			name := strings.TrimSuffix(nameWithExt, ".up.sql")
			name = strings.TrimSuffix(name, ".down.sql")

			migration = &Migration{
				Version: version,
				Name:    name,
			}
			migrations[version] = migration
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Store content and calculate checksum
		contentStr := string(content)
		if isUp {
			migration.UpFile = path
			migration.UpSQL = contentStr
			migration.Checksum = calculateChecksum(contentStr)
		} else {
			migration.DownFile = path
			migration.DownSQL = contentStr
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to read migration files: %w", err)
	}

	// Convert map to sorted slice
	var result []Migration
	for _, migration := range migrations {
		// Ensure both up and down files were found
		if migration.UpSQL == "" || migration.DownSQL == "" {
			continue // Skip incomplete migrations
		}
		result = append(result, *migration)
	}

	// Sort by version
	sort.Slice(result, func(i, j int) bool {
		return result[i].Version < result[j].Version
	})

	return result, nil
}

// getAppliedMigrations returns a map of applied migration versions
func getAppliedMigrations() (map[int]*time.Time, error) {
	query := "SELECT version, applied_at FROM schema_migrations ORDER BY version"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[int]*time.Time)
	for rows.Next() {
		var versionStr string
		var appliedAt time.Time

		err := rows.Scan(&versionStr, &appliedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan migration row: %w", err)
		}

		version, err := strconv.Atoi(versionStr)
		if err != nil {
			continue // Skip invalid versions
		}

		applied[version] = &appliedAt
	}

	return applied, nil
}

// markMigrationAsApplied records a migration as applied
func markMigrationAsApplied(migration Migration) error {
	query := `
		INSERT INTO schema_migrations (version, applied_at, checksum) 
		VALUES ($1, NOW(), $2)
		ON CONFLICT (version) DO NOTHING`

	_, err := DB.Exec(query, fmt.Sprintf("%06d", migration.Version), migration.Checksum)
	if err != nil {
		return fmt.Errorf("failed to mark migration %d as applied: %w", migration.Version, err)
	}

	return nil
}

// removeMigrationRecord removes a migration record (for rollback)
func removeMigrationRecord(version int) error {
	query := "DELETE FROM schema_migrations WHERE version = $1"
	_, err := DB.Exec(query, fmt.Sprintf("%06d", version))
	if err != nil {
		return fmt.Errorf("failed to remove migration record %d: %w", version, err)
	}
	return nil
}

// executeMigrationSQL executes a migration SQL in a transaction
func executeMigrationSQL(sql string) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(sql)
	if err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	return tx.Commit()
}

// NewMigrate initializes the migration system
func NewMigrate() error {
	return ensureMigrationsTable()
}

// RunMigrations applies all pending migrations
func RunMigrations() error {
	// Ensure migrations table exists
	if err := ensureMigrationsTable(); err != nil {
		return err
	}

	// Read all migration files
	migrations, err := readMigrationFiles()
	if err != nil {
		return err
	}

	// Get applied migrations
	applied, err := getAppliedMigrations()
	if err != nil {
		return err
	}

	// Find pending migrations
	var pending []Migration
	for _, migration := range migrations {
		if _, isApplied := applied[migration.Version]; !isApplied {
			pending = append(pending, migration)
		}
	}

	if len(pending) == 0 {
		fmt.Println("No pending migrations")
		return nil
	}

	// Apply pending migrations
	for _, migration := range pending {
		fmt.Printf("Applying migration %d: %s\n", migration.Version, migration.Name)

		// Execute the up migration
		err := executeMigrationSQL(migration.UpSQL)
		if err != nil {
			return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
		}

		// Mark as applied
		err = markMigrationAsApplied(migration)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Applied migration %d successfully\n", migration.Version)
	}

	return nil
}

// RollbackMigrations rolls back to a specific version
func RollbackMigrations(targetVersion int) error {
	// Read all migration files
	migrations, err := readMigrationFiles()
	if err != nil {
		return err
	}

	// Get applied migrations
	applied, err := getAppliedMigrations()
	if err != nil {
		return err
	}

	// Find migrations to rollback (in reverse order)
	var toRollback []Migration
	for i := len(migrations) - 1; i >= 0; i-- {
		migration := migrations[i]
		if migration.Version > targetVersion {
			if _, isApplied := applied[migration.Version]; isApplied {
				toRollback = append(toRollback, migration)
			}
		}
	}

	// Rollback migrations
	for _, migration := range toRollback {
		fmt.Printf("Rolling back migration %d: %s\n", migration.Version, migration.Name)

		// Execute down migration
		err := executeMigrationSQL(migration.DownSQL)
		if err != nil {
			return fmt.Errorf("failed to rollback migration %d: %w", migration.Version, err)
		}

		// Remove from applied migrations
		err = removeMigrationRecord(migration.Version)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Rolled back migration %d successfully\n", migration.Version)
	}

	return nil
}

// GetMigrationsStatus returns the status of all migrations
func GetMigrationsStatus() ([]MigrationStatus, error) {
	// Read all migration files
	migrations, err := readMigrationFiles()
	if err != nil {
		return nil, err
	}

	// Get applied migrations
	applied, err := getAppliedMigrations()
	if err != nil {
		return nil, err
	}

	// Build status list
	var status []MigrationStatus
	for _, migration := range migrations {
		appliedAt, isApplied := applied[migration.Version]
		status = append(status, MigrationStatus{
			Version:   migration.Version,
			Name:      migration.Name,
			Applied:   isApplied,
			AppliedAt: appliedAt,
		})
	}

	return status, nil
}

// ValidateMigrations validates all migration files
func ValidateMigrations() error {
	migrations, err := readMigrationFiles()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		// Check that both up and down files exist
		if migration.UpFile == "" {
			return fmt.Errorf("missing up migration file for version %d", migration.Version)
		}
		if migration.DownFile == "" {
			return fmt.Errorf("missing down migration file for version %d", migration.Version)
		}

		// Check that SQL content is not empty
		if strings.TrimSpace(migration.UpSQL) == "" {
			return fmt.Errorf("empty up migration content for version %d", migration.Version)
		}
		if strings.TrimSpace(migration.DownSQL) == "" {
			return fmt.Errorf("empty down migration content for version %d", migration.Version)
		}
	}

	fmt.Printf("✓ Validated %d migrations\n", len(migrations))
	return nil
}
