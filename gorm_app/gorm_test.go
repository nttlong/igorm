package gormapp

import (
	"fmt"
	"testing"
	"time"

	// _ "github.com/microsoft/go-mssqldb"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// Order đại diện cho bảng orders
type Order struct {
	OrderId   int64      `gorm:"column:order_id;primaryKey;type:bigint;autoIncrement"` // Khóa chính, tự tăng
	Version   int        `gorm:"column:version;primaryKey;type:int;not null"`          // Khóa chính thứ hai
	CreatedAt time.Time  `gorm:"column:created_at;type:datetime2(7);not null"`         // Thời gian tạo
	CreatedBy string     `gorm:"column:created_by;type:nvarchar(100);not null"`        // Người tạo
	UpdatedAt *time.Time `gorm:"column:updated_at;type:datetime2(7)"`                  // Thời gian cập nhật (NULL được)
	UpdatedBy *string    `gorm:"column:updated_by;type:nvarchar(100)"`                 // Người cập nhật (NULL được)
	Note      string     `gorm:"column:note;type:nvarchar(200);not null"`              // Ghi chú
	DeletedAt *time.Time `gorm:"column:updated_at;type:datetime2(7)"`
	// Tùy chỉnh tên bảng
	gorm.Model
}

// TableName định nghĩa tên bảng trong cơ sở dữ liệu
func (Order) TableName() string {
	return "orders"
}

func TestOrder(t *testing.T) {
	dsn := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := gorm.Open(sqlserver.Open(dsn))
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	db.AutoMigrate(&Order{})
	// Tự động tạo bảng (nếu chưa tồn tại)
	// db.AutoMigrate(&Order{})

	// Thực hiện truy vấn SELECT với OFFSET và FETCH
	var orders []Order

	result := db.
		Select([]string{"order_id AS OrderId",
							"created_at AS CreatedAt",
							"created_by AS CreatedBy",
							"note AS Note",

							"updated_at AS UpdatedAt",
							"updated_by AS UpdatedBy",
							"version AS Version"}).
		Order("order_id ASC, version ASC"). // Sửa ở đây: dùng chuỗi trực tiếp
		Limit(10000).
		Offset(0).
		Find(&orders)

	if result.Error != nil {
		panic("failed to query: " + result.Error.Error())
	}

	if result.Error != nil {
		panic("failed to query: " + result.Error.Error())
	}

	// In kết quả
	fmt.Println("Số lượng bản ghi:", len(orders))
}
func BenchmarkGormWithWhereAndLimitAndOffset(b *testing.B) {
	dsn := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := gorm.Open(sqlserver.Open(dsn))
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	avg := int64(0)
	sqlSmt := db.
		Select([]string{"order_id AS OrderId",
			"created_at AS CreatedAt",
			"created_by AS CreatedBy",
			"note AS Note",

			"updated_at AS UpdatedAt",
			"updated_by AS UpdatedBy",
			"version AS Version"}).
		Where("order_id > ?", 1000).
		Order("order_id ASC, version ASC"). // Sửa ở đây: dùng chuỗi trực tiếp
		Limit(10000).
		Offset(0)
	b.ResetTimer()
	fmt.Println("\n---------------------------------------------------")
	for i := 0; i < b.N; i++ {

		var orders []Order

		start := time.Now()
		result := sqlSmt.Find(&orders)
		n := time.Since(start).Nanoseconds()
		fmt.Printf("Tong thoi gian thuc hien lenh find cua GORM : %d ns\n", n)
		fmt.Println("---------------------------------------------------")
		avg += n

		//fmt.Println("Tong thoi gian thuc hien lenh find cua GORM:", time.Since(start).Nanoseconds())

		if result.Error != nil {
			panic("failed to query: " + result.Error.Error())
		}

		if result.Error != nil {
			panic("failed to query: " + result.Error.Error())
		}
	}
	//b.Log("Avg time: ", avg/int64(b.N))
	// In kết quả
	//fmt.Println("Số lượng bản ghi:", len(orders))
}
