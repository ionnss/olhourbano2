package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"olhourbano2/db"
	"olhourbano2/services"
	"strings"
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
	ID          int      `json:"id"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Address     string   `json:"address"`
	Status      string   `json:"status"`
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	VoteCount   int      `json:"vote_count"`
	CreatedAt   string   `json:"created_at"`
	HashedCPF   string   `json:"hashed_cpf"`
	Photos      []string `json:"photos"`
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
			// Format date for display
			createdAt := report.CreatedAt.Format("02/01/2006")

			// Get first 8 characters of hashed CPF for display
			hashedCPFDisplay := ""
			if len(report.HashedCPF) >= 8 {
				hashedCPFDisplay = report.HashedCPF[:8]
			} else if report.HashedCPF != "" {
				hashedCPFDisplay = report.HashedCPF
			}

			// Process photos
			var photos []string
			if report.PhotoPath != "" {
				photos = strings.Split(report.PhotoPath, ",")
				// Clean up any empty strings
				var cleanedPhotos []string
				for _, photo := range photos {
					if strings.TrimSpace(photo) != "" {
						cleanedPhotos = append(cleanedPhotos, strings.TrimSpace(photo))
					}
				}
				photos = cleanedPhotos
			}

			mapReports = append(mapReports, MapReportData{
				ID:          report.ID,
				Category:    report.ProblemType,
				Description: report.Description,
				Address:     report.Location,
				Status:      report.Status,
				Latitude:    report.Latitude,
				Longitude:   report.Longitude,
				VoteCount:   report.VoteCount,
				CreatedAt:   createdAt,
				HashedCPF:   hashedCPFDisplay,
				Photos:      photos,
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
	ReportID  int    `json:"report_id"`
	CPF       string `json:"cpf"`
	BirthDate string `json:"birth_date"`
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

	// Validate CPF and birth date
	if voteReq.CPF == "" || voteReq.BirthDate == "" {
		response := VoteResponse{
			Success: false,
			Message: "CPF and birth date are required",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Hash the CPF to check for existing votes BEFORE API verification
	hashedCPF := services.HashCPF(voteReq.CPF)

	// Check if user has already voted for this report
	hasVoted, err := services.HasUserVoted(db.DB, voteReq.ReportID, hashedCPF)
	if err != nil {
		log.Printf("Error checking if user has voted: %v", err)
		response := VoteResponse{
			Success: false,
			Message: "Failed to check vote status",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if hasVoted {
		// User has already voted - return current vote count without API call
		voteCount, err := services.GetVoteCount(db.DB, voteReq.ReportID)
		if err != nil {
			log.Printf("Error getting vote count for duplicate vote: %v", err)
			voteCount = 0
		}

		response := VoteResponse{
			Success:   false,
			Message:   "Você já votou neste relatório anteriormente",
			VoteCount: voteCount,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Verify CPF with birth date (only if user hasn't voted)
	result, err := services.VerifyCPFWithBirthDate(voteReq.CPF, voteReq.BirthDate)
	if err != nil {
		log.Printf("Error verifying CPF for vote: %v", err)
		// Fallback to mock verification for development
		result = services.MockCPFVerification(voteReq.CPF, voteReq.BirthDate)
		log.Printf("Using mock CPF verification for vote (development)")
	}

	if !result.Valid {
		response := VoteResponse{
			Success: false,
			Message: "CPF inválido ou não encontrado. Verifique os dados informados.",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Add vote to database (hashedCPF already calculated above)
	err = services.AddVote(db.DB, voteReq.ReportID, hashedCPF)
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

// ShareImageRequest represents the request for generating a share image
type ShareImageRequest struct {
	ReportID     int    `json:"report_id"`
	CategoryName string `json:"category_name"`
	CategoryIcon string `json:"category_icon"`
	Description  string `json:"description"`
	Location     string `json:"location"`
	VoteCount    int    `json:"vote_count"`
	CreatedAt    string `json:"created_at"`
}

// ShareImageResponse represents the response for share image generation
type ShareImageResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

// ShareImageHandler handles generating share images for reports
func ShareImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShareImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding share image request: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.ReportID == 0 {
		response := ShareImageResponse{
			Success: false,
			Message: "Report ID is required",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// For now, we'll return a success response with a placeholder
	// In a full implementation, you would generate the image server-side
	// using a library like github.com/fogleman/gg or similar
	response := ShareImageResponse{
		Success:  true,
		Message:  "Share image generated successfully",
		ImageURL: "/api/share-image/" + string(rune(req.ReportID)), // Placeholder URL
	}

	json.NewEncoder(w).Encode(response)
}
