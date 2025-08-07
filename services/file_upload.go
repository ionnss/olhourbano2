package services

import (
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"olhourbano2/models"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	MaxFileSize = 50 * 1024 * 1024 // 50MB for videos
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
	maxSize := getMaxFileSize(category, result.ContentType)
	if header.Size > maxSize {
		result.Error = fmt.Errorf("arquivo muito grande. Máximo permitido: %dMB", maxSize/(1024*1024))
		return result, result.Error
	}

	// Validate file type using models package
	if !models.IsFileTypeAllowed(category, result.ContentType) {
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

// getMaxFileSize returns the maximum file size based on file type and category
func getMaxFileSize(category, contentType string) int64 {
	// Videos get larger size limit
	if strings.HasPrefix(contentType, "video/") {
		return 50 * 1024 * 1024 // 50MB for videos
	}

	// Default size for other files
	return 10 * 1024 * 1024 // 10MB
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
	case strings.HasPrefix(contentType, "video/"):
		return cleanVideoMetadata(inputPath, outputPath)
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

// cleanVideoMetadata removes metadata from video files using ffmpeg
func cleanVideoMetadata(inputPath, outputPath string) error {
	// Check if ffmpeg is available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		// Fallback: just copy the file if ffmpeg is not available
		return copyFile(inputPath, outputPath)
	}

	// Use ffmpeg to strip metadata and re-encode
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-map_metadata", "-1", // Remove all metadata
		"-c:v", "copy", // Copy video stream without re-encoding
		"-c:a", "copy", // Copy audio stream without re-encoding
		"-y", // Overwrite output file
		outputPath)

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
	case "video/mp4":
		return ".mp4"
	case "video/avi":
		return ".avi"
	case "video/mov":
		return ".mov"
	case "video/wmv":
		return ".wmv"
	case "video/flv":
		return ".flv"
	case "video/webm":
		return ".webm"
	default:
		return ""
	}
}
