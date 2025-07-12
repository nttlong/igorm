package migrate

import (
	"eorm/tenantDB"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
)

type migratorMssql struct {
	loader IMigratorLoader
	db     *tenantDB.TenantDB
}

func (m *migratorMssql) Quote(names ...string) string {
	return "[" + strings.Join(names, "].[") + "]"
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

	// Load database schema hiện tại
	schema, err := m.loader.LoadFullSchema(m.db.DB)
	if err != nil {
		return "", err
	}

	// Lấy entity đã đăng ký
	entityItem := ModelRegistry.GetModelByType(typ)
	if entityItem == nil {
		return "", fmt.Errorf("model %s not found", typ.Name())
	}
	tableName := entityItem.tableName

	// Nếu bảng đã tồn tại → không tạo
	if _, ok := schema.Tables[tableName]; ok {
		return "", nil
	}

	// Danh sách các cột để tạo bảng
	strCols := []string{}
	for _, col := range entityItem.entity.cols {
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
			defaultVal := col.Default
			if strings.EqualFold(defaultVal, "now") {
				defaultVal = "GETDATE()"
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
	return sql, nil
}

func (m *migratorMssql) GetSqlAddColumn(typ reflect.Type) (string, error) {
	panic("not implemented")

}
func (m *migratorMssql) GetSqlAddIndex(typ reflect.Type) (string, error) {
	panic("not implemented")

}
func (m *migratorMssql) GetSqlAddUniqueIndex(typ reflect.Type) (string, error) {
	panic("not implemented")

}
