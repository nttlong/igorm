package models

// Attendance struct GORM
type Attendance struct {
	ID int `gorm:"primaryKey;autoIncrement"` // vdb: db:"pk;auto"
	// UserID là khóa ngoại, GORM sẽ tự động hiểu nếu có mối quan hệ BelongsTo/HasMany
	// GORM: foreignKey và references để rõ ràng hơn, hoặc GORM tự suy luận nếu tên khớp
	UserID    int    `gorm:"index:idx_att_user"`           // vdb: db:"idx:idx_att_user"
	Date      string `gorm:"type:date;index:idx_att_date"` // vdb: db:"type:date;idx:idx_att_date"
	CheckIn   string `gorm:"type:time"`                    // vdb: db:"type:time"
	CheckOut  string `gorm:"type:time"`                    // vdb: db:"type:time"
	BaseModel        // Nhúng BaseModel
	// Quan hệ với User (cho GORM)
	User User `gorm:"foreignKey:UserID"`
}
