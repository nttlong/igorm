package migrate

import (
	"eorm/tenantDB"
	"sync"
)

type MigratorLoaderMssql struct {
	cacheLoadFullSchema sync.Map
}

func (m *MigratorLoaderMssql) GetDbName(db *tenantDB.TenantDB) string {
	var dbName string
	err := db.QueryRow("SELECT DB_NAME()").Scan(&dbName)
	if err != nil {
		return ""
	}
	return dbName
}

func (m *MigratorLoaderMssql) LoadAllTable(db *tenantDB.TenantDB) (map[string]map[string]ColumnInfo, error) {
	query := `
	SELECT
		t.name AS TableName,
		c.name AS ColumnName,
		ty.name AS DataType,
		c.is_nullable,
		c.max_length
	FROM sys.columns c
	JOIN sys.tables t ON c.object_id = t.object_id
	JOIN sys.types ty ON c.user_type_id = ty.user_type_id`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tables := make(map[string]map[string]ColumnInfo)
	for rows.Next() {
		var table, column, dbType string
		var nullable bool
		var length int
		if err := rows.Scan(&table, &column, &dbType, &nullable, &length); err != nil {
			return nil, err
		}
		if _, ok := tables[table]; !ok {
			tables[table] = make(map[string]ColumnInfo)
		}
		tables[table][column] = ColumnInfo{
			Name:     column,
			DbType:   dbType,
			Nullable: nullable,
			Length:   length,
		}
	}
	return tables, nil
}

func (m *MigratorLoaderMssql) LoadAllPrimaryKey(db *tenantDB.TenantDB) (map[string]ColumnsInfo, error) {
	query := `
	SELECT
		KCU.table_name,
		KCU.column_name,
		TC.constraint_name
	FROM information_schema.table_constraints AS TC
	JOIN information_schema.key_column_usage AS KCU
		ON TC.constraint_name = KCU.constraint_name
	WHERE TC.constraint_type = 'PRIMARY KEY'`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]ColumnsInfo)
	for rows.Next() {
		var table, column, constraint string
		if err := rows.Scan(&table, &column, &constraint); err != nil {
			return nil, err
		}
		info := result[constraint]
		info.TableName = table
		info.Columns = append(info.Columns, ColumnInfo{Name: column})
		result[constraint] = info
	}
	return result, nil
}

func (m *MigratorLoaderMssql) LoadAllUniIndex(db *tenantDB.TenantDB) (map[string]ColumnsInfo, error) {
	query := `
	SELECT
		t.name AS TableName,
		i.name AS IndexName,
		c.name AS ColumnName
	FROM sys.indexes i
	JOIN sys.tables t ON i.object_id = t.object_id
	JOIN sys.index_columns ic ON i.object_id = ic.object_id AND i.index_id = ic.index_id
	JOIN sys.columns c ON ic.object_id = c.object_id AND ic.column_id = c.column_id
	WHERE i.type_desc = 'NONCLUSTERED' AND is_unique_constraint = 1`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]ColumnsInfo)
	for rows.Next() {
		var table, index, column string
		if err := rows.Scan(&table, &index, &column); err != nil {
			return nil, err
		}
		info := result[index]
		info.TableName = table
		info.Columns = append(info.Columns, ColumnInfo{Name: column})
		result[index] = info
	}
	return result, nil
}

func (m *MigratorLoaderMssql) LoadAllIndex(db *tenantDB.TenantDB) (map[string]ColumnsInfo, error) {
	query := `
	SELECT
		t.name AS TableName,
		i.name AS IndexName,
		c.name AS ColumnName
	FROM sys.indexes i
	JOIN sys.tables t ON i.object_id = t.object_id
	JOIN sys.index_columns ic ON i.object_id = ic.object_id AND i.index_id = ic.index_id
	JOIN sys.columns c ON ic.object_id = c.object_id AND ic.column_id = c.column_id
	WHERE i.is_primary_key = 0 AND is_unique_constraint = 0`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]ColumnsInfo)
	for rows.Next() {
		var table, index, column string
		if err := rows.Scan(&table, &index, &column); err != nil {
			return nil, err
		}
		info := result[index]
		info.TableName = table
		info.Columns = append(info.Columns, ColumnInfo{Name: column})
		result[index] = info
	}
	return result, nil
}

func (m *MigratorLoaderMssql) LoadFullSchema(db *tenantDB.TenantDB) (*DbSchema, error) {
	cacheKey := db.GetDBName()
	if val, ok := m.cacheLoadFullSchema.Load(cacheKey); ok {
		return val.(*DbSchema), nil
	}
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
	for table, columns := range tables {
		cols := make(map[string]bool)
		for col := range columns {
			cols[col] = true
		}
		schema.Tables[table] = cols
	}
	m.cacheLoadFullSchema.Store(cacheKey, schema)
	return schema, nil
}
