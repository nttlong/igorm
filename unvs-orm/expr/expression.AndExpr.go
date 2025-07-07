package expr

import (
	"strings"

	"github.com/xwb1989/sqlparser"
)

func (e *expression) AndExpr(expr *sqlparser.AndExpr, context *ResolveContext, isFunctionParamCompiler bool) (*expressionCompileResult, error) {
	left, err := e.compile(expr.Left, context, isFunctionParamCompiler)
	if err != nil {
		return nil, err
	}

	right, err := e.compile(expr.Right, context, isFunctionParamCompiler)
	if err != nil {
		return nil, err
	}

	ret := []string{left.Syntax + " AND " + right.Syntax}
	return &expressionCompileResult{Syntax: strings.Join(ret, " AND ")}, nil
}
func (e *expression) OrExpr(expr *sqlparser.OrExpr, context *ResolveContext, isFunctionParamCompiler bool) (*expressionCompileResult, error) {
	left, err := e.compile(expr.Left, context, isFunctionParamCompiler)
	if err != nil {
		return nil, err
	}

	right, err := e.compile(expr.Right, context, isFunctionParamCompiler)
	if err != nil {
		return nil, err
	}

	ret := []string{left.Syntax + " OR " + right.Syntax}
	return &expressionCompileResult{Syntax: strings.Join(ret, " OR ")}, nil
}
