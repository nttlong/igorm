package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// HashPassword tạo hash theo format:
// $argon2id$v=19$m=<mem>,t=<time>,p=<threads>$<b64(salt)>$<b64(hash)>
func (authService *AuthService) HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("empty password")
	}

	// tạo salt ngẫu nhiên
	salt := make([]byte, authService.saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("salt gen failed: %w", err)
	}

	hash := argon2.IDKey([]byte(password),
		salt, authService.argonTime,
		authService.argonMemory,
		authService.argonThreads,
		authService.keyLen,
	)

	// encode theo tiêu chuẩn phổ biến để dễ lưu + verify
	b64 := base64.RawStdEncoding
	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		authService.argonMemory, authService.argonTime, authService.argonThreads,
		b64.EncodeToString(salt),
		b64.EncodeToString(hash),
	)
	return encoded, nil
}
