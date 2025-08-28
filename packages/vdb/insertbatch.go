package vdb

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"vdb/migrate"
	"vdb/tenantDB"

	mssql "github.com/microsoft/go-mssqldb"
)

func doInsertBatchV1(db *tenantDB.TenantDB, data interface{}) error {
	maxSQLParams := 2100 - 1 //<-- khi sua lai thi no chay duoc

	migrator, err := migrate.NewMigrator(db)
	if err != nil {
		return err
	}
	if err = migrator.DoMigrates(); err != nil {
		return err
	}

	dialect := dialectFactory.Create(db.GetDriverName())

	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		typ := val.Type().String()
		return fmt.Errorf("data must be a slice, but got %s", typ)
	}
	if val.Len() == 0 {
		return nil // Không có gì để insert
	}

	typ := val.Index(0).Type()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	modelReg := ModelRegistry.GetModelByType(typ)
	if modelReg == nil {
		return NewModelError(typ)
	}
	tableName := dialect.Quote(modelReg.GetTableName())

	// Lấy danh sách cột (bỏ qua auto-increment)
	colNames := []string{}
	cols := []migrate.ColumnDef{} // lưu lại để lặp nhiều lần<-- chua thay dinh nghia ColumnMeta
	for _, col := range modelReg.GetColumns() {
		if col.IsAuto {
			continue
		}
		colNames = append(colNames, dialect.Quote(col.Name))
		cols = append(cols, col)
	}
	numCols := len(colNames)
	if numCols == 0 {
		return fmt.Errorf("no columns to insert")
	}
	maxRowsPerBatch := maxSQLParams / numCols

	// Chia batch
	for start := 0; start < val.Len(); start += maxRowsPerBatch {
		end := start + maxRowsPerBatch
		if end > val.Len() {
			end = val.Len()
		}

		args := []interface{}{}
		valuePlaceholders := []string{}
		paramIndex := 1
		for i := start; i < end; i++ {
			item := val.Index(i)

			if item.Kind() == reflect.Ptr {
				item = item.Elem()
			}

			placeholders := []string{}

			for _, col := range cols {

				fieldVal := item.FieldByName(col.Field.Name)
				if fieldVal.IsValid() {
					args = append(args, fieldVal.Interface())
				} else {
					args = append(args, nil)
				}
				if dialect.Name() == "mssql" {
					placeholders = append(placeholders, "@p"+fmt.Sprintf("%d", paramIndex))

					if paramIndex > 2100 { //<-- kg thay vao cho nay
						return fmt.Errorf("Too many parameters: %d > 2100", paramIndex)
					}
					paramIndex++

				} else {
					placeholders = append(placeholders, "?")
				}
			}

			valuePlaceholders = append(valuePlaceholders, "("+strings.Join(placeholders, ",")+")")
		}

		sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
			tableName,
			strings.Join(colNames, ","),
			strings.Join(valuePlaceholders, ","),
		)
		fmt.Println(len(args)) //<-- cho nay dung la 2100, vi luc test da dng den 1000 user moi user co 12 cot
		if _, err := db.Exec(sql, args...); err != nil {
			fmt.Println(sql[len(sql)-100 : 
			return fmt.Errorf("insert batch failed: %w", err)
		}
	}

	return nil
}
func doInsertBatch(db *tenantDB.TenantDB, data interface{}) error {
	maxSQLParams := 2100 - 1 // SQL Server giới hạn 2100 params, trừ 1 cho an toàn

	migrator, err := migrate.NewMigrator(db)
	if err != nil {
		return err
	}
	if err = migrator.DoMigrates(); err != nil {
		return err
	}

	dialect := dialectFactory.Create(db.GetDriverName())

	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice, but got %s", val.Type())
	}
	if val.Len() == 0 {
		return nil
	}

	typ := val.Index(0).Type()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	modelReg := ModelRegistry.GetModelByType(typ)
	if modelReg == nil {
		return NewModelError(typ)
	}
	tableName := dialect.Quote(modelReg.GetTableName())

	colNames := []string{}
	cols := []migrate.ColumnDef{}
	for _, col := range modelReg.GetColumns() {
		if col.IsAuto {
			continue
		}
		colNames = append(colNames, dialect.Quote(col.Name))
		cols = append(cols, col)
	}
	numCols := len(colNames)
	if numCols == 0 {
		return fmt.Errorf("no columns to insert")
	}

	maxRowsPerBatch := maxSQLParams / numCols

	// Reusable slices (reused across batches)
	args := make([]interface{}, 0, maxRowsPerBatch*numCols)
	valuePlaceholders := make([]string, 0, maxRowsPerBatch)
	placeholders := make([]string, 0, numCols)
	placeholdersPool := make([]string, maxSQLParams)
	for i := 0; i < maxSQLParams; i++ {
		placeholdersPool[i] = "@p" + strconv.Itoa(i+1)
	}

	for start := 0; start < val.Len(); start += maxRowsPerBatch {
		end := start + maxRowsPerBatch
		if end > val.Len() {
			end = val.Len()
		}

		args = args[:0]
		valuePlaceholders = valuePlaceholders[:0]

		for i := start; i < end; i++ {
			item := val.Index(i)
			if item.Kind() == reflect.Ptr {
				item = item.Elem()
			}

			paramIndex := 0
			placeholders := placeholders[:0] // reset slice

			for _, col := range cols {
				fieldVal := item.FieldByIndex(col.IndexOfField)
				if fieldVal.IsValid() {
					args = append(args, fieldVal.Interface())
				} else {
					args = append(args, nil)
				}

				if dialect.Name() == "mssql" {
					placeholders = append(placeholders, placeholdersPool[paramIndex])
					paramIndex++
				} else {
					placeholders = append(placeholders, "?")
				}
			}
			valuePlaceholders = append(valuePlaceholders, "("+strings.Join(placeholders, ",")+")")
		}

		sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
			tableName,
			strings.Join(colNames, ","),
			strings.Join(valuePlaceholders, ","),
		)

		if _, err := db.Exec(sql, args...); err != nil {
			return fmt.Errorf("insert batch failed: %w", err)
		}
	}

	return nil
}
func InsertBatchOptimized(db *tenantDB.TenantDB, data interface{}) error {
	const maxSQLParams = 2100 - 1 // MSSQL limit

	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return fmt.Errorf("data must be slice")
	}
	if val.Len() == 0 {
		return nil
	}

	elemType := val.Index(0).Type()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	model := ModelRegistry.GetModelByType(elemType)
	if model == nil {
		return NewModelError(elemType)
	}

	dialect := dialectFactory.Create(db.GetDriverName())
	tableName := dialect.Quote(model.GetTableName())

	// Chỉ lấy cột không phải auto
	var cols []migrate.ColumnDef
	var colNames []string
	for _, col := range model.GetColumns() {
		if col.IsAuto {
			continue
		}
		cols = append(cols, col)
		colNames = append(colNames, dialect.Quote(col.Name))
	}
	numCols := len(cols)
	if numCols == 0 {
		return fmt.Errorf("no insertable columns")
	}

	// MSSQL param limit
	maxRows := maxSQLParams / numCols
	total := val.Len()

	// Pre-generate MSSQL placeholder
	var preGenPlaceholders []string
	if dialect.Name() == "mssql" {
		preGenPlaceholders = make([]string, maxRows*numCols+1) // +1 để tránh out of range
		for i := 1; i < len(preGenPlaceholders); i++ {
			preGenPlaceholders[i] = "@p" + strconv.Itoa(i)
		}
	}

	args := make([]interface{}, 0, maxRows*numCols)
	valuePlaceholders := make([]string, 0, maxRows)
	placeholders := make([]string, 0, numCols)

	for start := 0; start < total; start += maxRows {
		end := start + maxRows
		if end > total {
			end = total
		}

		args = args[:0]
		valuePlaceholders = valuePlaceholders[:0]
		paramIndex := 1

		for i := start; i < end; i++ {
			item := val.Index(i)
			if item.Kind() == reflect.Ptr {
				item = item.Elem()
			}

			placeholders = placeholders[:0]
			for _, col := range cols {
				fieldVal := item.FieldByIndex(col.IndexOfField)
				if fieldVal.IsValid() {
					args = append(args, fieldVal.Interface())
				} else {
					args = append(args, nil)
				}

				if dialect.Name() == "mssql" {
					placeholders = append(placeholders, preGenPlaceholders[paramIndex])
					paramIndex++
				} else {
					placeholders = append(placeholders, "?")
				}
			}

			valuePlaceholders = append(valuePlaceholders, "("+strings.Join(placeholders, ",")+")")
		}

		query := fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES %s",
			tableName,
			strings.Join(colNames, ","),
			strings.Join(valuePlaceholders, ","),
		)

		if _, err := db.Exec(query, args...); err != nil {
			return fmt.Errorf("insert batch failed: %w", err)
		}
	}

	return nil
}
func InsertBatchTVP(db *sql.DB, data interface{}, typeName string, procName string) error {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice, got %s", val.Type())
	}
	if val.Len() == 0 {
		return nil
	}

	// Lấy kiểu struct gốc
	elemType := val.Type().Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	// Ánh xạ schema
	model := ModelRegistry.GetModelByType(elemType)
	if model == nil {
		return NewModelError(elemType)
	}
	columns := model.GetColumns()

	// Chuẩn bị dữ liệu dạng [][]interface{}
	tvRows := make([][]interface{}, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		row := []interface{}{}
		elem := val.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		for _, col := range columns {
			if col.IsAuto {
				continue
			}
			fv := elem.FieldByIndex(col.IndexOfField)
			if fv.IsValid() && fv.CanInterface() {
				row = append(row, fv.Interface())
			} else {
				row = append(row, nil)
			}
		}
		tvRows = append(tvRows, row)
	}

	// Tên TYPE phải ở dạng [schema].[type] như "dbo.MyUserType"
	tvp := mssql.TVP{
		TypeName: typeName,
		Value:    tvRows,
	}

	// Gọi stored procedure dạng: EXEC procName @ParamName = @tvp
	_, err := db.Exec(fmt.Sprintf("EXEC %s @tvp", procName),
		sql.Named("tvp", tvp),
	)

	if err != nil {
		return fmt.Errorf("InsertBatchTVP failed: %w", err)
	}

	return nil
}

func init() {
	tenantDB.OnDbInsertBatchFunc = InsertBatchOptimized
}
