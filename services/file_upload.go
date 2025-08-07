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
	MaxFileSize   = 50 * 1024 * 1024 // 50MB for videos
	UploadDir     = "./uploads"
	ThumbnailDir  = "./uploads/thumbnails"
	ThumbnailSize = "150x150" // Small thumbnails for performance
)

// FileUploadResult represents the result of a file upload
type FileUploadResult struct {
	OriginalName  string
	SavedPath     string
	ThumbnailPath string
	FileSize      int64
	ContentType   string
	Error         error
}

// ProcessFileUpload handles file upload with metadata cleaning and thumbnail generation
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

	// Ensure thumbnail directory exists
	if err := os.MkdirAll(ThumbnailDir, 0755); err != nil {
		result.Error = fmt.Errorf("erro ao criar diretório de thumbnails: %v", err)
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

	// Generate thumbnail asynchronously (don't block upload)
	go generateThumbnailAsync(result.SavedPath, result.ContentType, hash)

	return result, nil
}

// generateThumbnailAsync generates thumbnail for the uploaded file
func generateThumbnailAsync(filePath, contentType, hash string) {
	// Determine if we should generate thumbnail
	if !shouldGenerateThumbnail(contentType) {
		return
	}

	thumbnailPath := filepath.Join(ThumbnailDir, hash+"_thumb.jpg")

	var err error
	if strings.HasPrefix(contentType, "video/") {
		err = generateVideoThumbnail(filePath, thumbnailPath)
	} else if contentType == "application/pdf" {
		err = generatePDFThumbnail(filePath, thumbnailPath)
	}

	if err != nil {
		// Log error but don't fail the upload
		fmt.Printf("Warning: Failed to generate thumbnail for %s: %v\n", filePath, err)
	}
}

// shouldGenerateThumbnail determines if we should generate a thumbnail for this file type
func shouldGenerateThumbnail(contentType string) bool {
	return strings.HasPrefix(contentType, "video/") || contentType == "application/pdf"
}

// generateVideoThumbnail generates a thumbnail from the first frame of a video
func generateVideoThumbnail(videoPath, thumbnailPath string) error {
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-ss", "00:00:01", // Start at 1 second to avoid black frames
		"-vframes", "1",
		"-vf", "scale=150:150:force_original_aspect_ratio=decrease,pad=150:150:(ow-iw)/2:(oh-ih)/2",
		"-y", // Overwrite output file
		thumbnailPath,
	)

	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// generatePDFThumbnail generates a thumbnail from the first page of a PDF
func generatePDFThumbnail(pdfPath, thumbnailPath string) error {
	cmd := exec.Command("convert",
		"-density", "150", // Higher DPI for better quality
		"-resize", ThumbnailSize,
		"-background", "white",
		"-alpha", "remove",
		"-alpha", "off",
		fmt.Sprintf("%s[0]", pdfPath), // First page only
		thumbnailPath,
	)

	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetThumbnailPath returns the thumbnail path for a given file
func GetThumbnailPath(filePath string) string {
	if filePath == "" {
		return ""
	}

	// Extract hash from filename
	filename := filepath.Base(filePath)
	ext := filepath.Ext(filename)
	hash := strings.TrimSuffix(filename, ext)

	thumbnailPath := filepath.Join(ThumbnailDir, hash+"_thumb.jpg")

	// Check if thumbnail exists
	if _, err := os.Stat(thumbnailPath); err == nil {
		return thumbnailPath
	}

	// Return empty string if thumbnail doesn't exist
	return ""
}

// GetFileTypeIcon returns the appropriate icon for a file type
func GetFileTypeIcon(contentType string) string {
	switch {
	case strings.HasPrefix(contentType, "image/"):
		return "bi-image"
	case strings.HasPrefix(contentType, "video/"):
		return "bi-camera-video"
	case contentType == "application/pdf":
		return "bi-file-pdf"
	case strings.HasPrefix(contentType, "text/"):
		return "bi-file-text"
	case strings.HasPrefix(contentType, "application/msword") ||
		strings.HasPrefix(contentType, "application/vnd.openxmlformats-officedocument.wordprocessingml"):
		return "bi-file-word"
	default:
		return "bi-file-earmark"
	}
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

// CleanupThumbnail removes the thumbnail for a given file
func CleanupThumbnail(filePath string) error {
	if filePath == "" {
		return nil
	}

	thumbnailPath := GetThumbnailPath(filePath)
	if thumbnailPath != "" {
		return os.Remove(thumbnailPath)
	}
	return nil
}

// CleanupThumbnailsForReport removes thumbnails for all files in a report
func CleanupThumbnailsForReport(photoPath string) error {
	if photoPath == "" {
		return nil
	}

	// Split by comma and clean up each thumbnail
	rawPhotos := strings.Split(photoPath, ",")
	for _, photo := range rawPhotos {
		trimmed := strings.TrimSpace(photo)
		if trimmed != "" {
			if err := CleanupThumbnail(trimmed); err != nil {
				// Log error but continue with other files
				fmt.Printf("Warning: Failed to cleanup thumbnail for %s: %v\n", trimmed, err)
			}
		}
	}
	return nil
}
