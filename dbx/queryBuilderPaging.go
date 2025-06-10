package dbx

import (
	"context"
	"reflect"
	"strconv"
	"strings"
)

//	type PagerItem[T any] struct {
//		Data      T
//		Index     int
//		PageIndex int
//		HasNext   bool
//		NextPage  int
//	}
type QrPager[T any] struct {
	dbx          *DBXTenant
	selectFields []string
	selector     string
	where        string
	orders       []string
	sorting      string
	pagesize     int //num of row per page
	pageIndex    int
	from         string
	args         []interface{}
	ctx          context.Context
	entityType   *EntityType
}

// func getUnderEntityType(t reflect.Type){

// }
func Pager[T any](dbx *DBXTenant, ctx context.Context) *QrPager[T] {
	entityType, err := newEntityType(reflect.TypeFor[T]())
	if err != nil {
		panic(err)
	}
	selectFields := []string{}
	// for _, field := range entityType.EntityFields {
	// 	if strings.ToLower(field.Name) != "index" {
	// 		selectFields = append(selectFields, "`"+field.Name+"`")
	// 	}

	// }
	// selectFields = append(selectFields, "row_number() `Index`")
	return &QrPager[T]{
		dbx:          dbx,
		selector:     "*",
		selectFields: selectFields,
		entityType:   entityType,
		from:         entityType.TableName,
		ctx:          ctx,
	}

}
func (q *QrPager[T]) Where(where string, args ...interface{}) *QrPager[T] {
	q.where = where
	q.args = append(q.args, args...)
	return q
}
func (q *QrPager[T]) Sort(OrderBy string) *QrPager[T] {
	q.sorting = OrderBy
	return q
}
func (q *QrPager[T]) OrderBy(fieldName string, isAsc bool) *QrPager[T] {
	if isAsc {
		q.orders = append(q.orders, fieldName)
	} else {
		q.orders = append(q.orders, fieldName+" DESC")
	}
	return q

}
func (q *QrPager[T]) Size(size int) *QrPager[T] {
	q.pagesize = size
	return q
}
func (q *QrPager[T]) Page(pageIndex int) *QrPager[T] {
	q.pageIndex = pageIndex
	return q
}
func (q *QrPager[T]) toSQL() (string, error) {
	if len(q.selectFields) > 0 {
		q.selector = strings.Join(q.selectFields, ",")
	}
	sqlRet := "SELECT " + q.selector + " FROM " + q.from

	fromIndex := strconv.Itoa(q.pageIndex * q.pagesize)
	toIndex := strconv.Itoa(q.pagesize)

	if q.where != "" {
		sqlRet += " WHERE " + q.where
		// q.where += " AND `index`>=" + fromIndex + " AND `index`<" + toIndex
	} else {
		// q.where += "`index`>=" + fromIndex + " AND `index`<" + toIndex
	}

	if len(q.orders) > 0 {
		if q.sorting != "" {
			q.sorting += ","
		}
		q.sorting += strings.Join(q.orders, ",")

	}
	if q.sorting != "" {
		sqlRet += " ORDER BY " + q.sorting
	}
	sqlRet += " LIMIT " + fromIndex + "," + toIndex
	// if q.limit == 0 {
	// 	q.limit = 50
	// }
	// if q.limit > 0 {
	// 	sqlRet += " LIMIT " + strconv.Itoa(q.limit)
	// }
	// if q.offset > 0 {
	// 	sqlRet += " OFFSET " + strconv.Itoa(q.offset)
	// }
	// fromIndex := q.offset
	execSQl, err := q.dbx.compiler.Parse(sqlRet, q.args...)
	if err != nil {
		return "", err
	}
	return execSQl, nil

}
func (q *QrPager[T]) Query() ([]T, error) {
	execSQl, err := q.toSQL()
	if err != nil {
		return nil, err
	}
	rows, err := q.dbx.QueryContext(q.ctx, execSQl, q.args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//ret := PagerItem[T]{}

	r, e := fetchAllRows(rows, reflect.TypeFor[T]())
	if e != nil {
		return nil, e
	}
	return r.([]T), nil
	// for rows.Next() {
	// 	item := PagerItem[T]{}

	// 	err = rows.Scan(
	// 		&item.Data,
	// 		&item.Index,
	// 		&item.PageIndex,
	// 		&item.HasNext,
	// 		&item.NextPage,
	// 	)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	ret = append(ret, item)
	// }
	// return ret, nil
}
