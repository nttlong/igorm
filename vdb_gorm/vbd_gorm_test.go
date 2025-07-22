package vdbgorm_test

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	models "vdb_gorm/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert" // Hoặc driver bạn đang dùng, ví dụ: postgres.Open("...")

	//"github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gormDb *gorm.DB

func createTenantDbGorm() error {
	var err error

	// 1. Chuỗi kết nối đến database tenant cụ thể ("test001")
	// Lưu ý: GORM không tự động TẠO database vật lý.
	// Bạn cần đảm bảo database "test001" đã được tạo trên MySQL server của bạn TRƯỚC KHI chạy code này.
	// Ví dụ: Bạn có thể chạy lệnh SQL "CREATE DATABASE test001;" trong MySQL client.
	//tenantDbName := "gorm_test004" // Tên database của tenant
	tenantDbName := "vdb_test005"
	dsnTenant := fmt.Sprintf("root:123456@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=True", tenantDbName)

	gormDb, err = gorm.Open(mysql.Open(dsnTenant), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info), // Thay đổi logger.Info thành logger.Silent để tắt log queries
		// SỬA ĐOẠN NÀY ĐỂ TẮT LOG QUERIES:
		Logger: logger.Default.LogMode(logger.Silent), // Thay đổi logger.Info thành logger.Silent
	})
	if err != nil {
		// Nếu database chưa tồn tại hoặc thông tin kết nối sai, bạn sẽ nhận được lỗi ở đây.
		return fmt.Errorf("failed to connect to tenant DB '%s': %w", tenantDbName, err)
	}
	log.Printf("Connected to tenant DB '%s' successfully.", tenantDbName)

	// 2. Quan trọng: AutoMigrate các model cho database của tenant này
	// Điều này sẽ tạo bảng `users` (và các bảng khác nếu bạn có model khác)
	// trong database `test001` nếu chúng chưa tồn tại, hoặc cập nhật schema.
	//err = models.MigrateAllModels(gormDb)
	if err != nil {
		return fmt.Errorf("failed to auto migrate models for tenant DB '%s': %w", tenantDbName, err)
	}
	log.Printf("AutoMigrate for tenant DB '%s' completed successfully.", tenantDbName)

	return nil
}

