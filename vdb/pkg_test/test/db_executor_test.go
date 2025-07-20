package test

import (
	"strconv"
	"testing"
	"time"
	"vdb"

	"vdb/pkg_test/models"
	_ "vdb/pkg_test/models" //<-- This is important to load models

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var db *vdb.TenantDB //<-- khai báo database quản lý tenant
/*
hàm này sẽ được gọi từ các test case để khởi tạo database quản lý tenant
@param driver: driver của database, vdb hiện tại hỗ trợ sqlserver mysql và postgres
*/
func initDb(driver string, conn string) {
	var err error
	db, err = vdb.Open(driver, conn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

var testDb *vdb.TenantDB //<-- khai báo database quản lý tenant
func TestCreateTenantDb(t *testing.T) {
	var err error
	vdb.SetManagerDb("mysql", "tenant_manager") //<--- Cài đặt database quản lý tenannt
	// 	// Data base quản lý tenant phai co trước, đặc điểm của nó là kg migrate các model dựng sẵn,
	// 	// Nó chỉ tập trung vào việc quản lý tenant, không có migrate các model dựng sẵn.
	// 	// Việc chỉ định database quản lý tenant , bằng cách gọi hàm vdb.SetManagerDb("mysql", "tenantManager"), là rất quan trọng
	// 	// Nó giúp vdb biết database quản lý tenant là database nào để thực hiện các thao tác liên quan đến tenant.
	initDb("mysql", "root:123456@tcp(127.0.0.1:3306)/tenant_manager?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=True")
	testDb, err = db.CreateDB("test001") //<--- Tạo database tenant tên là test001
	assert.NoError(t, err)
	assert.Equal(t, "test001", testDb.GetDBName())

}
func TestCreateUser(t *testing.T) {
	TestCreateTenantDb(t) //<--- chạy test trước khi test này
	user := &models.User{
		UserId:       vdb.Ptr(uuid.NewString()),
		Email:        "test@test.com",
		Phone:        "0987654321",
		Username:     vdb.Ptr("test001"), //<-- hàm Ptr() được dùng để truyền tham số thành pointer
		HashPassword: vdb.Ptr("123456"),
		BaseModel: models.BaseModel{
			Description: vdb.Ptr("test"),
			CreatedAt:   vdb.Ptr(time.Now()),
		},
	}
	err := testDb.Insert(user)
	if err != nil {
		if vdbErr, ok := err.(*vdb.DialectError); ok {
			/*
				 vdb se có 1 số lội sau khi thực hiện thao tác insert,update hoặc deltete trên database
				 	DIALECT_DB_ERROR_TYPE_UNKNOWN DIALECT_DB_ERROR_TYPE = iota //<-- không xác định được lỗi gì
					DIALECT_DB_ERROR_TYPE_DUPLICATE //<-- duplicate record
					DIALECT_DB_ERROR_TYPE_REFERENCES // ✅ <-- vi phạm ràng buộc khi thực hiện thao tác insert,update hoặc delete
					DIALECT_DB_ERROR_TYPE_REQUIRED // <-- thiếu các trường cần thiết khi thực hiện thao tác insert,update hoặc delete
					DIALECT_DB_ERROR_TYPE_LIMIT_SIZE //<-- vượt qua kích thước của các trường khi thực hiện thao tác insert,update hoặc delete
			*/
			if vdbErr.ErrorType == vdb.DIALECT_DB_ERROR_TYPE_DUPLICATE { //<-- nếu có lỗi duplicate thì sẽ báo lỗi
				assert.Equal(t, []string{"Email"}, vdbErr.Fields) //<-- nếu có lỗi duplicate thì sẽ báo lỗi cu thể tren Feild nao của struct
				assert.Equal(t, []string{"email"}, vdbErr.DbCols) // <-- và cu cột nao của database
				assert.Equal(t, "users", vdbErr.Table)            //<-- và cụ thể tên của các bạng có liên quan đến lỗi duplicate

			}
		}
	} else {
		assert.NoError(t, err)
	}

	assert.Equal(t, "test001", *user.Username)
	assert.Equal(t, "test@test.com", user.Email)
	assert.Equal(t, "0987654321", user.Phone)

}
func TestGetUser(t *testing.T) {
	TestCreateTenantDb(t) //<--- chạy test trước khi test này
	user := &models.User{}
	err := testDb.First(user, "id = ?", 1)
	assert.NoError(t, err)
	assert.Equal(t, "test001", *user.Username)
	assert.Equal(t, "test@test.com", user.Email)
	assert.Equal(t, "0987654321", user.Phone)
}
func TestGetUpdateUser(t *testing.T) {
	TestCreateTenantDb(t) //<--- chạy test trước khi test này
	result := testDb.Model(&models.User{}).Where("id = ?", 1).Update("Username", "test002")
	assert.NoError(t, result.Error)
	assert.Equal(t, int64(1), result.RowsAffected)

}
func TestGetUpdateUserByMap(t *testing.T) {
	TestCreateTenantDb(t) //<--- chạy test trước khi test này
	result := testDb.Model(&models.User{}).Where("id = ?", 1).Update(
		map[string]interface{}{
			"Username": "test003",
			"Email":    "william.henry.harrison@example-pet-store.com",
		},
	)
	assert.NoError(t, result.Error)
	assert.Equal(t, int64(1), result.RowsAffected)

}
func TestUpdateUserByCallDbFunc(t *testing.T) {
	TestCreateTenantDb(t) //<--- chạy test trước khi test này
	TestCreateTenantDb(t) //<--- chạy test trước khi test này
	/*
			Sometime update set the value of a field with a function, we can use the DbFunCall() function to pass the function as a parameter to the update function.
		 use DbFunCall(expr string, args ...interface{})
	*/
	result := testDb.Model(&models.User{}).Where("id = ?", 1).Update(
		"Username", vdb.DbFunCall("CONCAT(UPPER(Username),?)", "test003")) //<-- hàm CONCAT() được dùng để tạo ra một chuỗi mới từ các giá trị truyền vào
	assert.NoError(t, result.Error)
	assert.Equal(t, int64(1), result.RowsAffected)

}
func TestUpdateUserByMapAndCallDbFunc(t *testing.T) {
	TestCreateTenantDb(t) //<--- chạy test trước khi test này
	//testDb.LikeValue("*.com") se chuyen thanh '%.com' neu chay tren mysql
	//vdb lay chuan sqlserver cho tat ca cac ham va toan tu sau do bien dich ra theo dialect
	//vdb sẽ tự động sửa tên field đúng với tên fiel trong database , không phân biệt chữ hoa chữ thường của tên field
	result := testDb.Model(&models.User{}).Where("email like ?", testDb.LikeValue("*.edu")).Update(
		map[string]interface{}{
			"Username":    vdb.DbFunCall("lower(Username)"),
			"Email":       vdb.DbFunCall("CONCAT(UPPER(Email),?)", ".com"),
			"phone":       vdb.DbFunCall("CONCAT(LEFT(Phone,3),?)", "-123456"),
			"description": "Hệ thống sẽ tự động sửa tên field đúng với tên fiel trong database , không phân biệt chữ hoa chữ thường của tên field",
		},
	)
	assert.NoError(t, result.Error)
	assert.Equal(t, int64(1), result.RowsAffected)
}
func TestInsertPositionAndDepartmentOnce(t *testing.T) {
	TestCreateTenantDb(t) //<--- chạy test trước khi test này
	i := 4
	strIndex := strconv.Itoa(i)
	position := &models.Position{

		Name:  "CEO" + strIndex,
		Code:  "CEO0" + strIndex,
		Title: "Chief Executive Officer " + strIndex,
		Level: 1,
		BaseModel: models.BaseModel{
			Description: vdb.Ptr("test"),
			CreatedAt:   vdb.Ptr(time.Now()),
		},
	}
	dept := &models.Department{
		Name: "CEO" + strIndex,
		Code: "CEO0" + strIndex,
		BaseModel: models.BaseModel{
			Description: vdb.Ptr("test"),
			CreatedAt:   vdb.Ptr(time.Now()),
		},
	}
	user := &models.User{
		UserId:       vdb.Ptr(uuid.NewString()),
		Email:        "test@test.com" + strIndex,
		Phone:        "0987654321",
		Username:     vdb.Ptr("test001" + strIndex), //<-- hàm Ptr() được dùng để truyền tham số thành pointer
		HashPassword: vdb.Ptr("123456" + strIndex),
		BaseModel: models.BaseModel{
			Description: vdb.Ptr("test"),
			CreatedAt:   vdb.Ptr(time.Now()),
		},
	}

	tx, err := testDb.Begin()
	if err != nil {
		t.Error(err)
	}
	err = testDb.Insert(position, dept, user) //<-- them 3 bang vao database trong 1 transaction
	if err != nil {
		if vdbErr, ok := err.(*vdb.DialectError); ok { //<-- neu detect duoc loi
			if vdbErr.ErrorType == vdb.DIALECT_DB_ERROR_TYPE_DUPLICATE { //<-- neu loi duoc detect la duplicate
				assert.Equal(t, "position", vdbErr.Table)
				assert.Equal(t, "models.Position", vdbErr.StructName) //<--       // kiem tra bang nao gay ra loi duplicate
				assert.Equal(t, []string{"Name"}, vdbErr.Fields)      //<-- kiem tra field nao trong struct gay co loi duplicate
				assert.Equal(t, []string{"name"}, vdbErr.DbCols)      //<-- kiem tra cot gay co loi duplicate
			}
		} else {
		}
		err = tx.Rollback()
		if err != nil {
			t.Error(err)
		}
	} else {
		err = tx.Commit()
		if err != nil {
			t.Error(err)
		} else {
			tx, err := testDb.Begin()
			if err != nil {
				t.Error(err)
			}
			emp := &models.Employee{

				PositionID:   position.ID, //<-- Id là số auto increment của position
				DepartmentID: dept.ID,     //<-- Id là số auto increment của dept
				UserID:       user.ID,     //<-- Id là số auto increment của User
				FirstName:    "John",
				LastName:     "Doe",
				BaseModel: models.BaseModel{
					Description: vdb.Ptr("test"),
					CreatedAt:   vdb.Ptr(time.Now()),
				},
			}
			err = tx.Insert(emp)
			if err != nil {
				err = tx.Rollback()
				if err != nil {
					t.Error(err)
				}
			} else {
				err = tx.Commit()
				if err != nil {
					t.Error(err)
				}
			}

		}
	}

}
func TestSelectAllepmloyeeAndUser(t *testing.T) {
	TestCreateTenantDb(t) //<--- chạy test trước khi test này
	type QueryResult struct {
		FullName     *string
		PositionID   *int64
		DepartmentID *int64
		Email        *string
		Phone        *string
	}
	items := []QueryResult{}

	qr := testDb.From((&models.Employee{}).As("e")).LeftJoin( //<-- inner join, left join ,right join and full join just change function
		(&models.User{}).As("u"), "e.id = u.userId",
	).Select(
		"concat(e.FirstName,' ', e.LastName) as fullName",
		"e.positionId",
		"e.departmentId",
		"u.email",
		"u.phone",
	)
	sql, _ := qr.BuildSql()
	t.Log(sql)
	//<-- count tong so user
	err := qr.ToArray(&items)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(items))
	assert.Equal(t, "John Doe", *items[0].FullName)

}
func TestSelfJoin(t *testing.T) {
	TestCreateTenantDb(t) //<--- chạy test trước khi test này
	type QueryResult struct {
		Name      *string
		ChildName *string
	}
	items := []QueryResult{}

	qr := testDb.From(
		(&models.Department{}).As("d"),
	).LeftJoin(
		(&models.Department{}).As("c"), "d.id = c.parentId",
	).Select(
		"d.name",              //<-- autho change to pascal case, because Go can not fill value to non-pascal case field
		"c.name as childName", //<-- autho change to pascal case, because Go can not fill value to non-pascal case field
	)

	err := qr.ToArray(&items)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(items))
	assert.Equal(t, "CEO01", *items[0].Name)
	assert.Equal(t, "CEO02", *items[0].ChildName)

}
