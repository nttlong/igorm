package migrate

import (
	"fmt"
	"strings"
)

func (m *migratorMssql) GetSqlAddForeignKey() ([]string, error) {
	ret := []string{}
	schema, err := m.loader.LoadFullSchema(m.db)
	if err != nil {
		return nil, err
	}

	for fk, info := range ForeignKeyRegistry.fkMap {
		if _, ok := schema.ForeignKeys[fk]; !ok {

			formCols := "[" + strings.Join(info.FromCols, "],[") + "]"
			toCols := "[" + strings.Join(info.ToCols, "],[") + "]"
			script := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)", m.Quote(info.FromTable), m.Quote(fk), formCols, m.Quote(info.ToTable), toCols)
			schema.ForeignKeys[fk] = DbForeignKeyInfo{
				ConstraintName: fk,
				Table:          info.ToTable,
				Columns:        info.FromCols,
				RefTable:       info.ToTable,
				RefColumns:     info.ToCols,
			}
			if info.Cascade.OnDelete {
				script += " ON DELETE CASCADE"
			}
			if info.Cascade.OnUpdate {
				script += " ON UPDATE CASCADE"
			}
			ret = append(ret, script)
		}
	}
	return ret, nil
}
