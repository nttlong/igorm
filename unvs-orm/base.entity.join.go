package orm

func (b *Base) Join(expr string, args ...interface{}) *BoolField {
	return &BoolField{
		UnderField: &JoinExpr{
			joinType: "INNER",
			joinExprText: &joinExprText{
				Expr: expr,
				Args: args,
			},
		},
	}

}
func (b *Base) LeftJoin(expr string, args ...interface{}) *BoolField {
	return &BoolField{
		UnderField: &JoinExpr{
			joinType: "LEFT",
			joinExprText: &joinExprText{
				Expr: expr,
				Args: args,
			},
		},
	}
}
func (b *Base) RightJoin(expr string, args ...interface{}) *BoolField {
	return &BoolField{
		UnderField: &JoinExpr{
			joinType: "RIGHT",
			joinExprText: &joinExprText{
				Expr: expr,
				Args: args,
			},
		},
	}
}
func (b *Base) FullJoin(expr string, args ...interface{}) *BoolField {
	return &BoolField{
		UnderField: &JoinExpr{
			joinType: "FULL",
			joinExprText: &joinExprText{
				Expr: expr,
				Args: args,
			},
		},
	}
}
