package test

import (
	"strconv"
	"sync"
	"testing"
	"time"
	"vdb"
	"vdb/pkg_test/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createTenantDb() error {
	var err error
	vdb.SetManagerDb("mysql", "tenant_manager") //<--- Cài đặt database quản lý tenannt
	// 	// Data base quản lý tenant phai co trước, đặc điểm của nó là kg migrate các model dựng sẵn,
	// 	// Nó chỉ tập trung vào việc quản lý tenant, không có migrate các model dựng sẵn.
	// 	// Việc chỉ định database quản lý tenant , bằng cách gọi hàm vdb.SetManagerDb("mysql", "tenantManager"), là rất quan trọng
	// 	// Nó giúp vdb biết database quản lý tenant là database nào để thực hiện các thao tác liên quan đến tenant.
	initDb("mysql", "root:123456@tcp(127.0.0.1:3306)/tenant_manager?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=True")
	testDb, err = db.CreateDB("vdb_test005") //<--- Tạo database tenant tên là test004 dong thoi migrate các model dựng sẵn
	return err

}
func Benchmark_TestCreateUser(t *testing.B) {
	assert.NoError(t, createTenantDb())
	name := "test" + uuid.NewString()
	user := &models.User{
		UserId:       vdb.Ptr(uuid.NewString()),
		Email:        name + "@test.com",
		Phone:        "0987654321",
		Username:     vdb.Ptr(name), //<-- hàm Ptr() được dùng để truyền tham số thành pointer
		HashPassword: vdb.Ptr("123456"),
		BaseModel: models.BaseModel{
			Description: vdb.Ptr("test"),
			CreatedAt:   vdb.Ptr(time.Now()),
		},
	}
	testDb.Insert(user)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		name := "test" + uuid.NewString()
		user := &models.User{
			UserId:       vdb.Ptr(uuid.NewString()),
			Email:        name + "@test.com",
			Phone:        "0987654321",
			Username:     vdb.Ptr(name), //<-- hàm Ptr() được dùng để truyền tham số thành pointer
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

		assert.Equal(t, name, *user.Username)
		assert.Equal(t, name+"@test.com", user.Email)
		assert.Equal(t, "0987654321", user.Phone)
	}

}
func Benchmark_TestGetUpdateUserByMap(t *testing.B) {
	assert.NoError(t, createTenantDb()) //<--- chạy test trước khi test này
	for i := 0; i < t.N; i++ {
		name := "test" + uuid.NewString()
		result := testDb.Model(&models.User{}).Where("id = ?", 1).Update(
			map[string]interface{}{
				"Username": name,
				"Email":    "william.henry.harrison@example-pet-store.com",
			},
		)
		assert.NoError(t, result.Error)
		assert.Equal(t, int64(1), result.RowsAffected)
	}

}
func Benchmark_TestUpdateUserByCallDbFunc(t *testing.B) {
	assert.NoError(t, createTenantDb()) //<--- chạy test trước khi test này
	for i := 0; i < t.N; i++ {
		/*
				Sometime update set the value of a field with a function, we can use the DbFunCall() function to pass the function as a parameter to the update function.
			 use DbFunCall(expr string, args ...interface{})
		*/
		result := testDb.Model(&models.User{}).Where("id = ?", 1).Update(
			"Username", vdb.Expr("LEFT(CONCAT(?,UPPER(Username)),50)", uuid.NewString())) //<-- hàm CONCAT() được dùng để tạo ra một chuỗi mới từ các giá trị truyền vào
		assert.NoError(t, result.Error)
		assert.Equal(t, int64(1), result.RowsAffected)
	}

}
func Benchmark_TestUpdateUserByMapAndCallDbFunc(t *testing.B) {
	assert.NoError(t, createTenantDb()) //<--- chạy test trước khi test này
	//testDb.LikeValue("*.com") se chuyen thanh '%.com' neu chay tren mysql
	//vdb lay chuan sqlserver cho tat ca cac ham va toan tu sau do bien dich ra theo dialect
	//vdb sẽ tự động sửa tên field đúng với tên fiel trong database , không phân biệt chữ hoa chữ thường của tên field
	for i := 0; i < t.N; i++ {
		result := testDb.Model(&models.User{}).Where("email not like ?", testDb.LikeValue("*.edu")).Update(
			map[string]interface{}{
				"Username":    vdb.Expr("CONCAT(left(UPPER(Username),len(Username)-1), ?)", strconv.Itoa(i)), //<-- hàm LEFT() được dùng để lấy một phần của chuỗi
				"Email":       vdb.Expr("CONCAT(left(UPPER(Email),len(Email)-40), ?)", ".com"+uuid.NewString()),
				"phone":       vdb.Expr("CONCAT(LEFT(Phone,3),?)", "-123456"),
				"description": "Hệ thống sẽ tự động sửa tên field đúng với tên fiel trong database , không phân biệt chữ hoa chữ thường của tên field",
			},
		)
		assert.NoError(t, result.Error)

	}
}

// go test -bench=^Benchmark_TestInsertPositionAndDepartmentOnce$ -benchmem -count=5 vdb/pkg_test/test
var setupOnce sync.Once

// go test -bench=Benchmark_TestInsertPositionAndDepartmentOnce -run=^$ -benchmem -benchtime=5s -count=10 > vdb5.txt
func Benchmark_TestInsertPositionAndDepartmentOnce(b *testing.B) {

	setupOnce.Do(func() {
		assert.NoError(b, createTenantDb())
	})

	for i := 0; i < b.N; i++ {
		strIndex := uuid.NewString()
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
			Username:     vdb.Ptr("test001" + strIndex),
			HashPassword: vdb.Ptr("123456" + strIndex),
			BaseModel: models.BaseModel{
				Description: vdb.Ptr("test"),
				CreatedAt:   vdb.Ptr(time.Now()),
			},
		}

		// 👉 Bắt đầu đo tại đây (sau phần chuẩn bị struct)
		b.StartTimer()

		tx, err := testDb.Begin()
		if err != nil {
			b.Error(err)
		}
		if err = tx.Insert(position, dept, user); err != nil {
			b.Error(err)
		}
		emp := &models.Employee{
			PositionID:   position.ID,
			DepartmentID: dept.ID,
			UserID:       user.ID,
			FirstName:    "John",
			LastName:     "Doe",
			BaseModel: models.BaseModel{
				Description: vdb.Ptr("test"),
				CreatedAt:   vdb.Ptr(time.Now()),
			},
		}
		if err = tx.Insert(emp); err != nil {
			b.Error(err)
		}
		if err = tx.Commit(); err != nil {
			b.Error(err)
		}

		b.StopTimer()
	}
}
