package orm

func (f *Base) Expr(expr string, args ...interface{}) *exprField {
	return &exprField{
		UnderField: &ExprBase{
			Stmt: expr,
			Args: args,
		},
	}
}
func (m *Model[T]) Select(fields ...interface{}) *SqlSelectBuilder {
	return &SqlSelectBuilder{
		source:  m,
		selects: fields,
	}
}
func (m *Model[T]) Where(expr *BoolField) *SqlSelectBuilder {
	return &SqlSelectBuilder{
		source:    m,
		condition: expr,
	}
}
