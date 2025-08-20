package libs

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

type PasswordService struct {
}

// generateSalt tạo một salt ngẫu nhiên với n byte
func (utils *PasswordService) GenerateSalt(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// hashPassword hash password + salt
func (utils *PasswordService) HashPassword(password string) (string, error) {
	h := sha256.New()
	salt, err := utils.GenerateSalt(16)
	if err != nil {
		return "", err
	}
	h.Write([]byte(password + salt))
	return hex.EncodeToString(h.Sum(nil)), nil
}
