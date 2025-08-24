package services

type AuthService struct {
	argonTime    uint32 // số vòng (t)
	argonMemory  uint32 // KiB (m): 64MB
	argonThreads uint8  // số luồng (p)
	saltLen      int    // bytes
	keyLen       uint32 // bytes (256-bit)
}

func (authService *AuthService) New() error {

	authService.argonTime = 3           // số vòng (t)
	authService.argonTime = 3           // số vòng (t)
	authService.argonMemory = 64 * 1024 // KiB (m): 64MB
	authService.argonThreads = 2        // số luồng (p)
	authService.saltLen = 16            // bytes
	authService.keyLen = 32
	return nil
}
