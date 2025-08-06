package routes

import (
	"net/http"
	"olhourbano2/handlers"

	"github.com/gorilla/mux"
)

func CreateRoutes() *mux.Router {
	r := mux.NewRouter()

	// Initialize middlewares

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Serve static files
	fileServer := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// Serve uploaded files
	uploadServer := http.FileServer(http.Dir("./uploads/"))
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", uploadServer))

	// Serve templates
	fileServer = http.FileServer(http.Dir("./templates/pages/"))
	r.PathPrefix("/templates/pages/").Handler(http.StripPrefix("/templates/pages/", fileServer))

	fileServer = http.FileServer(http.Dir("./templates/components/"))
	r.PathPrefix("/templates/components/").Handler(http.StripPrefix("/templates/components/", fileServer))

	// Page routes
	r.HandleFunc("/", handlers.IndexHandler).Methods("GET")

	// Report routes - two-step process
	r.HandleFunc("/report", handlers.ReportHandler).Methods("GET", "POST")                          // Step 1: Category selection
	r.HandleFunc("/report/category/{category}", handlers.ReportStep2Handler).Methods("GET", "POST") // Step 2: Report details
	r.HandleFunc("/report/success/{id:[0-9]+}", handlers.ReportSuccessHandler).Methods("GET")       // Success page
	r.HandleFunc("/report/{id:[0-9]+}", handlers.ReportDetailHandler).Methods("GET")                // View existing report

	// API routes
	r.HandleFunc("/api/googlemaps", handlers.GoogleMapsAPIHandler).Methods("GET") // Google Maps config
	r.HandleFunc("/api/verify-cpf", handlers.VerifyCPFHandler).Methods("POST")    // CPF verification
	r.HandleFunc("/api/reports/map", handlers.MapReportsHandler).Methods("GET")   // Map reports data
	r.HandleFunc("/api/reports/cities", handlers.CitiesHandler).Methods("GET")    // Cities data
	r.HandleFunc("/api/vote", handlers.VoteHandler).Methods("POST")               // Vote on reports

	r.HandleFunc("/feed", handlers.FeedHandler).Methods("GET")
	r.HandleFunc("/map", handlers.MapHandler).Methods("GET")

	return r
}
