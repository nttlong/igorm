package expr

type DB_TYPE int

const (
	DB_TYPE_UNKNOWN DB_TYPE = iota
	DB_TYPE_MYSQL
	DB_TYPE_POSTGRES
	DB_TYPE_MSSQL
)

func (t DB_TYPE) FromString(str string) DB_TYPE {
	switch str {
	case "mysql":
		return DB_TYPE_MYSQL
	case "postgres":
		return DB_TYPE_POSTGRES
	case "mssql":
		return DB_TYPE_MSSQL
	default:
		return -1
	}
}
func (t DB_TYPE) String() string {
	switch t {
	case DB_TYPE_MYSQL:
		return "mysql"
	case DB_TYPE_POSTGRES:
		return "postgres"
	case DB_TYPE_MSSQL:
		return "mssql"
	default:
		return "unknown"
	}
}

type expression struct {
	keywords    []string
	specialChar []byte
	DbDriver    DB_TYPE
}
type ResolverResult struct {
	Syntax      string
	Args        []interface{}
	AliasSource *map[string]string
}

func (e *expression) Quote(str ...string) string {
	return OnGetQuoteFunc(e.DbDriver, str...)
}
func (e *expression) resolve(tables *[]string, context *map[string]string, caller interface{}, requireAlias bool) (*ResolverResult, error) {
	return OnCompileFunc(e.DbDriver, tables, context, caller, requireAlias)
	//return nil, nil
}

type OnCompile = func(dbDriver DB_TYPE, tables *[]string, context *map[string]string, caller interface{}, requireAlias bool) (*ResolverResult, error)
type ExpressionTest = expression
type OnGetQuote = func(dbDriver DB_TYPE, str ...string) string

var OnGetQuoteFunc OnGetQuote
var OnCompileFunc OnCompile
