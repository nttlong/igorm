package migrate

import (
	"fmt"
	"reflect"
	"strings"
)

func (m *migratorPostgres) GetSqlAddIndex(typ reflect.Type) (string, error) {
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

	for _, cols := range entityItem.entity.getIndexConstraints() {
		var colNames []string
		var colNameInConstraint []string

		for _, col := range cols {
			colNames = append(colNames, m.Quote(col.Name))
			colNameInConstraint = append(colNameInConstraint, col.Name)
		}

		constraintName := fmt.Sprintf("IDX_%s__%s", entityItem.tableName, strings.Join(colNameInConstraint, "_"))

		// Nếu index chưa tồn tại trong schema
		if _, ok := schema.Indexes[constraintName]; !ok {
			// PostgreSQL mặc định dùng BTREE, có thể thêm USING nếu cần
			sql := fmt.Sprintf(
				"CREATE INDEX %s ON %s (%s)",
				m.Quote(constraintName),
				m.Quote(entityItem.tableName),
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
