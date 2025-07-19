package unvsauth

import (
	"fmt"
	"testing"
	"vdb"
	di "vdi"

	"github.com/stretchr/testify/assert"
)

func TestAuthService(t *testing.T) {

	service := initAuthService()

	for i := 0; i < 10; i++ {

		user, err := service.Register("alice@example.com", "alice", "123456")
		fmt.Println("Register:", user, err)

		// Đăng nhập
		token, err := service.Login("alice@example.com", "123456")
		fmt.Println("Login token:", token)
	}
}
func BenchmarkAuthService(b *testing.B) {

	service := initAuthService()

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
		svc.UserRepo.Init = func(owner *AuthService) UserRepository { //👈 cung phai sua
			mssqlDns := "sqlserver://sa:123456@localhost?database=a004"
			db, err := vdb.Open("sqlserver", mssqlDns)
			if err != nil {
				panic(err)
			}
			return NewUserRepositorySQL(db) // 👈 real repo
		}
		svc.JwtProvider.Init = func(owner *AuthService) *JwtProvider {
			return &JwtProvider{SecretKey: "your-secret-key"}
		}
		svc.Hasher.Init = func(owner *AuthService) Hasher {
			return &Argon2Hasher{
				Time:    1,
				Memory:  16 * 1024, // 16MB → hợp lý hơn
				Threads: 1,
				KeyLen:  32,
				SaltLen: 16,
			}
		}
		svc.UserCache.Init = func(owner *AuthService) *UserRepositoryCached {
			return &UserRepositoryCached{
				db: *svc.UserRepo.Get(),
			}
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
