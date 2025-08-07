package models

import (
	"time"
)

// Comment represents a comment on a report
type Comment struct {
	ID        int       `json:"id" db:"id"`
	ReportID  int       `json:"report_id" db:"report_id"`
	HashedCPF string    `json:"-" db:"hashed_cpf"` // Don't expose in JSON
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CommentDisplay represents a comment with display information
type CommentDisplay struct {
	ID              int       `json:"id"`
	ReportID        int       `json:"report_id"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"created_at"`
	HashedCPFDisplay string   `json:"hashed_cpf_display"`
}

// CommentFormData represents the form data for creating a comment
type CommentFormData struct {
	ReportID  int    `form:"report_id" validate:"required"`
	CPF       string `form:"cpf" validate:"required,cpf"`
	BirthDate string `form:"birth_date" validate:"required"`
	Content   string `form:"content" validate:"required,min=1,max=500"`
}



// GetHashedCPFDisplay returns the display format for hashed CPF
func (c *Comment) GetHashedCPFDisplay() string {
	if len(c.HashedCPF) >= 8 {
		return c.HashedCPF[:8]
	}
	return c.HashedCPF
}
