package db

import (
	"database/sql"
	"fmt"
	"olhourbano2/config"

	_ "github.com/lib/pq"
)

// DB is a global variable to hold the database connection
var DB *sql.DB

// ConnectDB connects to the database using the config package
func ConnectDB() (*sql.DB, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Open a connection to the database using the DSN from config
	DB, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Ping the database to ensure the connection is established
	if err := DB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Return the database connection
	return DB, nil
}
