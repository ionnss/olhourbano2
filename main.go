package main

import (
	"fmt"
	"net/http"
	"olhourbano2/config"
	"olhourbano2/db"
	"olhourbano2/routes"
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
