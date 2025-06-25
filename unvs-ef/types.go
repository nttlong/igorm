package unvsef

type DbField[Field any] struct {
	TableName string
	ColName   string
}

func (f DbField[Field]) ToSQL(d Dialect) (string, []interface{}) {
	return d.QuoteIdent(f.TableName, f.ColName), nil
}

func (f DbField[Field]) Eq(value interface{}) *BinaryExpr {
	return AddExpr(f, "=", toExpr(value))
}
