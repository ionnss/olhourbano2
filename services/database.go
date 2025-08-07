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
	// Extract city name from location for better filtering
	city := ExtractCityFromLocation(report.Location)

	query := `
		INSERT INTO reports (problem_type, hashed_cpf, birth_date, email, location, city, latitude, longitude, description, photo_path, transport_type, transport_data, created_at, vote_count, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
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
		city,
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
		SELECT id, problem_type, hashed_cpf, birth_date, email, location, city, latitude, longitude, description, photo_path, transport_type, transport_data, created_at, vote_count, status
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
		&report.City,
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
func GetReports(db *sql.DB, page int, category, status, city, sort string, limit int) ([]*models.Report, error) {
	offset := (page - 1) * limit

	query := `
		SELECT id, problem_type, hashed_cpf, birth_date, email, location, city, latitude, longitude, description, photo_path, transport_type, transport_data, created_at, vote_count, status
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
		query += fmt.Sprintf(" AND city ILIKE $%d", argCount)
		args = append(args, "%"+city+"%")
	}

	// Add ORDER BY clause based on sort parameter
	switch sort {
	case "votes":
		query += " ORDER BY vote_count DESC, created_at DESC"
	case "oldest":
		query += " ORDER BY created_at ASC"
	case "recent":
		fallthrough
	default:
		query += " ORDER BY created_at DESC"
	}

	argCount++
	query += fmt.Sprintf(" LIMIT $%d", argCount)
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
			&report.City,
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
		query += fmt.Sprintf(" AND city ILIKE $%d", argCount)
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

// AddVote adds a vote to a report
func AddVote(db *sql.DB, reportID int, hashedCPF string) error {
	// First, try to insert the vote
	_, err := db.Exec(`
		INSERT INTO votes (report_id, vote_hashed_cpf, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (vote_hashed_cpf, report_id) DO NOTHING
	`, reportID, hashedCPF)

	if err != nil {
		return fmt.Errorf("error adding vote: %w", err)
	}

	// Update the vote count in the reports table
	_, err = db.Exec(`
		UPDATE reports 
		SET vote_count = (
			SELECT COUNT(*) 
			FROM votes 
			WHERE report_id = $1
		)
		WHERE id = $1
	`, reportID)

	if err != nil {
		return fmt.Errorf("error updating vote count: %w", err)
	}

	return nil
}

// GetVoteCount returns the vote count for a report
func GetVoteCount(db *sql.DB, reportID int) (int, error) {
	var count int
	err := db.QueryRow(`
		SELECT vote_count 
		FROM reports 
		WHERE id = $1
	`, reportID).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("error getting vote count: %w", err)
	}

	return count, nil
}

// HasUserVoted checks if a user (by hashed CPF) has already voted for a specific report
func HasUserVoted(db *sql.DB, reportID int, hashedCPF string) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM votes 
		WHERE report_id = $1 AND vote_hashed_cpf = $2
	`, reportID, hashedCPF).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("error checking if user has voted: %w", err)
	}

	return count > 0, nil
}

// GetReportsForMap retrieves reports with location data for map display
func GetReportsForMap(db *sql.DB, category, status, city string) ([]*models.Report, error) {
	query := `
		SELECT id, problem_type, hashed_cpf, birth_date, email, location, city, latitude, longitude, description, photo_path, transport_type, transport_data, created_at, vote_count, status
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
		query += fmt.Sprintf(" AND city ILIKE $%d", argCount)
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
			&report.City,
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

// GetCitiesFromReports retrieves unique cities from reports using reverse geocoding
func GetCitiesFromReports(db *sql.DB) ([]string, error) {
	// First, try to get cities from a dedicated city column if it exists
	// This is more efficient than reverse geocoding every time
	query := `
		SELECT DISTINCT city 
		FROM reports 
		WHERE city IS NOT NULL AND city != '' 
		ORDER BY city
	`

	rows, err := db.Query(query)
	if err != nil {
		// If city column doesn't exist, fall back to extracting from location
		return getCitiesFromLocation(db)
	}
	defer rows.Close()

	var cities []string
	for rows.Next() {
		var city string
		err := rows.Scan(&city)
		if err != nil {
			return nil, err
		}
		if city != "" && !contains(cities, city) {
			cities = append(cities, city)
		}
	}

	return cities, nil
}

