package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
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
		transportRequired := config.IsTransportRequiredGlobal(categoryID)
		transportTypes := config.GetTransportTypesGlobal()

		data := map[string]interface{}{
			"Step":              2,
			"Category":          category,
			"LocationRequired":  locationRequired,
			"TransportRequired": transportRequired,
			"TransportTypes":    transportTypes,
			"PageTitle":         "Nova Denúncia - " + category.Name,
			"MaxFiles":          maxFiles,
			"AllowedTypes":      allowedTypes,
			"GoogleMapsAPIKey":  getGoogleMapsAPIKey(),
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
	birthDateInput := r.FormValue("birth_date")
	email := r.FormValue("email")

	// Convert birth date format
	birthDate, err := services.ConvertBirthDateToDBFormat(birthDateInput)
	if err != nil {
		http.Error(w, "Data de nascimento: "+err.Error(), http.StatusBadRequest)
		return
	}
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

	// Extract transport data if required
	var transportType string
	var transportData *models.TransportData

	if config.IsTransportRequiredGlobal(category.ID) {
		transportType = r.FormValue("transport_type")

		// Only create transport data if transport type is selected
		if transportType != "" {
			transportData = &models.TransportData{}

			switch transportType {
			case "bus":
				transportData.BusNumber = r.FormValue("bus_number")
				transportData.BusLine = r.FormValue("bus_line")
				transportData.BusStop = r.FormValue("bus_stop")
				transportData.BusCompany = r.FormValue("bus_company")
			case "metro":
				transportData.MetroLine = r.FormValue("metro_line")
				transportData.MetroStation = r.FormValue("metro_station")
				transportData.MetroWagon = r.FormValue("metro_wagon")
				transportData.MetroCard = r.FormValue("metro_card")
			case "train":
				transportData.TrainLine = r.FormValue("train_line")
				transportData.TrainStation = r.FormValue("train_station")
				transportData.TrainWagon = r.FormValue("train_wagon")
			case "other":
				transportData.TransportDetails = r.FormValue("transport_details")
			}

			// Check if any transport data was actually provided
			hasData := false
			if transportData.BusNumber != "" || transportData.BusLine != "" || transportData.BusStop != "" || transportData.BusCompany != "" ||
				transportData.MetroLine != "" || transportData.MetroStation != "" || transportData.MetroWagon != "" || transportData.MetroCard != "" ||
				transportData.TrainLine != "" || transportData.TrainStation != "" || transportData.TrainWagon != "" ||
				transportData.TransportDetails != "" {
				hasData = true
			}

			// If no meaningful data was provided, set transportData to nil
			if !hasData {
				transportData = nil
			}
		}
	}

	// Validate form data
	validationErrors := services.ValidateForm(category.ID, cpf, birthDate, email, emailConfirmation, location, description, latitude, longitude)

	// Validate file uploads - check if at least one file is provided
	files := r.MultipartForm.File["files"]
	if files == nil {
		files = make([]*multipart.FileHeader, 0)
	}
	fileValidationErrors := services.ValidateFiles(len(files))
	validationErrors = append(validationErrors, fileValidationErrors...)

	if len(validationErrors) > 0 {
		// Return to form with errors
		data := map[string]interface{}{
			"Step":              2,
			"Category":          category,
			"LocationRequired":  config.IsLocationRequiredGlobal(category.ID),
			"TransportRequired": config.IsTransportRequiredGlobal(category.ID),
			"TransportTypes":    config.GetTransportTypesGlobal(),
			"PageTitle":         "Nova Denúncia - " + category.Name,
			"MaxFiles":          models.GetMaxFiles(category.ID),
			"AllowedTypes":      models.GetAllowedFileTypes(category.ID),
			"GoogleMapsAPIKey":  getGoogleMapsAPIKey(),
			"Errors":            validationErrors,
			"FormData": map[string]interface{}{
				"CPF":           cpf,
				"BirthDate":     birthDate,
				"Email":         email,
				"Location":      location,
				"Description":   description,
				"Latitude":      latitude,
				"Longitude":     longitude,
				"TransportType": transportType,
				"TransportData": transportData,
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
		ProblemType:   category.ID,
		HashedCPF:     services.HashCPF(cpf),
		BirthDate:     birthDate, // Store birth date (will be hashed in production)
		Email:         email,
		Location:      location,
		Latitude:      latitude,
		Longitude:     longitude,
		Description:   description,
		PhotoPath:     strings.Join(uploadedFiles, ","), // Store multiple paths comma-separated
		TransportType: transportType,
	}

	// Set transport data if available
	if transportData != nil {
		err = report.SetTransportData(transportData)
		if err != nil {
			log.Printf("Error setting transport data: %v", err)
			// Continue without transport data rather than failing
		}
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

	// Process photos
	var photos []string
	if report.PhotoPath != "" {
		// Split by comma and trim whitespace
		rawPhotos := strings.Split(report.PhotoPath, ",")
		for _, photo := range rawPhotos {
			trimmed := strings.TrimSpace(photo)
			if trimmed != "" {
				photos = append(photos, trimmed)
			}
		}
	}

	// Get first 8 characters of hashed CPF for display
	hashedCPFDisplay := ""
	log.Printf("DEBUG: Report %d - Raw HashedCPF: '%s', Length: %d", reportID, report.HashedCPF, len(report.HashedCPF))

	if len(report.HashedCPF) >= 8 {
		hashedCPFDisplay = report.HashedCPF[:8]
		log.Printf("DEBUG: Using first 8 chars: '%s'", hashedCPFDisplay)
	} else if report.HashedCPF != "" {
		// If hashed CPF is shorter than 8 characters, use what's available
		hashedCPFDisplay = report.HashedCPF
		log.Printf("DEBUG: Using full hash (short): '%s'", hashedCPFDisplay)
	} else {
		log.Printf("DEBUG: HashedCPF is empty, will show [Anônimo]")
	}

	// Get Google Maps API key
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Error loading config for Google Maps API key: %v", err)
	}

	// Get initial comments for the report
	comments, err := services.GetCommentsForReport(db.DB, reportID, "recent", 10, 0)
	if err != nil {
		log.Printf("Error fetching comments for report %d: %v", reportID, err)
		// Continue without comments
		comments = []*models.CommentDisplay{}
	}

	// Process status text for display
	statusText := ""
	switch report.Status {
	case "pending":
		statusText = "Pendente"
	case "approved":
		statusText = "Resolvida"
	default:
		statusText = "Pendente"
	}

	// Process transport details for display
	transportDetails := ""
	transportTypeName := ""
	if report.TransportType != "" && report.TransportData != nil {
		// Get human-readable transport type name
		switch report.TransportType {
		case "bus":
			transportTypeName = "Ônibus"
		case "metro":
			transportTypeName = "Metrô"
		case "train":
			transportTypeName = "Trem"
		case "other":
			transportTypeName = "Outro"
		default:
			transportTypeName = "Transporte"
		}

		transportData, err := report.GetTransportData()
		if err != nil {
			log.Printf("Error parsing transport data for report %d: %v", reportID, err)
		} else if transportData != nil {
			var details []string

			// Format bus details
			if transportData.BusNumber != "" {
				details = append(details, "Ônibus: "+transportData.BusNumber)
			}
			if transportData.BusLine != "" {
				details = append(details, "Linha: "+transportData.BusLine)
			}
			if transportData.BusStop != "" {
				details = append(details, "Ponto: "+transportData.BusStop)
			}
			if transportData.BusCompany != "" {
				details = append(details, "Empresa: "+transportData.BusCompany)
			}

			// Format metro details
			if transportData.MetroLine != "" {
				details = append(details, "Linha: "+transportData.MetroLine)
			}
			if transportData.MetroStation != "" {
				details = append(details, "Estação: "+transportData.MetroStation)
			}
			if transportData.MetroWagon != "" {
				details = append(details, "Vagão: "+transportData.MetroWagon)
			}
			if transportData.MetroCard != "" {
				details = append(details, "Cartão: "+transportData.MetroCard)
			}

			// Format train details
			if transportData.TrainLine != "" {
				details = append(details, "Linha: "+transportData.TrainLine)
			}
			if transportData.TrainStation != "" {
				details = append(details, "Estação: "+transportData.TrainStation)
			}
			if transportData.TrainWagon != "" {
				details = append(details, "Vagão: "+transportData.TrainWagon)
			}

			// Format other transport details
			if transportData.TransportDetails != "" {
				details = append(details, transportData.TransportDetails)
			}

			transportDetails = strings.Join(details, " • ")
		}
	}

	data := map[string]interface{}{
		"ReportID":          reportID,
		"Report":            report,
		"Category":          category,
		"Photos":            photos,
		"HashedCPFDisplay":  hashedCPFDisplay,
		"PageTitle":         "Denúncia #" + reportIDStr,
		"GoogleMapsAPIKey":  cfg.GoogleMapsAPIKey,
		"Comments":          comments,
		"TransportDetails":  transportDetails,
		"TransportTypeName": transportTypeName,
		"StatusText":        statusText,
	}

	if err := renderTemplate(w, "04_report_detail.html", data); err != nil {
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
