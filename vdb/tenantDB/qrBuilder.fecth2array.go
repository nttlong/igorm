package tenantDB

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type initMakeScanInfo struct {
	once sync.Once
	info [][]int
}

var scanInfoCache sync.Map

func makeScanInfo(t reflect.Type, cols []string) [][]int {
	key := t.String() + strings.Join(cols, ",")
	actual, _ := scanInfoCache.LoadOrStore(key, &initMakeScanInfo{})
	initInfo := actual.(*initMakeScanInfo)
	initInfo.once.Do(func() {
		initInfo.info = makeScanInfoNoCache(t, cols)
	})
	return initInfo.info
}
func makeScanInfoNoCache(t reflect.Type, cols []string) [][]int {
	ret := make([][]int, len(cols))
	for i, col := range cols {
		field, ok := t.FieldByNameFunc(func(s string) bool {

			return strings.EqualFold(s, col)
		})
		if !ok {
			continue
		}
		ret[i] = field.Index
	}
	return ret
}
func doScan(rows *sql.Rows, dest interface{}) error {
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer to a struct or slice")
	}
	if destVal.Elem().Kind() == reflect.Slice {
		sliceVal := destVal.Elem()
		elemType := sliceVal.Type().Elem()
		if elemType.Kind() != reflect.Struct {
			return fmt.Errorf("only slice of struct is supported")
		}

		cols, err := rows.Columns()
		if err != nil {
			return err
		}

		scanInfo := makeScanInfo(elemType, cols)
		for rows.Next() {
			scanArgs := make([]interface{}, len(cols))
			elemVal := reflect.New(elemType)
			for i, index := range scanInfo {
				if len(index) > 0 {
					fieldVal := elemVal.Elem().FieldByIndex(index)
					if fieldVal.CanSet() {
						scanArgs[i] = fieldVal.Addr().Interface()
					}
				} else {
					var dummy interface{}
					scanArgs[i] = &dummy // For columns that do not match any field
				}
			}
			if err := rows.Scan(scanArgs...); err != nil {
				return err
			}
			sliceVal.Set(reflect.Append(sliceVal, elemVal.Elem()))

		}

	}
	return nil

}
func (q *query) ToArray(items interface{}) error {
	sql, args := q.BuildSql()
	rows, err := q.db.Query(sql, args...)
	if err != nil {
		return err
	}

	//err = ToArrayUnsafeFast(rows, items)
	err = doScan(rows, items)

	if err != nil {
		return err
	}
	return nil

}
