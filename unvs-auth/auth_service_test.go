package unvsauth

import (
	"fmt"
	"testing"
	di "vdi"

	"github.com/stretchr/testify/assert"
)

func TestAuthService(t *testing.T) {
	di.RegisterContainer(func(svc *AuthService) error {
		svc.UserRepo.Init = func(owner *AuthService) *UserRepository {
			return NewUserRepository()
		}
		svc.JwtProvider.Init = func(owner *AuthService) *JwtProvider {
			return &JwtProvider{SecretKey: "your-secret-key"}
		}
		return nil
	})

	service, _ := di.Resolve[AuthService]()
	for i := 0; i < 10; i++ {

		user, err := service.Register("alice@example.com", "alice", "123456")
		fmt.Println("Register:", user, err)

		// Đăng nhập
		token, err := service.Login("alice@example.com", "123456")
		fmt.Println("Login token:", token)
	}
}
func BenchmarkAuthService(b *testing.B) {
	di.RegisterContainer(func(svc *AuthService) error {
		svc.UserRepo.Init = func(owner *AuthService) *UserRepository {
			return NewUserRepository()
		}
		svc.JwtProvider.Init = func(owner *AuthService) *JwtProvider {
			return &JwtProvider{SecretKey: "your-secret-key"}
		}
		return nil
	})

	service, err := di.Resolve[AuthService]()
	assert.NoError(b, err)
	for i := 0; i < b.N; i++ {
		email := fmt.Sprintf("alice%d@example.com", i)
		username := fmt.Sprintf("alice%d", i)

		user, err := service.Register(email, username, "123456")
		assert.NoError(b, err)
		assert.Equal(b, email, user.Email)

		token, err := service.Login(email, "123456")
		assert.NoError(b, err)
		assert.NotEmpty(b, token)
	}

}
func initAuthService() *AuthService {
	di.RegisterContainer(func(svc *AuthService) error {
		svc.UserRepo.Init = func(owner *AuthService) *UserRepository {
			return NewUserRepository()
		}
		svc.JwtProvider.Init = func(owner *AuthService) *JwtProvider {
			return &JwtProvider{SecretKey: "your-secret-key"}
		}
		return nil
	})
	service, _ := di.Resolve[AuthService]()
	return service
}
func BenchmarkRegisterOnly(b *testing.B) {
	service := initAuthService()
	for i := 0; i < b.N; i++ {
		email := fmt.Sprintf("user%d@example.com", i)
		username := fmt.Sprintf("user%d", i)
		_, err := service.Register(email, username, "123456")
		assert.NoError(b, err)
	}
}
func BenchmarkLoginOnly(b *testing.B) {
	service := initAuthService()
	service.Register("user@example.com", "user", "123456")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.Login("user@example.com", "123456")
		assert.NoError(b, err)
	}
}
func BenchmarkVerifyToken(b *testing.B) {
	service := initAuthService()
	service.Register("demo@example.com", "demo", "123456")
	token, _ := service.Login("demo@example.com", "123456")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.JwtProvider.Get().VerifyToken(token)
		assert.NoError(b, err)
	}
}
