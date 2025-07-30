package routes

import (
	"net/http"

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

	// Serve templates
	fileServer = http.FileServer(http.Dir("./templates/"))
	r.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", fileServer))

	// Page routes

	return r
}
