package test

import (
	"dbv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// import (
// 	"dbv"
// 	"dbv/pkg_test/models"
// 	"reflect"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// )

func TestMssqlQueryBuilder(t *testing.T) {
	mssqlDns := "sqlserver://sa:123456@localhost?database=a003"
	db, err := dbv.Open("sqlserver", mssqlDns)
	// dialect := dbv.DialectFactory.Create(db.GetDriverName())
	assert.NoError(t, err)
	err = db.Ping()
	assert.NoError(t, err)

	type OutPut struct { //<-- vi du day la dau ra
		Username       string
		DepartmentName string
		PositionName   string
		Email          string
		UserBirthday   time.Time
	}
	for i := 0; i < 10; i++ {
		qr := dbv.Qr().From("user").Select(
			"concat(firstName, ?,lastName) as FullName",
			"Code",
			"salary*? as Salary", dbv.Lit(" "), 0.12)
		qr.Where("id =? or id =?", 10, 100)

		sql, args := qr.BuildSQL(db)
		expectSql := "SELECT CONCAT([T1].[first_name], @p1, [T1].[last_name]), [T1].[code] AS [Code], [T1].[salary] * @p2 AS [Salary] FROM [users] AS [T1] WHERE [T1].[id] = @p3 OR [T1].[id] = @p4"
		assert.Equal(t, expectSql, sql)
		assert.Equal(t, []interface{}{" ", 0.12, 10, 100}, args)
	}

}
func BenchmarkTestMssqlQueryBuilder(t *testing.B) {
	mssqlDns := "sqlserver://sa:123456@localhost?database=a003"
	db, err := dbv.Open("sqlserver", mssqlDns)
	// dialect := dbv.DialectFactory.Create(db.GetDriverName())
	assert.NoError(t, err)
	err = db.Ping()
	assert.NoError(t, err)
	dbv.NewExprCompiler(db)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {

		qr := dbv.Qr().From("user").Select(
			"concat(firstName, ?,lastName) as FullName",
			"Code",
			"salary*? as Salary", dbv.Lit(" "), 0.12)
		qr.Where("id =? or id =?", 10, 100)
		sql, args := qr.BuildSQL(db)
		expectSql := "SELECT CONCAT([T1].[first_name], @p1, [T1].[last_name]), [T1].[code] AS [Code], [T1].[salary] * @p2 AS [Salary] FROM [users] AS [T1] WHERE [T1].[id] = @p3 OR [T1].[id] = @p4"
		assert.Equal(t, expectSql, sql)
		assert.Equal(t, []interface{}{" ", 0.12, 10, 100}, args)
	}
}

// func TestExecToArray(t *testing.T) {
// 	mssqlDns := "sqlserver://sa:123456@localhost?database=a003"
// 	db, err := dbv.Open("sqlserver", mssqlDns)
// 	assert.NoError(t, err)
// 	err = db.Ping()
// 	sql := `SELECT        users.username AS Username, departments.name AS DepartmentName,
// positions.name AS PositionName, users.email AS Email, users.[address] AS [Address],
// users.phone AS Phone, users.birthday AS BirthDay,
//                          departments.code AS DepartmentCode
// FROM            users INNER JOIN
//                          positions ON users.id = positions.id INNER JOIN
//                          departments ON users.id = departments.id`
// 	type UserFullInfo struct {
// 		Username       *string
// 		DepartmentName *string
// 		PositionName   *string
// 		Email          *string
// 		Address        *string
// 		Phone          *string
// 		BirthDay       *time.Time
// 		DepartmentCode *string
// 	}
// 	items, err := db.ExecToArray(reflect.TypeOf(UserFullInfo{}), sql)
// 	assert.NoError(t, err)

// 	assert.Equal(t, 9604, len(items))

// }
// func BenchmarkTestExecToArray(t *testing.B) {
// 	mssqlDns := "sqlserver://sa:123456@localhost?database=a003"
// 	db, err := dbv.Open("sqlserver", mssqlDns)
// 	assert.NoError(t, err)
// 	err = db.Ping()
// 	sql := `SELECT        users.username AS Username, departments.name AS DepartmentName,
// positions.name AS PositionName, users.email AS Email, users.[address] AS [Address],
// users.phone AS Phone, users.birthday AS BirthDay,
//                          departments.code AS DepartmentCode
// FROM            users INNER JOIN
//                          positions ON users.id = positions.id INNER JOIN
//                          departments ON users.id = departments.id`
// 	type UserFullInfo struct {
// 		Username       string //<-- neu cho nay thay doi thanh Username string bi loi, vi username trong cau sql tren bang null, fx duoc kg
// 		DepartmentName string
// 		PositionName   string
// 		Email          string
// 		Address        string
// 		Phone          string
// 		BirthDay       time.Time
// 		DepartmentCode string
// 	}
// 	for i := 0; i < t.N; i++ {
// 		items, err := db.ExecToArray(reflect.TypeOf(UserFullInfo{}), sql)
// 		assert.NoError(t, err)

// 		assert.Equal(t, 9604, len(items))
// 	}
// }
// func TestMssqlSelect(t *testing.T) {
// 	mssqlDns := "sqlserver://sa:123456@localhost?database=a003"
// 	db, err := dbv.Open("sqlserver", mssqlDns)
// 	assert.NoError(t, err)
// 	err = db.Ping()
// 	assert.NoError(t, err)
// 	m, err := dbv.NewMigrator(db)
// 	assert.NoError(t, err)
// 	err = m.DoMigrates()
// 	assert.NoError(t, err)
// 	items, err := dbv.SelectAll[models.User](db)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 24004, len(items))

// 	// err = db.Select(&users, "SELECT * FROM users")
// 	// assert.NoError(t, err)
// 	// assert.Equal(t, 1, len(users))
// }

// func BenchmarkTestMssqlSelect(t *testing.B) {
// 	mssqlDns := "sqlserver://sa:123456@localhost?database=a003"
// 	db, err := dbv.Open("sqlserver", mssqlDns)
// 	assert.NoError(t, err)
// 	err = db.Ping()
// 	assert.NoError(t, err)
// 	m, err := dbv.NewMigrator(db)
// 	assert.NoError(t, err)
// 	err = m.DoMigrates()
// 	assert.NoError(t, err)
// 	for i := 0; i < t.N; i++ {
// 		items, err := dbv.SelectAll[models.User](db)
// 		assert.NoError(t, err)
// 		assert.Equal(t, 24004, len(items))
// 	}

// }
// func TestMssqlUserInsertBatch(t *testing.T) {
// 	msssqlDns := "sqlserver://sa:123456@localhost?database=a003"
// 	db, err := dbv.Open("sqlserver", msssqlDns)
// 	assert.NoError(t, err)
// 	users := []models.User{}
// 	for i := 0; i < 1000; i++ {
// 		user := models.User{
// 			Name:       "Dylan",
// 			PositionID: 1,
// 			DeptID:     1,
// 			BaseModel: models.BaseModel{
// 				CreatedAt: time.Now(),
// 			},
// 		}
// 		users = append(users, user)

// 	}
// 	// dbv.InsertBatch(db, users)
// 	err = db.InsertBatch(users)
// 	assert.NoError(t, err)
// }
// func BenchmarkTestMssqlUserInsertBatch(t *testing.B) {
// 	msssqlDns := "sqlserver://sa:123456@localhost?database=a003"
// 	db, err := dbv.Open("sqlserver", msssqlDns)
// 	assert.NoError(t, err)
// 	t.ResetTimer()
// 	for i := 0; i < t.N; i++ {
// 		users := []models.User{}
// 		for i := 0; i < 1000; i++ {
// 			user := models.User{
// 				Name:       "Dylan",
// 				PositionID: 1,
// 				DeptID:     1,
// 				BaseModel: models.BaseModel{
// 					CreatedAt: time.Now(),
// 				},
// 			}
// 			users = append(users, user)

// 		}
// 		// dbv.InsertBatch(db, users)
// 		err = db.InsertBatch(users) /*<-- tr den ham moi bang cach dat
// 				args, valuePlaceholders, và placeholders được khai báo một lần bên ngoài vòng lặp batch.

// 		Dùng slice = slice[:0] để reset mà không cấp phát lại bộ nhớ.

// 		Sử dụng reflect.Value.FieldByIndex() đúng cách để lấy giá trị field (kể cả field embedded trong struct lồng nhau).
// 		*/
// 		assert.NoError(t, err)
// 	}
// }
// func TestMssqlInsertDepartment(t *testing.T) {
// 	msssqlDns := "sqlserver://sa:123456@localhost?database=a002"
// 	db, err := dbv.Open("sqlserver", msssqlDns)
// 	assert.NoError(t, err)
// 	err = db.Ping()

// 	assert.NoError(t, err)
// 	m, err := dbv.NewMigrator(db)
// 	assert.NoError(t, err)
// 	err = m.DoMigrates()
// 	assert.NoError(t, err)

// 	dep := models.Department{
// 		Name: "HRM",
// 		Code: "HR Department",
// 		BaseModel: models.BaseModel{
// 			CreatedAt: time.Now(),
// 		},
// 	}
// 	err = db.Insert(&dep)
// 	assert.NoError(t, err)
// }

// func TestMssqlInsertDepartmentRaw(t *testing.T) {
// 	msssqlDns := "sqlserver://sa:123456@localhost?database=a002"
// 	db, err := dbv.Open("sqlserver", msssqlDns)
// 	assert.NoError(t, err)
// 	err = db.Ping()
// 	assert.NoError(t, err)
// 	dep := models.Department{
// 		Name: "HRM",
// 		Code: "HR Department",
// 		BaseModel: models.BaseModel{
// 			CreatedAt: time.Now(),
// 		},
// 	}
// 	sql := "INSERT INTO [departments] ([name], [code], [parent_id], [created_at], [updated_at], [description]) OUTPUT INSERTED.[id] VALUES (@p1, @p2, @p3, @p4, @p5, @p6)"
// 	row := db.QueryRow(sql, dep.Name, dep.Code, dep.ParentID, dep.CreatedAt, dep.UpdatedAt, dep.Description)
// 	assert.NoError(t, err)
// 	var id int
// 	err = row.Scan(&id)
// 	assert.NoError(t, err)
// 	dep.ID = id
// 	assert.Equal(t, dep.ID, id)

// }
// func BenchmarkTestMssqlInsertDepartmentRaw(t *testing.B) {
// 	msssqlDns := "sqlserver://sa:123456@localhost?database=a002"
// 	db, err := dbv.Open("sqlserver", msssqlDns)
// 	assert.NoError(t, err)
// 	err = db.Ping()

// 	assert.NoError(t, err)
// 	m, err := dbv.NewMigrator(db)
// 	assert.NoError(t, err)
// 	err = m.DoMigrates()
// 	assert.NoError(t, err)
// 	for i := 0; i < t.N; i++ {

// 		dep := models.Department{
// 			Name: "HRM",
// 			Code: "HR Department",
// 			BaseModel: models.BaseModel{
// 				CreatedAt: time.Now(),
// 			},
// 		}
// 		err = db.Insert(&dep)
// 		assert.NoError(t, err)
// 	}
// }
// func BenchmarkTestMssqlInsertDepartment(t *testing.B) {
// 	msssqlDns := "sqlserver://sa:123456@localhost?database=a002"
// 	db, err := dbv.Open("sqlserver", msssqlDns)
// 	assert.NoError(t, err)
// 	err = db.Ping()

// 	assert.NoError(t, err)
// 	m, err := dbv.NewMigrator(db)
// 	assert.NoError(t, err)
// 	err = m.DoMigrates()
// 	assert.NoError(t, err)
// 	t.ResetTimer()
// 	for i := 0; i < t.N; i++ {
// 		dep := models.Department{
// 			Code: "HRM",
// 			Name: "HR Department",
// 			BaseModel: models.BaseModel{
// 				CreatedAt: time.Now(),
// 			},
// 		}
// 		err = db.Insert(&dep)
// 		assert.NoError(t, err)
// 	}

// }
// func BenchmarkTestMssqlInsertDepartmentByTx(t *testing.B) {
// 	msssqlDns := "sqlserver://sa:123456@localhost?database=a003"
// 	db, err := dbv.Open("sqlserver", msssqlDns)
// 	assert.NoError(t, err)
// 	err = db.Ping()

// 	assert.NoError(t, err)
// 	m, err := dbv.NewMigrator(db)
// 	assert.NoError(t, err)
// 	err = m.DoMigrates()
// 	assert.NoError(t, err)
// 	t.ResetTimer()
// 	tx, err := db.Begin()
// 	assert.NoError(t, err)
// 	for i := 0; i < t.N; i++ {
// 		dep := models.Department{
// 			Code: "HRM",
// 			Name: "HR Department",
// 			BaseModel: models.BaseModel{
// 				CreatedAt: time.Now(),
// 			},
// 		}
// 		pos := models.Position{
// 			Code:  "MNG",
// 			Name:  "HR Manager",
// 			Title: "Manager",
// 			Level: 1,
// 			BaseModel: models.BaseModel{
// 				CreatedAt: time.Now(),
// 			},
// 		}

// 		err = tx.Insert(&dep, &pos)
// 		assert.NoError(t, err)
// 		if err != nil {
// 			err = tx.Rollback()
// 			assert.NoError(t, err)
// 			return
// 		}
// 	}
// 	err = tx.Commit()
// 	assert.NoError(t, err)

// }
