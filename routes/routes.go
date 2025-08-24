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

	// SEO routes
	r.HandleFunc("/sitemap.xml", handlers.SitemapHandler).Methods("GET")
	r.HandleFunc("/robots.txt", handlers.RobotsHandler).Methods("GET")

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
	r.HandleFunc("/api/share-image", handlers.ShareImageHandler).Methods("POST")  // Share image generation
	r.HandleFunc("/api/stats", handlers.StatsHandler).Methods("GET")              // Statistics data

	// Comment routes
	r.HandleFunc("/api/comments", handlers.CreateCommentHandler).Methods("POST") // Create comment
	r.HandleFunc("/api/comments", handlers.GetCommentsHandler).Methods("GET")    // Get comments

	r.HandleFunc("/feed", handlers.FeedHandler).Methods("GET")
	r.HandleFunc("/map", handlers.MapHandler).Methods("GET")

	// Article routes
	r.HandleFunc("/articles", handlers.ArticlesHandler).Methods("GET")
	r.HandleFunc("/articles/{slug}", handlers.ArticleHandler).Methods("GET")

	// Footer pages routes
	r.HandleFunc("/sobre", handlers.SobreHandler).Methods("GET")
	r.HandleFunc("/status", handlers.StatusHandler).Methods("GET")
	r.HandleFunc("/termos", handlers.TermosHandler).Methods("GET")
	r.HandleFunc("/ajuda", handlers.AjudaHandler).Methods("GET")
	r.HandleFunc("/blog", handlers.BlogHandler).Methods("GET")
	r.HandleFunc("/governos", handlers.GovernosHandler).Methods("GET")
	r.HandleFunc("/empresas", handlers.EmpresasHandler).Methods("GET")
	r.HandleFunc("/pesquisadores", handlers.PesquisadoresHandler).Methods("GET")
	r.HandleFunc("/transparencia", handlers.TransparenciaHandler).Methods("GET")

	return r
}
