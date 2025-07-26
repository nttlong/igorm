package security

import (
	ctxSvc "vapi/internal/context_service"
	"vapi/internal/security/models"
	"vcache"
	"vdb"
)

type SecurityPolicyService struct {
	Cache  vcache.Cache
	DB     *vdb.TenantDB
	CtxSvc *ctxSvc.ContextService
}

const cachePrefix = "security:policy:"

func NewSecurityPolicyService(cache vcache.Cache, db *vdb.TenantDB) *SecurityPolicyService {
	return &SecurityPolicyService{
		Cache: cache,
		DB:    db,
	}
}

func (s *SecurityPolicyService) GetPolicy(tenantID string) (*models.SecurityPolicy, error) {
	// Check cache
	data := &models.SecurityPolicy{}

	if s.Cache.Get(s.CtxSvc.GetContext(), tenantID, data) {
		return data, nil
	}

	// Query DB

	err := s.DB.First(&data, "tenant_id = ?", tenantID)
	if err != nil {
		return nil, err
	}

	// Cache lại
	s.Cache.Set(s.CtxSvc.GetContext(), tenantID, data, 0)
	return data, nil
}

func (s *SecurityPolicyService) SetJwtSecret(tenantID, newSecret string) error {
	// Cập nhật trong DB
	err := s.DB.Model(&models.SecurityPolicy{}).
		Context(s.CtxSvc.GetContext()).
		Where("tenant_id = ?", tenantID).
		Update("jwt_secret", newSecret).Error
	if err != nil {
		return err
	}

	// Xoá cache để lần sau load lại
	s.Cache.Delete(s.CtxSvc.GetContext(), tenantID)
	return nil
}

func (s *SecurityPolicyService) CreateOrUpdate(policy *models.SecurityPolicy) error {
	// Nếu đã tồn tại thì cập nhật
	existing := &models.SecurityPolicy{}
	err := s.DB.FirstWithContext(s.CtxSvc.GetContext(), existing, "tenant_id = ?", policy.TenantID)
	if err == nil {
		policy.ID = existing.ID
		rs := s.DB.Update(&policy)
		if rs.Error != nil {
			return rs.Error
		}
		return nil
	}
	// Không tồn tại thì tạo mới
	return s.DB.Insert(policy)
}
