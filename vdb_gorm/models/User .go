package models

// User struct GORM
type User struct {
	ID           int     `gorm:"primaryKey;autoIncrement"`          // vdb: db:"pk;auto"
	UserId       *string `gorm:"type:varchar(36);unique"`           // vdb: db:"size:36;unique"
	Email        string  `gorm:"type:varchar(150);unique;not null"` // vdb: db:"uk:uq_email;size:150"
	Phone        string  `gorm:"type:varchar(20)"`                  // vdb: db:"size:20"
	Username     *string `gorm:"type:varchar(150);unique;not null"` // vdb: db:"size:50;unique"
	HashPassword *string `gorm:"type:varchar(100)"`                 // vdb: db:"size:100"
	BaseModel            // Nhúng BaseModel
}
