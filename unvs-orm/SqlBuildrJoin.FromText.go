package orm

import (
	EXPR "unvs-orm/expr"
)

func (c *JoinCompilerUtils) fromExprString(expr *JoinExpr) (*resolverResult, error) {
	if v, ok := c.cacheFromExprString.Load(expr.Expr); ok {
		return v.(*resolverResult), nil
	}
	result, err := c.fromExprStringNoCache(expr)
	if err != nil {
		return nil, err
	}
	c.cacheFromExprString.Store(expr.Expr, result)
	return result, nil
}
func (c *JoinCompilerUtils) fromExprStringNoCache(expr *JoinExpr) (*resolverResult, error) {
	/*
		e := EXPR.ExpressionTest{
			DbDriver: EXPR.DB_TYPE_MSSQL,
		}
		cmd := "ORDER.OrderID,Order.Note"

		compiled, err := e.CompileSelect(cmd)

		assert.NoError(t, err)
		compiledExpected := "[orders].[order_id] AS [OrderID], [orders].[note] AS [Note]"
		assert.Equal(t, compiledExpected, compiled)
	*/
	e := EXPR.ExpressionTest{
		DbDriver: EXPR.DB_TYPE_UNKNOWN.FromString(c.dialect.driverName()),
	}
	compiled, err := e.CompileSelect(expr.Expr)
	if err != nil {
		return nil, err
	}

	return &resolverResult{
		Syntax:       compiled.Syntax,
		Args:         expr.Args,
		buildContext: &compiled.Context.Map,
		Tables:       compiled.Context.Tables,
	}, nil
}
