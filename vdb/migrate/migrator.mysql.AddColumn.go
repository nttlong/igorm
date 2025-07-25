package migrate

import (
	"fmt"
	"reflect"
	"strings"
)

func (m *migratorMySql) GetSqlAddColumn(typ reflect.Type) (string, error) {
	mapType := m.GetColumnDataTypeMapping()
	defaultValueByFromDbTag := m.GetGetDefaultValueByFromDbTag()

	schema, err := m.loader.LoadFullSchema(m.db)
	if err != nil {
		return "", err
	}

	entityItem := ModelRegistry.GetModelByType(typ)
	if entityItem == nil {
		return "", fmt.Errorf("model %s not found, please register model first by call ModelRegistry.Add(%s)", typ.String(), typ.String())
	}

	scripts := []string{}
	tableName := entityItem.tableName

	for _, col := range entityItem.entity.cols {
		if _, ok := schema.Tables[tableName][col.Name]; !ok {
			fieldType := col.Field.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}

			sqlType, ok := mapType[fieldType]
			if !ok {
				panic(fmt.Sprintf("unsupported field type %s, check GetColumnDataTypeMapping()", fieldType.String()))
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
					panic(fmt.Errorf("unsupported default value from %s, check GetGetDefaultValueByFromDbTag()", col.Default))
				}
			}

			stmt := fmt.Sprintf("ALTER TABLE %s ADD %s", m.Quote(tableName), colDef)
			scripts = append(scripts, stmt)

			// Update schema cache
			schema.Tables[tableName][col.Name] = true
		}
	}

	return strings.Join(scripts, ";\n"), nil
}
