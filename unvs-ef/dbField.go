package unvsef

type imp struct {
	Key string
}
type DbField struct {
	TableName string
	FieldName string
	ColName   string
}

func (dbd *DbField) ToSqlExpr(d Dialect) (string, []interface{}) {

	return d.QuoteIdent(dbd.TableName, dbd.ColName), nil
}
func (dbd DbField) ToSqlExpr2(d Dialect) (string, []interface{}) {

	return d.QuoteIdent(dbd.TableName, dbd.ColName), nil
}

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

type SortField struct {
	Field interface{}
	Sort  string
}
type Field[TField any] struct {
	*DbField
	*AliasField
	*BinaryField
	*FuncField
	*SortField
	Op string

	val *TField
}

func (f *DbField) clone() *DbField {
	return &DbField{
		TableName: f.TableName,
		ColName:   f.ColName,
		FieldName: f.FieldName,
	}
}
func (f *Field[TField]) Set(val *TField) {
	f.val = val
}
func (f *Field[TField]) Get() *TField {
	return f.val
}

func (f *Field[TField]) Eq(other interface{}) *BinaryField {
	return &BinaryField{

		Left:  f,
		Op:    "=",
		Right: other,
	}
}
func (f *Field[TField]) Ne(other interface{}) *Field[TField] {
	return &Field[TField]{
		DbField: f.DbField,
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "!=",
			Right: other,
		},
	}
}
func (f *Field[TField]) Gt(other interface{}) *Field[TField] {
	return &Field[TField]{
		DbField: f.DbField,
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
		DbField: f.DbField,
		BinaryField: &BinaryField{
			Left:  f,
			Op:    ">=",
			Right: other,
		},
	}
}
func (f *Field[TField]) Lte(other interface{}) *Field[TField] {
	return &Field[TField]{
		DbField: f.DbField,
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "<=",
			Right: other,
		},
	}
}
func (f *Field[TField]) Add(other interface{}) *Field[TField] {
	return &Field[TField]{
		DbField: f.DbField,
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "+",
			Right: other,
		},
	}
}
func (f *Field[TField]) Sub(other interface{}) *Field[TField] {
	return &Field[TField]{
		DbField: f.DbField,
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "-",
			Right: other,
		},
	}
}

func (f *Field[TField]) Mul(other interface{}) *Field[TField] {
	return &Field[TField]{
		DbField: f.DbField,
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "*",
			Right: other,
		},
	}
}

func (f *Field[TField]) Div(other interface{}) *Field[TField] {
	return &Field[TField]{
		DbField: f.DbField,
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "/",
			Right: other,
		},
	}
}
func (f *Field[TField]) Len() *Field[int] {
	return &Field[int]{
		DbField: f.DbField,
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
func Lit[TField any](val TField) *TField {
	return &val
}
