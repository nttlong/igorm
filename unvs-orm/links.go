package orm

import (
	"fmt"
	expression "unvs-orm/expr"
)

func init() {
	expression.OnGetQuoteFunc = func(dbDriver expression.DB_TYPE, str ...string) string {

		if dbDriver == expression.DB_TYPE_MSSQL {
			panic("not implemented mysql dialect")
		}
		if dbDriver == expression.DB_TYPE_POSTGRES {
			panic("not implemented postgres dialect")
		}
		if dbDriver == expression.DB_TYPE_MSSQL {
			return NewMssqlDialect().getCompiler().Quote(str...)
		}

		panic(fmt.Sprintf("not support dialect for %s, file orm/links.go, line %d", dbDriver, 21))

	}
	expression.OnCompileFunc = func(dbDriver expression.DB_TYPE, aliasSource *map[string]string, caller interface{}) (*expression.ResolverResult, error) {
		if dbDriver == expression.DB_TYPE_MSSQL {
			panic("not implemented mysql dialect")
		}
		if dbDriver == expression.DB_TYPE_POSTGRES {
			panic("not implemented postgres dialect")
		}
		if dbDriver == expression.DB_TYPE_MSSQL {
			//if exprCall,ok:=caller.(*expression.ExpressionTest);ok{
			result, err := NewMssqlDialect().getCompiler().Resolve(aliasSource, caller)
			if err != nil {
				return nil, err
			}
			return &expression.ResolverResult{
				Syntax:      result.Syntax,
				Args:        result.Args,
				AliasSource: result.AliasSource,
			}, nil

		}

		panic(fmt.Sprintf("not support dialect for %s, file orm/links.go, line %d", dbDriver, 45))
	}
}
