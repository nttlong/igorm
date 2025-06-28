package unvsef

type FieldBool Field[bool]

func (where *FieldBool) ToSqlExpr(d Dialect) (string, []interface{}) {
	return compiler.Compile(where, d)
}
func (f Field[TField]) ToSqlExpr2(d Dialect) (string, []interface{}) {
	return (&f).ToSqlExpr(d)
}
func (f *FuncField) ToSqlExpr(d Dialect) (string, []interface{}) {
	return compiler.Compile(f, d)
}
