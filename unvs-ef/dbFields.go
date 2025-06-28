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

//	func Lit[T any](val T) *Field[T] {
//		return &Field[T]{
//			ValueField: &ValueField[T]{
//				Value: val,
//			},
//		}
//	}
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

func (f *Field[TField]) Mod(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "%",
			Right: other,
		},
	}
}
func (f *Field[TField]) Like(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "LIKE",
			Right: other,
		},
	}
}
func (f *Field[TField]) NotLike(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "LIKE",
			Right: other,
		},
	}
}
func (f *Field[TField]) In(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "IN",
			Right: other,
		},
	}
}
func (f *Field[TField]) NotIn(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "NOT IN",
			Right: other,
		},
	}
}
func (f *Field[TField]) Between(left interface{}, right interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "BETWEEN",
			Right: []interface{}{left, right},
		},
	}
}
func (f *FieldDateTime) Between(left interface{}, right interface{}) *FieldBool {
	return &FieldBool{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "BETWEEN",
			Right: []interface{}{left, right},
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
func (f *Field[TField]) Case(cases []interface{}, elseVal interface{}) *Field[TField] {
	return &Field[TField]{
		FuncField: &FuncField{
			FuncName: "CASE",
			Args:     cases,
		},
	}
}

func (f *FieldString) Len() *Field[int] {
	return &Field[int]{
		FuncField: &FuncField{
			FuncName: "LEN",
			Args:     []interface{}{f},
		},
	}
}

func (f *FieldDateTime) Year() *Field[int] {
	return &Field[int]{
		FuncField: &FuncField{
			FuncName: "YEAR",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldDateTime) Month() *Field[int] {
	return &Field[int]{
		FuncField: &FuncField{
			FuncName: "MONTH",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldDateTime) Day() *Field[int] {
	return &Field[int]{
		FuncField: &FuncField{
			FuncName: "DAY",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldDateTime) Hour() *Field[int] {
	return &Field[int]{
		FuncField: &FuncField{
			FuncName: "HOUR",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldDateTime) Minute() *Field[int] {
	return &Field[int]{
		FuncField: &FuncField{
			FuncName: "MINUTE",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldDateTime) Second() *Field[int] {
	return &Field[int]{
		FuncField: &FuncField{
			FuncName: "SECOND",
			Args:     []interface{}{f},
		},
	}
}
func (f *BinaryField) ToSqlExpr(d Dialect) (string, []interface{}) {
	return compiler.Compile(f, d)
}
func (f *Field[TField]) ToSqlExpr(d Dialect) (string, []interface{}) {
	return compiler.Compile(f, d)
}
func (where *FieldBool) ToSqlExpr(d Dialect) (string, []interface{}) {
	return compiler.Compile(where, d)
}
func (f Field[TField]) ToSqlExpr2(d Dialect) (string, []interface{}) {
	return (&f).ToSqlExpr(d)
}
