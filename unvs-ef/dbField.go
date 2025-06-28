package unvsef

type BinaryField struct {
	Left  interface{}
	Op    string
	Right interface{}
}
type AliasField struct {
	Field interface{}
	Alias string
}
type FuncField struct {
	FuncName string
	Args     []interface{}
}
type DbField struct {
	TableName string
	ColName   string
}

//	type ValueField[TField any] struct {
//		Value TField
//	}
type Field[TField any] struct {
	*DbField
	*AliasField
	*BinaryField
	*FuncField
	Op string
	// *ValueField[TField]
}

func (f *Field[TField]) Eq(other interface{}) *FieldBool {
	return &FieldBool{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "=",
			Right: other,
		},
	}
}
func (f *Field[TField]) Ne(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "!=",
			Right: other,
		},
	}
}
func (f *Field[TField]) Gt(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    ">",
			Right: other,
		},
	}
}
func (f *Field[TField]) Lt(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "<",
			Right: other,
		},
	}
}
func (f *Field[TField]) Gte(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    ">=",
			Right: other,
		},
	}
}
func (f *Field[TField]) Lte(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "<=",
			Right: other,
		},
	}
}
func (f *Field[TField]) Add(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "+",
			Right: other,
		},
	}
}
func (f *Field[TField]) Sub(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "-",
			Right: other,
		},
	}
}

func (f *Field[TField]) Mul(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "*",
			Right: other,
		},
	}
}

func (f *Field[TField]) Div(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "/",
			Right: other,
		},
	}
}
func (f *Field[TField]) Len() *Field[int] {
	return &Field[int]{
		FuncField: &FuncField{
			FuncName: "LEN",
			Args:     []interface{}{f},
		},
	}
}

/*
Logic Operand
*/
func (f *FieldBool) And(other interface{}) *FieldBool {
	return &FieldBool{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "AND",
			Right: other,
		},
	}
}

func (f *FieldBool) Or(other interface{}) *FieldBool {
	return &FieldBool{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "OR",
			Right: other,
		},
	}
}

func (f *FieldBool) Not() *FieldBool {
	return &FieldBool{
		BinaryField: &BinaryField{
			Left:  nil,
			Op:    "NOT",
			Right: f,
		},
	}
}
func (f *Field[TField]) IsNull() *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "IS NULL",
			Right: nil,
		},
	}
}

func (f *Field[TField]) Sum() *FuncField {
	return &FuncField{
		FuncName: "SUM",
		Args:     []interface{}{f},
	}
}
func (f *Field[TField]) Min() *FuncField {
	return &FuncField{
		FuncName: "MIN",
		Args:     []interface{}{f},
	}
}
func (f *Field[TField]) Max() *FuncField {
	return &FuncField{
		FuncName: "MAX",
		Args:     []interface{}{f},
	}

}
func (f *Field[TField]) Avg() *FuncField {
	return &FuncField{
		FuncName: "AVG",
		Args:     []interface{}{f},
	}
}
func (f *Field[TField]) Count() *FuncField {
	return &FuncField{
		FuncName: "COUNT",
		Args:     []interface{}{f},
	}
}
