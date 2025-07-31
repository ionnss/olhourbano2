package handlers

import (
	"log"
	"net/http"
	"olhourbano2/config"
	"strconv"
)

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

	// Get categories for filter dropdown
	categories := config.GetAllCategories()

	// Here you would fetch filtered and paginated reports
	data := map[string]interface{}{
		"Page":       page,
		"Category":   category,
		"Status":     status,
		"Categories": categories,
		"PageTitle":  "Feed de Denúncias",
		// "Reports": getReports(page, category, status),
		// "TotalPages": getTotalPages(category, status),
		// "HasNext": hasNextPage(page, category, status),
		// "HasPrev": page > 1,
	}

	if err := renderTemplate(w, "02_feed.html", data); err != nil {
		log.Printf("Error rendering feed template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
