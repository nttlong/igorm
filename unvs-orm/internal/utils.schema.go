package internal

import "database/sql"

func (u *utilsPackage) extractSchema(db *sql.DB, dbName string, dialect Dialect) (*schemaMap, error) {

	dialect.RefreshSchemaCache(db, dbName)
	schema, err := dialect.GetSchema(db, dbName)
	if err != nil {
		return nil, err
	}
	ret := &schemaMap{
		table:  make(map[string]bool),
		unique: make(map[string]bool),
		index:  make(map[string]bool),
		fk:     make(map[string]bool),
	}
	for tableName, table := range schema {
		ret.table[tableName] = true

		for _, constraintName := range table.UniqueConstraints {
			ret.unique[constraintName] = true
		}
		for _, constraintName := range table.IndexConstraints {
			ret.index[constraintName] = true
		}
	}
	u.schemaCache.Store(dbName, ret)
	return ret, nil

}
