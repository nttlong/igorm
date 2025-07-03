package internal

import "database/sql"

type SqlServerDialect struct {
	baseDialect
	db               sql.DB
	paramPlaceholder string
}

func NewMssqlDialect() Dialect {
	return &SqlServerDialect{
		paramPlaceholder: "@p",
		baseDialect: baseDialect{
			schema: map[string]map[string]TableSchema{},
			mapGoTypeToDb: map[string]string{
				"string":    "NVARCHAR(MAX)",
				"int":       "INT",
				"int32":     "INT",
				"int64":     "BIGINT",
				"uint":      "BIGINT",
				"uint32":    "BIGINT",
				"uint64":    "BIGINT",
				"int16":     "SMALLINT",
				"int8":      "TINYINT",
				"uint8":     "TINYINT",
				"bool":      "BIT",
				"float32":   "REAL",
				"float64":   "FLOAT",
				"time.Time": "DATETIME2",
			},
			mapDefaultValue: map[string]string{
				"string":    "''",
				"int":       "0",
				"int32":     "0",
				"int64":     "0",
				"uint":      "0",
				"uint32":    "0",
				"uint64":    "0",
				"int16":     "0",
				"int8":      "0",
				"uint8":     "0",
				"bool":      "0",
				"float32":   "0",
				"float64":   "0",
				"time.Time": "now()",
				"true":      "1",
				"false":     "0",
				"now()":     "GETDATE()",
			},
		},
	}
}
