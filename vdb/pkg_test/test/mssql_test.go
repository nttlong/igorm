package test

import (
	"testing"
	"time"
	"vdb"
	"vdb/pkg_test/models"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	mssqlDns := "sqlserver://sa:123456@localhost?database=a004"
	db, err := vdb.Open("sqlserver", mssqlDns)

	assert.NoError(t, err)
	defer db.Close()
	r, err := db.Delete(&models.User{}, "userId=?", "d8cbde8c-9e9f-4ff7-9c35-50fb7d408ef9")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), r)

}
func TestMssqlDbTestFind(t *testing.T) {
	mssqlDns := "sqlserver://sa:123456@localhost?database=a001"
	db, err := vdb.Open("sqlserver", mssqlDns)

	assert.NoError(t, err)
	defer db.Close()
	user := models.User{}
	err = db.First(&user, "id=?", 1)
	assert.NoError(t, err)

}
func TestMssqlDbQuery(t *testing.T) {
	mssqlDns := "sqlserver://sa:123456@localhost?database=a003"
	db, err := vdb.Open("sqlserver", mssqlDns)

	// dialect := vdb.DialectFactory.Create(db.GetDriverName())
	assert.NoError(t, err)
	err = db.Ping()
	assert.NoError(t, err)
	qr := db.From("users").OrderBy("id").Select("ID, Name").OrderBy("id").OffsetLimit(0, 10000)
	sql, _ := qr.BuildSql()
	expectSql := "SELECT [T1].[id] AS [id], [T1].[name] AS [name] FROM [users] AS [T1] OFFSET 0 ROWS FETCH NEXT 100 ROWS ONLY"
	assert.Equal(t, expectSql, sql)
	users := []models.User{}
	err = qr.ToArray(&users)
	assert.NoError(t, err)
	assert.Equal(t, 100, len(users))

}
func BenchmarkTestMssqlDbQuery(t *testing.B) {
	mssqlDns := "sqlserver://sa:123456@localhost?database=a003"
	db, err := vdb.Open("sqlserver", mssqlDns)
	// dialect := vdb.DialectFactory.Create(db.GetDriverName())
	assert.NoError(t, err)
	err = db.Ping()
	assert.NoError(t, err)
	for i := 0; i < t.N; i++ {
		qr := db.From("users").Select(
			"id",
			"Name",
			"concat(Name, ?,Name) as Email",
			"Birthday",
			db.Lit(" "),
		).OrderBy("ID").OffsetLimit(0, 10000)

		users := []models.User{}
		err = qr.ToArray(&users)
		assert.NoError(t, err)
		assert.Equal(t, 10000, len(users))
	}
}
func TestMssqlSelectJoinExprAndSort(t *testing.T) {
	mssqlDns := "sqlserver://sa:123456@localhost?database=a003"
	db, err := vdb.Open("sqlserver", mssqlDns)
	// dialect := vdb.DialectFactory.Create(db.GetDriverName())
	assert.NoError(t, err)
	err = db.Ping()
	assert.NoError(t, err)

	type OutPut struct { //<-- vi du day la dau ra
		Username       *string
		DepartmentName *string
		PositionName   *string
		Email          *string
		UserBirthday   *time.Time
	}
	qr := vdb.Qr().From("user u").LeftJoin("department d", "u.deptId= d.id")
	qr.LeftJoin("position p", "u.positionId= p.id")
	qr.Select(
		"u.username",
		"d.name DepartmentName",
		"p.name PositionName",
		"u.email",
		"u.birthday UserBirthday",
	)
	qr.OrderBy("u.id")
	qr.OffsetLimit(0, 100)

	sql, _ := qr.BuildSQL(db.TenantDB)
	expectSql := "SELECT [u].[username] AS [username], [d].[name] AS [DepartmentName], [p].[name] AS [PositionName], [u].[email] AS [email], [u].[birthday] AS [UserBirthday] FROM [users] AS [u] LEFT JOIN [departments] AS [d] ON [u].[dept_id] = [d].[id] LEFT JOIN [positions] AS [p] ON [u].[position_id] = [p].[id] ORDER BY [u].[id] ASC OFFSET 0 ROWS FETCH NEXT 100 ROWS ONLY"
	assert.Equal(t, expectSql, sql)
	items := []OutPut{}
	err = db.ExecToArray(&items, sql)
	assert.NoError(t, err)
	assert.Equal(t, 100, len(items))

}
func BenchmarkTestMssqlSelectJoinExpr(t *testing.B) {

	//go test -bench='^BenchmarkTestMssqlSelectJoinExpr$' -run=^$ -benchmem -count=10

	mssqlDns := "sqlserver://sa:123456@localhost?database=a003"
	db, err := vdb.Open("sqlserver", mssqlDns)
	// dialect := vdb.DialectFactory.Create(db.GetDriverName())
	assert.NoError(t, err)
	err = db.Ping()
	assert.NoError(t, err)
	for i := 0; i < t.N; i++ {
		qr := vdb.Qr().From("user u").LeftJoin("department d", "u.deptId= d.id")
		qr.LeftJoin("position p", "u.positionId= p.id")
		qr.Select(
			"u.username",
			"d.name DepartmentName",
			"p.name PositionName",
			"u.email",
			"u.birthday UserBirthday",
		).OrderBy("u.id").OffsetLimit(0, 10000)

		sql, _ := qr.BuildSQL(db.TenantDB)
		// expectSql := "SELECT [u].[username] AS [username], [d].[name] AS [DepartmentName], [p].[name] AS [PositionName], [u].[email] AS [email], [u].[birthday] AS [UserBirthday] FROM [users] AS [u] LEFT JOIN [departments] AS [d] ON [u].[dept_id] = [d].[id] LEFT JOIN [positions] AS [p] ON [u].[position_id] = [p].[id]"
		// assert.Equal(t, expectSql, sql)
		type OutPut struct { //<-- vi du day la dau ra
			Username       *string
			DepartmentName *string
			PositionName   *string
			Email          *string
			UserBirthday   *time.Time
		}
		items := []OutPut{}
		err := db.ExecToArray(&items, sql)
		assert.NoError(t, err)
		assert.Equal(t, 10000, len(items))
	}
}
func TestMssqlSelectJoinExpr(t *testing.T) {
	mssqlDns := "sqlserver://sa:123456@localhost?database=a003"
	db, err := vdb.Open("sqlserver", mssqlDns)
	// dialect := vdb.DialectFactory.Create(db.GetDriverName())
	assert.NoError(t, err)
	err = db.Ping()
	assert.NoError(t, err)

	type OutPut struct { //<-- vi du day la dau ra
		Username       *string
		DepartmentName *string
		PositionName   *string
		Email          *string
		UserBirthday   *time.Time
	}
	qr := vdb.Qr().From("user u").LeftJoin("department d", "u.deptId= d.id")
	qr.LeftJoin("position p", "u.positionId= p.id")
	qr.Select(
		"u.username",
		"d.name DepartmentName",
		"p.name PositionName",
		"u.email",
		"u.birthday UserBirthday",
	)

	sql, _ := qr.BuildSQL(db.TenantDB)
	expectSql := "SELECT [u].[username] AS [username], [d].[name] AS [DepartmentName], [p].[name] AS [PositionName], [u].[email] AS [email], [u].[birthday] AS [UserBirthday] FROM [users] AS [u] LEFT JOIN [departments] AS [d] ON [u].[dept_id] = [d].[id] LEFT JOIN [positions] AS [p] ON [u].[position_id] = [p].[id]"
	assert.Equal(t, expectSql, sql)
	items := []OutPut{}

	err = db.ExecToArray(items, sql)
	assert.NoError(t, err)
	assert.Equal(t, 26004, len(items))

}
func TestMssqlQueryBuilder(t *testing.T) {
	mssqlDns := "sqlserver://sa:123456@localhost?database=a003"
	db, err := vdb.Open("sqlserver", mssqlDns)
	// dialect := vdb.DialectFactory.Create(db.GetDriverName())
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
		qr := vdb.Qr().From("user").Select(
			"concat(firstName, ?,lastName) as FullName",
			"Code",
			"salary*? as Salary", vdb.Lit(" "), 0.12)
		qr.Where("id =? or id =?", 10, 100)

		sql, args := qr.BuildSQL(db.TenantDB)
		expectSql := "SELECT CONCAT([T1].[first_name], @p1, [T1].[last_name]), [T1].[code] AS [Code], [T1].[salary] * @p2 AS [Salary] FROM [users] AS [T1] WHERE [T1].[id] = @p3 OR [T1].[id] = @p4"
		assert.Equal(t, expectSql, sql)
		assert.Equal(t, []interface{}{" ", 0.12, 10, 100}, args)
		items := []OutPut{}

		err := db.ExecToArray(items, sql)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(items))
	}

}
func BenchmarkTestMssqlQueryBuilder(t *testing.B) {
	mssqlDns := "sqlserver://sa:123456@localhost?database=a003"
	db, err := vdb.Open("sqlserver", mssqlDns)
	// dialect := vdb.DialectFactory.Create(db.GetDriverName())
	assert.NoError(t, err)
	err = db.Ping()
	assert.NoError(t, err)
	vdb.NewExprCompiler(db.TenantDB)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {

		qr := vdb.Qr().From("user").Select(
			"concat(firstName, ?,lastName) as FullName",
			"Code",
			"salary*? as Salary", vdb.Lit(" "), 0.12)
		qr.Where("id =? or id =?", 10, 100)
		sql, args := qr.BuildSQL(db.TenantDB)
		expectSql := "SELECT CONCAT([T1].[first_name], @p1, [T1].[last_name]), [T1].[code] AS [Code], [T1].[salary] * @p2 AS [Salary] FROM [users] AS [T1] WHERE [T1].[id] = @p3 OR [T1].[id] = @p4"
		assert.Equal(t, expectSql, sql)
		assert.Equal(t, []interface{}{" ", 0.12, 10, 100}, args)
	}
}

