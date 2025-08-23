package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    uint32 = 3         // số vòng (t)
	argonMemory  uint32 = 64 * 1024 // KiB (m): 64MB
	argonThreads uint8  = 2         // số luồng (p)
	saltLen             = 16        // bytes
	keyLen              = 32        // bytes (256-bit)
)

// HashPassword tạo hash theo format:
// $argon2id$v=19$m=<mem>,t=<time>,p=<threads>$<b64(salt)>$<b64(hash)>
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("empty password")
	}

	// tạo salt ngẫu nhiên
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("salt gen failed: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, keyLen)

	// encode theo tiêu chuẩn phổ biến để dễ lưu + verify
	b64 := base64.RawStdEncoding
	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argonMemory, argonTime, argonThreads,
		b64.EncodeToString(salt),
		b64.EncodeToString(hash),
	)
	return encoded, nil
}

// VerifyPassword kiểm tra password với hash đã lưu
func VerifyPassword(encodedHash, password string) (bool, error) {
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
	if subtleCompare(gotHash, wantHash) {
		return true, nil
	}
	return false, nil
}

// subtleCompare so sánh hằng thời gian
func subtleCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var v byte
	for i := 0; i < len(a); i++ {
		v |= a[i] ^ b[i]
	}
	return v == 0
}
