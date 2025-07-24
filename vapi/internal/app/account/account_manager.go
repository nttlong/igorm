package account

import (
	"sync"
	"time"
	"vcache"
	unvsdi "vdi"
)

type AccountScope struct {
	Container  unvsdi.Container
	LastAccess time.Time
}

type AccountManager struct {
	accountMap sync.Map // map[string]*AccountScope
}

func NewAccountManager() *AccountManager {
	return &AccountManager{}
}

func (m *AccountManager) getOrCreateScope(accountID string) *AccountScope {
	scope := &AccountScope{
		Container:  unvsdi.NewScopedContainer(),
		LastAccess: time.Now(),
	}
	return scope
}

func (m *AccountManager) Login(accountID string) {
	scope := m.getOrCreateScope(accountID)

	// Register cache for this user (could use Redis, Memcached, etc.)
	scope.Container.RegisterScoped(func(c unvsdi.Container) vcache.Cache {
		return vcache.NewInMemoryCache(10*time.Minute, 1*time.Minute)
	})

	m.accountMap.Store(accountID, scope)
}

func (m *AccountManager) Logout(accountID string) {
	if scope, ok := m.accountMap.Load(accountID); ok {
		scope.(*AccountScope).Container.Dispose()
		m.accountMap.Delete(accountID)
	}
}

func (m *AccountManager) GetScope(accountID string) *AccountScope {
	if v, ok := m.accountMap.Load(accountID); ok {
		scope := v.(*AccountScope)
		scope.LastAccess = time.Now()
		return scope
	}
	return nil
}
