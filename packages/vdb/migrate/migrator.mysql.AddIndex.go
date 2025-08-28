package migrate

import (
	"fmt"
	"reflect"
	"strings"
)

func (m *migratorMySql) GetSqlAddIndex(typ reflect.Type) (string, error) {
	scripts := []string{}

	// Load schema hiện tại
	schema, err := m.loader.LoadFullSchema(m.db)
	if err != nil {
		return "", err
	}

	// Lấy entity đã đăng ký
	entityItem := ModelRegistry.GetModelByType(typ)
	if entityItem == nil {
		return "",  NewModelError(typ)
	}

	for _, cols := range entityItem.entity.getIndexConstraints() {
		var colNames []string
		var colNameInConstraint []string

		for _, col := range cols {
			colNames = append(colNames, m.Quote(col.Name))
			colNameInConstraint = append(colNameInConstraint, col.Name)
		}

		constraintName := fmt.Sprintf("IDX_%s__%s", entityItem.tableName, strings.Join(colNameInConstraint, "_"))

		// Nếu chưa tồn tại index này trong schema
		if _, ok := schema.Indexes[constraintName]; !ok {
			stmt := fmt.Sprintf(
				"CREATE INDEX %s ON %s (%s)",
				m.Quote(constraintName),
				m.Quote(entityItem.tableName),
				strings.Join(colNames, ", "),
			)
			scripts = append(scripts, stmt)
		}
	}

	return strings.Join(scripts, ";\n"), nil
}
