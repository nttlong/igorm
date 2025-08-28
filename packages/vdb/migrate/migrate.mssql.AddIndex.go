package migrate

import (
	"fmt"
	"reflect"
	"strings"
)

var NewModelError func(typ reflect.Type) error

func (m *migratorMssql) GetSqlAddIndex(typ reflect.Type) (string, error) {
	scripts := []string{}

	// Load database schema hiện tại
	schema, err := m.loader.LoadFullSchema(m.db)
	if err != nil {
		return "", err
	}
	fmt.Println(typ.String())
	// Lấy entity đã đăng ký
	entityItem := ModelRegistry.GetModelByType(typ)
	if entityItem == nil {
		return "", NewModelError(typ)
	}
	for _, cols := range entityItem.entity.getIndexConstraints() {
		var colNames []string
		colNameInConstraint := []string{}
		for _, col := range cols {
			colNames = append(colNames, m.Quote(col.Name))
			colNameInConstraint = append(colNameInConstraint, col.Name)
		}
		constraintName := fmt.Sprintf("IDX_%s__%s", entityItem.tableName, strings.Join(colNameInConstraint, "_"))
		if _, ok := schema.UniqueKeys[constraintName]; !ok {
			constraint := fmt.Sprintf("CREATE INDEX %s ON %s (%s)", m.Quote(constraintName), m.Quote(entityItem.tableName), strings.Join(colNames, ", "))
			scripts = append(scripts, constraint)

		}
	}
	return strings.Join(scripts, ";\n"), nil

}
