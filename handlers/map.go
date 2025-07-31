package handlers

import (
	"log"
	"net/http"
	"olhourbano2/config"
)

// MapHandler handles the interactive map page showing reports by location
func MapHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for map filtering
	category := r.URL.Query().Get("category")
	status := r.URL.Query().Get("status")

	// Get categories for filter dropdown
	categories := config.GetAllCategories()

	// Here you would fetch reports with location data
	data := map[string]interface{}{
		"Category":   category,
		"Status":     status,
		"Categories": categories,
		"PageTitle":  "Mapa de Denúncias",
		// "Reports": getReportsWithLocation(category, status),
		// "MapCenter": getDefaultMapCenter(),
		// "MapZoom": getDefaultMapZoom(),
	}

	if err := renderTemplate(w, "03_map.html", data); err != nil {
		log.Printf("Error rendering map template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
