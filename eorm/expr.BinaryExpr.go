package eorm

import (
	"github.com/xwb1989/sqlparser"
)

func (compiler *exprReceiver) BinaryExpr(context *exprCompileContext, expr *sqlparser.BinaryExpr) (string, error) {
	left, err := compiler.compile(context, expr.Left)
	if err != nil {
		return "", err
	}
	right, err := compiler.compile(context, expr.Right)
	if err != nil {
		return "", err
	}

	return left + " " + expr.Operator + " " + right, nil

}
