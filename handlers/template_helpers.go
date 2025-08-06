package handlers

import (
	"html/template"
	"strings"
)

// TemplateFuncs returns the template functions map
func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"subtract": func(a, b int) int {
			return a - b
		},
		"join": func(slice []string, sep string) string {
			return strings.Join(slice, sep)
		},
	}
}
