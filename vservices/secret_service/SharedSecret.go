package secret_service

import (
	"encoding/base64"
	"fmt"
	"math/rand"
)

// SharedSecretService sinh ra SharedSecret an toàn (dài, ngẫu nhiên).
type SharedSecretService interface {
	Generate() string
}

type sharedSecretService struct{}

func NewSharedSecretService() SharedSecretService {
	return &sharedSecretService{}
}

// Generate tạo ra chuỗi shared secret ngẫu nhiên, base64 encoded.
func (s *sharedSecretService) Generate() string {
	var min, max = 32, 64
	length := rand.Intn(max-min+1) + min
	if length <= 0 {
		panic(fmt.Errorf("invalid secret length: %d", length))
	}

	// Base64 sẽ làm dài chuỗi ra ~4/3, nên sinh ít hơn
	byteLen := (length * 3) / 4
	buf := make([]byte, byteLen)

	_, err := rand.Read(buf)
	if err != nil {
		panic(fmt.Errorf("failed to generate random bytes: %w", err))
	}

	secret := base64.URLEncoding.EncodeToString(buf)
	if len(secret) > length {
		secret = secret[:length]
	}
	return secret
}
