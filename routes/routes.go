package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func CreateRoutes() *mux.Router {
	r := mux.NewRouter()

	// Initialize middlewares

	// Serve static files
	fileServer := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// Page routes

	return r
}
