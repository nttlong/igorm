package migrate

import (
	"fmt"
	"reflect"
	"strings"
)

func (m *migratorMssql) GetSqlAddColumn(typ reflect.Type) (string, error) {
	mapType := m.GetColumnDataTypeMapping()
	defaultValueByFromDbTag := m.GetGetDefaultValueByFromDbTag()

	// Load database schema hiện tại
	schema, err := m.loader.LoadFullSchema(m.db)
	if err != nil {
		return "", err
	}

	// Lấy entity đã đăng ký
	entityItem := ModelRegistry.GetModelByType(typ)
	if entityItem == nil {
		return "", NewModelError(typ)
	}
	scripts := []string{}
	for _, col := range entityItem.entity.cols {
		if _, ok := schema.Tables[entityItem.tableName][col.Name]; !ok {
			fieldType := col.Field.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}

			sqlType := mapType[fieldType]
			if col.Length != nil {
				sqlType = fmt.Sprintf("%s(%d)", sqlType, *col.Length)
			}

			colDef := m.Quote(col.Name) + " " + sqlType

			if col.IsAuto {
				colDef += " IDENTITY(1,1)"
			}

			if col.Nullable {
				colDef += " NULL"
			} else {
				colDef += " NOT NULL"
			}

			if col.Default != "" {
				df, err := typeUtils.GetDefaultValue(col.Default, defaultValueByFromDbTag)
				if err != nil {
					err = fmt.Errorf("not support default value from %s, review GetGetDefaultValueByFromDbTag() function in %s ", col.Default, "vdb/migrate/migrator.mssql.AddColumn.go")
					panic(err)
				}
				colDef += " DEFAULT " + df

				colDef += fmt.Sprintf(" DEFAULT %s", colDef)
			}

			scripts = append(scripts, fmt.Sprintf("ALTER TABLE %s ADD %s", m.Quote(entityItem.tableName), colDef))

			schema.Tables[entityItem.tableName][col.Name] = true
		}
	}

	return strings.Join(scripts, ";\n"), nil

}
