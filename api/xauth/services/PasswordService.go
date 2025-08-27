package services

type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(encodedHash, password string) (bool, error)
}
type PasswordArgon struct {
	argonTime    uint32 // số vòng (t)
	argonMemory  uint32 // KiB (m): 64MB
	argonThreads uint8  // số luồng (p)
	saltLen      int    // bytes
	keyLen       uint32 // bytes (256-bit)

}

func NewAuthServiceArgon() *PasswordArgon {
	return &PasswordArgon{
		argonTime:    3,
		argonMemory:  64 * 1024,
		argonThreads: 2,
		saltLen:      16,
		keyLen:       32,
	}
}
