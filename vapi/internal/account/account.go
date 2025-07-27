package account

import (
	"context"
	"errors"
	"fmt"
	"time"
	ps "vapi/internal/PasswordService"
	"vapi/internal/account/models"
	jwtservice "vapi/internal/jwt_service"
	"vapi/internal/security"
	"vcache"
	"vdb"
)

type AccountService struct {
	Cache          vcache.Cache
	GetDb          func() *vdb.TenantDB
	GetCtx         func() context.Context
	PasswordHasher *ps.BcryptPasswordService

	PolicySvc *security.SecurityPolicyService
	JwtSvc    jwtservice.JWTService
}

func (s *AccountService) CreateOrUpdate(account *models.Account) error {
	ctx := s.GetCtx()
	existing := &models.Account{}
	err := s.GetDb().FirstWithContext(ctx, existing, "username = ?", account.Username)
	if err != nil {
		var nfErr *vdb.ErrRecordNotFound
		if errors.As(err, &nfErr) {
			return s.GetDb().Insert(account)
		}
		return err
	}
	account.ID = existing.ID
	return s.GetDb().Update(account).Error
}

type LoginResult struct {
	AccountID int
	Token     string
	PolicySvc security.SecurityPolicyService
}
type ErrInvalidCredentials struct {
	Err error
}
type ErrAccountLocked struct {
	Err error
}

func (e *ErrInvalidCredentials) Error() string {
	return "invalid credentials"
}
func (e *ErrAccountLocked) Error() string {
	return "account locked"
}

func (s *AccountService) Login(tenantID, username, password string) (*LoginResult, error) {
	ctx := s.GetCtx()

	// 1. Tìm account
	account := &models.Account{}
	err := s.GetDb().FirstWithContext(ctx, account, "tenant_id = ? AND username = ?", tenantID, username)
	if err != nil {
		if errors.As(err, &vdb.ErrRecordNotFound{}) {
			return nil, &ErrInvalidCredentials{Err: err}
		}
		return nil, err
	}

	// 2. Kiểm tra có bị khóa không
	lockKey := fmt.Sprintf("lock:%s:%s", tenantID, username)

	if locked, _ := s.Cache.GetBool(ctx, lockKey); locked {
		return nil, &ErrAccountLocked{Err: err}
	}

	// 3. So sánh mật khẩu
	if !s.PasswordHasher.Verify(password, account.HashedPassword) {
		// 4. Tăng số lần sai
		// failKey := fmt.Sprintf("fail:%s:%s", tenantID, username)
		// failures, _ := s.Cache.Incr(ctx, failKey)
		// s.Cache.Expire(ctx, failKey, time.Minute*30) // tự hết hạn nếu user không cố thêm

		// 5. Lấy policy để kiểm tra số lần cho phép
		policy, err := s.PolicySvc.Get()
		// if err == nil && failures >= int64(policy.MaxLoginFailures) {
		// 	s.Cache.SetBool(ctx, lockKey, true, time.Minute*time.Duration(policy.LockoutMinutes))
		// 	return nil, &ErrAccountLocked{Err: fmt.Errorf("account locked for %d minutes", policy.LockoutMinutes)}
		// }
		if err == nil {
			s.Cache.Set(ctx, lockKey, true, time.Minute*time.Duration(policy.LockoutMinutes))
			return nil, &ErrAccountLocked{Err: fmt.Errorf("account locked for %d minutes", policy.LockoutMinutes)}
		}
		return nil, &ErrInvalidCredentials{
			Err: fmt.Errorf("invalid password"),
		}
	}

	// 6. Đúng mật khẩu → reset số lần sai
	failKey := fmt.Sprintf("fail:%s:%s", tenantID, username)
	s.Cache.Delete(ctx, failKey)

	// 7. Tạo JWT
	policy, err := s.PolicySvc.Get()
	if err != nil {
		return nil, err
	}
	token, err := s.JwtSvc.GenerateToken(account.ID, tenantID, policy.JwtSecret, time.Minute*time.Duration(policy.JwtExpireMinutes))
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccountID: account.ID,
		Token:     token,
	}, nil
}
