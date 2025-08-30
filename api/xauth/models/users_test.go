package models

import (
	"testing"
	"vdb"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	user, err := vdb.NewFromModel[User]()
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.True(t, user.Active)
}
func BenchmarkCreateUser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		user, err := vdb.NewFromModel[User]()
		if err != nil {
			b.Fatal(err)
		}
		if user == nil {
			b.Fatal("user is nil")
		}
	}

}
