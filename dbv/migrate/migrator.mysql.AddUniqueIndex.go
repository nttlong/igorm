package migrate

import (
	"fmt"
	"reflect"
	"strings"
)

func (m *migratorMySql) GetSqlAddUniqueIndex(typ reflect.Type) (string, error) {
	scripts := []string{}

	// Load database schema hiện tại
	schema, err := m.loader.LoadFullSchema(m.db)
	if err != nil {
		return "", err
	}

	// Lấy entity đã đăng ký
	entityItem := ModelRegistry.GetModelByType(typ)
	if entityItem == nil {
		return "", fmt.Errorf("model %s not found, please register model first by call ModelRegistry.Add(%s)", typ.String(), typ.String())
	}

	uk := entityItem.entity.getUniqueConstraints()

	for _, cols := range uk {
		var colNames []string
		var colNameInConstraint []string
		for _, col := range cols {
			colNames = append(colNames, m.Quote(col.Name))
			colNameInConstraint = append(colNameInConstraint, col.Name)
		}

		constraintName := fmt.Sprintf("UQ_%s__%s", entityItem.tableName, strings.Join(colNameInConstraint, "___"))

		if _, ok := schema.UniqueKeys[constraintName]; !ok {
			script := fmt.Sprintf(
				"CREATE UNIQUE INDEX %s ON %s (%s)",
				m.Quote(constraintName),
				m.Quote(entityItem.tableName),
				strings.Join(colNames, ", "),
			)
			scripts = append(scripts, script)
		}
	}

	return strings.Join(scripts, ";\n"), nil
}
