package services

import (
	"database/sql"
	"fmt"
	"olhourbano2/models"
	"strings"
	"time"
)

// ReportStats represents report statistics
type ReportStats struct {
	TotalReports int
	ThisMonth    int
	Resolved     int
}

// CreateReport inserts a new report into the database
func CreateReport(db *sql.DB, report *models.Report) (int, error) {
	query := `
		INSERT INTO reports (problem_type, hashed_cpf, birth_date, email, location, latitude, longitude, description, photo_path, transport_type, transport_data, created_at, vote_count, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`

	// Handle nullable transport fields
	var transportType, transportData sql.NullString

	if report.TransportType != "" {
		transportType.String = report.TransportType
		transportType.Valid = true
	}

	if report.TransportData != nil {
		transportData.String = string(report.TransportData)
		transportData.Valid = true
	}

	var id int
	err := db.QueryRow(
		query,
		report.ProblemType,
		report.HashedCPF,
		report.BirthDate,
		report.Email,
		report.Location,
		report.Latitude,
		report.Longitude,
		report.Description,
		report.PhotoPath,
		transportType,
		transportData,
		time.Now(),
		0, // vote_count
		models.StatusPending,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetReportByID retrieves a report by its ID
func GetReportByID(db *sql.DB, id int) (*models.Report, error) {
	query := `
		SELECT id, problem_type, hashed_cpf, birth_date, email, location, latitude, longitude, description, photo_path, transport_type, transport_data, created_at, vote_count, status
		FROM reports
		WHERE id = $1
	`

	report := &models.Report{}
	var transportType, transportData sql.NullString

	err := db.QueryRow(query, id).Scan(
		&report.ID,
		&report.ProblemType,
		&report.HashedCPF,
		&report.BirthDate,
		&report.Email,
		&report.Location,
		&report.Latitude,
		&report.Longitude,
		&report.Description,
		&report.PhotoPath,
		&transportType,
		&transportData,
		&report.CreatedAt,
		&report.VoteCount,
		&report.Status,
	)

	if err != nil {
		return nil, err
	}

	// Convert nullable strings to regular strings
	if transportType.Valid {
		report.TransportType = transportType.String
	}
	if transportData.Valid {
		report.TransportData = []byte(transportData.String)
	}

	return report, nil
}

// GetReports retrieves reports with pagination and filtering
func GetReports(db *sql.DB, page int, category, status, city string, limit int) ([]*models.Report, error) {
	offset := (page - 1) * limit

	query := `
		SELECT id, problem_type, hashed_cpf, birth_date, email, location, latitude, longitude, description, photo_path, transport_type, transport_data, created_at, vote_count, status
		FROM reports
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 0

	if category != "" {
		argCount++
		query += fmt.Sprintf(" AND problem_type = $%d", argCount)
		args = append(args, category)
	}

	if status != "" {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
	}

	if city != "" {
		argCount++
		query += fmt.Sprintf(" AND location ILIKE $%d", argCount)
		args = append(args, "%"+city+"%")
	}

	argCount++
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", argCount)
	args = append(args, limit)

	argCount++
	query += fmt.Sprintf(" OFFSET $%d", argCount)
	args = append(args, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*models.Report
	for rows.Next() {
		report := &models.Report{}
		var transportType, transportData sql.NullString

		err := rows.Scan(
			&report.ID,
			&report.ProblemType,
			&report.HashedCPF,
			&report.BirthDate,
			&report.Email,
			&report.Location,
			&report.Latitude,
			&report.Longitude,
			&report.Description,
			&report.PhotoPath,
			&transportType,
			&transportData,
			&report.CreatedAt,
			&report.VoteCount,
			&report.Status,
		)
		if err != nil {
			return nil, err
		}

		// Convert nullable strings to regular strings
		if transportType.Valid {
			report.TransportType = transportType.String
		}
		if transportData.Valid {
			report.TransportData = []byte(transportData.String)
		}

		reports = append(reports, report)
	}

	return reports, nil
}

// GetTotalReports returns the total number of reports with optional filtering
func GetTotalReports(db *sql.DB, category, status, city string) (int, error) {
	query := `SELECT COUNT(*) FROM reports WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if category != "" {
		argCount++
		query += fmt.Sprintf(" AND problem_type = $%d", argCount)
		args = append(args, category)
	}

	if status != "" {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
	}

	if city != "" {
		argCount++
		query += fmt.Sprintf(" AND location ILIKE $%d", argCount)
		args = append(args, "%"+city+"%")
	}

	var count int
	err := db.QueryRow(query, args...).Scan(&count)
	return count, err
}

// GetReportStats returns report statistics
func GetReportStats(db *sql.DB) (*ReportStats, error) {
	stats := &ReportStats{}

	// Total reports
	err := db.QueryRow("SELECT COUNT(*) FROM reports").Scan(&stats.TotalReports)
	if err != nil {
		return nil, err
	}

	// Reports this month
	err = db.QueryRow(`
		SELECT COUNT(*) FROM reports 
		WHERE created_at >= date_trunc('month', CURRENT_DATE)
	`).Scan(&stats.ThisMonth)
	if err != nil {
		return nil, err
	}

	// Resolved reports (approved)
	err = db.QueryRow(`
		SELECT COUNT(*) FROM reports 
		WHERE status = 'approved'
	`).Scan(&stats.Resolved)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// GetReportsForMap retrieves reports with location data for map display
func GetReportsForMap(db *sql.DB, category, status, city string) ([]*models.Report, error) {
	query := `
		SELECT id, problem_type, hashed_cpf, birth_date, email, location, latitude, longitude, description, photo_path, transport_type, transport_data, created_at, vote_count, status
		FROM reports
		WHERE latitude IS NOT NULL AND longitude IS NOT NULL AND latitude != 0 AND longitude != 0
	`
	args := []interface{}{}
	argCount := 0

	if category != "" {
		argCount++
		query += fmt.Sprintf(" AND problem_type = $%d", argCount)
		args = append(args, category)
	}

	if status != "" {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
	}

	if city != "" {
		argCount++
		query += fmt.Sprintf(" AND location ILIKE $%d", argCount)
		args = append(args, "%"+city+"%")
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*models.Report
	for rows.Next() {
		report := &models.Report{}
		var transportType, transportData sql.NullString

		err := rows.Scan(
			&report.ID,
			&report.ProblemType,
			&report.HashedCPF,
			&report.BirthDate,
			&report.Email,
			&report.Location,
			&report.Latitude,
			&report.Longitude,
			&report.Description,
			&report.PhotoPath,
			&transportType,
			&transportData,
			&report.CreatedAt,
			&report.VoteCount,
			&report.Status,
		)
		if err != nil {
			return nil, err
		}

		// Convert nullable strings to regular strings
		if transportType.Valid {
			report.TransportType = transportType.String
		}
		if transportData.Valid {
			report.TransportData = []byte(transportData.String)
		}

		reports = append(reports, report)
	}

	return reports, nil
}

// GetCitiesFromReports retrieves unique cities from reports
func GetCitiesFromReports(db *sql.DB) ([]string, error) {
	query := `
		SELECT DISTINCT location 
		FROM reports 
		WHERE location IS NOT NULL AND location != '' 
		ORDER BY location
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []string
	for rows.Next() {
		var location string
		err := rows.Scan(&location)
		if err != nil {
			return nil, err
		}

		// Extract city name from location (assuming format like "Street, City, State")
		parts := strings.Split(location, ",")
		if len(parts) >= 2 {
			city := strings.TrimSpace(parts[len(parts)-2]) // Second to last part is usually city
			if city != "" && !contains(cities, city) {
				cities = append(cities, city)
			}
		}
	}

	return cities, nil
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