// func TestExecToArray(t *testing.T) {
// 	mssqlDns := "sqlserver://sa:123456@localhost?database=a003"
// 	db, err := vdb.Open("sqlserver", mssqlDns)
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
// 	db, err := vdb.Open("sqlserver", mssqlDns)
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
// 	db, err := vdb.Open("sqlserver", mssqlDns)
// 	assert.NoError(t, err)
// 	err = db.Ping()
// 	assert.NoError(t, err)
// 	m, err := vdb.NewMigrator(db)
// 	assert.NoError(t, err)
// 	err = m.DoMigrates()
// 	assert.NoError(t, err)
// 	items, err := vdb.SelectAll[models.User](db)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 24004, len(items))

// 	// err = db.Select(&users, "SELECT * FROM users")
// 	// assert.NoError(t, err)
// 	// assert.Equal(t, 1, len(users))
// }

// func BenchmarkTestMssqlSelect(t *testing.B) {
// 	mssqlDns := "sqlserver://sa:123456@localhost?database=a003"
// 	db, err := vdb.Open("sqlserver", mssqlDns)
// 	assert.NoError(t, err)
// 	err = db.Ping()
// 	assert.NoError(t, err)
// 	m, err := vdb.NewMigrator(db)
// 	assert.NoError(t, err)
// 	err = m.DoMigrates()
// 	assert.NoError(t, err)
// 	for i := 0; i < t.N; i++ {
// 		items, err := vdb.SelectAll[models.User](db)
// 		assert.NoError(t, err)
// 		assert.Equal(t, 24004, len(items))
// 	}

