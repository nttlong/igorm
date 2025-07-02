package orm

func (f *NumberField[T]) Add(other interface{}) *fieldBinary {
	return &fieldBinary{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "+",
	}
}
func (f *NumberField[T]) Sub(other interface{}) *fieldBinary {
	return &fieldBinary{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "-",
	}
}
func (f *NumberField[T]) Mul(other interface{}) *fieldBinary {
	return &fieldBinary{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "*",
	}
}
func (f *NumberField[T]) Div(other interface{}) *fieldBinary {
	return &fieldBinary{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "/",
	}
}
func (f *NumberField[T]) Mod(other interface{}) *fieldBinary {
	return &fieldBinary{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "%",
	}
}
func (f *NumberField[T]) Pow(other interface{}) *fieldBinary {
	return &fieldBinary{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "^",
	}
}
