package services

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct {
	CacheService
}
type cacheVerifyPassword struct {
	Username       string
	Password       string
	HashedPassword string
}

func (p *PasswordService) validate() string {

	typ := reflect.TypeOf(*p)
	if p.Cache == nil {

		panic(fmt.Sprintf("%s.PasswordService.Cache is not nil", typ.Name()))
	}
	if p.Context == nil {
		panic(fmt.Sprintf("%s.PasswordService.Context is not nil", typ.Name()))
	}
	return typ.PkgPath() + "/" + typ.Name()
}
func (p *PasswordService) VerifyPassword(username, password, hashedPassword string) error {

	hasContent := password + "@" + strings.ToLower(username)
	cacheKey := p.validate() + "/" + hasContent
	ret := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(hasContent))
	if ret == nil {
		p.Cache.Set(
			p.Context,
			cacheKey,
			&cacheVerifyPassword{
				Username:       username,
				Password:       password,
				HashedPassword: hashedPassword,
			}, time.Minute*5)
		return nil

	}
	return ret
}
func (p *PasswordService) RemoveCache(username, password string) error {

	hasContent := password + "@" + strings.ToLower(username)
	cacheKey := p.validate() + "/" + hasContent
	p.Cache.Delete(p.Context, cacheKey)
	return nil
}

// hashPasswordWithSalt băm mật khẩu với muối sử dụng bcrypt
func (p *PasswordService) HashPassword(password, username string) (string, error) {
	// Chuyển mật khẩu thành []byte
	passwordBytes := []byte(password + "@" + strings.ToLower(username))

	// Tạo hash với bcrypt, sử dụng cost factor mặc định (10)
	// bcrypt tự động tạo muối ngẫu nhiên
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Chuyển hash thành chuỗi để trả về
	return string(hash), nil
}