// }
// func TestMssqlUserInsertBatch(t *testing.T) {
// 	msssqlDns := "sqlserver://sa:123456@localhost?database=a003"
// 	db, err := vdb.Open("sqlserver", msssqlDns)
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
// 	// vdb.InsertBatch(db, users)
// 	err = db.InsertBatch(users)
// 	assert.NoError(t, err)
// }
// func BenchmarkTestMssqlUserInsertBatch(t *testing.B) {
// 	msssqlDns := "sqlserver://sa:123456@localhost?database=a003"
// 	db, err := vdb.Open("sqlserver", msssqlDns)
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
// 		// vdb.InsertBatch(db, users)
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
// 	db, err := vdb.Open("sqlserver", msssqlDns)
// 	assert.NoError(t, err)
// 	err = db.Ping()

// 	assert.NoError(t, err)
// 	m, err := vdb.NewMigrator(db)
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
// 	db, err := vdb.Open("sqlserver", msssqlDns)
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
// 	db, err := vdb.Open("sqlserver", msssqlDns)
// 	assert.NoError(t, err)
// 	err = db.Ping()

// 	assert.NoError(t, err)
// 	m, err := vdb.NewMigrator(db)
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
// 	db, err := vdb.Open("sqlserver", msssqlDns)
// 	assert.NoError(t, err)
// 	err = db.Ping()

// 	assert.NoError(t, err)
// 	m, err := vdb.NewMigrator(db)
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
// 	db, err := vdb.Open("sqlserver", msssqlDns)
// 	assert.NoError(t, err)
// 	err = db.Ping()

// 	assert.NoError(t, err)
// 	m, err := vdb.NewMigrator(db)
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
