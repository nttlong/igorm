package unvsauth

import "golang.org/x/crypto/bcrypt"

func HashPassword(pw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(hash, pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	return err == nil
}
