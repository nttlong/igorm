package expr

import (
	"fmt"
	"strings"
	"unvs-orm/internal"

	"github.com/xwb1989/sqlparser"
)

type expressionCompileResult struct {
	Syntax string
	Tables []string
}
type ResolveContext struct {
	Tables []string

	Map map[string]string
}

func (c *ResolveContext) addTables(tableNames ...string) {
	for _, tableName := range tableNames {
		if _, ok := c.Map[tableName]; !ok { // Đọc từ map thông thường
			c.Tables = append(c.Tables, tableName)
			c.Map[tableName] = "T" + fmt.Sprintf("%d", len(c.Tables))
			// Ghi vào map thông thường
		}
	}
}
func (e *expression) compile(expr interface{}, context *ResolveContext, isFunctionParamCompiler bool) (*expressionCompileResult, error) {
	if context == nil {
		context = &ResolveContext{Tables: []string{}, Map: map[string]string{}}
	}
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
			context.addTables(tableName)
			tableAlias := context.Map[tableName]
			ret := e.Quote(tableAlias, fieldName) + " AS " + e.Quote(expr.Name.CompliantName())
			return &expressionCompileResult{Syntax: ret, Tables: context.Tables}, nil
		} else {
			context.addTables(tableName)
			tableAlias := context.Map[tableName]
			return &expressionCompileResult{Syntax: e.Quote(tableAlias, fieldName), Tables: context.Tables}, nil
		}

	case *sqlparser.AliasedExpr:
		ret, err := e.compile(expr.Expr, context, true)
		if err != nil {
			return nil, err
		}
		if isFunctionParamCompiler {
			return ret, nil
		} else {
			if expr.As.CompliantName() != "" {
				synTax := ret.Syntax + " AS " + e.Quote(expr.As.CompliantName())

				return &expressionCompileResult{Syntax: synTax}, nil
			} else {
				return ret, nil
			}
		}
	case *sqlparser.FuncExpr:
		return e.funcExpr(expr, context)
	case sqlparser.SelectExprs:

		retStr := []string{}

		for _, expr := range expr {
			ret, err := e.compile(expr, context, isFunctionParamCompiler)
			if err != nil {
				return nil, err
			}

			retStr = append(retStr, ret.Syntax)
		}
		return &expressionCompileResult{Syntax: strings.Join(retStr, ", ")}, nil

	case *sqlparser.SelectExprs:
		return e.compile(&expr, context, isFunctionParamCompiler)
	case sqlparser.BinaryExpr:
		return e.compileBinaryExpr(&expr, context, false)
	case *sqlparser.BinaryExpr:
		return e.compileBinaryExpr(expr, context, false)
	case *sqlparser.AndExpr:
		return e.AndExpr(expr, context, isFunctionParamCompiler)
	case sqlparser.AndExpr:
		return e.AndExpr(&expr, context, isFunctionParamCompiler)
	case *sqlparser.OrExpr:
		return e.OrExpr(expr, context, isFunctionParamCompiler)
	case sqlparser.OrExpr:
		return e.OrExpr(&expr, context, isFunctionParamCompiler)
	case *sqlparser.ComparisonExpr:
		return e.ComparisonExpr(expr, context, isFunctionParamCompiler)
	case sqlparser.ComparisonExpr:
		return e.ComparisonExpr(&expr, context, isFunctionParamCompiler)
	default:
		return nil, fmt.Errorf("not support %s in file orm/expression.getField.go", expr)
	}

}
