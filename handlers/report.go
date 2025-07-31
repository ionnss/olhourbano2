package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"olhourbano2/config"
	"olhourbano2/db"
	"olhourbano2/models"
	"olhourbano2/services"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// ReportHandler handles step 1: category selection
func ReportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// GET: Show category selection
		categories := config.GetAllCategories()

		data := map[string]interface{}{
			"Step":       1,
			"Categories": categories,
			"PageTitle":  "Nova Denúncia - Selecionar Categoria",
		}

		if err := renderTemplate(w, "01_report.html", data); err != nil {
			log.Printf("Error rendering report step 1 template: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// POST: Process category selection and redirect to step 2
	if r.Method == "POST" {
		categoryID := r.FormValue("category")
		if categoryID == "" {
			http.Error(w, "Categoria é obrigatória", http.StatusBadRequest)
			return
		}

		// Validate category exists
		if config.GetCategory(categoryID) == nil {
			http.Error(w, "Categoria inválida", http.StatusBadRequest)
			return
		}

		// Redirect to step 2
		http.Redirect(w, r, "/report/category/"+categoryID, http.StatusSeeOther)
		return
	}
}

// ReportStep2Handler handles step 2: report details form
func ReportStep2Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID := vars["category"]

	// Validate category exists
	category := config.GetCategory(categoryID)
	if category == nil {
		http.NotFound(w, r)
		return
	}

	if r.Method == "GET" {
		// GET: Show form
		locationRequired := config.IsLocationRequiredGlobal(categoryID)
		maxFiles := models.GetMaxFiles(categoryID)
		allowedTypes := models.GetAllowedFileTypes(categoryID)

		data := map[string]interface{}{
			"Step":             2,
			"Category":         category,
			"LocationRequired": locationRequired,
			"PageTitle":        "Nova Denúncia - " + category.Name,
			"MaxFiles":         maxFiles,
			"AllowedTypes":     allowedTypes,
			"GoogleMapsAPIKey": getGoogleMapsAPIKey(),
		}

		if err := renderTemplate(w, "01_report.html", data); err != nil {
			log.Printf("Error rendering report step 2 template: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// POST: Process form submission
	if r.Method == "POST" {
		handleReportSubmission(w, r, category)
		return
	}
}

// handleReportSubmission processes the report form submission
func handleReportSubmission(w http.ResponseWriter, r *http.Request, category *config.Category) {
	// Parse multipart form (for file uploads)
	err := r.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	// Extract form data
	cpf := r.FormValue("cpf")
	birthDate := r.FormValue("birth_date")
	email := r.FormValue("email")
	emailConfirmation := r.FormValue("email_confirmation")
	location := r.FormValue("location")
	description := r.FormValue("description")

	// Parse coordinates
	latitude, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
	if err != nil {
		http.Error(w, "Latitude inválida", http.StatusBadRequest)
		return
	}

	longitude, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
	if err != nil {
		http.Error(w, "Longitude inválida", http.StatusBadRequest)
		return
	}

	// Validate form data
	validationErrors := services.ValidateForm(category.ID, cpf, birthDate, email, emailConfirmation, location, description, latitude, longitude)
	if len(validationErrors) > 0 {
		// Return to form with errors
		data := map[string]interface{}{
			"Step":             2,
			"Category":         category,
			"LocationRequired": config.IsLocationRequiredGlobal(category.ID),
			"PageTitle":        "Nova Denúncia - " + category.Name,
			"MaxFiles":         models.GetMaxFiles(category.ID),
			"AllowedTypes":     models.GetAllowedFileTypes(category.ID),
			"GoogleMapsAPIKey": getGoogleMapsAPIKey(),
			"Errors":           validationErrors,
			"FormData": map[string]interface{}{
				"CPF":         cpf,
				"BirthDate":   birthDate,
				"Email":       email,
				"Location":    location,
				"Description": description,
				"Latitude":    latitude,
				"Longitude":   longitude,
			},
		}

		if err := renderTemplate(w, "01_report.html", data); err != nil {
			log.Printf("Error rendering report template with errors: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Process file uploads
	var uploadedFiles []string
	files := r.MultipartForm.File["files"]

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		result, err := services.ProcessFileUpload(file, fileHeader, category.ID)
		if err != nil {
			log.Printf("Error uploading file %s: %v", fileHeader.Filename, err)
			continue
		}

		uploadedFiles = append(uploadedFiles, result.SavedPath)
	}

	// Create report record
	report := &models.Report{
		ProblemType: category.ID,
		HashedCPF:   services.HashCPF(cpf),
		BirthDate:   birthDate, // Store birth date (will be hashed in production)
		Email:       email,
		Location:    location,
		Latitude:    latitude,
		Longitude:   longitude,
		Description: description,
		PhotoPath:   strings.Join(uploadedFiles, ","), // Store multiple paths comma-separated
	}

	// Save to database
	reportID, err := services.CreateReport(db.DB, report)
	if err != nil {
		log.Printf("Error creating report: %v", err)
		http.Error(w, "Erro ao salvar denúncia", http.StatusInternalServerError)
		return
	}

	// Send confirmation email (async)
	go services.SendConfirmationEmail(email, reportID, category.Name)

	// Redirect to success page
	http.Redirect(w, r, fmt.Sprintf("/report/success/%d", reportID), http.StatusSeeOther)
}

// ReportSuccessHandler shows the success page after report submission
func ReportSuccessHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reportIDStr := vars["id"]

	reportID, err := strconv.Atoi(reportIDStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"ReportID":  reportID,
		"PageTitle": "Denúncia Enviada com Sucesso",
		"Success":   true,
		"Now":       time.Now(),
	}

	if err := renderTemplate(w, "01_report.html", data); err != nil {
		log.Printf("Error rendering success template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// ReportDetailHandler handles viewing a specific report by ID
func ReportDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reportIDStr := vars["id"]

	reportID, err := strconv.Atoi(reportIDStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Fetch report from database
	report, err := services.GetReportByID(db.DB, reportID)
	if err != nil {
		log.Printf("Error fetching report %d: %v", reportID, err)
		http.NotFound(w, r)
		return
	}

	// Get category info
	category := config.GetCategory(report.ProblemType)

	data := map[string]interface{}{
		"ReportID":  reportID,
		"Report":    report,
		"Category":  category,
		"PageTitle": "Denúncia #" + reportIDStr,
	}

	if err := renderTemplate(w, "01_report.html", data); err != nil {
		log.Printf("Error rendering report detail template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// GoogleMapsAPIHandler provides API endpoints for Google Maps integration
func GoogleMapsAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"apiKey": getGoogleMapsAPIKey(),
		"config": map[string]interface{}{
			"defaultCenter": map[string]float64{
				"lat": -23.5505, // São Paulo coordinates as default
				"lng": -46.6333,
			},
			"defaultZoom": 12,
		},
	}

	json.NewEncoder(w).Encode(response)
}

// Helper functions

func getGoogleMapsAPIKey() string {
	// Get from config
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Error loading config for Google Maps API key: %v", err)
		return ""
	}
	return cfg.GoogleMapsAPIKey
}
