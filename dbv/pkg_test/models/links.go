package models

import (
	"dbv"
)

func init() {
	dbv.ModelRegistry.Add(
		&Contract{},
		&User{},

		&Department{},
		&Position{},
		&Contract{},
		&User{},
	)
}
