package migrate

import (
	"eorm/tenantDB"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type migratorMssql struct {
	loader             IMigratorLoader
	cacheGetFullScript sync.Map

	db *tenantDB.TenantDB
}

func (m *migratorMssql) Quote(names ...string) string {
	return "[" + strings.Join(names, "].[") + "]"
}
func (m *migratorMssql) GetGetDefaultValueByFromDbTag() map[string]string {
	return map[string]string{
		"now": "GETDATE()",
	}
}
func (m *migratorMssql) GetColumnDataTypeMapping() map[reflect.Type]string {
	return map[reflect.Type]string{
		reflect.TypeOf(""):          "nvarchar",
		reflect.TypeOf(int(0)):      "int",
		reflect.TypeOf(int8(0)):     "tinyint",
		reflect.TypeOf(int16(0)):    "smallint",
		reflect.TypeOf(int32(0)):    "int",
		reflect.TypeOf(int64(0)):    "bigint",
		reflect.TypeOf(uint(0)):     "bigint", // SQL Server doesn't support unsigned
		reflect.TypeOf(uint8(0)):    "tinyint",
		reflect.TypeOf(uint16(0)):   "int",
		reflect.TypeOf(uint32(0)):   "bigint",
		reflect.TypeOf(uint64(0)):   "bigint",
		reflect.TypeOf(float32(0)):  "real",
		reflect.TypeOf(float64(0)):  "float",
		reflect.TypeOf(bool(false)): "bit",
		reflect.TypeOf([]byte{}):    "varbinary",
		reflect.TypeOf(time.Time{}): "datetime2",
		reflect.TypeOf(uuid.UUID{}): "uniqueidentifier",
	}
}
func (m *migratorMssql) GetSqlCreateTable(typ reflect.Type) (string, error) {
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
		return "", fmt.Errorf("model %s not found, please register model first by call ModelRegistry.Add(%s)", typ.String(), typ.String())
	}
	if entityItem == nil {
		return "", fmt.Errorf("model %s not found", typ.Name())
	}
	tableName := entityItem.tableName
	if _, ok := schema.Tables[tableName]; ok {
		return "", nil
	}

	// Nếu bảng đã tồn tại → không tạo
	if _, ok := schema.Tables[tableName]; ok {
		return "", nil
	}

	// Danh sách các cột để tạo bảng
	strCols := []string{}
	newTableMap := map[string]bool{}
	for _, col := range entityItem.entity.cols {
		newTableMap[col.Name] = true
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
			defaultVal := ""
			if val, ok := defaultValueByFromDbTag[col.Default]; ok {
				defaultVal = val
			} else {
				err = fmt.Errorf("not support default value from %s, review GetGetDefaultValueByFromDbTag() function in %s ", col.Default, reflect.TypeOf(m).Elem())
				panic(err)
			}

			colDef += fmt.Sprintf(" DEFAULT %s", defaultVal)
		}

		strCols = append(strCols, colDef)
	}

	// Xử lý PRIMARY KEY constraint
	for _, cols := range entityItem.entity.primaryConstraints {
		//var colNames []string
		// colNameInConstraint := []string{}

		var pkCols []string
		var pkColName []string
		for _, col := range cols {
			if col.PKName != "" {
				pkCols = append(pkCols, m.Quote(col.Name))
				pkColName = append(pkColName, col.Name)
			}
		}
		pkConstraintName := ""
		if len(pkCols) > 0 {
			// Constraint name theo chuẩn PK_<table>__<col1>_<col2>
			pkConstraintName = fmt.Sprintf("PK_%s__%s", tableName, strings.Join(pkColName, "_"))
			constraint := fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", m.Quote(pkConstraintName), strings.Join(pkCols, ", "))
			strCols = append(strCols, constraint)
		}
		// constraintName = fmt.Sprintf("PK_%s__%s", tableName, strings.Join(colNameInConstraint, "___"))
		//constraint := fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", m.Quote(pkConstraintName), strings.Join(colNames, ", "))
		strCols = append(strCols)

	}

	// Kết hợp thành câu lệnh CREATE TABLE
	sql := fmt.Sprintf("CREATE TABLE %s (\n  %s\n)", m.Quote(tableName), strings.Join(strCols, ",\n  "))
	schema.Tables[tableName] = newTableMap

	return sql, nil
}

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
		return "", fmt.Errorf("model %s not found, please register model first by call ModelRegistry.Add(%s)", typ.String(), typ.String())
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
				defaultVal := ""
				if val, ok := defaultValueByFromDbTag[col.Default]; ok {
					defaultVal = val
				} else {
					err = fmt.Errorf("not support default value from %s, review GetGetDefaultValueByFromDbTag() function in %s ", col.Default, reflect.TypeOf(m).Elem())
					panic(err)
				}

				colDef += fmt.Sprintf(" DEFAULT %s", defaultVal)
			}

			scripts = append(scripts, fmt.Sprintf("ALTER TABLE %s ADD %s", m.Quote(entityItem.tableName), colDef))

			schema.Tables[entityItem.tableName][col.Name] = true
		}
	}

	return strings.Join(scripts, ";\n"), nil

}
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
		return "", fmt.Errorf("model %s not found, please register model first by call ModelRegistry.Add(%s)", typ.String(), typ.String())
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
func (m *migratorMssql) GetSqlMigrate(entityType reflect.Type) ([]string, error) {
	scripts := []string{}
	scriptTable, err := m.GetSqlCreateTable(entityType)
	if err != nil {
		return nil, err
	}
	if scriptTable == "" {
		scriptAddColumn, err := m.GetSqlAddColumn(entityType)
		if err != nil {
			return nil, err
		}
		scripts = append(scripts, scriptTable, scriptAddColumn)
	}

	scriptAddUniqueIndex, err := m.GetSqlAddUniqueIndex(entityType)
	if err != nil {
		return nil, err
	}
	scripts = append(scripts, scriptTable, scriptAddUniqueIndex)
	return scripts, nil

}
func (m *migratorMssql) DoMigrate(entityType reflect.Type) error {
	scripts, err := m.GetSqlMigrate(entityType)
	if err != nil {
		return err
	}
	for _, script := range scripts {
		_, err := m.db.Exec(script)
		if err != nil {
			return err
		}
	}
	return nil

}

