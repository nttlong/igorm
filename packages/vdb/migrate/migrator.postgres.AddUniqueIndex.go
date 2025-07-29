package migrate

import (
	"fmt"
	"reflect"
	"strings"
)

func (m *migratorPostgres) GetSqlAddUniqueIndex(typ reflect.Type) (string, error) {
	scripts := []string{}

	// Load current schema
	schema, err := m.loader.LoadFullSchema(m.db)
	if err != nil {
		return "", err
	}

	// Get registered entity
	entityItem := ModelRegistry.GetModelByType(typ)
	if entityItem == nil {
		return "", fmt.Errorf("model %s not found, please register model first by call ModelRegistry.Add(%s)", typ.String(), typ.String())
	}

	// Duyệt các unique constraint
	for _, cols := range entityItem.entity.getUniqueConstraints() {
		var colNames []string
		var colNameInConstraint []string

		for _, col := range cols {
			colNames = append(colNames, m.Quote(col.Name))
			colNameInConstraint = append(colNameInConstraint, col.Name)
		}

		constraintName := fmt.Sprintf("UQ_%s__%s", entityItem.tableName, strings.Join(colNameInConstraint, "_"))

		// Nếu chưa có trong schema
		if _, ok := schema.UniqueKeys[constraintName]; !ok {
			sql := fmt.Sprintf(
				`ALTER TABLE %s ADD CONSTRAINT %s UNIQUE (%s)`,
				m.Quote(entityItem.tableName),
				m.Quote(constraintName),
				strings.Join(colNames, ", "),
			)
			scripts = append(scripts, sql)
		}
	}

	if len(scripts) == 0 {
		return "", nil
	}

	return strings.Join(scripts, ";\n") + ";", nil
}
