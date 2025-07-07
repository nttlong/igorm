package orm

func (f *NumberField[T]) makeArithmetic(other interface{}, op string) *fieldBinary {
	return &fieldBinary{
		left:  f,
		right: other,
		op:    op,
	}
}
func (f *NumberField[T]) Add(other interface{}) *NumberField[T] {
	return &NumberField[T]{
		UnderField: f.makeArithmetic(other, "+"),
	}

}
func (f *NumberField[T]) Sub(other interface{}) *NumberField[T] {
	return &NumberField[T]{
		UnderField: f.makeArithmetic(other, "-"),
	}
}
func (f *NumberField[T]) Mul(other interface{}) *NumberField[T] {
	return &NumberField[T]{
		UnderField: f.makeArithmetic(other, "*"),
	}
}
func (f *NumberField[T]) Div(other interface{}) *NumberField[T] {
	return &NumberField[T]{
		UnderField: f.makeArithmetic(other, "/"),
	}
}
func (f *NumberField[T]) Mod(other interface{}) *NumberField[T] {
	return &NumberField[T]{
		UnderField: f.makeArithmetic(other, "%"),
	}
}
func (f *NumberField[T]) Pow(other interface{}) *NumberField[T] {
	return &NumberField[T]{
		UnderField: f.makeArithmetic(other, "^"),
	}
}
