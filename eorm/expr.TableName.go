package eorm

import (
	"strconv"

	"eorm/sqlparser"
)

func (compiler *exprReceiver) TableName(context *exprCompileContext, expr *sqlparser.TableName) (string, error) {
	if expr.Name.String() == "$$$$$$$$$$$$$$" {
		return "", nil
	}
	tableName := expr.Name.String()

	if context.purpose == build_purpose_join {
		if aliasTableName, ok := context.stackAliasTables.Pop(); ok {
			if _, ok := context.alias[aliasTableName]; !ok {
				context.tables = append(context.tables, aliasTableName)
				context.alias[aliasTableName] = aliasTableName
			}
			return context.dialect.Quote(tableName) + " AS " + context.dialect.Quote(aliasTableName), nil
		} else {

			if _, ok := context.alias[tableName]; !ok {
				// if _, ok := (*context.schema)[tableName]; !ok { // not found in schema, try to pluralize it
				// 	tableName = utils.Plural(tableName)
				// }
				if _, ok := context.alias[tableName]; !ok {
					context.tables = append(context.tables, tableName)
					context.alias[tableName] = "T" + strconv.Itoa(len(context.tables))
				}
			}
			return context.dialect.Quote(tableName) + " AS " + context.dialect.Quote(context.alias[tableName]), nil
		}
	} else {
		if _, ok := (*context.schema)[tableName]; ok {
			return context.dialect.Quote(tableName), nil
		}
		tableName = utils.Plural(tableName)
		return context.dialect.Quote(tableName), nil
	}
	panic("not implemented")

}
