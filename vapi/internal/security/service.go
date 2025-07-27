package security

import (
	"context"
	"vapi/internal/security/models"
	"vcache"
	"vdb"
)

type SecurityPolicyService struct {
	Cache vcache.Cache
	Db    *vdb.TenantDB
	Ctx   context.Context
}

func NewSecurityPolicyService(
	ctx context.Context,
	cache vcache.Cache,
	db *vdb.TenantDB,
) *SecurityPolicyService {
	return &SecurityPolicyService{
		Cache: cache,
		Ctx:   ctx,
		Db:    db,
	}

}

func (s *SecurityPolicyService) GetPolicy(tenantID string) (*models.SecurityPolicy, error) {
	// Check cache
	data := &models.SecurityPolicy{}

	if s.Cache.Get(s.Ctx, tenantID, data) {
		return data, nil
	}

	// Query DB

	err := s.Db.First(&data, "tenant_id = ?", tenantID)
	if err != nil {
		return nil, err
	}

	// Cache lại
	s.Cache.Set(s.Ctx, tenantID, data, 0)
	return data, nil
}

func (s *SecurityPolicyService) SetJwtSecret(tenantID, newSecret string) error {
	// Cập nhật trong DB
	err := s.Db.Model(&models.SecurityPolicy{}).
		Context(s.Ctx).
		Where("tenant_id = ?", tenantID).
		Update("jwt_secret", newSecret).Error
	if err != nil {
		return err
	}

	// Xoá cache để lần sau load lại
	s.Cache.Delete(s.Ctx, tenantID)
	return nil
}

func (s *SecurityPolicyService) CreateOrUpdate(policy *models.SecurityPolicy) error {
	ctx := s.Ctx
	existing := &models.SecurityPolicy{}

	err := s.Db.FirstWithContext(ctx, existing, "tenant_id = ?", policy.TenantID)
	if err != nil {
		if _, ok := err.(*vdb.ErrRecordNotFound); ok {
			// Không tồn tại -> tạo mới
			return s.Db.Insert(policy)
		}
		// Lỗi khác
		return err
	}

	// Đã tồn tại -> cập nhật
	policy.ID = existing.ID
	if rs := s.Db.Update(policy); rs.Error != nil {
		return rs.Error
	}
	return nil
}
