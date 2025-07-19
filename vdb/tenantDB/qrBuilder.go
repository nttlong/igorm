package tenantDB

import (
	"reflect"
)

type query struct {
	db         *TenantDB
	qrInstance interface{}
	sql        string
	args       []interface{}
}
type OnNewQr func() interface{}

var OnNewQrFn OnNewQr
var OnFrom OnExpr

type onLiterals = func(string) interface{}

var OnLiterals onLiterals

func (db *TenantDB) Lit(str string) interface{} {
	return OnLiterals(str)
}
func (db *TenantDB) From(table string) *query {
	ret := createQr(db)
	ret.qrInstance = OnFrom(ret.qrInstance, table)

	return ret
}
func createQr(db *TenantDB) *query {
	ret := &query{
		db: db,
	}
	ret.qrInstance = OnNewQrFn()
	return ret
}

type selectExpr struct {
	Expr string
	Args []interface{}
}
type litOfString struct {
	val string
}

/*
Example:

	Select("id", "name", "age") <-- no args

	Select("id", "name", "age", 123) <-- with args

	In case arg is string select can not recognize it as parameter and will be treated as string literal
	Select("concat(firstName, ?,lastName)",Lit(" ")) <-- with args
*/
type onQrSelect func(qrInstance interface{}, exprsAndArgs ...interface{}) interface{}

var OnQrSelect onQrSelect

func (q *query) Select(exprsAndArgs ...interface{}) *query {

	q.qrInstance = OnQrSelect(q.qrInstance, exprsAndArgs...)
	return q
}

type onJoin = func(qrInstance interface{}, table string, onExpr string, args ...interface{}) interface{}

var OnInnerJoin onJoin
var OnLeftJoin onJoin
var OnRightJoin onJoin
var OnFullJoin onJoin

func (q *query) InnerJoin(table string, onExpr string, args ...interface{}) *query {
	q.qrInstance = OnInnerJoin(q.qrInstance, table, onExpr, args...)
	return q
}
func (q *query) LeftJoin(table string, onExpr string, args ...interface{}) *query {
	q.qrInstance = OnLeftJoin(q.qrInstance, table, onExpr, args...)
	return q
}

type OnExpr = func(qrInstance interface{}, expr string, args ...interface{}) interface{}

var OnWhere OnExpr

func (q *query) Where(expr string, args ...interface{}) *query {
	q.qrInstance = OnWhere(q.qrInstance, expr, args...)
	return q
}

type onQueryArgs = func(qrInstance interface{}, args ...interface{}) interface{}

var OnOrderBy onQueryArgs

func (q *query) OrderBy(args ...interface{}) *query {
	q.qrInstance = OnOrderBy(q.qrInstance, args...)
	return q
}

var OnGroupBy onQueryArgs

func (q *query) GroupBy(args ...interface{}) *query {
	q.qrInstance = OnGroupBy(q.qrInstance, args...)
	return q
}

type onOffsetLimit = func(qrInstance interface{}, offset, limit int) interface{}

var OnOffsetLimit onOffsetLimit

func (q *query) OffsetLimit(offset, limit int) *query {
	q.qrInstance = OnOffsetLimit(q.qrInstance, offset, limit)
	return q
}

var OnHaving OnExpr

func (q *query) Having(expr string, args ...interface{}) *query {
	q.qrInstance = OnHaving(q.qrInstance, expr, args...)
	return q
}

type onBuildSql = func(qrInstance interface{}, db *TenantDB) (string, []interface{})

var OnBuildSql onBuildSql

func (q *query) BuildSql() (string, []interface{}) {
	if q.sql != "" {
		return q.sql, q.args
	}
	sql, args := OnBuildSql(q.qrInstance, q.db)
	q.sql = sql
	q.args = args
	return sql, args
}
func (q *query) ToArray(items interface{}) error {
	sql, args := q.BuildSql()

	err := q.db.ExecToArray(items, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

type onBuildSQL func(typ reflect.Type, db *TenantDB, filter string) (string, error)

var OnBuildSQL onBuildSQL

/*
Example:

	r.db.Find(&users, "email = ? OR username = ?", identifier, identifier)
*/
func (db *TenantDB) Find(entity interface{}, filter string, args ...interface{}) error {
	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	sql, err := OnBuildSQL(typ, db, filter)
	if err != nil {
		return err
	}
	return db.ExecToArray(entity, sql, args...)

}

type onBuildSQLFirstItem func(typ reflect.Type, db *TenantDB, filter string) (string, error)

var OnBuildSQLFirstItem onBuildSQLFirstItem

func (db *TenantDB) First(entity interface{}, filter string, args ...interface{}) error {
	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	sql, err := OnBuildSQLFirstItem(typ, db, filter)
	if err != nil {
		return err
	}
	return db.ExecToItem(entity, sql, args...)

}

type whereExpr struct {
	expr string
	args []interface{}
	db   *TenantDB
}

func (db *TenantDB) Where(expr string, args ...interface{}) *whereExpr {
	return &whereExpr{
		db:   db,
		expr: expr,
		args: args,
	}
}
func (we *whereExpr) And(expr string, args ...interface{}) *whereExpr {
	we.expr += " AND " + expr
	we.args = append(we.args, args...)
	return we
}
func (we *whereExpr) Or(expr string, args ...interface{}) *whereExpr {
	we.expr += " OR " + expr
	we.args = append(we.args, args...)
	return we
}
func (we *whereExpr) First(entity interface{}) error {
	return we.db.First(entity, we.expr, we.args...)
}

type onCreateEntity = func(db *TenantDB, entity interface{}) error

var OnCreateEntity onCreateEntity

func (db *TenantDB) Create(entity interface{}) error {
	return OnCreateEntity(db, entity)
}
