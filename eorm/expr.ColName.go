package eorm

import (
	"eorm/sqlparser"
	"strconv"
)

// ComparisonExpr
func (e *exprReceiver) ColName(context *exprCompileContext, expr sqlparser.ColName) (string, error) {
	if context.aliasToDbTable == nil {
		context.aliasToDbTable = map[string]string{}
	}

	tableName := expr.Qualifier.Name.String()
	fieldName := expr.Name.String()
	aliasFieldName := expr.Name.String()
	if context.schema == nil {
		context.schema = &map[string]bool{}
	}

	if _, ok := (*context.schema)[tableName]; !ok {
		/*
			if not found in database calculate alias table name , field name and alias field name
		*/

		if _, ok := context.aliasToDbTable[tableName]; !ok {
			fieldName = utils.ToSnakeCase(fieldName)
		}
		if aliasTable, ok := context.alias[tableName]; ok {
			tableName = aliasTable
		} else {
			if context.purpose == build_purpose_join {
				/*
					if purpose is join, the compiling process need
					extract tables if they were found when compiling the query
				*/
				context.tables = append(context.tables, tableName)
				context.alias[tableName] = "T" + strconv.Itoa(len(context.tables))
				tableName = context.alias[tableName]
			} else {
				tableName = utils.Plural(tableName)
			}
		}
		if aliasFieldFromStack, ok := context.stackAliasFields.Pop(); ok {
			aliasFieldName = aliasFieldFromStack
		}
	} else {
		/*
			if found in database calculate alias field name
			tableName from database schema
			column name no change because it is already in SQL statement where DEV declared in their code
		*/
		/*
			But important that we need to check if alias field name is already in stack, if yes then use it, otherwise use field name
		*/
		if aliasFieldFromStack, ok := context.stackAliasFields.Pop(); ok {
			aliasFieldName = aliasFieldFromStack
		} else {
			aliasFieldName = fieldName
		}

	}
	if context.purpose == build_purpose_select {
		/*
			if purpose is select, then return tablename.fieldname as aliasfieldname
			Heed: quote all the things
		*/
		return context.dialect.Quote(tableName, fieldName) + " AS " + context.dialect.Quote(aliasFieldName), nil

	} else {
		return context.dialect.Quote(tableName, fieldName), nil
	}

	// if not found in database schema, then assume it is a plural table name
	// if context.schema == nil {
	// 	context.schema = &map[string]bool{}
	// }
	// tableName := expr.Qualifier.Name.String()
	// fieldName := expr.Name.String()
	// var fullName string
	// if aliasField, ok := context.stackAliasFields.Pop(); ok {
	// 	if context.purpose == build_purpose_select {
	// 		fullName = context.dialect.Quote(tableName, fieldName) + " AS " + context.dialect.Quote(aliasField)
	// 	}
	// } else {
	// 	if context.purpose == build_purpose_select {
	// 		fullName = context.dialect.Quote(context.alias[tableName], utils.ToSnakeCase(fieldName)) + " AS " + context.dialect.Quote(expr.Name.String())
	// 	} else {
	// 		fullName = context.dialect.Quote(context.alias[tableName], utils.ToSnakeCase(fieldName))
	// 	}

	// }
	// return fullName, nil
	// if context.purpose == build_purpose_select {
	// 	return fullName + " AS " + context.dialect.Quote(expr.Name.String()), nil

	// } else if context.purpose == build_purpose_select {
	// 	return fullName, nil
	// } else {
	// 	return fullName, nil
	// }

	// if _, ok := context.alias[tableName]; !ok { // if not found in alias, then check if it is a schema table
	// 	if _, ok := (*context.schema)[tableName]; !ok { // if not found in database schema, then assume it is a plural table name

	// 		tableName = utils.Plural(tableName)
	// 		//fieldName = utils.ToSnakeCase(expr.Name.String())
	// 	}
	// 	if _, ok := context.alias[tableName]; !ok {

	// 		context.tables = append(context.tables, tableName)
	// 		context.alias[tableName] = "T" + strconv.Itoa(len(context.tables))
	// 	}
	// }
	// if context.purpose == build_purpose_for_function || context.purpose == build_purpose_join {
	// 	compileTableName := context.pluralTableName(tableName)
	// 	compileFieldName := utils.ToSnakeCase(expr.Name.String())

	// 	return context.dialect.Quote(compileTableName, compileFieldName), nil
	// }
	// if context.purpose == build_purpose_select {
	// 	if aliasField, ok := context.stackAliasFields.Pop(); ok {

	// 		compileFieldName := utils.ToSnakeCase(expr.Name.String())
	// 		ret := context.dialect.Quote(tableName, compileFieldName) + " AS " + context.dialect.Quote(aliasField)

	// 		return ret, nil
	// 	}
	// 	compileTableName := tableName
	// 	if aliasTable, ok := context.alias[tableName]; ok {
	// 		compileTableName = aliasTable
	// 	}
	// 	compileTableName = context.pluralTableName(tableName)
	// 	compileFieldName := utils.ToSnakeCase(expr.Name.String())
	// 	return context.dialect.Quote(compileTableName, compileFieldName) + " AS " + context.dialect.Quote(expr.Name.String()), nil
	// }
	// compileTableName := context.pluralTableName(tableName)
	// compileFieldName := utils.ToSnakeCase(expr.Name.String())
	// return context.dialect.Quote(compileTableName, compileFieldName), nil

}
