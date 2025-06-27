package unvsef

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/jinzhu/inflection"
)

type DBType string

const (
	DBPostgres  DBType = "postgres"
	DBMySQL     DBType = "mysql"
	DBMariaDB   DBType = "mariadb"
	DBMSSQL     DBType = "sqlserver"
	DBSQLite    DBType = "sqlite"
	DBOracle    DBType = "oracle"
	DBTiDB      DBType = "tidb"
	DBCockroach DBType = "cockroach"
	DBGreenplum DBType = "greenplum"
	DBUnknown   DBType = "unknown"
)

// --------------------- Tag Metadata ---------------------
// FieldTag holds parsed metadata from struct field tags.
type FieldTag struct {
	PrimaryKey    bool
	AutoIncrement bool
	Unique        bool
	/*
		can be field name if no unique index name in tag else name of unique index in tag
	*/
	UniqueName string
	Index      bool
	/*
		can be field name if no  index name in tag else name of  index in tag
	*/
	IndexName string
	Length    *int
	FTSName   string
	DBType    string
	TableName string
	Check     string
	Nullable  bool
	Field     reflect.StructField
	Default   string
}
type utilsPackage struct {
	cacheGetMetaInfo                        sync.Map
	CacheTableNameFromStruct                sync.Map
	cacheGetPkFromMeta                      sync.Map
	cacheGetUniqueConstraintsFromMetaByType sync.Map
	cacheGetIndexConstraintsFromMetaByType  sync.Map
	schemaCache                             sync.Map
	// future: add cache or shared state here
}

var utils = &utilsPackage{}

// ParseDBTag parses the `db` struct tag into a FieldTag struct.
func (u *utilsPackage) ParseDBTag(field reflect.StructField) FieldTag {

	tag := strings.TrimSpace(field.Tag.Get("db"))
	t := FieldTag{
		Field: field,
	}

	t.Nullable = strings.Contains(field.Type.String(), ".DbField[*")
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return t
	}
	parts := strings.Split(tag, ";")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		switch {

		case p == "primaryKey":
			t.PrimaryKey = true
		case p == "autoIncrement":
			t.AutoIncrement = true
		case p == "unique":
			t.Unique = true
			t.UniqueName = u.ToSnakeCase(field.Name)
		case strings.HasPrefix(p, "unique("):
			t.Unique = true
			t.UniqueName = u.extractName(p)
		case p == "index":
			t.Index = true
			t.IndexName = u.ToSnakeCase(field.Name)
		case strings.HasPrefix(p, "index("):
			t.Index = true
			t.IndexName = u.extractName(p)
		case strings.HasPrefix(p, "table("):
			t.TableName = u.extractName(p)
		case strings.HasPrefix(p, "length("):
			if s := u.extractName(p); s != "" {
				if n, err := strconv.Atoi(s); err == nil {
					t.Length = &n
				}
			}
		case strings.HasPrefix(p, "check("):
			t.Check = u.extractName(p)

			if s := u.extractName(p); s != "" {
				if n, err := strconv.Atoi(s); err == nil {
					t.Length = &n
				}
			}
		case strings.HasPrefix(p, "FTS("):
			t.FTSName = u.extractName(p)
		case strings.HasPrefix(p, "type:"):
			t.DBType = strings.TrimPrefix(p, "type:")
		case strings.HasPrefix(p, "default:"):
			t.Default = strings.TrimPrefix(p, "default:")
		}

	}
	return t
}

func (u *utilsPackage) extractName(s string) string {
	re := regexp.MustCompile(`\((.*?)\)`)
	matches := re.FindStringSubmatch(s)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}
func (u *utilsPackage) Contains(list []string, item string) bool {
	item = strings.ToLower(item)
	for _, v := range list {
		if strings.ToLower(v) == item {
			return true
		}
	}
	return false
}
func (u *utilsPackage) ToSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) &&
			(unicode.IsLower(rune(str[i-1])) || (i+1 < len(str) && unicode.IsLower(rune(str[i+1])))) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

