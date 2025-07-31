package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"olhourbano2/config"
	"time"
)

// CPFVerificationRequest represents the request to CPFHub API
type CPFVerificationRequest struct {
	CPF       string `json:"cpf"`
	BirthDate string `json:"birthDate"`
}

// CPFVerificationResponse represents the response from CPFHub API
type CPFVerificationResponse struct {
	Success bool   `json:"success"`
	Valid   bool   `json:"valid,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    struct {
		Name      string `json:"name,omitempty"`
		Status    string `json:"status,omitempty"`
		Situation string `json:"situation,omitempty"`
		BirthDate string `json:"birthDate,omitempty"`
		CPFNumber string `json:"cpfNumber,omitempty"`
	} `json:"data,omitempty"`
}

// VerifyCPFWithBirthDate verifies CPF with birth date using CPFHub API
func VerifyCPFWithBirthDate(cpf, birthDate string) (*CPFVerificationResponse, error) {
	// First, do local CPF validation
	if !ValidateCPF(cpf) {
		return &CPFVerificationResponse{
			Valid:   false,
			Message: "CPF format is invalid",
		}, nil
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}

	// Prepare request
	normalizedCPF := NormalizeCPF(cpf)

	// Convert date from YYYY-MM-DD to DD/MM/YYYY format
	formattedBirthDate, err := convertDateFormat(birthDate)
	if err != nil {
		return &CPFVerificationResponse{
			Valid:   false,
			Message: "Invalid birth date format",
		}, nil
	}

	requestData := CPFVerificationRequest{
		CPF:       normalizedCPF,
		BirthDate: formattedBirthDate,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", cfg.CPFHubAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", cfg.CPFHubAPIKey)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to CPFHub: %v", err)
	}
	defer resp.Body.Close()

	// Parse response
	var cpfResponse CPFVerificationResponse
	if err := json.NewDecoder(resp.Body).Decode(&cpfResponse); err != nil {
		return nil, fmt.Errorf("error parsing CPFHub response: %v", err)
	}

	// Handle API errors
	if resp.StatusCode != http.StatusOK {
		return &CPFVerificationResponse{
			Valid:   false,
			Message: fmt.Sprintf("CPF verification failed (HTTP %d): %s", resp.StatusCode, cpfResponse.Error),
		}, nil
	}

	// Process the response according to CPFHub format
	if cpfResponse.Success {
		// Check if the CPF is valid and active
		isValid := cpfResponse.Data.Status != "Rejeitado" &&
			cpfResponse.Data.Situation != "TITULAR FALECIDO" &&
			cpfResponse.Data.Situation != "CPF CANCELADO" &&
			cpfResponse.Data.Situation != "CPF SUSPENSO"

		return &CPFVerificationResponse{
			Success: true,
			Valid:   isValid,
			Message: fmt.Sprintf("Status: %s - %s", cpfResponse.Data.Status, cpfResponse.Data.Situation),
		}, nil
	} else {
		return &CPFVerificationResponse{
			Success: false,
			Valid:   false,
			Message: "CPF verification failed",
		}, nil
	}
}

// convertDateFormat converts date from YYYY-MM-DD to DD/MM/YYYY
func convertDateFormat(dateStr string) (string, error) {
	// Parse the input date (YYYY-MM-DD)
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}

	// Format as DD/MM/YYYY
	return parsedDate.Format("02/01/2006"), nil
}

// MockCPFVerification is a fallback when CPFHub is not available
func MockCPFVerification(cpf, birthDate string) *CPFVerificationResponse {
	// Basic validation
	if !ValidateCPF(cpf) {
		return &CPFVerificationResponse{
			Success: false,
			Valid:   false,
			Message: "CPF format is invalid",
		}
	}

	// Parse birth date
	_, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return &CPFVerificationResponse{
			Success: false,
			Valid:   false,
			Message: "Invalid birth date format",
		}
	}

	// For development, always return true if CPF is valid
	return &CPFVerificationResponse{
		Success: true,
		Valid:   true,
		Message: "CPF validated (mock mode)",
	}
}
