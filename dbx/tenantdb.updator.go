package dbx

import (
	"database/sql"
	"errors"
	"fmt"
)

type updaterWhere struct {
	tableName string
	where     string
	args      []interface{}
	db        *DBXTenant
	err       error
}
type updaterWhereSet struct {
	where   *updaterWhere
	setters []string
	values  []interface{}
	err     error
}
type updaterWhereExecutor struct {
	setter *updaterWhereSet
}

func (t *DBXTenant) Update(entity interface{}) *updaterWhere {

	if t == nil {
		return nil
	}
	tblInfo, err := Entities.GetEntityType(entity)
	if err != nil {
		return &updaterWhere{
			tableName: "",
			db:        t,
			err:       err,
		}
	}
	tableName := tblInfo.TableName

	return &updaterWhere{
		tableName: tableName,
		db:        t,
		err:       err,
	}

}
func (u *updaterWhere) Where(where string, args ...interface{}) *updaterWhereSet {
	u.where = where
	u.args = args
	return &updaterWhereSet{
		where: u,
	}
}
func (u *updaterWhereSet) Set(args ...interface{}) *updaterWhereExecutor {
	ret := &updaterWhereExecutor{
		setter: u,
	}
	lenOfArgs := len(args)
	if lenOfArgs%2 != 0 {
		u.err = errors.New("invalid number of arguments")
		return ret
	}
	for i := 0; i < lenOfArgs; i += 2 {
		field := args[i]
		if fieldName, ok := field.(string); ok {
			if Entities.CheckField(u.where.tableName, fieldName) {
				u.setters = append(u.setters, fieldName)
				u.values = append(u.values, args[i+1])
			} else {
				u.err = fmt.Errorf("field %s not found in table %s", fieldName, u.where.tableName)
				return ret
			}

		} else {
			u.err = errors.New("invalid field name")
			return ret
		}

	}
	return ret
}
func (exec *updaterWhereExecutor) GetSql() (string, []interface{}, error) {
	args := []interface{}{}
	if exec.setter.err != nil {
		return "", args, exec.setter.err
	}
	if exec.setter.where.err != nil {
		return "", args, exec.setter.where.err
	}
	sql := "UPDATE " + exec.setter.where.tableName + " SET "
	for i, setter := range exec.setter.setters {
		sql += setter + " = ?"
		args = append(args, exec.setter.values[i])
		if i < len(exec.setter.setters)-1 {
			sql += ", "
		}
	}
	sql += " WHERE " + exec.setter.where.where
	args = append(args, exec.setter.where.args...)
	return sql, args, nil
}
func (exec *updaterWhereExecutor) Execute() (sql.Result, error) {

	sql, args, err := exec.GetSql()
	if err != nil {
		return nil, err
	}
	return exec.setter.where.db.Exec(sql, args...)

}