type initDoMigrates struct {
	once sync.Once
}

var cacheDoMigrates sync.Map

func (m *migratorMssql) DoMigrates() error {
	key := fmt.Sprintf("%s_%s", m.db.GetDBName(), m.db.GetDbType())
	actual, _ := cacheDoMigrates.LoadOrStore(key, &initDoMigrates{})

	mi := actual.(*initDoMigrates)
	var err error
	mi.once.Do(func() {

		scripts, err := m.GetFullScript()
		if err != nil {
			return
		}
		for _, script := range scripts {
			_, err := m.db.Exec(script)
			if err != nil {
				break
			}
		}
		// for _, entity := range ModelRegistry.GetAllModels() {
		// 	err = m.DoMigrate(entity.entity.entityType)

		// }
	})
	return err
}
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
			ret = append(ret, script)
		}
	}
	return ret, nil
}

type getFullScriptInit struct {
	once sync.Once
	ret  []string
}

func (m *migratorMssql) GetFullScript() ([]string, error) {
	key := fmt.Sprintf("%s_%s", m.db.GetDBName(), m.db.GetDbType())
	actual, _ := m.cacheGetFullScript.LoadOrStore(key, &getFullScriptInit{})
	init := actual.(*getFullScriptInit)
	var err error
	init.once.Do(func() {
		init.ret, err = m.getFullScript()
	})
	return init.ret, err
}
func (m *migratorMssql) getFullScript() ([]string, error) {

	scripts := []string{}
	for _, entity := range ModelRegistry.GetAllModels() {
		script, err := m.GetSqlCreateTable(entity.entity.entityType)
		if err != nil {
			return nil, err
		}
		if script != "" {
			scripts = append(scripts, script)
		}

	}
	for _, entity := range ModelRegistry.GetAllModels() {
		script, err := m.GetSqlAddColumn(entity.entity.entityType)
		if err != nil {
			return nil, err
		}
		if script != "" {
			scripts = append(scripts, script)
		}
	}
	scriptForeignKey, err := m.GetSqlAddForeignKey()
	if err != nil {
		return nil, err
	}
	scripts = append(scripts, scriptForeignKey...)

	return scripts, nil
}
