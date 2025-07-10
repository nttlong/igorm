package eorm

import (
	"strings"
	"sync"
)

type cacheSqlExecutorBuildSourceItem struct {
	tables []string
	alias  map[string]string
	source string
}

var cacheSqlExecutorBuildSource sync.Map

type SqlBuilderReceiver struct {
}
type SqlExecutor struct {
	source    string
	selectors string
	tables    []string
	alias     map[string]string
	args      []interface{}
	Err       error
}

func (s *SqlBuilderReceiver) From(expr string, args ...interface{}) *SqlExecutor {
	return &SqlExecutor{
		source: expr,
		tables: []string{},
		alias:  map[string]string{},
		args:   args,
	}

	// }
	// err := ej.build("Departments INNER JOIN User ON User.Code = Departments.Code INNER JOIN Check ON Check.Name = 'John'")
	// assert.NoError(t, err)
	// assert.Equal(t, "[departments] AS [T1] INNER JOIN [User] AS [T2] ON [T2].[Code] = [T1].[code] INNER JOIN [checks] AS [T3] ON [T3].[name] = N'John'", ej.content)
}
func (s *SqlExecutor) Select(args ...interface{}) *SqlExecutor {
	if len(args) == 0 {
		return s
	} else if len(args) == 1 {
		if str, ok := args[0].(string); ok {
			s.selectors = str
		}

	} else {
		if str, ok := args[0].(string); ok {
			s.selectors = str
		}
		s.args = args[1:]
	}
	return s
}

func (s *SqlExecutor) buildSource(dialect Dialect) error {
	key := dialect.Name() + "://" + s.source
	if v, ok := cacheSqlExecutorBuildSource.Load(key); ok {
		item := v.(cacheSqlExecutorBuildSourceItem)
		s.tables = item.tables
		s.alias = item.alias
		s.source = item.source
		return nil
	}
	joinCompilerSource := &exprJoin{
		context: &exprCompileContext{
			tables:      s.tables,
			alias:       s.alias,
			dialect:     dialect,
			IsBuildJoin: false,
		},
	}
	err := joinCompilerSource.build(s.source)
	if err != nil {
		return err
	}
	s.source = joinCompilerSource.content

	item := cacheSqlExecutorBuildSourceItem{
		tables: []string{},
		alias:  map[string]string{},

		source: s.source,
	}
	for _, table := range joinCompilerSource.context.alias {
		item.tables = append(item.tables, table)
	}
	for k, v := range joinCompilerSource.context.alias {
		item.alias[k] = v
	}
	cacheSqlExecutorBuildSource.Store(key, item)

	s.source = joinCompilerSource.content
	s.alias = item.alias
	s.tables = item.tables
	return nil
}
func (s *SqlExecutor) buildSelectors(dialect Dialect) error {
	if s.selectors == "" {
		if len(s.tables) == 1 {
			s.selectors = "*"
		} else {
			selectAll := []string{}
			for _, table := range s.tables {
				selectAll = append(selectAll, dialect.Quote(table)+".*")
			}
			s.selectors = strings.Join(selectAll, ", ")
		}

	}
	return nil

}
func (s *SqlExecutor) ToSql(dialect Dialect) (string, []interface{}) {
	err := s.buildSource(dialect)
	if err != nil {
		s.Err = err
		return "", nil
	}
	err = s.buildSelectors(dialect)
	if err != nil {
		s.Err = err
		return "", nil
	}
	sql := "SELECT " + s.selectors + " FROM " + s.source
	return sql, s.args

}

var SqlBuilder = &SqlBuilderReceiver{}
