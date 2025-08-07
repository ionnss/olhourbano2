package models

import (
	"encoding/json"
	"time"
)

// Report represents a report in the database
type Report struct {
	ID            int             `json:"id" db:"id"`
	ProblemType   string          `json:"problem_type" db:"problem_type"`
	HashedCPF     string          `json:"-" db:"hashed_cpf"` // Don't expose in JSON
	BirthDate     string          `json:"-" db:"birth_date"` // Don't expose in JSON
	Email         string          `json:"email" db:"email"`
	Location      string          `json:"location" db:"location"`
	City          string          `json:"city" db:"city"`
	Latitude      float64         `json:"latitude" db:"latitude"`
	Longitude     float64         `json:"longitude" db:"longitude"`
	Description   string          `json:"description" db:"description"`
	PhotoPath     string          `json:"photo_path" db:"photo_path"`
	TransportType string          `json:"transport_type,omitempty" db:"transport_type"`
	TransportData json.RawMessage `json:"transport_data,omitempty" db:"transport_data"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	VoteCount     int             `json:"vote_count" db:"vote_count"`
	Status        string          `json:"status" db:"status"`
}

// TransportData represents the transport-specific information
type TransportData struct {
	// Bus fields
	BusNumber  string `json:"bus_number,omitempty"`
	BusLine    string `json:"bus_line,omitempty"`
	BusStop    string `json:"bus_stop,omitempty"`
	BusCompany string `json:"bus_company,omitempty"`

	// Metro fields
	MetroLine    string `json:"metro_line,omitempty"`
	MetroStation string `json:"metro_station,omitempty"`
	MetroWagon   string `json:"metro_wagon,omitempty"`
	MetroCard    string `json:"metro_card,omitempty"`

	// Train fields
	TrainLine    string `json:"train_line,omitempty"`
	TrainStation string `json:"train_station,omitempty"`
	TrainWagon   string `json:"train_wagon,omitempty"`

	// Other transport fields
	TransportDetails string `json:"transport_details,omitempty"`
}

// ReportFormData represents the form data for creating a report
type ReportFormData struct {
	Category          string        `form:"category" validate:"required"`
	CPF               string        `form:"cpf" validate:"required,cpf"`
	BirthDate         string        `form:"birth_date" validate:"required"`
	Email             string        `form:"email" validate:"required,email"`
	EmailConfirmation string        `form:"email_confirmation" validate:"required,eqfield=Email"`
	Location          string        `form:"location" validate:"required"`
	Latitude          float64       `form:"latitude" validate:"required"`
	Longitude         float64       `form:"longitude" validate:"required"`
	Description       string        `form:"description" validate:"required,min=10,max=1000"`
	TransportType     string        `form:"transport_type"`
	TransportData     TransportData `form:"transport_data"`
}

// FileUpload represents an uploaded file
type FileUpload struct {
	OriginalName string
	SavedPath    string
	FileSize     int64
	ContentType  string
}

// ReportStatus constants
const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"
	StatusInReview = "in_review"
)

// AllowedFileTypes defines allowed file types per category
var AllowedFileTypes = map[string][]string{
	"default": {
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/webp",
		"video/mp4",
		"video/avi",
		"video/mov",
		"video/wmv",
		"video/flv",
		"video/webm",
	},
	"corrupcao_gestao_publica": {
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/webp",
		"application/pdf",
		"text/plain",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"video/mp4",
		"video/avi",
		"video/mov",
		"video/wmv",
		"video/flv",
		"video/webm",
	},
	"outros": {
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/webp",
		"application/pdf",
		"text/plain",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"video/mp4",
		"video/avi",
		"video/mov",
		"video/wmv",
		"video/flv",
		"video/webm",
	},
}

// GetAllowedFileTypes returns allowed file types for a category
func GetAllowedFileTypes(category string) []string {
	if types, exists := AllowedFileTypes[category]; exists {
		return types
	}
	return AllowedFileTypes["default"]
}

// IsFileTypeAllowed checks if a file type is allowed for a category
func IsFileTypeAllowed(category, contentType string) bool {
	allowedTypes := GetAllowedFileTypes(category)
	for _, allowedType := range allowedTypes {
		if allowedType == contentType {
			return true
		}
	}
	return false
}

// MaxFiles defines maximum files per category
var MaxFiles = map[string]int{
	"corrupcao_gestao_publica": 5,
	"outros":                   5,
	"default":                  2,
}

// GetMaxFiles returns maximum allowed files for a category
func GetMaxFiles(category string) int {
	if max, exists := MaxFiles[category]; exists {
		return max
	}
	return MaxFiles["default"]
}

// GetTransportDataFromReport extracts TransportData from a Report
func (r *Report) GetTransportData() (*TransportData, error) {
	if r.TransportData == nil {
		return nil, nil
	}

	var transportData TransportData
	err := json.Unmarshal(r.TransportData, &transportData)
	if err != nil {
		return nil, err
	}
	return &transportData, nil
}

// SetTransportData sets TransportData in a Report
func (r *Report) SetTransportData(data *TransportData) error {
	if data == nil {
		r.TransportData = nil
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r.TransportData = jsonData
	return nil
}
