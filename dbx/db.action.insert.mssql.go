package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

func (ctx *DBXTenant) mssqlInsert(cntx context.Context, tblInfo *EntityType, entity interface{}) error {
	// Kiểm tra ctx.DB, nếu chưa mở thì panic (giữ nguyên)
	if ctx.DB == nil {
		panic("please open TenantDbContext first")
	}

	// Thực hiện migrate (giữ nguyên)
	err := mssqlSqlMigrateEntity(ctx.DB, ctx.TenantDbName, entity)
	if err != nil {
		return err
	}

	// Tạo lệnh INSERT (giữ nguyên)
	dataInsert, err := createInsertCommand(entity, tblInfo)
	if err != nil {
		return err
	}

	// Parse SQL (giữ nguyên)
	execSql, err := ctx.compiler.Parse(dataInsert.Sql)
	if err != nil {
		return err
	}

	// Chuẩn bị SQL cuối cùng với SCOPE_IDENTITY() (giữ nguyên)
	// execSql2 là biến chứa câu lệnh INSERT; SELECT ID = convert(bigint, SCOPE_IDENTITY());
	execSql2, err := ctx.compiler.parseInsertSQL(parseInsertInfo{
		TableName:        tblInfo.TableName,
		DefaultValueCols: tblInfo.getDefaultValueColsNames(),
		SqlInsert:        execSql, // execSql là câu lệnh INSERT gốc
		keyColsNames:     tblInfo.GetPrimaryKeyName(),
	})
	if err != nil {
		return err
	}

	// === THAY ĐỔI: LOẠI BỎ TRANSACTION ===
	// Không còn Begin(), Rollback(), Commit()

	// Sử dụng sql.NullInt64 để chứa giá trị ID có thể là NULL
	var insertedNullId sql.NullInt64

	var qr *sql.Row
	if cntx == nil {
		qr = ctx.DB.QueryRow(*execSql2, dataInsert.Params...)
	} else {

		qr = ctx.DB.QueryRowContext(cntx, *execSql2, dataInsert.Params...)

	}

	if qr.Err() != nil {
		if dbxErr := MssqlErrorParser.ParseError(cntx, ctx.DB, qr.Err()); dbxErr != nil {
			return dbxErr
		}
		return qr.Err()
	}
	err = qr.Scan(&insertedNullId)
	if err != nil {
		if dbxErr := MssqlErrorParser.ParseError(cntx, ctx.DB, err); dbxErr != nil {
			return dbxErr
		}
		return err
	}

	// Kiểm tra giá trị NullInt64
	if insertedNullId.Valid {
		insertedID := insertedNullId.Int64

		// Gán ID vào entity của bạn (giữ nguyên)
		entityVal := reflect.ValueOf(entity)
		if entityVal.Kind() == reflect.Ptr {
			entityVal = entityVal.Elem()
		}
		idField := entityVal.FieldByName(tblInfo.GetPrimaryKeyName()[0]) // Lấy tên cột khóa chính
		if idField.IsValid() && idField.CanSet() {
			switch idField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				idField.SetInt(insertedID)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				// Kiểm tra xem insertedID có phải là số âm không.
				// Nếu là số âm, có thể đó là một lỗi từ cơ sở dữ liệu
				// hoặc một trường hợp không mong muốn.
				if insertedID < 0 {
					// Error handling: Log, return error, or set default value 0
					// Example: log.Printf("Warning: Negative ID received from DB: %d", insertedID)
					// Or: return errors.New("negative ID received")
					// For this case, you may want to set 0 or ignore
					idField.SetUint(0) // Gán 0 nếu bạn muốn bỏ qua ID âm

				} else {
					idField.SetUint(uint64(insertedID))
				}
			}
		} else {
			return fmt.Errorf("cannot set '%s' field", tblInfo.GetPrimaryKeyName()[0])
		}
	}

	return nil // Trả về nil nếu mọi thứ thành công
}
