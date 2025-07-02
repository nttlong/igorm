package unvsef

type NumberValue interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		complex64 | complex128
}
type FieldNumber[TField NumberValue] struct {
	*DbField
	*AliasField
	*BinaryField
	*FuncField
	Op  string
	val *TField
}

func (f *FieldNumber[TField]) ToSqlExpr(d Dialect) (string, []interface{}) {
	return compiler.Compile(f, d)
}
func (f FieldNumber[TField]) ToSqlExpr2(d Dialect) (string, []interface{}) {
	return compiler.Compile(f, d)
}
func (f *FieldNumber[TField]) Set(val *TField) {
	f.val = val
}
func (f *FieldNumber[TField]) Desc() *SortField {
	return &SortField{
		Field: f,
		Sort:  "ASC",
	}

}
func (f *FieldNumber[TField]) Asc() *SortField {
	return &SortField{
		Field: f,
		Sort:  "ASC",
	}

}
func (f *FieldNumber[TField]) Eq(val interface{}) *FieldBool {
	return &FieldBool{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "=",
			Right: val,
		},
	}

}
func (f *FieldNumber[TField]) Ne(val interface{}) *FieldBool {
	return &FieldBool{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "!=",
			Right: val,
		},
	}

}
func (f *FieldNumber[TField]) Gt(val interface{}) *FieldBool {
	return &FieldBool{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    ">",
			Right: val,
		},
	}

}
func (f *FieldNumber[TField]) Lt(val interface{}) *FieldBool {
	return &FieldBool{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "<",
			Right: val,
		},
	}

}
func (f *FieldNumber[TField]) Ge(val interface{}) *FieldBool {
	return &FieldBool{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    ">=",
			Right: val,
		},
	}

}
func (f *FieldNumber[TField]) Le(val interface{}) *FieldBool {
	return &FieldBool{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "<=",
			Right: val,
		},
	}

}
func (f *FieldNumber[TField]) In(vals ...interface{}) *FieldBool {
	return &FieldBool{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "IN",
			Right: vals,
		},
	}

}
func (f *FieldNumber[TField]) NotIn(vals ...interface{}) *FieldBool {
	return &FieldBool{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "NOT IN",
			Right: vals,
		},
	}

}
func (f *FieldNumber[TField]) IsNull() *FieldBool {
	return &FieldBool{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "IS NULL",
			Right: nil,
		},
	}

}
func (f *FieldNumber[TField]) IsNotNull() *FieldBool {
	return &FieldBool{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "IS NOT NULL",
			Right: nil,
		},
	}

}
func (f *FieldNumber[TField]) Add(val interface{}) *FieldNumber[TField] {
	return &FieldNumber[TField]{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "+",
			Right: val,
		},
	}

}
func (f *FieldNumber[TField]) Sub(val interface{}) *FieldNumber[TField] {
	return &FieldNumber[TField]{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "-",
			Right: val,
		},
	}

}
func (f *FieldNumber[TField]) Mul(val interface{}) *FieldNumber[TField] {
	return &FieldNumber[TField]{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "*",
			Right: val,
		},
	}

}
func (f *FieldNumber[TField]) Div(val interface{}) *FieldNumber[TField] {
	return &FieldNumber[TField]{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "/",
			Right: val,
		},
	}

}
func (f *FieldNumber[TField]) Mod(val interface{}) *FieldNumber[TField] {
	return &FieldNumber[TField]{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "%",
			Right: val,
		},
	}

}
func (f *FieldNumber[TField]) As(alias string) *AliasField {
	return &AliasField{
		Field: *f.DbField,
		Alias: alias,
	}
}
