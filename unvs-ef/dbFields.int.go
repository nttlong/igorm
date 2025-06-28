package unvsef

func (f *Field[TField]) ToSqlExpr(d Dialect) (string, []interface{}) {
	return compiler.Compile(f, d)
}
