package migrate

import (
	"fmt"
	"sync"
)

type postgresGetFullScriptInit struct {
	once sync.Once
	ret  []string
}

func (m *migratorPostgres) GetFullScript() ([]string, error) {
	key := fmt.Sprintf("%s_%s", m.db.GetDBName(), m.db.GetDbType())
	actual, _ := m.cacheGetFullScript.LoadOrStore(key, &postgresGetFullScriptInit{})
	init := actual.(*postgresGetFullScriptInit)
	var err error
	init.once.Do(func() {
		init.ret, err = m.getFullScript()
	})
	return init.ret, err
}
func (m *migratorPostgres) getFullScript() ([]string, error) {
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
	scriptForeignKey, err := m.GetSqlAddForeignKey()
	if err != nil {
		return nil, err
	}
	scripts = append(scripts, scriptForeignKey...)

	return scripts, nil
}
