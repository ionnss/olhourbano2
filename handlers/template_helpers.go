package handlers

import (
	"olhourbano2/services"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateFuncs returns template functions
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
		"isImageFile": func(filename string) bool {
			ext := strings.ToLower(filepath.Ext(filename))
			imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp"}
			for _, imgExt := range imageExtensions {
				if ext == imgExt {
					return true
				}
			}
			return false
		},
		"isVideoFile": func(filename string) bool {
			ext := strings.ToLower(filepath.Ext(filename))
			videoExtensions := []string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm"}
			for _, vidExt := range videoExtensions {
				if ext == vidExt {
					return true
				}
			}
			return false
		},
		"isPdfFile": func(filename string) bool {
			return strings.ToLower(filepath.Ext(filename)) == ".pdf"
		},
		"getFileType": func(filename string) string {
			ext := strings.ToLower(filepath.Ext(filename))
			imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp"}
			videoExtensions := []string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm"}

			for _, imgExt := range imageExtensions {
				if ext == imgExt {
					return "image"
				}
			}
			for _, vidExt := range videoExtensions {
				if ext == vidExt {
					return "video"
				}
			}
			if ext == ".pdf" {
				return "pdf"
			}
			return "document"
		},
		"getThumbnailPath": func(filePath string) string {
			return services.GetThumbnailPath(filePath)
		},
		"getThumbnailFilename": func(filePath string) string {
			if filePath == "" {
				return ""
			}

			// Extract hash from filename
			filename := filepath.Base(filePath)
			ext := filepath.Ext(filename)
			hash := strings.TrimSuffix(filename, ext)

			return hash + "_thumb.jpg"
		},
		"getFileTypeIcon": func(filename string) string {
			ext := strings.ToLower(filepath.Ext(filename))

			// Map extensions to content types for icon detection
			contentType := ""
			switch ext {
			case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp":
				contentType = "image/jpeg"
			case ".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm":
				contentType = "video/mp4"
			case ".pdf":
				contentType = "application/pdf"
			case ".txt":
				contentType = "text/plain"
			case ".doc":
				contentType = "application/msword"
			case ".docx":
				contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
			default:
				contentType = "application/octet-stream"
			}

			return services.GetFileTypeIcon(contentType)
		},
	}
}
