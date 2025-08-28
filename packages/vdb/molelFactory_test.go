package vdb

import (
	"testing"
	"time"
)

type UserData struct {
	Model[UserData]
	Id        string    `db:"default:uuid()"`
	CreateOn  time.Time "db:default:now()"
	Check     bool      `db:"default:true"`
	CreatedBy string    `db:"default:'admin'"`
	Maximun   float64   `db:"default:100"`
}

func TestNewModel(t *testing.T) {
	ModelRegistry.Add(&UserData{})
	user, err := NewFromModel[UserData]()
	if err != nil {
		t.Error(err)
	}
	t.Log(user)
}
func BenchmarkNewModel(b *testing.B) {
	ModelRegistry.Add(&UserData{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewFromModel[UserData]()
	}
}
