package orm

type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		complex64 | complex128
}
type NumberField[T Number] struct {
	*dbField
	val        *T
	callMethod *methodCall
}

func (f *NumberField[T]) Get() *T {
	return f.val
}
func (f *NumberField[T]) Set(val *T) {
	f.val = val
}

func (f *NumberField[T]) As(name string) aliasField {
	return *f.dbField.As(name)

}

func (f *NumberField[T]) Eq(other interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "=",
	}
}
func (f *NumberField[T]) Ne(other interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "!=",
	}
}
func (f *NumberField[T]) Gt(other interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      ">",
	}
}
func (f *NumberField[T]) Lt(other interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "<",
	}
}
func (f *NumberField[T]) Ge(other interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      ">=",
	}
}
func (f *NumberField[T]) Le(other interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "<=",
	}
}
func (f *NumberField[T]) IsNull() *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		op:      "IS NULL",
	}
}
func (f *NumberField[T]) IsNotNull() *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		op:      "IS NOT NULL",
	}
}
func (f *NumberField[T]) In(others ...interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   others,
		op:      "IN",
	}
}
func (f *NumberField[T]) NotIn(others ...interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   others,
		op:      "NOT IN",
	}
}
func (f *NumberField[T]) Between(min, max interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   []interface{}{min, max},
		op:      "BETWEEN",
	}
}
func (f *NumberField[T]) NotBetween(min, max interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   []interface{}{min, max},
		op:      "NOT BETWEEN",
	}
}

/*
DateTimeField_test.go:47: method In not found in type *orm.DateTimeField

	DateTimeField_test.go:47: method NotIn not found in type *orm.DateTimeField
	DateTimeField_test.go:47: method IsNull not found in type *orm.DateTimeField
	DateTimeField_test.go:47: method IsNotNull not found in type *orm.DateTimeField
	DateTimeField_test.go:47: method NotBetween not found in type *orm.DateTimeField
	DateTimeField_test.go:47: method Day not found in type *orm.DateTimeField
	DateTimeField_test.go:47: method Month not found in type *orm.DateTimeField
	DateTimeField_test.go:47: method Year not found in type *orm.DateTimeField
	DateTimeField_test.go:47: method Hour not found in type *orm.DateTimeField
	DateTimeField_test.go:47: method Minute not found in type *orm.DateTimeField
	DateTimeField_test.go:47: method Second not found in type *orm.DateTimeField
	DateTimeField_test.go:47: method Format not found in type *orm.DateTimeField
*/
