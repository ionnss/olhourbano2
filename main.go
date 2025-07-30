package main

import (
	"fmt"
	"log"
	"net/http"
	"olhourbano2/config"
	"olhourbano2/db"
	"olhourbano2/routes"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Olho Urbano Aberto")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return
	}
	fmt.Printf("Configuration loaded successfully (App Version: %s)\n", cfg.AppVersion)

	// Connect to the database
	db.DB, err = db.ConnectDB()
	if err != nil {
		fmt.Printf("Error connecting to the database: %v\n", err)
		return
	}
	defer db.DB.Close()
	fmt.Printf("Database connection established successfully (Host: %s:%s)\n", cfg.DBHost, cfg.DBPort)

	// Handle migration commands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			if err := db.RunMigrations(); err != nil {
				log.Fatalf("Error running migrations: %v\n", err)
			}
			fmt.Println("Migrations applied successfully")
			return

		case "migrate:status":
			status, err := db.GetMigrationsStatus()
			if err != nil {
				log.Fatalf("Error getting migration status: %v\n", err)
			}

			fmt.Println("Migration Status:")
			fmt.Println("================")
			for _, s := range status {
				statusStr := "❌ Pending"
				if s.Applied {
					statusStr = fmt.Sprintf("✅ Applied at %s", s.AppliedAt.Format("2006-01-02 15:04:05"))
				}
				fmt.Printf("Version %d: %s - %s\n", s.Version, s.Name, statusStr)
			}
			return

		case "migrate:rollback":
			if len(os.Args) < 3 {
				log.Fatalf("Usage: %s migrate:rollback <target_version>\n", os.Args[0])
			}

			targetVersion, err := strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatalf("Invalid target version: %v\n", err)
			}

			if err := db.RollbackMigrations(targetVersion); err != nil {
				log.Fatalf("Error rolling back migrations: %v\n", err)
			}
			fmt.Printf("Rolled back to version %d successfully\n", targetVersion)
			return

		case "migrate:validate":
			if err := db.ValidateMigrations(); err != nil {
				log.Fatalf("Migration validation failed: %v\n", err)
			}
			fmt.Println("All migrations are valid")
			return

		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			fmt.Println("Available commands:")
			fmt.Println("  migrate           - Apply pending migrations")
			fmt.Println("  migrate:status    - Show migration status")
			fmt.Println("  migrate:rollback <version> - Rollback to specific version")
			fmt.Println("  migrate:validate  - Validate migration files")
			return
		}
	}

	// Create routes
	r := routes.CreateRoutes()

	// Start server
	fmt.Println("Starting server...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
	fmt.Println("Server is running on port 8080")
}
