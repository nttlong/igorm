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
		source:  m.TableName,
		selects: fields,
		noAlias: true,
	}
}

func (m *Model[T]) Filter(expr *BoolField) *SqlSelectBuilder {
	return &SqlSelectBuilder{
		source:    m.TableName,
		condition: expr,
	}
}
