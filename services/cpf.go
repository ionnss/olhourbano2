package services

import (
	"crypto/sha256"
	"fmt"
)

// HashCPF creates a secure hash of the CPF for database storage
func HashCPF(cpf string) string {
	// Normalize CPF (remove formatting)
	normalizedCPF := NormalizeCPF(cpf)

	// Create SHA-256 hash
	hash := sha256.Sum256([]byte(normalizedCPF))

	// Return hex representation
	return fmt.Sprintf("%x", hash)
}

// VerifyCPF checks if a CPF matches a stored hash
func VerifyCPF(cpf, hashedCPF string) bool {
	return HashCPF(cpf) == hashedCPF
}