/*
Get table name from struct

	 Example :

	 type User struct {
		  _ DbField[any] `db:table(my_user)`

	 }

	 return my_user

	 type User struct {} -> "users"
*/
func (u *utilsPackage) TableNameFromStruct(typ reflect.Type) string {
	// Check for override via table(...) tag
	if v, ok := u.CacheTableNameFromStruct.Load(typ); ok {
		return v.(string)
	}
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Name == "_" {
			parsed := u.ParseDBTag(typ.Field(i))
			if parsed.TableName != "" {
				return parsed.TableName
			}
		}
	}
	base := u.ToSnakeCase(typ.Name())
	ret := inflection.Plural(base)
	u.CacheTableNameFromStruct.Store(typ, ret)
	return ret
}

/*
	   Get Go type of DbField in name

	   Example:
	    var testType= DbField[time.Time]{}
		utils.ResolveFieldKind(reflect.TypeOf(&testType))-> "time.Time"
*/
func (u *utilsPackage) ResolveFieldKind(field reflect.StructField) string {
	strFt := field.Type.String()

	if strings.Contains(strFt, ".DbField[") {
		typeParam := strings.Split(strFt, ".DbField[")[1]
		typeParam = strings.Split(typeParam, "]")[0]
		typeParam = strings.TrimLeft(typeParam, "*")
		return typeParam
	}
	return field.Type.Kind().String()
}

/*
Extract all info of reflect.Type
After fetch all info in reflect.Type the meta info will be cache
Next call will get from cache instead of fetch again

# Example 1:

	type User struct {
			_ DbField[any] `db:"table(MyUser)"` //<-- Optional
			Id   DbField[uint64] `db:"primaryKey;autoIncrement"`
			Code DbField[string] `db:"unique;length(50)"`
			Name DbField[string] `db:"index;length(50)"`
		}

		return map {
				users: {
					id:{
						PrimaryKey    : true
					},
					code: {
							Unique:true,
							UniqueName: code_uk
					},
					name: {
							Index: true,
							IndexName; name_uk
					}
				}
		}
		# Exmaple 2
		type struct UserRole {
			Id   DbField[uint64]     `db:"primaryKey;autoIncrement"`
			UserId  DbField[uint64]  `db:"unique(user_role)"` //<-- unique constraint 2 columns
			RoleId  DbField[uint64]  `db:"unique(user_role)"` //<-- unique constraint 2 columns
		}
		return
		{
			users: {
				id: {
					PrimaryKey: true
				},
				user_id: {
					Unique: true,
					UniqueName: user_role_uk
				},
				role_id: {
					Unique: true,
					UniqueName: user_role_uk
				}
			}
		}
*/
func (u *utilsPackage) GetMetaInfo(typ reflect.Type) map[string]map[string]FieldTag {
	// 1. Kiểm tra cache trước (check cache first)
	if metaInfo, ok := u.cacheGetMetaInfo.Load(typ); ok {
		return metaInfo.(map[string]map[string]FieldTag)
	}

	// 2. Tạo mới metadata
	metaInfo := make(map[string]map[string]FieldTag)
	tableName := utils.TableNameFromStruct(typ)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Bỏ qua field đặc biệt "_" (used for table(...) override)
		if field.Name == "_" {

			continue
		}

		// 3. Nếu là embedded struct (anonymous), đệ quy lấy metadata của nó
		if field.Anonymous {
			embeddedMeta := u.GetMetaInfo(field.Type)
			for tableName, fields := range embeddedMeta {
				if _, ok := metaInfo[tableName]; !ok {
					metaInfo[tableName] = make(map[string]FieldTag)
				}
				for fieldName, fieldTag := range fields {

					metaInfo[tableName][u.ToSnakeCase(fieldName)] = fieldTag
				}
			}
			continue
		}

		if _, ok := metaInfo[tableName]; !ok {
			metaInfo[tableName] = make(map[string]FieldTag)
		}

		// 5. Gán tag metadata cho field
		metaInfo[tableName][u.ToSnakeCase(field.Name)] = u.ParseDBTag(field)
	}

	// 6. Cache lại và trả về
	u.cacheGetMetaInfo.Store(typ, metaInfo)
	return metaInfo
}

