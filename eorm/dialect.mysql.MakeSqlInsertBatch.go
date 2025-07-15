package eorm

import "eorm/migrate"

func (d *mySqlDialect) MakeSqlInsertBatch(tableName string, columns []migrate.ColumnDef, data interface{}) (string, []interface{}) {
	panic("not implemented")
}
