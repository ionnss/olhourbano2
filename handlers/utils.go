package handlers

import (
	"fmt"
	"net/http"
	"text/template"
)

// renderTemplate is a shared utility function for rendering page templates
// It parses all component templates first, then the specific page template
func renderTemplate(w http.ResponseWriter, pageName string, data interface{}) error {
	// Parse all component templates first
	tmpl, err := template.ParseGlob("./templates/components/*.html")
	if err != nil {
		return fmt.Errorf("error parsing component templates: %w", err)
	}

	// Parse the specific page template
	tmpl, err = tmpl.ParseFiles(fmt.Sprintf("./templates/pages/%s", pageName))
	if err != nil {
		return fmt.Errorf("error parsing page template %s: %w", pageName, err)
	}

	// Execute the template
	err = tmpl.ExecuteTemplate(w, pageName, data)
	if err != nil {
		return fmt.Errorf("error executing template %s: %w", pageName, err)
	}

	return nil
}
