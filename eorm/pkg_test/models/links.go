package models

import (
	"eorm"
)

func init() {
	eorm.ModelRegistry.Add(
		&Contract{},
		&User{},

		&Department{},
		&Position{},
		&Contract{},
		&User{},
	)
}
