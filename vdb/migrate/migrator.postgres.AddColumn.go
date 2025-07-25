package migrate

import (
	"fmt"
	"reflect"
	"strings"
)

func (m *migratorPostgres) GetSqlAddColumn(typ reflect.Type) (string, error) {
	mapType := m.GetColumnDataTypeMapping()
	defaultValueByFromDbTag := m.GetGetDefaultValueByFromDbTag()

	// Load current schema
	schema, err := m.loader.LoadFullSchema(m.db)
	if err != nil {
		return "", err
	}

	// Get registered model
	entityItem := ModelRegistry.GetModelByType(typ)
	if entityItem == nil {
		return "", fmt.Errorf("model %s not found, please register model first by call ModelRegistry.Add(%s)", typ.String(), typ.String())
	}

	scripts := []string{}
	checkLengthScripts := []string{}

	for _, col := range entityItem.entity.cols {
		// Column chưa tồn tại thì mới thêm
		if _, ok := schema.Tables[entityItem.tableName][col.Name]; !ok {
			fieldType := col.Field.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}

			sqlType, ok := mapType[fieldType]
			if !ok {
				return "", fmt.Errorf("unsupported field type %s, check GetColumnDataTypeMapping", fieldType.String())
			}

			if col.Length != nil {
				checkLengthScripts = append(checkLengthScripts, m.createCheckLenConstraint(entityItem.tableName, col))
			}

			colDef := m.Quote(col.Name)

			// Xử lý auto increment trong PostgreSQL
			if col.IsAuto {
				if fieldType.Kind() == reflect.Int || fieldType.Kind() == reflect.Int64 {
					colDef += " BIGSERIAL"
				} else {
					colDef += fmt.Sprintf(" %s GENERATED ALWAYS AS IDENTITY", sqlType)
				}
			} else {
				colDef += " " + sqlType
			}

			if !col.Nullable {
				colDef += " NOT NULL"
			}

			if col.Default != "" {
				defaultVal := ""
				if val, ok := defaultValueByFromDbTag[col.Default]; ok {
					defaultVal = val
				} else {
					panic(fmt.Errorf("unsupported default tag: %s in %s", col.Default, reflect.TypeOf(m).Elem()))
				}
				colDef += fmt.Sprintf(" DEFAULT %s", defaultVal)
			}

			script := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s", m.Quote(entityItem.tableName), colDef)
			scripts = append(scripts, script)

			// Update schema cache
			schema.Tables[entityItem.tableName][col.Name] = true
		}
	}

	if len(scripts) == 0 {
		return "", nil
	}
	scripts = append(scripts, checkLengthScripts...)

	return strings.Join(scripts, ";\n"), nil
}
