package dbx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

func Insert(db *DBXTenant, entity interface{}) error {

	return db.Insert(entity)
}
func InsertWithContext(ctx context.Context, db *DBXTenant, entity interface{}) error {

	return db.InsertWithContext(ctx, entity)
}
func (ctx *DBXTenant) Insert(entity interface{}) error {
	validateSize(entity)
	if ctx.DB == nil {
		panic("please open TenantDbContext first")
	}
	if ctx.cfg.Driver == "postgres" {
		err := postgresMigrateEntity(ctx.DB, ctx.TenantDbName, entity)
		if err != nil {
			return err
		}
	} else if ctx.cfg.Driver == "mysql" {

		err := mySqlMigrateEntity(ctx.DB, ctx.TenantDbName, entity)
		if err != nil {
			return err
		}
	} else if ctx.cfg.Driver == "mssql" {

		err := mssqlSqlMigrateEntity(ctx.DB, ctx.TenantDbName, entity)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("not support driver %s", ctx.cfg.Driver)
	}

	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("entity must be a struct or a pointer to a struct")
	}
	tblInfo, err := newEntityType(typ)
	if err != nil {
		return err
	}
	if ctx.cfg.Driver == "postgres" {
		return ctx.pgInsert(nil, tblInfo, entity)
	} else if ctx.cfg.Driver == "mysql" {
		return ctx.mysqlInsert(nil, tblInfo, entity)
		//return ctx.myInsert(tblInfo, entity)
	} else if ctx.cfg.Driver == "mssql" {
		return ctx.mssqlInsert(nil, tblInfo, entity)
		//return ctx.myInsert(tblInfo, entity)
	} else {
		return fmt.Errorf("not support driver %s", ctx.cfg.Driver)
	}
}
func (ctx *DBXTenant) InsertWithContext(cntx context.Context, entity interface{}) error {
	if errSize := validateSize(entity); errSize != nil {
		return errSize
	}
	if ctx.DB == nil {
		panic("please open TenantDbContext first")
	}
	if ctx.cfg.Driver == "postgres" {
		err := postgresMigrateEntity(ctx.DB, ctx.TenantDbName, entity)
		if err != nil {
			return err
		}
	} else if ctx.cfg.Driver == "mysql" {

		err := mySqlMigrateEntity(ctx.DB, ctx.TenantDbName, entity)
		if err != nil {
			return err
		}
	} else if ctx.cfg.Driver == "mssql" {

		err := mssqlSqlMigrateEntity(ctx.DB, ctx.TenantDbName, entity)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("not support driver %s", ctx.cfg.Driver)
	}

	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("entity must be a struct or a pointer to a struct")
	}
	tblInfo, err := newEntityType(typ)
	if err != nil {
		return err
	}
	if ctx.cfg.Driver == "postgres" {
		return ctx.pgInsert(cntx, tblInfo, entity)
	} else if ctx.cfg.Driver == "mysql" {
		return ctx.mysqlInsert(cntx, tblInfo, entity)
		//return ctx.myInsert(tblInfo, entity)
	} else if ctx.cfg.Driver == "mssql" {
		return ctx.mssqlInsert(cntx, tblInfo, entity)
		//return ctx.myInsert(tblInfo, entity)
	} else {
		return fmt.Errorf("not support driver %s", ctx.cfg.Driver)
	}
}
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

	execSql2, err := ctx.compiler.parseInsertSQL(parseInsertInfo{
		TableName:        tblInfo.TableName,
		DefaultValueCols: tblInfo.getDefaultValueColsNames(),
		// ReturnColAfterInsert: tblInfo.autoValueColsName,
		SqlInsert:    execSql,
		keyColsNames: tblInfo.GetPrimaryKeyName(),
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
		return err
	}
	insertedId, err := result.LastInsertId()
	if err != nil {
		// tx.Rollback()
		return err
	}
	v := reflect.ValueOf(entity)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("entity must be a non-nil pointer")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("entity must point to a struct")
	}

	idField := v.FieldByName(tblInfo.primaryKeyNames[0])
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
func (ctx *DBXTenant) pgInsert(cntx context.Context, tblInfo *EntityType, entity interface{}) error {
	if errSize := validateSize(entity); errSize != nil {
		return errSize
	}

	err := postgresMigrateEntity(ctx.DB, ctx.TenantDbName, entity)

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

	execSql2, err := ctx.compiler.parseInsertSQL(parseInsertInfo{
		TableName:        tblInfo.TableName,
		DefaultValueCols: tblInfo.getDefaultValueColsNames(),
		// ReturnColAfterInsert: tblInfo.autoValueColsName,
		SqlInsert:    execSql,
		keyColsNames: tblInfo.GetPrimaryKeyName(),
	})
	//.OnParseInsertSQL(walker, execSql, tblInfo.AutoValueColsName, []string{})
	if err != nil {
		return err
	}
	// resultArray := []interface{}{}
	//ctx.Open()
	var rw *sql.Rows
	var errQr error
	if cntx == nil {

		rw, errQr = ctx.DB.Query((*execSql2), dataInsert.Params...)
	} else {
		start := time.Now()
		rw, errQr = ctx.DB.QueryContext(cntx, (*execSql2), dataInsert.Params...)
		n := time.Since(start).Milliseconds()
		defer func() {
			fmt.Println(red, "time", n, reset, green, (*execSql2), reset)
			log.Println(red, "time", n, reset, green, (*execSql2), reset)
		}()
		if errQr != nil {
			return PostgresErrorParser.ParseError(cntx, ctx.DB, errQr)

		}

	}

	if errQr != nil {

		return PostgresErrorParser.ParseError(cntx, ctx.DB, errQr)
	}
	defer rw.Close()
	cols, err := rw.Columns()
	if err != nil {
		return err
	}
	// if len(cols) != 1 {
	// 	return fmt.Errorf("insert failed, expect 1 column, but got %d", len(cols))
	// }
	// insertedId := 0
	// for rw.Next() {
	// 	err := rw.Scan(&insertedId)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// if err != nil {
	// 	return err
	// }

	for rw.Next() {
		err := scanRowToStruct(rw, entity, cols) // thay may cai vong lap o duoi ban ham nay chay OK
		if err != nil {
			return err
		}

	}

	if err != nil {
		return err
	}
	return nil
}
func (ctx *DBXTenant) mssqlInsertDelete(tblInfo *EntityType, entity interface{}) error {
	err := mssqlSqlMigrateEntity(ctx.DB, ctx.TenantDbName, entity)

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

	execSql2, err := ctx.compiler.parseInsertSQL(parseInsertInfo{
		TableName:        tblInfo.TableName,
		DefaultValueCols: tblInfo.getDefaultValueColsNames(),
		// ReturnColAfterInsert: tblInfo.autoValueColsName,
		SqlInsert:    execSql,
		keyColsNames: tblInfo.GetPrimaryKeyName(),
	})
	//.OnParseInsertSQL(walker, execSql, tblInfo.AutoValueColsName, []string{})
	if err != nil {
		return err
	}
	// resultArray := []interface{}{}
	//ctx.Open()

	tx, err := ctx.Begin()
	if err != nil {
		return err
	}

	var finalErr error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if finalErr != nil {
			tx.Rollback()
		}
	}()

	// === Sửa lỗi Scan NULL: Sử dụng sql.NullInt64 ===
	var insertedNullId sql.NullInt64 // Sử dụng sql.NullInt64 để chứa giá trị có thể là NULL

	// Thực thi truy vấn và scan kết quả
	err = tx.QueryRow(*execSql2, dataInsert.Params...).Scan(&insertedNullId)
	if err != nil {
		finalErr = err
		return fmt.Errorf("lỗi khi thực thi INSERT hoặc scan ID: %w", err)
	}

	// === Kiểm tra giá trị NullInt64 ===
	if insertedNullId.Valid {
		// Nếu giá trị hợp lệ (không phải NULL)
		insertedID := insertedNullId.Int64 // Lấy giá trị int64 thực tế
		fmt.Printf("Đã chèn và lấy được EmployeeId: %d\n", insertedID)

		// --- Gán ID vào entity của bạn ---
		entityVal := reflect.ValueOf(entity)
		if entityVal.Kind() == reflect.Ptr {
			entityVal = entityVal.Elem()
		}
		idField := entityVal.FieldByName(tblInfo.GetPrimaryKeyName()[0])
		if idField.IsValid() && idField.CanSet() {
			switch idField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				idField.SetInt(insertedID)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				idField.SetUint(uint64(insertedID))
			}
		} else {
			fmt.Printf("Cảnh báo: Không thể gán insertedID %d vào trường '%s' của entity\n", insertedID, tblInfo.GetPrimaryKeyName())
		}
	} else {
		// === QUAN TRỌNG: Nếu Valid là FALSE, có nghĩa là SCOPE_IDENTITY() trả về NULL ===
		// Điều này HẦU HẾT các lần là do câu lệnh INSERT bị lỗi.
		// Bạn cần kiểm tra xem các tham số dataInsert.Params có đúng không,
		// và bảng Employees có vi phạm ràng buộc NOT NULL, UNIQUE, PRIMARY KEY, FOREIGN KEY nào không.
		finalErr = fmt.Errorf("không thể lấy EmployeeId: SCOPE_IDENTITY() trả về NULL. Có thể INSERT thất bại do vi phạm ràng buộc.")
		return finalErr
	}

	// Commit transaction sau khi đã lấy được ID và xử lý kết quả.
	err = tx.Commit()
	if err != nil {
		finalErr = err
		return err
	}

	return nil
}
