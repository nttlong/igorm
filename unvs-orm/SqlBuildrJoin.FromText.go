package orm

import (
	EXPR "unvs-orm/expr"
)

func (c *JoinCompilerUtils) fromExprString(sourceCache string, expr *JoinExpr, tables *[]string, context *map[string]string) (*resolverResult, error) {
	// if v, ok := c.cacheFromExprString.Load(expr.Expr); ok {
	// 	return v.(*resolverResult), nil
	// }
	result, err := c.fromExprStringNoCache(sourceCache, expr, tables, context)
	if err != nil {
		return nil, err
	}
	// c.cacheFromExprString.Store(expr.Expr, result)
	return result, nil
}
func (c *JoinCompilerUtils) fromExprStringNoCache(sourceCache string, expr *JoinExpr, tables *[]string, context *map[string]string) (*resolverResult, error) {

	e := EXPR.ExpressionTest{
		DbDriver: EXPR.DB_TYPE_UNKNOWN.FromString(c.dialect.driverName()),
	}

	compiled, err := e.Compile(sourceCache, tables, context, expr.Expr, true, true)

	if err != nil {
		return nil, err
	}

	return &resolverResult{
		Syntax: compiled.Syntax,
		Args:   expr.Args,
	}, nil
}
