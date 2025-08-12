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
		"PageTitle":    "Sobre Nós",
		"PageSubtitle": "Conheça mais sobre o Olho Urbano",
		"Content":      "sobre_content",
	}

	if err := renderTemplate(w, "05_footer_pages.html", data); err != nil {
		log.Printf("Error rendering sobre template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle":    "Status do Sistema",
		"PageSubtitle": "Acompanhe o status dos nossos serviços",
		"Content":      "status_content",
	}

	if err := renderTemplate(w, "05_footer_pages.html", data); err != nil {
		log.Printf("Error rendering status template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func TermosHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle":    "Termos e Privacidade",
		"PageSubtitle": "Nossos termos de uso e política de privacidade",
		"Content":      "termos_content",
	}

	if err := renderTemplate(w, "05_footer_pages.html", data); err != nil {
		log.Printf("Error rendering termos template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func AjudaHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle":    "Ajuda",
		"PageSubtitle": "Encontre respostas para suas dúvidas",
		"Content":      "ajuda_content",
	}

	if err := renderTemplate(w, "05_footer_pages.html", data); err != nil {
		log.Printf("Error rendering ajuda template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func BlogHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := loadArticles()
	if err != nil {
		log.Printf("Error loading articles: %v", err)
		http.Error(w, "Error loading articles", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"PageTitle":    "Blog Olho Urbano",
		"PageSubtitle": "Artigos, insights e histórias sobre cidades inteligentes",
		"Content":      "blog_content",
		"Articles":     articles,
		"Total":        len(articles),
	}

	if err := renderTemplate(w, "05_footer_pages.html", data); err != nil {
		log.Printf("Error rendering blog template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func GovernosHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle":    "Prefeituras e Governos",
		"PageSubtitle": "Soluções para gestão pública",
		"Content":      "governos_content",
	}

	if err := renderTemplate(w, "05_footer_pages.html", data); err != nil {
		log.Printf("Error rendering governos template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func EmpresasHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle":    "Empresas",
		"PageSubtitle": "Soluções para o setor privado",
		"Content":      "empresas_content",
	}

	if err := renderTemplate(w, "05_footer_pages.html", data); err != nil {
		log.Printf("Error rendering empresas template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func PesquisadoresHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle":    "Universidades e Pesquisadores",
		"PageSubtitle": "Ferramentas para pesquisa e desenvolvimento",
		"Content":      "pesquisadores_content",
	}

	if err := renderTemplate(w, "05_footer_pages.html", data); err != nil {
		log.Printf("Error rendering pesquisadores template: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
