package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"olhourbano2/db"
	"olhourbano2/services"
)

// CPFVerificationRequest represents the incoming JSON request
type CPFVerificationRequest struct {
	CPF       string `json:"cpf"`
	BirthDate string `json:"birth_date"`
}

// MapReportsResponse represents the response for map reports
type MapReportsResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Reports []MapReportData `json:"reports,omitempty"`
}

// MapReportData represents a report for the map
type MapReportData struct {
	ID          int     `json:"id"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Address     string  `json:"address"`
	Status      string  `json:"status"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
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

// MapReportsHandler handles fetching reports for the map
func MapReportsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get query parameters
	category := r.URL.Query().Get("category")
	status := r.URL.Query().Get("status")
	city := r.URL.Query().Get("city")

	// Fetch reports with location data
	reports, err := services.GetReportsForMap(db.DB, category, status, city)
	if err != nil {
		log.Printf("Error fetching reports for map: %v", err)
		response := MapReportsResponse{
			Success: false,
			Message: "Failed to fetch reports",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to map format
	mapReports := make([]MapReportData, 0, len(reports))
	for _, report := range reports {
		if report.Latitude != 0 && report.Longitude != 0 {
			mapReports = append(mapReports, MapReportData{
				ID:          report.ID,
				Category:    report.ProblemType,
				Description: report.Description,
				Address:     report.Location,
				Status:      report.Status,
				Latitude:    report.Latitude,
				Longitude:   report.Longitude,
			})
		}
	}

	response := MapReportsResponse{
		Success: true,
		Reports: mapReports,
	}

	json.NewEncoder(w).Encode(response)
}

// CitiesResponse represents the response for cities
type CitiesResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Cities  []string `json:"cities,omitempty"`
}

// CitiesHandler handles fetching cities from reports
func CitiesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Fetch cities from reports
	cities, err := services.GetCitiesFromReports(db.DB)
	if err != nil {
		log.Printf("Error fetching cities: %v", err)
		response := CitiesResponse{
			Success: false,
			Message: "Failed to fetch cities",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := CitiesResponse{
		Success: true,
		Cities:  cities,
	}

	json.NewEncoder(w).Encode(response)
}
