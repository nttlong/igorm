package migrate

import (
	"fmt"
	"reflect"
	"strings"
)

func (m *migratorMySql) GetSqlCreateTable(typ reflect.Type) (string, error) {
	mapType := m.GetColumnDataTypeMapping()
	defaultValueByFromDbTag := m.GetGetDefaultValueByFromDbTag()
	schemaLoader := m.GetLoader()
	if schemaLoader == nil {
		return "", fmt.Errorf("schema loader is nil, please set it by call SetLoader() function in %s", reflect.TypeOf(m).Elem())
	}

	schema, err := schemaLoader.LoadFullSchema(m.db)
	if err != nil {
		return "", err
	}

	entityItem := ModelRegistry.GetModelByType(typ)
	if entityItem == nil {
		return "", fmt.Errorf("model %s not found, please register model first by call ModelRegistry.Add(%s)", typ.String(), typ.String())
	}

	tableName := entityItem.tableName
	if _, ok := schema.Tables[tableName]; ok {
		return "", nil // table already exists
	}

	strCols := []string{}
	newTableMap := map[string]bool{}
	for _, col := range entityItem.entity.cols {
		newTableMap[col.Name] = true
		fieldType := col.Field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		sqlType, ok := mapType[fieldType]
		if !ok {
			panic(fmt.Sprintf("not support field type %s, review GetColumnDataTypeMapping() function in %s", fieldType.String(), reflect.TypeOf(m).Elem()))
		}

		if col.Length != nil {
			sqlType = fmt.Sprintf("%s(%d)", sqlType, *col.Length)
		}

		colDef := m.Quote(col.Name) + " " + sqlType

		if col.IsAuto {
			colDef += " AUTO_INCREMENT"
		}

		if col.Nullable {
			colDef += " NULL"
		} else {
			colDef += " NOT NULL"
		}

		if col.Default != "" {
			if val, ok := defaultValueByFromDbTag[col.Default]; ok {
				colDef += fmt.Sprintf(" DEFAULT %s", val)
			} else {
				panic(fmt.Errorf("not support default value from %s, review GetGetDefaultValueByFromDbTag() function in %s", col.Default, reflect.TypeOf(m).Elem()))
			}
		}

		strCols = append(strCols, colDef)
	}

	for _, cols := range entityItem.entity.primaryConstraints {
		var pkCols []string
		var pkColNames []string
		for _, col := range cols {
			if col.PKName != "" {
				pkCols = append(pkCols, m.Quote(col.Name))
				pkColNames = append(pkColNames, col.Name)
			}
		}

		if len(pkCols) > 0 {
			pkConstraintName := fmt.Sprintf("PK_%s__%s", tableName, strings.Join(pkColNames, "_"))
			// MySQL thường không cần đặt tên, nhưng bạn vẫn có thể nếu muốn
			constraint := fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", m.Quote(pkConstraintName), strings.Join(pkCols, ", "))
			strCols = append(strCols, constraint)
		}
	}

	sql := fmt.Sprintf("CREATE TABLE %s (\n  %s\n)", m.Quote(tableName), strings.Join(strCols, ",\n  "))
	schema.Tables[tableName] = newTableMap

	return sql, nil
}
