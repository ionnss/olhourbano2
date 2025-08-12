package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

// SitemapURL represents a URL in the sitemap
type SitemapURL struct {
	Loc        string `xml:"loc"`
	Lastmod    string `xml:"lastmod"`
	Changefreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

// Sitemap represents the sitemap structure
type Sitemap struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

// SitemapHandler generates and serves the sitemap.xml
func SitemapHandler(w http.ResponseWriter, r *http.Request) {
	baseURL := "https://olhourbano.com.br"
	currentTime := time.Now().Format("2006-01-02")

	sitemap := Sitemap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs: []SitemapURL{
			{
				Loc:        baseURL + "/",
				Lastmod:    currentTime,
				Changefreq: "daily",
				Priority:   "1.0",
			},
			{
				Loc:        baseURL + "/report",
				Lastmod:    currentTime,
				Changefreq: "weekly",
				Priority:   "0.9",
			},
			{
				Loc:        baseURL + "/feed",
				Lastmod:    currentTime,
				Changefreq: "hourly",
				Priority:   "0.8",
			},
			{
				Loc:        baseURL + "/map",
				Lastmod:    currentTime,
				Changefreq: "daily",
				Priority:   "0.8",
			},
			{
				Loc:        baseURL + "/sobre",
				Lastmod:    currentTime,
				Changefreq: "monthly",
				Priority:   "0.6",
			},
			{
				Loc:        baseURL + "/ajuda",
				Lastmod:    currentTime,
				Changefreq: "monthly",
				Priority:   "0.6",
			},
			{
				Loc:        baseURL + "/termos",
				Lastmod:    currentTime,
				Changefreq: "monthly",
				Priority:   "0.5",
			},
			{
				Loc:        baseURL + "/status",
				Lastmod:    currentTime,
				Changefreq: "weekly",
				Priority:   "0.7",
			},
			{
				Loc:        baseURL + "/blog",
				Lastmod:    currentTime,
				Changefreq: "weekly",
				Priority:   "0.7",
			},
			{
				Loc:        baseURL + "/governos",
				Lastmod:    currentTime,
				Changefreq: "monthly",
				Priority:   "0.6",
			},
			{
				Loc:        baseURL + "/empresas",
				Lastmod:    currentTime,
				Changefreq: "monthly",
				Priority:   "0.6",
			},
			{
				Loc:        baseURL + "/pesquisadores",
				Lastmod:    currentTime,
				Changefreq: "monthly",
				Priority:   "0.6",
			},
		},
	}

	// Add recent reports to sitemap (you can enhance this by querying your database)
	// For now, we'll add a few example report URLs
	reports := []string{"123", "456", "789"} // Replace with actual report IDs from your database
	for _, reportID := range reports {
		sitemap.URLs = append(sitemap.URLs, SitemapURL{
			Loc:        fmt.Sprintf("%s/report/%s", baseURL, reportID),
			Lastmod:    currentTime,
			Changefreq: "weekly",
			Priority:   "0.7",
		})
	}

	// Set content type
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour

	// Encode and write the sitemap
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	if err := encoder.Encode(sitemap); err != nil {
		http.Error(w, "Error generating sitemap", http.StatusInternalServerError)
		return
	}
}