// Benchmark tương tự với GORM
func Benchmark_GORM_TestCreateUser(b *testing.B) {
	// Khởi tạo database cho GORM, chỉ chạy một lần
	if err := createTenantDbGorm(); err != nil {
		log.Fatalf("Failed to setup GORM DB: %v", err)
	}

	b.ResetTimer() // Đặt lại bộ đếm thời gian trước khi bắt đầu benchmark

	for i := 0; i < b.N; i++ {
		name := "gorm_test_user_1x" + uuid.NewString()
		email := name + "@gorm.com"
		phone := "0987654321"
		password := "gorm_123456"
		description := "GORM test user"
		now := time.Now()

		user := &models.User{
			UserId:       &[]string{uuid.NewString()}[0], // Lấy địa chỉ của phần tử đầu tiên trong slice để tạo *string
			Email:        email,
			Phone:        phone,
			Username:     &name,
			HashPassword: &password,
			BaseModel: models.BaseModel{
				Description: &description,
				CreatedAt:   &now, // Sử dụng hàm NowFunc của GORM để lấy thời gian hiện tại
			},
		}

		// Thao tác Insert chính được benchmark
		err := gormDb.Create(user).Error
		if err != nil {
			// GORM không có cơ chế lỗi chi tiết như vdb.DialectError
			// Bạn sẽ cần kiểm tra chuỗi lỗi hoặc mã lỗi SQL cụ thể
			// Ví dụ: kiểm tra lỗi trùng lặp UNIQUE constraint
			if gormDb.Error != nil && gormDb.Error.Error() == "UNIQUE constraint failed: users.email" {
				// Đây là một cách đơn giản để kiểm tra lỗi trùng lặp cho SQLite
				// Đối với các DB khác, lỗi có thể khác (ví dụ: "duplicate key value violates unique constraint")
				assert.Fail(b, "Duplicate entry error for email, expected to be unique within test run")
			} else {
				assert.NoError(b, err, "Unexpected error during GORM user creation")
			}
		} else {
			assert.NoError(b, err) // Đảm bảo không có lỗi
		}

		// Kiểm tra dữ liệu sau khi chèn
		assert.Equal(b, name, *user.Username)
		assert.Equal(b, email, user.Email)
		assert.Equal(b, phone, user.Phone)
	}
}
func Benchmark_TestGetUpdateUserByMap(t *testing.B) {
	// Khởi tạo database cho GORM, chỉ chạy một lần
	if err := createTenantDbGorm(); err != nil {
		log.Fatalf("Failed to setup GORM DB: %v", err)
	}

	t.ResetTimer() // Đặt lại bộ đếm thời gian trước khi bắt đầu benchmark
	for i := 0; i < t.N; i++ {
		name := "test" + uuid.NewString()
		result := gormDb.Model(&models.User{}).Where("id = ?", 1).Updates(
			map[string]interface{}{
				"Username": name,
				"Email":    "william.henry.harrison@example-pet-store.com",
			},
		)
		assert.NoError(t, result.Error)
		assert.Equal(t, int64(1), result.RowsAffected)
	}

}
func Benchmark_GORM_TestUpdateUserByCallDbFunc(b *testing.B) {
	// Khởi tạo database và đảm bảo có dữ liệu mẫu User với ID=1
	if err := createTenantDbGorm(); err != nil {
		log.Fatalf("Failed to setup GORM DB: %v", err)
	}

	b.ResetTimer() // Đặt lại bộ đếm thời gian trước khi bắt đầu benchmark

	for i := 0; i < b.N; i++ {
		// Sử dụng gorm.Expr để truyền biểu thức SQL
		// Tương đương với vdb.DbFunCall("LEFT(CONCAT(UPPER(Username),?),50)", "test003")
		result := gormDb.Model(&models.User{}).Where("id = ?", 1).Update(
			"Username", gorm.Expr("LEFT(CONCAT(?,UPPER(username)),50)", uuid.NewString()),
		)
		// Lưu ý: GORM sẽ tự động chuyển đổi tên trường 'Username' sang 'username'
		// hoặc tên cột đã định nghĩa trong tag gorm:`column:username`.
		// Tôi đã sử dụng 'username' (lowercase) trong gorm.Expr để khớp với MySQL convention.

		assert.NoError(b, result.Error)
		assert.Equal(b, int64(1), result.RowsAffected)
	}
}
func Benchmark_TestUpdateUserByMapAndCallDbFunc_GORM(b *testing.B) {
	if err := createTenantDbGorm(); err != nil {
		log.Fatalf("Failed to setup GORM DB: %v", err)
	} // giả định tạo DB test thành công

	for i := 0; i < b.N; i++ {
		likePattern := strings.ReplaceAll("*.edu", "*", "%") // LikeValue tương đương

		// MySQL/MSSQL dùng LEN hoặc LENGTH tùy dialect, ví dụ ở đây là MySQL:
		usernameExpr := gorm.Expr("CONCAT(LEFT(UPPER(username), CHAR_LENGTH(username) - 1), ?)", strconv.Itoa(i))

		//CONCAT(left(UPPER(Email),len(Email)-40), ?)", ".com"+uuid.NewString()
		emailExpr := gorm.Expr("CONCAT(LEFT(UPPER(email), CHAR_LENGTH(email) - 40), ?)", ".com"+uuid.NewString())
		phoneExpr := gorm.Expr("CONCAT(LEFT(phone, 3), ?)", "-123456")

		result := gormDb.Model(&models.User{}).
			Where("email NOT LIKE ?", likePattern).
			Updates(map[string]interface{}{
				"username":    usernameExpr,
				"email":       emailExpr,
				"phone":       phoneExpr,
				"description": "Hệ thống sẽ tự động sửa tên field đúng với tên fiel trong database , không phân biệt chữ hoa chữ thường của tên field",
			})

		assert.NoError(b, result.Error)
	}
}
func ptr[T any](v T) *T {
	return &v
}
func handleDuplicateError(t testing.TB, err error, table, structName string, fields, dbCols []string) error {
	// Tùy driver, bạn có thể match lỗi từ `err.Error()` hoặc `errors.As(err, *mysql.MySQLError)`
	if strings.Contains(err.Error(), "duplicate") {
		t.Logf("Detected duplicate on table %s, struct %s, field %v, column %v",
			table, structName, fields, dbCols)
		return nil // skip commit
	}
	return err
}