/*
The purpose support for Dialects is to provide a way to customize the SQL generation for different databases.
*/
func (u *utilsPackage) Quote(strQuote string, str ...string) string {
	left := strQuote[0:1]
	right := strQuote[1:2]
	ret := left + strings.Join(str, left+"."+right) + right
	return ret
}

/*
The function will get all primary key cols from type (the  info will be stored in cache . The next call will return form cache)
return

	{
		"<primary key constraint name": {//<-- The value is obtained by the combination of the table name, double underscores ("__"), and the column name.
			"<key field name 1>": {...},//<--FieldTag Info
			...
			"<key field name n>": {...},//<--FieldTag Info
		}
	}
*/
func (u *utilsPackage) GetPkFromMetaByType(typ reflect.Type) map[string]map[string]FieldTag {
	// 1. Kiểm tra cache trước (check cache first)
	if pk, ok := u.cacheGetPkFromMeta.Load(typ); ok {
		return pk.(map[string]map[string]FieldTag)
	}
	metaInfo := u.GetMetaInfo(typ)
	ret := make(map[string]map[string]FieldTag)
	fieldsNames := []string{}
	for tableName, fields := range metaInfo {
		pkMap := make(map[string]FieldTag)
		for fieldName, fieldTag := range fields {
			if fieldTag.PrimaryKey {
				pkMap[fieldName] = fieldTag
				fieldsNames = append(fieldsNames, fieldName)

			}
		}
		ret[tableName+"_"+strings.Join(fieldsNames, "_")] = pkMap

	}
	u.cacheGetPkFromMeta.Store(typ, ret)
	return ret
}

/*
The method will get all Unique Constraint from declarification type
#Exmaple

	type struct UserRole {
			Id   DbField[uint64]     `db:"primaryKey;autoIncrement"`
			UserId  DbField[uint64]  `db:"unique(user_role_idx)"` //<-- unique constraint 2 columns
			RoleId  DbField[uint64]  `db:"unique(user_role_idx)"` //<-- unique constraint 2 columns
			OwnerId DbField[uint64] `db:"unique"` //<-- unique constraint 1 column
	}

	return {
		"user_roles":{ //<-- table name
			"user_role_idx____user_roles___user_id__role_id":{ //<--- convention of constraint name is the combination of constraint name in tag, four underscores and table name (snake case)
				"user_id": {...} //<-- col 1
				"role_id": {...} // <--col 2
			},
			"user_role_idx____user_roles__owner_id": { //<--- convention of constraint name is the combination of field name (snake case) four underscore and table name (snake case)
				"owner_id": {...} //<-- field name
			}

		}
	}
*/
func (u *utilsPackage) GetUniqueConstraintsFromMetaByType(typ reflect.Type) map[string]map[string]FieldTag {
	// 1. Kiểm tra cache trước (check cache first)
	if unique, ok := u.cacheGetUniqueConstraintsFromMetaByType.Load(typ); ok {
		return unique.(map[string]map[string]FieldTag)
	}
	metaInfo := u.GetMetaInfo(typ)
	info := make(map[string]map[string]FieldTag)
	tableName := u.TableNameFromStruct(typ)

	for _, fields := range metaInfo {
		for fieldName, fieldTag := range fields {

			if fieldTag.Unique {
				ukName := fieldTag.UniqueName //<-- use fieldTag.UniqueName can be field name if no unique index name in tag else name of unique index in tag
				if _, ok := info[ukName]; !ok {
					info[ukName] = make(map[string]FieldTag)

				}
				info[ukName][fieldName] = fieldTag

			}
		}
	}
	ret := make(map[string]map[string]FieldTag)
	for ukName, fields := range info {

		refFields := []string{}
		for fieldName := range fields {
			refFields = append(refFields, fieldName)
		}
		constraintName := ukName + "____" + tableName + "___" + strings.Join(refFields, "__")
		ret[constraintName] = fields
	}
	u.cacheGetUniqueConstraintsFromMetaByType.Store(typ, ret)
	return ret
}

