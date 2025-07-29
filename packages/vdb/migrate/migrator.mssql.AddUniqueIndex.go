package migrate

import (
	"fmt"
	"reflect"
	"strings"
)

func (m *migratorMssql) GetSqlAddUniqueIndex(typ reflect.Type) (string, error) {
	scripts := []string{}

	// Load database schema hiện tại
	schema, err := m.loader.LoadFullSchema(m.db)
	if err != nil {
		return "", err
	}

	// Lấy entity đã đăng ký
	entityItem := ModelRegistry.GetModelByType(typ)
	uk := entityItem.entity.getUniqueConstraints()
	for _, cols := range uk {
		var colNames []string
		colNameInConstraint := []string{}
		for _, col := range cols {
			colNames = append(colNames, m.Quote(col.Name))
			colNameInConstraint = append(colNameInConstraint, col.Name)
		}
		constraintName := fmt.Sprintf("UQ_%s__%s", entityItem.tableName, strings.Join(colNameInConstraint, "___"))
		if _, ok := schema.UniqueKeys[constraintName]; !ok {
			constraint := fmt.Sprintf("CONSTRAINT %s UNIQUE (%s)", m.Quote(constraintName), strings.Join(colNames, ", "))
			script := fmt.Sprintf("ALTER TABLE %s ADD %s", m.Quote(entityItem.tableName), constraint)
			scripts = append(scripts, script)
		}
	}
	return strings.Join(scripts, ";\n"), nil

}