// getCitiesFromLocation extracts cities from location field as fallback
func getCitiesFromLocation(db *sql.DB) ([]string, error) {
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

		// Try to extract city name from location using multiple strategies
		city := ExtractCityFromLocation(location)
		if city != "" && !contains(cities, city) {
			cities = append(cities, city)
		}
	}

	return cities, nil
}

// ExtractCityFromLocation tries to extract city name from location string
func ExtractCityFromLocation(location string) string {
	if location == "" {
		return ""
	}

	// Clean up the location string
	location = strings.TrimSpace(location)

	// Remove common suffixes and prefixes that might interfere
	location = strings.ReplaceAll(location, " - Brasil", "")
	location = strings.ReplaceAll(location, ", Brazil", "")
	location = strings.ReplaceAll(location, ", Brasil", "")
	location = strings.ReplaceAll(location, "Brasil", "")
	location = strings.ReplaceAll(location, "Brazil", "")

	// Remove extra whitespace
	location = strings.Join(strings.Fields(location), " ")

	// Split by comma and clean up
	parts := strings.Split(location, ",")
	if len(parts) < 2 {
		return ""
	}

	// For Brazilian addresses, the city is usually the part that contains " - " followed by state
	// Format: "Street, Number - Neighborhood, City - State, CEP"
	for i := len(parts) - 1; i >= 0; i-- {
		part := strings.TrimSpace(parts[i])

		// Skip empty parts
		if part == "" {
			continue
		}

		// Skip parts that look like postal codes (CEP)
		if isPostalCode(part) {
			continue
		}

		// Look for parts that contain " - " which usually indicates "City - State"
		if strings.Contains(part, " - ") {
			subParts := strings.Split(part, " - ")
			if len(subParts) >= 2 {
				cityPart := strings.TrimSpace(subParts[0])
				statePart := strings.TrimSpace(subParts[1])

				// If the second part looks like a state abbreviation, the first part is the city
				if isStateAbbreviation(statePart) && isValidCityName(cityPart) {
					return cleanCityName(cityPart)
				}
			}
		}

		// Skip parts that look like state abbreviations
		if isStateAbbreviation(part) {
			continue
		}

		// Skip parts that are too short or don't look like city names
		if !isValidCityName(part) {
			continue
		}

		// If we haven't found a city yet and this part doesn't look like a postal code or state,
		// and it's not the first part (which is usually the street), it might be the city
		if i > 0 && !isPostalCode(part) && !isStateAbbreviation(part) {
			// Check if this part contains a state abbreviation that should be removed
			if strings.Contains(part, " - ") {
				subParts := strings.Split(part, " - ")
				if len(subParts) >= 2 {
					cityPart := strings.TrimSpace(subParts[0])
					statePart := strings.TrimSpace(subParts[1])

					// If the second part looks like a state abbreviation, return just the city part
					if isStateAbbreviation(statePart) && isValidCityName(cityPart) {
						return cleanCityName(cityPart)
					}
				}
			}

			// Additional validation: check if this part looks like a city name
			if isValidCityName(part) {
				return cleanCityName(part)
			}
		}
	}

	return ""
}

// isPostalCode checks if a string looks like a Brazilian postal code
func isPostalCode(s string) bool {
	// Remove non-digits
	digits := strings.ReplaceAll(s, "-", "")
	digits = strings.ReplaceAll(digits, " ", "")

	// Brazilian CEP format: 5 digits + 3 digits (optional)
	return len(digits) == 8 || len(digits) == 5
}

// isStateAbbreviation checks if a string looks like a Brazilian state abbreviation
func isStateAbbreviation(s string) bool {
	stateAbbreviations := []string{
		"AC", "AL", "AP", "AM", "BA", "CE", "DF", "ES", "GO", "MA",
		"MT", "MS", "MG", "PA", "PB", "PR", "PE", "PI", "RJ", "RN",
		"RS", "RO", "RR", "SC", "SP", "SE", "TO",
	}

	for _, abbr := range stateAbbreviations {
		if strings.EqualFold(s, abbr) {
			return true
		}
	}
	return false
}

