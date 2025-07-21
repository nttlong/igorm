package models

// Department struct GORM
type Department struct {
	ID        int    `gorm:"primaryKey;autoIncrement"` // vdb: db:"pk;auto"
	Name      string `gorm:"type:varchar(100);unique"` // vdb: db:"size:100;uk:uq_dept_name"
	Code      string `gorm:"type:varchar(120);unique"` // vdb: db:"size:20;uk:uq_dept_code"
	ParentID  *int   `gorm:"index"`                    // vdb: ParentID *int, AddForeignKey("ParentID", &Department{}, "ID")
	BaseModel        // Nhúng BaseModel
	// Quan hệ đệ quy (self-referencing relationship)
	// Parent Department (HasOne hoặc BelongsTo)
	// Children Departments (HasMany)
	Parent   *Department  `gorm:"foreignKey:ParentID"` // Quan hệ với chính nó cho Parent
	Children []Department `gorm:"foreignKey:ParentID"` // Quan hệ với chính nó cho Children
}

// GORM sẽ tự động tạo khóa ngoại khi AutoMigrate nếu các mối quan hệ được định nghĩa đúng.
// Bạn không cần init() function riêng cho khóa ngoại đệ quy như vdb.
