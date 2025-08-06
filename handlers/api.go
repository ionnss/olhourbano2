package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"olhourbano2/db"
	"olhourbano2/services"
	"time"
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

// VoteRequest represents a vote request
type VoteRequest struct {
	ReportID int `json:"report_id"`
}

// VoteResponse represents the response for voting
type VoteResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	VoteCount int    `json:"vote_count,omitempty"`
}

// VoteHandler handles voting on reports
func VoteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var voteReq VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&voteReq); err != nil {
		log.Printf("Error decoding vote request: %v", err)
		response := VoteResponse{
			Success: false,
			Message: "Invalid request format",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate report ID
	if voteReq.ReportID <= 0 {
		response := VoteResponse{
			Success: false,
			Message: "Invalid report ID",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get hashed CPF from session or generate a temporary one
	// For now, we'll use a simple approach - in production, you'd want proper user authentication
	hashedCPF := "temp_user_" + fmt.Sprintf("%d", time.Now().Unix())

	// Add vote to database
	err := services.AddVote(db.DB, voteReq.ReportID, hashedCPF)
	if err != nil {
		log.Printf("Error adding vote: %v", err)
		response := VoteResponse{
			Success: false,
			Message: "Failed to register vote",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get updated vote count
	voteCount, err := services.GetVoteCount(db.DB, voteReq.ReportID)
	if err != nil {
		log.Printf("Error getting vote count: %v", err)
		// Continue anyway, just return 0 for vote count
		voteCount = 0
	}

	response := VoteResponse{
		Success:   true,
		Message:   "Vote registered successfully",
		VoteCount: voteCount,
	}

	json.NewEncoder(w).Encode(response)
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