// isValidCityName checks if a string looks like a valid city name
func isValidCityName(s string) bool {
	// Must be at least 3 characters
	if len(s) < 3 {
		return false
	}

	// Must not contain numbers (except for cities like "São João")
	if strings.ContainsAny(s, "0123456789") {
		return false
	}

	// Must not contain special characters that are not common in city names
	// Allow letters, spaces, hyphens, apostrophes, and common accented characters
	invalidChars := []string{"@", "#", "$", "%", "^", "&", "*", "(", ")", "+", "=", "[", "]", "{", "}", "|", "\\", ":", ";", "\"", "'", "<", ">", ",", ".", "?", "/"}
	for _, char := range invalidChars {
		if strings.Contains(s, char) {
			return false
		}
	}

	// Must not be all uppercase (likely an abbreviation)
	if s == strings.ToUpper(s) && len(s) <= 5 {
		return false
	}

	// Must not contain common non-city words
	nonCityWords := []string{"rua", "avenida", "av", "r", "n", "s", "l", "o", "e", "norte", "sul", "leste", "oeste", "centro", "bairro", "distrito", "zona"}
	for _, word := range nonCityWords {
		if strings.EqualFold(s, word) {
			return false
		}
	}

	return true
}

// cleanCityName removes extra text and normalizes city names
func cleanCityName(s string) string {
	// Remove extra whitespace
	s = strings.Join(strings.Fields(s), " ")

	// Remove common suffixes that might be attached to city names
	suffixes := []string{" - ", " -", "- ", " - SP", " - RJ", " - MG", " - RS", " - PR", " - SC", " - BA", " - CE", " - PE", " - GO", " - MT", " - MS", " - MA", " - PA", " - PB", " - PI", " - RN", " - AL", " - SE", " - TO", " - RO", " - AC", " - RR", " - AP", " - AM", " - DF", " - ES"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			s = strings.TrimSuffix(s, suffix)
			break
		}
	}

	// Remove common prefixes
	prefixes := []string{"Cidade de ", "Município de ", "Municipio de "}
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			s = strings.TrimPrefix(s, prefix)
			break
		}
	}

	// Capitalize properly (first letter of each word)
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			// Handle common prepositions and articles
			lowerWords := []string{"de", "da", "do", "das", "dos", "e", "em", "na", "no", "nas", "nos", "para", "por", "com", "sem", "sob", "sobre", "entre", "contra", "desde", "até", "atras", "após", "ante", "perante", "segundo", "conforme", "mediante", "salvo", "exceto", "menos", "fora", "tirante", "senão", "que", "qual", "quais", "cujo", "cuja", "cujos", "cujas", "onde", "quando", "como", "porque", "pois", "logo", "portanto", "então", "assim", "também", "tambem", "ainda", "já", "ja", "nunca", "sempre", "só", "so", "apenas", "sómente", "somente", "unicamente", "exclusivamente", "inclusive", "inclusivamente", "exclusive", "exclusivamente", "salvo", "exceto", "menos", "fora", "tirante", "senão", "que", "qual", "quais", "cujo", "cuja", "cujos", "cujas", "onde", "quando", "como", "porque", "pois", "logo", "portanto", "então", "assim", "também", "tambem", "ainda", "já", "ja", "nunca", "sempre", "só", "so", "apenas", "sómente", "somente", "unicamente", "exclusivamente", "inclusive", "inclusivamente", "exclusive", "exclusivamente"}
			isLower := false
			for _, lowerWord := range lowerWords {
				if strings.EqualFold(word, lowerWord) {
					isLower = true
					break
				}
			}

			if !isLower {
				// Capitalize first letter, keep rest as is
				if len(word) > 1 {
					words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
				} else {
					words[i] = strings.ToUpper(word)
				}
			} else {
				words[i] = strings.ToLower(word)
			}
		}
	}

	return strings.Join(words, " ")
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

// UpdateExistingReportsWithCity extracts and updates city field for existing reports
func UpdateExistingReportsWithCity(db *sql.DB) error {
	query := `
		UPDATE reports 
		SET city = $1
		WHERE id = $2
	`

	// Get all reports that have location data
	rows, err := db.Query(`
		SELECT id, location 
		FROM reports 
		WHERE location IS NOT NULL AND location != ''
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var location string
		err := rows.Scan(&id, &location)
		if err != nil {
			return err
		}

		city := ExtractCityFromLocation(location)
		if city != "" {
			_, err = db.Exec(query, city, id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
