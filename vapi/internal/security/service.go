package security

import (
	"context"
	"vapi/internal/security/models"
	"vcache"
	"vdb"
)

type SecurityPolicyService struct {
	Cache  vcache.Cache
	GetDb  func() *vdb.TenantDB
	GetCtx func() context.Context
}

func (s *SecurityPolicyService) GetPolicy(tenantID string) (*models.SecurityPolicy, error) {
	// Check cache
	data := &models.SecurityPolicy{}

	if s.Cache.Get(s.GetCtx(), tenantID, data) {
		return data, nil
	}

	// Query DB

	err := s.GetDb().First(&data, "tenant_id = ?", tenantID)
	if err != nil {
		return nil, err
	}

	// Cache lại
	s.Cache.Set(s.GetCtx(), tenantID, data, 0)
	return data, nil
}

func (s *SecurityPolicyService) SetJwtSecret(tenantID, newSecret string) error {
	// Cập nhật trong DB
	err := s.GetDb().Model(&models.SecurityPolicy{}).
		Context(s.GetCtx()).
		Where("tenant_id = ?", tenantID).
		Update("jwt_secret", newSecret).Error
	if err != nil {
		return err
	}

	// Xoá cache để lần sau load lại
	s.Cache.Delete(s.GetCtx(), tenantID)
	return nil
}

func (s *SecurityPolicyService) CreateOrUpdate(policy *models.SecurityPolicy) error {
	ctx := s.GetCtx()
	existing := &models.SecurityPolicy{}

	err := s.GetDb().FirstWithContext(ctx, existing, "tenant_id = ?", policy.TenantID)
	if err != nil {
		if _, ok := err.(*vdb.ErrRecordNotFound); ok {
			// Không tồn tại -> tạo mới
			return s.GetDb().Insert(policy)
		}
		// Lỗi khác
		return err
	}

	// Đã tồn tại -> cập nhật
	policy.ID = existing.ID
	if rs := s.GetDb().Update(policy); rs.Error != nil {
		return rs.Error
	}
	return nil
}
func (s *SecurityPolicyService) Get() (*models.SecurityPolicy, error) {
	ctx := s.GetCtx()

	policy := &models.SecurityPolicy{}
	err := s.GetDb().FirstWithContext(ctx, policy)
	if err != nil {
		return nil, err
	}
	return policy, nil
}
