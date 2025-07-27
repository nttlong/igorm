package vdb

import (
	"fmt"
	"reflect"
	"strings"
	"vdb/tenantDB"
)

func buildBasicSqlFirstItemNoFilterNoCache(typ reflect.Type, db *tenantDB.TenantDB) (string, string, [][]int, error) {
	dialect := dialectFactory.create(db.GetDriverName())

	repoType := inserterObj.getEntityInfo(typ)
	tableName := repoType.tableName
	compiler, err := CompileJoin(tableName, db)
	if err != nil {
		return "", "", nil, err
	}
	tableName = compiler.content
	columns := repoType.entity.GetColumns()

	fieldsSelect := make([]string, len(columns))
	filterFields := []string{}
	keyFieldIndex := [][]int{}
	for i, col := range columns {
		if col.PKName != "" {
			filterFields = append(filterFields, repoType.tableName+"."+col.Name+" =?")
			keyFieldIndex = append(keyFieldIndex, col.IndexOfField)
		}
		fieldsSelect[i] = repoType.tableName + "." + col.Field.Name + " AS " + col.Field.Name
	}
	filter := strings.Join(filterFields, " AND ")
	compiler.context.purpose = build_purpose_select
	err = compiler.buildSelectField(strings.Join(fieldsSelect, ", "))
	if err != nil {
		return "", "", nil, err
	}
	strField := compiler.content

	sql := fmt.Sprintf("SELECT %s FROM %s", strField, tableName)
	if filter != "" {
		compiler.context.purpose = build_purpose_where
		err = compiler.buildWhere(filter)
		if err != nil {
			return "", "", nil, err
		}

	}
	sql = dialect.MakeSelectTop(sql, 1)
	return sql, compiler.content, keyFieldIndex, nil
}
