package dbv

import "dbv/migrate"

type QueryBuilder struct {
	source string
	cols   []migrate.ColumnDef
}

func (m *Model[T]) query() *QueryBuilder {

	return &QueryBuilder{
		source: m.Entity.TableName(),
		cols:   m.Entity.GetColumns(),
	}
}
func (m *Model[T]) SelectAll() *QueryBuilder {
	return m.query()
}
