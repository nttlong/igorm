package migrate

import (
	"eorm/tenantDB"
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

func (m *migratorPostgres) GetColumnDataTypeMapping() map[reflect.Type]string {
	panic("not implemented")
}
func (m *migratorPostgres) GetGetDefaultValueByFromDbTag() map[string]string {
	panic("not implemented")
}
func (m *migratorPostgres) GetSqlCreateTable(entityType reflect.Type) (string, error) {
	panic("not implemented")
}
func (m *migratorPostgres) GetSqlAddColumn(entityType reflect.Type) (string, error) {
	panic("not implemented")
}
func (m *migratorPostgres) GetSqlAddIndex(entityType reflect.Type) (string, error) {
	panic("not implemented")
}
func (m *migratorPostgres) GetSqlAddUniqueIndex(entityType reflect.Type) (string, error) {
	panic("not implemented")
}
func (m *migratorPostgres) GetSqlMigrate(entityType reflect.Type) ([]string, error) {
	panic("not implemented")
}
func (m *migratorPostgres) GetSqlAddForeignKey() ([]string, error) {
	panic("not implemented")
}
func (m *migratorPostgres) GetFullScript() ([]string, error) {
	panic("not implemented")
}
func (m *migratorPostgres) DoMigrate(entityType reflect.Type) error {
	panic("not implemented")
}
func (m *migratorPostgres) DoMigrates() error {
	panic("not implemented")
}