/*
The method will get all Unique Constraint from declarification type
#Exmaple

	type struct UserRole {
			Id   DbField[uint64]     `db:"primaryKey;autoIncrement"`
			UserId  DbField[uint64]  `db:"index(user_role_idx)"` //<-- unique constraint 2 columns
			RoleId  DbField[uint64]  `db:"index(user_role_idx)"` //<-- unique constraint 2 columns
			OwnerId DbField[uint64] `db:"index"` //<-- unique constraint 1 column
	}

	return {
		"user_roles":{ //<-- table name
			"user_role_idx____user_roles":{ //<--- convention of constraint name is constraint name in tag +"___"+ table name
				"user_id": {...} //<-- col 1
				"role_id": {...} // <--col 2
			},
			"user_roles____owner_id_idx": { //<--- constraint name
				"owner_id": {...} //<-- field name
			}

		}
	}
*/
func (u *utilsPackage) GetIndexConstraintsFromMetaByType(typ reflect.Type) map[string]map[string]FieldTag {
	// 1. Kiểm tra cache trước (check cache first)
	if unique, ok := u.cacheGetIndexConstraintsFromMetaByType.Load(typ); ok {
		return unique.(map[string]map[string]FieldTag)
	}
	metaInfo := u.GetMetaInfo(typ)
	info := make(map[string]map[string]FieldTag)
	tableName := u.TableNameFromStruct(typ)

	for _, fields := range metaInfo {
		for fieldName, fieldTag := range fields {

			if fieldTag.Index {
				indexName := fieldTag.IndexName
				if _, ok := info[indexName]; !ok {
					info[indexName] = make(map[string]FieldTag)

				}
				info[indexName][fieldName] = fieldTag

			}
		}
	}
	ret := make(map[string]map[string]FieldTag)
	for idxName, fields := range info {

		refFields := []string{}
		for fieldName := range fields {
			refFields = append(refFields, fieldName)
		}
		constraintName := idxName + "____" + tableName + "___" + strings.Join(refFields, "__")
		ret[constraintName] = fields
	}
	u.cacheGetIndexConstraintsFromMetaByType.Store(typ, ret)
	return ret
}

type fkInfo struct {
	FromTable string
	FromField []string
	ToTable   string
	ToField   []string
}

type schemaMap struct {
	table  map[string]bool
	unique map[string]bool
	index  map[string]bool
	fk     map[string]bool
}

func (u *utilsPackage) extractSchema(db *sql.DB, dbName string, dialect Dialect) (*schemaMap, error) {

	dialect.RefreshSchemaCache(db, dbName)
	schema, err := dialect.GetSchema(db, dbName)
	if err != nil {
		return nil, err
	}
	ret := &schemaMap{
		table:  make(map[string]bool),
		unique: make(map[string]bool),
		index:  make(map[string]bool),
		fk:     make(map[string]bool),
	}
	for tableName, table := range schema {
		ret.table[tableName] = true

		for _, constraintName := range table.UniqueConstraints {
			ret.unique[constraintName] = true
		}
		for _, constraintName := range table.IndexConstraints {
			ret.index[constraintName] = true
		}
	}
	u.schemaCache.Store(dbName, ret)
	return ret, nil

}

