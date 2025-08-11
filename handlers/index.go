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

// Footer page handlers
func SobreHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle": "Sobre Nós - Olho Urbano",
	}

	if err := renderTemplate(w, "footer_pages/sobre.html", data); err != nil {
		log.Printf("Error rendering sobre template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle": "Status do Sistema - Olho Urbano",
	}

	if err := renderTemplate(w, "footer_pages/status.html", data); err != nil {
		log.Printf("Error rendering status template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func TermosHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle": "Termos e Privacidade - Olho Urbano",
	}

	if err := renderTemplate(w, "footer_pages/termos.html", data); err != nil {
		log.Printf("Error rendering termos template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func AjudaHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle": "Ajuda - Olho Urbano",
	}

	if err := renderTemplate(w, "footer_pages/ajuda.html", data); err != nil {
		log.Printf("Error rendering ajuda template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func BlogHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle": "Blog - Olho Urbano",
	}

	if err := renderTemplate(w, "footer_pages/blog.html", data); err != nil {
		log.Printf("Error rendering blog template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func GovernosHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle": "Prefeituras e Governos - Olho Urbano",
	}

	if err := renderTemplate(w, "footer_pages/governos.html", data); err != nil {
		log.Printf("Error rendering governos template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func EmpresasHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle": "Empresas - Olho Urbano",
	}

	if err := renderTemplate(w, "footer_pages/empresas.html", data); err != nil {
		log.Printf("Error rendering empresas template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func PesquisadoresHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle": "Universidades e Pesquisadores - Olho Urbano",
	}

	if err := renderTemplate(w, "footer_pages/pesquisadores.html", data); err != nil {
		log.Printf("Error rendering pesquisadores template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
