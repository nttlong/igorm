package migrate

import (
	"eorm/tenantDB"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type migratorMySql struct {
	loader             IMigratorLoader
	cacheGetFullScript sync.Map

	db *tenantDB.TenantDB
}

func (m *migratorMySql) GetLoader() IMigratorLoader {
	return m.loader
}
func (m *migratorMySql) Quote(names ...string) string {
	return "`" + strings.Join(names, "`.`") + "`"
}

func (m *migratorMySql) GetSqlMigrate(entityType reflect.Type) ([]string, error) {
	panic("implement me")
}

func (m *migratorMySql) DoMigrate(entityType reflect.Type) error {
	panic("implement me")
}

func (m *migratorMySql) DoMigrates() error {

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
				mi.err = createError(script, err)
				break
			}
		}

	})
	return mi.err
}
