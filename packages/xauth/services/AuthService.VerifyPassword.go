package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// VerifyPassword kiểm tra password với hash đã lưu
func (authService *AuthService) VerifyPassword(encodedHash, password string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, errors.New("invalid hash format")
	}

	// parts[2] == v=19 (không dùng ở đây nhưng giữ để tương thích)
	var mem uint32
	var time uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &mem, &time, &threads); err != nil {
		return false, fmt.Errorf("invalid params: %w", err)
	}

	b64 := base64.RawStdEncoding
	salt, err := b64.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("invalid salt b64: %w", err)
	}
	wantHash, err := b64.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("invalid hash b64: %w", err)
	}

	gotHash := argon2.IDKey([]byte(password), salt, time, mem, threads, uint32(len(wantHash)))

	// hằng số so sánh để tránh timing leak
	if authService.subtleCompare(gotHash, wantHash) {
		return true, nil
	}
	return false, nil
}

// subtleCompare so sánh hằng thời gian
func (authService *AuthService) subtleCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var v byte
	for i := 0; i < len(a); i++ {
		v |= a[i] ^ b[i]
	}
	return v == 0
}
