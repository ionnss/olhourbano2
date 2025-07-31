package handlers

import (
	"log"
	"net/http"
	"olhourbano2/config"
)

// IndexHandler handles the main index/home page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Get categories for display on index page
	categories := config.GetAllCategories()

	data := map[string]interface{}{
		"Categories": categories,
		"PageTitle":  "Olho Urbano - Denúncias Públicas",
		// "FeaturedReports": getFeaturedReports(),
		// "Statistics": getStatistics(),
	}

	if err := renderTemplate(w, "00_index.html", data); err != nil {
		log.Printf("Error rendering index template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
