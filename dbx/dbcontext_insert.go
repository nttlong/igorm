package dbx

import (
	"context"

	"github.com/google/uuid"
)

func Insert(db *DBXTenant, entity interface{}) error {
	return db.Insert(entity)
}
func InsertWithContext(ctx context.Context, db *DBXTenant, entity interface{}) error {
	return db.InsertWithContext(ctx, entity)
}
func NewUUID() string {
	return uuid.New().String()
}
