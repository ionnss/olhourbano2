package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

// renderTemplate is a shared utility function for rendering page templates
// It parses all component templates first, then the specific page template
func renderTemplate(w http.ResponseWriter, pageName string, data interface{}) error {
	// Create template with functions
	tmpl := template.New("").Funcs(TemplateFuncs())

	// Parse all component templates first
	tmpl, err := tmpl.ParseGlob("./templates/components/*.html")
	if err != nil {
		return fmt.Errorf("error parsing component templates: %w", err)
	}

	// Parse all footer page templates (needed for 05_footer_pages.html)
	tmpl, err = tmpl.ParseGlob("./templates/footer_pages/*.html")
	if err != nil {
		return fmt.Errorf("error parsing footer page templates: %w", err)
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

// renderFooterTemplate is a utility function for rendering footer page templates
// It parses all component templates first, then the specific footer page template
func renderFooterTemplate(w http.ResponseWriter, pageName string, data interface{}) error {
	// Create template with functions
	tmpl := template.New("").Funcs(TemplateFuncs())

	// Parse all component templates first
	tmpl, err := tmpl.ParseGlob("./templates/components/*.html")
	if err != nil {
		return fmt.Errorf("error parsing component templates: %w", err)
	}

	// Parse the specific footer page template
	tmpl, err = tmpl.ParseFiles(fmt.Sprintf("./templates/footer_pages/%s", pageName))
	if err != nil {
		return fmt.Errorf("error parsing footer page template %s: %w", pageName, err)
	}

	// Execute the template - use the filename as the template name (with .html extension)
	err = tmpl.ExecuteTemplate(w, pageName, data)
	if err != nil {
		return fmt.Errorf("error executing footer page template %s: %w", pageName, err)
	}

	return nil
}
