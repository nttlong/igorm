package unvsef

type FieldUint64 Field[uint64]

func (f *FieldUint64) Eq(other interface{}) *BinaryField {
	return &BinaryField{
		Left:  f,
		Op:    "==",
		Right: other,
	}
}
func (f *FieldUint64) Gt(other interface{}) *BinaryField {
	return &BinaryField{
		Left:  f,
		Op:    ">",
		Right: other,
	}
}
func (f *FieldUint64) Lt(other interface{}) *BinaryField {
	return &BinaryField{
		Left:  f,
		Op:    "<",
		Right: other,
	}

}
func (f *FieldUint64) Gte(other interface{}) *BinaryField {
	return &BinaryField{
		Left:  f,
		Op:    ">=",
		Right: other,
	}
}
func (f *FieldUint64) Lte(other interface{}) *BinaryField {
	return &BinaryField{
		Left:  f,
		Op:    "<=",
		Right: other,
	}
}
func (f *FieldUint64) Add(other interface{}) *BinaryField {
	return &BinaryField{
		Left:  f,
		Op:    "+",
		Right: other,
	}
}
func (f *FieldUint64) Sub(other interface{}) *BinaryField {
	return &BinaryField{
		Left:  f,
		Op:    "-",
		Right: other,
	}

}
func (f *FieldUint64) Mul(other interface{}) *BinaryField {
	return &BinaryField{
		Left:  f,
		Op:    "*",
		Right: other,
	}

}
func (f *FieldUint64) Div(other interface{}) *BinaryField {
	return &BinaryField{
		Left:  f,
		Op:    "/",
		Right: other,
	}

}
func (f *FieldUint64) Mod(other interface{}) *BinaryField {
	return &BinaryField{
		Left:  f,
		Op:    "%",
		Right: other,
	}

}
func (f *FieldUint64) Sum() *FuncField {
	return &FuncField{
		FuncName: "SUM",
		Args:     []interface{}{f},
	}
}
func (f *FieldUint64) Count() *FuncField {
	return &FuncField{
		FuncName: "COUNT",
		Args:     []interface{}{f},
	}
}
func (f *FieldUint64) Avg() *FuncField {
	return &FuncField{
		FuncName: "AVG",
		Args:     []interface{}{f},
	}
}
func (f *FieldUint64) Min() *FuncField {
	return &FuncField{
		FuncName: "MIN",
		Args:     []interface{}{f},
	}
}
func (f *FieldUint64) Max() *FuncField {
	return &FuncField{
		FuncName: "MAX",
		Args:     []interface{}{f},
	}
}

func (f *FieldUint64) ToSqlExpr(d Dialect) (string, []interface{}) {
	return compiler.Compile(f, d)
}
func (f FieldUint64) ToSqlExpr2(d Dialect) (string, []interface{}) {
	return compiler.Compile(f, d)
}