// go test -bench=^Benchmark_TestInsertPositionAndDepartmentOnce_GORM$ -benchmem -benchtime=5s -count=10 > gorm4.txt
var setupOnce sync.Once

func Benchmark_TestInsertPositionAndDepartmentOnce_GORM(b *testing.B) {

	setupOnce.Do(func() {
		if err := createTenantDbGorm(); err != nil {
			log.Fatalf("Failed to setup GORM DB: %v", err)
		}
	})

	b.ResetTimer() // <--- Reset lại bộ đếm thời gian trước benchmark

	for i := 0; i < b.N; i++ {

		strIndex := uuid.NewString()

		// Khởi tạo dữ liệu
		position := &models.Position{
			Name:  "CEO" + strIndex,
			Code:  "CEO0" + strIndex,
			Title: "Chief Executive Officer " + strIndex,
			Level: 1,
			BaseModel: models.BaseModel{
				Description: ptr("test"),
				CreatedAt:   ptr(time.Now()),
			},
		}

		dept := &models.Department{
			Name: "CEO" + strIndex,
			Code: "CEO0" + strIndex,
			BaseModel: models.BaseModel{
				Description: ptr("test"),
				CreatedAt:   ptr(time.Now()),
			},
		}

		user := &models.User{
			UserId:       ptr(uuid.NewString()),
			Email:        "test@test.com" + strIndex,
			Phone:        "0987654321",
			Username:     ptr("test001" + strIndex),
			HashPassword: ptr("123456" + strIndex),
			BaseModel: models.BaseModel{
				Description: ptr("test"),
				CreatedAt:   ptr(time.Now()),
			},
		}

		err := gormDb.Transaction(func(tx *gorm.DB) error {
			// Insert Position
			if err := tx.Create(position).Error; err != nil {
				return err
			}
			// Insert Department
			if err := tx.Create(dept).Error; err != nil {
				return err
			}
			// Insert User
			if err := tx.Create(user).Error; err != nil {
				return err
			}

			// Insert Employee (dùng ID từ các bản ghi trước)
			emp := &models.Employee{
				PositionID:   position.ID,
				DepartmentID: dept.ID,
				UserID:       user.ID,
				FirstName:    "John",
				LastName:     "Doe",
				BaseModel: models.BaseModel{
					Description: ptr("test"),
					CreatedAt:   ptr(time.Now()),
				},
			}
			if err := tx.Create(emp).Error; err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			b.Error(err)
		}

	}
}
func Benchmark_TestSelectEmployeeAndUser_GORM(b *testing.B) {
	//go test -bench=Benchmark_TestSelectEmployeeAndUser_GORM -run=^$ -benchmem -benchtime=5s -count=10 > gorm5.txt
	// Chạy setup một lần
	dsn := "root:123456@tcp(127.0.0.1:3306)/vdb_test005?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=True"

	gormDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info), // Thay đổi logger.Info thành logger.Silent để tắt log queries
		// SỬA ĐOẠN NÀY ĐỂ TẮT LOG QUERIES:
		Logger: logger.Default.LogMode(logger.Silent), // Thay đổi logger.Info thành logger.Silent
	})
	assert.NoError(b, err, "Failed to connect to GORM DB")
	type BaseModel struct {
		FullName   *string
		PositionID *int64
	}
	type QueryResult struct {
		BaseModel
		DepartmentID *int64
		Email        *string
		Phone        *string
	}

	for i := 0; i < b.N; i++ {
		var items []QueryResult

		qr := gormDb.
			Table("employees AS e").
			Select(`
				concat(e.first_name, ' ', e.last_name) AS full_name,
				e.position_id,
				e.department_id,
				u.email,
				u.phone`).
			Joins("LEFT JOIN users u ON e.user_id = u.id").
			Order("e.id").
			Offset(0). // Giả sử bạn muốn lấy từ bản ghi thứ 1000
			Limit(1000)
		// start := time.Now()
		err := qr.Scan(&items).Error
		// n := time.Since(start).Milliseconds()
		// fmt.Println("Time: ", n) // In ra thời gian thực hiện truy vấn 6ms
		assert.NoError(b, err)
		assert.Equal(b, 1000, len(items))
		assert.Equal(b, "John Doe", *items[0].FullName)
	}
}
