package handlers

import (
	"fmt"
	"net/http"
)

// RobotsHandler serves the robots.txt file
func RobotsHandler(w http.ResponseWriter, r *http.Request) {
	robotsContent := `User-agent: *
Allow: /
Allow: /feed
Allow: /map
Allow: /sobre
Allow: /ajuda
Allow: /blog
Allow: /governos
Allow: /empresas
Allow: /pesquisadores
Allow: /status
Allow: /termos
Allow: /articles

# Disallow admin and API endpoints
Disallow: /api/
Disallow: /admin/
Disallow: /uploads/
Disallow: /templates/

# Sitemap location
Sitemap: https://olhourbano.com.br/sitemap.xml

# Crawl delay (optional - be respectful to the server)
Crawl-delay: 1`

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 24 hours
	fmt.Fprint(w, robotsContent)
}
