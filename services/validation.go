package services

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ValidateCPF validates Brazilian CPF format and check digit
func ValidateCPF(cpf string) bool {
	// Remove non-numeric characters
	cpf = regexp.MustCompile(`\D`).ReplaceAllString(cpf, "")

	// Check length
	if len(cpf) != 11 {
		return false
	}

	// Check for known invalid patterns
	invalidCPFs := []string{
		"00000000000", "11111111111", "22222222222", "33333333333",
		"44444444444", "55555555555", "66666666666", "77777777777",
		"88888888888", "99999999999",
	}

	for _, invalid := range invalidCPFs {
		if cpf == invalid {
			return false
		}
	}

	// Validate check digits
	return validateCPFCheckDigits(cpf)
}

// validateCPFCheckDigits validates CPF check digits using the official algorithm
func validateCPFCheckDigits(cpf string) bool {
	// Convert to integer slice
	digits := make([]int, 11)
	for i, char := range cpf {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return false
		}
		digits[i] = digit
	}

	// Validate first check digit
	sum := 0
	for i := 0; i < 9; i++ {
		sum += digits[i] * (10 - i)
	}
	remainder := sum % 11
	firstCheckDigit := 0
	if remainder >= 2 {
		firstCheckDigit = 11 - remainder
	}

	if digits[9] != firstCheckDigit {
		return false
	}

	// Validate second check digit
	sum = 0
	for i := 0; i < 10; i++ {
		sum += digits[i] * (11 - i)
	}
	remainder = sum % 11
	secondCheckDigit := 0
	if remainder >= 2 {
		secondCheckDigit = 11 - remainder
	}

	return digits[10] == secondCheckDigit
}

// ValidateEmail validates email format using regex
func ValidateEmail(email string) bool {
	// Basic email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// NormalizeCPF removes formatting from CPF
func NormalizeCPF(cpf string) string {
	return regexp.MustCompile(`\D`).ReplaceAllString(cpf, "")
}

// FormatCPF formats CPF with standard Brazilian formatting
func FormatCPF(cpf string) string {
	cpf = NormalizeCPF(cpf)
	if len(cpf) != 11 {
		return cpf
	}
	return cpf[:3] + "." + cpf[3:6] + "." + cpf[6:9] + "-" + cpf[9:]
}

// ValidateForm validates the complete report form
func ValidateForm(category, cpf, birthDate, email, emailConfirmation, location, description string, latitude, longitude float64) []string {
	var errors []string

	// Validate CPF
	if !ValidateCPF(cpf) {
		errors = append(errors, "CPF inválido")
	}

	// Validate birth date
	if birthDate == "" {
		errors = append(errors, "Data de nascimento é obrigatória")
	}

	// Validate email
	if !ValidateEmail(email) {
		errors = append(errors, "Email inválido")
	}

	// Validate email confirmation
	if !strings.EqualFold(email, emailConfirmation) {
		errors = append(errors, "Confirmação de email não confere")
	}

	// Validate location
	if strings.TrimSpace(location) == "" {
		errors = append(errors, "Localização é obrigatória")
	}

	// Validate coordinates
	if latitude == 0 && longitude == 0 {
		errors = append(errors, "Coordenadas de localização são obrigatórias")
	}

	// Validate latitude range
	if latitude < -90 || latitude > 90 {
		errors = append(errors, "Latitude inválida")
	}

	// Validate longitude range
	if longitude < -180 || longitude > 180 {
		errors = append(errors, "Longitude inválida")
	}

	// Validate description
	description = strings.TrimSpace(description)
	if len(description) < 10 {
		errors = append(errors, "Descrição deve ter pelo menos 10 caracteres")
	}
	if len(description) > 1000 {
		errors = append(errors, "Descrição deve ter no máximo 1000 caracteres")
	}

	return errors
}

// ValidateFiles validates that at least one file is uploaded for a report
func ValidateFiles(fileCount int) []string {
	var errors []string

	if fileCount == 0 {
		errors = append(errors, "Pelo menos um arquivo (foto, vídeo ou documento) é obrigatório para comprovar a denúncia")
	}

	return errors
}

// ConvertBirthDateToDBFormat converts dd/mm/aaaa format to YYYY-MM-DD for database storage
func ConvertBirthDateToDBFormat(birthDate string) (string, error) {
	// Check if it's already in YYYY-MM-DD format (from old date inputs)
	if len(birthDate) == 10 && birthDate[4] == '-' && birthDate[7] == '-' {
		return birthDate, nil
	}

	// Parse dd/mm/aaaa format
	if len(birthDate) != 10 || birthDate[2] != '/' || birthDate[5] != '/' {
		return "", fmt.Errorf("formato inválido, use dd/mm/aaaa")
	}

	day := birthDate[0:2]
	month := birthDate[3:5]
	year := birthDate[6:10]

	// Validate components
	dayNum, err := strconv.Atoi(day)
	if err != nil || dayNum < 1 || dayNum > 31 {
		return "", fmt.Errorf("dia inválido")
	}

	monthNum, err := strconv.Atoi(month)
	if err != nil || monthNum < 1 || monthNum > 12 {
		return "", fmt.Errorf("mês inválido")
	}

	yearNum, err := strconv.Atoi(year)
	if err != nil || yearNum < 1900 || yearNum > time.Now().Year() {
		return "", fmt.Errorf("ano inválido")
	}

	// Validate the actual date
	date := time.Date(yearNum, time.Month(monthNum), dayNum, 0, 0, 0, 0, time.UTC)
	if date.Day() != dayNum || date.Month() != time.Month(monthNum) || date.Year() != yearNum {
		return "", fmt.Errorf("data inválida")
	}

	// Check if date is not in the future
	if date.After(time.Now()) {
		return "", fmt.Errorf("data não pode ser no futuro")
	}

	// Check minimum age (16 years)
	minAge := 16
	minDate := time.Now().AddDate(-minAge, 0, 0)
	if date.After(minDate) {
		return "", fmt.Errorf("idade mínima: %d anos", minAge)
	}

	// Convert to YYYY-MM-DD format
	return fmt.Sprintf("%04d-%02d-%02d", yearNum, monthNum, dayNum), nil
}
