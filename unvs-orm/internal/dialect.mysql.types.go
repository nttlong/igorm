package internal

import "database/sql"

type MySqlDialect struct {
	baseDialect
	db               *sql.DB
	paramPlaceholder string
}

func NewMysqlDialect() Dialect {
	return &MySqlDialect{
		paramPlaceholder: "$",
		baseDialect: baseDialect{
			schema: map[string]map[string]TableSchema{},
			mapGoTypeToDb: map[string]string{
				"string":    "CITEXT",
				"int":       "INTEGER",
				"int32":     "INTEGER",
				"int64":     "BIGINT",
				"uint":      "BIGINT",
				"uint32":    "BIGINT",
				"uint64":    "BIGINT",
				"int16":     "SMALLINT",
				"int8":      "SMALLINT",
				"uint8":     "SMALLINT",
				"bool":      "BOOLEAN",
				"float32":   "REAL",
				"float64":   "DOUBLE PRECISION",
				"time.Time": "TIMESTAMP",
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
				"bool":      "FALSE",
				"float32":   "0",
				"float64":   "0",
				"time.Time": "CURRENT_TIMESTAMP",
				"true":      "TRUE",
				"false":     "FALSE",
				"now()":     "CURRENT_TIMESTAMP",
			},
		},
	}
}
