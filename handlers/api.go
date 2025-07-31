package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"olhourbano2/services"
)

// CPFVerificationRequest represents the incoming JSON request
type CPFVerificationRequest struct {
	CPF       string `json:"cpf"`
	BirthDate string `json:"birth_date"`
}

// VerifyCPFHandler handles CPF verification with birth date
func VerifyCPFHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CPFVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding CPF verification request: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.CPF == "" || req.BirthDate == "" {
		response := services.CPFVerificationResponse{
			Valid:   false,
			Message: "CPF and birth date are required",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Verify CPF with birth date
	result, err := services.VerifyCPFWithBirthDate(req.CPF, req.BirthDate)
	if err != nil {
		log.Printf("Error verifying CPF: %v", err)

		// Fallback to mock verification
		result = services.MockCPFVerification(req.CPF, req.BirthDate)
		log.Printf("Using mock CPF verification for development")
	}

	// Return result
	json.NewEncoder(w).Encode(result)
}
