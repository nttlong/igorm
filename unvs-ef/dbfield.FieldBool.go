package unvsef

type FieldBool Field[bool]

func (f *FieldBool) As(alias string) *AliasField {
	return &AliasField{
		Field: f,
		Alias: alias,
	}
}
func (where *FieldBool) ToSqlExpr(d Dialect) (string, []interface{}) {
	return compiler.Compile(where, d)
}
func (f Field[TField]) ToSqlExpr2(d Dialect) (string, []interface{}) {
	return compiler.Compile(&f, d)
}
func (f *FuncField) ToSqlExpr(d Dialect) (string, []interface{}) {
	return compiler.Compile(f, d)
}
func (f FuncField) ToSqlExpr2(d Dialect) (string, []interface{}) {
	return compiler.Compile(f, d)
}
