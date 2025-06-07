package config

var jwtSecret = []byte("super_secret_test_key_for_development_and_testing_only_1234567890ABCDEF")

func GetJWTSecret() []byte {
	return jwtSecret
}
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}
