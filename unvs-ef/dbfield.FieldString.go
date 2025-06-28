package unvsef

type FieldString Field[string]

func (f *FieldString) ToSqlExpr(d Dialect) (string, []interface{}) {
	return compiler.Compile(f, d)
}
