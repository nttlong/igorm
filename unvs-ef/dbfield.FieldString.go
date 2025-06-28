package unvsef

type FieldString Field[string]

func (f *FieldString) ToSqlExpr(d Dialect) (string, []interface{}) {
	return compiler.Compile(f, d)
}
func (f FieldString) ToSqlExpr2(d Dialect) (string, []interface{}) {
	return compiler.Compile(f, d)
}
func (f *FieldString) Len() *Field[int] {
	return &Field[int]{
		FuncField: &FuncField{
			FuncName: "LEN",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldString) Like(other interface{}) *FieldBool {
	return &FieldBool{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "LIKE",
			Right: other,
		},
	}
}
func (f *Field[TField]) NotLike(other interface{}) *FieldBool {
	return &FieldBool{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "NOT LIKE",
			Right: other,
		},
	}
}
