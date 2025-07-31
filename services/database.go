package services

import (
	"database/sql"
	"fmt"
	"olhourbano2/models"
	"time"
)

// CreateReport inserts a new report into the database
func CreateReport(db *sql.DB, report *models.Report) (int, error) {
	query := `
		INSERT INTO reports (problem_type, hashed_cpf, birth_date, email, location, latitude, longitude, description, photo_path, created_at, vote_count, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`

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
		SELECT id, problem_type, hashed_cpf, birth_date, email, location, latitude, longitude, description, photo_path, created_at, vote_count, status
		FROM reports
		WHERE id = $1
	`

	report := &models.Report{}
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
		&report.CreatedAt,
		&report.VoteCount,
		&report.Status,
	)

	if err != nil {
		return nil, err
	}

	return report, nil
}

// GetReports retrieves reports with pagination and filtering
func GetReports(db *sql.DB, page int, category, status string, limit int) ([]*models.Report, error) {
	offset := (page - 1) * limit

	query := `
		SELECT id, problem_type, hashed_cpf, birth_date, email, location, latitude, longitude, description, photo_path, created_at, vote_count, status
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
			&report.CreatedAt,
			&report.VoteCount,
			&report.Status,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}