func (u *utilsPackage) GetScriptMigrate(db *sql.DB, dbName string, dialect Dialect, typ ...*reflect.Type) ([]string, error) {
	dbSchema, err := utils.extractSchema(db, dbName, dialect)
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for _, t := range typ {
		sqlCmds, err := dialect.GenerateCreateTableSql(dbName, *t)
		if err != nil {
			return nil, err
		}
		if sqlCmds != "" {
			ret = append(ret, sqlCmds)
		} else { //<--- table is existing in Database, just add columns
			sqlAddCols, err := dialect.GenerateAlterTableSql(dbName, *t)
			if err != nil {
				return nil, err
			}
			ret = append(ret, sqlAddCols...)
		}
		sqlUniqueConstraints := dialect.GenerateUniqueConstraintsSql(*t)

		if err != nil {
			return nil, err
		}

		for constraintName, sql := range sqlUniqueConstraints {
			if !dbSchema.unique[constraintName] {
				ret = append(ret, sql)
			}
		}
		sqlIndexConstraints := dialect.GenerateIndexConstraintsSql(*t)
		if err != nil {
			return nil, err
		}
		for constraintName, sql := range sqlIndexConstraints {
			if !dbSchema.index[constraintName] {
				ret = append(ret, sql)
			}

		}

	}
	return ret, nil
}
func (u *utilsPackage) DetectDatabaseType(db *sql.DB) (DBType, string, error) {
	var version string

	queries := []struct {
		query string
	}{
		{"SELECT version();"},        // PostgreSQL, MySQL, Cockroach, Greenplum
		{"SELECT @@VERSION;"},        // SQL Server, Sybase
		{"SELECT sqlite_version();"}, // SQLite
		{"SELECT tidb_version();"},   // TiDB
		{"SELECT * FROM v$version"},  // Oracle
	}

	for _, q := range queries {
		err := db.QueryRow(q.query).Scan(&version)
		if err == nil && version != "" {
			v := strings.ToLower(version)

			switch {
			case strings.Contains(v, "postgres"):
				if strings.Contains(v, "greenplum") {
					return DBGreenplum, version, nil
				}
				return DBPostgres, version, nil
			case strings.Contains(v, "cockroach"):
				return DBCockroach, version, nil
			case strings.Contains(v, "mysql"):
				if strings.Contains(v, "mariadb") {
					return DBMariaDB, version, nil
				}
				return DBMySQL, version, nil
			case strings.Contains(v, "mariadb"):
				return DBMariaDB, version, nil
			case strings.Contains(v, "microsoft") || strings.Contains(v, "sql server"):
				return DBMSSQL, version, nil
			case strings.Contains(v, "sqlite"):
				return DBSQLite, version, nil
			case strings.Contains(v, "tidb"):
				return DBTiDB, version, nil
			case strings.Contains(v, "oracle"):
				return DBOracle, version, nil
			}
		}
	}

	return DBUnknown, version, errors.New("unable to detect database type")
}
func (u *utilsPackage) GetCurrentDatabaseName(db *sql.DB, dbType DBType) (string, error) {
	var query string
	var dbName string

	switch dbType {
	case DBPostgres, DBGreenplum, DBCockroach:
		query = "SELECT current_database();"
	case DBMySQL, DBMariaDB, DBTiDB:
		query = "SELECT DATABASE();"
	case DBMSSQL:
		query = "SELECT DB_NAME();"
	case DBSQLite:
		query = "PRAGMA database_list;" // SQLite đặc biệt hơn, xem dưới
	case DBOracle:
		query = "SELECT SYS_CONTEXT('USERENV','DB_NAME') FROM dual;"
	default:
		return "", fmt.Errorf("unsupported db type: %s", dbType)
	}

	if dbType == DBSQLite {
		type sqliteEntry struct {
			Seq  int
			Name string
			File string
		}
		rows, err := db.Query(query)
		if err != nil {
			return "", err
		}
		defer rows.Close()

		for rows.Next() {
			var entry sqliteEntry
			if err := rows.Scan(&entry.Seq, &entry.Name, &entry.File); err != nil {
				return "", err
			}
			if entry.Seq == 0 {
				return entry.Name, nil // thường là "main"
			}
		}
		return "", fmt.Errorf("no database found in sqlite PRAGMA list")
	}

	// Các DB bình thường chỉ cần query trả 1 giá trị
	err := db.QueryRow(query).Scan(&dbName)
	if err != nil {
		return "", err
	}

	return dbName, nil
}
