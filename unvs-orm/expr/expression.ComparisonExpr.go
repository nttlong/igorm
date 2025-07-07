package expr

import "github.com/xwb1989/sqlparser"

func (e *expression) ComparisonExpr(expr *sqlparser.ComparisonExpr, context *ResolveContext, isFunctionParamCompiler bool) (*expressionCompileResult, error) {
	// TODO: implement this function
	left, err := e.compile(expr.Left, context, isFunctionParamCompiler)
	if err != nil {
		return nil, err
	}

	right, err := e.compile(expr.Right, context, isFunctionParamCompiler)
	if err != nil {
		return nil, err
	}

	syntax := left.Syntax + " " + string(expr.Operator) + " " + right.Syntax
	return &expressionCompileResult{Syntax: syntax}, nil

}
