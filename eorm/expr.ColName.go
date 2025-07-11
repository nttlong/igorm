package eorm

import (
	"strconv"

	"eorm/sqlparser"
)

// ComparisonExpr
func (e *exprReceiver) ColName(context *exprCompileContext, expr sqlparser.ColName) (string, error) {

	tableName := expr.Qualifier.Name.String()
	fieldName := expr.Name.String()

	if _, ok := context.alias[tableName]; !ok { // if not found in alias, then check if it is a schema table
		if _, ok := (*context.schema)[tableName]; !ok { // if not found in database schema, then assume it is a plural table name

			tableName = utils.Plural(tableName)
			fieldName = utils.ToSnakeCase(expr.Name.String())
		}
		if _, ok := context.alias[tableName]; !ok {

			context.tables = append(context.tables, tableName)
			context.alias[tableName] = "T" + strconv.Itoa(len(context.tables))
		}
	}
	if context.purpose == build_purpose_for_function {
		return context.dialect.Quote(context.alias[tableName], fieldName), nil
	}
	if context.purpose == build_purpose_select {
		if aliasField, ok := context.stackAliasFields.Pop(); ok {
			ret := context.dialect.Quote(context.alias[tableName], fieldName) + " AS " + context.dialect.Quote(aliasField)

			return ret, nil
		}
		return context.dialect.Quote(context.alias[tableName], fieldName) + " AS " + context.dialect.Quote(expr.Name.String()), nil
	}
	return context.dialect.Quote(context.alias[tableName], fieldName), nil

}
