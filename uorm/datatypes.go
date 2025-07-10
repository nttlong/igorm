package uorm

import (
	"fmt"
	"reflect"
	"sync"
)

type DB_TYPE int

const (
	DB_TYPE_MYSQL DB_TYPE = iota
	DB_TYPE_POSTGRES
	DB_TYPE_SQLITE
	DB_TYPE_MSSQL
)

func (t DB_TYPE) String() string {
	switch t {
	case DB_TYPE_MYSQL:
		return "mysql"
	case DB_TYPE_POSTGRES:
		return "postgres"
	case DB_TYPE_SQLITE:
		return "sqlite"
	case DB_TYPE_MSSQL:
		return "mssql"
	default:
		return "unknown"
	}
}

type quote struct {
	left  string
	right string
}

func (q *quote) Quote(s string) string {
	return q.left + s + q.right
}

type iField interface {
	getKey() string
}
type Field struct {
	iField
	expr     string
	table    *Table
	name     string
	args     []interface{}
	cacheKey string
}
type Table struct {
	name string
}
type Model struct {
	table  *Table
	entity reflect.Type
	dbType DB_TYPE
}

func (m *Model) As(alias string) interface{} {
	ret := utils.createQueryableFormType(m.entity, alias, m.dbType, true)

	return ret
}

func (f *Field) String() string {
	if f.table != nil {
		return f.table.name + "." + f.expr
	}
	return f.expr
}

var cacheMmakeBinaryExpr sync.Map

func makeBinaryExpr(left interface{}, right interface{}, op string) *Field {
	key := op

	var leftStr, rightStr string
	var args []interface{}

	// Xử lý left
	switch l := left.(type) {
	case *Field:
		if l.table != nil {
			key += ":" + l.table.name + ":" + l.expr
			if _leftStr, ok := cacheMmakeBinaryExpr.Load(key); ok {
				leftStr = _leftStr.(string)
			} else {
				leftStr = fmt.Sprintf("%s.%s", l.table.name, l.expr)
				cacheMmakeBinaryExpr.Store(key, leftStr)
			}
		} else {
			key += ":?"
			leftStr = l.expr
		}
		args = append(args, l.args...)
	case Field:
		if l.table != nil {
			key += ":" + l.table.name + l.expr
			if _leftStr, ok := cacheMmakeBinaryExpr.Load(key); ok {
				leftStr = _leftStr.(string)
			} else {
				leftStr = fmt.Sprintf("%s.%s", l.table.name, l.expr)
				cacheMmakeBinaryExpr.Store(key, leftStr)
			}

		} else {
			key += ":?"
			leftStr = l.expr
		}
		args = append(args, l.args...)
	default:
		leftStr = "?"
		key += ":?"
		args = append(args, l)
	}

	// Xử lý right
	switch r := right.(type) {
	case *Field:
		if r.table != nil {
			key += ":" + r.table.name + r.expr

			if _rightStr, ok := cacheMmakeBinaryExpr.Load(key); ok {
				rightStr = _rightStr.(string)
			} else {
				rightStr = fmt.Sprintf("%s.%s", r.table.name, r.expr)
				cacheMmakeBinaryExpr.Store(key, rightStr)
			}

		} else {
			key += ":?"
			rightStr = r.expr
		}
		args = append(args, r.args...)
	case Field:
		if r.table != nil {
			key += ":" + r.table.name + r.expr

			if _rightStr, ok := cacheMmakeBinaryExpr.Load(key); ok {
				rightStr = _rightStr.(string)
			} else {
				rightStr = fmt.Sprintf("%s.%s", r.table.name, r.expr)
				cacheMmakeBinaryExpr.Store(key, rightStr)
			}

		} else {
			key += ":?"
			rightStr = r.expr
		}
		args = append(args, r.args...)
	default:
		rightStr = "?"
		key += ":?"
		args = append(args, r)
	}
	if _strSyntax, ok := cacheMmakeBinaryExpr.Load(key); ok {
		return &Field{
			expr: _strSyntax.(string),
			args: args,
		}
	}
	strSyntax := fmt.Sprintf("%s %s %s", leftStr, op, rightStr)
	cacheMmakeBinaryExpr.Store(key, strSyntax)
	return &Field{
		expr: strSyntax,
		args: args,
	}
}

func Queryable[T any](dbType DB_TYPE, tableName string) T {

	ret := utils.createQueryableFormType(reflect.TypeFor[T](), tableName, dbType, false)
	return ret.(T)
}
