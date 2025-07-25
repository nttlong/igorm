package migrate

import (
	"sync"
	"vdb/tenantDB"
)

type MigratorLoaderMysql struct {
	cacheLoadFullSchema sync.Map
}

func (m *MigratorLoaderMysql) GetDbName(db *tenantDB.TenantDB) string {
	return db.GetDBName()
}

/*
This structure is used to ensure that each database runs only once,
regardless of the Go Routines
*/
type initMySqlLoadFullSchema struct {
	once   sync.Once
	err    error
	schema *DbSchema
}

func (m *MigratorLoaderMysql) LoadFullSchema(db *tenantDB.TenantDB) (*DbSchema, error) {
	cacheKey := db.GetDBName()
	actual, _ := m.cacheLoadFullSchema.LoadOrStore(cacheKey, &initMySqlLoadFullSchema{})
	initSchema := actual.(*initMySqlLoadFullSchema)
	initSchema.once.Do(func() {
		initSchema.schema, initSchema.err = m.loadFullSchema(db)
	})
	return initSchema.schema, initSchema.err
}
func (m *MigratorLoaderMysql) loadFullSchema(db *tenantDB.TenantDB) (*DbSchema, error) {

	tables, err := m.LoadAllTable(db)
	if err != nil {
		return nil, err
	}
	pks, _ := m.LoadAllPrimaryKey(db)
	uks, _ := m.LoadAllUniIndex(db)
	idxs, _ := m.LoadAllIndex(db)

	dbName := m.GetDbName(db)
	schema := &DbSchema{
		DbName:      dbName,
		Tables:      make(map[string]map[string]bool),
		PrimaryKeys: pks,
		UniqueKeys:  uks,
		Indexes:     idxs,
	}
	foreignKeys, err := m.LoadForeignKey(db)
	if err != nil {
		return nil, err
	}
	schema.ForeignKeys = map[string]DbForeignKeyInfo{}
	for _, fk := range foreignKeys {
		schema.ForeignKeys[fk.ConstraintName] = fk
	}
	for table, columns := range tables {
		cols := make(map[string]bool)
		for col := range columns {
			cols[col] = true
		}
		schema.Tables[table] = cols
	}

	return schema, nil
}
