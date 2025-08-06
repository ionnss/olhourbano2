package handlers

import (
	"log"
	"net/http"
	"olhourbano2/config"
	"olhourbano2/db"
	"olhourbano2/services"
)

// MapHandler handles the interactive map page showing reports by location
func MapHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for map filtering
	category := r.URL.Query().Get("category")
	status := r.URL.Query().Get("status")
	city := r.URL.Query().Get("city")
	
	// Get categories for filter dropdown
	categories := config.GetAllCategories()
	
	// Get cities from reports
	cities, err := services.GetCitiesFromReports(db.DB)
	if err != nil {
		log.Printf("Error fetching cities: %v", err)
		cities = []string{}
	}
	
	// Get statistics
	stats, err := services.GetReportStats(db.DB)
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		// Continue without stats rather than failing
		stats = &services.ReportStats{
			TotalReports: 0,
			ThisMonth:    0,
			Resolved:     0,
		}
	}
	
	// Calculate additional statistics
	totalReports := stats.TotalReports
	pendingReports := totalReports - stats.Resolved
	resolutionRate := 0
	if totalReports > 0 {
		resolutionRate = int(float64(stats.Resolved) / float64(totalReports) * 100)
	}
	
	data := map[string]interface{}{
		"Category":         category,
		"Status":           status,
		"City":             city,
		"Categories":       categories,
		"Cities":           cities,
		"PageTitle":        "Mapa de Denúncias",
		"GoogleMapsAPIKey": getGoogleMapsAPIKey(),
		"TotalReports":     totalReports,
		"ThisMonth":        stats.ThisMonth,
		"Resolved":         stats.Resolved,
		"Pending":          pendingReports,
		"ResolutionRate":   resolutionRate,
	}
	
	if err := renderTemplate(w, "03_map.html", data); err != nil {
		log.Printf("Error rendering map template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
