package internal

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

// // --------------------- Tag Metadata ---------------------
// // FieldTag holds parsed metadata from struct field tags.
// type FieldTag struct {
// 	PrimaryKey    bool
// 	AutoIncrement bool
// 	Unique        bool
// 	/*
// 		can be field name if no unique index name in tag else name of unique index in tag
// 	*/
// 	UniqueName string
// 	Index      bool
// 	/*
// 		can be field name if no  index name in tag else name of  index in tag
// 	*/
// 	IndexName string
// 	Length    *int
// 	FTSName   string
// 	DBType    string
// 	TableName string
// 	Check     string
// 	Nullable  bool
// 	Field     reflect.StructField
// 	Default   string
// }
