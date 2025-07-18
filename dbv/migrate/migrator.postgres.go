package migrate

import (
	"dbv/tenantDB"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type migratorPostgres struct {
	loader             IMigratorLoader
	cacheGetFullScript sync.Map

	db *tenantDB.TenantDB
}

func (m *migratorPostgres) GetLoader() IMigratorLoader {
	return m.loader
}
func (m *migratorPostgres) Quote(names ...string) string {
	return "\"" + strings.Join(names, "\".\"") + "\""
}

func (m *migratorPostgres) GetSqlMigrate(entityType reflect.Type) ([]string, error) {
	panic("not implemented")
}

func (m *migratorPostgres) DoMigrate(entityType reflect.Type) error {
	panic("not implemented")
}

type postgresInitDoMigrates struct {
	once sync.Once
}

func (m *migratorPostgres) DoMigrates() error {
	key := fmt.Sprintf("%s_%s", m.db.GetDBName(), m.db.GetDbType())
	actual, _ := cacheDoMigrates.LoadOrStore(key, &postgresInitDoMigrates{})

	mi := actual.(*postgresInitDoMigrates)
	var err error
	mi.once.Do(func() {

		scripts, err := m.GetFullScript()
		if err != nil {
			return
		}
		for _, script := range scripts {
			_, err := m.db.Exec(script)
			if err != nil {
				break
			}
		}

	})
	return err
}
