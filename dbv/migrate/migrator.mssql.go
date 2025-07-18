package migrate

import (
	"dbv/tenantDB"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type migratorMssql struct {
	loader             IMigratorLoader
	cacheGetFullScript sync.Map

	db *tenantDB.TenantDB
}

func (m *migratorMssql) Quote(names ...string) string {
	return "[" + strings.Join(names, "].[") + "]"
}

func (m *migratorMssql) GetSqlMigrate(entityType reflect.Type) ([]string, error) {
	scripts := []string{}
	scriptTable, err := m.GetSqlCreateTable(entityType)
	if err != nil {
		return nil, err
	}
	if scriptTable == "" {
		scriptAddColumn, err := m.GetSqlAddColumn(entityType)
		if err != nil {
			return nil, err
		}
		scripts = append(scripts, scriptTable, scriptAddColumn)
	}

	scriptAddUniqueIndex, err := m.GetSqlAddUniqueIndex(entityType)
	if err != nil {
		return nil, err
	}
	scripts = append(scripts, scriptTable, scriptAddUniqueIndex)
	return scripts, nil

}
func (m *migratorMssql) DoMigrate(entityType reflect.Type) error {
	scripts, err := m.GetSqlMigrate(entityType)
	if err != nil {
		return err
	}
	for _, script := range scripts {
		_, err := m.db.Exec(script)
		if err != nil {
			return err
		}
	}
	return nil

}

type mssqlInitDoMigrates struct {
	once sync.Once
	err  error
}

var cacheDoMigrates sync.Map

func (m *migratorMssql) DoMigrates() error {

	key := fmt.Sprintf("%s_%s", m.db.GetDBName(), m.db.GetDbType())
	actual, _ := cacheDoMigrates.LoadOrStore(key, &mssqlInitDoMigrates{})

	mi := actual.(*mssqlInitDoMigrates)

	mi.once.Do(func() {

		scripts, err := m.GetFullScript()
		if err != nil {
			return
		}
		for _, script := range scripts {
			_, err := m.db.Exec(script)
			if err != nil {
				mi.err = err
				break
			}
		}
		// for _, entity := range ModelRegistry.GetAllModels() {
		// 	err = m.DoMigrate(entity.entity.entityType)

		// }
	})
	return mi.err
}

func (m *migratorMssql) GetLoader() IMigratorLoader {
	return m.loader
}
