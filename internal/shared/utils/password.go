package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// DefaultCost is the bcrypt hashing cost.
const DefaultCost = bcrypt.DefaultCost

// HashPassword hashes a plaintext password using bcrypt with standard cost.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a plaintext password with a bcrypt hashed password.
// Returns true if password matches, false otherwise.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
