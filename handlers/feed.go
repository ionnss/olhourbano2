package handlers

import (
	"log"
	"net/http"
	"olhourbano2/config"
	"olhourbano2/db"
	"olhourbano2/models"
	"olhourbano2/services"
	"strconv"
	"strings"
)

const ReportsPerPage = 9

// FeedHandler handles the reports feed page with filtering and pagination
func FeedHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering and pagination
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	category := r.URL.Query().Get("category")
	status := r.URL.Query().Get("status")
	city := r.URL.Query().Get("city")
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "recent" // default sort
	}

	// Get categories for filter dropdown
	categories := config.GetAllCategories()

	// Get cities for filter dropdown
	cities, err := services.GetCitiesFromReports(db.DB)
	if err != nil {
		log.Printf("Error fetching cities: %v", err)
		cities = []string{} // Continue with empty cities list
	}

	// Fetch reports from database
	reports, err := services.GetReports(db.DB, page, category, status, city, sort, ReportsPerPage)
	if err != nil {
		log.Printf("Error fetching reports: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get total count for pagination
	totalReports, err := services.GetTotalReports(db.DB, category, status, city)
	if err != nil {
		log.Printf("Error getting total reports: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Calculate pagination
	totalPages := (totalReports + ReportsPerPage - 1) / ReportsPerPage
	hasNext := page < totalPages
	hasPrev := page > 1
	prevPage := page - 1
	nextPage := page + 1

	// Get statistics
	stats, err := services.GetReportStats(db.DB)
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		// Continue without stats rather than failing
		stats = &services.ReportStats{
			TotalReports: totalReports,
			ThisMonth:    0,
			Resolved:     0,
		}
	}

	// Process reports for template
	processedReports := processReportsForTemplate(reports)

	// Calculate additional statistics
	pendingReports := totalReports - stats.Resolved
	resolutionRate := 0
	if totalReports > 0 {
		resolutionRate = int(float64(stats.Resolved) / float64(totalReports) * 100)
	}

	data := map[string]interface{}{
		"Page":           page,
		"Category":       category,
		"Status":         status,
		"City":           city,
		"Sort":           sort,
		"Categories":     categories,
		"Cities":         cities,
		"Reports":        processedReports,
		"TotalReports":   totalReports,
		"TotalPages":     totalPages,
		"HasNext":        hasNext,
		"HasPrev":        hasPrev,
		"PrevPage":       prevPage,
		"NextPage":       nextPage,
		"PageTitle":      "Feed de Denúncias",
		"ThisMonth":      stats.ThisMonth,
		"Resolved":       stats.Resolved,
		"Pending":        pendingReports,
		"ResolutionRate": resolutionRate,
		"CurrentPage":    "feed",
	}

	if err := renderTemplate(w, "02_feed.html", data); err != nil {
		log.Printf("Error rendering feed template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// processReportsForTemplate converts database reports to template-friendly format
func processReportsForTemplate(reports []*models.Report) []map[string]interface{} {
	var processed []map[string]interface{}

	for _, report := range reports {
		// Get category info
		category := config.GetCategory(report.ProblemType)
		categoryIcon := "❓"
		categoryName := "Desconhecida"
		if category != nil {
			categoryIcon = category.Icon
			categoryName = category.Name
		}

		// Get status text
		statusText := getStatusText(report.Status)

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

		// Process transport info
		transportTypeName := ""
		transportDetails := ""
		if report.TransportType != "" {
			transportTypeName = getTransportTypeName(report.TransportType)
			transportDetails = getTransportDetails(report)
		}

		// Format date
		createdAt := report.CreatedAt.Format("02/01/2006 às 15:04")

		// Get first 8 characters of hashed CPF for display
		hashedCPFDisplay := ""
		if len(report.HashedCPF) >= 8 {
			hashedCPFDisplay = report.HashedCPF[:8]
		}

		processed = append(processed, map[string]interface{}{
			"ID":                report.ID,
			"CategoryIcon":      categoryIcon,
			"CategoryName":      categoryName,
			"Status":            report.Status,
			"StatusText":        statusText,
			"Location":          report.Location,
			"Description":       report.Description,
			"PhotoPath":         report.PhotoPath,
			"Photos":            photos,
			"TransportType":     report.TransportType,
			"TransportTypeName": transportTypeName,
			"TransportDetails":  transportDetails,
			"CreatedAt":         createdAt,
			"VoteCount":         report.VoteCount,
			"HashedCPFDisplay":  hashedCPFDisplay,
		})
	}

	return processed
}

// getStatusText returns human-readable status text
func getStatusText(status string) string {
	switch status {
	case "pending":
		return "Pendente"
	case "approved":
		return "Resolvida"
	default:
		return "Pendente"
	}
}

// getTransportTypeName returns human-readable transport type name
func getTransportTypeName(transportType string) string {
	switch transportType {
	case "bus":
		return "Ônibus"
	case "metro":
		return "Metrô"
	case "train":
		return "Trem"
	case "other":
		return "Outro"
	default:
		return "Transporte"
	}
}

// getTransportDetails returns formatted transport details
func getTransportDetails(report *models.Report) string {
	// This would need to be implemented to parse the JSON transport data
	// For now, return empty string
	return ""
}
