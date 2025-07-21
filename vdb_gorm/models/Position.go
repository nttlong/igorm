package models

// Position struct GORM
type Position struct {
	ID        int    `gorm:"primaryKey;autoIncrement"` // vdb: db:"pk;auto"
	Code      string `gorm:"type:varchar(100);unique"` // vdb: db:"size:100;uk:uq_pos_code"
	Name      string `gorm:"type:varchar(100);unique"` // vdb: db:"size:100;uk:uq_pos_name"
	Title     string `gorm:"type:varchar(100);unique"` // vdb: db:"size:100;uk:uq_pos_title"
	Level     int
	BaseModel // Nhúng BaseModel
}
