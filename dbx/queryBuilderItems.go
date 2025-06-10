package dbx

import "reflect"

func (q *QrBuilder[T]) Items() ([]T, error) {
	var zero T
	et := reflect.TypeOf(zero)
	entityType, err := newEntityType(et)
	if err != nil {
		return nil, err
	}
	sqlSelect := ""
	if q.where == "" {
		sqlSelect = "SELECT * FROM " + entityType.TableName
	} else {
		sqlSelect = "SELECT * FROM " + entityType.TableName + " WHERE " + q.where
	}
	rows, err := q.dbx.Query(sqlSelect, q.args...)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return nil, nil
	}
	ret, err := fetchAllRows(rows.Rows, et)
	if err != nil {
		return nil, err
	}

	if len(ret.([]T)) == 0 {
		return nil, nil
	}

	return ret.([]T), nil
}
func (q *QrBuilder[T]) Item() ([]T, error) {
	var zero T
	et := reflect.TypeOf(zero)
	entityType, err := newEntityType(et)
	if err != nil {
		return nil, err
	}
	sqlSelect := ""
	if q.where == "" {
		sqlSelect = "SELECT * FROM " + entityType.TableName + " LIMIT 1"
	} else {
		sqlSelect = "SELECT * FROM " + entityType.TableName + " WHERE " + q.where + " LIMIT 1"
	}
	rows, err := q.dbx.Query(sqlSelect, q.args...)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return nil, nil
	}
	ret, err := fetchAllRows(rows.Rows, et)
	if err != nil {
		return nil, err
	}

	if len(ret.([]T)) == 0 {
		return nil, nil
	}

	return ret.([]T), nil
}
