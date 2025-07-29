package models

import "vdb"

func init() {
	vdb.ModelRegistry.Add(
		&Contract{},
		&User{},

		&Department{},
		&Position{},
		&Contract{},
		&User{},
	)
}
