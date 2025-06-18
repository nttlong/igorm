package dbx

type DBXErrorCode int

const (
	DBXErrorCodeUnknown DBXErrorCode = iota
	DBXErrorCodeDuplicate
	DBXErrorCodeInvalidSize
	DBXErrorCodeReferenceConstraint
	DBXErrorCodeMissingRequiredField
)

func (e DBXErrorCode) String() string {
	switch e {
	case DBXErrorCodeUnknown:
		return "UNKNOWN"
	case DBXErrorCodeDuplicate:
		return "DUPLICATE"
	case DBXErrorCodeInvalidSize:
		return "INVALID_SIZE"
	case DBXErrorCodeReferenceConstraint:
		return "REFERENCE_CONSTRAINT"
	case DBXErrorCodeMissingRequiredField:
		return "MISSING_REQUIRED_FIELD"
	default:
		return "UNKNOWN"
	}
}

type DBXError struct {
	// Error code the value is one of DBXErrorCode
	Code DBXErrorCode `json:"code"`
	// Error message
	Message string `json:"message"`
	// table name
	TableName string `json:"tableName"`
	//constraint name
	ConstraintName string `json:"constraintName"`
	// list of column names caused the error
	Fields []string `json:"fields"`
	// values of columns caused the error
	Values  []string `json:"values"`
	MaxSize int      `json:"maxSize"`
}
type DBXMigrationError struct {
	Message   string `json:"message"`
	Err       error  `json:"error"`
	DBName    string `json:"dbName"`
	TableName string `json:"tableName"`
	Code      string `json:"code"`
	Sql       string `json:"sql"`
}

func (e *DBXError) Error() string {
	return e.Message
}
func (e DBXMigrationError) Error() string {
	return e.Message
}
