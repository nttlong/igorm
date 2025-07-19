package migrate

import (
	"fmt"
	"sync"
)

type mysqlGetFullScriptInit struct {
	once sync.Once
	ret  []string
}

func (m *migratorMySql) GetFullScript() ([]string, error) {
	key := fmt.Sprintf("%s_%s", m.db.GetDBName(), m.db.GetDbType())
	actual, _ := m.cacheGetFullScript.LoadOrStore(key, &mysqlGetFullScriptInit{})
	init := actual.(*mysqlGetFullScriptInit)
	var err error
	init.once.Do(func() {
		init.ret, err = m.getFullScript()
	})
	return init.ret, err
}
func (m *migratorMySql) getFullScript() ([]string, error) {
	sqlInstall, err := m.GetSqlInstallDb()
	if err != nil {
		return nil, err
	}
	scripts := sqlInstall
	for _, entity := range ModelRegistry.GetAllModels() {
		script, err := m.GetSqlCreateTable(entity.entity.entityType)
		if err != nil {
			return nil, err
		}
		if script != "" {
			scripts = append(scripts, script)
		}

	}
	for _, entity := range ModelRegistry.GetAllModels() {
		script, err := m.GetSqlAddColumn(entity.entity.entityType)
		if err != nil {
			return nil, err
		}
		if script != "" {
			scripts = append(scripts, script)
		}
	}
	for _, entity := range ModelRegistry.GetAllModels() {
		script, err := m.GetSqlAddIndex(entity.entity.entityType)
		if err != nil {
			return nil, err
		}
		if script != "" {
			scripts = append(scripts, script)
		}
	}
	for _, entity := range ModelRegistry.GetAllModels() {
		script, err := m.GetSqlAddUniqueIndex(entity.entity.entityType)
		if err != nil {
			return nil, err
		}
		if script != "" {
			scripts = append(scripts, script)
		}
	}
	scriptForeignKey, err := m.GetSqlAddForeignKey()
	if err != nil {
		return nil, err
	}
	scripts = append(scripts, scriptForeignKey...)

	return scripts, nil
}
