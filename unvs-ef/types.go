package unvsef

type DbField[Field any] struct {
	TableName string
	ColName   string
	Value     *Field
}

func (f DbField[Field]) ToSQL(d Dialect) (string, []interface{}) {
	return d.QuoteIdent(f.TableName, f.ColName), nil
}
func (f *DbField[T]) Set(v T) {
	f.Value = &v
}

func (f DbField[T]) Get() T {
	if f.Value == nil {
		var zero T
		return zero
	}
	return *f.Value
}
func (f DbField[Field]) Eq(value interface{}) *BinaryExpr {
	return AddExpr(f, "=", toExpr(value))
}
func (f DbField[Field]) Ne(value interface{}) *BinaryExpr {
	return AddExpr(f, "<>", toExpr(value))
}

func (f DbField[Field]) Gt(value interface{}) *BinaryExpr {
	return AddExpr(f, ">", toExpr(value))
}

func (f DbField[Field]) Lt(value interface{}) *BinaryExpr {
	return AddExpr(f, "<", toExpr(value))
}

func (f DbField[Field]) Gte(value interface{}) *BinaryExpr {
	return AddExpr(f, ">=", toExpr(value))
}

func (f DbField[Field]) Lte(value interface{}) *BinaryExpr {
	return AddExpr(f, "<=", toExpr(value))
}

func (f DbField[Field]) Like(value interface{}) *BinaryExpr {
	return AddExpr(f, "LIKE", toExpr(value))
}

func (f DbField[Field]) In(values ...interface{}) *BinaryExpr {
	return AddExpr(f, "IN", Literal[[]interface{}]{Value: values})
}

func (f DbField[Field]) IsNull() *BinaryExpr {
	return AddExpr(f, "IS", Literal[string]{Value: "NULL"})
}

func (f DbField[Field]) IsNotNull() *BinaryExpr {
	return AddExpr(f, "IS NOT", Literal[string]{Value: "NULL"})
}
