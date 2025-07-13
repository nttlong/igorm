package eorm

import (
	"sort"
	"strings"
	"sync"
)

type cacheSqlExecutorBuildSourceItem struct {
	tables []string
	alias  map[string]string
	source string
}

type SqlBuilderReceiver struct {
}
type SqlExecutor struct {
	source          string
	selectors       string
	tables          []string
	alias           map[string]string
	args            []interface{}
	Err             error
	tableInDataBase *[]string
	dbName          string
}

/*
While compiling the code keep original table name and field name in database.
if table found in tableInDataBase
*/
func (s *SqlExecutor) SetTableInDataBase(databaseName string, tableInDataBase *[]string) {
	s.tableInDataBase = tableInDataBase
	s.dbName = databaseName
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

func (s *SqlBuilderReceiver) Select(args ...interface{}) *SqlExecutor {

	ret := &SqlExecutor{}
	ret.Select(args...)
	return ret

}
func (s *SqlExecutor) Select(args ...interface{}) *SqlExecutor {
	strSelectors := []string{}
	argsSelectors := []interface{}{}
	for _, arg := range args {
		if str, ok := arg.(string); ok {
			strSelectors = append(strSelectors, str)
		} else if p, ok := arg.(*builderParam); ok {
			argsSelectors = append(argsSelectors, p.Params...)
		}
	}
	s.selectors = strings.Join(strSelectors, ", ")
	s.args = append(argsSelectors, s.args...)
	return s
}

type builderParam struct {
	Params []interface{}
}

func Paras(args ...interface{}) *builderParam {
	return &builderParam{Params: args}
}
func (s *SqlExecutor) From(args ...interface{}) *SqlExecutor {
	strSource := []string{}
	argsOfBuilder := []interface{}{}
	for _, arg := range args {
		if str, ok := arg.(string); ok {
			strSource = append(strSource, str)
		} else if p, ok := arg.(*builderParam); ok {
			argsOfBuilder = append(argsOfBuilder, p.Params...)
		}
	}

	s.source = strings.Join(strSource, " ")
	s.args = append(s.args, argsOfBuilder...)
	return s
}
func (s *SqlExecutor) buildSource(exprCompile *exprCompiler, dialect Dialect) error {

	err := exprCompile.build(s.source)
	if err != nil {
		return err
	}
	s.source = exprCompile.content

	item := cacheSqlExecutorBuildSourceItem{
		tables: []string{},
		alias:  map[string]string{},

		source: s.source,
	}
	for _, table := range exprCompile.context.alias {
		item.tables = append(item.tables, table)
	}
	for k, v := range exprCompile.context.alias {
		item.alias[k] = v
	}

	s.source = exprCompile.content
	s.alias = item.alias
	s.tables = item.tables
	return nil
}
func (s *SqlExecutor) buildSelectors(compiler *exprCompiler, dialect Dialect) error {
	if s.selectors == "" {
		if len(s.tables) == 1 {
			s.selectors = "*"
		} else {
			selectAll := []string{}
			sort.Strings(s.tables)
			for _, table := range s.tables {
				selectAll = append(selectAll, dialect.Quote(table)+".*")
			}
			s.selectors = strings.Join(selectAll, ", ")
			return nil
		}

	}

	err := compiler.buildSelectField(s.selectors)
	if err != nil {
		return err
	}
	s.selectors = compiler.content

	return nil

}

type toSqlInit struct {
	once sync.Once
	val  string
}

var cacheSqlExecutorToSql sync.Map

func (s *SqlExecutor) ToSql(dialect Dialect) (string, []interface{}) {
	key := dialect.Name() + "://" + s.source + "?" + s.selectors
	if s.tableInDataBase != nil {
		key += "//" + strings.Join(*s.tableInDataBase, ",")
	}
	if s.dbName != "" {
		key += "@" + s.dbName
	}

	actual, _ := cacheSqlExecutorToSql.LoadOrStore(key, &toSqlInit{})
	v := actual.(*toSqlInit)
	var err error
	var sql string
	v.once.Do(func() {
		sql, err = s.toSql(dialect)
		v.val = sql

	})
	if err != nil {
		s.Err = err
		return "", s.args
	}
	return v.val, s.args

}
func (s *SqlExecutor) toSql(dialect Dialect) (string, error) {

	dbSchema := map[string]bool{}
	if s.tableInDataBase != nil {
		for _, table := range *s.tableInDataBase {
			dbSchema[table] = true
		}
	}
	exprCompile := &exprCompiler{
		context: &exprCompileContext{
			tables:  s.tables,
			alias:   s.alias,
			dialect: dialect,
			purpose: build_purpose_join,
			schema:  &dbSchema,
		},
	}
	err := s.buildSource(exprCompile, dialect)
	if err != nil {

		return "", err
	}
	err = s.buildSelectors(exprCompile, dialect)
	if err != nil {

		return "", err
	}
	sql := "SELECT " + s.selectors + " FROM " + s.source

	return sql, nil

}

var SqlBuilder = &SqlBuilderReceiver{}
