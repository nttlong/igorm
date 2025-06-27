package unvsef

import "database/sql"

/*
This struct is definition of tenant db
*/
type TenantDb struct {
	DB         *sql.DB
	Dialect    Dialect
	DBType     DBType
	DBTypeName string
	SqlMigrate []string
	DbName     string
}

func (t *TenantDb) DoMigrate() error {
	for _, sql := range t.SqlMigrate {
		_, err := t.DB.Exec(sql)
		if err != nil {
			return err
		}
	}
	return nil
}
