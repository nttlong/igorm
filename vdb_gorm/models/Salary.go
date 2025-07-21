package models

// Salary struct GORM
type Salary struct {
	ID        int     `gorm:"primaryKey;autoIncrement"`            // vdb: db:"pk;auto"
	UserID    int     `gorm:"index:idx_salary_user"`               // vdb: db:"idx:idx_salary_user"
	Month     string  `gorm:"type:char(7);index:idx_salary_month"` // vdb: db:"type:char(7);idx:idx_salary_month"
	Base      float64 `gorm:"type:decimal(15,2)"`                  // vdb: db:"type:decimal(15,2)"
	Bonus     float64 `gorm:"type:decimal(15,2)"`                  // vdb: db:"type:decimal(15,2)"
	Deduction float64 `gorm:"type:decimal(15,2)"`                  // vdb: db:"type:decimal(15,2)"
	BaseModel         // Nhúng BaseModel
	// Quan hệ với User (cho GORM)
	User User `gorm:"foreignKey:UserID"`
}
