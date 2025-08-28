package migrate

import (
	"fmt"
	"reflect"
)

func (m *migratorPostgres) createCheckLenConstraint(tableName string, col ColumnDef) string {
	if col.Field.Type == reflect.TypeFor[string]() || col.Field.Type == reflect.TypeFor[*string]() {
		checkSyntax := fmt.Sprintf("CHECK (char_length(%s) <= %d)", m.Quote(col.Name), *col.Length)

		constraintCheckName := fmt.Sprintf("%s_chk_%s_length", tableName, col.Name)
		sqlCreateCheckLen := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s", m.Quote(tableName), m.Quote(constraintCheckName), checkSyntax)
		return sqlCreateCheckLen
	}
	return ""

}
