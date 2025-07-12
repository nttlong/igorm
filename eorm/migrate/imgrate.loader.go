package migrate

import (
	"database/sql"
	"eorm/tenantDB"
	"fmt"
)

type ColumnInfo struct {
	/*
		Db field name
	*/
	Name string
	/*
		Go field name
	*/
	DbType string
	/*
		Is allow null?
	*/
	Nullable bool
	/*
		Length is the length of the column
	*/
	Length int
}
type ColumnsInfo struct {
	TableName string
	Columns   []ColumnInfo
}
type DbSchema struct {
	/*
		Database name
	*/
	DbName string
	/*
		map[<table name>]map[<column name>]bool
	*/
	Tables map[string]map[string]bool
	/*
		map[<primary key constraint name>]ColumnsInfo
	*/
	PrimaryKeys map[string]ColumnsInfo
	/*
		map[<Unique Keys constraint name>]ColumnsInfo
	*/
	UniqueKeys map[string]ColumnsInfo
	/*
		map[<Index name>]ColumnsInfo
	*/
	Indexes map[string]ColumnsInfo
}
type IMigratorLoader interface {
	GetError() error

	LoadAllTable(db *sql.DB) (map[string]map[string]ColumnInfo, error)
	LoadAllPrimaryKey(db *sql.DB) (map[string]ColumnsInfo, error)
	/*
		Heed: for SQL Server, we need to use the following query to get the unique keys:
			SELECT
			t.name AS TableName,
			i.name AS IndexName
			FROM sys.indexes i
			JOIN sys.tables t ON i.object_id = t.object_id
			WHERE i.type_desc = 'NONCLUSTERED' and is_unique_constraint=1
	*/
	LoadAllUniIndex(db *sql.DB) (map[string]ColumnsInfo, error)
	/*

	 */
	LoadAllIndex(db *sql.DB) (map[string]ColumnsInfo, error)
	LoadFullSchema(db *sql.DB) (*DbSchema, error)
}

func MigratorLoader(db *tenantDB.TenantDB) (IMigratorLoader, error) {
	err := db.Detect()
	if err != nil {
		return nil, err
	}
	switch db.DbType {
	case tenantDB.DB_DRIVER_MSSQL:
		return &MigratorLoaderMssql{}, nil
	default:
		panic(fmt.Errorf("unsupported database type: %s", db.DbType))
	}
}
