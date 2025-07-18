package unvsauth

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct {
	Cost int // ví dụ: 10
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.Cost)
	return string(bytes), err
}

func (h *BcryptHasher) Verify(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
