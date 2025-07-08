package orm

func (b *Base) Join(expr string, args ...interface{}) *BoolField {
	return &BoolField{
		underField: &JoinExpr{
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
		underField: &JoinExpr{
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
		underField: &JoinExpr{
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
		underField: &JoinExpr{
			joinType: "FULL",
			joinExprText: &joinExprText{
				Expr: expr,
				Args: args,
			},
		},
	}
}
