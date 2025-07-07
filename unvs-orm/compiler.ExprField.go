package orm

import (
	EXPR "unvs-orm/expr"
)

func (c *CompilerUtils) resolveExprField(context *map[string]string, f *exprField) (*resolverResult, error) {
	if bf, ok := f.UnderField.(*ExprBase); ok {
		r, err := EXPR.NewExpressionCompiler(c.dialect.driverName()).CompileSelect(bf.Stmt)
		if err != nil {
			return nil, err
		}
		return &resolverResult{
			Syntax:       r.Syntax,
			Args:         r.Args,
			buildContext: context,
			Tables:       r.Context.Tables,
		}, nil
	} else {
		return c.Resolve(context, f.UnderField)
	}
}
