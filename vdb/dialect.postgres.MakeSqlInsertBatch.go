package vdb

import "vdb/migrate"

func (d *postgresDialect) MakeSqlInsertBatch(tableName string, columns []migrate.ColumnDef, data interface{}) (string, []interface{}) {
	panic("not implemented")
}
