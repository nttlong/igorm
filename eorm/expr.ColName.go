package eorm

import (
	"strconv"

	"github.com/xwb1989/sqlparser"
)

// ComparisonExpr
func (e *exprReceiver) ColName(context *exprCompileContext, expr sqlparser.ColName) (string, error) {

	tableName := expr.Qualifier.Name.String()
	fieldName := expr.Name.String()
	if _, ok := context.schema[tableName]; !ok {

		tableName = utils.Plural(tableName)
		fieldName = utils.ToSnakeCase(expr.Name.String())
	}

	if _, ok := context.alias[tableName]; !ok {
		context.tables = append(context.tables, tableName)
		context.alias[tableName] = "T" + strconv.Itoa(len(context.tables))
	}
	return context.dialect.Quote(context.alias[tableName], fieldName), nil

}
