package models

// Employee struct GORM
type Employee struct {
	ID           int    `json:"id" gorm:"primaryKey;autoIncrement"`     // vdb: db:"pk;auto;"
	FirstName    string `json:"name" gorm:"type:varchar(50);index"`     // vdb: db:"size(50);idx"
	LastName     string `json:"lastName" gorm:"type:varchar(50);index"` // vdb: db:"size(50);idx"
	DepartmentID int    `json:"departmentId" gorm:"index"`              // vdb: db:"fk(Department.ID)"
	PositionID   int    `json:"positionId" gorm:"index"`                // vdb: db:"fk(Position.ID)"
	UserID       int    `json:"userId" gorm:"index"`                    // vdb: db:"fk(User.ID)"
	BaseModel           // Nhúng BaseModel

	// Định nghĩa mối quan hệ (Associations) trong GORM
	// GORM sẽ tạo khóa ngoại và cho phép preload dữ liệu.
	Department Department `gorm:"foreignKey:DepartmentID"`
	Position   Position   `gorm:"foreignKey:PositionID"`
	User       User       `gorm:"foreignKey:UserID"`
}

// Trong GORM, bạn không cần hàm init() để đăng ký model hoặc thêm khóa ngoại như vdb.
// GORM sẽ tự động phát hiện các model và tạo khóa ngoại khi bạn gọi AutoMigrate.
// Ví dụ: db.AutoMigrate(&Employee{}, &Department{}, &Position{}, &User{})
// Đảm bảo rằng các model liên quan đã được định nghĩa và có các trường ID tương ứng.
