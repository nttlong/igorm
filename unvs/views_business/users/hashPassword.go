package users

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Create hashes the password with a random salt, combining it with lowercase username
func (v *User) hashPassword(userInfo userInfo) (string, error) {
	if userInfo.Username == "" || userInfo.Password == "" {
		return "", errors.New("username and password are required")
	}

	// Convert username to lowercase
	lowerUsername := strings.ToLower(userInfo.Username)

	// Combine username and password into one value
	combined := lowerUsername + userInfo.Password

	// Generate a random salt (16 bytes is a common size)
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Convert salt to base64 for safe storage
	saltStr := base64.StdEncoding.EncodeToString(salt)

	// Combine the combined value with salt for hashing
	dataToHash := combined + saltStr

	// Hash the combined data with bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(dataToHash), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Return the hashed password (including salt) as a string
	return string(hashed), nil
}
