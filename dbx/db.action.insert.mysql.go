package dbx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func (ctx *DBXTenant) mysqlInsert(cntx context.Context, tblInfo *EntityType, entity interface{}) error {
	if errSize := validateSize(entity); errSize != nil {
		return errSize
	}
	err := mySqlMigrateEntity(ctx.DB, ctx.TenantDbName, entity)

	if err != nil {
		return err
	}
	dataInsert, err := createInsertCommand(entity, tblInfo)

	if err != nil {
		return err
	}

	execSql, err := ctx.compiler.Parse(dataInsert.Sql)
	if err != nil {
		return err
	}
	pkAutoCols := tblInfo.GetPkAutoCos()

	execSql2, err := ctx.compiler.parseInsertSQL(parseInsertInfo{
		TableName:        tblInfo.TableName,
		DefaultValueCols: tblInfo.getDefaultValueColsNames(),
		// ReturnColAfterInsert: tblInfo.autoValueColsName,
		SqlInsert:    execSql,
		keyColsNames: pkAutoCols,
	})
	//.OnParseInsertSQL(walker, execSql, tblInfo.AutoValueColsName, []string{})
	if err != nil {
		return err
	}
	// resultArray := []interface{}{}
	//ctx.Open()
	sqlInsert := strings.Split(*execSql2, "\n")[0]
	// sqlSelect := strings.Split(*execSql2, "\n")[1]
	db := ctx.DB
	// tx, err := db.Begin()
	if err != nil {
		return err
	}
	// start := time.Now()
	var result sql.Result
	if cntx == nil {
		result, err = db.Exec(sqlInsert, dataInsert.Params...)
	} else {
		result, err = db.ExecContext(cntx, sqlInsert, dataInsert.Params...)
	}

	// fmt.Println("Insert time: ", time.Since(start).Milliseconds())
	if err != nil {
		// tx.Rollback()
		return mySqlErrorParser.ParseError(cntx, db, err)
	}
	if len(pkAutoCols) == 0 { //nothing update to entity
		return nil
	}
	insertedId, err := result.LastInsertId()
	if err != nil {
		// tx.Rollback()
		return mySqlErrorParser.ParseError(cntx, db, err)
	}
	v := reflect.ValueOf(entity)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("entity must be a non-nil pointer")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("entity must point to a struct")
	}

	idField := v.FieldByName(pkAutoCols[0])
	if !idField.IsValid() {
		return errors.New("field 'Id' not found in struct")
	}
	if !idField.CanSet() {
		return errors.New("cannot set 'Id' field")
	}

	switch idField.Kind() {
	case reflect.Int, reflect.Int64:
		idField.SetInt(insertedId)
	default:
		return fmt.Errorf("unsupported 'Id' field type: %s", idField.Kind())
	}

	return nil

}
