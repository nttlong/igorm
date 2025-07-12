package modelstest

import (
	"eorm"
	"time"
)

type User struct {
	eorm.Model `db:"table:users"`

	ID        int       `db:"pk" auto:"true"`                 // primary key, auto increment
	Name      string    `db:"column:name" type:"string(100)"` // mapped column name
	Email     string    `db:"unique" type:"string(255)"`      // unique constraint
	Profile   string    `db:"column:profile" type:"json"`     // abstract JSON/document field
	CreatedAt time.Time `db:"default(now)" type:"datetime"`   // default timestamp
}

func init() {
	eorm.ModelRegistry.Add(&User{})
}
