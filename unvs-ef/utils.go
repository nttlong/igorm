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
	cacheGetOrCreateRepository              sync.Map
	cacheDbFunctions                        sync.Map
	// future: add cache or shared state here
}

var utils = &utilsPackage{}

// ParseDBTag parses the `db` struct tag into a FieldTag struct.
func (u *utilsPackage) ParseDBTag(field reflect.StructField) FieldTag {

	tag := strings.TrimSpace(field.Tag.Get("db"))
	t := FieldTag{
		Field: field,
	}

	t.Nullable = field.Type.Kind() == reflect.Ptr
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

func (u *utilsPackage) GetScriptMigrate(db *sql.DB, dbName string, dialect Dialect, typ ...reflect.Type) ([]string, error) {
	dbSchema, err := utils.extractSchema(db, dbName, dialect)
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for _, t := range typ {
		sqlCmds, err := dialect.GenerateCreateTableSql(dbName, t)
		if err != nil {
			return nil, err
		}
		if sqlCmds != "" {
			ret = append(ret, sqlCmds)
		} else { //<--- table is existing in Database, just add columns
			sqlAddCols, err := dialect.GenerateAlterTableSql(dbName, t)
			if err != nil {
				return nil, err
			}
			ret = append(ret, sqlAddCols...)
		}
		sqlUniqueConstraints := dialect.GenerateUniqueConstraintsSql(t)

		if err != nil {
			return nil, err
		}

		for constraintName, sql := range sqlUniqueConstraints {
			if !dbSchema.unique[constraintName] {
				ret = append(ret, sql)
			}
		}
		sqlIndexConstraints := dialect.GenerateIndexConstraintsSql(t)
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
func (u *utilsPackage) GetDbName(db *sql.DB) (string, error) {
	dbType, _, err := u.DetectDatabaseType(db)
	if err != nil {
		return "", err
	}
	return u.GetCurrentDatabaseName(db, dbType)
}

func (u *utilsPackage) getTenantDb(db *sql.DB, typ reflect.Type) (*TenantDb, error) {
	baseType := reflect.TypeOf(TenantDb{})

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if baseType.String() == field.Type.String() {
			_dbSchema, err := u.newTenantDb(db)
			if err != nil {
				return nil, err
			}
			return _dbSchema, nil

		} else {
			_dbSchema, err := u.getTenantDb(db, field.Type)
			if err != nil {
				return nil, err
			} else if _dbSchema != nil {
				return _dbSchema, nil
			} else {
				continue
			}
		}
	}
	return nil, nil
}
func (u *utilsPackage) newTenantDb(db *sql.DB) (*TenantDb, error) {
	ret := &TenantDb{}
	ret.DB = db
	dbDetect, dbTypeName, err := utils.DetectDatabaseType(db)

	if err != nil {
		return nil, err
	}
	dbName, err := utils.GetCurrentDatabaseName(db, dbDetect)
	if err != nil {
		return nil, err
	}
	ret.DbName = dbName
	if dbDetect == DBMSSQL {
		ret.Dialect = NewSqlServerDialect(db)
		ret.DBType = DBMSSQL
		ret.DBTypeName = dbTypeName
	} else if dbDetect == DBMySQL {
		ret.Dialect = NewSqlServerDialect(db)
		ret.DBType = DBMySQL
		ret.DBTypeName = dbTypeName
	} else if dbDetect == DBPostgres {
		ret.DBType = DBPostgres
		ret.DBTypeName = dbTypeName
	} else {
		return nil, fmt.Errorf("Unsupported database type '%s'", dbTypeName)

	}
	return ret, nil
}

/*
This struct is only used for the function buildRepositoryFromType of the "utilsPackage"
*/
type repositoryValueStruct struct {
	/*
		The buildRepositoryFromType function of utilsPackage
		will analyze the type of repo: During the analysis process,
		it will use reflect.New(type) to create a value for this field
	*/
	ValueOfRepo reflect.Value
	/*
		During the analysis of the entity type, these are fields that have struct types,
		and those structs have fields declared with types like DbField[<type>], including in embedded struct
	*/
	EntityTypes []reflect.Type
}

/*
This function will read information from @typ and create a structure similar to the one described in
"repositoryValueStruct" if no error occurs
*/
func (u *utilsPackage) buildRepositoryFromType(typ reflect.Type) (*repositoryValueStruct, error) {
	valueOfRepo := reflect.New(typ).Elem()

	baseType := reflect.TypeOf(TenantDb{})
	entityTypes := []reflect.Type{}
	for i := 0; i < typ.NumField(); i++ {

		field := typ.Field(i)
		if field.Anonymous {
			fmt.Println(field.Type.Name())
			if baseType.Name() == field.Type.Name() {

				continue

			} else {
				repoVal, err := u.buildRepositoryFromType(field.Type) //<-- do not gen sql migrate for inner entity
				if err != nil {
					return nil, err
				}
				entityTypes = append(entityTypes, repoVal.EntityTypes...)
				valueOfRepo.Field(i).Set(repoVal.ValueOfRepo.Addr())
				entityType := field.Type
				if entityType.Kind() == reflect.Ptr {
					entityType = entityType.Elem()
				}

				continue
			}
		}
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		entityVal := u.EntityFromType(fieldType)

		valueOfRepo.Field(i).Set(entityVal.Addr())
		entityType := field.Type
		if entityType.Kind() == reflect.Ptr {
			entityType = entityType.Elem()
		}
		entityTypes = append(entityTypes, entityType)
	}
	ret := &repositoryValueStruct{
		ValueOfRepo: valueOfRepo,
		EntityTypes: entityTypes,
	}
	return ret, nil
}

func (u *utilsPackage) GetOrCreateRepository(typ reflect.Type) (*repositoryValueStruct, error) {
	//check cache
	key := typ.String()
	if val, ok := u.cacheGetOrCreateRepository.Load(key); ok {
		return val.(*repositoryValueStruct), nil
	}
	repoVal, err := u.buildRepositoryFromType(typ)
	if err != nil {
		return nil, err
	}
	u.cacheGetOrCreateRepository.Store(key, repoVal)
	return repoVal, nil
}
func (u *utilsPackage) EntityFromType(typ reflect.Type) reflect.Value {
	val := reflect.New(typ).Elem()

	for i := 0; i < typ.NumField(); i++ {
		valField := val.Field(i)
		typField := typ.Field(i).Type

		if valField.Kind() == reflect.Ptr {
			typField = typField.Elem()
			valField = reflect.New(typField)
			val.Field(i).Set(valField)
			valField = valField.Elem()

		}
		// Locate and set the "TableName" field inside each struct field
		tableNameField := valField.FieldByName("TableName")
		if tableNameField.IsValid() && tableNameField.CanSet() {
			tableNameField.SetString(utils.TableNameFromStruct(typ))
		}

		// Locate and set the "ColName" field inside each struct field
		columnNameField := valField.FieldByName("ColName")
		if columnNameField.IsValid() && columnNameField.CanSet() {
			columnNameField.SetString(utils.ToSnakeCase(typ.Field(i).Name))
		}
	}
	return val
}

func (u *utilsPackage) exprToSQLDelete(v interface{}, d Dialect) (string, []interface{}) {

	switch val := v.(type) {
	case *Field[any]:
		return val.ToSqlExpr(d)
	case *Field[string]:
		return val.ToSqlExpr(d)
	case *Field[int]:
		return val.ToSqlExpr(d)
	case *Field[float64]:
		return val.ToSqlExpr(d)
	case *Field[float32]:
		return val.ToSqlExpr(d)
	case *Field[uint64]:
		return val.ToSqlExpr(d)
	case *Field[bool]:
		return val.ToSqlExpr(d)
	default:
		return "?", []interface{}{val}
	}
}

func (u *utilsPackage) join(parts []string, sep string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += sep
		}
		out += p
	}
	return out
}
func (u *utilsPackage) PtrToInterface(v interface{}) interface{} {
	val := reflect.ValueOf(v)
	ptr := reflect.New(reflect.TypeOf((*interface{})(nil)).Elem()) // tạo interface{}
	ptr.Elem().Set(val)                                            // gán giá trị
	return ptr.Interface()                                         // trả *interface{}
}

// Entity creates an instance of struct T and auto-populates any fields
// named "TableName" and "ColName" with the corresponding struct and field names.

// This is useful for initializing DbField[TTable, TField] fields in a model,
// so that the table and column names can be inferred via reflection without manual assignment.

// Requirements:
// - Each field inside struct T must be a struct that contains fields named "TableName" and "ColName".
// - Those inner fields must be settable (exported and addressable).
func Entity[T any]() T {
	var v T

	// Get the type name of T to use as table name
	typ := reflect.TypeOf(v)

	ret := utils.EntityFromType(typ)
	return ret.Interface().(T)
}
