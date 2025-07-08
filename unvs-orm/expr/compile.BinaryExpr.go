package expr

import (
	"fmt"

	"github.com/xwb1989/sqlparser"
)

func (e *expression) compileBinaryExpr(expr *sqlparser.BinaryExpr, tables *[]string, context *map[string]string, isFunctionParamCompiler bool, requireAlias bool) (*expressionCompileResult, error) {
	left, err := e.compile(expr.Left, tables, context, true, requireAlias)

	if err != nil {
		return nil, err
	}

	right, err := e.compile(expr.Right, tables, context, true, requireAlias)
	if err != nil {
		return nil, err
	}

	ret := fmt.Sprintf("%s %s %s", left.Syntax, expr.Operator, right.Syntax)
	return &expressionCompileResult{Syntax: ret}, nil
}
