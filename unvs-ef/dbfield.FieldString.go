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
		DbField: &DbField{
			TableName: f.TableName,
			ColName:   f.ColName,
		},
		FuncField: &FuncField{
			FuncName: "LEN",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldString) Like(other interface{}) *FieldBool {
	return &FieldBool{
		DbField: &DbField{
			TableName: f.TableName,
			ColName:   f.ColName,
		},
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "LIKE",
			Right: other,
		},
	}
}
func (f *Field[TField]) NotLike(other interface{}) *FieldBool {
	return &FieldBool{
		DbField: &DbField{
			TableName: f.TableName,
			ColName:   f.ColName,
		},
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "NOT LIKE",
			Right: other,
		},
	}
}
func (f *FieldString) In(other interface{}) *FieldBool {
	return &FieldBool{
		DbField: &DbField{
			TableName: f.TableName,
			ColName:   f.ColName,
		},
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "IN",
			Right: other,
		},
	}
}
func (f *FieldString) NotIn(other interface{}) *FieldBool {
	return &FieldBool{
		DbField: &DbField{
			TableName: f.TableName,
			ColName:   f.ColName,
		},
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "NOT IN",
			Right: other,
		},
	}
}
func (f *FieldString) Eq(other interface{}) *FieldBool {
	return &FieldBool{
		DbField: &DbField{
			TableName: f.TableName,
			ColName:   f.ColName,
		},
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "=",
			Right: other,
		},
	}
}
func (f *FieldString) NotEq(other interface{}) *FieldBool {
	return &FieldBool{
		DbField: &DbField{
			TableName: f.TableName,
			ColName:   f.ColName,
		},
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "!=",
			Right: other,
		},
	}
}
func (f *FieldString) IsNull() *FieldBool {
	return &FieldBool{
		DbField: &DbField{
			TableName: f.TableName,
			ColName:   f.ColName,
		},
		BinaryField: &BinaryField{
			Left: f,
			Op:   "IS NULL",
		},
	}
}
func (f *FieldString) IsNotNull() *FieldBool {
	return &FieldBool{
		DbField: &DbField{
			TableName: f.TableName,
			ColName:   f.ColName,
		},
		BinaryField: &BinaryField{
			Left: f,
			Op:   "IS NOT NULL",
		},
	}
}
func (f *FieldString) Asc() *SortField {
	return &SortField{
		Field: f,
		Sort:  "ASC",
	}
}

func (f *FieldString) Desc() *SortField {
	return &SortField{
		Field: f,
		Sort:  "DESC",
	}
}
func (f *FieldString) Get() string {
	return *f.val
}
func (f *FieldString) Set(val *string) {
	f.val = val
}
func (f *FieldString) As(alias string) *AliasField {
	return &AliasField{
		Field: f,
		Alias: alias,
	}
}
