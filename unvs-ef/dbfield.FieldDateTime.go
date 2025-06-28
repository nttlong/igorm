package unvsef

import "time"

type FieldDateTime Field[time.Time]

func (f *FieldDateTime) ToSqlExpr(d Dialect) (string, []interface{}) {
	return compiler.Compile(f, d)
}
