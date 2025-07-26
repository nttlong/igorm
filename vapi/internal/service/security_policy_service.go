package service

// import (
// 	"context"
// 	"time"
// 	"vcache"
// 	"vdb"
// )

// // Cấu trúc lưu trữ chính sách bảo mật
// type SecurityPolicy struct {
// 	TenantID         string
// 	MaxLoginFailures int
// 	LockoutDuration  time.Duration
// 	JwtSecretKey     string
// 	JwtExpiration    time.Duration
// 	LastUpdated      time.Time
// }



// func NewSecurityPolicyService(db *vdb.TenantDB, cache vcache.Cache) *SecurityPolicyService {
// 	return &SecurityPolicyService{
// 		db:    db,
// 		cache: cache,
// 	}
// }

// func cacheKey(tenantID string) string {
// 	return "secpol:" + tenantID
// }

// // Lấy chính sách bảo mật của tenant (ưu tiên cache)
// func (s *SecurityPolicyService) GetPolicy(ctx context.Context, tenantID string) (*SecurityPolicy, error) {
// 	key := cacheKey(tenantID)

// 	var policy SecurityPolicy
// 	found, err := s.cache.Get(ctx, key, &policy)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if found {
// 		return &policy, nil
// 	}

// 	// Giả lập lấy từ DB (bạn có thể thay bằng truy vấn thực)
// 	row := s.db.QueryRowContext(ctx, `
// 		SELECT max_login_failures, lockout_minutes, jwt_secret, jwt_expiration_minute
// 		FROM security_policies
// 		WHERE tenant_id = ?
// 	`, tenantID)

// 	err = row.Scan(&policy.MaxLoginFailures, &policy.LockoutDuration, &policy.JwtSecretKey, &policy.JwtExpiration)
// 	if err != nil {
// 		return nil, err
// 	}

// 	policy.TenantID = tenantID
// 	policy.LastUpdated = time.Now()

// 	// Lưu vào cache
// 	_ = s.cache.Set(ctx, key, policy, 5*time.Minute)
// 	return &policy, nil
// }

// // Tạo hoặc cập nhật chính sách bảo mật
// func (s *SecurityPolicyService) SavePolicy(ctx context.Context, policy *SecurityPolicy) error {
// 	// Cập nhật DB (giả lập)
// 	_, err := s.db.ExecContext(ctx, `
// 		INSERT INTO security_policies
// 			(tenant_id, max_login_failures, lockout_minutes, jwt_secret, jwt_expiration_minute)
// 		VALUES (?, ?, ?, ?, ?)
// 		ON DUPLICATE KEY UPDATE
// 			max_login_failures = VALUES(max_login_failures),
// 			lockout_minutes = VALUES(lockout_minutes),
// 			jwt_secret = VALUES(jwt_secret),
// 			jwt_expiration_minute = VALUES(jwt_expiration_minute)
// 	`, policy.TenantID, policy.MaxLoginFailures, int(policy.LockoutDuration.Minutes()), policy.JwtSecretKey, int(policy.JwtExpiration.Minutes()))
// 	if err != nil {
// 		return err
// 	}

// 	// Cập nhật cache
// 	return s.cache.Set(ctx, cacheKey(policy.TenantID), policy, 5*time.Minute)
// }

// // Xóa secret key (reset chính sách)
// func (s *SecurityPolicyService) DeletePolicy(ctx context.Context, tenantID string) error {
// 	_, err := s.db.ExecContext(ctx, `
// 		DELETE FROM security_policies WHERE tenant_id = ?
// 	`, tenantID)
// 	if err != nil {
// 		return err
// 	}
// 	return s.cache.Delete(ctx, cacheKey(tenantID))
// }
