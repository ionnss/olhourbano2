package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"olhourbano2/db"
	"olhourbano2/models"
	"olhourbano2/services"
	"strconv"
)

// CommentRequest represents a request to create a comment
type CommentRequest struct {
	ReportID  int    `json:"report_id"`
	CPF       string `json:"cpf"`
	BirthDate string `json:"birth_date"`
	Content   string `json:"content"`
}

// CommentResponse represents the response for comment operations
type CommentResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message,omitempty"`
	Comment  *models.Comment    `json:"comment,omitempty"`
	Comments []*models.CommentDisplay `json:"comments,omitempty"`
	Total    int                `json:"total,omitempty"`
}

// CreateCommentHandler handles comment creation
func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.ReportID <= 0 || req.CPF == "" || req.BirthDate == "" || req.Content == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Validate content length
	if len(req.Content) > 500 {
		http.Error(w, "Comment content exceeds 500 character limit", http.StatusBadRequest)
		return
	}

	// Verify CPF with birth date
	verification, err := services.VerifyCPFWithBirthDate(req.CPF, req.BirthDate)
	if err != nil {
		log.Printf("Error verifying CPF: %v", err)
		http.Error(w, "Error verifying CPF", http.StatusInternalServerError)
		return
	}

	if !verification.Success || !verification.Valid {
		http.Error(w, "Invalid CPF or birth date", http.StatusBadRequest)
		return
	}

	// Hash the CPF
	hashedCPF := services.HashCPF(req.CPF)

	// Create the comment
	comment, err := services.CreateComment(db.DB, req.ReportID, hashedCPF, req.Content)
	if err != nil {
		log.Printf("Error creating comment: %v", err)
		http.Error(w, "Error creating comment", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := CommentResponse{
		Success: true,
		Message: "Comment created successfully",
		Comment:  comment,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCommentsHandler handles retrieving comments for a report
func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get report ID from query parameters
	reportIDStr := r.URL.Query().Get("report_id")
	if reportIDStr == "" {
		http.Error(w, "Missing report_id parameter", http.StatusBadRequest)
		return
	}

	reportID, err := strconv.Atoi(reportIDStr)
	if err != nil {
		http.Error(w, "Invalid report_id", http.StatusBadRequest)
		return
	}

	// Get optional parameters
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "recent"
	}
	if sort != "recent" && sort != "votes" {
		sort = "recent"
	}

	limit := 20 // Default limit for infinite scroll
	offset := 0

	var comments []*models.CommentDisplay
	comments, err = services.GetCommentsForReport(db.DB, reportID, sort, limit, offset)

	if err != nil {
		log.Printf("Error getting comments: %v", err)
		http.Error(w, "Error retrieving comments", http.StatusInternalServerError)
		return
	}

	// Get total comment count
	total, err := services.GetCommentCountForReport(db.DB, reportID)
	if err != nil {
		log.Printf("Error getting comment count: %v", err)
		// Continue without total count
		total = len(comments)
	}

	response := CommentResponse{
		Success:  true,
		Comments: comments,
		Total:    total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


