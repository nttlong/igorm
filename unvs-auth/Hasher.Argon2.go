package unvsauth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type Argon2Hasher struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

func (a *Argon2Hasher) Hash(password string) (string, error) {
	salt := make([]byte, a.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, a.Time, a.Memory, a.Threads, a.KeyLen)
	encoded := fmt.Sprintf("$argon2id$%s$%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
	return encoded, nil
}

func (a *Argon2Hasher) Verify(hashed, password string) bool {
	// Đơn giản: parse salt + hash từ chuỗi (đã encode theo định dạng trên)
	var saltB64, hashB64 string
	_, err := fmt.Sscanf(hashed, "$argon2id$%s$%s", &saltB64, &hashB64)
	if err != nil {
		return false
	}
	salt, _ := base64.RawStdEncoding.DecodeString(saltB64)
	expectedHash, _ := base64.RawStdEncoding.DecodeString(hashB64)

	computed := argon2.IDKey([]byte(password), salt, a.Time, a.Memory, a.Threads, a.KeyLen)
	return string(computed) == string(expectedHash)
}
