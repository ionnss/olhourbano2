package services

import (
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB
	UploadDir   = "./uploads"
)

// FileUploadResult represents the result of a file upload
type FileUploadResult struct {
	OriginalName string
	SavedPath    string
	FileSize     int64
	ContentType  string
	Error        error
}

// ProcessFileUpload handles file upload with metadata cleaning
func ProcessFileUpload(file multipart.File, header *multipart.FileHeader, category string) (*FileUploadResult, error) {
	result := &FileUploadResult{
		OriginalName: header.Filename,
		FileSize:     header.Size,
		ContentType:  header.Header.Get("Content-Type"),
	}

	// Validate file size
	if header.Size > MaxFileSize {
		result.Error = fmt.Errorf("arquivo muito grande. Máximo permitido: 10MB")
		return result, result.Error
	}

	// Validate file type
	if !isFileTypeAllowed(category, result.ContentType) {
		result.Error = fmt.Errorf("tipo de arquivo não permitido para esta categoria")
		return result, result.Error
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	hash := generateFileHash(header.Filename, time.Now().String())
	filename := fmt.Sprintf("%s%s", hash, ext)

	// Ensure upload directory exists
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		result.Error = fmt.Errorf("erro ao criar diretório de upload: %v", err)
		return result, result.Error
	}

	// Create temporary file path
	tempPath := filepath.Join(UploadDir, "temp_"+filename)
	finalPath := filepath.Join(UploadDir, filename)

	// Save file temporarily
	tempFile, err := os.Create(tempPath)
	if err != nil {
		result.Error = fmt.Errorf("erro ao criar arquivo temporário: %v", err)
		return result, result.Error
	}
	defer tempFile.Close()

	// Copy file content
	_, err = io.Copy(tempFile, file)
	if err != nil {
		os.Remove(tempPath)
		result.Error = fmt.Errorf("erro ao salvar arquivo: %v", err)
		return result, result.Error
	}

	// Clean metadata
	err = cleanFileMetadata(tempPath, finalPath, result.ContentType)
	if err != nil {
		os.Remove(tempPath)
		result.Error = fmt.Errorf("erro ao limpar metadados: %v", err)
		return result, result.Error
	}

	// Remove temporary file
	os.Remove(tempPath)

	result.SavedPath = finalPath
	return result, nil
}

// isFileTypeAllowed checks if file type is allowed (using our models function)
func isFileTypeAllowed(category, contentType string) bool {
	// Import from models package would be better, but for simplicity:
	allowedTypes := map[string][]string{
		"default": {
			"image/jpeg", "image/jpg", "image/png", "image/webp",
		},
		"corrupcao_gestao_publica": {
			"image/jpeg", "image/jpg", "image/png", "image/webp",
			"application/pdf", "text/plain",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		},
		"outros": {
			"image/jpeg", "image/jpg", "image/png", "image/webp",
			"application/pdf", "text/plain",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		},
	}

	categoryTypes, exists := allowedTypes[category]
	if !exists {
		categoryTypes = allowedTypes["default"]
	}

	for _, allowedType := range categoryTypes {
		if allowedType == contentType {
			return true
		}
	}
	return false
}

// generateFileHash generates a unique hash for the filename
func generateFileHash(filename, timestamp string) string {
	hash := sha256.Sum256([]byte(filename + timestamp))
	return fmt.Sprintf("%x", hash)[:16]
}

// cleanFileMetadata removes metadata from uploaded files
func cleanFileMetadata(inputPath, outputPath, contentType string) error {
	switch {
	case strings.HasPrefix(contentType, "image/"):
		return cleanImageMetadata(inputPath, outputPath)
	case contentType == "application/pdf":
		return cleanPDFMetadata(inputPath, outputPath)
	default:
		// For other file types, just copy the file
		return copyFile(inputPath, outputPath)
	}
}

// cleanImageMetadata removes EXIF data from images using ImageMagick
func cleanImageMetadata(inputPath, outputPath string) error {
	// Check if ImageMagick is available
	if _, err := exec.LookPath("convert"); err != nil {
		// Fallback: just copy the file if ImageMagick is not available
		return copyFile(inputPath, outputPath)
	}

	// Use ImageMagick to strip metadata
	cmd := exec.Command("convert", inputPath, "-strip", outputPath)
	return cmd.Run()
}

// cleanPDFMetadata removes metadata from PDF files using qpdf
func cleanPDFMetadata(inputPath, outputPath string) error {
	// Check if qpdf is available
	if _, err := exec.LookPath("qpdf"); err != nil {
		// Fallback: just copy the file if qpdf is not available
		return copyFile(inputPath, outputPath)
	}

	// Use qpdf to clean metadata
	cmd := exec.Command("qpdf", "--linearize", "--deterministic-id", inputPath, outputPath)
	return cmd.Run()
}

// copyFile copies a file from source to destination
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// GetFileExtension returns the file extension from content type
func GetFileExtension(contentType string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "application/pdf":
		return ".pdf"
	case "text/plain":
		return ".txt"
	case "application/msword":
		return ".doc"
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return ".docx"
	default:
		return ""
	}
}
