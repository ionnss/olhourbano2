package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all application configuration
type Config struct {
	// Database Configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBName     string
	DBPassword string

	// Email Configuration
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string

	// Security Configuration
	SessionKey   string
	CookieDomain string

	// App Configuration
	AppVersion string

	// API Keys
	CPFHubAPIKey     string
	CPFHubAPIURL     string
	GoogleMapsAPIKey string
	GoogleMapsAPIURL string
}

// readSecretFile reads a secret from a file path
// Handles both local development (./secrets/) and Docker secrets (/run/secrets/)
func readSecretFile(filePath string) (string, error) {
	if filePath == "" {
		return "", fmt.Errorf("secret file path is empty")
	}

	// Try to read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		// If file doesn't exist in Docker secrets path, try local secrets
		if strings.HasPrefix(filePath, "/run/secrets/") {
			localPath := strings.Replace(filePath, "/run/secrets/", "./secrets/", 1)
			content, err = os.ReadFile(localPath)
			if err != nil {
				return "", fmt.Errorf("failed to read secret file: %w", err)
			}
		} else {
			return "", fmt.Errorf("failed to read secret file: %w", err)
		}
	}

	// Remove any trailing newlines or whitespace
	return strings.TrimSpace(string(content)), nil
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsIntOrDefault returns environment variable as int or default if not set/invalid
func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Load loads all configuration from environment variables and secret files
func Load() (*Config, error) {
	config := &Config{}

	// Database Configuration
	config.DBHost = getEnvOrDefault("DB_HOST", "localhost")
	config.DBPort = getEnvOrDefault("DB_PORT", "5432")
	config.DBUser = getEnvOrDefault("DB_USER", "postgres")
	config.DBName = getEnvOrDefault("DB_NAME", "olhourbanovault")

	// Read database password from secret file
	dbPasswordFile := getEnvOrDefault("DB_PASSWORD_FILE", "/run/secrets/db_password")
	var err error
	config.DBPassword, err = readSecretFile(dbPasswordFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load database password")
	}

	// Email Configuration
	config.SMTPHost = getEnvOrDefault("SMTP_HOST", "smtp.gmail.com")
	config.SMTPPort = getEnvAsIntOrDefault("SMTP_PORT", 587)
	config.SMTPUsername = getEnvOrDefault("SMTP_USERNAME", "")

	// Read SMTP password from secret file
	smtpPasswordFile := getEnvOrDefault("SMTP_PASSWORD_FILE", "/run/secrets/smtp_password")
	config.SMTPPassword, err = readSecretFile(smtpPasswordFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load SMTP password")
	}

	// Security Configuration
	sessionKeyFile := getEnvOrDefault("SESSION_KEY_FILE", "/run/secrets/session_key")
	config.SessionKey, err = readSecretFile(sessionKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load session key")
	}

	config.CookieDomain = getEnvOrDefault("COOKIE_DOMAIN", ".olhourbano.com")

	// App Configuration
	config.AppVersion = getEnvOrDefault("APP_VERSION", "2.0.0")

	// API Keys
	config.CPFHubAPIURL = getEnvOrDefault("CPFHUB_API_URL", "https://api.cpfhub.io/api/cpf")
	config.GoogleMapsAPIURL = getEnvOrDefault("GOOGLE_MAPS_API_URL", "https://maps.googleapis.com/maps/api")

	cpfhubAPIKeyFile := getEnvOrDefault("CPFHUB_API_KEY_FILE", "/run/secrets/cpfhub_api_key")
	config.CPFHubAPIKey, err = readSecretFile(cpfhubAPIKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load CPFHub API key")
	}

	googleMapsAPIKeyFile := getEnvOrDefault("GOOGLE_MAPS_API_KEY_FILE", "/run/secrets/google_maps_api_key")
	config.GoogleMapsAPIKey, err = readSecretFile(googleMapsAPIKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load Google Maps API key")
	}

	return config, nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		c.DBUser, c.DBPassword, c.DBName, c.DBHost, c.DBPort)
}

// String returns a safe representation of the config (no secrets)
func (c *Config) String() string {
	return fmt.Sprintf("Config{DBHost:%s, DBPort:%s, DBUser:%s, DBName:%s, SMTPHost:%s, SMTPPort:%d, SMTPUsername:%s, CookieDomain:%s, AppVersion:%s, CPFHubAPIURL:%s, GoogleMapsAPIURL:%s}",
		c.DBHost, c.DBPort, c.DBUser, c.DBName, c.SMTPHost, c.SMTPPort, c.SMTPUsername, c.CookieDomain, c.AppVersion, c.CPFHubAPIURL, c.GoogleMapsAPIURL)
}
