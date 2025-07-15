package eorm

import "eorm/migrate"

func (d *mssqlDialect) MakeSqlInsertBatch(tableName string, columns []migrate.ColumnDef, data interface{}) (string, []interface{}) {
	panic("not implemented")
}
