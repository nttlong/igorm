package expr

import (
	"fmt"
	"strings"
	"unvs-orm/internal"

	"github.com/xwb1989/sqlparser"
)

func (e *expression) compile(expr interface{}, isFunctionParamCompiler bool) ([]string, error) {

	switch expr := expr.(type) {
	case *sqlparser.SQLVal:
		return e.compileSQLVal(expr)

	case *sqlparser.ColName:
		tableNameFromSyntax := expr.Qualifier.ToViewName().Name.CompliantName()
		tableName := tableNameFromSyntax
		if dbTableName := internal.Utils.GetDbTableName(tableNameFromSyntax); dbTableName != "" {
			tableName = dbTableName
		}
		if tableName == "" {
			return nil, fmt.Errorf("table name not found for %s", expr)
		}
		fieldName := expr.Name.CompliantName()
		metaInfo := internal.Utils.GetMetaInfoByTableName(tableName)
		if metaInfo != nil {
			if _, ok := metaInfo[strings.ToLower(fieldName)]; !ok {
				return nil, fmt.Errorf("field %s not found in table %s", fieldName, tableName)
			} else {
				fieldName = internal.Utils.ToSnakeCase(fieldName)
			}
		}
		if !isFunctionParamCompiler {
			ret := e.Quote(tableName, fieldName) + " AS " + e.Quote(expr.Name.CompliantName())
			return []string{ret}, nil
		} else {
			return []string{e.Quote(tableName, fieldName)}, nil
		}

	case *sqlparser.AliasedExpr:
		ret, err := e.compile(expr.Expr, true)
		if err != nil {
			return nil, err
		}
		if isFunctionParamCompiler {
			return ret, nil
		} else {
			if expr.As.CompliantName() != "" {
				return []string{ret[0] + " AS " + e.Quote(expr.As.CompliantName())}, nil
			} else {
				return ret, nil
			}
		}
	case *sqlparser.FuncExpr:
		return e.funcExpr(expr)
	case sqlparser.SelectExprs:

		retStr := []string{}

		for _, expr := range expr {
			ret, err := e.compile(expr, isFunctionParamCompiler)
			if err != nil {
				return nil, err
			}
			retStr = append(retStr, ret...)
		}
		return retStr, nil

	case *sqlparser.SelectExprs:

		return e.compile(expr, isFunctionParamCompiler)
	case sqlparser.BinaryExpr:
		return e.compileBinaryExpr(&expr, false)
	case *sqlparser.BinaryExpr:
		return e.compileBinaryExpr(expr, false)
	default:
		return nil, fmt.Errorf("not support %s", expr)
	}

}
